# Changelog

## [1.1.0] — 2026-06-29

### Added
- `routes` collector — route table size by type (total, active, static, BGP, OSPF, connected, dynamic)
- `conntrack` collector — connection tracking table (total, TCP/UDP/ICMP, established, time-wait)
- Multi-arch Docker image published to `ghcr.io` (amd64, arm64, arm/v7)
- GitHub issue templates (bug report, feature request)
- Pull request template
- Security policy (`SECURITY.md`)
- CI builds linux/arm/v7 binary

## [1.0.0] — 2026-06-29

### Added
- Initial release
- 10 collectors: system, interface, bgp, ospf, dhcp, firewall, queue, wireguard, capsman, health
- Multi-device support — monitor multiple routers from one exporter
- RouterOS v6 and v7 BGP API auto-detection
- Docker Compose stack with Prometheus and auto-provisioned Grafana dashboard
- TLS support for encrypted RouterOS API connections
- Per-device collector filtering
- Graceful error handling — one failed device or collector does not affect others
- Cross-platform builds: linux/amd64, linux/arm64, darwin/arm64, windows/amd64
