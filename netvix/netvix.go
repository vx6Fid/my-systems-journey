package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

func MeasurePerformance(rawURL string) *Metrics {
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Invalid URL:", rawURL)
		return nil
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		fmt.Println("Failed to resolve IP for:", u.Hostname())
		return nil
	}

	var ipv4 string
	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4 = ip.String()
			break
		}
	}
	if ipv4 == "" {
		ipv4 = "N/A"
	}

	var (
		dnsStart, dnsDone             time.Time
		connStart, connDone           time.Time
		tlsStart, tlsDone             time.Time
		gotConn, gotFirstResponseByte time.Time
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone: func(httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart: func(_, _ string) { connStart = time.Now() },
		ConnectDone: func(_, _ string, _ error) { connDone = time.Now() },
		TLSHandshakeStart: func() { tlsStart = time.Now() },
		TLSHandshakeDone: func(tls.ConnectionState, error) { tlsDone = time.Now() },
		GotConn: func(httptrace.GotConnInfo) { gotConn = time.Now() },
		GotFirstResponseByte: func() { gotFirstResponseByte = time.Now() },
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		fmt.Println("Request creation failed:", rawURL)
		return nil
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	start := time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		fmt.Println("Request failed for:", rawURL)
		return nil
	}
	defer resp.Body.Close()

	total := time.Since(start)
	body, _ := io.ReadAll(resp.Body)
	sizeKB := float64(len(body)) / 1024

	return &Metrics{
		URL:        rawURL,
		IPv4:       ipv4,
		StatusCode: resp.StatusCode,
		DNS:        float64(dnsDone.Sub(dnsStart).Milliseconds()),
		TCP:        float64(connDone.Sub(connStart).Milliseconds()),
		TLS:        float64(tlsDone.Sub(tlsStart).Milliseconds()),
		TTFB:       float64(gotFirstResponseByte.Sub(gotConn).Milliseconds()),
		Total:      float64(total.Milliseconds()),
		BodySize:   sizeKB,
	}
}

