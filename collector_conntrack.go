package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type conntrackCollector struct {
	total    *prometheus.Desc
	estab    *prometheus.Desc
	timeWait *prometheus.Desc
	udp      *prometheus.Desc
	tcp      *prometheus.Desc
	icmp     *prometheus.Desc
}

func newConntrackCollector() subCollector {
	l := []string{"device", "address"}
	d := func(n, h string) *prometheus.Desc {
		return prometheus.NewDesc("mikrotik_conntrack_"+n, h, l, nil)
	}
	return &conntrackCollector{
		total:    d("total", "Total active connections in connection tracking table"),
		estab:    d("established", "Established TCP connections"),
		timeWait: d("time_wait", "TCP connections in TIME-WAIT state"),
		tcp:      d("tcp", "Total TCP connections"),
		udp:      d("udp", "Total UDP connections"),
		icmp:     d("icmp", "Total ICMP connections"),
	}
}

func (c *conntrackCollector) name() string { return "conntrack" }

func (c *conntrackCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.estab
	ch <- c.timeWait
	ch <- c.tcp
	ch <- c.udp
	ch <- c.icmp
}

func (c *conntrackCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/ip/firewall/connection/print",
		"=.proplist=protocol,tcp-state")
	if err != nil {
		log.Printf("[%s] firewall/connection: %v", dev.Name, err)
		return
	}

	var total, estab, timeWait, tcp, udp, icmp float64
	for _, re := range reply.Re {
		total++
		switch re.Map["protocol"] {
		case "tcp":
			tcp++
			switch re.Map["tcp-state"] {
			case "established":
				estab++
			case "time-wait":
				timeWait++
			}
		case "udp":
			udp++
		case "icmp":
			icmp++
		}
	}

	l := []string{dev.Name, dev.Address}
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, total, l...)
	ch <- prometheus.MustNewConstMetric(c.estab, prometheus.GaugeValue, estab, l...)
	ch <- prometheus.MustNewConstMetric(c.timeWait, prometheus.GaugeValue, timeWait, l...)
	ch <- prometheus.MustNewConstMetric(c.tcp, prometheus.GaugeValue, tcp, l...)
	ch <- prometheus.MustNewConstMetric(c.udp, prometheus.GaugeValue, udp, l...)
	ch <- prometheus.MustNewConstMetric(c.icmp, prometheus.GaugeValue, icmp, l...)
}
