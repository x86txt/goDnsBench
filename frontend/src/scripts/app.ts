// This file will be updated with actual Wails bindings after wails generate
// For now, we'll create the structure and mock the API calls

interface Server {
  name: string;
  dns: string;
  doh: string;
  dot: string;
  doq: string;
}

interface BenchmarkResult {
  serverName: string;
  protocol: string;
  metrics: {
    min: string;
    max: string;
    mean: string;
    median: string;
    p95: string;
    p99: string;
    success: number;
    failed: number;
    total: number;
    successRate: string;
  };
}

interface BenchmarkResults {
  results: BenchmarkResult[];
  startTime: string;
  endTime: string;
  duration: string;
}

class App {
  private servers: Server[] = [];
  private results: BenchmarkResults | null = null;
  private chart: any = null;

  async init() {
    await this.loadServers();
    this.setupEventListeners();
    this.renderServerList();
  }

  async loadServers() {
    try {
      // Will be replaced with: this.servers = await window.go.main.App.GetServers();
      this.servers = await this.mockGetServers();
    } catch (error) {
      console.error('Failed to load servers:', error);
    }
  }

  setupEventListeners() {
    // Run Benchmark Button
    document.getElementById('run-benchmark-btn')?.addEventListener('click', () => {
      this.runBenchmark();
    });

    // Refresh Servers Button
    document.getElementById('refresh-servers-btn')?.addEventListener('click', () => {
      this.refreshServers();
    });

    // Export Buttons
    document.getElementById('export-json-btn')?.addEventListener('click', () => {
      this.exportResults('json');
    });

    document.getElementById('export-csv-btn')?.addEventListener('click', () => {
      this.exportResults('csv');
    });

    // Add Server Button
    document.getElementById('add-server-btn')?.addEventListener('click', () => {
      this.showAddServerDialog();
    });

    // Import Button
    document.getElementById('import-btn')?.addEventListener('click', () => {
      this.importServers();
    });
  }

  renderServerList() {
    const container = document.getElementById('server-list');
    if (!container) return;

    container.innerHTML = this.servers.map(server => `
      <div class="flex items-center justify-between px-3 py-2 rounded hover:bg-lighter-gray transition">
        <label class="flex items-center space-x-2 cursor-pointer flex-1">
          <input type="checkbox" class="form-checkbox text-electric-blue server-checkbox" data-server="${server.name}" checked>
          <span class="text-sm">${server.name}</span>
        </label>
      </div>
    `).join('');
  }

  getSelectedProtocols(): string[] {
    const protocols: string[] = [];
    if ((document.getElementById('protocol-dns') as HTMLInputElement)?.checked) protocols.push('DNS');
    if ((document.getElementById('protocol-doh') as HTMLInputElement)?.checked) protocols.push('DoH');
    if ((document.getElementById('protocol-dot') as HTMLInputElement)?.checked) protocols.push('DoT');
    if ((document.getElementById('protocol-doq') as HTMLInputElement)?.checked) protocols.push('DoQ');
    return protocols;
  }

  getSelectedServers(): Server[] {
    const checkboxes = document.querySelectorAll('.server-checkbox:checked');
    const selectedNames = Array.from(checkboxes).map(cb => (cb as HTMLInputElement).dataset.server);
    return this.servers.filter(s => selectedNames.includes(s.name));
  }

  async runBenchmark() {
    const protocols = this.getSelectedProtocols();
    const selectedServers = this.getSelectedServers();

    if (protocols.length === 0) {
      alert('Please select at least one protocol');
      return;
    }

    if (selectedServers.length === 0) {
      alert('Please select at least one server');
      return;
    }

    // Show results container and progress
    this.showElement('results-container');
    this.showElement('progress-container');
    this.hideElement('welcome-message');
    this.hideElement('results-table-container');
    this.hideElement('charts-container');

    try {
      // Will be replaced with: await window.go.main.App.RunBenchmark(protocols);
      await this.mockRunBenchmark(protocols);

      // Get results
      // Will be replaced with: this.results = await window.go.main.App.GetResults();
      this.results = await this.mockGetResults();

      // Hide progress, show results
      this.hideElement('progress-container');
      this.showElement('results-table-container');
      this.showElement('charts-container');

      this.renderResults();
      this.renderChart();
    } catch (error) {
      console.error('Benchmark failed:', error);
      alert('Benchmark failed: ' + error);
    }
  }

  renderResults() {
    if (!this.results) return;

    const tbody = document.getElementById('results-tbody');
    if (!tbody) return;

    tbody.innerHTML = this.results.results.map(result => `
      <tr class="border-t border-lighter-gray hover:bg-lighter-gray/50">
        <td class="px-4 py-3 text-sm">${result.serverName}</td>
        <td class="px-4 py-3 text-sm">
          <span class="px-2 py-1 bg-electric-blue/20 text-electric-blue rounded text-xs font-semibold">
            ${result.protocol}
          </span>
        </td>
        <td class="px-4 py-3 text-sm text-right font-mono">${this.formatLatency(result.metrics.min)}</td>
        <td class="px-4 py-3 text-sm text-right font-mono">${this.formatLatency(result.metrics.mean)}</td>
        <td class="px-4 py-3 text-sm text-right font-mono">${this.formatLatency(result.metrics.p95)}</td>
        <td class="px-4 py-3 text-sm text-right">
          <span class="${this.getSuccessRateClass(result.metrics.successRate)}">
            ${result.metrics.successRate}
          </span>
        </td>
      </tr>
    `).join('');
  }

  formatLatency(latencyStr: string): string {
    // Convert nanosecond strings to ms
    const match = latencyStr.match(/(\d+\.?\d*)/);
    if (match) {
      const ns = parseFloat(match[1]);
      return (ns / 1000000).toFixed(2);
    }
    return latencyStr;
  }

  getSuccessRateClass(rate: string): string {
    const percentage = parseFloat(rate);
    if (percentage >= 95) return 'text-green-400 font-semibold';
    if (percentage >= 80) return 'text-yellow-400 font-semibold';
    return 'text-red-400 font-semibold';
  }

  async renderChart() {
    if (!this.results) return;

    const canvas = document.getElementById('latency-chart') as HTMLCanvasElement;
    if (!canvas) return;

    // Import Chart.js dynamically
    const { Chart, registerables } = await import('chart.js');
    Chart.register(...registerables);

    // Prepare data
    const labels = this.results.results.map(r => `${r.serverName} (${r.protocol})`);
    const data = this.results.results.map(r => {
      const match = r.metrics.mean.match(/(\d+\.?\d*)/);
      return match ? parseFloat(match[1]) / 1000000 : 0; // Convert ns to ms
    });

    // Destroy existing chart
    if (this.chart) {
      this.chart.destroy();
    }

    // Create new chart
    this.chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels,
        datasets: [{
          label: 'Mean Latency (ms)',
          data,
          backgroundColor: '#0066ff',
          borderColor: '#00ffff',
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        scales: {
          y: {
            beginAtZero: true,
            grid: {
              color: '#3a3a3a'
            },
            ticks: {
              color: '#ffffff'
            }
          },
          x: {
            grid: {
              color: '#3a3a3a'
            },
            ticks: {
              color: '#ffffff',
              maxRotation: 45,
              minRotation: 45
            }
          }
        },
        plugins: {
          legend: {
            labels: {
              color: '#ffffff'
            }
          }
        }
      }
    });
  }

  async refreshServers() {
    try {
      // Will be replaced with: await window.go.main.App.RefreshServerList();
      await this.mockRefreshServers();
      await this.loadServers();
      this.renderServerList();
      alert('Server list refreshed successfully');
    } catch (error) {
      console.error('Failed to refresh servers:', error);
      alert('Failed to refresh server list');
    }
  }

  async exportResults(format: 'json' | 'csv') {
    if (!this.results) {
      alert('No results to export');
      return;
    }

    try {
      // Will be replaced with Wails file dialog
      const filename = `goDnsBench_results_${Date.now()}.${format}`;

      if (format === 'json') {
        // await window.go.main.App.ExportResultsJSON(filename);
        console.log('Export JSON:', filename);
      } else {
        // await window.go.main.App.ExportResultsCSV(filename);
        console.log('Export CSV:', filename);
      }

      alert(`Results exported to ${filename}`);
    } catch (error) {
      console.error('Export failed:', error);
      alert('Export failed');
    }
  }

  showAddServerDialog() {
    // This will be implemented with a proper modal
    const name = prompt('Server name:');
    if (!name) return;

    const dns = prompt('DNS address (e.g., 1.1.1.1):') || '';
    const doh = prompt('DoH URL (optional):') || '';
    const dot = prompt('DoT address (optional):') || '';
    const doq = prompt('DoQ address (optional):') || '';

    const server: Server = { name, dns, doh, dot, doq };

    // Will be replaced with: await window.go.main.App.AddServer(server);
    this.servers.push(server);
    this.renderServerList();
  }

  async importServers() {
    // Will be replaced with Wails file dialog
    alert('Import functionality will be implemented with Wails file dialog');
  }

  showElement(id: string) {
    document.getElementById(id)?.classList.remove('hidden');
  }

  hideElement(id: string) {
    document.getElementById(id)?.classList.add('hidden');
  }

  // Mock functions - will be replaced with actual Wails bindings
  async mockGetServers(): Promise<Server[]> {
    return [
      { name: 'Cloudflare Primary', dns: '1.1.1.1', doh: 'https://cloudflare-dns.com/dns-query', dot: '1.1.1.1:853', doq: '1.1.1.1:8853' },
      { name: 'Google Primary', dns: '8.8.8.8', doh: 'https://dns.google/dns-query', dot: '8.8.8.8:853', doq: '' },
    ];
  }

  async mockRunBenchmark(protocols: string[]): Promise<void> {
    // Simulate progress
    const progressBar = document.getElementById('progress-bar');
    const progressPercentage = document.getElementById('progress-percentage');
    const currentServer = document.getElementById('current-server');

    for (let i = 0; i <= 100; i += 10) {
      await new Promise(resolve => setTimeout(resolve, 200));
      if (progressBar) progressBar.style.width = `${i}%`;
      if (progressPercentage) progressPercentage.textContent = `${i}%`;
      if (currentServer) currentServer.textContent = `Testing server ${Math.floor(i / 10)}...`;
    }
  }

  async mockGetResults(): Promise<BenchmarkResults> {
    return {
      startTime: new Date().toISOString(),
      endTime: new Date().toISOString(),
      duration: '5s',
      results: [
        {
          serverName: 'Cloudflare Primary',
          protocol: 'DNS',
          metrics: {
            min: '15000000',
            max: '25000000',
            mean: '20000000',
            median: '19000000',
            p95: '24000000',
            p99: '25000000',
            success: 10,
            failed: 0,
            total: 10,
            successRate: '100.00%'
          }
        },
        {
          serverName: 'Cloudflare Primary',
          protocol: 'DoH',
          metrics: {
            min: '25000000',
            max: '45000000',
            mean: '35000000',
            median: '34000000',
            p95: '42000000',
            p99: '44000000',
            success: 10,
            failed: 0,
            total: 10,
            successRate: '100.00%'
          }
        }
      ]
    };
  }

  async mockRefreshServers(): Promise<void> {
    await new Promise(resolve => setTimeout(resolve, 500));
  }
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', () => {
    const app = new App();
    app.init();
  });
} else {
  const app = new App();
  app.init();
}
