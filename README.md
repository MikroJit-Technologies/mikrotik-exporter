<div align="center">

# mikrotik-exporter

**Prometheus exporter for MikroTik RouterOS — the one that actually covers everything**

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker&logoColor=white)](https://hub.docker.com)
[![RouterOS](https://img.shields.io/badge/RouterOS-v6%20%2F%20v7-293239?style=flat-square&logo=mikrotik&logoColor=white)](https://mikrotik.com)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

</div>

---

## What This Is

A single-binary Prometheus exporter for MikroTik RouterOS that covers what other exporters miss:

| Collector | What It Exports |
|---|---|
| `system` | CPU load, memory, storage, uptime |
| `interface` | RX/TX bytes, packets, errors, drops, link state |
| `bgp` | Peer state, uptime, prefix count — RouterOS v6 and v7 |
| `ospf` | Neighbor state |
| `dhcp` | Lease count by server and status |
| `firewall` | Rule packet/byte counters |
| `queue` | Simple queue rate, max-limit, drops |
| `wireguard` | Peer last-handshake, RX/TX bytes |
| `capsman` | Client count, signal, TX/RX rate per AP |
| `health` | Temperature, voltage, fan speed |

Multi-device: monitor as many routers as you want from one exporter.  
Comes with a Grafana dashboard, Prometheus config, and full Docker Compose stack.

---

## Quick Start

**1. Create a read-only API user on your MikroTik:**

```
/user group add name=prometheus policy=read,api,!local,!telnet,!ssh,!ftp,!reboot,!write,!policy,!test,!winbox,!password,!web,!sniff,!sensitive,!romon
/user add name=prometheus group=prometheus password=your-password
```

**2. Copy the example config:**

```bash
cp config.example.yml config.yml
# edit config.yml — set device address, user, password
```

**3. Start the full stack:**

```bash
docker compose up -d
```

| Service | URL |
|---|---|
| Exporter metrics | http://localhost:9090/metrics |
| Prometheus | http://localhost:9091 |
| Grafana | http://localhost:3000 (admin / admin) |

The Grafana dashboard is provisioned automatically.

---

## Configuration

```yaml
listen_address: ":9090"

devices:
  - name: "core-router"
    address: "192.168.88.1:8728"   # port 8728 = plain, 8729 = TLS
    user: "prometheus"
    password: "your-password"
    tls: false

  - name: "edge-router"
    address: "10.0.0.1:8729"
    user: "prometheus"
    password: "your-password"
    tls: true
    skip_verify: false
    collectors: [system, interface, bgp, wireguard]  # omit = all enabled
```

---

## Metrics Reference

All metrics include `device` and `address` labels.

```
mikrotik_up                                    # 1=reachable, 0=down
mikrotik_system_cpu_load_percent
mikrotik_system_memory_free_bytes
mikrotik_system_memory_total_bytes
mikrotik_system_storage_free_bytes
mikrotik_system_storage_total_bytes
mikrotik_system_uptime_seconds
mikrotik_system_cpu_count

mikrotik_interface_running                     # labels: interface, type
mikrotik_interface_rx_bytes_total
mikrotik_interface_tx_bytes_total
mikrotik_interface_rx_packets_total
mikrotik_interface_tx_packets_total
mikrotik_interface_rx_errors_total
mikrotik_interface_tx_errors_total
mikrotik_interface_rx_drops_total
mikrotik_interface_tx_drops_total

mikrotik_bgp_peer_established                  # labels: peer, remote_as, remote_address
mikrotik_bgp_peer_uptime_seconds
mikrotik_bgp_peer_prefix_count

mikrotik_ospf_neighbor_full                    # labels: router_id, interface, neighbor_address

mikrotik_dhcp_lease_count                      # labels: server, status

mikrotik_firewall_rule_packets_total           # labels: chain, action, comment
mikrotik_firewall_rule_bytes_total

mikrotik_queue_rate_bytes                      # labels: queue, direction (upload/download)
mikrotik_queue_max_limit_bytes
mikrotik_queue_bytes_total
mikrotik_queue_dropped_packets_total

mikrotik_wireguard_peer_last_handshake_seconds # labels: interface, public_key, endpoint
mikrotik_wireguard_peer_rx_bytes_total
mikrotik_wireguard_peer_tx_bytes_total

mikrotik_capsman_client_count                  # labels: interface, ssid
mikrotik_capsman_client_tx_rate_bytes          # labels: interface, ssid, mac
mikrotik_capsman_client_rx_rate_bytes
mikrotik_capsman_client_signal_dbm

mikrotik_health_value                          # labels: name, type
```

---

## Build From Source

```bash
git clone https://github.com/MikroJit-Technologies/mikrotik-exporter.git
cd mikrotik-exporter
go build -o mikrotik-exporter .
CONFIG_FILE=config.yml ./mikrotik-exporter
```

---

## Run Without Docker

```bash
./mikrotik-exporter                   # reads config.yml by default
CONFIG_FILE=/etc/mt-exporter.yml ./mikrotik-exporter
```

---

## Stack Layout

```
mikrotik-exporter   :9090   Prometheus exporter
prometheus          :9091   Prometheus (scrapes exporter every 30s, 30d retention)
grafana             :3000   Grafana (dashboard auto-provisioned)
```

---

<div align="center">

Built by **[MikroJit Technologies](https://github.com/MikroJit-Technologies)** · Bangkok, Thailand

</div>
