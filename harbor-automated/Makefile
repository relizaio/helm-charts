.PHONY: build clean setup test install help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=harbor-modifier
BINARY_PATH=bin/$(BINARY_NAME)

# Harbor parameters
HARBOR_VERSION?=1.18.0
NAMESPACE?=harbor

all: build

## build: Build the harbor-modifier binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY_PATH) ./cmd/harbor-modifier
	@echo "✅ Built: $(BINARY_PATH)"

## clean: Remove build artifacts and generated chart
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	rm -rf harbor-helm/
	@echo "✅ Cleaned"

## setup: Pull and modify Harbor chart
setup: build
	@echo "Setting up Harbor chart..."
	./$(BINARY_PATH) -version=$(HARBOR_VERSION)

## test: Run tests
test:
	$(GOTEST) -v ./...

## install: Install Harbor to Kubernetes
install: setup
	@echo "Installing Harbor to $(NAMESPACE)..."
	helm install harbor ./harbor-helm -n $(NAMESPACE) --create-namespace \
		-f examples/values-reliza-postgresql.yaml

## uninstall: Uninstall Harbor from Kubernetes
uninstall:
	@echo "Uninstalling Harbor from $(NAMESPACE)..."
	helm uninstall harbor -n $(NAMESPACE) || true
	kubectl delete pvc -n $(NAMESPACE) --all || true

## lint: Run helm lint
lint: setup
	helm lint ./harbor-helm

## template: Render templates
template: setup
	helm template test ./harbor-helm --set expose.type=clusterIP --set expose.tls.auto.commonName=harbor.local

## deps: Download Go dependencies
deps:
	$(GOGET) -v -t -d ./...
	go mod tidy

## help: Show this help message
help:
	@echo "Harbor Chart Automation - Makefile Commands"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
