package main

import (
	"log"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type ospfCollector struct {
	neighborFull *prometheus.Desc
}

func newOSPFCollector() subCollector {
	l := []string{"device", "address", "router_id", "interface", "neighbor_address"}
	return &ospfCollector{
		neighborFull: prometheus.NewDesc("mikrotik_ospf_neighbor_full",
			"OSPF neighbor in Full state (1=full)", l, nil),
	}
}

func (c *ospfCollector) name() string { return "ospf" }

func (c *ospfCollector) describe(ch chan<- *prometheus.Desc) { ch <- c.neighborFull }

func (c *ospfCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/routing/ospf/neighbor/print")
	if err != nil {
		log.Printf("[%s] ospf/neighbor: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		l := []string{dev.Name, dev.Address, re.Map["router-id"], re.Map["interface"], re.Map["address"]}
		full := 0.0
		if strings.ToLower(re.Map["state"]) == "full" {
			full = 1.0
		}
		ch <- prometheus.MustNewConstMetric(c.neighborFull, prometheus.GaugeValue, full, l...)
	}
}
