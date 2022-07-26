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
	"os"
	"time"

	"github.com/observiq/observiq-otel-collector/packagestate"
	"github.com/observiq/observiq-otel-collector/updater/internal/install"
	"github.com/observiq/observiq-otel-collector/updater/internal/rollback"
	"github.com/observiq/observiq-otel-collector/updater/internal/state"
	"github.com/observiq/observiq-otel-collector/updater/internal/version"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

// Unimplemented
func main() {
	var showVersion = pflag.BoolP("version", "v", false, "Prints the version of the collector and exits, if specified.")
	var tmpDir = pflag.String("tmpdir", "", "Temporary directory for artifacts. Parent of the 'rollback' directory.")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector updater version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	if *tmpDir == "" {
		log.Println("The --tmpdir flag must be specified!")
		pflag.PrintDefaults()
		os.Exit(1)
	}

	// Create a monitor and load the package status file
	// TODO replace nop logger with real one
	monitor, err := state.NewCollectorMonitor(zap.NewNop())
	if err != nil {
		log.Fatalln("Failed to create monitor:", err)
	}

	installer, err := install.NewInstaller(*tmpDir)
	if err != nil {
		log.Fatalf("Failed to create installer: %s", err)
	}

	rb, err := rollback.NewRollbacker(*tmpDir)
	if err != nil {
		log.Fatalf("Failed to create rollbacker: %s", err)
	}

	if err := rb.Backup(); err != nil {
		log.Fatalf("Failed to backup: %s", err)
	}

	if err := installer.Install(rb); err != nil {
		log.Default().Printf("Failed to install: %s", err)

		// Set the state to failed before rollback so collector knows it failed
		if setErr := monitor.SetState(packagestate.DefaultFileName, protobufs.PackageStatus_InstallFailed, err); setErr != nil {
			log.Println("Failed to set state on install failure:", setErr)
		}
		rb.Rollback()
		log.Default().Fatalf("Rollback complete")
	}

	// Create a context with timeout to wait for a success or failed status
	checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Monitor the install state
	if err := monitor.MonitorForSuccess(checkCtx, packagestate.DefaultFileName); err != nil {
		log.Println("Failed to install:", err)

		// If this is not an error due to the collector setting a failed status we need to set a failed status
		if !errors.Is(err, state.ErrFailedStatus) {
			// Set the state to failed before rollback so collector knows it failed
			if setErr := monitor.SetState(packagestate.DefaultFileName, protobufs.PackageStatus_InstallFailed, err); setErr != nil {
				log.Println("Failed to set state on install failure:", setErr)
			}
		}

		rb.Rollback()
		log.Fatalln("Rollback complete")
	}

	// Successful update
	log.Println("Update Complete")
}
