package dns

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

// DoHClient implements DNS over HTTPS (RFC 8484)
type DoHClient struct {
	url     string
	timeout time.Duration
	client  *http.Client
}

// NewDoHClient creates a new DNS over HTTPS client
func NewDoHClient(url string, timeout time.Duration) *DoHClient {
	return &DoHClient{
		url:     url,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Query executes a DoH query and returns the result
func (c *DoHClient) Query(query Query) (*QueryResult, error) {
	msg := new(dns.Msg)

	// Set DNSSEC OK flag if this is a DNSSEC query
	if query.Type == QueryTypeDNSSEC {
		msg.SetEdns0(4096, true)
	}

	// Build the query
	qtype := c.getQueryType(query.Type)
	msg.SetQuestion(dns.Fqdn(query.Domain), qtype)

	// Pack the DNS message
	packed, err := msg.Pack()
	if err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to pack DNS message: %w", err),
		}, nil
	}

	// Execute the query and measure latency
	start := time.Now()

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(packed))
	if err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to create HTTP request: %w", err),
		}, nil
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := c.client.Do(req)
	latency := time.Since(start)

	result := &QueryResult{
		Query:   query,
		Latency: latency,
	}

	if err != nil {
		result.Success = false
		result.Error = err
		return result, nil
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		result.Success = false
		result.Error = fmt.Errorf("HTTP error: %d", resp.StatusCode)
		return result, nil
	}

	// Read and unpack the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to read response body: %w", err)
		return result, nil
	}

	dnsResp := new(dns.Msg)
	if err := dnsResp.Unpack(body); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("failed to unpack DNS response: %w", err)
		return result, nil
	}

	result.Success = true
	result.Response = dnsResp

	// Check DNSSEC validity if requested
	if query.Type == QueryTypeDNSSEC {
		result.DNSSECValid = dnsResp.AuthenticatedData
	}

	return result, nil
}

// SupportsProtocol returns the protocol this client supports
func (c *DoHClient) SupportsProtocol() Protocol {
	return ProtocolDoH
}

// Close closes the HTTP client
func (c *DoHClient) Close() error {
	c.client.CloseIdleConnections()
	return nil
}

// getQueryType converts our QueryType to dns.Type
func (c *DoHClient) getQueryType(qt QueryType) uint16 {
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

// TestConnection tests if the DoH server is reachable
func (c *DoHClient) TestConnection() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	packed, err := msg.Pack()
	if err != nil {
		return fmt.Errorf("failed to pack test message: %w", err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewReader(packed))
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to DoH server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DoH server returned status %d", resp.StatusCode)
	}

	return nil
}
