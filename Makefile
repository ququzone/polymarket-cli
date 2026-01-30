.PHONY: build run clean test install

APP_NAME=polymarket-cli
VERSION?=dev
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

build:
	go build $(LDFLAGS) -o bin/$(APP_NAME)

run:
	go run main.go

clean:
	rm -rf bin/

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

install:
	go install $(LDFLAGS)

deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

fmt:
	go fmt ./...
