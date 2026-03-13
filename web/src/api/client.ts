const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080/api/v1"

export function getToken(): string | null {
  return localStorage.getItem("access_token")
}

export function setToken(token: string) {
  localStorage.setItem("access_token", token)
}

export function clearToken() {
  localStorage.removeItem("access_token")
}

export async function fetchWithAuth(
  url: string,
  options: RequestInit = {}
): Promise<Response> {
  const token = getToken()
  if (!token) throw new Error("Not authenticated")

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
  })
  return response
}

export async function authFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const response = await fetchWithAuth(`${API_BASE}${endpoint}`, options)
  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || `HTTP error! status: ${response.status}`)
  }
  return response.json()
}

export { API_BASE }
