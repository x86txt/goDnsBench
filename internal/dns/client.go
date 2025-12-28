package dns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// StandardClient implements DNS queries over standard UDP/TCP
type StandardClient struct {
	server  string
	timeout time.Duration
}

// NewStandardClient creates a new standard DNS client
func NewStandardClient(server string, timeout time.Duration) *StandardClient {
	return &StandardClient{
		server:  server,
		timeout: timeout,
	}
}

// Query executes a DNS query and returns the result
func (c *StandardClient) Query(query Query) (*QueryResult, error) {
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
		Timeout: query.Timeout,
	}

	resp, _, err := client.Exchange(msg, c.server+":53")
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
func (c *StandardClient) SupportsProtocol() Protocol {
	return ProtocolDNS
}

// Close closes the client (no-op for standard DNS)
func (c *StandardClient) Close() error {
	return nil
}

// getQueryType converts our QueryType to dns.Type
func (c *StandardClient) getQueryType(qt QueryType) uint16 {
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

// TestConnection tests if the DNS server is reachable
func (c *StandardClient) TestConnection() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	client := &dns.Client{
		Timeout: c.timeout,
	}

	_, _, err := client.Exchange(msg, c.server+":53")
	if err != nil {
		return fmt.Errorf("failed to connect to DNS server: %w", err)
	}

	return nil
}
