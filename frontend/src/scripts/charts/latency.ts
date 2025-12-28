import type { BenchmarkResults } from '../api/types';

let chart: any | null = null;

export async function renderLatencyChart(results: BenchmarkResults, canvasId = 'latency-chart'): Promise<void> {
  const canvas = document.getElementById(canvasId) as HTMLCanvasElement | null;
  if (!canvas) return;

  const { Chart, registerables } = await import('chart.js');
  Chart.register(...registerables);

  const labels = results.results.map((r) => `${r.serverName} (${r.protocol})`);
  const data = results.results.map((r) => r.metrics.mean);

  if (chart) {
    chart.destroy();
  }

  chart = new Chart(canvas, {
    type: 'bar',
    data: {
      labels,
      datasets: [
        {
          label: 'Mean Latency (ms)',
          data,
          backgroundColor: '#0066ff',
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
          grid: {
            color: '#3a3a3a',
          },
          ticks: {
            color: '#ffffff',
          },
        },
        x: {
          grid: {
            color: '#3a3a3a',
          },
          ticks: {
            color: '#ffffff',
            maxRotation: 45,
            minRotation: 45,
          },
        },
      },
      plugins: {
        legend: {
          labels: {
            color: '#ffffff',
          },
        },
      },
    },
  });
}

