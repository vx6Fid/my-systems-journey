package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: netviz <url1> <url2> ...")
		os.Exit(1)
	}

	var metricsList []Metrics
	for _, url := range os.Args[1:] {
		metrics := MeasurePerformance(url)
		if metrics != nil {
			metricsList = append(metricsList, *metrics)
		}
	}

	PrintMetricsTable(metricsList)
}

