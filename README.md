<div align="center">

# mikrotik-exporter

**The most complete Prometheus exporter for MikroTik RouterOS**

[![CI](https://github.com/MikroJit-Technologies/mikrotik-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/MikroJit-Technologies/mikrotik-exporter/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![RouterOS](https://img.shields.io/badge/RouterOS-v6%20%2F%20v7-293239?style=flat-square&logo=mikrotik&logoColor=white)](https://mikrotik.com)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![Grafana](https://img.shields.io/badge/Grafana-dashboard%20included-F46800?style=flat-square&logo=grafana&logoColor=white)]()
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/MikroJit-Technologies/mikrotik-exporter?style=flat-square)](https://github.com/MikroJit-Technologies/mikrotik-exporter/releases)

Single binary · Multi-device · RouterOS v6 + v7 · Grafana dashboard included

</div>

---

## Why This One

Most MikroTik exporters cover interface stats and stop there. This one covers everything you actually need to monitor a real network:

| Collector | Metrics |
|---|---|
| **system** | CPU load, memory, storage, uptime, CPU count |
| **interface** | RX/TX bytes, packets, errors, drops, link state — all interfaces |
| **bgp** | Peer established state, session uptime, prefix count — v6 and v7 API |
| **ospf** | Neighbor full state |
| **dhcp** | Lease count by server and status |
| **firewall** | Per-rule packet and byte counters |
| **queue** | Simple queue rate, max-limit, bytes, drops — upload and download |
| **wireguard** | Peer last-handshake age, RX/TX bytes per peer |
| **capsman** | Client count per AP/SSID, signal strength, TX/RX rate |
| **health** | Board temperature, voltage, fan speed |

**Multi-device** — one exporter, as many routers as you have.  
**Graceful** — if one collector or one device fails, everything else still runs.  
**Zero deps** — single static binary, runs in scratch container or bare metal.

---

## Quick Start

### 1. Create a read-only API user on your MikroTik

```routeros
/user group add name=prometheus policy=read,api,!local,!telnet,!ssh,!ftp,!reboot,!write,!policy,!test,!winbox,!password,!web,!sniff,!sensitive,!romon
/user add name=prometheus group=prometheus password=your-password
```

Enable the API service if not already on:

```routeros
/ip service enable api
```

### 2. Configure

```bash
cp config.example.yml config.yml
```

```yaml
# config.yml
listen_address: ":9090"

devices:
  - name: "core-router"
    address: "192.168.88.1:8728"
    user: "prometheus"
    password: "your-password"

  - name: "edge-router"
    address: "10.0.0.1:8729"
    user: "prometheus"
    password: "your-password"
    tls: true
    skip_verify: false
    collectors: [system, interface, bgp, wireguard]  # omit = all collectors
```

### 3. Run

```bash
docker compose up -d
```

| Service | URL | Notes |
|---|---|---|
| Metrics | `http://localhost:9090/metrics` | Prometheus scrape target |
| Prometheus | `http://localhost:9091` | 30-day retention |
| Grafana | `http://localhost:3000` | `admin` / `admin` — dashboard pre-loaded |

---

## Configuration Reference

| Field | Default | Description |
|---|---|---|
| `listen_address` | `:9090` | Exporter HTTP listen address |
| `devices[].name` | — | Label used in all metrics |
| `devices[].address` | — | `host:port` — 8728 plaintext, 8729 TLS |
| `devices[].user` | — | RouterOS API username |
| `devices[].password` | — | RouterOS API password |
| `devices[].tls` | `false` | Use encrypted API connection |
| `devices[].skip_verify` | `false` | Skip TLS certificate verification |
| `devices[].collectors` | all | List of collectors to enable for this device |

Available collectors: `system` `interface` `bgp` `ospf` `dhcp` `firewall` `queue` `wireguard` `capsman` `health`

Set `CONFIG_FILE` env var to use a different config path (default: `config.yml`).

---

## Metrics Reference

All metrics carry `device` and `address` labels.

```
# Connectivity
mikrotik_up{device, address}                                          gauge

# System
mikrotik_system_cpu_load_percent{device, address}                     gauge
mikrotik_system_memory_free_bytes{device, address}                    gauge
mikrotik_system_memory_total_bytes{device, address}                   gauge
mikrotik_system_storage_free_bytes{device, address}                   gauge
mikrotik_system_storage_total_bytes{device, address}                  gauge
mikrotik_system_uptime_seconds{device, address}                       counter
mikrotik_system_cpu_count{device, address}                            gauge

# Interfaces
mikrotik_interface_running{..., interface, type}                      gauge
mikrotik_interface_rx_bytes_total{..., interface, type}               counter
mikrotik_interface_tx_bytes_total{..., interface, type}               counter
mikrotik_interface_rx_packets_total{..., interface, type}             counter
mikrotik_interface_tx_packets_total{..., interface, type}             counter
mikrotik_interface_rx_errors_total{..., interface, type}              counter
mikrotik_interface_tx_errors_total{..., interface, type}              counter
mikrotik_interface_rx_drops_total{..., interface, type}               counter
mikrotik_interface_tx_drops_total{..., interface, type}               counter

# BGP
mikrotik_bgp_peer_established{..., peer, remote_as, remote_address}  gauge
mikrotik_bgp_peer_uptime_seconds{..., peer, remote_as, remote_address} gauge
mikrotik_bgp_peer_prefix_count{..., peer, remote_as, remote_address} gauge

# OSPF
mikrotik_ospf_neighbor_full{..., router_id, interface, neighbor_address} gauge

# DHCP
mikrotik_dhcp_lease_count{..., server, status}                        gauge

# Firewall
mikrotik_firewall_rule_packets_total{..., chain, action, comment}     counter
mikrotik_firewall_rule_bytes_total{..., chain, action, comment}       counter

# Queues
mikrotik_queue_rate_bytes{..., queue, direction}                      gauge
mikrotik_queue_max_limit_bytes{..., queue, direction}                 gauge
mikrotik_queue_bytes_total{..., queue, direction}                     counter
mikrotik_queue_dropped_packets_total{..., queue, direction}           counter

# WireGuard
mikrotik_wireguard_peer_last_handshake_seconds{..., interface, public_key, endpoint} gauge
mikrotik_wireguard_peer_rx_bytes_total{..., interface, public_key, endpoint}         counter
mikrotik_wireguard_peer_tx_bytes_total{..., interface, public_key, endpoint}         counter

# CAPsMAN
mikrotik_capsman_client_count{..., interface, ssid}                   gauge
mikrotik_capsman_client_tx_rate_bytes{..., interface, ssid, mac}      gauge
mikrotik_capsman_client_rx_rate_bytes{..., interface, ssid, mac}      gauge
mikrotik_capsman_client_signal_dbm{..., interface, ssid, mac}         gauge

# Health
mikrotik_health_value{..., name, type}                                gauge
```

---

## Useful PromQL

```promql
# Memory usage %
100 * (1 - mikrotik_system_memory_free_bytes / mikrotik_system_memory_total_bytes)

# Interface throughput (bytes/s)
rate(mikrotik_interface_rx_bytes_total[2m])
rate(mikrotik_interface_tx_bytes_total[2m])

# BGP sessions down
mikrotik_bgp_peer_established == 0

# WireGuard peers with stale handshake (>3 min)
mikrotik_wireguard_peer_last_handshake_seconds > 180

# DHCP active leases
mikrotik_dhcp_lease_count{status="bound"}

# Queue utilization %
100 * mikrotik_queue_rate_bytes / mikrotik_queue_max_limit_bytes
```

---

## Build From Source

```bash
git clone https://github.com/MikroJit-Technologies/mikrotik-exporter.git
cd mikrotik-exporter
go build -o mikrotik-exporter .
CONFIG_FILE=config.yml ./mikrotik-exporter
```

Requires Go 1.22+. No CGO. Cross-compiles to any platform:

```bash
GOOS=linux GOARCH=arm64 go build -o mikrotik-exporter-arm64 .
GOOS=linux GOARCH=amd64 go build -o mikrotik-exporter-amd64 .
```

---

## Deploy Without Docker

```bash
# Copy binary and config
cp mikrotik-exporter /usr/local/bin/
cp config.example.yml /etc/mikrotik-exporter.yml

# Run
CONFIG_FILE=/etc/mikrotik-exporter.yml mikrotik-exporter
```

Systemd unit file:

```ini
[Unit]
Description=MikroTik Prometheus Exporter
After=network.target

[Service]
ExecStart=/usr/local/bin/mikrotik-exporter
Environment=CONFIG_FILE=/etc/mikrotik-exporter.yml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## Stack Layout

```
mikrotik-exporter   :9090   scrapes all configured devices on each Prometheus poll
prometheus          :9091   stores 30 days of time-series data
grafana             :3000   pre-provisioned dashboard — open and it just works
```

---

<div align="center">

Built by **[MikroJit Technologies](https://github.com/MikroJit-Technologies)** · Bangkok, Thailand

</div>
