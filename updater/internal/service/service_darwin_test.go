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

//go:build darwin && integration

package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"github.com/stretchr/testify/require"
)

func TestDarwinServiceInstall(t *testing.T) {
	t.Run("Test install + uninstall", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = d.uninstall()
		require.NoError(t, err)
		require.NoFileExists(t, installedServicePath)

		// Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test stop + start", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		// TODO: Do this automagically
		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = d.Start()
		require.NoError(t, err)

		requireServiceRunning(t)

		err = d.Stop()
		require.NoError(t, err)

		requireServiceLoadedStatus(t, false)

		err = d.uninstall()
		require.NoError(t, err)
		require.NoFileExists(t, installedServicePath)

		// Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test invalid path for input file", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "does-not-exist.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.install()
		require.ErrorContains(t, err, "failed to open input file")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test invalid path for output file for install", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "directory-does-not-exist", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.install()
		require.ErrorContains(t, err, "failed to write service file")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Uninstall fails if not installed", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.uninstall()
		require.ErrorContains(t, err, "failed to stat installed service file")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Start fails if service not found", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.Start()
		require.ErrorContains(t, err, "failed to stat installed service file")
	})

	t.Run("Stop fails if service not found", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		d := &darwinService{
			newServiceFilePath:       filepath.Join("testdata", "darwin-service.plist"),
			installedServiceFilePath: installedServicePath,
		}

		err := d.Stop()
		require.ErrorContains(t, err, "failed to stat installed service file")
	})

	t.Run("Backup installed service succeeds", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		newServiceFile := filepath.Join("testdata", "darwin-service.plist")
		serviceFileContents, err := os.ReadFile(newServiceFile)
		require.NoError(t, err)

		d := &darwinService{
			newServiceFilePath:       newServiceFile,
			installedServiceFilePath: installedServicePath,
		}

		err = d.install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		require.NoError(t, d.Stop())

		backupServiceDir := t.TempDir()
		err = d.Backup(backupServiceDir)
		require.NoError(t, err)
		require.FileExists(t, path.BackupServiceFile(backupServiceDir))

		backupServiceContents, err := os.ReadFile(path.BackupServiceFile(backupServiceDir))

		require.Equal(t, serviceFileContents, backupServiceContents)
		require.NoError(t, d.uninstall())
	})

	t.Run("Backup installed service fails if not installed", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		newServiceFile := filepath.Join("testdata", "darwin-service.plist")

		d := &darwinService{
			newServiceFilePath:       newServiceFile,
			installedServiceFilePath: installedServicePath,
		}

		backupServiceDir := t.TempDir()
		err := d.Backup(backupServiceDir)
		require.ErrorContains(t, err, "failed to copy service file")
	})

	t.Run("Backup installed service fails if output file already exists", func(t *testing.T) {
		installedServicePath := filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents", "darwin-service.plist")

		uninstallService(t, installedServicePath)

		newServiceFile := filepath.Join("testdata", "darwin-service.plist")

		d := &darwinService{
			newServiceFilePath:       newServiceFile,
			installedServiceFilePath: installedServicePath,
		}

		err := d.install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		require.NoError(t, d.Stop())

		backupServiceDir := t.TempDir()
		// Write the backup file before creating it; Backup should
		// not ever overwrite an existing file
		os.WriteFile(path.BackupServiceFile(backupServiceDir), []byte("file exists"), 0600)

		err = d.Backup(backupServiceDir)
		require.ErrorContains(t, err, "failed to copy service file")
	})
}

// uninstallService is a helper that uninstalls the service manually for test setup, in case it is somehow leftover.
func uninstallService(t *testing.T, installedPath string) {
	t.Helper()

	cmd := exec.Command("launchctl", "unload", installedPath)
	// May already be unloaded; We'll ignore the error.
	_ = cmd.Run()

	err := os.RemoveAll(installedPath)
	require.NoError(t, err)
}

const exitCodeServiceNotFound = 113

func requireServiceLoadedStatus(t *testing.T, loaded bool) {
	t.Helper()

	cmd := exec.Command("launchctl", "list", "darwin-service")
	err := cmd.Run()
	if loaded {
		// If the service should be loaded, then we expect a 0 exit code, so no error is given
		require.NoError(t, err)
		return
	}

	eErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "launchctl list exited with non-ExitError: %s", eErr)
	require.Equal(t, exitCodeServiceNotFound, eErr.ExitCode(), "unexpected exit code when asserting service is unloaded: %d", eErr.ExitCode())
}

var descriptionPIDRegex = regexp.MustCompile(`\s*"PID" = \d+;`)

func requireServiceRunning(t *testing.T) {
	t.Helper()

	cmd := exec.Command("launchctl", "list", "darwin-service")
	out, err := cmd.Output()
	require.NoError(t, err)
	matches := descriptionPIDRegex.Match(out)
	require.True(t, matches, "Service should be running, but it was not found in launchctl list")
}
