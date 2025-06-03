# AI Terminal Assistant - Daemon Mode

The AI Terminal Assistant can run as a macOS background service (daemon), providing a persistent HTTP API for AI interactions without needing to run the terminal application.

## Overview

The daemon mode provides:
- **Background service**: Runs continuously as a macOS LaunchAgent
- **HTTP API**: RESTful endpoints for chat interactions
- **Automatic startup**: Starts on system boot
- **Service management**: Easy install/uninstall commands
- **Logging**: Centralized log files for monitoring

## Quick Start

### 1. Build and Install

```bash
# Build the daemon
make build-daemon

# Install as system service
make install

# Set up the service
ai-terminal-assistant -install

# Start the service
ai-terminal-assistant -start
```

### 2. Configure Environment

Create or edit `~/.ai-terminal-assistant.env`:

```bash
# Required: OpenAI API Key
OPENAI_API_KEY=your-openai-api-key-here

# Optional: Custom port (default: 8080)
PORT=8080

# Optional: Custom log file location
LOG_FILE=/usr/local/var/log/ai-terminal-assistant.log
```

### 3. Test the API

```bash
# Check service health
curl http://localhost:8080/health

# Send a chat message
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how are you?"}'
```

## Installation Methods

### Method 1: Using Makefile (Recommended)

```bash
# Build and install
make build-daemon install

# Set up service
ai-terminal-assistant -install
ai-terminal-assistant -start
```

### Method 2: Using Installation Script

```bash
# From release package
./install.sh

# Or build and run script
make build-daemon
cd build
../scripts/install.sh
```

### Method 3: Manual Installation

```bash
# Build binary
go build -o ai-terminal-assistant ./cmd/daemon

# Copy to system location
sudo cp ai-terminal-assistant /usr/local/bin/

# Install service
ai-terminal-assistant -install
```

## Service Management

### Start/Stop Commands

```bash
# Start the service
ai-terminal-assistant -start

# Stop the service
ai-terminal-assistant -stop

# Check status
ai-terminal-assistant -status

# Restart (stop then start)
ai-terminal-assistant -stop && ai-terminal-assistant -start
```

### Install/Uninstall Service

```bash
# Install LaunchAgent service
ai-terminal-assistant -install

# Uninstall service
ai-terminal-assistant -uninstall
```

### Service Information

- **Service Name**: `com.johnjallday.ai-terminal-assistant`
- **Plist Location**: `~/Library/LaunchAgents/com.johnjallday.ai-terminal-assistant.plist`
- **Default Port**: 8080
- **Log File**: `/usr/local/var/log/ai-terminal-assistant.log`
- **PID File**: `/usr/local/var/run/ai-terminal-assistant.pid`

## HTTP API Reference

### Health Check

```bash
GET /health
```

Response:
```json
{
  "status": "ok",
  "message": "AI Terminal Assistant daemon is running"
}
```

### Chat Endpoint

```bash
POST /chat
Content-Type: application/json

{
  "message": "Your message here",
  "model": "gpt-4.1-nano"  // optional, defaults to gpt-4.1-nano
}
```

Response:
```json
{
  "response": "AI response here",
  "model": "gpt-4.1-nano",
  "timestamp": "2025-06-02T15:32:02Z"
}
```

### Example API Usage

```bash
# Simple chat
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "What is the weather like?"}'

# With specific model
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Explain quantum computing", "model": "gpt-4"}'

# Health check
curl http://localhost:8080/health
```

## Configuration

### Environment Variables

The daemon supports configuration via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `OPENAI_API_KEY` | OpenAI API key (required) | - |
| `PORT` | Server port | 8080 |
| `LOG_FILE` | Log file path | `/usr/local/var/log/ai-terminal-assistant.log` |

### Configuration File

Create `~/.ai-terminal-assistant.env`:

```bash
# AI Terminal Assistant Environment Configuration

# Required: OpenAI API Key
OPENAI_API_KEY=sk-your-key-here

# Optional: Custom port
PORT=8080

# Optional: Custom log file
LOG_FILE=/usr/local/var/log/ai-terminal-assistant.log
```

## Logging

### Log Locations

- **Default log file**: `/usr/local/var/log/ai-terminal-assistant.log`
- **System logs**: Check with `log show --predicate 'subsystem == "com.johnjallday.ai-terminal-assistant"'`

### Viewing Logs

```bash
# Tail the log file
tail -f /usr/local/var/log/ai-terminal-assistant.log

# View recent logs
tail -n 100 /usr/local/var/log/ai-terminal-assistant.log

# System logs
log show --predicate 'subsystem == "com.johnjallday.ai-terminal-assistant"' --last 1h
```

## Troubleshooting

### Service Won't Start

1. **Check if already running**:
   ```bash
   ai-terminal-assistant -status
   ```

2. **Check logs**:
   ```bash
   tail -f /usr/local/var/log/ai-terminal-assistant.log
   ```

3. **Verify API key**:
   ```bash
   echo $OPENAI_API_KEY
   cat ~/.ai-terminal-assistant.env
   ```

4. **Check port availability**:
   ```bash
   lsof -i :8080
   ```

### Permission Issues

```bash
# Fix log directory permissions
sudo mkdir -p /usr/local/var/log
sudo chown $(whoami) /usr/local/var/log

# Fix PID directory permissions
sudo mkdir -p /usr/local/var/run
sudo chown $(whoami) /usr/local/var/run
```

### Service Management Issues

```bash
# Manually unload service
launchctl unload ~/Library/LaunchAgents/com.johnjallday.ai-terminal-assistant.plist

# Manually remove plist
rm ~/Library/LaunchAgents/com.johnjallday.ai-terminal-assistant.plist

# Reinstall
ai-terminal-assistant -install
```

### API Connection Issues

```bash
# Test local connection
curl -v http://localhost:8080/health

# Check if service is listening
netstat -an | grep 8080

# Check firewall (if applicable)
sudo pfctl -sr | grep 8080
```

## Development

### Running Locally

```bash
# Run daemon directly (not as service)
make run-daemon

# Or with custom port
./build/ai-terminal-assistant -port 9090
```

### Building from Source

```bash
# Build daemon only
make build-daemon

# Build all versions
make build-all

# Clean and rebuild
make clean build-daemon
```

### Testing

```bash
# Run tests
make test

# Test API endpoints
curl http://localhost:8080/health
curl -X POST http://localhost:8080/chat -H "Content-Type: application/json" -d '{"message": "test"}'
```

## Uninstallation

### Using Uninstall Script

```bash
# From installed location
scripts/uninstall.sh

# Or from release package
./uninstall.sh
```

### Manual Uninstallation

```bash
# Stop and uninstall service
ai-terminal-assistant -stop
ai-terminal-assistant -uninstall

# Remove binary
sudo rm /usr/local/bin/ai-terminal-assistant

# Remove config (optional)
rm ~/.ai-terminal-assistant.env

# Remove logs (optional)
sudo rm /usr/local/var/log/ai-terminal-assistant.log
```

## Security Considerations

- **API Key**: Store securely in environment file with appropriate permissions
- **Network**: The daemon binds to localhost by default (not accessible externally)
- **Logs**: May contain conversation data - secure appropriately
- **Permissions**: Runs as user service, not root

## Performance

- **Memory**: Typically uses 10-50MB RAM
- **CPU**: Minimal when idle, depends on API usage
- **Network**: Only outbound HTTPS to OpenAI API
- **Disk**: Log files may grow over time

## Integration Examples

### Shell Script Integration

```bash
#!/bin/bash
# Quick AI query script

QUERY="$*"
if [ -z "$QUERY" ]; then
    echo "Usage: ai-query <your question>"
    exit 1
fi

curl -s -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d "{\"message\": \"$QUERY\"}" | \
  jq -r '.response'
```

### Python Integration

```python
import requests
import json

def ask_ai(question, model="gpt-4.1-nano"):
    response = requests.post(
        "http://localhost:8080/chat",
        json={"message": question, "model": model}
    )
    return response.json()["response"]

# Usage
answer = ask_ai("What is the capital of France?")
print(answer)
```

## Support

- **Repository**: https://github.com/johnjallday/go-ai-terminal-assistant
- **Issues**: https://github.com/johnjallday/go-ai-terminal-assistant/issues
- **Documentation**: See `docs/` directory for more information