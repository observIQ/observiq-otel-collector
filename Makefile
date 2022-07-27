# All source code and documents, used when checking for misspellings
ALLDOC := $(shell find . \( -name "*.md" -o -name "*.yaml" \) \
                                -type f | sort)
ALL_MODULES := $(shell find . -type f -name "go.mod" -exec dirname {} \; | sort )

# All source code files
ALL_SRC := $(shell find . -name '*.go' -o -name '*.sh' -o -name 'Dockerfile' -type f | sort)

OUTDIR=./dist
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

INTEGRATION_TEST_ARGS?=-tags integration

ifeq ($(GOOS), windows)
EXT?=.exe
else
EXT?=
endif

PREVIOUS_TAG := $(shell git tag --sort=v:refname --no-contains HEAD | grep -E "[0-9]+\.[0-9]+\.[0-9]+$$" | tail -n1)
CURRENT_TAG := $(shell git tag --sort=v:refname --points-at HEAD | grep -E "v[0-9]+\.[0-9]+\.[0-9]+$$" | tail -n1)
# Version will be the tag pointing to the current commit, or the previous version tag if there is no such tag
VERSION ?= $(if $(CURRENT_TAG),$(CURRENT_TAG),$(PREVIOUS_TAG))

# Build binaries for current GOOS/GOARCH by default
.DEFAULT_GOAL := build-binaries

# Builds just the collector for current GOOS/GOARCH pair
.PHONY: collector
collector:
	go build -ldflags "-s -w -X github.com/observiq/observiq-otel-collector/internal/version.version=$(VERSION)" -o $(OUTDIR)/collector_$(GOOS)_$(GOARCH)$(EXT) ./cmd/collector

# Builds just the updater for current GOOS/GOARCH pair
.PHONY: updater
updater:
	cd ./updater/; go build -ldflags "-s -w -X github.com/observiq/observiq-otel-collector/internal/version.version=$(VERSION)" -o ../$(OUTDIR)/updater_$(GOOS)_$(GOARCH)$(EXT) ./cmd/updater

# Builds the updater + collector for current GOOS/GOARCH pair
.PHONY: build-binaries
build-binaries: collector updater

.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux: build-linux-amd64 build-linux-arm64 build-linux-arm

.PHONY: build-darwin
build-darwin: build-darwin-amd64 build-darwin-arm64

.PHONY: build-windows
build-windows: build-windows-amd64

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(MAKE) build-binaries -j2

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(MAKE) build-binaries -j2

.PHONY: build-linux-arm
build-linux-arm:
	GOOS=linux GOARCH=arm $(MAKE) build-binaries -j2

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(MAKE) build-binaries -j2

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(MAKE) build-binaries -j2

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(MAKE) build-binaries -j2

# tool-related commands
.PHONY: install-tools
install-tools:
	go install github.com/mgechev/revive@v1.2.0
	go install github.com/google/addlicense@v1.0.0
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install github.com/sigstore/cosign/cmd/cosign@v1.5.2
	go install github.com/goreleaser/goreleaser@v1.9.1
	go install github.com/securego/gosec/v2/cmd/gosec@v2.10.0
	go install github.com/uw-labs/lichen@v0.1.5
	go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@v0.47.0
	

.PHONY: lint
lint:
	revive -config revive/config.toml -formatter friendly ./...

.PHONY: misspell
misspell:
	misspell $(ALLDOC)

.PHONY: misspell-fix
misspell-fix:
	misspell -w $(ALLDOC)

.PHONY: test
test:
	$(MAKE) for-all CMD="go test -race ./..."

.PHONY: test-with-cover
test-with-cover:
	$(MAKE) for-all CMD="go test -coverprofile=cover.out ./..."
	$(MAKE) for-all CMD="go tool cover -html=cover.out -o cover.html"

.PHONY: test-updater-integration 
test-updater-integration:
	cd updater; go test $(INTEGRATION_TEST_ARGS) -race ./...

.PHONY: bench
bench:
	$(MAKE) for-all CMD="go test -benchmem -run=^$$ -bench ^* ./..."

.PHONY: check-fmt
check-fmt:
	goimports -d ./ | diff -u /dev/null -

.PHONY: fmt
fmt:
	goimports -w .

.PHONY: tidy
tidy:
	$(MAKE) for-all CMD="go mod tidy -compat=1.17"

.PHONY: gosec
gosec:
	gosec -exclude-dir updater  ./...
# exclude the testdata dir; it contains a go program for testing.
	cd updater; gosec -exclude-dir internal/service/testdata ./...

# This target performs all checks that CI will do (excluding the build itself)
.PHONY: ci-checks
ci-checks: check-fmt check-license misspell lint gosec test

# This target checks that license copyright header is on every source file
.PHONY: check-license
check-license:
	@ADDLICENSEOUT=`addlicense -check $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENSEOUT" ]; then \
			echo "addlicense FAILED => add License errors:\n"; \
			echo "$$ADDLICENSEOUT\n"; \
			echo "Use 'make add-license' to fix this."; \
			exit 1; \
		else \
			echo "Check License finished successfully"; \
		fi

# This target adds a license copyright header is on every source file that is missing one
.PHONY: add-license
add-license:
	@ADDLICENSEOUT=`addlicense -y "" -c "observIQ, Inc." $(ALL_SRC) 2>&1`; \
		if [ "$$ADDLICENSEOUT" ]; then \
			echo "addlicense FAILED => add License errors:\n"; \
			echo "$$ADDLICENSEOUT\n"; \
			exit 1; \
		else \
			echo "Add License finished successfully"; \
		fi

# Downloads and setups dependencies that are packaged with binary
.PHONY: release-prep
release-prep:
	@rm -rf release_deps
	@mkdir release_deps
	@echo 'v$(CURR_VERSION)' > release_deps/VERSION.txt
	./buildscripts/download-dependencies.sh release_deps
	@cp -r ./plugins release_deps/
	@cp config/example.yaml release_deps/config.yaml
	@cp config/logging.yaml release_deps/logging.yaml
	@cp service/com.observiq.collector.plist release_deps/com.observiq.collector.plist
	@jq ".files[] | select(.service != null)" windows/wix.json >> release_deps/windows_service.json
	@cp service/observiq-otel-collector.service release_deps/observiq-otel-collector.service

# Build, sign, and release
.PHONY: release
release:
	goreleaser release --parallelism 4 --rm-dist

# Build and sign, skip release and ignore dirty git tree
.PHONY: release-test
release-test:
	goreleaser release --parallelism 4 --skip-validate --skip-publish --skip-sign --rm-dist --snapshot

.PHONY: for-all
for-all:
	@echo "running $${CMD} in root"
	@$${CMD}
	@set -e; for dir in $(ALL_MODULES); do \
	  (cd "$${dir}" && \
	  	echo "running $${CMD} in $${dir}" && \
	 	$${CMD} ); \
	done

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: scan-licenses
scan-licenses:
	lichen --config=./license.yaml $$(find dist/collector_* | grep -v 'sig\|json\|CHANGELOG.md\|yaml\|SHA256' | xargs)

.PHONY: generate
generate:
	$(MAKE) for-all CMD="go generate ./..."
