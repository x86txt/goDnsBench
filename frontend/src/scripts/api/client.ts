import type { BenchmarkResults, Protocol, Server } from './types';

declare global {
  interface Window {
    go?: any;
    EventsOn?: (eventName: string, callback: (data?: any) => void) => void;
    EventsOff?: (eventName: string) => void;
  }
}

const PROTOCOL_BY_ENUM: Protocol[] = ['DNS', 'DoH', 'DoT', 'DoQ'];

function getWailsAppBinding(): any | undefined {
  if (typeof window === 'undefined') return undefined;
  return window.go?.main?.App;
}

export function isWailsRuntime(): boolean {
  return !!getWailsAppBinding();
}

function normalizeProtocol(value: unknown): Protocol {
  if (typeof value === 'string') {
    if (value === 'DNS' || value === 'DoH' || value === 'DoT' || value === 'DoQ') return value;
    return 'UNKNOWN';
  }
  if (typeof value === 'number') {
    return PROTOCOL_BY_ENUM[value] ?? 'UNKNOWN';
  }
  return 'UNKNOWN';
}

function toNumber(value: unknown): number {
  if (typeof value === 'number') return value;
  if (typeof value === 'string') {
    const n = Number(value);
    if (!Number.isNaN(n)) return n;
    const match = value.match(/-?\d+(\.\d+)?/);
    if (match) return Number(match[0]);
  }
  return 0;
}

// Convert a duration-like value into ms.
// Supports: ns numbers (Go time.Duration JSON), ms numbers (DTO), or strings.
function toMs(value: unknown): number {
  if (value == null) return 0;

  if (typeof value === 'number') {
    // Heuristic: ns values are typically >= 1e6 even for 1ms.
    return value > 100_000 ? value / 1_000_000 : value;
  }

  if (typeof value === 'string') {
    const trimmed = value.trim();
    const n = toNumber(trimmed);
    if (trimmed.endsWith('ms')) return n;
    if (trimmed.endsWith('ns')) return n / 1_000_000;
    if (trimmed.endsWith('s') && !trimmed.endsWith('ms')) return n * 1000;
    return n > 100_000 ? n / 1_000_000 : n;
  }

  return 0;
}

function normalizeServer(raw: any): Server {
  return {
    name: raw?.name ?? raw?.Name ?? '',
    dns: raw?.dns ?? raw?.DNS ?? '',
    doh: raw?.doh ?? raw?.DoH ?? '',
    dot: raw?.dot ?? raw?.DoT ?? '',
    doq: raw?.doq ?? raw?.DoQ ?? '',
  };
}

export async function getServers(): Promise<Server[]> {
  const app = getWailsAppBinding();
  if (app?.GetServers) {
    const raw = await app.GetServers();
    if (!Array.isArray(raw)) return [];
    return raw.map(normalizeServer);
  }
  return mockGetServers();
}

export async function refreshServers(): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.RefreshServerList) {
    await app.RefreshServerList();
    return;
  }
  await mockRefreshServers();
}

export async function addServer(server: Server): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.AddServer) {
    await app.AddServer(server);
    return;
  }
  // Browser dev fallback: no persistence, just succeed.
}

export async function loadServersFromFile(filepath: string): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.LoadServersFromFile) {
    await app.LoadServersFromFile(filepath);
    return;
  }
  throw new Error('LoadServersFromFile is only available in the desktop app');
}

export async function loadServersFromContent(content: string, fileType: string): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.LoadServersFromContent) {
    await app.LoadServersFromContent(content, fileType);
    return;
  }
  throw new Error('LoadServersFromContent is only available in the desktop app');
}

export async function setSelectedServers(serverNames: string[]): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.SetSelectedServers) {
    await app.SetSelectedServers(serverNames);
    return;
  }
  // Browser dev fallback: no-op
}

export async function runBenchmark(protocols: string[]): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.RunBenchmark) {
    await app.RunBenchmark(protocols);
    return;
  }
  await mockRunBenchmark(protocols);
}

function normalizeResults(raw: any): BenchmarkResults | null {
  if (!raw) return null;

  const rawResults = raw.results ?? raw.Results ?? [];
  const items = Array.isArray(rawResults) ? rawResults : [];

  const results = items.map((r: any) => {
    const metrics = r?.metrics ?? r?.Metrics ?? {};

    // DTOs use MinMs, MeanMs, etc. (already in milliseconds)
    // Fallback to old format (ns) for backward compatibility
    const hasMsFields = metrics.minMs !== undefined || metrics.MinMs !== undefined;
    
    const success = toNumber(metrics.success ?? metrics.Success);
    const failed = toNumber(metrics.failed ?? metrics.Failed);
    const total = toNumber(metrics.total ?? metrics.Total) || success + failed;

    return {
      serverName: r?.serverName ?? r?.ServerName ?? '',
      protocol: normalizeProtocol(r?.protocol ?? r?.Protocol),
      metrics: {
        min: hasMsFields 
          ? toNumber(metrics.minMs ?? metrics.MinMs)
          : toMs(metrics.min ?? metrics.Min),
        max: hasMsFields
          ? toNumber(metrics.maxMs ?? metrics.MaxMs)
          : toMs(metrics.max ?? metrics.Max),
        mean: hasMsFields
          ? toNumber(metrics.meanMs ?? metrics.MeanMs)
          : toMs(metrics.mean ?? metrics.Mean),
        median: hasMsFields
          ? toNumber(metrics.medianMs ?? metrics.MedianMs)
          : toMs(metrics.median ?? metrics.Median),
        p95: hasMsFields
          ? toNumber(metrics.p95Ms ?? metrics.P95Ms)
          : toMs(metrics.p95 ?? metrics.P95),
        p99: hasMsFields
          ? toNumber(metrics.p99Ms ?? metrics.P99Ms)
          : toMs(metrics.p99 ?? metrics.P99),
        success,
        failed,
        total,
      },
    };
  });

  const startTime = raw.startTime ?? raw.StartTime;
  const endTime = raw.endTime ?? raw.EndTime;
  const durationMs = raw.durationMs ?? raw.DurationMs;

  return {
    results,
    startTime: typeof startTime === 'string' ? startTime : undefined,
    endTime: typeof endTime === 'string' ? endTime : undefined,
    durationMs: typeof durationMs === 'number' ? durationMs : undefined,
  };
}

export async function getResults(): Promise<BenchmarkResults | null> {
  const app = getWailsAppBinding();
  if (app?.GetResults) {
    const raw = await app.GetResults();
    return normalizeResults(raw);
  }
  return normalizeResults(await mockGetResults());
}

export async function exportResultsJSON(filepath: string): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.ExportResultsJSON) {
    await app.ExportResultsJSON(filepath);
    return;
  }
  throw new Error('ExportResultsJSON is only available in the desktop app');
}

export async function exportResultsCSV(filepath: string): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.ExportResultsCSV) {
    await app.ExportResultsCSV(filepath);
    return;
  }
  throw new Error('ExportResultsCSV is only available in the desktop app');
}

export interface Settings {
  queryTimeoutMs: number; // milliseconds
  maxConcurrent: number;
  serverListUrl: string;
  lastServerUpdate?: string;
  enabledProtocols: string[];
  selectedDomains: {
    a: string[];
    mx: string[];
    txt: string[];
    dnssec: string[];
  };
}

export async function getSettings(): Promise<Settings> {
  const app = getWailsAppBinding();
  if (app?.GetSettings) {
    const raw = await app.GetSettings();
    return {
      queryTimeoutMs: raw?.queryTimeoutMs ?? raw?.QueryTimeoutMs ?? 1000,
      maxConcurrent: raw?.maxConcurrent ?? raw?.MaxConcurrent ?? 10,
      serverListUrl: raw?.serverListUrl ?? raw?.ServerListURL ?? '',
      lastServerUpdate: raw?.lastServerUpdate ?? raw?.LastServerUpdate,
      enabledProtocols: raw?.enabledProtocols ?? raw?.EnabledProtocols ?? ['DNS', 'DoH', 'DoT', 'DoQ'],
      selectedDomains: {
        a: raw?.selectedDomains?.a ?? raw?.SelectedDomains?.A ?? [],
        mx: raw?.selectedDomains?.mx ?? raw?.SelectedDomains?.MX ?? [],
        txt: raw?.selectedDomains?.txt ?? raw?.SelectedDomains?.TXT ?? [],
        dnssec: raw?.selectedDomains?.dnssec ?? raw?.SelectedDomains?.DNSSEC ?? [],
      },
    };
  }
  // Fallback defaults
  return {
    queryTimeoutMs: 1000,
    maxConcurrent: 10,
    serverListUrl: '',
    enabledProtocols: ['DNS', 'DoH', 'DoT', 'DoQ'],
    selectedDomains: {
      a: ['google.com', 'cloudflare.com', 'amazon.com', 'microsoft.com'],
      mx: ['gmail.com', 'microsoft.com'],
      txt: ['_dmarc.google.com', 'google.com'],
      dnssec: ['cloudflare.com', 'google.com'],
    },
  };
}

export async function updateSettings(settings: Settings): Promise<void> {
  const app = getWailsAppBinding();
  if (app?.UpdateSettings) {
    await app.UpdateSettings(settings);
    return;
  }
  throw new Error('UpdateSettings is only available in the desktop app');
}

// Wails event listener functions
export type ProgressCallback = (data: {
  currentServer: string;
  currentProtocol: string;
  completedTests: number;
  totalTests: number;
  percentage: number;
}) => void;

export type ResultsCallback = (data: BenchmarkResults) => void;

export type ErrorCallback = (data: { error: string }) => void;

export function onBenchmarkProgress(callback: ProgressCallback): () => void {
  if (typeof window !== 'undefined' && window.EventsOn) {
    window.EventsOn('benchmark:progress', callback);
    return () => {
      if (window.EventsOff) {
        window.EventsOff('benchmark:progress');
      }
    };
  }
  return () => {}; // No-op cleanup
}

export function onBenchmarkDone(callback: ResultsCallback): () => void {
  if (typeof window !== 'undefined' && window.EventsOn) {
    window.EventsOn('benchmark:done', (data: any) => {
      const normalized = normalizeResults(data);
      if (normalized) {
        callback(normalized);
      }
    });
    return () => {
      if (window.EventsOff) {
        window.EventsOff('benchmark:done');
      }
    };
  }
  return () => {}; // No-op cleanup
}

export function onBenchmarkError(callback: ErrorCallback): () => void {
  if (typeof window !== 'undefined' && window.EventsOn) {
    window.EventsOn('benchmark:error', callback);
    return () => {
      if (window.EventsOff) {
        window.EventsOff('benchmark:error');
      }
    };
  }
  return () => {}; // No-op cleanup
}

// -----------------------
// Mock fallbacks (browser)
// -----------------------

async function mockGetServers(): Promise<Server[]> {
  return [
    {
      name: 'Cloudflare Primary',
      dns: '1.1.1.1',
      doh: 'https://cloudflare-dns.com/dns-query',
      dot: '1.1.1.1:853',
      doq: '1.1.1.1:8853',
    },
    {
      name: 'Google Primary',
      dns: '8.8.8.8',
      doh: 'https://dns.google/dns-query',
      dot: '8.8.8.8:853',
      doq: '',
    },
  ];
}

async function mockRefreshServers(): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 300));
}

async function mockRunBenchmark(_protocols: string[]): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 1000));
}

async function mockGetResults(): Promise<any> {
  // Use ns-like values to exercise normalization.
  return {
    startTime: new Date().toISOString(),
    endTime: new Date().toISOString(),
    results: [
      {
        serverName: 'Cloudflare Primary',
        protocol: 'DNS',
        metrics: {
          min: 15_000_000,
          max: 25_000_000,
          mean: 20_000_000,
          median: 19_000_000,
          p95: 24_000_000,
          p99: 25_000_000,
          success: 10,
          failed: 0,
          total: 10,
        },
      },
      {
        serverName: 'Cloudflare Primary',
        protocol: 'DoH',
        metrics: {
          min: 25_000_000,
          max: 45_000_000,
          mean: 35_000_000,
          median: 34_000_000,
          p95: 42_000_000,
          p99: 44_000_000,
          success: 10,
          failed: 0,
          total: 10,
        },
      },
    ],
  };
}

