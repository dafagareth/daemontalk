# Production VPS Deployment Guide

Step-by-step instructions for deploying `daemontalk.com` to a Debian VPS using Docker, Caddy (automated TLS/HTTPS), and systemd.

---

## 1. Configure DNS Records

In your domain registrar DNS management panel (Cloudflare, Namecheap, or Rumahweb), configure the following records:

| Type | Host / Name | Target / IP Address | TTL |
| :--- | :--- | :--- | :--- |
| **A** | `@` (root) | `YOUR_SERVER_PUBLIC_IP` | 300 / Auto |
| **CNAME** | `www` | `daemontalk.com` | 300 / Auto |

DNS changes typically propagate within 5 to 15 minutes.

---

## 2. Server Provisioning (Debian 12 Bookworm)

SSH into your remote server as root:

```bash
ssh root@YOUR_SERVER_PUBLIC_IP
```

Run the automated provisioning script to install Docker, Docker Compose, Caddy, and configure the UFW firewall:

```bash
curl -sL https://raw.githubusercontent.com/dafagareth/daemontalk/main/scripts/setup-vps.sh | bash
```

The script automatically executes the following:
- Updates system packages (`apt update && apt upgrade`).
- Installs Docker Engine and Docker Compose plugin.
- Installs the Caddy web server from the official Cloudsmith repository.
- Configures UFW firewall rules: allows ports `22` (SSH), `80` (HTTP), `443` (HTTPS), and `2222` (TUI SSH).

---

## 3. Clone Repository and Set Environment Variables

```bash
# Clone repository
git clone https://github.com/dafagareth/daemontalk.git /opt/daemontalk
cd /opt/daemontalk

# Initialize environment configuration
cp .env.example .env
nano .env
```

Configure the production environment parameters:

```env
PORT=8080
SSH_PORT=2222
ENV=production
ADMIN_TOKEN=generate_a_secure_random_token_here
BASE_URL=https://daemontalk.com
```

---

## 4. Configure Caddy for Automated HTTPS

Copy the bundled `Caddyfile` to the system configuration path and reload Caddy:

```bash
sudo cp /opt/daemontalk/Caddyfile /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

Caddy will automatically obtain and renew TLS certificates from Let's Encrypt and ZeroSSL for `daemontalk.com` and `www.daemontalk.com`.

---

## 5. Launch the Application Container

Start the service using Docker Compose:

```bash
cd /opt/daemontalk
docker compose up -d --build
```

Verify that the container is running and healthy:

```bash
docker compose ps
docker compose logs -f
```

---

## 6. Access Verification

Verify that all access endpoints are operational:

- **Web Interface (HTTPS)**: Open `https://daemontalk.com` in your browser.
- **SSH TUI Client**: Connect directly from your terminal:
  ```bash
  ssh daemontalk.com -p 2222
  ```
- **CLI Feed Stream**:
  ```bash
  curl -sL https://daemontalk.com/daily
  ```

---

## 7. Operational Scripts and Automation

The repository includes dedicated automation scripts for operational maintenance:

| Script | Purpose and Usage |
| :--- | :--- |
| **`scripts/deploy.sh`** | **1-Click Production Update**: Automatically creates a pre-deploy safety backup, pulls git commits, rebuilds Docker containers, prunes stale images, and validates health status.<br>`./scripts/deploy.sh` (fast cached build)<br>`./scripts/deploy.sh --fresh` (100% clean rebuild without Docker cache) |
| **`scripts/backup.sh`** | **Automated Database & Content Backup**: Creates timestamped `.tar.gz` archives of SQLite databases (safe for WAL mode) and markdown content with a 7-day retention policy.<br>`./scripts/backup.sh` |
| **`scripts/restore.sh`** | **Disaster Recovery**: Safely restores databases and content from a backup archive with pre-restore safety snapshots.<br>`./scripts/restore.sh /opt/daemontalk/backups/daemontalk_backup_YYYYMMDD_HHMMSS.tar.gz` |
| **`scripts/healthcheck.sh`** | **Watchdog and Auto-Recovery**: Monitors the `/healthz` endpoint and automatically restarts the service if it becomes unresponsive.<br>`./scripts/healthcheck.sh` |

### Recommended Crontab Setup

Open your system crontab (`crontab -e`) and add the following scheduled jobs:

```bash
# Automated daily backup at 03:00 AM
0 3 * * * /opt/daemontalk/scripts/backup.sh > /dev/null 2>&1

# Health monitoring watchdog every 5 minutes (auto-restarts on failure)
*/5 * * * * /opt/daemontalk/scripts/healthcheck.sh > /dev/null 2>&1
```

---

## 8. Continuous Deployment with GitHub Actions (Optional)

To enable 100% automated deployment whenever you push code or new articles to GitHub:

1. In your GitHub repository, go to **Settings > Secrets and variables > Actions**.
2. Add the following repository secrets:
   - `VPS_HOST`: Your server public IP (e.g. `103.xxx.xxx.xxx` or `daemontalk.com`).
   - `VPS_USER`: `root` (or your sudo deploy user).
   - `VPS_SSH_KEY`: Your private SSH key content (`~/.ssh/id_ed25519`).
   - `VPS_PORT`: `22` (default).

Whenever you push to the `main` branch, `.github/workflows/deploy.yml` will automatically trigger `./scripts/deploy.sh` on your VPS.

