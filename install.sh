#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "Error: please run as root (sudo ./install.sh)"
  exit 1
fi

echo "[1/4] Compiling Drawbridge..."
go build -o drawbridge cmd/sniffer/main.go
if [ $? -ne 0 ]; then
    echo "Compilation failed."
    exit 1
fi

echo "[2/4] Installing binary..."
mv drawbridge /usr/local/bin/drawbridge
chmod +x /usr/local/bin/drawbridge

echo "[3/4] Setting up configuration..."
mkdir -p /etc/drawbridge
if [ ! -f /etc/drawbridge/config.yaml ]; then
    cp example.yaml /etc/drawbridge/config.yaml
    echo "      Default config installed."
else
    echo "      Config already exists. Skipping."
fi

echo "[4/4] Registering systemd service..."
cat <<EOF > /etc/systemd/system/drawbridge.service
[Unit]
Description=Drawbridge Port Knocking Daemon
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/drawbridge -config /etc/drawbridge/config.yaml
Restart=on-failure
RestartSec=5
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable drawbridge
systemctl restart drawbridge

echo ""
echo "Installation complete."
