// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package observiq contains OpAmp structures compatible with the observiq client
package observiq

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

var (
	// ErrUnsupportedURL is error returned when creating a client with an unsupported URL scheme
	ErrUnsupportedURL = errors.New("unsupported URL")
)

// Ensure interface is satisfied
var _ opamp.Client = (*Client)(nil)

// mainPackageName is the name for the top level packages for this collector
const mainPackageName = "observiq-otel-collector"

// Client represents a client that is connected to Iris via OpAmp
type Client struct {
	opampClient           client.OpAMPClient
	logger                *zap.Logger
	ident                 *identity
	configManager         opamp.ConfigManager
	collector             collector.Collector
	packagesStateProvider types.PackagesStateProvider

	currentConfig opamp.Config
}

// NewClientArgs arguments passed when creating a new client
type NewClientArgs struct {
	DefaultLogger *zap.Logger
	Config        opamp.Config
	Collector     collector.Collector

	ManagerConfigPath   string
	CollectorConfigPath string
	LoggerConfigPath    string
}

// NewClient creates a new OpAmp client
func NewClient(args *NewClientArgs) (opamp.Client, error) {
	clientLogger := args.DefaultLogger.Named("opamp")

	configManager := NewAgentConfigManager(args.DefaultLogger)

	packagesStateProvider := newPackagesStateProvider(clientLogger, "package_statuses.yaml")

	observiqClient := &Client{
		logger:                clientLogger,
		ident:                 newIdentity(clientLogger, args.Config),
		configManager:         configManager,
		collector:             args.Collector,
		currentConfig:         args.Config,
		packagesStateProvider: packagesStateProvider,
	}

	// Parse URL to determin scheme
	opampURL, err := url.Parse(args.Config.Endpoint)
	if err != nil {
		return nil, err
	}

	// Add managed configs
	if err := observiqClient.addManagedConfigs(args); err != nil {
		return nil, err
	}

	// Create collect client based on URL scheme
	switch opampURL.Scheme {
	case "ws", "wss":
		observiqClient.opampClient = client.NewWebSocket(clientLogger.Sugar())
	default:
		return nil, ErrUnsupportedURL
	}

	return observiqClient, nil
}

func (c *Client) addManagedConfigs(args *NewClientArgs) error {
	// Add configs to config manager
	managerManagedConfig, err := opamp.NewManagedConfig(args.ManagerConfigPath, managerReload(c, args.ManagerConfigPath))
	if err != nil {
		return fmt.Errorf("failed to create manager managed config: %w", err)
	}
	c.configManager.AddConfig(ManagerConfigName, managerManagedConfig)

	collectorManagedConfig, err := opamp.NewManagedConfig(args.CollectorConfigPath, collectorReload(c, args.CollectorConfigPath))
	if err != nil {
		return fmt.Errorf("failed to create collector managed config: %w", err)
	}
	c.configManager.AddConfig(CollectorConfigName, collectorManagedConfig)

	loggerManagedConfig, err := opamp.NewManagedConfig(args.LoggerConfigPath, loggerReload(c, args.LoggerConfigPath))
	if err != nil {
		return fmt.Errorf("failed to create logger managed config: %w", err)
	}
	c.configManager.AddConfig(LoggingConfigName, loggerManagedConfig)

	return nil
}

// Connect initiates a connection to the OpAmp server
func (c *Client) Connect(ctx context.Context) error {
	// Compose and set the agent description
	if err := c.opampClient.SetAgentDescription(c.ident.ToAgentDescription()); err != nil {
		c.logger.Error("Error while setting agent description", zap.Error(err))
		return err
	}

	tlsCfg, err := c.currentConfig.ToTLS()
	if err != nil {
		return fmt.Errorf("failed creating TLS config: %w", err)
	}

	settings := types.StartSettings{
		OpAMPServerURL: c.currentConfig.Endpoint,
		Header: http.Header{
			"Authorization":  []string{fmt.Sprintf("Secret-Key %s", c.currentConfig.GetSecretKey())},
			"User-Agent":     []string{fmt.Sprintf("observiq-otel-collector/%s", version.Version())},
			"OpAMP-Version":  []string{opamp.Version()},
			"Agent-ID":       []string{c.ident.agentID},
			"Agent-Version":  []string{version.Version()},
			"Agent-Hostname": []string{c.ident.hostname},
		},
		TLSConfig:   tlsCfg,
		InstanceUid: c.ident.agentID,
		Callbacks: types.CallbacksStruct{
			OnConnectFunc:          c.onConnectHandler,
			OnConnectFailedFunc:    c.onConnectFailedHandler,
			OnErrorFunc:            c.onErrorHandler,
			OnMessageFunc:          c.onMessageFuncHandler,
			GetEffectiveConfigFunc: c.onGetEffectiveConfigHandler,
			// Unimplemented handlers
			// OnOpampConnectionSettingsFunc
			// OnOpampConnectionSettingsAcceptedFunc
			// OnCommandFunc
			// SaveRemoteConfigStatusFunc
		},
		PackagesStateProvider: c.packagesStateProvider,
	}

	// Start the embedded collector
	// Pass in the background context here so it's clear we need to shutdown the collector instead
	// of the context shutting it down via a cancel.
	if err := c.collector.Run(context.Background()); err != nil {
		return fmt.Errorf("collector failed to start: %w", err)
	}

	return c.opampClient.Start(ctx, settings)
}

// Disconnect disconnects from the server
func (c *Client) Disconnect(ctx context.Context) error {
	c.collector.Stop()
	return c.opampClient.Stop(ctx)
}

// client callbacks

func (c *Client) onConnectHandler() {
	c.logger.Info("Successfully connected to server")
}

func (c *Client) onConnectFailedHandler(err error) {
	c.logger.Error("Failed to connect to server", zap.Error(err))
}

func (c *Client) onErrorHandler(errResp *protobufs.ServerErrorResponse) {
	c.logger.Error("Server returned an error response", zap.String("Error", errResp.GetErrorMessage()))
}

func (c *Client) onMessageFuncHandler(ctx context.Context, msg *types.MessageData) {
	c.logger.Debug("On message handler")
	if msg.RemoteConfig != nil {
		if err := c.onRemoteConfigHandler(ctx, msg.RemoteConfig); err != nil {
			c.logger.Error("Error while processing Remote Config Change", zap.Error(err))
		}
	}
	if msg.PackagesAvailable != nil {
		if err := c.onPackagesAvailableHandler(msg.PackagesAvailable); err != nil {
			c.logger.Error("Error while processing Packages Available Change", zap.Error(err))
		}
	}
}

func (c *Client) onRemoteConfigHandler(ctx context.Context, remoteConfig *protobufs.AgentRemoteConfig) error {
	c.logger.Debug("Remote config handler")

	changed, err := c.configManager.ApplyConfigChanges(remoteConfig)
	remoteCfgStatus := &protobufs.RemoteConfigStatus{
		LastRemoteConfigHash: remoteConfig.GetConfigHash(),
		Status:               protobufs.RemoteConfigStatus_APPLIED,
	}

	// If we received and error apply it to the config
	if err != nil {
		c.logger.Error("Failed applying remote config", zap.Error(err))

		remoteCfgStatus.Status = protobufs.RemoteConfigStatus_FAILED
		remoteCfgStatus.ErrorMessage = fmt.Sprintf("Failed to apply config changes: %s", err.Error())
	}

	// Set the remote config status
	if err := c.opampClient.SetRemoteConfigStatus(remoteCfgStatus); err != nil {
		return fmt.Errorf("failed to set remote config status: %w", err)
	}

	// If we changed the config call UpdateEffectiveConfig
	if changed {
		if err := c.opampClient.UpdateEffectiveConfig(ctx); err != nil {
			return fmt.Errorf("failed to update effective config: %w", err)
		}
	}
	return nil
}

func (c *Client) onPackagesAvailableHandler(packagesAvailable *protobufs.PackagesAvailable) error {
	c.logger.Debug("Packages available handler")

	// Initialize PackageStatuses that will eventually be sent back to server
	curPackageStatuses := &protobufs.PackageStatuses{
		ServerProvidedAllPackagesHash: packagesAvailable.GetAllPackagesHash(),
		Packages:                      map[string]*protobufs.PackageStatus{},
	}

	// Retrieve last known status (this should return with minimal info even on first time)
	lastPackageStatuses, err := c.packagesStateProvider.LastReportedStatuses()

	// If there is a problem retrieving the last saved PackageStatuses, we will log the error
	// but continue on as the only thing missing will be the agent package hash.
	if err != nil {
		c.logger.Warn("Failed to retrieve last reported package statuses", zap.Error(err))
	}

	var lastPkgStatusMap map[string]*protobufs.PackageStatus
	if lastPackageStatuses != nil {
		lastPkgStatusMap = lastPackageStatuses.GetPackages()
	}

	curPackages := c.createPackageStatusMap(packagesAvailable.GetPackages(), lastPkgStatusMap)
	curPackageStatuses.Packages = curPackages

	// We may have to modify this in the future to return an error if this causes issues for the actual Update
	if err = c.packagesStateProvider.SetLastReportedStatuses(curPackageStatuses); err != nil {
		c.logger.Warn("Failed to save last reported package statuses", zap.Error(err))
	}

	if err = c.opampClient.SetPackageStatuses(curPackageStatuses); err != nil {
		return fmt.Errorf("failed to set package statuses: %w", err)
	}

	return nil
}

func (c *Client) createPackageStatusMap(
	pkgAvailMap map[string]*protobufs.PackageAvailable,
	lastPkgStatusMap map[string]*protobufs.PackageStatus) map[string]*protobufs.PackageStatus {
	pkgStatusMap := map[string]*protobufs.PackageStatus{}

	// Loop through all of the available packages sent from the server
	for name, availPkg := range pkgAvailMap {
		switch name {
		// If it's an expected package, return an installing status
		case mainPackageName:
			var agentHash []byte
			if lastPkgStatusMap != nil && lastPkgStatusMap[name] != nil {
				if lastPkgStatusMap[name].GetAgentHasVersion() != version.Version() {
					c.logger.Debug(fmt.Sprintf(
						"Version: %s and last reported package status version: %s differ",
						version.Version(),
						lastPkgStatusMap[name].GetAgentHasVersion()))
				} else {
					agentHash = lastPkgStatusMap[name].GetAgentHasHash()
				}
			}

			pkgStatusMap[name] = &protobufs.PackageStatus{
				Name:                 name,
				AgentHasVersion:      version.Version(),
				AgentHasHash:         agentHash,
				ServerOfferedVersion: availPkg.GetVersion(),
				ServerOfferedHash:    availPkg.GetHash(),
				Status:               protobufs.PackageStatus_Installed,
			}

			if version.Version() == availPkg.GetVersion() {
				if agentHash == nil {
					pkgStatusMap[name].AgentHasHash = availPkg.GetHash()
				}
			} else {
				if availPkg.GetVersion() != "" {
					pkgStatusMap[name].Status = protobufs.PackageStatus_Installing
				}
			}

			if availPkg.GetVersion() != "" && version.Version() != availPkg.GetVersion() {
				pkgStatusMap[name].Status = protobufs.PackageStatus_Installing
			}
		// If it's not an expected package, return a failed status
		default:
			pkgStatusMap[name] = &protobufs.PackageStatus{
				Name:                 name,
				ServerOfferedVersion: availPkg.GetVersion(),
				ServerOfferedHash:    availPkg.GetHash(),
				Status:               protobufs.PackageStatus_InstallFailed,
				ErrorMessage:         fmt.Sprintf("Package %s not supported", name),
			}
		}
	}

	return pkgStatusMap
}

func (c *Client) onGetEffectiveConfigHandler(_ context.Context) (*protobufs.EffectiveConfig, error) {
	c.logger.Debug("Remote Compose Effective config handler")
	return c.configManager.ComposeEffectiveConfig()
}
