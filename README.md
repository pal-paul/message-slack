# Message Slack GitHub Action

A GitHub Action to send custom messages to Slack channels with rich formatting support.

## Features

- 🚀 Send custom messages to any Slack channel
- 📝 Rich text formatting with title and message body
- 🎨 Beautiful block-based message layout
- 🔒 Secure token-based authentication
- ⚡ Fast and lightweight Go implementation

## Usage

### Basic Example

```yaml
- name: Send Slack Message
  uses: pal-paul/message-slack@v1.4.0
  with:
    title: "Deployment Successful"
    text: "The application has been successfully deployed to production!"
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: "general"
```

### Conditional Notifications

```yaml
- name: Notify on Failure
  if: failure()
  uses: pal-paul/message-slack@v1.4.0
  with:
    title: "Build Failed"
    text: "❌ The build process has failed. Please check the logs for more details."
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: "alerts"
```

## Inputs

| Parameter | Description | Required | Example |
|-----------|-------------|----------|---------|
| `title` | Title of the message (displayed as header) | ✅ | `"Deployment Status"` |
| `text` | Main content of the message (supports Markdown) | ✅ | `"Build completed successfully!"` |
| `slack_token` | Slack Bot Token (store in secrets) | ✅ | `${{ secrets.SLACK_TOKEN }}` |
| `slack_channel` | Slack channel name (without #) | ✅ | `"general"` |

## Setup

### 1. Create a Slack App

1. Go to [Slack API](https://api.slack.com/apps)
2. Click "Create New App" → "From scratch"
3. Give your app a name and select your workspace

### 2. Configure Bot Permissions

Add the following OAuth scopes under "OAuth & Permissions":

- `chat:write` - Send messages to channels
- `chat:write.public` - Send messages to public channels

### 3. Install App to Workspace

1. Click "Install to Workspace"
2. Copy the "Bot User OAuth Token" (starts with `xoxb-`)

### 4. Add Token to GitHub Secrets

1. Go to your repository → Settings → Secrets and variables → Actions
2. Create a new secret named `SLACK_TOKEN`
3. Paste your Bot User OAuth Token as the value

### 5. Invite Bot to Channel

In your Slack channel, type: `/invite @YourBotName`

## Advanced Examples

### Multi-line Message with Formatting

```yaml
- name: Deployment Summary
  uses: pal-paul/message-slack@v1.4.0
  with:
    title: "🚀 Production Deployment"
    text: |
      **Environment:** Production
      **Version:** ${{ github.sha }}
      **Branch:** ${{ github.ref_name }}
      **Status:** ✅ Success
      
      [View Workflow](${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }})
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: "deployments"
```

### Using Environment Variables

```yaml
- name: Send Environment Info
  uses: pal-paul/message-slack@v1.4.0
  with:
    title: "Environment Status"
    text: "Deployment to ${{ vars.ENVIRONMENT_NAME }} completed"
    slack_token: ${{ secrets.SLACK_TOKEN }}
    slack_channel: ${{ vars.SLACK_CHANNEL }}
```

## Message Format

The action creates beautiful Slack messages with:

- **Header Block**: Contains the title in bold
- **Section Block**: Contains the main text with Markdown support

Messages support standard Slack markdown formatting:

- `*bold*` for **bold text**
- `_italic_` for *italic text*
- `[link text](url)` for links
- `` `code` `` for inline code
- Multi-line text and lists

## Troubleshooting

### Common Issues

**❌ "channel_not_found"**

- Ensure the bot is invited to the channel
- Use channel name without the # symbol

**❌ "not_authed" or "invalid_auth"**

- Check your Slack token is correct
- Ensure the token has proper permissions

**❌ "missing_scope"**

- Add required OAuth scopes to your Slack app
- Reinstall the app to workspace

## Development

This action is built with Go and uses composite actions for fast execution.

### Local Development

```bash
# Clone the repository
git clone https://github.com/pal-paul/message-slack.git
cd message-slack

# Build the binary
make build

# Run tests
make test

# Run all tests (unit + integration)
make test-all

# Run tests with coverage
make test-coverage

# Run benchmark tests
make test-bench

# Run quick manual test
./test.sh
```

## Testing

This action includes comprehensive test coverage:

### Test Types


- **Unit Tests**: Core functionality testing
- **Integration Tests**: End-to-end workflow testing
- **Validation Tests**: Input validation and error handling
- **Performance Tests**: Benchmark and load testing
- **Concurrency Tests**: Thread-safety validation
- **Security Tests**: Input sanitization testing


### Coverage

- **65.2%** statement coverage
- All critical paths tested
- Edge cases and error scenarios covered


### Running Tests

```bash
# Quick test run
make test

# Full test suite with coverage
make test-coverage

# Performance benchmarks
make test-bench

# Manual testing script
./test.sh
```

## Contributing


Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines

1. Write tests for new functionality
2. Ensure all tests pass before submitting
3. Follow Go best practices
4. Update documentation as needed

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/pal-paul/message-slack/issues) on GitHub.
