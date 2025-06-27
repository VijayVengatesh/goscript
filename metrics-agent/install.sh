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
echo "🔍 Step 1: Detected architecture: $ARCH"

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

# Step 2: Download the binary
echo "⬇️ Step 2: Downloading agent from $AGENT_URL"
curl -L --fail -o metrics-agent "$AGENT_URL"
if [ $? -ne 0 ]; then
    echo "❌ Error: Failed to download the agent binary. Please check the URL or internet connection."
    exit 1
fi

# Step 3: Make executable and move
echo "⚙️ Step 3: Installing agent binary..."
export APP_ENV=production
chmod +x metrics-agent || { echo "❌ Failed to make binary executable"; exit 1; }
sudo mv metrics-agent /usr/local/bin/metrics-agent || { echo "❌ Failed to move binary to /usr/local/bin"; exit 1; }
echo "✅ Binary installed at /usr/local/bin/metrics-agent"

# Step 4: Create config
echo "🛠 Step 4: Creating config..."
CONFIG_DIR="/etc/metrics-agent"
sudo mkdir -p "$CONFIG_DIR" || { echo "❌ Failed to create config directory"; exit 1; }
echo "{\"user_id\": \"$USER_ID\"}" | sudo tee "$CONFIG_DIR/config.json" > /dev/null
if [ $? -ne 0 ]; then
    echo "❌ Failed to write config file"
    exit 1
fi
echo "✅ Config created at $CONFIG_DIR/config.json"

# Step 5: Create log file
LOG_FILE="/var/log/metrics-agent.log"
echo "📝 Step 5: Creating log file..."
sudo touch "$LOG_FILE" && sudo chmod 644 "$LOG_FILE" || { echo "❌ Failed to create log file"; exit 1; }
echo "✅ Log file created at $LOG_FILE"

# Step 6: Configure logrotate
echo "♻️ Step 6: Setting up log rotation..."
cat <<EOF | sudo tee /etc/logrotate.d/metrics-agent > /dev/null
/var/log/metrics-agent.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    copytruncate
}
EOF


if [ $? -ne 0 ]; then
    echo "❌ Failed to create logrotate config"
    exit 1
fi
echo "✅ Logrotate config set at /etc/logrotate.d/metrics-agent"
export APP_ENV=production

# Step 7: Start agent
echo "🚀 Step 7: Starting agent..."
nohup /usr/local/bin/metrics-agent >> "$LOG_FILE" 2>&1 &
sleep 2

# Step 8: Confirm it's running
pgrep -f metrics-agent > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ Installation complete. Agent is running in background. Logs: $LOG_FILE"
else
    echo "❌ Agent failed to start. Check the log file at $LOG_FILE"
    exit 1
fi
