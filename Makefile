BINARY_NAME=finflow
PACKAGE=github.com/josiastomasnanez/finflow

.PHONY: all build test lint fmt sonar frontend-install frontend-build run docker

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/finflow

test:
	go test ./...

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

frontend-install:
	cd frontend && npm install

frontend-build: frontend-install
	cd frontend && npm run build

run: frontend-build build
	./$(BINARY_NAME)

sonar:
	sonar-scanner

docker:
	docker build -t finflow:latest .
