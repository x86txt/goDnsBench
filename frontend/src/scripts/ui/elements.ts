export function showElement(id: string): void {
  document.getElementById(id)?.classList.remove('hidden');
}

export function hideElement(id: string): void {
  document.getElementById(id)?.classList.add('hidden');
}

