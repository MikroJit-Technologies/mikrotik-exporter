package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type queueCollector struct {
	rate     *prometheus.Desc
	maxLimit *prometheus.Desc
	dropped  *prometheus.Desc
	bytes    *prometheus.Desc
}

func newQueueCollector() subCollector {
	l := []string{"device", "address", "queue", "direction"}
	d := func(n, h string) *prometheus.Desc {
		return prometheus.NewDesc("mikrotik_queue_"+n, h, l, nil)
	}
	return &queueCollector{
		rate:     d("rate_bytes", "Current queue rate (bytes/s)"),
		maxLimit: d("max_limit_bytes", "Queue max-limit (bytes/s)"),
		dropped:  d("dropped_packets_total", "Total dropped packets"),
		bytes:    d("bytes_total", "Total bytes passed through queue"),
	}
}

func (c *queueCollector) name() string { return "queue" }

func (c *queueCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.rate
	ch <- c.maxLimit
	ch <- c.dropped
	ch <- c.bytes
}

func (c *queueCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/queue/simple/print",
		"=.proplist=name,rate,max-limit,dropped,bytes")
	if err != nil {
		log.Printf("[%s] queue/simple: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		qname := re.Map["name"]
		dirs := []string{"upload", "download"}
		for i, dir := range dirs {
			l := []string{dev.Name, dev.Address, qname, dir}
			if v, ok := splitPair(re.Map["rate"], i); ok {
				ch <- prometheus.MustNewConstMetric(c.rate, prometheus.GaugeValue, v, l...)
			}
			if v, ok := splitPair(re.Map["max-limit"], i); ok {
				ch <- prometheus.MustNewConstMetric(c.maxLimit, prometheus.GaugeValue, v, l...)
			}
			if v, ok := splitPair(re.Map["bytes"], i); ok {
				ch <- prometheus.MustNewConstMetric(c.bytes, prometheus.CounterValue, v, l...)
			}
		}
		if v, ok := parseFloat(re.Map["dropped"]); ok {
			ch <- prometheus.MustNewConstMetric(c.dropped, prometheus.CounterValue, v,
				dev.Name, dev.Address, qname, "total")
		}
	}
}
