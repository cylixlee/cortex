import { getToken, setToken, clearToken, API_BASE } from "./client"
import type { User, LoginResponse } from "@/types"

export { getToken, setToken, clearToken }

export async function register(
  email: string,
  password: string
): Promise<{ user_id: string }> {
  const response = await fetch(`${API_BASE}/auth/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || "Registration failed")
  }
  return response.json()
}

export async function login(
  email: string,
  password: string
): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error || "Login failed")
  }
  const data = await response.json()
  setToken(data.access_token)
  localStorage.setItem("refresh_token", data.refresh_token)
  return data
}

export async function logout() {
  clearToken()
  localStorage.removeItem("refresh_token")
  localStorage.removeItem("user")
}

export async function getCurrentUser(): Promise<User> {
  const token = getToken()
  if (!token) throw new Error("Not authenticated")

  const response = await fetch(`${API_BASE}/auth/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!response.ok) {
    throw new Error("Failed to get user")
  }
  const user = await response.json()
  localStorage.setItem("user", JSON.stringify(user))
  return user
}

export async function refreshToken(): Promise<string> {
  const refreshToken = localStorage.getItem("refresh_token")
  if (!refreshToken) throw new Error("No refresh token")

  const response = await fetch(`${API_BASE}/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
  if (!response.ok) {
    logout()
    throw new Error("Token refresh failed")
  }
  const data = await response.json()
  setToken(data.access_token)
  return data.access_token
}

export function isAuthenticated(): boolean {
  return !!getToken()
}

export function getStoredUser(): User | null {
  const user = localStorage.getItem("user")
  return user ? JSON.parse(user) : null
}
