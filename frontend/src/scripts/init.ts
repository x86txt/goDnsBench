import * as api from './api/client';
import type { BenchmarkResults, Server } from './api/types';
import { renderLatencyChart } from './charts/latency';
import { renderComparisonCharts } from './charts/comparison';
import { saveToHistory } from './ui/history';
import { showExportDialog, showImportDialog, readFileAsText, showSuccess, showError } from './ui/dialogs';
import { hideElement, showElement } from './ui/elements';
import type { FilterOptions } from './ui/filtering';
import { resetProgress, setProgress } from './ui/progress';
import { renderResultsTable, setFilters, clearFilters } from './ui/results';
import { getSelectedProtocols, getSelectedServers, renderServerList } from './ui/servers';
import { showSettingsModal, hideSettingsModal, populateSettingsForm, collectSettingsFromForm } from './ui/settings';

type ExportFormat = 'json' | 'csv';

class AppState {
  servers: Server[] = [];
  results: BenchmarkResults | null = null;
}

const state = new AppState();

async function loadServers(): Promise<void> {
  try {
    state.servers = await api.getServers();
    renderServerList(state.servers);
  } catch (err) {
    console.error('Failed to load servers:', err);
    state.servers = [];
    renderServerList(state.servers);
  }
}

async function runBenchmark(): Promise<void> {
  const protocols = getSelectedProtocols();
  const selectedServers = getSelectedServers(state.servers);

  if (protocols.length === 0) {
    alert('Please select at least one protocol');
    return;
  }

  if (selectedServers.length === 0) {
    alert('Please select at least one server');
    return;
  }

  showElement('results-container');
  showElement('progress-container');
  hideElement('welcome-message');
  hideElement('results-table-container');
  hideElement('charts-container');

  resetProgress();
  setProgress({ currentServer: 'Running benchmark...', percentage: 5 });

  try {
    // Set selected servers before running benchmark
    if (api.isWailsRuntime()) {
      const serverNames = selectedServers.map(s => s.name);
      await api.setSelectedServers(serverNames);
    }

    // Set up event listeners if in Wails runtime
    let cleanupProgress: (() => void) | null = null;
    let cleanupDone: (() => void) | null = null;
    let cleanupError: (() => void) | null = null;

    if (api.isWailsRuntime()) {
      cleanupProgress = api.onBenchmarkProgress((data) => {
        setProgress({
          currentServer: data.currentServer,
          currentProtocol: data.currentProtocol,
          completedTests: data.completedTests,
          totalTests: data.totalTests,
          percentage: data.percentage,
        });
      });

      cleanupDone = api.onBenchmarkDone((results) => {
        state.results = results;
        
        // Save to history
        saveToHistory(results);
        
        hideElement('progress-container');
        showElement('results-table-container');
        showElement('charts-container');
        renderResultsTable(results);
        void renderLatencyChart(results);
        void renderComparisonCharts(results);
        setupFilterHandlers();
        
        // Cleanup listeners
        if (cleanupProgress) cleanupProgress();
        if (cleanupDone) cleanupDone();
        if (cleanupError) cleanupError();
      });

      cleanupError = api.onBenchmarkError((data) => {
        throw new Error(data.error);
      });
    }

    // Run benchmark
    await api.runBenchmark(protocols);

    // If not using events (browser dev mode), poll for results
    if (!api.isWailsRuntime()) {
      const results = await api.getResults();
      if (!results || results.results.length === 0) {
        throw new Error('No results returned');
      }
      state.results = results;
      
      // Save to history
      saveToHistory(results);

      hideElement('progress-container');
      showElement('results-table-container');
      showElement('charts-container');

      renderResultsTable(results);
      await renderLatencyChart(results);
      await renderComparisonCharts(results);
      setupFilterHandlers();
    }
  } catch (err) {
    console.error('Benchmark failed:', err);
    alert(`Benchmark failed: ${err instanceof Error ? err.message : String(err)}`);
    hideElement('progress-container');
  }
}

async function refreshServers(): Promise<void> {
  try {
    await api.refreshServers();
    await loadServers();
    alert('Server list refreshed successfully');
  } catch (err) {
    console.error('Failed to refresh servers:', err);
    alert('Failed to refresh server list');
  }
}

async function exportResults(format: ExportFormat): Promise<void> {
  if (!state.results) {
    showError('No results to export. Please run a benchmark first.');
    return;
  }

  try {
    const defaultName = `goDnsBench_results_${Date.now()}.${format}`;
    const filepath = await showExportDialog(defaultName, format);
    if (!filepath) {
      return; // User cancelled
    }

    if (format === 'json') {
      await api.exportResultsJSON(filepath);
    } else {
      await api.exportResultsCSV(filepath);
    }
    showSuccess(`Results exported to ${filepath}`);
  } catch (err) {
    console.error('Export failed:', err);
    showError(`Export failed: ${err instanceof Error ? err.message : String(err)}`);
  }
}

async function showAddServerDialog(): Promise<void> {
  const name = prompt('Server name:');
  if (!name) return;

  const dns = prompt('DNS address (e.g., 1.1.1.1):') || '';
  const doh = prompt('DoH URL (optional):') || '';
  const dot = prompt('DoT address (optional):') || '';
  const doq = prompt('DoQ address (optional):') || '';

  const server: Server = { name, dns, doh, dot, doq };

  try {
    await api.addServer(server);

    if (api.isWailsRuntime()) {
      await loadServers();
    } else {
      state.servers.push(server);
      renderServerList(state.servers);
    }
  } catch (err) {
    console.error('Failed to add server:', err);
    alert(`Failed to add server: ${err instanceof Error ? err.message : String(err)}`);
  }
}

async function importServers(): Promise<void> {
  if (!api.isWailsRuntime()) {
    showError('Import functionality is only available in the desktop app');
    return;
  }

  try {
    // Show file picker dialog
    const file = await showImportDialog();
    if (!file) {
      return; // User cancelled
    }

    // Determine file type from extension
    const fileName = file.name.toLowerCase();
    const fileType = fileName.endsWith('.json') ? 'json' : 
                     fileName.endsWith('.csv') ? 'csv' : null;
    
    if (!fileType) {
      showError('Unsupported file format. Please select a .json or .csv file.');
      return;
    }

    // Read file content
    const content = await readFileAsText(file);

    // Call backend method
    await api.loadServersFromContent(content, fileType);

    // Reload servers and refresh UI
    await loadServers();
    showSuccess('Servers imported successfully');
  } catch (err) {
    console.error('Import failed:', err);
    showError(`Import failed: ${err instanceof Error ? err.message : String(err)}`);
  }
}

function setupEventListeners(): void {
  document.getElementById('run-benchmark-btn')?.addEventListener('click', () => {
    void runBenchmark();
  });

  document.getElementById('refresh-servers-btn')?.addEventListener('click', () => {
    void refreshServers();
  });

  document.getElementById('export-json-btn')?.addEventListener('click', () => {
    void exportResults('json');
  });

  document.getElementById('export-csv-btn')?.addEventListener('click', () => {
    void exportResults('csv');
  });

  document.getElementById('add-server-btn')?.addEventListener('click', () => {
    void showAddServerDialog();
  });

  document.getElementById('import-btn')?.addEventListener('click', () => {
    void importServers();
  });

  // Settings Button
  document.getElementById('settings-btn')?.addEventListener('click', () => {
    void openSettings();
  });

  // Settings Modal Buttons
  document.getElementById('settings-close-btn')?.addEventListener('click', () => {
    closeSettings();
  });

  document.getElementById('settings-cancel-btn')?.addEventListener('click', () => {
    closeSettings();
  });

  document.getElementById('settings-save-btn')?.addEventListener('click', () => {
    void saveSettings();
  });
}

function setupFilterHandlers(): void {
  // Filter button toggle
  const filterBtn = document.getElementById('filter-btn');
  const filterPanel = document.getElementById('filter-panel');
  
  filterBtn?.addEventListener('click', () => {
    if (filterPanel) {
      filterPanel.classList.toggle('hidden');
    }
  });

  // Apply filters
  document.getElementById('apply-filters-btn')?.addEventListener('click', () => {
    const filters: FilterOptions = {};
    
    const serverFilter = (document.getElementById('filter-server') as HTMLInputElement)?.value;
    if (serverFilter) {
      filters.server = serverFilter;
    }
    
    const protocolFilter = (document.getElementById('filter-protocol') as HTMLSelectElement)?.value;
    if (protocolFilter) {
      filters.protocol = protocolFilter;
    }
    
    const successRateFilter = (document.getElementById('filter-success-rate') as HTMLInputElement)?.value;
    if (successRateFilter) {
      filters.minSuccessRate = parseFloat(successRateFilter);
    }
    
    setFilters(filters);
    if (filterPanel) {
      filterPanel.classList.add('hidden');
    }
  });

  // Clear filters
  document.getElementById('clear-filters-btn')?.addEventListener('click', () => {
    (document.getElementById('filter-server') as HTMLInputElement).value = '';
    (document.getElementById('filter-protocol') as HTMLSelectElement).value = '';
    (document.getElementById('filter-success-rate') as HTMLInputElement).value = '';
    clearFilters();
  });
}

async function openSettings(): Promise<void> {
  try {
    // Load current settings
    const settings = await api.getSettings();

    // Populate form
    populateSettingsForm(settings);
    
    // Show modal
    showSettingsModal();
  } catch (err) {
    console.error('Failed to open settings:', err);
    showError(`Failed to load settings: ${err instanceof Error ? err.message : String(err)}`);
  }
}

function closeSettings(): void {
  hideSettingsModal();
}

async function saveSettings(): Promise<void> {
  try {
    // Collect settings from form
    const newSettings = collectSettingsFromForm();
    
    // Get current settings to preserve serverListUrl
    const currentSettings = await api.getSettings();
    newSettings.serverListUrl = currentSettings.serverListUrl || '';

    // Save settings
    await api.updateSettings(newSettings);
    showSuccess('Settings saved successfully');

    // Close modal
    hideSettingsModal();
  } catch (err) {
    console.error('Failed to save settings:', err);
    showError(`Failed to save settings: ${err instanceof Error ? err.message : String(err)}`);
  }
}

export async function init(): Promise<void> {
  await loadServers();
  setupEventListeners();
}

// Auto-init when loaded via <script src=".../init.ts"></script>
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    void init();
  });
} else {
  void init();
}

