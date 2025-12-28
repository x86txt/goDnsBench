package dns

import (
	"fmt"
	"time"
)

// ServerCapabilities represents which protocols a DNS server supports
type ServerCapabilities struct {
	DNS bool
	DoH bool
	DoT bool
	DoQ bool
}

// String returns a human-readable representation of capabilities
func (sc ServerCapabilities) String() string {
	protocols := []string{}
	if sc.DNS {
		protocols = append(protocols, "DNS")
	}
	if sc.DoH {
		protocols = append(protocols, "DoH")
	}
	if sc.DoT {
		protocols = append(protocols, "DoT")
	}
	if sc.DoQ {
		protocols = append(protocols, "DoQ")
	}

	if len(protocols) == 0 {
		return "No protocols supported"
	}

	result := ""
	for i, p := range protocols {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	return result
}

// CheckCapabilities tests which protocols a server supports
func CheckCapabilities(dnsAddr, dohURL, dotAddr, doqAddr string, timeout time.Duration) ServerCapabilities {
	caps := ServerCapabilities{}

	// Test standard DNS
	if dnsAddr != "" {
		client := NewStandardClient(dnsAddr, timeout)
		if err := client.TestConnection(); err == nil {
			caps.DNS = true
		}
	}

	// Test DoH
	if dohURL != "" {
		client := NewDoHClient(dohURL, timeout)
		if err := client.TestConnection(); err == nil {
			caps.DoH = true
		}
	}

	// Test DoT
	if dotAddr != "" {
		client := NewDoTClient(dotAddr, timeout)
		if err := client.TestConnection(); err == nil {
			caps.DoT = true
		}
	}

	// Test DoQ
	if doqAddr != "" {
		client := NewDoQClient(doqAddr, timeout)
		if err := client.TestConnection(); err == nil {
			caps.DoQ = true
		}
	}

	return caps
}

// ValidateServerConfig checks if a server configuration is valid
func ValidateServerConfig(dnsAddr, dohURL, dotAddr, doqAddr string) error {
	if dnsAddr == "" && dohURL == "" && dotAddr == "" && doqAddr == "" {
		return fmt.Errorf("at least one protocol endpoint must be provided")
	}
	return nil
}
