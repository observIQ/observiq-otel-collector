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

package install

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

// Installer allows you to install files from latestDir into installDir,
// as well as update the service configuration using the "Install" method.
type Installer struct {
	latestDir  string
	installDir string
	svc        Service
}

// NewInstaller returns a new instance of an Installer.
func NewInstaller(tempDir string) (*Installer, error) {
	latestDir := filepath.Join(tempDir, "latest")
	installDirPath, err := installDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine install dir: %w", err)
	}

	return &Installer{
		latestDir:  latestDir,
		svc:        newService(latestDir),
		installDir: installDirPath,
	}, nil
}

// Install installs the unpacked artifacts in latestDirPath to installDirPath,
// as well as installing the new service file using the provided Service interface
func (i Installer) Install() error {
	// Stop service
	if err := i.svc.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	// install files that go to installDirPath to their correct location,
	// excluding any config files (logging.yaml, config.yaml, manager.yaml)
	if err := moveFiles(i.latestDir, i.installDir); err != nil {
		return fmt.Errorf("failed to install new files: %w", err)
	}

	// Uninstall previous service
	if err := i.svc.Uninstall(); err != nil {
		return fmt.Errorf("failed to uninstall service: %w", err)
	}

	// Install new service
	if err := i.svc.Install(); err != nil {
		return fmt.Errorf("failed to install service: %w", err)
	}

	// Start service
	if err := i.svc.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// moveFiles moves the file tree rooted at latestDirPath to installDirPath,
// skipping configuration files
func moveFiles(latestDirPath, installDirPath string) error {
	err := filepath.WalkDir(latestDirPath, func(path string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			// if there was an error walking the directory, we want to bail out.
			return err
		case d.IsDir():
			// Skip directories, we'll create them when we get a file in the directory.
			return nil
		case skipFile(path):
			// Found a config file that we should skip copying.
			return nil
		}

		cleanPath := filepath.Clean(path)

		// We want the path relative to the directory we are walking in order to calculate where the file should be
		// mirrored in the destination directory.
		relPath, err := filepath.Rel(latestDirPath, cleanPath)
		if err != nil {
			return err
		}

		// use the relative path to get the outPath (where we should write the file), and
		// to get the out directory (which we will create if it does not exist).
		outPath := filepath.Clean(filepath.Join(installDirPath, relPath))
		outDir := filepath.Dir(outPath)

		if err := os.MkdirAll(outDir, 0750); err != nil {
			return fmt.Errorf("failed to create dir: %w", err)
		}

		// Open the output file, creating it if it does not exist and truncating it.
		outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}
		defer func() {
			err := outFile.Close()
			if err != nil {
				log.Default().Printf("installFiles: Failed to close output file: %s", err)
			}
		}()

		// Open the input file for reading.
		inFile, err := os.Open(cleanPath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() {
			err := inFile.Close()
			if err != nil {
				log.Default().Printf("installFiles: Failed to close input file: %s", err)
			}
		}()

		// Copy the input file to the output file.
		if _, err := io.Copy(outFile, inFile); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk latest dir: %w", err)
	}

	return nil
}

// skipFile returns true if the given path is a special config file.
// These files should not be overwritten.
func skipFile(path string) bool {
	var configFiles = []string{
		"config.yaml",
		"logging.yaml",
		"manager.yaml",
	}

	fileName := filepath.Base(path)

	for _, f := range configFiles {
		if fileName == f {
			return true
		}
	}

	return false
}
