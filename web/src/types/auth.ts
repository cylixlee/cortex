export interface User {
  id: string
  email: string
  role: string
  created_at: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
}

export interface ApiError {
  error: string
}
