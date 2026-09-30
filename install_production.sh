#!/usr/bin/env bash
# =============================================================================
# Production Installer for Hardened Oracle Grafana DataSource
# =============================================================================
set -euo pipefail

if [ "$EUID" -ne 0 ]; then
  echo "[-] ERROR: This script must be run as root (or with sudo)." >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="${SCRIPT_DIR}/dist"
PLUGIN_TARGET="/var/lib/grafana/plugins/albertowd-oraclegrafana-datasource"

if [ ! -d "$DIST_DIR" ]; then
  echo "[*] Dist directory not found. Running 'make dist'..."
  (cd "$SCRIPT_DIR" && make dist)
fi

echo "[+] Installing plugin to: ${PLUGIN_TARGET}..."
mkdir -p "$PLUGIN_TARGET"
cp -r "${DIST_DIR}"/* "$PLUGIN_TARGET"/

# Set correct ownership and permissions
if id -u grafana >/dev/null 2>&1; then
  chown -R grafana:grafana "$PLUGIN_TARGET"
fi
chmod -R u=rwX,go=rX "$PLUGIN_TARGET"

# Ensure binary is executable
if [ -f "$PLUGIN_TARGET/gpx_oracle_grafana_linux_amd64" ]; then
  chmod 755 "$PLUGIN_TARGET/gpx_oracle_grafana_linux_amd64"
fi
if [ -f "$PLUGIN_TARGET/gpx_oracle_grafana_linux_arm64" ]; then
  chmod 755 "$PLUGIN_TARGET/gpx_oracle_grafana_linux_arm64"
fi

# SELinux support
if command -v restorecon >/dev/null 2>&1; then
  restorecon -Rv "$PLUGIN_TARGET" || true
fi

# Check if signed
if [ ! -f "$PLUGIN_TARGET/MANIFEST.txt" ]; then
  echo "[!] WARNING: Plugin is not yet signed with a Grafana Access Policy Token."
  echo "[!] Configuring albertowd-oraclegrafana-datasource in unsigned allowed list..."
  SYSCONFIG="/etc/sysconfig/grafana-server"
  DEFAULTS="/etc/default/grafana-server"
  TARGET_ENV=""
  if [ -f "$SYSCONFIG" ]; then
    TARGET_ENV="$SYSCONFIG"
  elif [ -f "$DEFAULTS" ]; then
    TARGET_ENV="$DEFAULTS"
  fi

  if [ -n "$TARGET_ENV" ]; then
    if ! grep -q "GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS.*albertowd-oraclegrafana-datasource" "$TARGET_ENV"; then
      echo 'GF_PLUGINS_ALLOW_LOADING_UNSIGNED_PLUGINS=albertowd-oraclegrafana-datasource' >> "$TARGET_ENV"
    fi
  fi
else
  echo "[+] Plugin is cryptographically SIGNED (MANIFEST.txt present)."
fi

echo "[+] Restarting Grafana service..."
if command -v systemctl >/dev/null 2>&1; then
  systemctl restart grafana-server || systemctl restart grafana || true
  sleep 3
  echo "[+] Checking Grafana logs:"
  journalctl -u grafana-server --since "1 minute ago" --no-pager | grep -i -E "oracle|plugin" | tail -n 15 || true
fi

echo "[✓] Installation complete. Open Grafana > Connections > Data Sources > Oracle to configure."
