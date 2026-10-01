#!/usr/bin/env bash
set -e

# Oracle Database Data Source for Grafana - Automated Installer
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/VinayakSomvanshi/oracle-grafana/main/install.sh | sudo bash

PLUGIN_ID="oracle-grafana-datasource"
VERSION="2.3.0"
ZIP_URL="https://github.com/VinayakSomvanshi/oracle-grafana/releases/latest/download/${PLUGIN_ID}-${VERSION}.zip"
TAR_URL="https://github.com/VinayakSomvanshi/oracle-grafana/releases/latest/download/${PLUGIN_ID}-${VERSION}.tar.gz"

echo "==> Installing Oracle Database Data Source for Grafana (v${VERSION})..."

# Check root privileges
if [ "$(id -u)" -ne 0 ]; then
  echo "ERROR: This installer must be run as root (or via sudo)." >&2
  exit 1
fi

# Detect Grafana plugins directory
PLUGIN_DIR=""
for candidate in "/var/lib/grafana/plugins" "/usr/local/var/lib/grafana/plugins" "/opt/grafana/plugins"; do
  if [ -d "$candidate" ]; then
    PLUGIN_DIR="$candidate"
    break
  fi
done

if [ -z "$PLUGIN_DIR" ]; then
  PLUGIN_DIR="/var/lib/grafana/plugins"
  mkdir -p "$PLUGIN_DIR"
fi

TARGET_DIR="${PLUGIN_DIR}/${PLUGIN_ID}"

# Try installation via grafana cli or grafana-cli
INSTALLED=0
if command -v grafana >/dev/null 2>&1; then
  echo "==> Detected grafana command. Attempting installation via 'grafana cli'..."
  if grafana cli --pluginUrl "$ZIP_URL" plugins install "$PLUGIN_ID" 2>/dev/null; then
    INSTALLED=1
  fi
fi

if [ "$INSTALLED" -eq 0 ] && command -v grafana-cli >/dev/null 2>&1; then
  echo "==> Detected legacy grafana-cli. Attempting installation..."
  if grafana-cli --homepath /usr/share/grafana --pluginUrl "$ZIP_URL" plugins install "$PLUGIN_ID" 2>/dev/null; then
    INSTALLED=1
  elif grafana-cli --pluginUrl "$ZIP_URL" plugins install "$PLUGIN_ID" 2>/dev/null; then
    INSTALLED=1
  fi
fi

if [ "$INSTALLED" -eq 0 ]; then
  echo "NOTICE: CLI installation unavailable or not completed. Using direct release archive extraction..."
fi

if [ "$INSTALLED" -eq 0 ]; then
  echo "==> Downloading and extracting release archive into ${TARGET_DIR}..."
  TMP_TAR="/tmp/${PLUGIN_ID}-${VERSION}.tar.gz"
  rm -f "$TMP_TAR"
  curl -fSL "$TAR_URL" -o "$TMP_TAR"

  rm -rf "$TARGET_DIR"
  mkdir -p "$TARGET_DIR"
  tar -xzf "$TMP_TAR" -C "$TARGET_DIR" --strip-components=1
  rm -f "$TMP_TAR"
fi

# Set proper ownership and permissions
if id -u grafana >/dev/null 2>&1; then
  chown -R grafana:grafana "$TARGET_DIR"
fi
chmod -R 755 "$TARGET_DIR"

echo "==> Installation complete."
echo ""
echo "Next steps:"
echo "1. Restart Grafana:"
echo "   sudo systemctl restart grafana-server"
echo "   (or: docker restart <grafana-container-name>)"
echo ""
echo "2. Open Grafana and add the data source:"
echo "   Connections > Data Sources > Add data source > Oracle"
