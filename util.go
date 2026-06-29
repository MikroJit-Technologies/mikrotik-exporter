package main

import (
	"regexp"
	"strconv"
	"strings"
)

var uptimeRe = regexp.MustCompile(`(?:(\d+)w)?(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?`)

func parseUptime(s string) float64 {
	m := uptimeRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return 0
	}
	mult := []float64{7 * 24 * 3600, 24 * 3600, 3600, 60, 1}
	var total float64
	for i, f := range mult {
		if m[i+1] != "" {
			v, _ := strconv.ParseFloat(m[i+1], 64)
			total += v * f
		}
	}
	return total
}

func parseFloat(s string) (float64, bool) {
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}

func boolToFloat(s string) float64 {
	if s == "true" {
		return 1
	}
	return 0
}

// splitPair splits "upload/download" RouterOS pair values e.g. "1000000/2000000"
func splitPair(s string, idx int) (float64, bool) {
	parts := strings.SplitN(s, "/", 2)
	if idx >= len(parts) {
		return 0, false
	}
	return parseFloat(parts[idx])
}
