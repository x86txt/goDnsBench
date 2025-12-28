import type { BenchmarkResults, BenchmarkResult } from '../api/types';
import type { FilterOptions, SortOptions } from './filtering';
import { filterResults, sortResults } from './filtering';

let currentResults: BenchmarkResult[] = [];
let currentFilters: FilterOptions = {};
let currentSort: SortOptions = { field: 'mean', direction: 'asc' };

function formatMs(ms: number): string {
  if (!Number.isFinite(ms)) return '0.00';
  return ms.toFixed(2);
}

function successRatePercent(success: number, total: number): number {
  if (!total) return 0;
  return (success / total) * 100;
}

function getSuccessRateClass(percent: number): string {
  if (percent >= 95) return 'text-green-400 font-semibold';
  if (percent >= 80) return 'text-yellow-400 font-semibold';
  return 'text-red-400 font-semibold';
}

function getSortIcon(field: string, currentField: string, direction: string): string {
  if (field !== currentField) {
    return '↕️';
  }
  return direction === 'asc' ? '↑' : '↓';
}

function renderTableHeader(sort: SortOptions): string {
  const sortIcon = (field: string) => getSortIcon(field, sort.field, sort.direction);
  
  return `
    <thead class="bg-lighter-gray">
      <tr>
        <th class="px-4 py-3 text-left text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="server">
          Server ${sortIcon('server')}
        </th>
        <th class="px-4 py-3 text-left text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="protocol">
          Protocol ${sortIcon('protocol')}
        </th>
        <th class="px-4 py-3 text-right text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="min">
          Min (ms) ${sortIcon('min')}
        </th>
        <th class="px-4 py-3 text-right text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="mean">
          Mean (ms) ${sortIcon('mean')}
        </th>
        <th class="px-4 py-3 text-right text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="p95">
          P95 (ms) ${sortIcon('p95')}
        </th>
        <th class="px-4 py-3 text-right text-sm font-semibold cursor-pointer hover:bg-lighter-gray/50" data-sort="successRate">
          Success Rate ${sortIcon('successRate')}
        </th>
      </tr>
    </thead>
  `;
}

export function renderResultsTable(results: BenchmarkResults, filters?: FilterOptions, sort?: SortOptions): void {
  const tbody = document.getElementById('results-tbody');
  const thead = document.querySelector('#results-table-container thead');
  if (!tbody) return;

  // Store current results
  currentResults = results.results;
  if (filters) currentFilters = filters;
  if (sort) currentSort = sort;

  // Apply filters and sorting
  let processed = [...currentResults];
  if (Object.keys(currentFilters).length > 0) {
    processed = filterResults(processed, currentFilters);
  }
  processed = sortResults(processed, currentSort);

  // Update header with sort indicators
  if (thead) {
    thead.outerHTML = renderTableHeader(currentSort);
    // Re-attach click handlers
    attachSortHandlers();
  }

  // Render table body
  tbody.innerHTML = processed
    .map((result) => {
      const rate = successRatePercent(result.metrics.success, result.metrics.total);
      const rateStr = `${rate.toFixed(2)}%`;

      return `
      <tr class="border-t border-lighter-gray hover:bg-lighter-gray/50">
        <td class="px-4 py-3 text-sm">${result.serverName}</td>
        <td class="px-4 py-3 text-sm">
          <span class="px-2 py-1 bg-electric-blue/20 text-electric-blue rounded text-xs font-semibold">
            ${result.protocol}
          </span>
        </td>
        <td class="px-4 py-3 text-sm text-right font-mono">${formatMs(result.metrics.min)}</td>
        <td class="px-4 py-3 text-sm text-right font-mono">${formatMs(result.metrics.mean)}</td>
        <td class="px-4 py-3 text-sm text-right font-mono">${formatMs(result.metrics.p95)}</td>
        <td class="px-4 py-3 text-sm text-right">
          <span class="${getSuccessRateClass(rate)}">
            ${rateStr}
          </span>
        </td>
      </tr>
    `;
    })
    .join('');

  // Update result count
  updateResultCount(processed.length, currentResults.length);
}

function attachSortHandlers(): void {
  const headers = document.querySelectorAll('#results-table-container th[data-sort]');
  headers.forEach((header) => {
    header.addEventListener('click', () => {
      const field = header.getAttribute('data-sort') as SortOptions['field'];
      if (!field) return;

      // Toggle sort direction if same field, otherwise set to asc
      if (currentSort.field === field) {
        currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
      } else {
        currentSort.field = field;
        currentSort.direction = 'asc';
      }

      // Re-render with new sort
      const results: BenchmarkResults = { results: currentResults };
      renderResultsTable(results, currentFilters, currentSort);
    });
  });
}

function updateResultCount(filtered: number, total: number): void {
  const countElement = document.getElementById('results-count');
  if (countElement) {
    if (filtered === total) {
      countElement.textContent = `Showing ${total} results`;
    } else {
      countElement.textContent = `Showing ${filtered} of ${total} results`;
    }
  }
}

export function setFilters(filters: FilterOptions): void {
  currentFilters = filters;
  const results: BenchmarkResults = { results: currentResults };
  renderResultsTable(results, currentFilters, currentSort);
}

export function clearFilters(): void {
  currentFilters = {};
  const results: BenchmarkResults = { results: currentResults };
  renderResultsTable(results, currentFilters, currentSort);
}

