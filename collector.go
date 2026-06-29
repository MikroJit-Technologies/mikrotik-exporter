package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/go-routeros/routeros"
	"github.com/prometheus/client_golang/prometheus"
)

type subCollector interface {
	name() string
	describe(ch chan<- *prometheus.Desc)
	collect(client *routeros.Client, dev DeviceConfig, ch chan<- prometheus.Metric)
}

type mikrotikCollector struct {
	cfg    *Config
	subs   []subCollector
	upDesc *prometheus.Desc
}

func newMikrotikCollector(cfg *Config) *mikrotikCollector {
	return &mikrotikCollector{
		cfg: cfg,
		subs: []subCollector{
			newSystemCollector(),
			newInterfaceCollector(),
			newBGPCollector(),
			newOSPFCollector(),
			newDHCPCollector(),
			newFirewallCollector(),
			newQueueCollector(),
			newWireGuardCollector(),
			newCapsManCollector(),
			newHealthCollector(),
		},
		upDesc: prometheus.NewDesc(
			"mikrotik_up",
			"Whether the MikroTik device is reachable (1=up, 0=down)",
			[]string{"device", "address"}, nil,
		),
	}
}

func (c *mikrotikCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upDesc
	for _, s := range c.subs {
		s.describe(ch)
	}
}

func (c *mikrotikCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	for _, dev := range c.cfg.Devices {
		wg.Add(1)
		go func(d DeviceConfig) {
			defer wg.Done()
			c.collectDevice(d, ch)
		}(dev)
	}
	wg.Wait()
}

func (c *mikrotikCollector) collectDevice(dev DeviceConfig, ch chan<- prometheus.Metric) {
	client, err := dial(dev)
	if err != nil {
		log.Printf("[%s] connect: %v", dev.Name, err)
		ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 0, dev.Name, dev.Address)
		return
	}
	defer client.Close()

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, 1, dev.Name, dev.Address)

	for _, s := range c.subs {
		if dev.collectorEnabled(s.name()) {
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[%s] %s panic: %v", dev.Name, s.name(), r)
					}
				}()
				s.collect(client, dev, ch)
			}()
		}
	}
}

func dial(dev DeviceConfig) (*routeros.Client, error) {
	if dev.TLS {
		return routeros.DialTLS(dev.Address, dev.User, dev.Password, &tls.Config{
			InsecureSkipVerify: dev.SkipVerify,
		})
	}
	client, err := routeros.Dial(dev.Address, dev.User, dev.Password)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", dev.Address, err)
	}
	return client, nil
}
