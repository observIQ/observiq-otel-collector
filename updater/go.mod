module github.com/observiq/observiq-otel-collector/updater

go 1.17

require (
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/observiq/observiq-otel-collector/packagestate v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opamp-go v0.2.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.21.0
	golang.org/x/sys v0.0.0-20220408201424-a24fb2fb8a0f
)

require (
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/observiq/observiq-otel-collector/packagestate => ../packagestate
