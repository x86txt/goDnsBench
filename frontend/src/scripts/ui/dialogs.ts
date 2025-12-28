// Dialog utilities for import/export functionality

interface Server {
  name: string;
  dns: string;
  doh: string;
  dot: string;
  doq: string;
}

/**
 * Shows a file picker dialog for importing servers
 * Returns the selected File object or null if cancelled
 */
export async function showImportDialog(): Promise<File | null> {
  return new Promise((resolve) => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json,.csv';
    input.style.display = 'none';

    input.addEventListener('change', (e) => {
      const target = e.target as HTMLInputElement;
      const file = target.files?.[0];
      resolve(file || null);
      document.body.removeChild(input);
    });

    input.addEventListener('cancel', () => {
      resolve(null);
      document.body.removeChild(input);
    });

    document.body.appendChild(input);
    input.click();
  });
}

/**
 * Reads a file as text
 */
export async function readFileAsText(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      resolve(e.target?.result as string);
    };
    reader.onerror = (e) => {
      reject(new Error('Failed to read file'));
    };
    reader.readAsText(file);
  });
}

/**
 * Shows a file save dialog for exporting results
 * Uses a prompt for filename since browser file inputs don't support save dialogs well
 * Returns the selected file path or null if cancelled
 */
export async function showExportDialog(defaultFilename: string, extension: 'json' | 'csv'): Promise<string | null> {
  const filename = prompt('Enter filename to save:', defaultFilename);
  if (!filename) {
    return null;
  }
  
  // Ensure correct extension
  let path = filename;
  if (!path.endsWith(`.${extension}`)) {
    path = `${path}.${extension}`;
  }
  
  return path;
}

/**
 * Shows a success message to the user
 */
export function showSuccess(message: string): void {
  // Simple alert for now, could be replaced with a toast notification
  alert(`Success: ${message}`);
}

/**
 * Shows an error message to the user
 */
export function showError(message: string): void {
  alert(`Error: ${message}`);
}
