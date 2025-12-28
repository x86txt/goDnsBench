package export

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/dns"
)

func createTestResultsForCSV() *benchmark.BenchmarkResults {
	startTime := time.Now()
	endTime := startTime.Add(5 * time.Second)

	return &benchmark.BenchmarkResults{
		StartTime: startTime,
		EndTime:   endTime,
		Results: []benchmark.ServerResult{
			{
				ServerName: "Test Server 1",
				Protocol:   dns.ProtocolDNS,
				Metrics: benchmark.Metrics{
					Min:     10 * time.Millisecond,
					Max:     50 * time.Millisecond,
					Mean:    30 * time.Millisecond,
					Median:  30 * time.Millisecond,
					P95:     45 * time.Millisecond,
					P99:     50 * time.Millisecond,
					Success: 8,
					Failed:  2,
					Total:   10,
				},
				StartTime: startTime,
				EndTime:   startTime.Add(1 * time.Second),
			},
			{
				ServerName: "Test Server 2",
				Protocol:   dns.ProtocolDoH,
				Metrics: benchmark.Metrics{
					Min:     20 * time.Millisecond,
					Max:     60 * time.Millisecond,
					Mean:    40 * time.Millisecond,
					Median:  40 * time.Millisecond,
					P95:     55 * time.Millisecond,
					P99:     60 * time.Millisecond,
					Success: 10,
					Failed:  0,
					Total:   10,
				},
				StartTime: startTime.Add(1 * time.Second),
				EndTime:   startTime.Add(2 * time.Second),
			},
		},
	}
}

func TestExportResultsCSV(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "results.csv")

	results := createTestResultsForCSV()

	if err := ExportResultsCSV(results, csvFile); err != nil {
		t.Fatalf("ExportResultsCSV() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Read and verify content
	file, err := os.Open(csvFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	// Verify header
	expectedHeader := []string{
		"Server",
		"Protocol",
		"Min (ms)",
		"Max (ms)",
		"Mean (ms)",
		"Median (ms)",
		"P95 (ms)",
		"P99 (ms)",
		"Success",
		"Failed",
		"Total",
		"Success Rate (%)",
	}

	if len(records) < 2 {
		t.Fatalf("Expected at least header + 2 data rows, got %d rows", len(records))
	}

	header := records[0]
	if len(header) != len(expectedHeader) {
		t.Errorf("Header length mismatch: expected %d, got %d", len(expectedHeader), len(header))
	}

	for i, expected := range expectedHeader {
		if header[i] != expected {
			t.Errorf("Header[%d] mismatch: expected %s, got %s", i, expected, header[i])
		}
	}

	// Verify first data row
	row1 := records[1]
	if len(row1) != len(expectedHeader) {
		t.Errorf("Row 1 length mismatch: expected %d, got %d", len(expectedHeader), len(row1))
	}

	if row1[0] != "Test Server 1" {
		t.Errorf("Expected server='Test Server 1', got %s", row1[0])
	}
	if row1[1] != "DNS" {
		t.Errorf("Expected protocol='DNS', got %s", row1[1])
	}

	// Verify numeric values
	success, err := strconv.Atoi(row1[8])
	if err != nil {
		t.Fatalf("Failed to parse Success: %v", err)
	}
	if success != 8 {
		t.Errorf("Expected Success=8, got %d", success)
	}

	failed, err := strconv.Atoi(row1[9])
	if err != nil {
		t.Fatalf("Failed to parse Failed: %v", err)
	}
	if failed != 2 {
		t.Errorf("Expected Failed=2, got %d", failed)
	}

	total, err := strconv.Atoi(row1[10])
	if err != nil {
		t.Fatalf("Failed to parse Total: %v", err)
	}
	if total != 10 {
		t.Errorf("Expected Total=10, got %d", total)
	}

	// Verify second data row
	if len(records) < 3 {
		t.Fatal("Expected at least 3 rows (header + 2 data)")
	}

	row2 := records[2]
	if row2[0] != "Test Server 2" {
		t.Errorf("Expected server='Test Server 2', got %s", row2[0])
	}
	if row2[1] != "DoH" {
		t.Errorf("Expected protocol='DoH', got %s", row2[1])
	}

	success2, err := strconv.Atoi(row2[8])
	if err != nil {
		t.Fatalf("Failed to parse Success: %v", err)
	}
	if success2 != 10 {
		t.Errorf("Expected Success=10, got %d", success2)
	}
}

func TestExportResultsCSV_NilResults(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "results.csv")

	err := ExportResultsCSV(nil, csvFile)
	if err == nil {
		t.Error("Expected error for nil results")
	}
	if err.Error() != "no results to export" {
		t.Errorf("Expected error message 'no results to export', got %v", err)
	}
}

func TestExportResultsCSV_EmptyResults(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "results.csv")

	results := &benchmark.BenchmarkResults{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Results:   []benchmark.ServerResult{},
	}

	if err := ExportResultsCSV(results, csvFile); err != nil {
		t.Fatalf("ExportResultsCSV() error = %v", err)
	}

	// File should still be created with just header
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		t.Fatal("CSV file was not created")
	}

	// Verify file has header only
	file, err := os.Open(csvFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if len(records) != 1 {
		t.Errorf("Expected 1 row (header only), got %d", len(records))
	}
}

func TestExportResultsCSV_LatencyFormatting(t *testing.T) {
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "results.csv")

	results := &benchmark.BenchmarkResults{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Results: []benchmark.ServerResult{
			{
				ServerName: "Test Server",
				Protocol:   dns.ProtocolDNS,
				Metrics: benchmark.Metrics{
					Min:     10 * time.Millisecond,
					Max:     50 * time.Millisecond,
					Mean:    30 * time.Millisecond,
					Median:  25 * time.Millisecond,
					P95:     45 * time.Millisecond,
					P99:     50 * time.Millisecond,
					Success: 10,
					Failed:  0,
					Total:   10,
				},
			},
		},
	}

	if err := ExportResultsCSV(results, csvFile); err != nil {
		t.Fatalf("ExportResultsCSV() error = %v", err)
	}

	// Read and verify latency values are formatted correctly
	file, err := os.Open(csvFile)
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV file: %v", err)
	}

	if len(records) < 2 {
		t.Fatal("Expected at least header + 1 data row")
	}

	row := records[1]
	// Min should be 10.00 ms
	if row[2] != "10.00" {
		t.Errorf("Expected Min='10.00', got %s", row[2])
	}
	// Max should be 50.00 ms
	if row[3] != "50.00" {
		t.Errorf("Expected Max='50.00', got %s", row[3])
	}
	// Mean should be 30.00 ms
	if row[4] != "30.00" {
		t.Errorf("Expected Mean='30.00', got %s", row[4])
	}
}
