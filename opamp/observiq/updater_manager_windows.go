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

package observiq

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Ensure interface is satisfied
var _ updaterManager = (*windowsUpdaterManager)(nil)

// windowsUpdaterManager handles starting a Updater binary and watching it for failure with a timeout
type windowsUpdaterManager struct {
	tmpPath string
	logger  *zap.Logger
}

// newUpdaterManager creates a new updaterManager
func newUpdaterManager(defaultLogger *zap.Logger, tmpPath string) updaterManager {
	updaterName = "updater.exe"
	return &windowsUpdaterManager{
		tmpPath: filepath.Clean(tmpPath),
		logger:  defaultLogger.Named("updater manager"),
	}
}

// StartAndMonitorUpdater will start the Updater binary and wait to see if it finishes unexpectedly.
// While waiting for Updater, it should kill the collector and we should never execute any code past running it
func (m windowsUpdaterManager) StartAndMonitorUpdater() error {
	updaterPath := filepath.Join(m.tmpPath, updaterDir, updaterName)
	absTmpPath, err := filepath.Abs(m.tmpPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of tmp dir: %w", err)
	}
	//#nosec G204 -- paths are not determined via user input
	cmd := exec.Command(updaterPath, "--tmpdir", absTmpPath)

	// Start does not block
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("updater had an issue while starting: %w", err)
	}

	// See if we're still alive after 5 seconds
	time.Sleep(5 * time.Second)

	// Updater might already be killed
	if err := cmd.Process.Kill(); err != nil {
		m.logger.Error("Failed to kill failed Updater", zap.Error(err))
	}

	// Ideally we should not get here as we will be killed by the updater.
	// Updater should either exit before us with error or we die before it does.
	return errors.New("updater failed to update collector")
}
