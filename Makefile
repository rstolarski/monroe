# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=monroe
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)

all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...

# Cross compilation
build_unix:
	$(GOBUILD) -o bin/$(BINARY_UNIX) -v
build_windows:
	$(GOBUILD) -o bin/$(BINARY_WINDOWS) -v