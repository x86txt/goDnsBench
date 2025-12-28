// Filtering and sorting utilities for results

import type { BenchmarkResult } from '../api/types';

export type SortField = 'server' | 'protocol' | 'min' | 'mean' | 'p95' | 'p99' | 'successRate';
export type SortDirection = 'asc' | 'desc';

export interface FilterOptions {
  server?: string;
  protocol?: string;
  minLatency?: number;
  maxLatency?: number;
  minSuccessRate?: number;
}

export interface SortOptions {
  field: SortField;
  direction: SortDirection;
}

/**
 * Filters results based on filter options
 */
export function filterResults(
  results: BenchmarkResult[],
  filters: FilterOptions
): BenchmarkResult[] {
  return results.filter((result) => {
    // Server filter
    if (filters.server && !result.serverName.toLowerCase().includes(filters.server.toLowerCase())) {
      return false;
    }

    // Protocol filter
    if (filters.protocol && result.protocol !== filters.protocol) {
      return false;
    }

    // Latency filters (using mean latency)
    if (filters.minLatency !== undefined && result.metrics.mean < filters.minLatency) {
      return false;
    }
    if (filters.maxLatency !== undefined && result.metrics.mean > filters.maxLatency) {
      return false;
    }

    // Success rate filter
    if (filters.minSuccessRate !== undefined) {
      const successRate = (result.metrics.success / result.metrics.total) * 100;
      if (successRate < filters.minSuccessRate) {
        return false;
      }
    }

    return true;
  });
}

/**
 * Sorts results based on sort options
 */
export function sortResults(
  results: BenchmarkResult[],
  sort: SortOptions
): BenchmarkResult[] {
  const sorted = [...results];

  sorted.sort((a, b) => {
    let aVal: string | number;
    let bVal: string | number;

    switch (sort.field) {
      case 'server':
        aVal = a.serverName.toLowerCase();
        bVal = b.serverName.toLowerCase();
        break;
      case 'protocol':
        aVal = a.protocol.toLowerCase();
        bVal = b.protocol.toLowerCase();
        break;
      case 'min':
        aVal = a.metrics.min;
        bVal = b.metrics.min;
        break;
      case 'mean':
        aVal = a.metrics.mean;
        bVal = b.metrics.mean;
        break;
      case 'p95':
        aVal = a.metrics.p95;
        bVal = b.metrics.p95;
        break;
      case 'p99':
        aVal = a.metrics.p99;
        bVal = b.metrics.p99;
        break;
      case 'successRate':
        aVal = (a.metrics.success / a.metrics.total) * 100;
        bVal = (b.metrics.success / b.metrics.total) * 100;
        break;
      default:
        return 0;
    }

    if (typeof aVal === 'string' && typeof bVal === 'string') {
      if (sort.direction === 'asc') {
        return aVal.localeCompare(bVal);
      } else {
        return bVal.localeCompare(aVal);
      }
    } else {
      if (sort.direction === 'asc') {
        return (aVal as number) - (bVal as number);
      } else {
        return (bVal as number) - (aVal as number);
      }
    }
  });

  return sorted;
}

/**
 * Gets unique server names from results
 */
export function getUniqueServers(results: BenchmarkResult[]): string[] {
  const servers = new Set<string>();
  results.forEach((r) => servers.add(r.serverName));
  return Array.from(servers).sort();
}

/**
 * Gets unique protocols from results
 */
export function getUniqueProtocols(results: BenchmarkResult[]): string[] {
  const protocols = new Set<string>();
  results.forEach((r) => protocols.add(r.protocol));
  return Array.from(protocols).sort();
}
