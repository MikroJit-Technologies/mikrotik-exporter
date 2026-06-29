package main

import (
	"log"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type routesCollector struct {
	total   *prometheus.Desc
	active  *prometheus.Desc
	dynamic *prometheus.Desc
	static  *prometheus.Desc
	bgp     *prometheus.Desc
	ospf    *prometheus.Desc
	connect *prometheus.Desc
}

func newRoutesCollector() subCollector {
	l := []string{"device", "address"}
	d := func(n, h string) *prometheus.Desc {
		return prometheus.NewDesc("mikrotik_routes_"+n, h, l, nil)
	}
	return &routesCollector{
		total:   d("total", "Total routes in routing table"),
		active:  d("active", "Active routes in routing table"),
		static:  d("static", "Static routes"),
		bgp:     d("bgp", "BGP routes"),
		ospf:    d("ospf", "OSPF routes"),
		dynamic: d("dynamic", "Dynamic routes (all protocols)"),
		connect: d("connect", "Connected routes"),
	}
}

func (c *routesCollector) name() string { return "routes" }

func (c *routesCollector) describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.active
	ch <- c.static
	ch <- c.bgp
	ch <- c.ospf
	ch <- c.dynamic
	ch <- c.connect
}

func (c *routesCollector) collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric) {
	reply, err := client.Run("/ip/route/print",
		"=.proplist=active,dynamic,static,bgp,ospf,connect")
	if err != nil {
		log.Printf("[%s] ip/route: %v", dev.Name, err)
		return
	}

	var total, active, dynamic, static, bgp, ospf, connect float64
	for _, re := range reply.Re {
		total++
		if re.Map["active"] == "true" {
			active++
		}
		if re.Map["static"] == "true" {
			static++
		}
		if re.Map["bgp"] == "true" {
			bgp++
		}
		if re.Map["ospf"] == "true" {
			ospf++
		}
		if re.Map["dynamic"] == "true" {
			dynamic++
		}
		if re.Map["connect"] == "true" {
			connect++
		}
	}

	l := []string{dev.Name, dev.Address}
	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, total, l...)
	ch <- prometheus.MustNewConstMetric(c.active, prometheus.GaugeValue, active, l...)
	ch <- prometheus.MustNewConstMetric(c.static, prometheus.GaugeValue, static, l...)
	ch <- prometheus.MustNewConstMetric(c.bgp, prometheus.GaugeValue, bgp, l...)
	ch <- prometheus.MustNewConstMetric(c.ospf, prometheus.GaugeValue, ospf, l...)
	ch <- prometheus.MustNewConstMetric(c.dynamic, prometheus.GaugeValue, dynamic, l...)
	ch <- prometheus.MustNewConstMetric(c.connect, prometheus.GaugeValue, connect, l...)
}
