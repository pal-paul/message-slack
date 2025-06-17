package main

import (
	"os"
	"strings"
)

// TestMode indicates if the application is running in test mode
var TestMode = false

// MockSlackResponse stores mock responses for testing
type MockSlackResponse struct {
	Success bool
	Error   error
	Message string
}

// MockSlackClient is a mock implementation of the Slack client for testing
type MockSlackClient struct {
	Responses   []MockSlackResponse
	CallCount   int
	LastMessage interface{}
}

// AddFormattedMessage mocks the Slack API call
func (m *MockSlackClient) AddFormattedMessage(channel string, message interface{}) (interface{}, error) {
	m.CallCount++
	m.LastMessage = message

	if len(m.Responses) > 0 {
		response := m.Responses[0]
		if len(m.Responses) > 1 {
			m.Responses = m.Responses[1:]
		}

		if response.Error != nil {
			return nil, response.Error
		}

		return map[string]interface{}{
			"ok":      response.Success,
			"message": response.Message,
		}, nil
	}

	// Default successful response
	return map[string]interface{}{
		"ok":      true,
		"message": "Message sent successfully",
	}, nil
}

// EnableTestMode enables test mode for the application
func EnableTestMode() {
	TestMode = true
}

// DisableTestMode disables test mode for the application
func DisableTestMode() {
	TestMode = false
}

// IsTestMode returns true if the application is running in test mode
func IsTestMode() bool {
	// Check explicit test mode flag first
	if TestMode {
		return true
	}

	// Check for test environment variables
	if os.Getenv("GO_TEST_MODE") == "true" || os.Getenv("TESTING") == "true" {
		return true
	}

	// Check if we're running from a test binary
	if len(os.Args) > 0 {
		executable := os.Args[0]
		if strings.Contains(executable, "test") || strings.HasSuffix(executable, ".test") {
			return true
		}
	}

	return false
}

// SetupTestEnvironment sets up a test environment with mock values
func SetupTestEnvironment() {
	os.Setenv("INPUT_TITLE", "Test Title")
	os.Setenv("INPUT_TEXT", "Test message content")
	os.Setenv("INPUT_SLACK_TOKEN", "xoxb-test-token")
	os.Setenv("INPUT_SLACK_CHANNEL", "test-channel")
	os.Setenv("GO_TEST_MODE", "true")
	EnableTestMode()
}

// CleanupTestEnvironment cleans up the test environment
func CleanupTestEnvironment() {
	testEnvVars := []string{
		"INPUT_TITLE",
		"INPUT_TEXT",
		"INPUT_SLACK_TOKEN",
		"INPUT_SLACK_CHANNEL",
		"GO_TEST_MODE",
		"TESTING",
	}

	for _, envVar := range testEnvVars {
		os.Unsetenv(envVar)
	}

	DisableTestMode()
}

// ValidateSlackToken validates the format of a Slack token
func ValidateSlackToken(token string) bool {
	if token == "" {
		return false
	}

	// Check for bot tokens (xoxb-) or user tokens (xoxp-)
	if strings.HasPrefix(token, "xoxb-") || strings.HasPrefix(token, "xoxp-") {
		return len(token) > 10 // Basic length check
	}

	return false
}

// ValidateSlackChannel validates the format of a Slack channel name
func ValidateSlackChannel(channel string) bool {
	if channel == "" {
		return false
	}

	// Channel names should not start with # (we handle that internally)
	if strings.HasPrefix(channel, "#") {
		return false
	}

	// Basic validation: alphanumeric, hyphens, underscores
	for _, char := range channel {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}

	return len(channel) >= 1 && len(channel) <= 80
}

// GetTestSlackClient returns a mock Slack client for testing
func GetTestSlackClient() *MockSlackClient {
	return &MockSlackClient{
		Responses: []MockSlackResponse{
			{Success: true, Error: nil, Message: "Test message sent"},
		},
		CallCount: 0,
	}
}
