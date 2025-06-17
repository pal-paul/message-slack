# Test Configuration for Message Slack Action

## Environment Variables for Testing

When running tests locally or in CI/CD, you may need to set the following environment variables:

### Required for Integration Tests
- `SLACK_TOKEN`: A valid Slack bot token (xoxb-...) for testing actual Slack API calls
- `TEST_SLACK_CHANNEL`: The Slack channel to use for test messages (without #)

### Optional for Local Testing
- `GO_TEST_TIMEOUT`: Timeout for Go tests (default: 10m)
- `GO_TEST_VERBOSE`: Enable verbose test output (default: false)

## Running Tests

### Unit Tests Only
```bash
make test
# or
go test -v ./cmd/
```

### Integration Tests
```bash
go test -v -tags=integration ./cmd/
```

### All Tests with Coverage
```bash
go test -v -coverprofile=coverage.out ./cmd/
go tool cover -html=coverage.out -o coverage.html
```

### Benchmarks
```bash
go test -bench=. -benchmem ./cmd/
```

## Test Structure

### Unit Tests (`cmd_test.go`)
- `TestSlackMessageBuilder`: Tests the message formatting logic
- `TestEnvironmentVariablesParsing`: Tests environment variable validation
- `TestSlackMessageBuilderEdgeCases`: Tests edge cases and special characters
- `BenchmarkSlackMessageBuilder`: Performance benchmarks

### Integration Tests (`integration_test.go`)
- `TestMainIntegration`: Tests the complete application flow
- `TestMainTimeout`: Tests timeout behavior

### GitHub Actions Tests (`.github/workflows/test.yml`)
- Go unit and integration tests
- Action functionality tests with different scenarios
- Performance tests with multiple concurrent executions
- Edge case testing (empty inputs, long messages, Unicode)

## Mocking Strategy

The tests use environment variable mocking to simulate different GitHub Actions scenarios without requiring actual Slack API calls during unit testing.

For integration tests that require actual Slack API interaction, use a test workspace and dedicated test channel.

## Continuous Integration

The GitHub Actions workflow automatically runs:
1. Unit tests with coverage reporting
2. Integration tests (if Slack credentials are available)
3. Action functionality tests with various message formats
4. Performance tests with different load scenarios
5. Edge case validation

## Test Data

Test cases cover:
- Basic message formatting
- Markdown and rich text formatting
- Unicode characters and emojis
- Long messages and titles
- Empty or invalid inputs
- Various GitHub Actions contexts (success, failure, deployment)

## Security Considerations

- Never commit real Slack tokens to the repository
- Use GitHub Secrets for sensitive test data
- Test tokens should have minimal permissions and access to test channels only
- Integration tests should use `continue-on-error: true` to prevent CI failures when tokens are unavailable
