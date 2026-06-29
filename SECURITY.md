# Security Policy

## Credentials

- Never commit `config.yml` — it contains passwords. The file is in `.gitignore`.
- The exporter connects to RouterOS using a **read-only API user**. Use the minimal policy shown in the README.
- Avoid exposing `:9090/metrics` publicly. Place it behind a firewall or Prometheus basic auth.
- For encrypted API connections use `tls: true` (RouterOS API port 8729).

## Reporting a Vulnerability

Please report security issues by email to **contact@thiraphat.work** rather than opening a public issue.

Include:
- Description of the vulnerability
- Steps to reproduce
- Impact assessment

We aim to respond within 48 hours and release a fix within 14 days.
