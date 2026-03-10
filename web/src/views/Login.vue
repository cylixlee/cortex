<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

defineOptions({
  name: 'AuthPage'
})

const router = useRouter()
const userStore = useUserStore()

const email = ref('')
const password = ref('')
const error = ref('')
const isRegister = ref(false)

async function handleSubmit() {
  error.value = ''
  try {
    if (isRegister.value) {
      await userStore.register(email.value, password.value)
      isRegister.value = false
      error.value = 'Registration successful! Please login.'
    } else {
      await userStore.login(email.value, password.value)
      router.push('/chat')
    }
  } catch (e) {
    error.value = (e as Error).message
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-header">
        <h1 class="logo">Cortex</h1>
        <p class="subtitle">{{ isRegister ? 'Create your account' : 'Welcome back' }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div class="form-group">
          <label for="email">Email</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="you@example.com"
            required
            autocomplete="email"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="••••••••"
            required
            minlength="6"
            autocomplete="current-password"
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <button type="submit" class="submit-btn" :disabled="userStore.isLoading">
          {{ userStore.isLoading ? 'Please wait...' : (isRegister ? 'Create account' : 'Sign in') }}
        </button>
      </form>

      <div class="auth-footer">
        <button class="toggle-btn" @click="isRegister = !isRegister">
          {{ isRegister ? 'Already have an account? Sign in' : "Don't have an account? Sign up" }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  width: 100%;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  padding: 20px;
}

.auth-container {
  width: 100%;
  max-width:400px;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  padding: 40px;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  font-size: 32px;
  font-weight: 600;
  color: #10a37f;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.subtitle {
  color: #6e6e80;
  font-size: 15px;
  margin: 0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  font-size: 15px;
  transition: border-color 0.2s, box-shadow 0.2s;
  background: #fff;
}

.form-group input:focus {
  outline: none;
  border-color: #10a37f;
  box-shadow: 0 0 0 3px rgba(16, 163, 127, 0.1);
}

.form-group input::placeholder {
  color: #9ca3af;
}

.error-message {
  padding: 12px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #dc2626;
  font-size: 14px;
  text-align: center;
}

.submit-btn {
  padding: 14px;
  background: #10a37f;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: #0d8c6d;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-footer {
  margin-top: 24px;
  text-align: center;
}

.toggle-btn {
  background: none;
  border: none;
  color: #10a37f;
  font-size: 14px;
  cursor: pointer;
  padding: 8px;
}

.toggle-btn:hover {
  text-decoration: underline;
}
</style>
