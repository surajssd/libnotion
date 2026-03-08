.PHONY: build test test-verbose test-cover test-coverprofile clean

# Build all packages
build:
	go build ./...

# Run all tests
test:
	go test ./...

# Run all tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage summary
test-cover:
	go test -cover ./...

# Run tests and generate coverage profile
test-coverprofile:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Open coverage report in browser
test-coverhtml: test-coverprofile
	go tool cover -html=coverage.out

# Clean generated files
clean:
	rm -f coverage.out
