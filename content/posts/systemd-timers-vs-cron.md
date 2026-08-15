---
title: "Replacing Cron with systemd Timers"
slug: f283819e
aliases: [systemd-timers-vs-cron]
date: 2025-11-03
tags: [linux, systemd, devops]
lang: en
draft: false
---

Cron is reliable and most Linux administrators know it by reflex. But it has gaps that become frustrating on modern systems: logs go nowhere useful by default, there is no dependency tracking, missed jobs during downtime are silently ignored, and debugging a failing cron job involves creative log spelunking.

systemd timers solve all of this. If your system runs systemd (and most Linux servers do), you already have everything you need.

## How systemd Timers Work

A timer is a pair of unit files: a `.timer` file that defines when to trigger, and a `.service` file that defines what to run. They link by sharing a name.

`backup.timer` activates `backup.service`. That is the entire relationship.

## A Simple Daily Backup

Create `/etc/systemd/system/backup.service`:

```ini
[Unit]
Description=Daily database backup

[Service]
Type=oneshot
User=deploy
ExecStart=/usr/local/bin/backup.sh
```

Create `/etc/systemd/system/backup.timer`:

```ini
[Unit]
Description=Run backup daily at 02:00

[Timer]
OnCalendar=*-*-* 02:00:00
Persistent=true

[Install]
WantedBy=timers.target
```

Enable and start:

```bash
systemctl daemon-reload
systemctl enable --now backup.timer
```

## OnCalendar Syntax

systemd uses its own calendar expression format. It is more readable than cron once you know it:

| Expression | Meaning |
|---|---|
| `daily` | Every day at midnight |
| `hourly` | Every hour |
| `weekly` | Every Monday at midnight |
| `*-*-* 02:00:00` | Every day at 02:00 |
| `Mon *-*-* 09:00:00` | Every Monday at 09:00 |
| `*-*-1 00:00:00` | 1st of every month |
| `*:0/15` | Every 15 minutes |

Validate any expression before using it:

```bash
systemd-analyze calendar "*-*-* 02:00:00"
```

This shows you the next several trigger times.

## Persistent: The Missed-Job Problem

The `Persistent=true` option in the timer tells systemd: if the timer was supposed to fire while the machine was off, run it as soon as the machine comes back up.

Cron has no equivalent. Jobs that fall during downtime are silently skipped.

## Logs That Actually Exist

Every timer invocation is logged to the journal:

```bash
journalctl -u backup.service
journalctl -u backup.service --since "1 hour ago"
```

You can see exit codes, stdout, stderr, and timestamps in one place. No more hunting through syslog for a cron entry.

## Listing Active Timers

```bash
systemctl list-timers
```

Output shows each timer, its last trigger time, next trigger time, and the associated unit. One command gives you a complete picture of all scheduled work on the system.

## When to Still Use Cron

Cron remains reasonable for:

- Per-user jobs (systemd user timers exist but are less common)
- Quick one-liners where writing two unit files is genuinely more overhead than the problem warrants
- Systems where the people maintaining it are more comfortable with cron

For anything that runs as a system service, requires reliable logging, or needs to handle missed executions gracefully, systemd timers are the better choice.
