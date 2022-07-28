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

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/observiq/observiq-otel-collector/packagestate"
	"github.com/observiq/observiq-otel-collector/updater/internal/install"
	"github.com/observiq/observiq-otel-collector/updater/internal/logging"
	"github.com/observiq/observiq-otel-collector/updater/internal/path"
	"github.com/observiq/observiq-otel-collector/updater/internal/rollback"
	"github.com/observiq/observiq-otel-collector/updater/internal/state"
	"github.com/observiq/observiq-otel-collector/updater/internal/version"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	var showVersion = pflag.BoolP("version", "v", false, "Prints the version of the updater and exits, if specified.")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector updater version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	// We can't create the zap logger yet, because we don't know the install dir, which is needed
	// to create the logger. So we pass a Nop logger here.
	installDir, err := path.InstallDir(zap.NewNop())
	if err != nil {
		log.Fatalf("Failed to determine install directory: %s", err)
	}

	logger, err := logging.NewLogger(installDir)
	if err != nil {
		log.Fatalf("Failed to create logger: %s", err)
	}

	// Create a monitor and load the package status file
	monitor, err := state.NewCollectorMonitor(logger, installDir)
	if err != nil {
		logger.Fatal("Failed to create monitor", zap.Error(err))
	}

	installer, err := install.NewInstaller(logger, installDir)
	if err != nil {
		logger.Fatal("Failed to create installer", zap.Error(err))
	}

	rb, err := rollback.NewRollbacker(logger, installDir)
	if err != nil {
		logger.Fatal("Failed to create rollbacker", zap.Error(err))
	}

	if err := rb.Backup(); err != nil {
		logger.Fatal("Failed to backup", zap.Error(err))
	}

	if err := installer.Install(rb); err != nil {
		logger.Error("Failed to install", zap.Error(err))

		// Set the state to failed before rollback so collector knows it failed
		if setErr := monitor.SetState(packagestate.CollectorPackageName, protobufs.PackageStatus_InstallFailed, err); setErr != nil {
			logger.Error("Failed to set state on install failure", zap.Error(setErr))
		}
		rb.Rollback()
		logger.Fatal("Rollback complete")
	}

	// Create a context with timeout to wait for a success or failed status
	checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug("Installation successful, begin monitor for success")

	// Monitor the install state
	if err := monitor.MonitorForSuccess(checkCtx, packagestate.CollectorPackageName); err != nil {
		logger.Error("Failed to install", zap.Error(err))

		// If this is not an error due to the collector setting a failed status we need to set a failed status
		if !errors.Is(err, state.ErrFailedStatus) {
			// Set the state to failed before rollback so collector knows it failed
			if setErr := monitor.SetState(packagestate.CollectorPackageName, protobufs.PackageStatus_InstallFailed, err); setErr != nil {
				logger.Error("Failed to set state on install failure", zap.Error(setErr))
			}
		}

		rb.Rollback()
		logger.Fatal("Rollback complete")
	}

	// Successful update
	logger.Info("Update Complete")
}
