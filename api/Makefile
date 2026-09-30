.PHONY: build library index test clean

BINARY_NAME := schola-theologiae-api

# FTS5 is not enabled in go-sqlite3 by default, and the search endpoint relies on
# the full-text index built by tools/indexer, so both binaries need this tag.
SQLITE_TAGS := sqlite_fts5

# Default to the host architecture so local builds work on ARM machines too;
# the container build still gets GOOS=linux from the Makefile default.
GOOS ?= linux
GOARCH ?= $(shell go env GOARCH)

build: library
	@# Create the bin directory
	mkdir --parents bin
	@# Install the dependencies
	go mod tidy
	@# Build the binary
	@# GOARCH=arm64 - for Apple Silicon or ARM64 archtectures
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags $(SQLITE_TAGS) -o bin/$(BINARY_NAME) main.go
	@# Add execute permissions
	chmod +x bin/$(BINARY_NAME)

library: clean
	@# Create the library directory if it doesn't exist
	mkdir -p ./data/library
	@# Copy the tar.gz from ../books to ./data/library
	cp ../books/*.tar.gz ./data/library
	@# Uncompress the tar.gz files
	for file in ./data/library/*.tar.gz; do \
		tar xzf "$$file" --directory ./data/library; \
	done
	@# Remove the tar.gz files
	rm -f ./data/library/*.tar.gz
	@# Build the full-text index over chapter titles and bodies
	$(MAKE) index

index:
	@# Index chapter titles + bodies for the search endpoint (accent-insensitive)
	go run -tags $(SQLITE_TAGS) ./tools/indexer -library ./data/library

test:
	@# Run the tests
	go test -v ./...
	@# Run the vet tool
	go vet ./...

clean:
	@# Remove the bin directory
	rm -rf ./bin/*
	find ./data/library/ ! -name '.keep' ! -name 'library' | xargs rm -rf
