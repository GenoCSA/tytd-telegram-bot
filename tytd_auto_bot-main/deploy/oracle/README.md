# TYTD – Oracle Always Free deployment

This package retains the upstream game engine, Telegram bot, built-in Mini App and dashboard. The new `/fleet` page is an account overview: the configured Telegram administrator can see all accounts; other users can see only their own account and accounts owned through the existing Leader/Đệ tử mechanism. **It does not turn every user into a leader**, and does not implement a new multi-account permission system. Existing `/members` controls require `IsLeader` in the database.

1. Provision Ubuntu on Oracle Cloud Always Free (ARM64 or AMD64) with sufficient available quota.
2. Upload/extract the ZIP to the VM. Run `sudo bash deploy/oracle/install.sh` from the project root.
3. Edit `/opt/tytd/config.yaml`: replace the example Telegram token and admin chat ID. Keep permissions `600`.
4. Run `sudo systemctl start tytd`; inspect `sudo journalctl -u tytd -f`.
5. Configure HTTPS using a reverse proxy (Caddy/nginx) to `127.0.0.1:8080`. Set Telegram Mini App URL to `https://YOUR_DOMAIN/miniapp` via BotFather. Enable the upstream `/web` feature using its documented command.
6. Use `https://YOUR_DOMAIN/fleet` to view authorized accounts. Existing leader members are managed at `/members`.

**Important:** `build/main-linux-*` are original upstream binaries, NOT rebuilt from modified source. To include `/fleet` in the running binary, rebuild from source with Go 1.26+ (`go build -o /opt/tytd/main ./cmd` after checking the upstream build commands and dependencies) and restart the service. Do not replace the executable while it is running. The installation script deliberately installs the upstream binary and does not claim otherwise.

Oracle Always Free is subject to resource availability, fair use and idle instance reclamation. It is not guaranteed perpetual uptime. The bot may be limited by upstream account restrictions and game rules. Never publish credentials, database or Telegram tokens.
