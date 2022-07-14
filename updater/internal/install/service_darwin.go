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

//go:build darwin

package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	darwinServiceFilePath = "/Library/LaunchDaemons/com.observiq.collector.plist"
	darwinInstallDir      = "/opt/observiq-otel-collector"
)

// newService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
func newService(latestPath string) Service {
	return &darwinService{
		newServiceFilePath:       filepath.Join(latestPath, "install", "com.observiq.collector.plist"),
		installedServiceFilePath: darwinServiceFilePath,
		installDir:               darwinInstallDir,
	}
}

type darwinService struct {
	// newServiceFilePath is the file path to the new plist file
	newServiceFilePath string
	// installedServiceFilePath is the file path to the installed plist file
	installedServiceFilePath string
	// installDir is the root directory of the main installation
	installDir string
}

// Start the service
func (d darwinService) Start() error {
	// Launchctl exits with error code 0 if the file does not exist.
	// We want to ensure that we error in this scenario.
	if _, err := os.Stat(d.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to stat installed service file: %w", err)
	}

	//#nosec G204 -- installedServiceFilePath is not determined by user input
	cmd := exec.Command("launchctl", "load", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (d darwinService) Stop() error {
	// Launchctl exits with error code 0 if the file does not exist.
	// We want to ensure that we error in this scenario.
	if _, err := os.Stat(d.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to stat installed service file: %w", err)
	}

	//#nosec G204 -- installedServiceFilePath is not determined by user input
	cmd := exec.Command("launchctl", "unload", d.installedServiceFilePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running launchctl failed: %w", err)
	}
	return nil
}

// Installs the service
func (d darwinService) Install() error {
	serviceFileBytes, err := os.ReadFile(d.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}

	expandedServiceFileBytes := replaceInstallDir(serviceFileBytes, d.installDir)
	if err := os.WriteFile(d.installedServiceFilePath, expandedServiceFileBytes, 0600); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return d.Start()
}

// Uninstalls the service
func (d darwinService) Uninstall() error {
	//#nosec G204 -- installedServiceFilePath is not determined by user input
	if err := d.Stop(); err != nil {
		return err
	}

	if err := os.Remove(d.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	return nil
}

// InstallDir returns the filepath to the install directory
func installDir() (string, error) {
	return darwinInstallDir, nil
}