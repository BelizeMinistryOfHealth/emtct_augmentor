SHELL := /bin/bash

export PROJECT = emtct_augmentor

# ==============================================================================
# Modules support

tidy:
	go mod tidy

# ==============================================================================

# ==============================================================================
# Build
build-linux:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/emtct cmd/server/main.go

build-macos:
	export GO111MODULE=on
	env GOOS=darwin go build -o bin/emtct cmd/server/main.go

clean:
	rm -rf ./bin Gopkg.lock



