# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -v
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
BINARY_NAME=imgurfetch
ROUTER_BINARY_NAME=kvrouter

.PHONY: all
all: test build

.PHONY: build
build:
	@mkdir -p $(BINARY_FOLDER)
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) ./cmd

.PHONY: test
test:
	$(GOTEST) -v -race ./...

.PHONY: clean
store-clean:
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(STORE_BINARY_NAME)
