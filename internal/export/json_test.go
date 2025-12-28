package export

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/x86txt/goDnsBench/internal/benchmark"
	"github.com/x86txt/goDnsBench/internal/dns"
)

func createTestResults() *benchmark.BenchmarkResults {
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

func TestExportResultsJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "results.json")

	results := createTestResults()

	if err := ExportResultsJSON(results, jsonFile); err != nil {
		t.Fatalf("ExportResultsJSON() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Fatal("JSON file was not created")
	}

	// Read and verify content
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var exported map[string]interface{}
	if err := json.Unmarshal(data, &exported); err != nil {
		t.Fatalf("Failed to parse exported JSON: %v", err)
	}

	// Verify structure
	if _, ok := exported["startTime"]; !ok {
		t.Error("Missing 'startTime' field")
	}
	if _, ok := exported["endTime"]; !ok {
		t.Error("Missing 'endTime' field")
	}
	if _, ok := exported["duration"]; !ok {
		t.Error("Missing 'duration' field")
	}
	if _, ok := exported["results"]; !ok {
		t.Error("Missing 'results' field")
	}

	// Verify results array
	resultsArray, ok := exported["results"].([]interface{})
	if !ok {
		t.Fatal("Results field is not an array")
	}

	if len(resultsArray) != 2 {
		t.Errorf("Expected 2 results, got %d", len(resultsArray))
	}

	// Verify first result
	result1, ok := resultsArray[0].(map[string]interface{})
	if !ok {
		t.Fatal("First result is not an object")
	}

	if result1["server"] != "Test Server 1" {
		t.Errorf("Expected server='Test Server 1', got %v", result1["server"])
	}
	if result1["protocol"] != "DNS" {
		t.Errorf("Expected protocol='DNS', got %v", result1["protocol"])
	}

	metrics1, ok := result1["metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("Metrics is not an object")
	}

	if metrics1["success"] != float64(8) {
		t.Errorf("Expected success=8, got %v", metrics1["success"])
	}
	if metrics1["failed"] != float64(2) {
		t.Errorf("Expected failed=2, got %v", metrics1["failed"])
	}
	if metrics1["total"] != float64(10) {
		t.Errorf("Expected total=10, got %v", metrics1["total"])
	}
}

func TestExportResultsJSON_NilResults(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "results.json")

	err := ExportResultsJSON(nil, jsonFile)
	if err == nil {
		t.Error("Expected error for nil results")
	}
}

func TestExportResultsJSON_EmptyResults(t *testing.T) {
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "results.json")

	results := &benchmark.BenchmarkResults{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Results:   []benchmark.ServerResult{},
	}

	if err := ExportResultsJSON(results, jsonFile); err != nil {
		t.Fatalf("ExportResultsJSON() error = %v", err)
	}

	// File should still be created
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Fatal("JSON file was not created")
	}
}

func TestConvertResults(t *testing.T) {
	results := []benchmark.ServerResult{
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
				Success: 8,
				Failed:  2,
				Total:   10,
			},
			StartTime: time.Now(),
			EndTime:   time.Now().Add(1 * time.Second),
		},
	}

	converted := convertResults(results)

	if len(converted) != 1 {
		t.Errorf("Expected 1 converted result, got %d", len(converted))
	}

	result := converted[0]
	if result["server"] != "Test Server" {
		t.Errorf("Expected server='Test Server', got %v", result["server"])
	}
	if result["protocol"] != "DNS" {
		t.Errorf("Expected protocol='DNS', got %v", result["protocol"])
	}

	metrics, ok := result["metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("Metrics is not a map")
	}

	// convertResults returns int values, not float64
	success, ok := metrics["success"].(int)
	if !ok {
		t.Fatalf("Success is not an int, got type %T, value %v", metrics["success"], metrics["success"])
	}
	if success != 8 {
		t.Errorf("Expected success=8, got %v", success)
	}

	failed, ok := metrics["failed"].(int)
	if !ok {
		t.Fatalf("Failed is not an int, got type %T, value %v", metrics["failed"], metrics["failed"])
	}
	if failed != 2 {
		t.Errorf("Expected failed=2, got %v", failed)
	}

	total, ok := metrics["total"].(int)
	if !ok {
		t.Fatalf("Total is not an int, got type %T, value %v", metrics["total"], metrics["total"])
	}
	if total != 10 {
		t.Errorf("Expected total=10, got %v", total)
	}

	// Verify success rate is formatted as percentage string
	successRate, ok := metrics["successRate"].(string)
	if !ok {
		t.Fatal("SuccessRate is not a string")
	}
	if successRate != "80.00%" {
		t.Errorf("Expected successRate='80.00%%', got %s", successRate)
	}
}
