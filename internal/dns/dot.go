package dns

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// DoTClient implements DNS over TLS (RFC 7858)
type DoTClient struct {
	server  string
	timeout time.Duration
}

// NewDoTClient creates a new DNS over TLS client
func NewDoTClient(server string, timeout time.Duration) *DoTClient {
	return &DoTClient{
		server:  server,
		timeout: timeout,
	}
}

// Query executes a DoT query and returns the result
func (c *DoTClient) Query(query Query) (*QueryResult, error) {
	msg := new(dns.Msg)

	// Set DNSSEC OK flag if this is a DNSSEC query
	if query.Type == QueryTypeDNSSEC {
		msg.SetEdns0(4096, true)
	}

	// Build the query
	qtype := c.getQueryType(query.Type)
	msg.SetQuestion(dns.Fqdn(query.Domain), qtype)

	// Execute the query and measure latency
	start := time.Now()

	client := &dns.Client{
		Net:     "tcp-tls",
		Timeout: query.Timeout,
		TLSConfig: &tls.Config{
			// Allow connection to any server for benchmarking purposes
			// In production, you might want to verify certificates
			InsecureSkipVerify: false,
		},
	}

	resp, _, err := client.Exchange(msg, c.server)
	latency := time.Since(start)

	result := &QueryResult{
		Query:    query,
		Latency:  latency,
		Response: resp,
	}

	if err != nil {
		result.Success = false
		result.Error = err
		return result, nil
	}

	result.Success = true

	// Check DNSSEC validity if requested
	if query.Type == QueryTypeDNSSEC && resp != nil {
		result.DNSSECValid = resp.AuthenticatedData
	}

	return result, nil
}

// SupportsProtocol returns the protocol this client supports
func (c *DoTClient) SupportsProtocol() Protocol {
	return ProtocolDoT
}

// Close closes the client (no-op for DoT)
func (c *DoTClient) Close() error {
	return nil
}

// getQueryType converts our QueryType to dns.Type
func (c *DoTClient) getQueryType(qt QueryType) uint16 {
	switch qt {
	case QueryTypeA, QueryTypeDNSSEC:
		return dns.TypeA
	case QueryTypeMX:
		return dns.TypeMX
	case QueryTypeTXT:
		return dns.TypeTXT
	default:
		return dns.TypeA
	}
}

// TestConnection tests if the DoT server is reachable
func (c *DoTClient) TestConnection() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	client := &dns.Client{
		Net:     "tcp-tls",
		Timeout: c.timeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	_, _, err := client.Exchange(msg, c.server)
	if err != nil {
		return fmt.Errorf("failed to connect to DoT server: %w", err)
	}

	return nil
}
