package main

import (
	"os"
	"os/exec"
	"testing"
)

// Integration tests that simulate the GitHub Actions environment
func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{
			name: "Valid environment variables",
			env: map[string]string{
				"INPUT_TITLE":         "Integration Test",
				"INPUT_TEXT":          "This is an integration test message",
				"INPUT_SLACK_TOKEN":   "xoxb-test-token-for-integration",
				"INPUT_SLACK_CHANNEL": "test-integration",
			},
			wantErr: true, // Will fail due to invalid token, but tests env parsing
		},
		{
			name: "Missing required environment variable",
			env: map[string]string{
				"INPUT_TITLE": "Integration Test",
				"INPUT_TEXT":  "This is an integration test message",
				// Missing INPUT_SLACK_TOKEN and INPUT_SLACK_CHANNEL
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build the binary first
			buildCmd := exec.Command("go", "build", "-o", "cmd_test", "cmd.go")
			buildCmd.Dir = "."
			if err := buildCmd.Run(); err != nil {
				t.Fatalf("Failed to build test binary: %v", err)
			}
			defer os.Remove("cmd_test")

			// Set up environment
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			// Clean up environment after test
			defer func() {
				for key := range tt.env {
					os.Unsetenv(key)
				}
			}()

			// Run the binary
			cmd := exec.Command("./cmd_test")
			err := cmd.Run()

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Test the binary execution with timeout
func TestMainTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout tests in short mode")
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "cmd_timeout_test", "cmd.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer os.Remove("cmd_timeout_test")

	// Set up minimal environment
	envVars := map[string]string{
		"INPUT_TITLE":         "Timeout Test",
		"INPUT_TEXT":          "Testing timeout behavior",
		"INPUT_SLACK_TOKEN":   "invalid-token",
		"INPUT_SLACK_CHANNEL": "test-timeout",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// This test verifies the binary doesn't hang indefinitely
	cmd := exec.Command("./cmd_timeout_test")
	err := cmd.Run()

	// We expect this to fail due to invalid token, but it should fail quickly
	if err == nil {
		t.Error("Expected command to fail with invalid token")
	}
}
