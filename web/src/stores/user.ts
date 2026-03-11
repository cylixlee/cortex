import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { User } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)

  const isLoggedIn = computed(() => !!user.value)

  async function initialize() {
    const storedUser = authApi.getStoredUser()
    if (storedUser && authApi.isAuthenticated()) {
      user.value = storedUser
    }
  }

  async function login(email: string, password: string) {
    isLoading.value = true
    try {
      await authApi.login(email, password)
      user.value = await authApi.getCurrentUser()
    } finally {
      isLoading.value = false
    }
  }

  async function register(email: string, password: string) {
    isLoading.value = true
    try {
      await authApi.register(email, password)
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    authApi.logout()
    user.value = null
  }

  return {
    user,
    isLoading,
    isLoggedIn,
    initialize,
    login,
    register,
    logout,
  }
})
