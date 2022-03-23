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

package varnishreceiver // import "github.com/observiq/observiq-otel-collector/receiver/varnishreceiver"

import (
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/observiq/observiq-otel-collector/receiver/varnishreceiver/internal/metadata"
)

// Config defines configuration for varnish metrics receiver.
type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	Metrics                                 metadata.MetricsSettings `mapstructure:"metrics"`
	InstanceName                            string                   `mapstructure:"instance_name"`
	ExecDir                                 string                   `mapstructure:"exec_dir"`
}

// Validate validates the config.
func (c *Config) Validate() error {
	if c.InstanceName != "" {
		if _, err := os.Stat(c.InstanceName); err != nil {
			return fmt.Errorf(`"instance_name" does not exists: %w`, err)
		} else {
			lastIndex := strings.LastIndex(c.InstanceName, "/")
			if lastIndex+1 == len(c.InstanceName) {
				return fmt.Errorf(`"instance_name" has invalid trailing "/" %s`, c.InstanceName)
			}
			c.InstanceName = c.InstanceName[lastIndex+1:]
		}
	}
	if c.ExecDir != "" {
		if _, err := os.Stat(c.ExecDir); err != nil {
			return fmt.Errorf(`"exec_dir" does not exists: %w`, err)
		}
	}

	return nil
}

// SetDefaultHostname sets a default hostname if config InstanceName is empty.
func (c *Config) SetDefaultHostname() error {
	if c.InstanceName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return err
		}
		c.InstanceName = hostname
	}

	return nil
}
