# Changelog

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
