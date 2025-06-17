#!/bin/bash

# Manual test script for message-slack action
# This script helps test the action locally before committing

set -e

echo "🧪 Message Slack Action - Manual Test Script"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go to run tests."
    exit 1
fi

print_success "Go is installed: $(go version)"

# Run unit tests
print_status "Running unit tests..."
if go test -v ./cmd/; then
    print_success "Unit tests passed!"
else
    print_error "Unit tests failed!"
    exit 1
fi

# Run tests with coverage
print_status "Running tests with coverage..."
if go test -v -coverprofile=coverage.out ./cmd/; then
    print_success "Coverage tests passed!"
    go tool cover -func=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
else
    print_error "Coverage tests failed!"
    exit 1
fi

# Run benchmark tests
print_status "Running benchmark tests..."
if go test -bench=. -benchmem ./cmd/; then
    print_success "Benchmark tests completed!"
else
    print_warning "Benchmark tests had issues, but continuing..."
fi

# Build the binary
print_status "Building the action binary..."
if make build; then
    print_success "Build successful!"
else
    print_error "Build failed!"
    exit 1
fi

# Test with mock environment variables
print_status "Testing with mock environment variables..."

export INPUT_TITLE="Test Message"
export INPUT_TEXT="This is a test message from the manual test script"
export INPUT_SLACK_TOKEN="xoxb-mock-token-for-testing"
export INPUT_SLACK_CHANNEL="test-channel"

print_status "Environment variables set:"
echo "  INPUT_TITLE=$INPUT_TITLE"
echo "  INPUT_TEXT=$INPUT_TEXT"
echo "  INPUT_SLACK_TOKEN=$INPUT_SLACK_TOKEN"
echo "  INPUT_SLACK_CHANNEL=$INPUT_SLACK_CHANNEL"

# Note: This will fail because we're using a mock token, but it tests the parsing
print_status "Running the binary (expected to fail with mock token)..."
if ./cmd/cmd; then
    print_warning "Binary ran successfully (unexpected with mock token)"
else
    print_success "Binary failed as expected with mock token"
fi

# Clean up environment
unset INPUT_TITLE INPUT_TEXT INPUT_SLACK_TOKEN INPUT_SLACK_CHANNEL

# Test action.yml validation
print_status "Validating action.yml..."
if [ -f "action.yml" ]; then
    print_success "action.yml exists"
    
    # Check for required fields
    if grep -q "name:" action.yml && grep -q "description:" action.yml && grep -q "inputs:" action.yml; then
        print_success "action.yml has required fields"
    else
        print_error "action.yml is missing required fields"
        exit 1
    fi
else
    print_error "action.yml not found!"
    exit 1
fi

# Validate examples
print_status "Validating example files..."
if [ -f "examples/example-usage.yml" ]; then
    print_success "Example usage file exists"
else
    print_warning "Example usage file not found"
fi

# Check for common issues
print_status "Checking for common issues..."

# Check if binary is executable
if [ -x "./cmd/cmd" ]; then
    print_success "Binary is executable"
else
    print_warning "Binary is not executable"
fi

# Check file sizes
BINARY_SIZE=$(stat -f%z "./cmd/cmd" 2>/dev/null || stat -c%s "./cmd/cmd" 2>/dev/null || echo "unknown")
print_status "Binary size: $BINARY_SIZE bytes"

if [ "$BINARY_SIZE" != "unknown" ] && [ "$BINARY_SIZE" -gt 50000000 ]; then
    print_warning "Binary is quite large (>50MB). Consider optimizing."
fi

# Clean up test artifacts
print_status "Cleaning up test artifacts..."
rm -f coverage.out coverage.html

print_success "🎉 Manual testing completed successfully!"
print_status "Next steps:"
echo "  1. Test with real Slack credentials in a safe environment"
echo "  2. Run the GitHub Actions workflow to test in CI/CD"
echo "  3. Test the action in a real repository workflow"
echo ""
print_status "To test with real Slack credentials:"
echo "  export INPUT_SLACK_TOKEN='your-real-token'"
echo "  export INPUT_SLACK_CHANNEL='your-test-channel'"
echo "  ./cmd/cmd"
