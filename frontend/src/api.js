import { logEvent } from './Logger';

function apiBase(raw) {
  const value = String(raw ?? '').trim().replace(/\/+$/, '');
  return value.replace(/\/api$/i, '');
}

const BASE_URL = apiBase(import.meta.env.VITE_API_URL);

export async function apiRequest(path, options = {}) {
  const response = await fetch(`${BASE_URL}${path}`, options);
  const data = await response.json().catch(() => null);

  return { response, data };
}

export async function apiFetch(path, options = {}) {
  const { response, data } = await apiRequest(path, options);

  if (!response.ok) {
    const errorMessage = `apiFetch error: ${path} - ${data?.error || response.statusText || 'Unknown error'}`;
    logEvent('ERROR', errorMessage, { status: response.status });
    throw new Error(errorMessage);
  }

  logEvent('INFO', `apiFetch OK: ${path}`, { status: response.status });
  return data;
}

export function jsonPost(body) {
  return {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body ?? {}),
  };
}

export { apiBase };
