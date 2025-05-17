.PHONY: build test clean

BINARY_NAME := schola-theologiae-api

build:
	@# Move the tar.gz from ../books to ./data
	cp ../books/*.tar.gz ./data
	@# Uncompress the tar.gz files
	tar -xzf ./data/*.tar.gz --directory .
	@# Remove the tar.gz files
	rm --force ./data/*.tar.gz
	@# Create the bin directory
	mkdir --parents bin
	@# Install the dependencies
	go mod tidy
	@# Build the binary
	@# GOARCH=arm64 - for Apple Silicon or ARM64 archtectures
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/schola-theologiae-api main.go
	@# Add execute permissions
	chmod +x bin/$(BINARY_NAME)

test:
	@# Run the tests
	go test -v ./...
	@# Run the vet tool
	go vet ./...

clean:
	@# Remove the bin directory
	rm --recursive --force ./bin/*
	find ./data/* ! -name '*.go' | xargs rm --recursive --force 
