import type { Protocol, Server } from '../api/types';

export function renderServerList(servers: Server[]): void {
  const container = document.getElementById('server-list');
  if (!container) return;

  container.innerHTML = servers
    .map(
      (server) => `
      <div class="flex items-center justify-between px-3 py-2 rounded hover:bg-lighter-gray transition">
        <label class="flex items-center space-x-2 cursor-pointer flex-1">
          <input type="checkbox" class="form-checkbox text-electric-blue server-checkbox" data-server="${server.name}" checked>
          <span class="text-sm">${server.name}</span>
        </label>
      </div>
    `,
    )
    .join('');
}

export function getSelectedProtocols(): Protocol[] {
  const protocols: Protocol[] = [];
  if ((document.getElementById('protocol-dns') as HTMLInputElement | null)?.checked) protocols.push('DNS');
  if ((document.getElementById('protocol-doh') as HTMLInputElement | null)?.checked) protocols.push('DoH');
  if ((document.getElementById('protocol-dot') as HTMLInputElement | null)?.checked) protocols.push('DoT');
  if ((document.getElementById('protocol-doq') as HTMLInputElement | null)?.checked) protocols.push('DoQ');
  return protocols;
}

export function getSelectedServerNames(): string[] {
  const checkboxes = document.querySelectorAll('.server-checkbox:checked');
  return Array.from(checkboxes)
    .map((cb) => (cb as HTMLInputElement).dataset.server)
    .filter((name): name is string => !!name);
}

export function getSelectedServers(servers: Server[]): Server[] {
  const selectedNames = getSelectedServerNames();
  return servers.filter((s) => selectedNames.includes(s.name));
}

