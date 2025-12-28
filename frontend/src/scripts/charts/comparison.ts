// Comparison charts for benchmark results

import type { BenchmarkResults } from '../api/types';

let charts: any[] = [];

export async function renderComparisonCharts(results: BenchmarkResults): Promise<void> {
  const { Chart, registerables } = await import('chart.js');
  Chart.register(...registerables);

  // Destroy existing charts
  charts.forEach((chart) => chart.destroy());
  charts = [];

  // Group results by server
  const byServer = new Map<string, any[]>();
  results.results.forEach((r) => {
    if (!byServer.has(r.serverName)) {
      byServer.set(r.serverName, []);
    }
    byServer.get(r.serverName)!.push(r);
  });

  // Group results by protocol
  const byProtocol = new Map<string, any[]>();
  results.results.forEach((r) => {
    if (!byProtocol.has(r.protocol)) {
      byProtocol.set(r.protocol, []);
    }
    byProtocol.get(r.protocol)!.push(r);
  });

  // Chart 1: Protocol comparison (mean latency by protocol)
  const protocolChartCanvas = document.getElementById('protocol-comparison-chart') as HTMLCanvasElement;
  if (protocolChartCanvas) {
    const protocolLabels = Array.from(byProtocol.keys());
    const protocolMeans = protocolLabels.map((proto) => {
      const protoResults = byProtocol.get(proto)!;
      const sum = protoResults.reduce((acc, r) => acc + r.metrics.mean, 0);
      return sum / protoResults.length;
    });

    charts.push(
      new Chart(protocolChartCanvas, {
        type: 'bar',
        data: {
          labels: protocolLabels,
          datasets: [
            {
              label: 'Mean Latency (ms)',
              data: protocolMeans,
              backgroundColor: ['#0066ff', '#00ffff', '#ff0066', '#00ff66'],
              borderColor: '#00ffff',
              borderWidth: 1,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          scales: {
            y: {
              beginAtZero: true,
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff' },
            },
            x: {
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff' },
            },
          },
          plugins: {
            legend: { labels: { color: '#ffffff' } },
          },
        },
      })
    );
  }

  // Chart 2: Server comparison (top 10 servers by mean latency)
  const serverChartCanvas = document.getElementById('server-comparison-chart') as HTMLCanvasElement;
  if (serverChartCanvas) {
    const serverData = Array.from(byServer.entries())
      .map(([name, results]) => ({
        name,
        mean: results.reduce((acc, r) => acc + r.metrics.mean, 0) / results.length,
      }))
      .sort((a, b) => a.mean - b.mean)
      .slice(0, 10);

    charts.push(
      new Chart(serverChartCanvas, {
        type: 'bar',
        data: {
          labels: serverData.map((s) => s.name),
          datasets: [
            {
              label: 'Mean Latency (ms)',
              data: serverData.map((s) => s.mean),
              backgroundColor: '#0066ff',
              borderColor: '#00ffff',
              borderWidth: 1,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          indexAxis: 'y',
          scales: {
            x: {
              beginAtZero: true,
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff' },
            },
            y: {
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff', maxRotation: 0 },
            },
          },
          plugins: {
            legend: { labels: { color: '#ffffff' } },
          },
        },
      })
    );
  }

  // Chart 3: Success rate comparison
  const successChartCanvas = document.getElementById('success-rate-chart') as HTMLCanvasElement;
  if (successChartCanvas) {
    const successData = results.results.map((r) => ({
      label: `${r.serverName} (${r.protocol})`,
      rate: (r.metrics.success / r.metrics.total) * 100,
    }));

    charts.push(
      new Chart(successChartCanvas, {
        type: 'line',
        data: {
          labels: successData.map((d) => d.label),
          datasets: [
            {
              label: 'Success Rate (%)',
              data: successData.map((d) => d.rate),
              borderColor: '#00ff66',
              backgroundColor: 'rgba(0, 255, 102, 0.1)',
              fill: true,
              tension: 0.4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: true,
          scales: {
            y: {
              beginAtZero: true,
              max: 100,
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff' },
            },
            x: {
              grid: { color: '#3a3a3a' },
              ticks: { color: '#ffffff', maxRotation: 45, minRotation: 45 },
            },
          },
          plugins: {
            legend: { labels: { color: '#ffffff' } },
          },
        },
      })
    );
  }
}

export function destroyComparisonCharts(): void {
  charts.forEach((chart) => chart.destroy());
  charts = [];
}
