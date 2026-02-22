#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "Error: please run as root (sudo ./uninstall.sh)"
  exit 1
fi

echo "[1/3] Stopping and disabling systemd service..."
systemctl stop drawbridge 2>/dev/null
systemctl disable drawbridge 2>/dev/null
rm -f /etc/systemd/system/drawbridge.service
systemctl daemon-reload

echo "[2/3] Removing binary..."
rm -f /usr/local/bin/drawbridge

echo "[3/3] Uninstallation complete."
echo "Note: configuration files (/etc/drawbridge) were preserved."
echo "To remove them completely, run:"
echo "sudo rm -rf /etc/drawbridge"
