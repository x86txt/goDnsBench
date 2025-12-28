// History management for benchmark results

import type { BenchmarkResults } from '../api/types';

const HISTORY_KEY = 'goDnsBench_history';
const MAX_HISTORY_ITEMS = 50;

export interface HistoryEntry {
  id: string;
  timestamp: string;
  results: BenchmarkResults;
}

/**
 * Saves results to history
 */
export function saveToHistory(results: BenchmarkResults): void {
  try {
    const history = loadHistory();
    const entry: HistoryEntry = {
      id: `benchmark_${Date.now()}`,
      timestamp: new Date().toISOString(),
      results,
    };

    // Add to beginning and limit size
    history.unshift(entry);
    if (history.length > MAX_HISTORY_ITEMS) {
      history.pop();
    }

    localStorage.setItem(HISTORY_KEY, JSON.stringify(history));
  } catch (err) {
    console.error('Failed to save to history:', err);
  }
}

/**
 * Loads history from localStorage
 */
export function loadHistory(): HistoryEntry[] {
  try {
    const stored = localStorage.getItem(HISTORY_KEY);
    if (!stored) return [];
    return JSON.parse(stored) as HistoryEntry[];
  } catch (err) {
    console.error('Failed to load history:', err);
    return [];
  }
}

/**
 * Clears all history
 */
export function clearHistory(): void {
  try {
    localStorage.removeItem(HISTORY_KEY);
  } catch (err) {
    console.error('Failed to clear history:', err);
  }
}

/**
 * Removes a specific history entry
 */
export function removeHistoryEntry(id: string): void {
  try {
    const history = loadHistory();
    const filtered = history.filter((entry) => entry.id !== id);
    localStorage.setItem(HISTORY_KEY, JSON.stringify(filtered));
  } catch (err) {
    console.error('Failed to remove history entry:', err);
  }
}

/**
 * Formats timestamp for display
 */
export function formatHistoryTimestamp(timestamp: string): string {
  try {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } catch {
    return timestamp;
  }
}
