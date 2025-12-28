// Settings UI component and utilities

import type { Settings } from '../api/client';

// Re-export Settings type
export type { Settings };

/**
 * Shows the settings modal
 */
export function showSettingsModal(): void {
  const modal = document.getElementById('settings-modal');
  if (modal) {
    modal.classList.remove('hidden');
  }
}

/**
 * Hides the settings modal
 */
export function hideSettingsModal(): void {
  const modal = document.getElementById('settings-modal');
  if (modal) {
    modal.classList.add('hidden');
  }
}

/**
 * Populates the settings form with current settings
 */
export function populateSettingsForm(settings: Settings): void {
  // Timeout (already in milliseconds)
  const timeoutInput = document.getElementById('settings-timeout') as HTMLInputElement;
  if (timeoutInput) {
    timeoutInput.value = settings.queryTimeoutMs.toString();
  }

  // Concurrency
  const concurrencyInput = document.getElementById('settings-concurrency') as HTMLInputElement;
  if (concurrencyInput) {
    concurrencyInput.value = settings.maxConcurrent.toString();
  }

  // Protocols
  const protocols = ['DNS', 'DoH', 'DoT', 'DoQ'];
  protocols.forEach(protocol => {
    const checkbox = document.getElementById(`settings-protocol-${protocol.toLowerCase()}`) as HTMLInputElement;
    if (checkbox) {
      checkbox.checked = settings.enabledProtocols.includes(protocol);
    }
  });

  // Domains
  const domainInputs = {
    a: document.getElementById('settings-domains-a') as HTMLTextAreaElement,
    mx: document.getElementById('settings-domains-mx') as HTMLTextAreaElement,
    txt: document.getElementById('settings-domains-txt') as HTMLTextAreaElement,
    dnssec: document.getElementById('settings-domains-dnssec') as HTMLTextAreaElement,
  };

  if (domainInputs.a) {
    domainInputs.a.value = settings.selectedDomains.a.join('\n');
  }
  if (domainInputs.mx) {
    domainInputs.mx.value = settings.selectedDomains.mx.join('\n');
  }
  if (domainInputs.txt) {
    domainInputs.txt.value = settings.selectedDomains.txt.join('\n');
  }
  if (domainInputs.dnssec) {
    domainInputs.dnssec.value = settings.selectedDomains.dnssec.join('\n');
  }
}

/**
 * Collects settings from the form and returns a Settings object
 */
export function collectSettingsFromForm(): Settings {
  // Timeout (already in milliseconds)
  const timeoutInput = document.getElementById('settings-timeout') as HTMLInputElement;
  const timeoutMs = timeoutInput ? parseInt(timeoutInput.value) || 1000 : 1000;

  // Concurrency
  const concurrencyInput = document.getElementById('settings-concurrency') as HTMLInputElement;
  const concurrency = concurrencyInput ? parseInt(concurrencyInput.value) || 10 : 10;

  // Protocols
  const protocols: string[] = [];
  const protocolCheckboxes = ['dns', 'doh', 'dot', 'doq'];
  protocolCheckboxes.forEach(proto => {
    const checkbox = document.getElementById(`settings-protocol-${proto}`) as HTMLInputElement;
    if (checkbox && checkbox.checked) {
      protocols.push(proto.toUpperCase());
    }
  });

  // Domains
  const domainInputs = {
    a: document.getElementById('settings-domains-a') as HTMLTextAreaElement,
    mx: document.getElementById('settings-domains-mx') as HTMLTextAreaElement,
    txt: document.getElementById('settings-domains-txt') as HTMLTextAreaElement,
    dnssec: document.getElementById('settings-domains-dnssec') as HTMLTextAreaElement,
  };

  const parseDomains = (textarea: HTMLTextAreaElement | null): string[] => {
    if (!textarea) return [];
    return textarea.value
      .split('\n')
      .map(line => line.trim())
      .filter(line => line.length > 0);
  };

  return {
    queryTimeoutMs: timeoutMs,
    maxConcurrent: concurrency,
    serverListUrl: '', // Will be preserved from existing settings
    enabledProtocols: protocols,
    selectedDomains: {
      a: parseDomains(domainInputs.a),
      mx: parseDomains(domainInputs.mx),
      txt: parseDomains(domainInputs.txt),
      dnssec: parseDomains(domainInputs.dnssec),
    },
  };
}
