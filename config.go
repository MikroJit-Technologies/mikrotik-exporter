package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenAddress string         `yaml:"listen_address"`
	Devices       []DeviceConfig `yaml:"devices"`
}

type DeviceConfig struct {
	Name       string   `yaml:"name"`
	Address    string   `yaml:"address"`
	User       string   `yaml:"user"`
	Password   string   `yaml:"password"`
	TLS        bool     `yaml:"tls"`
	SkipVerify bool     `yaml:"skip_verify"`
	Collectors []string `yaml:"collectors"` // empty = all enabled
}

func (d DeviceConfig) collectorEnabled(name string) bool {
	if len(d.Collectors) == 0 {
		return true
	}
	for _, c := range d.Collectors {
		if c == name {
			return true
		}
	}
	return false
}

func loadConfig() (*Config, error) {
	path := os.Getenv("CONFIG_FILE")
	if path == "" {
		path = "config.yml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{ListenAddress: ":9090"}
	return cfg, yaml.Unmarshal(data, cfg)
}
