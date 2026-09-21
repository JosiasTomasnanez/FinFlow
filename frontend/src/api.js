import { logEvent } from './Logger';

const BASE_URL = apiBase(import.meta.env.VITE_API_URL);

function apiBase(raw) {
  const value = String(raw ?? '').trim().replace(/\/+$/, '');
  return value.replace(/\/api$/i, '');
}

// Función auxiliar para generar un request_id rápido en el Frontend para correlación
function generateRequestId() {
  const bytes = new Uint8Array(16);
  window.crypto.getRandomValues(bytes);
  return Array.from(bytes, (byte) => byte.toString(16).padStart(2, '0')).join('');
}

export async function apiRequest(path, options = {}) {
  const requestId = generateRequestId();

  // Inyectamos un encabezado X-Request-ID para seguimiento en el backend
  const extendedOptions = {
    ...options,
    headers: {
      ...(options.headers ?? {}),
      'X-Request-ID': requestId,
    },
  };

  let response;
  let data;
  let errorOccurred = false;
  let errorMessage = '';

  try {
    response = await fetch(`${BASE_URL}${path}`, extendedOptions);
    data = await response.json().catch(() => null);

    if (!response.ok) {
      errorMessage = data?.error || response.statusText || 'Unknown error';
    }
  } catch (err) {
    errorOccurred = true;
    errorMessage = err.message || 'Network Error';
    throw err;
  } finally {

    const logAttributes = {
      request_id: requestId,
      method: options.method ?? 'GET',
      path: path,
      status: errorOccurred ? 0 : response?.status ?? 500,
      message: errorOccurred || !response?.ok ? `apiFetch KO: ${errorMessage}` : 'apiFetch OK'
    };

    if (errorOccurred || !response?.ok) {
      logEvent('ERROR', 'http_request', logAttributes);
    } else {
      logEvent('INFO', 'http_request', logAttributes);
    }
  }

  return { response, data, errorMessage };
}

export async function apiFetch(path, options = {}) {
  const { response, data, errorMessage } = await apiRequest(path, options);

  if (!response.ok) {
    throw new Error(errorMessage);
  }

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
