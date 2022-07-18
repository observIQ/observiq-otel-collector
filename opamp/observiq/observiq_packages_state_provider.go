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

// Package observiq contains OpAmp structures compatible with the observiq client
package observiq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/internal/version"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

// Ensure interface is satisfied
var _ types.PackagesStateProvider = (*packagesStateProvider)(nil)

// packagesStateProvider represents a PackagesStateProvider which uses a json file to persist PackageStatuses
type packagesStateProvider struct {
	jsonPath string
	logger   *zap.Logger
}

type packageState struct {
	Name          string                         `json:"name"`
	AgentVersion  string                         `json:"agent_version"`
	AgentHash     []byte                         `json:"agent_hash"`
	ServerVersion string                         `json:"server_version"`
	ServerHash    []byte                         `json:"server_hash"`
	Status        protobufs.PackageStatus_Status `json:"status"`
	ErrorMessage  string                         `json:"error_message"`
}
type packageStates struct {
	AllPackagesHash []byte                   `json:"all_packages_hash"`
	AllErrorMessage string                   `json:"all_error_message"`
	PackageStates   map[string]*packageState `json:"package_states"`
}

// newPackagesStateProvider creates a new OpAmp PackagesStateProvider
func newPackagesStateProvider(logger *zap.Logger, jsonPath string) types.PackagesStateProvider {
	return &packagesStateProvider{
		jsonPath: filepath.Clean(jsonPath),
		logger:   logger,
	}
}

// AllPackagesHash not implemented so returns an error with this info
func (p *packagesStateProvider) AllPackagesHash() ([]byte, error) {
	p.logger.Debug("Retrieve all packages hash")

	return nil, errors.New("method not implemented: PackageStateProvider AllPackagesHash")
}

// SetAllPackagesHash not implemented so returns an error with this info
func (p *packagesStateProvider) SetAllPackagesHash(_ []byte) error {
	p.logger.Debug("Set all packages hash")

	return errors.New("method not implemented: PackageStateProvider SetAllPackagesHash")
}

// Packages not implemented so returns an error with this info
func (p *packagesStateProvider) Packages() ([]string, error) {
	p.logger.Debug("Retrieve package names")

	return nil, errors.New("method not implemented: PackageStateProvider Packages")
}

// PackageState not implemented so returns an error with this info
func (p *packagesStateProvider) PackageState(_ string) (state types.PackageState, err error) {
	p.logger.Debug("Retrieve package state")

	packageState := types.PackageState{}

	return packageState, errors.New("method not implemented: PackageStateProvider PackageState")
}

// SetPackageState not implemented so returns an error with this info
func (p *packagesStateProvider) SetPackageState(_ string, _ types.PackageState) error {
	p.logger.Debug("Set package state")

	return errors.New("method not implemented: PackageStateProvider SetPackageState")
}

// CreatePackage not implemented so returns an error with this info
func (p *packagesStateProvider) CreatePackage(_ string, _ protobufs.PackageAvailable_PackageType) error {
	p.logger.Debug("Create package")

	return errors.New("method not implemented: PackageStateProvider CreatePackage")
}

// FileContentHash not implemented so returns an error with this info
func (p *packagesStateProvider) FileContentHash(_ string) ([]byte, error) {
	p.logger.Debug("Retrieve package content hash")

	return nil, errors.New("method not implemented: PackageStateProvider FileContentHash")
}

// UpdateContent not implemented so returns an error with this info
func (p *packagesStateProvider) UpdateContent(_ context.Context, _ string, _ io.Reader, _ []byte) error {
	p.logger.Debug("Update package content")

	return errors.New("method not implemented: PackageStateProvider UpdateContent")
}

// DeletePackage not implemented so returns an error with this info
func (p *packagesStateProvider) DeletePackage(_ string) error {
	p.logger.Debug("Delete package")

	return errors.New("method not implemented: PackageStateProvider DeletePackage")
}

// LastReportedStatuses retrieves the PackagesStatuses from a saved json file
func (p *packagesStateProvider) LastReportedStatuses() (*protobufs.PackageStatuses, error) {
	p.logger.Debug("Retrieve last reported package statuses")

	packages := map[string]*protobufs.PackageStatus{
		mainPackageName: {
			Name:            mainPackageName,
			AgentHasVersion: version.Version(),
			Status:          protobufs.PackageStatus_Installed,
		},
	}
	packageStatuses := &protobufs.PackageStatuses{
		Packages: packages,
	}

	// If there's a problem reading the file, we just return a barebones PackageStatuses.
	if _, err := os.Stat(p.jsonPath); errors.Is(err, os.ErrNotExist) {
		p.logger.Debug("Package statuses json doesn't exist")
		return packageStatuses, nil
	}

	statusesBytes, err := os.ReadFile(p.jsonPath)
	if err != nil {
		return packageStatuses, fmt.Errorf("failed to read package statuses json: %w", err)
	}

	var packageStates packageStates
	if err := json.Unmarshal(statusesBytes, &packageStates); err != nil {
		return packageStatuses, fmt.Errorf("failed to unmarshal package statuses: %w", err)
	}

	return packageStatesToStatuses(packageStates), nil
}

// SetLastReportedStatuses saves the given PackageStatuses into a json file
func (p *packagesStateProvider) SetLastReportedStatuses(statuses *protobufs.PackageStatuses) error {
	p.logger.Debug("Set last reported package statuses")

	// If there is any problem saving the new package statuses, make sure that we delete any existing file
	// in order to not have outdated data as its better to start fresh
	if err := os.Remove(p.jsonPath); err != nil {
		p.logger.Debug("Failed to delete package statuses json", zap.Error(err))
	}

	states := packageStatusesToStates(statuses)

	data, err := json.Marshal(states)
	if err != nil {
		return fmt.Errorf("failed to marshal package statuses: %w", err)
	}

	// Write data to a package_statuses.json file, with 0600 file permission
	if err := os.WriteFile(p.jsonPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write package statuses json: %w", err)
	}

	return nil
}

func packageStatusesToStates(statuses *protobufs.PackageStatuses) *packageStates {
	states := &packageStates{
		AllPackagesHash: statuses.GetServerProvidedAllPackagesHash(),
		AllErrorMessage: statuses.GetErrorMessage(),
	}

	packageStates := map[string]*packageState{}
	for name, packageStatus := range statuses.Packages {
		packageState := &packageState{
			Name:          packageStatus.GetName(),
			AgentVersion:  packageStatus.GetAgentHasVersion(),
			AgentHash:     packageStatus.GetAgentHasHash(),
			ServerVersion: packageStatus.GetServerOfferedVersion(),
			ServerHash:    packageStatus.GetServerOfferedHash(),
			Status:        packageStatus.GetStatus(),
			ErrorMessage:  packageStatus.GetErrorMessage(),
		}
		packageStates[name] = packageState
	}
	states.PackageStates = packageStates

	return states
}

func packageStatesToStatuses(states packageStates) *protobufs.PackageStatuses {
	statuses := &protobufs.PackageStatuses{
		ServerProvidedAllPackagesHash: states.AllPackagesHash,
		ErrorMessage:                  states.AllErrorMessage,
	}

	packages := map[string]*protobufs.PackageStatus{}
	for name, packageState := range states.PackageStates {
		packageStatus := &protobufs.PackageStatus{
			Name:                 packageState.Name,
			AgentHasVersion:      packageState.AgentVersion,
			AgentHasHash:         packageState.AgentHash,
			ServerOfferedVersion: packageState.ServerVersion,
			ServerOfferedHash:    packageState.ServerHash,
			Status:               packageState.Status,
			ErrorMessage:         packageState.ErrorMessage,
		}
		packages[name] = packageStatus
	}
	statuses.Packages = packages

	return statuses
}
