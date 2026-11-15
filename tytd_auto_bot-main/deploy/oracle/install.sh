#!/usr/bin/env bash
set -euo pipefail
# Run from repository root on Ubuntu ARM64/AMD64. No credentials embedded.
if [[ $(id -u) -ne 0 ]]; then echo 'Run: sudo bash deploy/oracle/install.sh'; exit 1; fi
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
ARCH="$(uname -m)"
case "$ARCH" in x86_64) BIN=main-linux-amd64;; aarch64|arm64) BIN=main-linux-arm64;; *) echo "Unsupported CPU: $ARCH"; exit 1;; esac
id tytd >/dev/null 2>&1 || useradd --system --home /opt/tytd --shell /usr/sbin/nologin tytd
install -d -o tytd -g tytd -m 700 /opt/tytd
install -m 755 "$ROOT/build/$BIN" /opt/tytd/main
if [[ ! -f /opt/tytd/config.yaml ]]; then install -m 600 -o tytd -g tytd "$ROOT/config.example.yaml" /opt/tytd/config.yaml; fi
chown tytd:tytd /opt/tytd/main /opt/tytd/config.yaml
cat >/etc/systemd/system/tytd.service <<'UNIT'
[Unit]
Description=TYTD Telegram game bot
After=network-online.target
Wants=network-online.target
[Service]
Type=simple
User=tytd
Group=tytd
WorkingDirectory=/opt/tytd
Environment=WEB_PORT=8080
ExecStart=/opt/tytd/main -config /opt/tytd/config.yaml
Restart=always
RestartSec=10
NoNewPrivileges=true
ProtectSystem=full
PrivateTmp=true
ReadWritePaths=/opt/tytd
[Install]
WantedBy=multi-user.target
UNIT
systemctl daemon-reload
systemctl enable tytd
printf '\nEdit /opt/tytd/config.yaml (token and admin ID), then: sudo systemctl start tytd\n'
