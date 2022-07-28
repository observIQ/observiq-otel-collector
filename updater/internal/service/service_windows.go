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

//go:build windows

package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/kballard/go-shellquote"
	"github.com/observiq/observiq-otel-collector/updater/internal/path"
)

const (
	defaultProductName           = "observIQ Distro for OpenTelemetry Collector"
	defaultServiceName           = "observiq-otel-collector"
	uninstallServicePollInterval = 50 * time.Millisecond
	serviceNotExistErrStr        = "The specified service does not exist as an installed service."
)

// Option is an extra option for creating a Service
type Option func(winSvc *windowsService)

// WithServiceFile returns an option setting the service file to use when updating using the service
func WithServiceFile(svcFilePath string) Option {
	return func(winSvc *windowsService) {
		winSvc.newServiceFilePath = svcFilePath
	}
}

// NewService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
func NewService(logger *zap.Logger, installDir string, opts ...Option) Service {
	winSvc := &windowsService{
		newServiceFilePath: filepath.Join(path.ServiceFileDir(installDir), "windows_service.json"),
		serviceName:        defaultServiceName,
		productName:        defaultProductName,
		installDir:         installDir,
		logger:             logger.Named("windows-service"),
	}

	for _, opt := range opts {
		opt(winSvc)
	}

	return winSvc
}

type windowsService struct {
	// newServiceFilePath is the file path to the new unit file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// productName is the name of the installed product
	productName string
	installDir  string
	logger      *zap.Logger
}

// Start the service
func (w windowsService) Start() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	if err := s.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Stop the service
func (w windowsService) Stop() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	if _, err := s.Control(svc.Stop); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Installs the service
func (w windowsService) install() error {
	// parse the service definition from disk
	wsc, err := readWindowsServiceConfig(w.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	// expand the arguments to be properly formatted (expand [INSTALLDIR], clean '&quot;' to be '"')
	expandArguments(wsc, w.installDir)

	// Split the arguments; Arguments are "shell-like", in that they may contain spaces, and can be quoted to indicate that.
	splitArgs, err := shellquote.Split(wsc.Service.Arguments)
	if err != nil {
		return fmt.Errorf("failed to parse arguments in service config: %w", err)
	}

	// Get the start type
	startType, delayed, err := winapiStartType(wsc.Service.Start)
	if err != nil {
		return fmt.Errorf("failed to parse start type in service config: %w", err)
	}

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	// Create the service using the service manager.
	s, err := m.CreateService(w.serviceName,
		filepath.Join(w.installDir, wsc.Path),
		mgr.Config{
			Description:      wsc.Service.Description,
			DisplayName:      wsc.Service.DisplayName,
			StartType:        startType,
			DelayedAutoStart: delayed,
		},
		splitArgs...,
	)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer s.Close()

	return nil
}

// Uninstalls the service
func (w windowsService) uninstall() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return fmt.Errorf("failed to open service: %w", err)
	}

	// Note on deleting services in windows:
	// Deleting the service is not immediate. If there are open handles to the service (e.g. you have services.msc open)
	// then the service deletion will be delayed, perhaps indefinitely. However, we want this logic to be synchronous, so
	// we will try to wait for the service to actually be deleted.
	if err = s.Delete(); err != nil {
		sCloseErr := s.Close()
		if sCloseErr != nil {
			log.Default().Printf("Failed to close service: %s\n", err)
		}
		return fmt.Errorf("failed to delete service: %w", err)
	}

	if err := s.Close(); err != nil {
		return fmt.Errorf("failed to close service: %w", err)
	}

	// Wait for the service to actually be deleted:
	for {
		s, err := m.OpenService(w.serviceName)
		if err != nil {
			if err.Error() == serviceNotExistErrStr {
				// This is expected when the service is uninstalled.
				break
			}
			return fmt.Errorf("got unexpected error when waiting for service deletion: %w", err)
		}

		if err := s.Close(); err != nil {
			return fmt.Errorf("failed to close service: %w", err)
		}
		// rest with the handle closed to let the service manager remove the service
		time.Sleep(uninstallServicePollInterval)
	}
	return nil
}

func (w windowsService) Update() error {
	if err := w.uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall old service: %w", err)
	}

	if err := w.install(); err != nil {
		return fmt.Errorf("failed to install new service: %w", err)
	}

	return nil
}

func (w windowsService) Backup() error {

	wsc, err := w.currentServiceConfig()
	if err != nil {
		return fmt.Errorf("failed to construct service config: %w", err)
	}

	// Marshal config as json
	wscBytes, err := json.Marshal(wsc)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Open with O_EXCL to fail if the file already exists
	f, err := os.OpenFile(path.BackupServiceFile(w.installDir), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to create backup service file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			w.logger.Error("Failed to close backup service file", zap.Error(err))
		}
	}()

	// finally, write the config out so we can rollback.
	if _, err := f.Write(wscBytes); err != nil {
		return fmt.Errorf("failed to write backup service config: %w", err)
	}

	return nil
}

// windowsServiceConfig defines how the service should be configured, including the entrypoint for the service.
type windowsServiceConfig struct {
	// Path is the file that will be executed for the service. It is relative to the install directory.
	Path string `json:"path"`
	// Configuration for the service (e.g. start type, display name, desc)
	Service windowsServiceDefinitionConfig `json:"service"`
}

// windowsServiceDefinitionConfig defines how the service should be configured.
// Name is a part of the on disk config, but we keep the service name hardcoded; We do not want to use a different service name.
type windowsServiceDefinitionConfig struct {
	// Start gives the start type of the service.
	// See: https://wixtoolset.org/documentation/manual/v3/xsd/wix/serviceinstall.html
	Start string `json:"start"`
	// DisplayName is the human-readable name of the service.
	DisplayName string `json:"display-name"`
	// Description is a human-readable description of the service.
	Description string `json:"description"`
	// Arguments is a list of space-separated
	Arguments string `json:"arguments"`
}

// readWindowsServiceConfig reads the service config from the file at the given path
func readWindowsServiceConfig(path string) (*windowsServiceConfig, error) {
	cleanPath := filepath.Clean(path)
	b, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var wsc windowsServiceConfig
	err = json.Unmarshal(b, &wsc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &wsc, nil
}

// expandArguments expands [INSTALLDIR] to the actual install directory and
// expands '&quot;' to the literal '"'
func expandArguments(wsc *windowsServiceConfig, installDir string) {
	wsc.Service.Arguments = string(replaceInstallDir([]byte(wsc.Service.Arguments), installDir))
	wsc.Service.Arguments = strings.ReplaceAll(wsc.Service.Arguments, "&quot;", `"`)
}

func (w windowsService) currentServiceConfig() (*windowsServiceConfig, error) {
	m, err := mgr.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(w.serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open service: %w", err)
	}
	defer s.Close()

	// Get the current config of the service
	conf, err := s.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to get service config: %w", err)
	}

	fullBinaryPath, argString, err := splitServiceBinaryName(conf.BinaryPathName)
	if err != nil {
		return nil, fmt.Errorf("failed to split service BinaryPathName: %w", err)
	}

	// In the original config, the Path is the main binary path, relative to the install directory.
	binaryPath, err := filepath.Rel(w.installDir, fullBinaryPath)
	if err != nil {
		return nil, fmt.Errorf("could not find service exe relative to install dir: %w", err)
	}

	// Convert windows api start type to the config file service type
	confStartType, err := configStartType(conf.StartType, conf.DelayedAutoStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get start type: %w", err)
	}

	// Construct the config
	return &windowsServiceConfig{
		Path: binaryPath,
		Service: windowsServiceDefinitionConfig{
			Start:       confStartType,
			DisplayName: conf.DisplayName,
			Description: conf.Description,
			Arguments:   argString,
		},
	}, nil
}

func splitServiceBinaryName(binaryPathName string) (binaryPath, argString string, err error) {
	// Split the service arguments into an array of arguments
	args, err := shellquote.Split(binaryPathName)
	if err != nil {
		return "", "", fmt.Errorf("failed to split service config args: %w", err)
	}

	// The first argument is always the binary name; If the length of the array is 0, we know this is an invalid argument list.
	if len(args) < 1 {
		return "", "", fmt.Errorf("no binary specified in service config")
	}

	// The absolute path to the binary is the first argument
	binaryPath = args[0]

	// Stored argument string doesn't include the binary path (first arg)
	args = args[1:]

	// Args should end up being a string, where literal quotes are "&quot;"
	argString = shellquote.Join(args...)
	// shellquote uses ' to quote, so we convert those to "&quot;"
	argString = strings.ReplaceAll(argString, "'", "&quot;")

	return binaryPath, argString, nil
}

// winapiStartType converts the start type from the windowsServiceConfig to a start type recognizable by the windows
// service API
func winapiStartType(cfgStartType string) (startType uint32, delayed bool, err error) {
	switch cfgStartType {
	case "auto":
		// Automatically starts on system bootup.
		startType = mgr.StartAutomatic
	case "demand":
		// Must be started manually
		startType = mgr.StartManual
	case "disabled":
		// Does not start, must be enabled to run.
		startType = mgr.StartDisabled
	case "delayed":
		// Boots automatically on start, but AFTER bootup has completed.
		startType = mgr.StartAutomatic
		delayed = true
	default:
		err = fmt.Errorf("invalid start type in service config: %s", cfgStartType)
	}
	return
}

func configStartType(winapiStartType uint32, delayed bool) (string, error) {
	switch winapiStartType {
	case mgr.StartAutomatic:
		if delayed {
			return "delayed", nil
		}
		return "auto", nil
	case mgr.StartDisabled:
		return "disabled", nil
	case mgr.StartManual:
		return "manual", nil
	default:
		return "", fmt.Errorf("invalid winapi start type: %d", winapiStartType)
	}
}
