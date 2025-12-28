package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/x86txt/goDnsBench/internal/benchmark"
)

// ExportResultsJSON exports benchmark results to a JSON file
func ExportResultsJSON(results *benchmark.BenchmarkResults, filepath string) error {
	// Create a simplified structure for export
	exportData := map[string]interface{}{
		"startTime": results.StartTime,
		"endTime":   results.EndTime,
		"duration":  results.Duration().String(),
		"results":   convertResults(results.Results),
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// convertResults converts ServerResults to a simplified format
func convertResults(results []benchmark.ServerResult) []map[string]interface{} {
	var converted []map[string]interface{}

	for _, result := range results {
		converted = append(converted, map[string]interface{}{
			"server":   result.ServerName,
			"protocol": result.Protocol.String(),
			"metrics": map[string]interface{}{
				"min":         result.Metrics.Min.String(),
				"max":         result.Metrics.Max.String(),
				"mean":        result.Metrics.Mean.String(),
				"median":      result.Metrics.Median.String(),
				"p95":         result.Metrics.P95.String(),
				"p99":         result.Metrics.P99.String(),
				"success":     result.Metrics.Success,
				"failed":      result.Metrics.Failed,
				"total":       result.Metrics.Total,
				"successRate": fmt.Sprintf("%.2f%%", result.Metrics.SuccessRate()),
			},
			"duration": result.Duration().String(),
		})
	}

	return converted
}
