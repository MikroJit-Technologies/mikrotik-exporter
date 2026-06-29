package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type healthCollector struct {
	value *prometheus.Desc
}

func newHealthCollector() subCollector {
	return &healthCollector{
		value: prometheus.NewDesc("mikrotik_health_value",
			"Hardware health sensor value (temperature °C, voltage V, fan RPM)",
			[]string{"device", "address", "name", "type"}, nil),
	}
}

func (c *healthCollector) name() string { return "health" }

func (c *healthCollector) describe(ch chan<- *prometheus.Desc) { ch <- c.value }

func (c *healthCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/system/health/print")
	if err != nil {
		log.Printf("[%s] system/health: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		if v, ok := parseFloat(re.Map["value"]); ok {
			l := []string{dev.Name, dev.Address, re.Map["name"], re.Map["type"]}
			ch <- prometheus.MustNewConstMetric(c.value, prometheus.GaugeValue, v, l...)
		}
	}
}
