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

package path

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

const defaultProductName = "observIQ Distro for OpenTelemetry Collector"

// InstallDirFromRegistry gets the installation dir of the given product from the Windows Registry
func InstallDirFromRegistry(productName string) (string, error) {
	// this key is created when installing using the MSI installer
	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Uninstall\%s`, productName)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.READ)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %w", err)
	}
	defer func() {
		err := key.Close()
		if err != nil {
			log.Default().Printf("InstallDirFromRegistry: failed to close registry key")
		}
	}()

	// This value ("InstallLocation") contains the path to the install folder.
	val, _, err := key.GetStringValue("InstallLocation")
	if err != nil {
		return "", fmt.Errorf("failed to read install dir: %w", err)
	}

	return val, nil
}

// InstallDir returns the filepath to the install directory
func InstallDir() (string, error) {
	return InstallDirFromRegistry(defaultProductName)
}
