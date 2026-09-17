function apiBase(raw) {
  const value = String(raw ?? '').trim().replace(/\/+$/, '')
  return value.replace(/\/api$/i, '')
}

const BASE_URL = apiBase(import.meta.env.VITE_API_URL)

export async function apiRequest(path, options = {}) {
  const response = await fetch(`${BASE_URL}${path}`, options)
  const data = await response.json().catch(() => null)
  return { ok: response.ok, status: response.status, data }
}

export async function apiFetch(path, options = {}) {
  const { ok, status, data } = await apiRequest(path, options)
  if (!ok) {
    throw new Error(data?.error || String(status))
  }
  return data
}

export function jsonPost(body) {
  return {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body ?? {}),
  }
}

export { apiBase }
