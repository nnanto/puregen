.PHONY: build test clean install example

BUILD_FILE=./build/protoc-gen-puregen
# Build the plugin
build:
	go build -o $(BUILD_FILE) ./cmd/protoc-gen-puregen

# Install the plugin to GOPATH/bin
install: build
	go install ./cmd/protoc-gen-puregen

# Clean build artifacts
clean:
	rm -f protoc-gen-puregen
	rm -rf examples/generated

# Test with example proto file
example: build
	mkdir -p examples/generated
	$(BUILD_FILE) --help || true
	protoc --plugin=$(BUILD_FILE) \
		--puregen_out=examples/generated \
		--puregen_opt=language=all \
		-I examples/proto \
		examples/proto/*.proto

# Test specific languages
example-go: build
	mkdir -p examples/generated
	protoc --plugin=$(BUILD_FILE) \
		--puregen_out=examples/generated \
		--puregen_opt=language=go \
		-I examples/proto \
		examples/proto/*.proto

example-java: build
	mkdir -p examples/generated
	protoc --plugin=$(BUILD_FILE) \
		--puregen_out=examples/generated \
		--puregen_opt=language=java \
		-I examples/proto \
		examples/proto/*.proto

example-python: build
	mkdir -p examples/generated
	protoc --plugin=$(BUILD_FILE) \
		--puregen_out=examples/generated \
		--puregen_opt=language=python \
		-I examples/proto \
		examples/proto/*.proto

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Lint code
lint:
	golangci-lint run

# Show version
version: build
	$(BUILD_FILE) --version
