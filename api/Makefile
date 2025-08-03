.PHONY: build library test clean

BINARY_NAME := schola-theologiae-api

build: library
	@# Create the bin directory
	mkdir --parents bin
	@# Install the dependencies
	go mod tidy
	@# Build the binary
	@# GOARCH=arm64 - for Apple Silicon or ARM64 archtectures
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/schola-theologiae-api main.go
	@# Add execute permissions
	chmod +x bin/$(BINARY_NAME)

library:
	@# Create the library directory if it doesn't exist
	mkdir -p ./data/library
	@# Copy the tar.gz from ../books to ./data/library
	cp ../books/*.tar.gz ./data/library
	@# Uncompress the tar.gz files
	for file in ./data/library/*.tar.gz; do \
		tar xzf "$$file" --directory ./data/library; \
	done
	@# Remove the tar.gz files
	rm --force ./data/library/*.tar.gz

test:
	@# Run the tests
	go test -v ./...
	@# Run the vet tool
	go vet ./...

clean:
	@# Remove the bin directory
	rm --recursive --force ./bin/*
	find ./data/library/ ! -name '.keep' ! -name 'library' | xargs rm --recursive --force 
