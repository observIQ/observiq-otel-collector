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
	opampClient             client.OpAMPClient
	logger                  *zap.Logger
	ident                   *identity
	configManager           opamp.ConfigManager
	downloadableFileManager opamp.DownloadableFileManager
	collector               collector.Collector
	packagesStateProvider   types.PackagesStateProvider
	updatingPackage         bool

	currentConfig opamp.Config
}

// NewClientArgs arguments passed when creating a new client
type NewClientArgs struct {
	DefaultLogger *zap.Logger
	Config        opamp.Config
	Collector     collector.Collector

	TmpPath             string
	ManagerConfigPath   string
	CollectorConfigPath string
	LoggerConfigPath    string
}

// NewClient creates a new OpAmp client
func NewClient(args *NewClientArgs) (opamp.Client, error) {
	clientLogger := args.DefaultLogger.Named("opamp")

	configManager := NewAgentConfigManager(args.DefaultLogger)

	observiqClient := &Client{
		logger:                  clientLogger,
		ident:                   newIdentity(clientLogger, args.Config),
		configManager:           configManager,
		downloadableFileManager: newDownloadableFileManager(clientLogger, args.TmpPath),
		collector:               args.Collector,
		currentConfig:           args.Config,
		packagesStateProvider:   newPackagesStateProvider(clientLogger, "package_statuses.json"),
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

		// Set package status file for error (for Updater to pick up), but do not force send to Server
		c.attemptFailedInstall(fmt.Sprintf("Error while setting agent description: %s", err))

		return err
	}

	tlsCfg, err := c.currentConfig.ToTLS()
	if err != nil {
		// Set package status file for error (for Updater to pick up), but do not force send to Server
		c.attemptFailedInstall(fmt.Sprintf("Failed creating TLS config: %s", err))

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
		// Set package status file for error (for Updater to pick up), but do not force send to Server
		c.attemptFailedInstall(fmt.Sprintf("Collector failed to start: %s", err))

		return fmt.Errorf("collector failed to start: %w", err)
	}

	err = c.opampClient.Start(ctx, settings)
	if err != nil {
		// Set package status file for error (for Updater to pick up), but do not force send to Server
		c.attemptFailedInstall(fmt.Sprintf("OpAMP client failed to start: %s", err))
	}

	return err
}

// Disconnect disconnects from the server
func (c *Client) Disconnect(ctx context.Context) error {
	c.collector.Stop()
	return c.opampClient.Stop(ctx)
}

// client callbacks

func (c *Client) onConnectHandler() {
	c.logger.Info("Successfully connected to server")

	// See if we can retrieve the PackageStatuses where the main package is in an installing state
	lastPackageStatuses := c.getMainPackageInstallingLastStatuses()
	if lastPackageStatuses == nil {
		return
	}

	lastMainPackageStatus := lastPackageStatuses.Packages[mainPackageName]
	// If in the middle of an install and we just connected, this is most likely becasue the collector was just spun up fresh by the Updater.
	// If the current version matches the server offered version, this implies a good install and so we should set the PackageStatuses and
	// send it to the OpAMP Server. If the version does not match, just change the PackageStatues JSON so that the Updater can start rollback.
	if lastMainPackageStatus.ServerOfferedVersion == version.Version() {
		lastMainPackageStatus.Status = protobufs.PackageStatus_Installed
		lastMainPackageStatus.AgentHasVersion = version.Version()
		lastMainPackageStatus.AgentHasHash = lastMainPackageStatus.ServerOfferedHash

		if err := c.packagesStateProvider.SetLastReportedStatuses(lastPackageStatuses); err != nil {
			c.logger.Error("Failed to set last reported package statuses", zap.Error(err))
		}

		// Only immediately send to server on success. Rollback will send this for failure.
		if err := c.opampClient.SetPackageStatuses(lastPackageStatuses); err != nil {
			c.logger.Error("OpAMP client failed to set package statuses", zap.Error(err))
		}
	} else {
		lastMainPackageStatus.Status = protobufs.PackageStatus_InstallFailed

		if err := c.packagesStateProvider.SetLastReportedStatuses(lastPackageStatuses); err != nil {
			c.logger.Error("Failed to set last reported package statuses", zap.Error(err))
		}
	}
}

func (c *Client) onConnectFailedHandler(err error) {
	c.logger.Error("Failed to connect to server", zap.Error(err))

	// Set package status file for error (for Updater to pick up), but do not force send to Server
	c.attemptFailedInstall(fmt.Sprintf("Failed to connect to BindPlane: %s", err))
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

	// Don't respond to PackagesAvailable messages while currently installing. We use this in memory data rather than the
	// PackageStatuses persistant data in order to ensure that we don't get stuck in a stuck state
	if c.updatingPackage {
		curPackageStatuses.ErrorMessage = "Already installing new packages"
		if err := c.opampClient.SetPackageStatuses(curPackageStatuses); err != nil {
			c.logger.Error("OpAMP client failed to set package statuses", zap.Error(err))
		}
		return fmt.Errorf("failed because already installing packages")
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

	curPackages, curPackageFiles := c.createPackageMaps(packagesAvailable.GetPackages(), lastPkgStatusMap)
	curPackageStatuses.Packages = curPackages

	// This is an error because we need this file for communication during the update
	if err = c.packagesStateProvider.SetLastReportedStatuses(curPackageStatuses); err != nil {
		return fmt.Errorf("failed to save last reported package statuses: %w", err)
	}

	if err = c.opampClient.SetPackageStatuses(curPackageStatuses); err != nil {
		return fmt.Errorf("opamp client failed to set package statuses: %w", err)
	}

	// Start update if applicable
	collectorDownloadableFile := curPackageFiles[mainPackageName]
	if collectorDownloadableFile != nil {
		c.updatingPackage = true
		go c.installPackageFromFile(collectorDownloadableFile, curPackageStatuses)
	}

	return nil
}

func (c *Client) createPackageMaps(
	pkgAvailMap map[string]*protobufs.PackageAvailable,
	lastPkgStatusMap map[string]*protobufs.PackageStatus) (map[string]*protobufs.PackageStatus, map[string]*protobufs.DownloadableFile) {
	pkgStatusMap := map[string]*protobufs.PackageStatus{}
	pkgFileMap := map[string]*protobufs.DownloadableFile{}

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
				break
			}

			if availPkg.GetVersion() != "" {
				if availPkg.File != nil {
					pkgStatusMap[name].Status = protobufs.PackageStatus_Installing
					pkgFileMap[name] = availPkg.GetFile()
				} else {
					pkgStatusMap[name].Status = protobufs.PackageStatus_InstallFailed
					pkgStatusMap[name].ErrorMessage = fmt.Sprintf("Package %s does not have a valid downloadable file", name)
				}
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

	return pkgStatusMap, pkgFileMap
}

// installPackageFromFile tries to download and extract the given tarball and then start up the new Updater binary that was
// inside of it
func (c *Client) installPackageFromFile(file *protobufs.DownloadableFile, curPackageStatuses *protobufs.PackageStatuses) {
	if fileManagerErr := c.downloadableFileManager.FetchAndExtractArchive(file); fileManagerErr != nil {
		// Change existing status to show that install failed and get ready to send
		curPackageStatuses.Packages[mainPackageName].Status = protobufs.PackageStatus_InstallFailed
		curPackageStatuses.Packages[mainPackageName].ErrorMessage =
			fmt.Sprintf("Failed to download and verify package %s's downloadable file", mainPackageName)

		if err := c.packagesStateProvider.SetLastReportedStatuses(curPackageStatuses); err != nil {
			c.logger.Error("Failed to save last reported package statuses", zap.Error(err))
		}

		if err := c.opampClient.SetPackageStatuses(curPackageStatuses); err != nil {
			c.logger.Error("OpAMP client failed to set package statuses", zap.Error(err))
		}

		c.updatingPackage = false
		return
	}

	// TODO: Start Updater

	return
}

func (c *Client) onGetEffectiveConfigHandler(_ context.Context) (*protobufs.EffectiveConfig, error) {
	c.logger.Debug("Remote Compose Effective config handler")
	return c.configManager.ComposeEffectiveConfig()
}

// attemptFailedInstall sets PackageStatuses status to failed and error message if we are in the middle of an install.
// This should allow the updater to pick this up and start the rollback process
func (c Client) attemptFailedInstall(errMsg string) {
	// See if we can retrieve the PackageStatuses where the main package is in an installing state
	lastPackageStatuses := c.getMainPackageInstallingLastStatuses()
	if lastPackageStatuses == nil {
		return
	}

	lastMainPackageStatus := lastPackageStatuses.Packages[mainPackageName]
	lastMainPackageStatus.Status = protobufs.PackageStatus_InstallFailed
	lastMainPackageStatus.ErrorMessage = errMsg

	if err := c.packagesStateProvider.SetLastReportedStatuses(lastPackageStatuses); err != nil {
		c.logger.Error("Failed to set last reported package statuses", zap.Error(err))
	}
}

func (c Client) getMainPackageInstallingLastStatuses() *protobufs.PackageStatuses {
	lastPackageStatuses, err := c.packagesStateProvider.LastReportedStatuses()
	if err != nil {
		c.logger.Error("Failed to retrieve last reported package statuses", zap.Error(err))
		return nil
	}

	// If we have no info on our main package, nothing else to do
	if lastPackageStatuses == nil || lastPackageStatuses.Packages == nil || lastPackageStatuses.Packages[mainPackageName] == nil {
		c.logger.Warn("Failed to retrieve last reported package statuses for main package")
		return nil
	}

	lastMainPackageStatus := lastPackageStatuses.Packages[mainPackageName]

	// If we were not installing before the connection, nothing else to do
	if lastMainPackageStatus.Status != protobufs.PackageStatus_Installing {
		return nil
	}

	return lastPackageStatuses
}
