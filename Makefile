include Makefile.defs

# Define the name of the output binary
NAME = go-tinycache

# Default target, builds dependencies and the server
all: deps server

# Manage Go module dependencies
deps:
	go mod tidy

# Build the server binary
server:
	go build -v -o $(NAME) cmd/go-tinycache/main.go

# Clean up generated files
clean:
	rm -f $(NAME)
	rm -f tools/memcached_client

# Runs all Go unit tests
test:
	go test ./...
	rm -f bolt.db

# Runs the standalone memcached protocol/integration test
protocol_test:
	go run tools/memcached_command_review.go

# Builds the memcached client tool
tools:
	go build -o tools/memcached_client tools/memcached_command_review.go

# Declare phony targets
.PHONY: all deps server clean test protocol_test tools
