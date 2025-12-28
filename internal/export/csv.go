package export

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/x86txt/goDnsBench/internal/benchmark"
)

// ExportResultsCSV exports benchmark results to a CSV file
func ExportResultsCSV(results *benchmark.BenchmarkResults, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
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

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write results
	for _, result := range results.Results {
		record := []string{
			result.ServerName,
			result.Protocol.String(),
			fmt.Sprintf("%.2f", float64(result.Metrics.Min.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.Metrics.Max.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.Metrics.Mean.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.Metrics.Median.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.Metrics.P95.Microseconds())/1000.0),
			fmt.Sprintf("%.2f", float64(result.Metrics.P99.Microseconds())/1000.0),
			fmt.Sprintf("%d", result.Metrics.Success),
			fmt.Sprintf("%d", result.Metrics.Failed),
			fmt.Sprintf("%d", result.Metrics.Total),
			fmt.Sprintf("%.2f", result.Metrics.SuccessRate()),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}
