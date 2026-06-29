package main

import (
	"log"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type bgpCollector struct {
	established *prometheus.Desc
	uptime      *prometheus.Desc
	prefixes    *prometheus.Desc
}

func newBGPCollector() subCollector {
	l := []string{"device", "address", "peer", "remote_as", "remote_address"}
	return &bgpCollector{
		established: prometheus.NewDesc("mikrotik_bgp_peer_established",
			"BGP peer in established state (1=yes)", l, nil),
		uptime: prometheus.NewDesc("mikrotik_bgp_peer_uptime_seconds",
			"BGP session uptime in seconds", l, nil),
		prefixes: prometheus.NewDesc("mikrotik_bgp_peer_prefix_count",
			"Number of prefixes received from peer", l, nil),
	}
}

func (c *bgpCollector) name() string { return "bgp" }

func (c *bgpCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.established
	ch <- c.uptime
	ch <- c.prefixes
}

func (c *bgpCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	// Try RouterOS v7 API first
	reply, err := client.Run("/routing/bgp/session/print")
	isV7 := err == nil && len(reply.Re) > 0

	if !isV7 {
		reply, err = client.Run("/routing/bgp/peer/print")
		if err != nil {
			log.Printf("[%s] bgp: %v", dev.Name, err)
			return
		}
	}

	for _, re := range reply.Re {
		var peerName, remoteAS, remoteAddr string
		if isV7 {
			peerName = re.Map["name"]
			remoteAS = re.Map["remote.as"]
			remoteAddr = re.Map["remote.address"]
		} else {
			peerName = re.Map["name"]
			remoteAS = re.Map["remote-as"]
			remoteAddr = re.Map["remote-address"]
		}

		l := []string{dev.Name, dev.Address, peerName, remoteAS, remoteAddr}

		up := 0.0
		if strings.ToLower(re.Map["state"]) == "established" {
			up = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.established, prometheus.GaugeValue, up, l...)

		if v := parseUptime(re.Map["uptime"]); v > 0 {
			ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.GaugeValue, v, l...)
		}
		if v, ok := parseFloat(re.Map["prefix-count"]); ok {
			ch <- prometheus.MustNewConstMetric(c.prefixes, prometheus.GaugeValue, v, l...)
		}
	}
}
