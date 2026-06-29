package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type systemCollector struct {
	cpuLoad  *prometheus.Desc
	memFree  *prometheus.Desc
	memTotal *prometheus.Desc
	hddFree  *prometheus.Desc
	hddTotal *prometheus.Desc
	uptime   *prometheus.Desc
	cpuCount *prometheus.Desc
}

func newSystemCollector() subCollector {
	l := []string{"device", "address"}
	d := func(n, h string) *prometheus.Desc {
		return prometheus.NewDesc("mikrotik_system_"+n, h, l, nil)
	}
	return &systemCollector{
		cpuLoad:  d("cpu_load_percent", "CPU load (%)"),
		memFree:  d("memory_free_bytes", "Free memory bytes"),
		memTotal: d("memory_total_bytes", "Total memory bytes"),
		hddFree:  d("storage_free_bytes", "Free storage bytes"),
		hddTotal: d("storage_total_bytes", "Total storage bytes"),
		uptime:   d("uptime_seconds", "System uptime in seconds"),
		cpuCount: d("cpu_count", "Number of CPUs"),
	}
}

func (c *systemCollector) name() string { return "system" }

func (c *systemCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.cpuLoad
	ch <- c.memFree
	ch <- c.memTotal
	ch <- c.hddFree
	ch <- c.hddTotal
	ch <- c.uptime
	ch <- c.cpuCount
}

func (c *systemCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/system/resource/print")
	if err != nil {
		log.Printf("[%s] system/resource: %v", dev.Name, err)
		return
	}
	for _, re := range reply.Re {
		l := []string{dev.Name, dev.Address}
		emit := func(desc *prometheus.Desc, vt prometheus.ValueType, key string) {
			if v, ok := parseFloat(re.Map[key]); ok {
				ch <- prometheus.MustNewConstMetric(desc, vt, v, l...)
			}
		}
		emit(c.cpuLoad, prometheus.GaugeValue, "cpu-load")
		emit(c.memFree, prometheus.GaugeValue, "free-memory")
		emit(c.memTotal, prometheus.GaugeValue, "total-memory")
		emit(c.hddFree, prometheus.GaugeValue, "free-hdd-space")
		emit(c.hddTotal, prometheus.GaugeValue, "total-hdd-space")
		emit(c.cpuCount, prometheus.GaugeValue, "cpu-count")
		if v := parseUptime(re.Map["uptime"]); v > 0 {
			ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.CounterValue, v, l...)
		}
	}
}
