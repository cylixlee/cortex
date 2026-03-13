import { create } from "zustand"
import type { User } from "@/types"
import { isAuthenticated, getStoredUser, logout as apiLogout } from "@/api/auth"

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  checkAuth: () => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: getStoredUser(),
  isAuthenticated: isAuthenticated(),
  isLoading: false,
  checkAuth: () => {
    const authenticated = isAuthenticated()
    const user = getStoredUser()
    set({ isAuthenticated: authenticated, user })
  },
  logout: () => {
    apiLogout()
    set({ user: null, isAuthenticated: false })
  },
}))
