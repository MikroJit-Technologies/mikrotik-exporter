package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type firewallCollector struct {
	packets *prometheus.Desc
	bytes   *prometheus.Desc
}

func newFirewallCollector() subCollector {
	l := []string{"device", "address", "chain", "action", "comment"}
	return &firewallCollector{
		packets: prometheus.NewDesc("mikrotik_firewall_rule_packets_total",
			"Total packets matched by firewall rule", l, nil),
		bytes: prometheus.NewDesc("mikrotik_firewall_rule_bytes_total",
			"Total bytes matched by firewall rule", l, nil),
	}
}

func (c *firewallCollector) name() string { return "firewall" }

func (c *firewallCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.packets
	ch <- c.bytes
}

func (c *firewallCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/ip/firewall/filter/print",
		"=.proplist=chain,action,packets,bytes,comment",
		"?bytes>0")
	if err != nil {
		log.Printf("[%s] firewall/filter: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		l := []string{dev.Name, dev.Address, re.Map["chain"], re.Map["action"], re.Map["comment"]}
		if v, ok := parseFloat(re.Map["packets"]); ok {
			ch <- prometheus.MustNewConstMetric(c.packets, prometheus.CounterValue, v, l...)
		}
		if v, ok := parseFloat(re.Map["bytes"]); ok {
			ch <- prometheus.MustNewConstMetric(c.bytes, prometheus.CounterValue, v, l...)
		}
	}
}
