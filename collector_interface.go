package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type interfaceCollector struct {
	rxBytes  *prometheus.Desc
	txBytes  *prometheus.Desc
	rxPkts   *prometheus.Desc
	txPkts   *prometheus.Desc
	rxErrors *prometheus.Desc
	txErrors *prometheus.Desc
	rxDrops  *prometheus.Desc
	txDrops  *prometheus.Desc
	running  *prometheus.Desc
}

func newInterfaceCollector() subCollector {
	l := []string{"device", "address", "interface", "type"}
	d := func(n, h string) *prometheus.Desc {
		return prometheus.NewDesc("mikrotik_interface_"+n, h, l, nil)
	}
	return &interfaceCollector{
		rxBytes:  d("rx_bytes_total", "Total received bytes"),
		txBytes:  d("tx_bytes_total", "Total transmitted bytes"),
		rxPkts:   d("rx_packets_total", "Total received packets"),
		txPkts:   d("tx_packets_total", "Total transmitted packets"),
		rxErrors: d("rx_errors_total", "Total receive errors"),
		txErrors: d("tx_errors_total", "Total transmit errors"),
		rxDrops:  d("rx_drops_total", "Total receive drops"),
		txDrops:  d("tx_drops_total", "Total transmit drops"),
		running:  d("running", "Interface running state (1=up, 0=down)"),
	}
}

func (c *interfaceCollector) name() string { return "interface" }

func (c *interfaceCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.rxBytes
	ch <- c.txBytes
	ch <- c.rxPkts
	ch <- c.txPkts
	ch <- c.rxErrors
	ch <- c.txErrors
	ch <- c.rxDrops
	ch <- c.txDrops
	ch <- c.running
}

func (c *interfaceCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/interface/print",
		"=.proplist=name,type,running,disabled,rx-byte,tx-byte,rx-packet,tx-packet,rx-error,tx-error,rx-drop,tx-drop")
	if err != nil {
		log.Printf("[%s] interface/print: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		if re.Map["disabled"] == "true" {
			continue
		}
		l := []string{dev.Name, dev.Address, re.Map["name"], re.Map["type"]}
		emit := func(desc *prometheus.Desc, vt prometheus.ValueType, key string) {
			if v, ok := parseFloat(re.Map[key]); ok {
				ch <- prometheus.MustNewConstMetric(desc, vt, v, l...)
			}
		}
		ch <- prometheus.MustNewConstMetric(c.running, prometheus.GaugeValue, boolToFloat(re.Map["running"]), l...)
		emit(c.rxBytes, prometheus.CounterValue, "rx-byte")
		emit(c.txBytes, prometheus.CounterValue, "tx-byte")
		emit(c.rxPkts, prometheus.CounterValue, "rx-packet")
		emit(c.txPkts, prometheus.CounterValue, "tx-packet")
		emit(c.rxErrors, prometheus.CounterValue, "rx-error")
		emit(c.txErrors, prometheus.CounterValue, "tx-error")
		emit(c.rxDrops, prometheus.CounterValue, "rx-drop")
		emit(c.txDrops, prometheus.CounterValue, "tx-drop")
	}
}
