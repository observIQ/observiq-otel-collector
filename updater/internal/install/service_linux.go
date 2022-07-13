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

//go:build linux

package install

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const linuxServiceName = "observiq-otel-collector"
const linuxServiceFilePath = "/usr/lib/systemd/system/observiq-otel-collector.service"

// newService returns an instance of the Service interface for managing the observiq-otel-collector service on the current OS.
func newService(latestPath string) Service {
	return linuxService{
		newServiceFilePath:       filepath.Join(latestPath, "install", "observiq-otel-collector.service"),
		serviceName:              linuxServiceName,
		installedServiceFilePath: linuxServiceFilePath,
	}
}

type linuxService struct {
	// newServiceFilePath is the file path to the new unit file
	newServiceFilePath string
	// serviceName is the name of the service
	serviceName string
	// installedServiceFilePath is the file path to the installed unit file
	installedServiceFilePath string
}

// Start the service
func (l linuxService) Start() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("systemctl", "start", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// Stop the service
func (l linuxService) Stop() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("systemctl", "stop", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running systemctl failed: %w", err)
	}
	return nil
}

// Installs the service
func (l linuxService) Install() error {
	inFile, err := os.Open(l.newServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			log.Default().Printf("Service Install: Failed to close input file: %s", err)
		}
	}()

	outFile, err := os.OpenFile(l.installedServiceFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer func() {
		err := outFile.Close()
		if err != nil {
			log.Default().Printf("Service Install: Failed to close output file: %s", err)
		}
	}()

	if _, err := io.Copy(outFile, inFile); err != nil {
		return fmt.Errorf("failed to copy service file: %w", err)
	}

	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reloading systemctl failed: %w", err)
	}

	//#nosec G204 -- serviceName is not determined by user input
	cmd = exec.Command("systemctl", "enable", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("enabling unit file failed: %w", err)
	}

	return nil
}

// Uninstalls the service
func (l linuxService) Uninstall() error {
	//#nosec G204 -- serviceName is not determined by user input
	cmd := exec.Command("systemctl", "disable", l.serviceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable unit: %w", err)
	}

	if err := os.Remove(l.installedServiceFilePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reloading systemctl failed: %w", err)
	}

	return nil
}

// installDir returns the filepath to the install directory
func installDir() (string, error) {
	return "/opt/observiq-otel-collector", nil
}
