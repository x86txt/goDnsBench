export type Protocol = 'DNS' | 'DoH' | 'DoT' | 'DoQ' | 'UNKNOWN';

export interface Server {
  name: string;
  dns?: string;
  doh?: string;
  dot?: string;
  doq?: string;
}

// All latency values are normalized to milliseconds (ms).
export interface LatencyMetrics {
  min: number;
  max: number;
  mean: number;
  median: number;
  p95: number;
  p99: number;
  success: number;
  failed: number;
  total: number;
}

export interface BenchmarkResult {
  serverName: string;
  protocol: Protocol;
  metrics: LatencyMetrics;
}

export interface BenchmarkResults {
  results: BenchmarkResult[];
  startTime?: string;
  endTime?: string;
  durationMs?: number;
}

export interface BenchmarkProgress {
  currentServer: string;
  currentProtocol: Protocol | string;
  completedTests: number;
  totalTests: number;
  percentage: number;
}

