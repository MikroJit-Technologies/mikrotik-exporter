package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type wireguardCollector struct {
	lastHandshake *prometheus.Desc
	rxBytes       *prometheus.Desc
	txBytes       *prometheus.Desc
}

func newWireGuardCollector() subCollector {
	l := []string{"device", "address", "interface", "public_key", "endpoint"}
	return &wireguardCollector{
		lastHandshake: prometheus.NewDesc("mikrotik_wireguard_peer_last_handshake_seconds",
			"Seconds elapsed since last WireGuard handshake", l, nil),
		rxBytes: prometheus.NewDesc("mikrotik_wireguard_peer_rx_bytes_total",
			"Total bytes received from WireGuard peer", l, nil),
		txBytes: prometheus.NewDesc("mikrotik_wireguard_peer_tx_bytes_total",
			"Total bytes transmitted to WireGuard peer", l, nil),
	}
}

func (c *wireguardCollector) name() string { return "wireguard" }

func (c *wireguardCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.lastHandshake
	ch <- c.rxBytes
	ch <- c.txBytes
}

func (c *wireguardCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/interface/wireguard/peers/print")
	if err != nil {
		log.Printf("[%s] wireguard/peers: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		endpoint := re.Map["endpoint-address"]
		if port := re.Map["endpoint-port"]; port != "" {
			endpoint += ":" + port
		}
		l := []string{dev.Name, dev.Address, re.Map["interface"], re.Map["public-key"], endpoint}

		if v := parseUptime(re.Map["last-handshake"]); v > 0 {
			ch <- prometheus.MustNewConstMetric(c.lastHandshake, prometheus.GaugeValue, v, l...)
		}
		if v, ok := parseFloat(re.Map["rx"]); ok {
			ch <- prometheus.MustNewConstMetric(c.rxBytes, prometheus.CounterValue, v, l...)
		}
		if v, ok := parseFloat(re.Map["tx"]); ok {
			ch <- prometheus.MustNewConstMetric(c.txBytes, prometheus.CounterValue, v, l...)
		}
	}
}
