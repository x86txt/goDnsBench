package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LoadServersFromCSV loads servers from a CSV file
// Expected format: name,dns,doh,dot,doq
func LoadServersFromCSV(filepath string) ([]Server, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return parseCSV(file)
}

// parseCSV parses a CSV reader into a list of servers
func parseCSV(r io.Reader) ([]Server, error) {
	reader := csv.NewReader(r)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	if len(header) < 1 {
		return nil, fmt.Errorf("CSV must have at least a name column")
	}

	var servers []Server

	// Read records
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		if len(record) == 0 {
			continue
		}

		server := Server{
			Name: strings.TrimSpace(record[0]),
		}

		// Handle optional columns
		if len(record) > 1 {
			server.DNS = strings.TrimSpace(record[1])
		}
		if len(record) > 2 {
			server.DoH = strings.TrimSpace(record[2])
		}
		if len(record) > 3 {
			server.DoT = strings.TrimSpace(record[3])
		}
		if len(record) > 4 {
			server.DoQ = strings.TrimSpace(record[4])
		}

		if err := server.Validate(); err != nil {
			return nil, fmt.Errorf("invalid server in CSV: %w", err)
		}

		servers = append(servers, server)
	}

	return servers, nil
}

// SaveServersToCSV saves servers to a CSV file
func SaveServersToCSV(filepath string, servers []Server) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"name", "dns", "doh", "dot", "doq"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write records
	for _, server := range servers {
		record := []string{
			server.Name,
			server.DNS,
			server.DoH,
			server.DoT,
			server.DoQ,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// FetchServersFromURL downloads a server list from a URL
func FetchServersFromURL(url string, timeout time.Duration) ([]Server, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch servers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Determine content type and parse accordingly
	contentType := resp.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") || strings.HasSuffix(url, ".json") {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var serverList ServerList
		if err := json.Unmarshal(data, &serverList); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}

		return serverList.Servers, nil
	} else if strings.Contains(contentType, "text/csv") || strings.HasSuffix(url, ".csv") {
		return parseCSV(resp.Body)
	}

	return nil, fmt.Errorf("unsupported content type: %s", contentType)
}
