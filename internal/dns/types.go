package dns

import (
	"time"

	"github.com/miekg/dns"
)

// QueryType represents the type of DNS query
type QueryType int

const (
	QueryTypeA QueryType = iota
	QueryTypeMX
	QueryTypeTXT
	QueryTypeDNSSEC
)

// String returns the string representation of QueryType
func (q QueryType) String() string {
	switch q {
	case QueryTypeA:
		return "A"
	case QueryTypeMX:
		return "MX"
	case QueryTypeTXT:
		return "TXT"
	case QueryTypeDNSSEC:
		return "DNSSEC"
	default:
		return "UNKNOWN"
	}
}

// Protocol represents the DNS protocol type
type Protocol int

const (
	ProtocolDNS Protocol = iota
	ProtocolDoH
	ProtocolDoT
	ProtocolDoQ
)

// String returns the string representation of Protocol
func (p Protocol) String() string {
	switch p {
	case ProtocolDNS:
		return "DNS"
	case ProtocolDoH:
		return "DoH"
	case ProtocolDoT:
		return "DoT"
	case ProtocolDoQ:
		return "DoQ"
	default:
		return "UNKNOWN"
	}
}

// Query represents a single DNS query
type Query struct {
	Domain    string
	Type      QueryType
	Protocol  Protocol
	Timeout   time.Duration
}

// QueryResult represents the result of a DNS query
type QueryResult struct {
	Query         Query
	Success       bool
	Latency       time.Duration
	Error         error
	Response      *dns.Msg
	DNSSECValid   bool
}

// Client interface defines the contract for all DNS clients
type Client interface {
	Query(query Query) (*QueryResult, error)
	SupportsProtocol() Protocol
	Close() error
}
