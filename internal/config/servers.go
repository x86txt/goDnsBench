package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Server represents a DNS server configuration
type Server struct {
	Name string `json:"name"`
	DNS  string `json:"dns,omitempty"`
	DoH  string `json:"doh,omitempty"`
	DoT  string `json:"dot,omitempty"`
	DoQ  string `json:"doq,omitempty"`
}

// Validate checks if the server configuration is valid
func (s Server) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	if s.DNS == "" && s.DoH == "" && s.DoT == "" && s.DoQ == "" {
		return fmt.Errorf("server must have at least one protocol endpoint")
	}
	return nil
}

// ServerList represents a collection of DNS servers
type ServerList struct {
	Servers []Server `json:"servers"`
}

// LoadServersFromJSON loads servers from a JSON file
func LoadServersFromJSON(filepath string) ([]Server, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var serverList ServerList
	if err := json.Unmarshal(data, &serverList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate all servers
	for _, server := range serverList.Servers {
		if err := server.Validate(); err != nil {
			return nil, fmt.Errorf("invalid server %s: %w", server.Name, err)
		}
	}

	return serverList.Servers, nil
}

// SaveServersToJSON saves servers to a JSON file
func SaveServersToJSON(filepath string, servers []Server) error {
	serverList := ServerList{Servers: servers}

	data, err := json.MarshalIndent(serverList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// DefaultServers returns the default list of DNS servers
func DefaultServers() []Server {
	return []Server{
		{
			Name: "Cloudflare Primary",
			DNS:  "1.1.1.1",
			DoH:  "https://cloudflare-dns.com/dns-query",
			DoT:  "1.1.1.1:853",
			DoQ:  "1.1.1.1:8853",
		},
		{
			Name: "Cloudflare Secondary",
			DNS:  "1.0.0.1",
			DoH:  "https://cloudflare-dns.com/dns-query",
			DoT:  "1.0.0.1:853",
			DoQ:  "1.0.0.1:8853",
		},
		{
			Name: "Google Primary",
			DNS:  "8.8.8.8",
			DoH:  "https://dns.google/dns-query",
			DoT:  "8.8.8.8:853",
		},
		{
			Name: "Google Secondary",
			DNS:  "8.8.4.4",
			DoH:  "https://dns.google/dns-query",
			DoT:  "8.8.4.4:853",
		},
		{
			Name: "Quad9",
			DNS:  "9.9.9.9",
			DoH:  "https://dns.quad9.net/dns-query",
			DoT:  "9.9.9.9:853",
		},
		{
			Name: "OpenDNS Primary",
			DNS:  "208.67.222.222",
			DoH:  "https://doh.opendns.com/dns-query",
		},
		{
			Name: "OpenDNS Secondary",
			DNS:  "208.67.220.220",
			DoH:  "https://doh.opendns.com/dns-query",
		},
		{
			Name: "AdGuard DNS",
			DNS:  "94.140.14.14",
			DoH:  "https://dns.adguard.com/dns-query",
			DoT:  "94.140.14.14:853",
			DoQ:  "94.140.14.14:8853",
		},
		{
			Name: "NextDNS",
			DoH:  "https://dns.nextdns.io/dns-query",
			DoT:  "dns.nextdns.io:853",
		},
		{
			Name: "CleanBrowsing",
			DNS:  "185.228.168.9",
			DoH:  "https://doh.cleanbrowsing.org/doh/family-filter",
			DoT:  "185.228.168.9:853",
		},
		{
			Name: "Comodo Secure DNS",
			DNS:  "8.26.56.26",
		},
		{
			Name: "Level3 Primary",
			DNS:  "4.2.2.1",
		},
		{
			Name: "Level3 Secondary",
			DNS:  "4.2.2.2",
		},
	}
}
