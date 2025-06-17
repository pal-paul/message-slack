package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	slack "github.com/pal-paul/go-libraries/pkg/slack"
)

func TestSlackMessageBuilder(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		text     string
		channel  string
		expected struct {
			blocksCount int
			headerText  string
			sectionText string
		}
	}{
		{
			name:    "Basic message with title and text",
			title:   "Test Title",
			text:    "Test message content",
			channel: "test-channel",
			expected: struct {
				blocksCount int
				headerText  string
				sectionText string
			}{
				blocksCount: 2,
				headerText:  "Test Title",
				sectionText: "Test message content",
			},
		},
		{
			name:    "Message with markdown text",
			title:   "Deployment Status",
			text:    "*Bold text* and _italic text_ with [link](https://example.com)",
			channel: "deployments",
			expected: struct {
				blocksCount int
				headerText  string
				sectionText string
			}{
				blocksCount: 2,
				headerText:  "Deployment Status",
				sectionText: "*Bold text* and _italic text_ with [link](https://example.com)",
			},
		},
		{
			name:    "Message with special characters",
			title:   "Build Failed ❌",
			text:    "Error: Cannot find module 'express'\nAt line 5:12",
			channel: "alerts",
			expected: struct {
				blocksCount int
				headerText  string
				sectionText string
			}{
				blocksCount: 2,
				headerText:  "Build Failed ❌",
				sectionText: "Error: Cannot find module 'express'\nAt line 5:12",
			},
		},
		{
			name:    "Empty strings",
			title:   "",
			text:    "",
			channel: "test",
			expected: struct {
				blocksCount int
				headerText  string
				sectionText string
			}{
				blocksCount: 2,
				headerText:  "",
				sectionText: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := SlackMessageBuilder(tt.title, tt.text, tt.channel)

			// Verify channel is set correctly
			if message.Channel != tt.channel {
				t.Errorf("Expected channel %s, got %s", tt.channel, message.Channel)
			}

			// Verify number of blocks
			if len(message.Blocks) != tt.expected.blocksCount {
				t.Errorf("Expected %d blocks, got %d", tt.expected.blocksCount, len(message.Blocks))
			}

			// Verify header block
			if len(message.Blocks) > 0 {
				headerBlock := message.Blocks[0]
				if headerBlock.Type != slack.HeaderBlock {
					t.Errorf("Expected first block to be HeaderBlock, got %s", headerBlock.Type)
				}
				if headerBlock.Text == nil {
					t.Error("Expected header block to have text")
				} else {
					if headerBlock.Text.Type != slack.PlainText {
						t.Errorf("Expected header text type to be PlainText, got %s", headerBlock.Text.Type)
					}
					if headerBlock.Text.Text != tt.expected.headerText {
						t.Errorf("Expected header text '%s', got '%s'", tt.expected.headerText, headerBlock.Text.Text)
					}
				}
			}

			// Verify section block
			if len(message.Blocks) > 1 {
				sectionBlock := message.Blocks[1]
				if sectionBlock.Type != slack.SectionBlock {
					t.Errorf("Expected second block to be SectionBlock, got %s", sectionBlock.Type)
				}
				if sectionBlock.Text == nil {
					t.Error("Expected section block to have text")
				} else {
					if sectionBlock.Text.Type != slack.Mrkdwn {
						t.Errorf("Expected section text type to be Mrkdwn, got %s", sectionBlock.Text.Type)
					}
					if sectionBlock.Text.Text != tt.expected.sectionText {
						t.Errorf("Expected section text '%s', got '%s'", tt.expected.sectionText, sectionBlock.Text.Text)
					}
				}
			}
		})
	}
}

func TestInitializeApp(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantErr  bool
		errorMsg string
	}{
		{
			name: "All required variables present",
			envVars: map[string]string{
				"INPUT_TITLE":         "Test Title",
				"INPUT_TEXT":          "Test Text",
				"INPUT_SLACK_TOKEN":   "xoxb-test-token",
				"INPUT_SLACK_CHANNEL": "test-channel",
			},
			wantErr: false,
		},
		{
			name: "Missing title",
			envVars: map[string]string{
				"INPUT_TEXT":          "Test Text",
				"INPUT_SLACK_TOKEN":   "xoxb-test-token",
				"INPUT_SLACK_CHANNEL": "test-channel",
			},
			wantErr: true,
		},
		{
			name: "Missing text",
			envVars: map[string]string{
				"INPUT_TITLE":         "Test Title",
				"INPUT_SLACK_TOKEN":   "xoxb-test-token",
				"INPUT_SLACK_CHANNEL": "test-channel",
			},
			wantErr: true,
		},
		{
			name: "Missing slack token",
			envVars: map[string]string{
				"INPUT_TITLE":         "Test Title",
				"INPUT_TEXT":          "Test Text",
				"INPUT_SLACK_CHANNEL": "test-channel",
			},
			wantErr: true,
		},
		{
			name: "Missing slack channel",
			envVars: map[string]string{
				"INPUT_TITLE":       "Test Title",
				"INPUT_TEXT":        "Test Text",
				"INPUT_SLACK_TOKEN": "xoxb-test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables
			clearEnvVars := []string{
				"INPUT_TITLE",
				"INPUT_TEXT",
				"INPUT_SLACK_TOKEN",
				"INPUT_SLACK_CHANNEL",
			}

			for _, envVar := range clearEnvVars {
				os.Unsetenv(envVar)
			}

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				for _, envVar := range clearEnvVars {
					os.Unsetenv(envVar)
				}
			}()

			// Test initialization
			err := initializeApp()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If no error expected, verify environment was parsed correctly
			if !tt.wantErr && err == nil {
				if envVar.Input.Title != tt.envVars["INPUT_TITLE"] {
					t.Errorf("Expected title %s, got %s", tt.envVars["INPUT_TITLE"], envVar.Input.Title)
				}
				if envVar.Input.Text != tt.envVars["INPUT_TEXT"] {
					t.Errorf("Expected text %s, got %s", tt.envVars["INPUT_TEXT"], envVar.Input.Text)
				}
				if envVar.Slack.Token != tt.envVars["INPUT_SLACK_TOKEN"] {
					t.Errorf("Expected token %s, got %s", tt.envVars["INPUT_SLACK_TOKEN"], envVar.Slack.Token)
				}
				if envVar.Slack.Channel != tt.envVars["INPUT_SLACK_CHANNEL"] {
					t.Errorf("Expected channel %s, got %s", tt.envVars["INPUT_SLACK_CHANNEL"], envVar.Slack.Channel)
				}
			}
		})
	}
}

func TestSlackMessageBuilderEdgeCases(t *testing.T) {
	// Set up environment
	envVar.Slack.Channel = "test-channel"

	t.Run("Very long title", func(t *testing.T) {
		longTitle := "This is a very long title that exceeds normal limits and might cause issues with Slack's API if not handled properly. It contains multiple sentences and goes on for quite a while to test edge cases."
		message := SlackMessageBuilder(longTitle, "Short text", "test-channel")

		if len(message.Blocks) != 2 {
			t.Errorf("Expected 2 blocks, got %d", len(message.Blocks))
		}

		if message.Blocks[0].Text.Text != longTitle {
			t.Error("Long title was not preserved correctly")
		}
	})

	t.Run("Very long text", func(t *testing.T) {
		longText := "This is a very long message that contains multiple paragraphs and exceeds normal message length limits. " +
			"It should still be handled correctly by the message builder function. " +
			"This tests whether the function can handle large amounts of text without issues. " +
			"The text continues with more content to ensure we're testing realistic scenarios where users might send " +
			"detailed deployment information, error logs, or comprehensive status updates through the Slack notification system."

		message := SlackMessageBuilder("Title", longText, "test-channel")

		if len(message.Blocks) != 2 {
			t.Errorf("Expected 2 blocks, got %d", len(message.Blocks))
		}

		if message.Blocks[1].Text.Text != longText {
			t.Error("Long text was not preserved correctly")
		}
	})

	t.Run("Unicode and emoji handling", func(t *testing.T) {
		title := "🚀 Deployment Status 🎉"
		text := "Successfully deployed 🌟 with 100% uptime ✅\nGreat job team! 👏"

		message := SlackMessageBuilder(title, text, "test-channel")

		if message.Blocks[0].Text.Text != title {
			t.Error("Unicode characters in title were not preserved")
		}

		if message.Blocks[1].Text.Text != text {
			t.Error("Unicode characters in text were not preserved")
		}
	})
}

// Additional test cases using test mode functionality
func TestValidateSlackToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "Empty token",
			token:    "",
			expected: false,
		},
		{
			name:     "Invalid prefix",
			token:    "invalid-token",
			expected: false,
		},
		{
			name:     "Too short bot token",
			token:    "anything",
			expected: false,
		},
		{
			name:     "Too short user token",
			token:    "anything",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSlackToken(tt.token)
			if result != tt.expected {
				t.Errorf("ValidateSlackToken(%s) = %v, expected %v", tt.token, result, tt.expected)
			}
		})
	}
}

func TestValidateSlackChannel(t *testing.T) {
	tests := []struct {
		name     string
		channel  string
		expected bool
	}{
		{
			name:     "Valid channel name",
			channel:  "general",
			expected: true,
		},
		{
			name:     "Valid channel with hyphen",
			channel:  "team-updates",
			expected: true,
		},
		{
			name:     "Valid channel with underscore",
			channel:  "dev_team",
			expected: true,
		},
		{
			name:     "Valid channel with numbers",
			channel:  "channel123",
			expected: true,
		},
		{
			name:     "Empty channel",
			channel:  "",
			expected: false,
		},
		{
			name:     "Channel with hash prefix",
			channel:  "#general",
			expected: false,
		},
		{
			name:     "Channel with special characters",
			channel:  "team@updates",
			expected: false,
		},
		{
			name:     "Channel with spaces",
			channel:  "team updates",
			expected: false,
		},
		{
			name:     "Very long channel name",
			channel:  "a" + strings.Repeat("b", 80),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSlackChannel(tt.channel)
			if result != tt.expected {
				t.Errorf("ValidateSlackChannel(%s) = %v, expected %v", tt.channel, result, tt.expected)
			}
		})
	}
}

func TestMockSlackClient(t *testing.T) {
	mockClient := GetTestSlackClient()

	message := SlackMessageBuilder("Test Title", "Test message", "test-channel")

	// Test successful message sending
	response, err := mockClient.AddFormattedMessage("test-channel", message)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if mockClient.CallCount != 1 {
		t.Errorf("Expected call count 1, got %d", mockClient.CallCount)
	}

	if response == nil {
		t.Error("Expected response, got nil")
	}

	// Test multiple calls
	mockClient.AddFormattedMessage("test-channel", message)
	if mockClient.CallCount != 2 {
		t.Errorf("Expected call count 2, got %d", mockClient.CallCount)
	}
}

// Test error handling in mock Slack client
func TestMockSlackClientErrorHandling(t *testing.T) {
	// Test with error response
	mockClient := &MockSlackClient{
		Responses: []MockSlackResponse{
			{Success: false, Error: errors.New("network error"), Message: ""},
		},
	}

	message := SlackMessageBuilder("Test Title", "Test message", "test-channel")

	_, err := mockClient.AddFormattedMessage("test-channel", message)
	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "network error" {
		t.Errorf("Expected 'network error', got '%v'", err)
	}
}

// Test concurrent usage of SlackMessageBuilder
func TestSlackMessageBuilderConcurrent(t *testing.T) {
	const numGoroutines = 100
	const numMessages = 10

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numMessages)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				title := fmt.Sprintf("Concurrent Test %d-%d", id, j)
				text := fmt.Sprintf("Message from goroutine %d, iteration %d", id, j)
				channel := fmt.Sprintf("test-channel-%d", id)

				message := SlackMessageBuilder(title, text, channel)

				// Verify the message was built correctly
				if message.Channel != channel {
					errors <- fmt.Errorf("wrong channel: expected %s, got %s", channel, message.Channel)
					return
				}

				if len(message.Blocks) != 2 {
					errors <- fmt.Errorf("wrong number of blocks: expected 2, got %d", len(message.Blocks))
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("Concurrent test error: %v", err)
		}
	}
}

// Test input sanitization
func TestInputSanitization(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		text    string
		channel string
	}{
		{
			name:    "XSS attempt in title",
			title:   "<script>alert('xss')</script>",
			text:    "Normal text",
			channel: "test-channel",
		},
		{
			name:    "SQL injection attempt in text",
			title:   "Normal title",
			text:    "'; DROP TABLE users; --",
			channel: "test-channel",
		},
		{
			name:    "Null bytes in input",
			title:   "Title with null\x00byte",
			text:    "Text with null\x00byte",
			channel: "test-channel",
		},
		{
			name:    "Control characters",
			title:   "Title\r\n\t",
			text:    "Text\r\n\t",
			channel: "test-channel",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic or cause issues
			message := SlackMessageBuilder(tt.title, tt.text, tt.channel)

			// Verify structure is maintained
			if len(message.Blocks) != 2 {
				t.Errorf("Expected 2 blocks, got %d", len(message.Blocks))
			}

			// Verify data is preserved (not sanitized - up to Slack API to handle)
			if message.Blocks[0].Text.Text != tt.title {
				t.Errorf("Title not preserved: expected '%s', got '%s'", tt.title, message.Blocks[0].Text.Text)
			}

			if message.Blocks[1].Text.Text != tt.text {
				t.Errorf("Text not preserved: expected '%s', got '%s'", tt.text, message.Blocks[1].Text.Text)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSlackMessageBuilder(b *testing.B) {
	title := "Benchmark Test"
	text := "This is a benchmark test message"
	channel := "benchmark-channel"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SlackMessageBuilder(title, text, channel)
	}
}

func BenchmarkSlackMessageBuilderLarge(b *testing.B) {
	title := "Large Benchmark Test with a very long title that simulates real-world usage"
	text := "This is a large benchmark test message that contains multiple lines of text and simulates " +
		"real-world usage scenarios where users might send detailed information through Slack notifications. " +
		"This includes deployment details, error messages, status updates, and other comprehensive information " +
		"that teams typically share through their CI/CD pipelines and monitoring systems."
	channel := "benchmark-channel"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SlackMessageBuilder(title, text, channel)
	}
}

func TestTestModeHelpers(t *testing.T) {
	// Save original state
	originalTestMode := TestMode
	defer func() {
		TestMode = originalTestMode
		CleanupTestEnvironment()
	}()

	// Start with clean state - note that IsTestMode() will return true during tests
	// because we're running from a test binary, so we test the explicit flag control
	TestMode = false
	CleanupTestEnvironment()

	// Test enabling test mode explicitly
	EnableTestMode()
	if !TestMode {
		t.Error("Expected TestMode flag to be true after enabling")
	}

	// Test disabling test mode explicitly
	DisableTestMode()
	if TestMode {
		t.Error("Expected TestMode flag to be false after disabling")
	}

	// Test environment setup and cleanup
	SetupTestEnvironment()
	if !TestMode {
		t.Error("Expected TestMode flag to be true after setup")
	}

	// Verify environment variables were set
	if os.Getenv("INPUT_TITLE") != "Test Title" {
		t.Error("Expected INPUT_TITLE to be set")
	}

	if os.Getenv("GO_TEST_MODE") != "true" {
		t.Error("Expected GO_TEST_MODE to be set")
	}

	CleanupTestEnvironment()
	if TestMode {
		t.Error("Expected TestMode flag to be false after cleanup")
	}

	if os.Getenv("INPUT_TITLE") != "" {
		t.Error("Expected INPUT_TITLE to be cleaned up")
	}

	if os.Getenv("GO_TEST_MODE") != "" {
		t.Error("Expected GO_TEST_MODE to be cleaned up")
	}
}
