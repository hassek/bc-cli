# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X github.com/hassek/bc-cli/api.Version=$(VERSION) \
           -X github.com/hassek/bc-cli/cmd.Version=$(VERSION) \
           -X github.com/hassek/bc-cli/cmd.GitCommit=$(GIT_COMMIT) \
           -X github.com/hassek/bc-cli/cmd.BuildDate=$(BUILD_DATE)

compile:
	go build -ldflags "$(LDFLAGS)" -o bc-cli .

install:
	pip install pip -U
	pip install -r requirements-dev.txt
	go install golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest
	pre-commit install
	pre-commit autoupdate
	pre-commit install-hooks
	pre-commit install --hook-type commit-msg

upgrade:
	go get -u ./...
	go mod tidy
	pip-upgrade

.PHONY:
.ONESHELL:
