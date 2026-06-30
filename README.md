<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=293239&height=120&section=header&text=mikrotik-exporter&fontSize=36&fontColor=ffffff&fontAlignY=38&desc=Prometheus%20exporter%20for%20MikroTik%20RouterOS&descAlignY=62&descColor=c9d1d9" width="100%"/>

[![CI](https://github.com/MikroJit-Technologies/mikrotik-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/MikroJit-Technologies/mikrotik-exporter/actions/workflows/ci.yml)
[![Go 1.22](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Docker](https://img.shields.io/badge/ghcr.io%2Fmikrotik--exporter-latest-2496ED?style=flat-square&logo=docker&logoColor=white)](https://ghcr.io/mikrojit-technologies/mikrotik-exporter)
[![Prometheus](https://img.shields.io/badge/Prometheus-compatible-E6522C?style=flat-square&logo=prometheus&logoColor=white)](https://prometheus.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-388bfd?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/MikroJit-Technologies/mikrotik-exporter?style=flat-square&color=3fb950)](https://github.com/MikroJit-Technologies/mikrotik-exporter/releases)

**The most complete Prometheus exporter for MikroTik RouterOS — v6 and v7 supported.**

[Collectors](#collectors) · [Quick Start](#quick-start) · [Configuration](#configuration) · [Metrics Reference](#metrics-reference) · [Grafana](#grafana-dashboard) · [Docker](#docker)

</div>

---

## Why mikrotik-exporter?

Most MikroTik Prometheus exporters expose 3–5 metrics. This one exports **10+ collector domains** via RouterOS SSH API — including BGP session state, OSPF neighbors, WireGuard peer stats, CAPsMAN access point details, and hardware health sensors.

```
mikrotik_interface_rx_bytes{device="main-router",name="ether1"} 1.23e+10
mikrotik_bgp_session_state{device="main-router",peer="10.0.0.1",state="established"} 1
mikrotik_wireguard_peer_rx_bytes{device="main-router",pubkey="xLgJk3N..."} 4.56e+08
mikrotik_dhcp_lease_count{device="main-router",server="dhcp1"} 47
mikrotik_system_cpu_load{device="main-router"} 12
mikrotik_system_free_memory{device="main-router"} 1.07e+08
```

---

## Collectors

| Domain | Metrics |
|---|---|
| **System** | CPU load, free/total memory, uptime, RouterOS version |
| **Interface** | RX/TX bytes, packets, errors, drops, link status, speed |
| **IP Routes** | Total, active, static, BGP, OSPF, connected, dynamic counts |
| **BGP** | Session state (established/idle/active), prefix counts per peer (v6 + v7) |
| **OSPF** | Neighbor state, adjacency count |
| **DHCP** | Active lease count per server |
| **Firewall** | Rule packet/byte counters, connection tracking total |
| **Queues** | Simple queue RX/TX bytes and queue drop counts |
| **WireGuard** | Peer RX/TX bytes, last handshake timestamp, peer count |
| **CAPsMAN** | AP registration count, client count, TX/RX per AP |
| **Health** | CPU temperature, board temperature, voltage, fan speed |
| **Connection Tracking** | Total connections, TCP established, UDP, ICMP counts |

---

## Quick Start

**Docker Compose (recommended)**

```bash
cp config.example.yml config.yml
# edit config.yml — add your devices
docker compose up -d
```

Add to your **Prometheus** `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: mikrotik
    static_configs:
      - targets: ['localhost:9436']
```

**Binary**

```bash
curl -Lo mikrotik-exporter https://github.com/MikroJit-Technologies/mikrotik-exporter/releases/latest/download/mikrotik-exporter-linux-amd64
chmod +x mikrotik-exporter
./mikrotik-exporter -config config.yml
```

Metrics available at **http://localhost:9436/metrics**

---

## Configuration

```yaml
listen_addr: ":9436"
interval: "60s"       # collection interval

devices:
  - name: main-router
    host: 10.33.1.1
    port: 22
    username: prometheus
    password: "secret"

  - name: branch-router
    host: 192.168.1.1
    username: prometheus
    identity: /etc/mikrotik-exporter/id_rsa   # SSH key auth

collectors:
  - system
  - interfaces
  - routes
  - bgp
  - ospf
  - dhcp
  - firewall
  - queues
  - wireguard
  - capsman
  - health
  - conntrack
```

### Reference

| Key | Default | Description |
|---|---|---|
| `listen_addr` | `:9436` | Prometheus scrape endpoint |
| `interval` | `60s` | How often to poll devices |
| `device.name` | — | Label name in metrics |
| `device.host` | — | RouterOS IP or hostname |
| `device.port` | `22` | SSH port |
| `device.username` | `admin` | RouterOS SSH user |
| `device.password` | — | Password auth |
| `device.identity` | — | SSH private key path (mutually exclusive with password) |
| `collectors` | all enabled | List of collectors to run |

### RouterOS user setup

Create a read-only API user on your MikroTik:

```routeros
/user group add name=prometheus policy=read,api,!local,!telnet,!ssh,!ftp,!reboot,!write,!policy,!test,!winbox,!password,!web,!sniff,!sensitive,!romon
/user add name=prometheus group=prometheus password=secret
```

For SSH-based collection (used by this exporter), use the `ssh` policy instead:

```routeros
/user group add name=exporter policy=read,ssh,!local,!telnet,!ftp,!reboot,!write,!policy,!test,!winbox,!password,!web,!sniff,!sensitive,!romon
/user add name=exporter group=exporter password=secret
```

---

## Metrics Reference

### System

| Metric | Type | Description |
|---|---|---|
| `mikrotik_system_cpu_load` | gauge | CPU load percentage |
| `mikrotik_system_free_memory` | gauge | Free memory in bytes |
| `mikrotik_system_total_memory` | gauge | Total memory in bytes |
| `mikrotik_system_uptime_seconds` | gauge | System uptime |

### Interfaces

| Metric | Labels | Description |
|---|---|---|
| `mikrotik_interface_rx_bytes` | `name`, `type` | Received bytes |
| `mikrotik_interface_tx_bytes` | `name`, `type` | Transmitted bytes |
| `mikrotik_interface_rx_errors` | `name` | RX error count |
| `mikrotik_interface_tx_drops` | `name` | TX drop count |
| `mikrotik_interface_running` | `name` | 1 if link is up |

### BGP

| Metric | Labels | Description |
|---|---|---|
| `mikrotik_bgp_session_state` | `peer`, `remote_as` | 1 if established |
| `mikrotik_bgp_prefix_count` | `peer` | Received prefix count |

### WireGuard

| Metric | Labels | Description |
|---|---|---|
| `mikrotik_wireguard_peer_rx_bytes` | `interface`, `pubkey` | Bytes received from peer |
| `mikrotik_wireguard_peer_tx_bytes` | `interface`, `pubkey` | Bytes sent to peer |
| `mikrotik_wireguard_peer_last_handshake` | `interface`, `pubkey` | Handshake Unix timestamp |
| `mikrotik_wireguard_peer_up` | `interface`, `pubkey` | 1 if handshake ≤ 180s ago |

> All metrics carry a `device` label matching the name in config.

---

## Grafana Dashboard

Import the bundled dashboard (`grafana/dashboard.json`) for instant visibility into:

- Interface traffic (per device, stacked)
- BGP session state history
- CPU / memory over time
- WireGuard peer connectivity
- DHCP lease count trend
- Firewall connection tracking

---

## Docker

```bash
docker run -d \
  --name mikrotik-exporter \
  -p 9436:9436 \
  -v $(pwd)/config.yml:/app/config.yml:ro \
  ghcr.io/mikrojit-technologies/mikrotik-exporter:latest
```

---

## Part of MikroJit Technologies

<div align="center">

| Tool | Description |
|---|---|
| **mikrotik-exporter** | **Prometheus exporter for MikroTik RouterOS** |
| [mikrotik-backup](https://github.com/MikroJit-Technologies/mikrotik-backup) | Auto-backup RouterOS configs to Git + Telegram diffs |
| [wireguard-manager](https://github.com/MikroJit-Technologies/wireguard-manager) | Web UI for WireGuard peer management |
| [netmon](https://github.com/MikroJit-Technologies/netmon) | Uptime monitor — HTTP / ping / TCP + Telegram alerts |
| [routeros-cli](https://github.com/MikroJit-Technologies/routeros-cli) | CLI tool for multi-device RouterOS command execution |

</div>

---

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=293239&height=80&section=footer" width="100%"/>

MIT License · [MikroJit Technologies](https://github.com/MikroJit-Technologies) · Bangkok, Thailand

</div>
