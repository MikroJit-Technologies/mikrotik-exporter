package main

import (
	"log"
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type capsmanCollector struct {
	clientCount *prometheus.Desc
	txRate      *prometheus.Desc
	rxRate      *prometheus.Desc
	signal      *prometheus.Desc
}

func newCapsManCollector() subCollector {
	lClient := []string{"device", "address", "interface", "ssid"}
	lPeer := []string{"device", "address", "interface", "ssid", "mac"}
	return &capsmanCollector{
		clientCount: prometheus.NewDesc("mikrotik_capsman_client_count",
			"Number of CAPsMAN clients per AP and SSID", lClient, nil),
		txRate: prometheus.NewDesc("mikrotik_capsman_client_tx_rate_bytes",
			"CAPsMAN client TX rate (bytes/s)", lPeer, nil),
		rxRate: prometheus.NewDesc("mikrotik_capsman_client_rx_rate_bytes",
			"CAPsMAN client RX rate (bytes/s)", lPeer, nil),
		signal: prometheus.NewDesc("mikrotik_capsman_client_signal_dbm",
			"CAPsMAN client signal strength (dBm)", lPeer, nil),
	}
}

func (c *capsmanCollector) name() string { return "capsman" }

func (c *capsmanCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.clientCount
	ch <- c.txRate
	ch <- c.rxRate
	ch <- c.signal
}

func (c *capsmanCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/caps-man/registration-table/print",
		"=.proplist=interface,ssid,mac-address,tx-rate,rx-rate,signal-strength")
	if err != nil {
		log.Printf("[%s] capsman: %v", dev.Name, err)
		return
	}

	counts := map[string]float64{}
	for _, re := range reply.Re {
		iface := re.Map["interface"]
		ssid := re.Map["ssid"]
		mac := re.Map["mac-address"]
		counts[iface+"|"+ssid]++

		lPeer := []string{dev.Name, dev.Address, iface, ssid, mac}
		if v, ok := parseFloat(re.Map["tx-rate"]); ok {
			ch <- prometheus.MustNewConstMetric(c.txRate, prometheus.GaugeValue, v, lPeer...)
		}
		if v, ok := parseFloat(re.Map["rx-rate"]); ok {
			ch <- prometheus.MustNewConstMetric(c.rxRate, prometheus.GaugeValue, v, lPeer...)
		}
		if v, ok := parseFloat(re.Map["signal-strength"]); ok {
			ch <- prometheus.MustNewConstMetric(c.signal, prometheus.GaugeValue, v, lPeer...)
		}
	}

	for key, count := range counts {
		parts := strings.SplitN(key, "|", 2)
		lClient := []string{dev.Name, dev.Address, parts[0], parts[1]}
		ch <- prometheus.MustNewConstMetric(c.clientCount, prometheus.GaugeValue, count, lClient...)
	}
}
