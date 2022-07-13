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

// an elevated user is needed to run the service tests
//go:build linux && integration

package install

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// NOTE: These tests must run as root in order to pass
func TestLinuxServiceInstall(t *testing.T) {
	t.Run("Test install + uninstall", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		//We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = l.Uninstall()
		require.NoError(t, err)
		require.NoFileExists(t, installedServicePath)

		//Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test stop + start", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Install()
		require.NoError(t, err)
		require.FileExists(t, installedServicePath)

		// We want to check that the service was actually loaded
		requireServiceLoadedStatus(t, true)

		err = l.Start()
		require.NoError(t, err)

		requireServiceRunningStatus(t, true)

		err = l.Stop()
		require.NoError(t, err)

		requireServiceRunningStatus(t, false)

		err = l.Uninstall()
		require.NoError(t, err)
		require.NoFileExists(t, installedServicePath)

		// Make sure the service is no longer listed
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test invalid path for input file", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "does-not-exist.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Install()
		require.ErrorContains(t, err, "failed to open input file")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Test invalid path for output file for install", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/dir-does-not-exist/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Install()
		require.ErrorContains(t, err, "failed to open output file")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Uninstall fails if not installed", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Uninstall()
		require.ErrorContains(t, err, "failed to disable unit")
		requireServiceLoadedStatus(t, false)
	})

	t.Run("Start fails if service not found", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Start()
		require.ErrorContains(t, err, "running systemctl failed")
	})

	t.Run("Stop fails if service not found", func(t *testing.T) {
		installedServicePath := "/usr/lib/systemd/system/linux-service.service"
		uninstallService(t, installedServicePath, "linux-service")

		l := &linuxService{
			newServiceFilePath:       filepath.Join("testdata", "linux-service.service"),
			serviceName:              "linux-service",
			installedServiceFilePath: installedServicePath,
		}

		err := l.Stop()
		require.ErrorContains(t, err, "running systemctl failed")
	})
}

// uninstallService is a helper that uninstalls the service manually for test setup, in case it is somehow leftover.
func uninstallService(t *testing.T, installedPath, serviceName string) {
	cmd := exec.Command("systemctl", "stop", serviceName)
	_ = cmd.Run()

	cmd = exec.Command("systemctl", "disable", serviceName)
	_ = cmd.Run()

	err := os.RemoveAll(installedPath)
	require.NoError(t, err)

	cmd = exec.Command("systemctl", "daemon-reload")
	_ = cmd.Run()
}

const exitCodeServiceNotFound = 4
const exitCodeServiceInactive = 3

func requireServiceLoadedStatus(t *testing.T, loaded bool) {
	t.Helper()

	cmd := exec.Command("systemctl", "status", "linux-service")
	err := cmd.Run()
	require.Error(t, err, "expected non-zero exit code from 'systemctl status linux-service'")

	eErr, ok := err.(*exec.ExitError)
	if loaded {
		// If the service should be loaded, then we expect a 0 exit code, so no error is given
		require.Equal(t, exitCodeServiceInactive, eErr.ExitCode(), "unexpected exit code when asserting service is unloaded: %d", eErr.ExitCode())
		return
	}

	require.True(t, ok, "systemctl status exited with non-ExitError: %s", eErr)
	require.Equal(t, exitCodeServiceNotFound, eErr.ExitCode(), "unexpected exit code when asserting service is unloaded: %d", eErr.ExitCode())
}

func requireServiceRunningStatus(t *testing.T, running bool) {
	cmd := exec.Command("systemctl", "status", "linux-service")
	err := cmd.Run()

	if running {
		// exit code 0 indicates service is loaded & running
		require.NoError(t, err)
		return
	}

	eErr, ok := err.(*exec.ExitError)
	require.True(t, ok, "systemctl status exited with non-ExitError: %s", eErr)
	require.Equal(t, exitCodeServiceInactive, eErr.ExitCode(), "unexpected exit code when asserting service is not running: %d", eErr.ExitCode())
}
