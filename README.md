# 🚇 Tunnel Agent

A local agent that connects to the tunnel server and forwards requests to local services. This agent establishes a WebSocket connection to the tunnel server and forwards incoming requests to services running on your local machine.

## 🏗️ Architecture

```
Local Service ← Tunnel Agent ← WebSocket ← Tunnel Server ← Internet Request
```

The tunnel agent:
- Connects to tunnel server via WebSocket
- Forwards requests to local services
- Handles automatic reconnection
- Provides local monitoring API

## 🚀 Quick Start

### Build from Source

```bash
# Install dependencies
go mod tidy

# Build for current platform
go build -o agent ./cmd/agent

# Build for Windows (from Linux/macOS)
./build_windows.sh
```

### Usage

```bash
# Connect to tunnel server
./agent -server wss://your-app.onrender.com/agent/connect?id=my-tunnel -local http://127.0.0.1:9000

# With environment variables
export AGENT_SERVER_URL="wss://your-app.onrender.com/agent/connect?id=my-tunnel"
export AGENT_LOCAL_URL="http://127.0.0.1:9000"
./agent
```

### Environment Configuration

Copy `env.example` to `config.env` and modify the values:

```bash
cp env.example config.env
# Edit config.env with your settings
```

## 🔧 Configuration

### Command Line Flags

| Flag | Description | Default | Example |
|------|-------------|---------|---------|
| `-server` | WebSocket server URL | Required | `wss://app.onrender.com/agent/connect?id=test` |
| `-local` | Local service URL | `http://127.0.0.1:9000` | `http://127.0.0.1:3000` |
| `-port` | Local API port | `5050` | `8080` |
| `-timeout` | Request timeout (seconds) | `30` | `60` |
| `-debug` | Enable debug mode | `false` | `true` |

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `AGENT_SERVER_URL` | WebSocket server URL | Required | `wss://app.onrender.com/agent/connect?id=test` |
| `AGENT_LOCAL_URL` | Local service URL | `http://127.0.0.1:9000` | `http://127.0.0.1:3000` |
| `AGENT_LOCAL_PORT` | Local API port | `5050` | `8080` |
| `AGENT_TIMEOUT` | Request timeout | `30` | `60` |
| `DEBUG_MODE` | Debug mode | `false` | `true` |
| `AUTH_TOKEN` | Authentication token | None | `your-secret-token` |

## 🖥️ Windows Usage

### Prerequisites

- Windows 10/11 or Windows Server 2016+
- A local service running on your Windows machine
- Internet connection

### Quick Setup

1. **Copy files to Windows:**
   - `agent.exe` (Windows executable)
   - `config.env` (Configuration file)
   - `start-windows-agent.bat` (Startup script)
   - `test-windows-setup.bat` (Optional test script)

2. **Update configuration:**
   Edit `config.env` with your server URL:
   ```env
   AGENT_SERVER_URL=wss://your-app-name.onrender.com/agent/connect?id=windows-test
   AGENT_LOCAL_URL=http://127.0.0.1:9000
   AGENT_LOCAL_PORT=5050
   AGENT_TIMEOUT=30
   DEBUG_MODE=false
   AUTH_TOKEN=your-secret-token-here
   ```

3. **Start a local service:**
   ```cmd
   # Example: Simple HTTP server
   python -m http.server 9000
   
   # Or start your existing application
   # Update AGENT_LOCAL_URL in config accordingly
   ```

4. **Run the agent:**
   ```cmd
   start-windows-agent.bat
   ```

### Windows Testing

1. **Test setup (optional):**
   ```cmd
   test-windows-setup.bat
   ```

2. **Check agent status:**
   ```cmd
   curl http://127.0.0.1:5050/status
   ```

3. **Test tunnel:**
   ```cmd
   curl https://your-app-name.onrender.com/tunnel/windows-test/
   ```

### Manual Windows Setup

```cmd
# Build Windows executable (from Linux/macOS)
./build_windows.sh

# Copy to Windows and run manually
agent.exe -server wss://your-app-name.onrender.com/agent/connect?id=test -local http://127.0.0.1:9000
```

### Windows Troubleshooting

**Common Issues:**

1. **Connection fails:**
   - Check Windows Firewall settings
   - Verify server URL is correct
   - Ensure internet connectivity

2. **Agent won't start:**
   - Ensure `agent.exe` and `windows-config.env` are in same folder
   - Check configuration file syntax
   - Run from Command Prompt to see error messages

3. **Local service not accessible:**
   - Verify local service is running
   - Check local URL in configuration
   - Ensure service is listening on correct port

**Windows Firewall:**
If you encounter connection issues, allow the agent through Windows Firewall:
1. Open Windows Defender Firewall
2. Click "Allow an app or feature through Windows Defender Firewall"
3. Click "Change settings" → "Allow another app..."
4. Browse to `agent.exe` and add it
5. Check both "Private" and "Public" networks

## 🧪 Testing

### Step 1: Start Local Service

```bash
# Start a simple HTTP server
python -m http.server 9000

# Or use any local service
# Example: http://127.0.0.1:3000 (React app)
# Example: http://127.0.0.1:8000 (Django app)
# Example: http://127.0.0.1:8080 (Spring Boot app)
```

### Step 2: Connect Agent

```bash
# Connect to tunnel server
./agent -server wss://your-app.onrender.com/agent/connect?id=test -local http://127.0.0.1:9000
```

### Step 3: Test Tunnel

```bash
# Test the tunnel (from any machine)
curl https://your-app.onrender.com/tunnel/test/

# Should return your local service content
```



## 🔍 Monitoring

The agent provides a local API for monitoring:

### Check Agent Status
```bash
curl http://127.0.0.1:5050/status
```

**Response:**
```json
{
  "connected": true,
  "tunnel_id": "test",
  "local": "http://127.0.0.1:9000",
  "server": "wss://your-app.onrender.com/agent/connect?id=test"
}
```

### Reload Configuration
```bash
curl -X POST http://127.0.0.1:5050/reload
```

## 🚀 Deployment Examples

### Local Development
```bash
# Start local service
python -m http.server 9000

# Connect agent
./agent -server wss://dev-server.onrender.com/agent/connect?id=dev -local http://127.0.0.1:9000
```

### Production Use
```bash
# Start production service
./my-app --port 3000

# Connect agent
./agent -server wss://prod-server.onrender.com/agent/connect?id=prod -local http://127.0.0.1:3000
```

### Multiple Tunnels
```bash
# Tunnel 1: Web app
./agent -server wss://server.onrender.com/agent/connect?id=web -local http://127.0.0.1:3000

# Tunnel 2: API server
./agent -server wss://server.onrender.com/agent/connect?id=api -local http://127.0.0.1:8000
```

## 📄 License

This project is open source and available under the MIT License.