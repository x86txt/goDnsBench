import type { BenchmarkProgress } from '../api/types';

export function setProgress(progress: Partial<BenchmarkProgress>): void {
  const progressBar = document.getElementById('progress-bar') as HTMLElement | null;
  const progressPercentage = document.getElementById('progress-percentage');
  const currentServer = document.getElementById('current-server');

  if (typeof progress.percentage === 'number') {
    const pct = Math.max(0, Math.min(100, progress.percentage));
    if (progressBar) progressBar.style.width = `${pct}%`;
    if (progressPercentage) progressPercentage.textContent = `${Math.round(pct)}%`;
  }

  if (currentServer) {
    const protocolSuffix = progress.currentProtocol ? ` (${progress.currentProtocol})` : '';
    currentServer.textContent = progress.currentServer ? `${progress.currentServer}${protocolSuffix}` : 'Testing...';
  }
}

export function resetProgress(): void {
  setProgress({
    currentServer: 'Starting...',
    currentProtocol: '',
    completedTests: 0,
    totalTests: 0,
    percentage: 0,
  });
}

