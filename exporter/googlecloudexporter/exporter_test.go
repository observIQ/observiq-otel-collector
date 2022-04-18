package googlecloudexporter

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
)

func TestExporterCapabilities(t *testing.T) {
	exporter := &Exporter{}
	capabilities := exporter.Capabilities()
	assert.True(t, capabilities.MutatesData)
}

func TestExporterWithConsumers(t *testing.T) {
	consumer := &MockProcessor{}
	consumer.On("ConsumeLogs", mock.Anything, mock.Anything).Return(nil).Once()
	consumer.On("ConsumeMetrics", mock.Anything, mock.Anything).Return(nil).Once()
	consumer.On("ConsumeTraces", mock.Anything, mock.Anything).Return(nil).Once()
	exporter := &Exporter{
		metricsConsumer: consumer,
		logsConsumer:    consumer,
		tracesConsumer:  consumer,
	}

	ctx := context.Background()
	err := exporter.ConsumeLogs(ctx, pdata.NewLogs())
	assert.Nil(t, err)

	err = exporter.ConsumeMetrics(ctx, pdata.NewMetrics())
	assert.Nil(t, err)

	err = exporter.ConsumeTraces(ctx, pdata.NewTraces())
	assert.Nil(t, err)

	consumer.AssertExpectations(t)
}

func TestExporterWithoutConsumers(t *testing.T) {
	exporter := &Exporter{}

	ctx := context.Background()
	err := exporter.ConsumeLogs(ctx, pdata.NewLogs())
	assert.Nil(t, err)

	err = exporter.ConsumeMetrics(ctx, pdata.NewMetrics())
	assert.Nil(t, err)

	err = exporter.ConsumeTraces(ctx, pdata.NewTraces())
	assert.Nil(t, err)
}

func TestExporterStart(t *testing.T) {
	testCases := []struct {
		name          string
		exporter      *Exporter
		expectedError error
	}{
		{
			name: "Successful metrics",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{createValidProcessor()},
				metricsExporter:   createValidExporter(),
			},
		},
		{
			name: "Successful traces",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{createValidProcessor()},
				tracesExporter:   createValidExporter(),
			},
		},
		{
			name: "Successful logs",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{createValidProcessor()},
				logsExporter:   createValidExporter(),
			},
		},
		{
			name: "Failing metrics processor",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				metricsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to start metrics processor"),
		},
		{
			name: "Failing traces processor",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				tracesExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to start traces processor"),
		},
		{
			name: "Failing logs processor",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				logsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to start logs processor"),
		},
		{
			name: "Failing metrics processor",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				metricsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to start metrics exporter"),
		},
		{
			name: "Failing traces processor",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				tracesExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to start traces exporter"),
		},
		{
			name: "Failing logs processor",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				logsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to start logs exporter"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.exporter.Start(context.Background(), nil)

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestExporterShutdown(t *testing.T) {
	testCases := []struct {
		name          string
		exporter      *Exporter
		expectedError error
	}{
		{
			name: "Successful metrics",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{createValidProcessor()},
				metricsExporter:   createValidExporter(),
			},
		},
		{
			name: "Successful traces",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{createValidProcessor()},
				tracesExporter:   createValidExporter(),
			},
		},
		{
			name: "Successful logs",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{createValidProcessor()},
				logsExporter:   createValidExporter(),
			},
		},
		{
			name: "Failing metrics processor",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				metricsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to shutdown metrics processor"),
		},
		{
			name: "Failing traces processor",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				tracesExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to shutdown traces processor"),
		},
		{
			name: "Failing logs processor",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{
					createValidProcessor(),
					createFailingProcessor(),
				},
				logsExporter: createValidExporter(),
			},
			expectedError: errors.New("failed to shutdown logs processor"),
		},
		{
			name: "Failing metrics processor",
			exporter: &Exporter{
				metricsProcessors: []component.MetricsProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				metricsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to shutdown metrics exporter"),
		},
		{
			name: "Failing traces processor",
			exporter: &Exporter{
				tracesProcessors: []component.TracesProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				tracesExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to shutdown traces exporter"),
		},
		{
			name: "Failing logs processor",
			exporter: &Exporter{
				logsProcessors: []component.LogsProcessor{
					createValidProcessor(),
					createValidProcessor(),
				},
				logsExporter: createFailingExporter(),
			},
			expectedError: errors.New("failed to shutdown logs exporter"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.exporter.Shutdown(context.Background())

			if tc.expectedError != nil {
				assert.Error(t, tc.expectedError, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func createValidProcessor() *MockProcessor {
	processor := &MockProcessor{}
	processor.On("Start", mock.Anything, mock.Anything).Return(nil)
	processor.On("Shutdown", mock.Anything).Return(nil)
	return processor
}

func createFailingProcessor() *MockProcessor {
	processor := &MockProcessor{}
	processor.On("Start", mock.Anything, mock.Anything).Return(errors.New("failure"))
	processor.On("Shutdown", mock.Anything).Return(errors.New("failure"))
	return processor
}

func createValidExporter() *MockExporter {
	exporter := &MockExporter{}
	exporter.On("Start", mock.Anything, mock.Anything).Return(nil)
	exporter.On("Shutdown", mock.Anything).Return(nil)
	return exporter
}

func createFailingExporter() *MockExporter {
	exporter := &MockExporter{}
	exporter.On("Start", mock.Anything, mock.Anything).Return(errors.New("failure"))
	exporter.On("Shutdown", mock.Anything).Return(errors.New("failure"))
	return exporter
}

// MockProcessor is an autogenerated mock type for the Processor type
type MockProcessor struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *MockProcessor) Capabilities() consumer.Capabilities {
	ret := _m.Called()

	var r0 consumer.Capabilities
	if rf, ok := ret.Get(0).(func() consumer.Capabilities); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(consumer.Capabilities)
	}

	return r0
}

// ConsumeLogs provides a mock function with given fields: ctx, ld
func (_m *MockProcessor) ConsumeLogs(ctx context.Context, ld pdata.Logs) error {
	ret := _m.Called(ctx, ld)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Logs) error); ok {
		r0 = rf(ctx, ld)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeMetrics provides a mock function with given fields: ctx, md
func (_m *MockProcessor) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	ret := _m.Called(ctx, md)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Metrics) error); ok {
		r0 = rf(ctx, md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeTraces provides a mock function with given fields: ctx, td
func (_m *MockProcessor) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	ret := _m.Called(ctx, td)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Traces) error); ok {
		r0 = rf(ctx, td)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *MockProcessor) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: ctx, host
func (_m *MockProcessor) Start(ctx context.Context, host component.Host) error {
	ret := _m.Called(ctx, host)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, component.Host) error); ok {
		r0 = rf(ctx, host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockExporter is an autogenerated mock type for the Exporter type
type MockExporter struct {
	mock.Mock
}

// Capabilities provides a mock function with given fields:
func (_m *MockExporter) Capabilities() consumer.Capabilities {
	ret := _m.Called()

	var r0 consumer.Capabilities
	if rf, ok := ret.Get(0).(func() consumer.Capabilities); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(consumer.Capabilities)
	}

	return r0
}

// ConsumeLogs provides a mock function with given fields: ctx, ld
func (_m *MockExporter) ConsumeLogs(ctx context.Context, ld pdata.Logs) error {
	ret := _m.Called(ctx, ld)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Logs) error); ok {
		r0 = rf(ctx, ld)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeMetrics provides a mock function with given fields: ctx, md
func (_m *MockExporter) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	ret := _m.Called(ctx, md)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Metrics) error); ok {
		r0 = rf(ctx, md)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumeTraces provides a mock function with given fields: ctx, td
func (_m *MockExporter) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	ret := _m.Called(ctx, td)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pdata.Traces) error); ok {
		r0 = rf(ctx, td)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *MockExporter) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: ctx, host
func (_m *MockExporter) Start(ctx context.Context, host component.Host) error {
	ret := _m.Called(ctx, host)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, component.Host) error); ok {
		r0 = rf(ctx, host)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}