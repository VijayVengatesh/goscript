#!/bin/bash

# Get user ID from arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -key) USER_ID="$2"; shift ;;
        *) echo "❌ Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [ -z "$USER_ID" ]; then
    echo "❌ Error: User ID is required. Use -key <your_user_id>"
    exit 1
fi

# Detect system architecture
ARCH=$(uname -m)
echo "🔍 Detected architecture: $ARCH"

# Map architecture to release binary
case "$ARCH" in
    x86_64)
        AGENT_URL="https://github.com/VijayVengatesh/goscript/releases/download/v1.0.0/metrics-agent-linux-amd64"
        ;;
    i386 | i686)
        AGENT_URL="https://github.com/VijayVengatesh/goscript/releases/download/v1.0.0/metrics-agent-linux-386"
        ;;
    aarch64 | arm64)
        AGENT_URL="https://github.com/VijayVengatesh/goscript/releases/download/v1.0.0/metrics-agent-linux-arm64"
        ;;
    armv7l | armv6l)
        AGENT_URL="https://github.com/VijayVengatesh/goscript/releases/download/v1.0.0/metrics-agent-linux-arm"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Download the binary
echo "⬇️ Downloading agent from $AGENT_URL"
curl -L -o metrics-agent "$AGENT_URL"

# Set permissions and move to /usr/local/bin
chmod +x metrics-agent
sudo mv metrics-agent /usr/local/bin/metrics-agent

# Create config directory and file
CONFIG_DIR="/etc/metrics-agent"
sudo mkdir -p "$CONFIG_DIR"
echo "{ \"user_id\": \"$USER_ID\" }" | sudo tee "$CONFIG_DIR/config.json" > /dev/null

# Create log file
LOG_FILE="/var/log/metrics-agent.log"
sudo touch "$LOG_FILE"
sudo chmod 644 "$LOG_FILE"

# Set up logrotate config for 7-day retention
cat <<EOF | sudo tee /etc/logrotate.d/metrics-agent > /dev/null
$LOG_FILE {
    daily
    rotate 7
    compress
    missingok
    notifempty
    copytruncate
}
EOF

# Start agent in background with logging
echo "🚀 Starting metrics-agent..."
nohup /usr/local/bin/metrics-agent >> "$LOG_FILE" 2>&1 &

echo "✅ Installation complete. Agent running in background. Logs: $LOG_FILE"
