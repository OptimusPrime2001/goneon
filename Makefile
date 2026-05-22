.PHONY: all run build test fmt tidy clean

# Default target
all: tidy fmt test build

# Run the application
run:
	go run main.go

# Build the binary
build:
	mkdir -p bin
	go build -o bin/goneon main.go

# Run unit tests
test:
	go test -v -race ./...

# Format the codebase
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
