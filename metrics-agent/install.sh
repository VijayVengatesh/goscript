#!/bin/bash

# --- Configuration ---
AGENT_VERSION="v1.0.0"
AGENT_NAME="metrics-agent"
CONFIG_DIR="/etc/$AGENT_NAME"
LOG_FILE="/var/log/$AGENT_NAME.log"
SERVICE_FILE="/etc/systemd/system/$AGENT_NAME.service"
LOGROTATE_FILE="/etc/logrotate.d/$AGENT_NAME"

# --- Argument Parsing ---
while [[ "$#" -gt 0 ]]; do
  case $1 in
    -key)
      USER_ID="$2"
      shift
      ;;
    *)
      echo "❌ Unknown parameter: $1"
      exit 1
      ;;
  esac
  shift
done

if [ -z "$USER_ID" ]; then
  echo "❌ Error: User ID is required. Use -key <your_user_id>"
  exit 1
fi

# --- Architecture Detection ---
ARCH=$(uname -m)
echo "🔍 Detected architecture: $ARCH"

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

# --- Download and Install Binary ---
echo "⬇️ Downloading $AGENT_NAME from $AGENT_URL"
curl -L --fail -o $AGENT_NAME "$AGENT_URL"
if [ $? -ne 0 ]; then
  echo "❌ Failed to download the agent binary"
  exit 1
fi

chmod +x $AGENT_NAME
sudo mv $AGENT_NAME /usr/local/bin/$AGENT_NAME

echo "✅ Binary installed to /usr/local/bin/$AGENT_NAME"

# --- Create Config Directory and File ---
sudo mkdir -p "$CONFIG_DIR"
echo "{\"user_id\": \"$USER_ID\"}" | sudo tee "$CONFIG_DIR/config.json" > /dev/null

echo "✅ Config file created at $CONFIG_DIR/config.json"

# --- Create Log File ---
sudo touch "$LOG_FILE"
sudo chmod 644 "$LOG_FILE"
echo "✅ Log file created at $LOG_FILE"

# --- Configure Logrotate ---
cat <<EOF | sudo tee "$LOGROTATE_FILE" > /dev/null
$LOG_FILE {
    daily
    rotate 7
    compress
    missingok
    notifempty
    copytruncate
}
EOF

echo "✅ Logrotate config created at $LOGROTATE_FILE"

# --- Setup systemd Service ---
cat <<EOF | sudo tee "$SERVICE_FILE" > /dev/null
[Unit]
Description=Metrics Agent Service
After=network.target

[Service]
ExecStart=/usr/local/bin/$AGENT_NAME
Restart=always
RestartSec=5
StandardOutput=append:$LOG_FILE
StandardError=append:$LOG_FILE

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reexec
sudo systemctl enable $AGENT_NAME --now

# --- Verify Service ---
if systemctl is-active --quiet $AGENT_NAME; then
  echo "✅ $AGENT_NAME service started successfully"
else
  echo "❌ Failed to start $AGENT_NAME service. Check logs at $LOG_FILE"
  exit 1
fi
