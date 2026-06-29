package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type dhcpCollector struct {
	leaseCount *prometheus.Desc
}

func newDHCPCollector() subCollector {
	return &dhcpCollector{
		leaseCount: prometheus.NewDesc("mikrotik_dhcp_lease_count",
			"DHCP lease count by server and status",
			[]string{"device", "address", "server", "status"}, nil),
	}
}

func (c *dhcpCollector) name() string { return "dhcp" }

func (c *dhcpCollector) describe(ch chan<- *prometheus.Desc) { ch <- c.leaseCount }

func (c *dhcpCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/ip/dhcp-server/lease/print", "=.proplist=server,status")
	if err != nil {
		log.Printf("[%s] dhcp/lease: %v", dev.Name, err)
		return
	}

	counts := map[string]map[string]float64{}
	for _, re := range reply.Re {
		server := re.Map["server"]
		status := re.Map["status"]
		if counts[server] == nil {
			counts[server] = map[string]float64{}
		}
		counts[server][status]++
	}

	for server, statuses := range counts {
		for status, count := range statuses {
			l := []string{dev.Name, dev.Address, server, status}
			ch <- prometheus.MustNewConstMetric(c.leaseCount, prometheus.GaugeValue, count, l...)
		}
	}
}
