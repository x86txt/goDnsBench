package benchmark

import (
	"math"
	"sort"
	"time"
)

// Metrics represents statistical metrics for benchmark results
type Metrics struct {
	Min     time.Duration
	Max     time.Duration
	Mean    time.Duration
	Median  time.Duration
	P95     time.Duration
	P99     time.Duration
	Success int
	Failed  int
	Total   int
}

// SuccessRate returns the success rate as a percentage
func (m *Metrics) SuccessRate() float64 {
	if m.Total == 0 {
		return 0
	}
	return float64(m.Success) / float64(m.Total) * 100
}

// FailureRate returns the failure rate as a percentage
func (m *Metrics) FailureRate() float64 {
	if m.Total == 0 {
		return 0
	}
	return float64(m.Failed) / float64(m.Total) * 100
}

// CalculateMetrics computes statistical metrics from a list of latencies
func CalculateMetrics(latencies []time.Duration, successCount, failedCount int) Metrics {
	metrics := Metrics{
		Success: successCount,
		Failed:  failedCount,
		Total:   successCount + failedCount,
	}

	if len(latencies) == 0 {
		return metrics
	}

	// Sort latencies for percentile calculations
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Min and Max
	metrics.Min = sorted[0]
	metrics.Max = sorted[len(sorted)-1]

	// Mean
	var sum time.Duration
	for _, d := range sorted {
		sum += d
	}
	metrics.Mean = time.Duration(int64(sum) / int64(len(sorted)))

	// Median
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		metrics.Median = (sorted[mid-1] + sorted[mid]) / 2
	} else {
		metrics.Median = sorted[mid]
	}

	// P95 (95th percentile)
	p95Idx := int(math.Ceil(float64(len(sorted)) * 0.95))
	if p95Idx >= len(sorted) {
		p95Idx = len(sorted) - 1
	}
	metrics.P95 = sorted[p95Idx]

	// P99 (99th percentile)
	p99Idx := int(math.Ceil(float64(len(sorted)) * 0.99))
	if p99Idx >= len(sorted) {
		p99Idx = len(sorted) - 1
	}
	metrics.P99 = sorted[p99Idx]

	return metrics
}

// CompareMetrics returns the percentage difference between two metrics
// Positive value means m1 is slower, negative means m1 is faster
func CompareMetrics(m1, m2 Metrics) float64 {
	if m2.Mean == 0 {
		return 0
	}
	return (float64(m1.Mean-m2.Mean) / float64(m2.Mean)) * 100
}
