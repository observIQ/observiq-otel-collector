package collector

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/extension/extensionhelper"
)

func TestCollectorRunValid(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	err = collector.Run()
	require.NoError(t, err)

	status := collector.Status()
	require.True(t, status.Running)
	require.NoError(t, status.Err)

	collector.Stop()
	status = collector.Status()
	require.False(t, status.Running)
}

func TestCollectorRunMultiple(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	for i := 1; i < 5; i++ {
		attempt := fmt.Sprintf("Attempt %d", i)
		t.Run(attempt, func(t *testing.T) {
			err = collector.Run()
			require.NoError(t, err)

			status := collector.Status()
			require.True(t, status.Running)
			require.NoError(t, status.Err)

			collector.Stop()
			status = collector.Status()
			require.False(t, status.Running)
		})
	}
}

func TestCollectorRunInvalidConfig(t *testing.T) {
	settings, err := NewSettings("./test/invalid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	err = collector.Run()
	require.Error(t, err)

	status := collector.Status()
	require.False(t, status.Running)
	require.Error(t, status.Err)
	require.Contains(t, status.Err.Error(), "cannot build pipelines")
}

func TestCollectorRunInvalidFactory(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	settings.Factories.Extensions["invalid"] = extensionhelper.NewFactory(
		"invalid",
		defaultInvalidConfig,
		createInvalidExtension,
	)

	collector := New(settings)
	err = collector.Run()
	require.Error(t, err)

	status := collector.Status()
	require.False(t, status.Running)
	require.Contains(t, status.Err.Error(), "invalid config settings")
}

func TestCollectorRunTwice(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	err = collector.Run()
	require.NoError(t, err)
	defer collector.Stop()

	err = collector.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "service already running")
}

func TestCollectorStopTwice(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)

	collector := New(settings)
	err = collector.Run()
	require.NoError(t, err)
	collector.Stop()

	status := collector.Status()
	require.False(t, status.Running)

	collector.Stop()
	require.False(t, status.Running)
}

// InvalidConfig is a config without a mapstructure tag.
type InvalidConfig struct {
	config.ExtensionSettings `mapstructure:",squash"`
	FieldWithoutTag          string
}

// defaultInvalidConfig creates the default invalid config.
func defaultInvalidConfig() config.Extension {
	return &InvalidConfig{
		ExtensionSettings: config.NewExtensionSettings(config.NewID("invalid")),
	}
}

// createInvalidExtension always errors when creating an extension.
func createInvalidExtension(_ context.Context, _ component.ExtensionCreateSettings, _ config.Extension) (component.Extension, error) {
	return nil, errors.New("invalid extension")
}
