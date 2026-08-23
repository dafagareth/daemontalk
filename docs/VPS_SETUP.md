# VPS Deployment Guide

Step-by-step instructions for deploying DaemonTalk to a Debian VPS using Docker, Caddy (automated TLS/HTTPS), and GitHub Actions.

## Domain Configuration

In your domain registrar DNS management panel, configure the following records:

- **A Record**: `@` (root) pointing to `YOUR_SERVER_PUBLIC_IP`
- **A Record**: `www` pointing to `YOUR_SERVER_PUBLIC_IP`

DNS changes typically propagate within 5 to 15 minutes.

## Server Provisioning

Log into your remote Debian server as `root` (or a sudo-enabled user):

```bash
$ ssh root@YOUR_SERVER_PUBLIC_IP
```

Run the automated provisioning script to install Docker, Docker Compose, Caddy, and configure the UFW firewall:

```bash
$ curl -sL https://raw.githubusercontent.com/dafagareth/daemontalk/main/scripts/setup-vps.sh | bash
```

The script automatically executes the following:
- Updates system packages (`apt update && apt upgrade`).
- Installs Docker Engine and Docker Compose plugin.
- Installs the Caddy web server from the official Cloudsmith repository.
- Configures UFW firewall rules: allows ports `22` (SSH), `80` (HTTP), `443` (HTTPS), and `2222` (TUI SSH).

## Initial Manual Deployment

For the very first time, you must clone the repository and set up the `.env` file manually. You can choose to place the project in either `/opt/daemontalk` (FHS standard) or `/var/www/daemontalk` (common web root).

```bash
# Choose your preferred path:
$ INSTALL_DIR="/opt/daemontalk"    # OR "/var/www/daemontalk"

$ git clone https://github.com/dafagareth/daemontalk.git $INSTALL_DIR
$ cd $INSTALL_DIR
$ cp .env.example .env
$ nano .env
```

Configure the production environment parameters in `.env`:

```env
PORT=8080
SSH_PORT=2222
ENV=production
ADMIN_TOKEN=generate_a_secure_random_token_here
BASE_URL=https://www.daemontalk.com
```

*(Tip: You can generate a secure `ADMIN_TOKEN` by running `openssl rand -hex 32` in your terminal)*

Copy the bundled `Caddyfile` to the system configuration path and reload Caddy:

```bash
$ sudo cp Caddyfile /etc/caddy/Caddyfile
$ sudo systemctl reload caddy
```

Start the service using Docker Compose. Since the VPS has limited RAM (1GB), we pull the pre-built image from GHCR rather than building it locally:

```bash
$ docker compose pull web
$ docker compose up -d
```

## Access Verification

Verify that all access endpoints are operational:

- **Web Interface (HTTPS)**: Open `https://daemontalk.com` in your browser.
- **SSH TUI Client**: Connect directly from your terminal:
  ```bash
  $ ssh daemontalk.com -p 2222
  ```

## Operational Scripts and Automation

The repository includes dedicated automation scripts for operational maintenance in the `scripts/` directory:

- **`deploy.sh`**: 1-Click Production Update. Pulls git commits, pulls the latest Docker image, and validates health status. **This is usually run automatically by GitHub Actions.**
- **`backup.sh`**: Automated Database & Content Backup. Creates timestamped `.tar.gz` archives of SQLite databases and markdown content with a 7-day retention policy.
- **`restore.sh <file>`**: Disaster Recovery. Safely restores databases and content from a backup archive with pre-restore safety snapshots.
- **`healthcheck.sh`**: Watchdog and Auto-Recovery. Monitors the `/healthz` endpoint and automatically restarts the service if it becomes unresponsive.

**Recommended Crontab Setup**
Open your system crontab (`crontab -e`) and add the following scheduled jobs (adjust the path to match where you installed the project):

```bash
# Automated daily backup at 03:00 AM
0 3 * * * /opt/daemontalk/scripts/backup.sh > /dev/null 2>&1

# Health monitoring watchdog every 5 minutes (auto-restarts on failure)
*/5 * * * * /opt/daemontalk/scripts/healthcheck.sh > /dev/null 2>&1
```

## Continuous Deployment (CI/CD)

Once the initial manual deployment is running, you can enable 100% automated deployment whenever you push code or new articles to GitHub:

1. In your GitHub repository, go to **Settings > Secrets and variables > Actions**.
2. Add the following repository secrets:
   - `VPS_HOST`: Your server public IP.
   - `VPS_USER`: `root` (or your sudo deploy user).
   - `VPS_SSH_KEY`: Your private SSH key content (`~/.ssh/id_ed25519`).
   - `VPS_PORT`: `22` (default).

Whenever you push to the `main` branch, GitHub Actions will build the Docker image, push it to GHCR, and automatically trigger `./scripts/deploy.sh` on your VPS to update the live server.

## Optional: Performance & Security with Cloudflare CDN

For a 1GB RAM VPS, utilizing a Content Delivery Network (CDN) like Cloudflare is highly recommended to offload bandwidth (especially for heavy media like `.mp4` videos) and protect against malicious traffic or DDoS attacks.

If you choose to route your traffic through Cloudflare, you must configure it correctly so it doesn't conflict with Caddy's automatic HTTPS/SSL provisioning:

1. **DNS Setup**: In your Cloudflare dashboard, configure your `A` records to point to your VPS IP and enable the **Proxied (Orange Cloud)** status.
2. **SSL/TLS Mode**: Navigate to **SSL/TLS > Overview** and set the encryption mode to **Full (Strict)**.
   - *Why?* Caddy automatically provisions valid Let's Encrypt certificates on your VPS. Setting Cloudflare to "Full (Strict)" ensures Cloudflare strictly encrypts the proxy connection and validates Caddy's underlying certificate without causing redirect loops.
3. **Bandwidth Savings**: Once enabled, Cloudflare will automatically cache static assets (like `/static/images` and CSS) at their global edge nodes. Your VPS will only be responsible for serving the lightweight HTMX/HTML payloads.
