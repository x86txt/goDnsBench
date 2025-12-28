package benchmark

import (
	"math"
	"testing"
	"time"
)

func TestCalculateMetrics_EmptyLatencies(t *testing.T) {
	metrics := CalculateMetrics([]time.Duration{}, 5, 2)

	if metrics.Success != 5 {
		t.Errorf("Expected Success=5, got %d", metrics.Success)
	}
	if metrics.Failed != 2 {
		t.Errorf("Expected Failed=2, got %d", metrics.Failed)
	}
	if metrics.Total != 7 {
		t.Errorf("Expected Total=7, got %d", metrics.Total)
	}
	if metrics.Min != 0 {
		t.Errorf("Expected Min=0, got %v", metrics.Min)
	}
}

func TestCalculateMetrics_SingleLatency(t *testing.T) {
	latencies := []time.Duration{100 * time.Millisecond}
	metrics := CalculateMetrics(latencies, 1, 0)

	if metrics.Min != 100*time.Millisecond {
		t.Errorf("Expected Min=100ms, got %v", metrics.Min)
	}
	if metrics.Max != 100*time.Millisecond {
		t.Errorf("Expected Max=100ms, got %v", metrics.Max)
	}
	if metrics.Mean != 100*time.Millisecond {
		t.Errorf("Expected Mean=100ms, got %v", metrics.Mean)
	}
	if metrics.Median != 100*time.Millisecond {
		t.Errorf("Expected Median=100ms, got %v", metrics.Median)
	}
	if metrics.P95 != 100*time.Millisecond {
		t.Errorf("Expected P95=100ms, got %v", metrics.P95)
	}
	if metrics.P99 != 100*time.Millisecond {
		t.Errorf("Expected P99=100ms, got %v", metrics.P99)
	}
}

func TestCalculateMetrics_MultipleLatencies(t *testing.T) {
	latencies := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
	}
	metrics := CalculateMetrics(latencies, 5, 0)

	if metrics.Min != 10*time.Millisecond {
		t.Errorf("Expected Min=10ms, got %v", metrics.Min)
	}
	if metrics.Max != 50*time.Millisecond {
		t.Errorf("Expected Max=50ms, got %v", metrics.Max)
	}
	if metrics.Mean != 30*time.Millisecond {
		t.Errorf("Expected Mean=30ms, got %v", metrics.Mean)
	}
	if metrics.Median != 30*time.Millisecond {
		t.Errorf("Expected Median=30ms, got %v", metrics.Median)
	}
}

func TestCalculateMetrics_EvenCount(t *testing.T) {
	latencies := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
	}
	metrics := CalculateMetrics(latencies, 4, 0)

	expectedMedian := (20*time.Millisecond + 30*time.Millisecond) / 2
	if metrics.Median != expectedMedian {
		t.Errorf("Expected Median=%v, got %v", expectedMedian, metrics.Median)
	}
}

func TestCalculateMetrics_Percentiles(t *testing.T) {
	// Create 100 latencies for accurate percentile testing
	latencies := make([]time.Duration, 100)
	for i := 0; i < 100; i++ {
		latencies[i] = time.Duration(i+1) * time.Millisecond
	}
	metrics := CalculateMetrics(latencies, 100, 0)

	// P95 should be around index 95 (96th value)
	expectedP95 := 96 * time.Millisecond
	if metrics.P95 != expectedP95 {
		t.Errorf("Expected P95=%v, got %v", expectedP95, metrics.P95)
	}

	// P99 should be around index 99 (100th value)
	expectedP99 := 100 * time.Millisecond
	if metrics.P99 != expectedP99 {
		t.Errorf("Expected P99=%v, got %v", expectedP99, metrics.P99)
	}
}

func TestCalculateMetrics_PercentilesSmallSample(t *testing.T) {
	// Test with small sample size (10 items)
	latencies := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		latencies[i] = time.Duration(i+1) * time.Millisecond
	}
	metrics := CalculateMetrics(latencies, 10, 0)

	// P95 should clamp to last index (9)
	expectedP95 := 10 * time.Millisecond
	if metrics.P95 != expectedP95 {
		t.Errorf("Expected P95=%v, got %v", expectedP95, metrics.P95)
	}

	// P99 should also clamp to last index
	expectedP99 := 10 * time.Millisecond
	if metrics.P99 != expectedP99 {
		t.Errorf("Expected P99=%v, got %v", expectedP99, metrics.P99)
	}
}

func TestSuccessRate(t *testing.T) {
	tests := []struct {
		name     string
		metrics  Metrics
		expected float64
	}{
		{
			name:     "100% success",
			metrics:  Metrics{Success: 10, Failed: 0, Total: 10},
			expected: 100.0,
		},
		{
			name:     "50% success",
			metrics:  Metrics{Success: 5, Failed: 5, Total: 10},
			expected: 50.0,
		},
		{
			name:     "0% success",
			metrics:  Metrics{Success: 0, Failed: 10, Total: 10},
			expected: 0.0,
		},
		{
			name:     "Zero total",
			metrics:  Metrics{Success: 0, Failed: 0, Total: 0},
			expected: 0.0,
		},
		{
			name:     "Partial success",
			metrics:  Metrics{Success: 7, Failed: 3, Total: 10},
			expected: 70.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.SuccessRate()
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("Expected SuccessRate=%.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestFailureRate(t *testing.T) {
	tests := []struct {
		name     string
		metrics  Metrics
		expected float64
	}{
		{
			name:     "100% failure",
			metrics:  Metrics{Success: 0, Failed: 10, Total: 10},
			expected: 100.0,
		},
		{
			name:     "50% failure",
			metrics:  Metrics{Success: 5, Failed: 5, Total: 10},
			expected: 50.0,
		},
		{
			name:     "0% failure",
			metrics:  Metrics{Success: 10, Failed: 0, Total: 10},
			expected: 0.0,
		},
		{
			name:     "Zero total",
			metrics:  Metrics{Success: 0, Failed: 0, Total: 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.FailureRate()
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("Expected FailureRate=%.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestCompareMetrics(t *testing.T) {
	tests := []struct {
		name     string
		m1       Metrics
		m2       Metrics
		expected float64 // Positive = m1 slower, Negative = m1 faster
	}{
		{
			name:     "m1 is 50% slower",
			m1:       Metrics{Mean: 150 * time.Millisecond},
			m2:       Metrics{Mean: 100 * time.Millisecond},
			expected: 50.0,
		},
		{
			name:     "m1 is 50% faster",
			m1:       Metrics{Mean: 50 * time.Millisecond},
			m2:       Metrics{Mean: 100 * time.Millisecond},
			expected: -50.0,
		},
		{
			name:     "m1 is same speed",
			m1:       Metrics{Mean: 100 * time.Millisecond},
			m2:       Metrics{Mean: 100 * time.Millisecond},
			expected: 0.0,
		},
		{
			name:     "m2 is zero (should return 0)",
			m1:       Metrics{Mean: 100 * time.Millisecond},
			m2:       Metrics{Mean: 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareMetrics(tt.m1, tt.m2)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("Expected CompareMetrics=%.2f, got %.2f", tt.expected, result)
			}
		})
	}
}
