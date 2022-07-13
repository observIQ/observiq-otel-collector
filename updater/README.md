# observIQ Distro for OpenTelemetry Updater

The updater is a separate binary that runs as a separate process to update collector artifacts (including the collector itself) when managed by [BindPlane OP](https://github.com/observIQ/bindplane-op).

Because the updater edits service configurations, it needs elevated privileges to run (root on Linux + macOS, administrative privileges on Windows).
