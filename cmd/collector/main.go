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
	"fmt"
	"log"

	"github.com/observiq/observiq-otel-collector/collector"
	"github.com/observiq/observiq-otel-collector/internal/logging"
	"github.com/observiq/observiq-otel-collector/internal/service"
	"github.com/observiq/observiq-otel-collector/internal/version"
	"go.uber.org/zap"

	"github.com/spf13/pflag"
)

func main() {
	var configPaths = pflag.StringSlice("config", []string{"./config.yaml"}, "the collector config path")
	var loggingPath = pflag.String("logging", "./logging.yaml", "the collector logging config path")
	_ = pflag.String("log-level", "", "not implemented") // TEMP(jsirianni): Required for OTEL k8s operator
	var showVersion = pflag.BoolP("version", "v", false, "prints the version of the collector")
	pflag.Parse()

	if *showVersion {
		fmt.Println("observiq-otel-collector version", version.Version())
		fmt.Println("commit:", version.GitHash())
		fmt.Println("built at:", version.Date())
		return
	}

	logOpts, err := logOptions(loggingPath)
	if err != nil {
		log.Fatalf("Failed to get log options: %v", err)
	}

	// logOpts will override options here
	logger, err := zap.NewProduction(logOpts...)
	if err != nil {
		log.Fatalf("Failed to set up logger: %v", err)
	}

	col := collector.New(*configPaths, version.Version(), logOpts)
	err = service.RunService(logger, service.NewStandaloneCollectorService(col))
	if err != nil {
		logger.Fatal("RunService returned error", zap.Error(err))
	}
}

func logOptions(loggingConfigPath *string) ([]zap.Option, error) {
	if loggingConfigPath == nil {
		return nil, nil
	}

	l, err := logging.NewLoggerConfig(*loggingConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger config: %w", err)
	}

	return l.Options()
}
