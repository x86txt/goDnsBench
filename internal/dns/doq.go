package dns

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"github.com/miekg/dns"
	"github.com/quic-go/quic-go"
)

// DoQClient implements DNS over QUIC (RFC 9250)
type DoQClient struct {
	server  string
	timeout time.Duration
}

// NewDoQClient creates a new DNS over QUIC client
func NewDoQClient(server string, timeout time.Duration) *DoQClient {
	return &DoQClient{
		server:  server,
		timeout: timeout,
	}
}

// Query executes a DoQ query and returns the result
func (c *DoQClient) Query(query Query) (*QueryResult, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), query.Timeout)
	defer cancel()

	tlsConf := &tls.Config{
		NextProtos: []string{"doq"},
		// For benchmarking, we'll accept any certificate
		// In production, consider proper certificate validation
		InsecureSkipVerify: false,
	}

	conn, err := quic.DialAddr(ctx, c.server, tlsConf, &quic.Config{
		MaxIdleTimeout: query.Timeout,
	})

	if err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to establish QUIC connection: %w", err),
			Latency: time.Since(start),
		}, nil
	}
	defer conn.CloseWithError(0, "")

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to open QUIC stream: %w", err),
			Latency: time.Since(start),
		}, nil
	}
	defer stream.Close()

	// Write DNS message with length prefix (2 bytes, big-endian)
	length := uint16(len(packed))
	lenBytes := []byte{byte(length >> 8), byte(length & 0xff)}

	if _, err := stream.Write(append(lenBytes, packed...)); err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to write query: %w", err),
			Latency: time.Since(start),
		}, nil
	}

	// Read response length prefix
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(stream, lenBuf); err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to read response length: %w", err),
			Latency: time.Since(start),
		}, nil
	}

	respLen := int(lenBuf[0])<<8 | int(lenBuf[1])

	// Read response
	respBuf := make([]byte, respLen)
	if _, err := io.ReadFull(stream, respBuf); err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to read response: %w", err),
			Latency: time.Since(start),
		}, nil
	}

	latency := time.Since(start)

	// Unpack response
	dnsResp := new(dns.Msg)
	if err := dnsResp.Unpack(respBuf); err != nil {
		return &QueryResult{
			Query:   query,
			Success: false,
			Error:   fmt.Errorf("failed to unpack DNS response: %w", err),
			Latency: latency,
		}, nil
	}

	result := &QueryResult{
		Query:    query,
		Success:  true,
		Latency:  latency,
		Response: dnsResp,
	}

	// Check DNSSEC validity if requested
	if query.Type == QueryTypeDNSSEC {
		result.DNSSECValid = dnsResp.AuthenticatedData
	}

	return result, nil
}

// SupportsProtocol returns the protocol this client supports
func (c *DoQClient) SupportsProtocol() Protocol {
	return ProtocolDoQ
}

// Close closes the client (no-op for DoQ)
func (c *DoQClient) Close() error {
	return nil
}

// getQueryType converts our QueryType to dns.Type
func (c *DoQClient) getQueryType(qt QueryType) uint16 {
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

// TestConnection tests if the DoQ server is reachable
func (c *DoQClient) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tlsConf := &tls.Config{
		NextProtos:         []string{"doq"},
		InsecureSkipVerify: false,
	}

	conn, err := quic.DialAddr(ctx, c.server, tlsConf, &quic.Config{
		MaxIdleTimeout: c.timeout,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to DoQ server: %w", err)
	}

	conn.CloseWithError(0, "")
	return nil
}
