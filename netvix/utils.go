package main

import (
	"fmt"
)

type Metrics struct {
	URL        string
	IPv4       string
	StatusCode int
	DNS        float64
	TCP        float64
	TLS        float64
	TTFB       float64
	Total      float64
	BodySize   float64
	Redirects  int
}

func PrintMetricsTable(results []Metrics) {
	fmt.Printf("\n%-30s %-15s %-6s %-4s %-4s %-4s %-5s %-6s %-8s %-5s\n",
		"URL", "IPv4", "Code", "DNS", "TCP", "TLS", "TTFB", "Total", "Size(KB)", " ")

	for _, m := range results {
		fmt.Printf("%-30s %-15s %-6d %-4.0f %-4.0f %-4.0f %-5.0f %-6.0f %-8.1f\n",
			m.URL, m.IPv4, m.StatusCode,
			m.DNS, m.TCP, m.TLS, m.TTFB, m.Total, m.BodySize)
	}
	fmt.Println()
}
