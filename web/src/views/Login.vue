<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

defineOptions({
  name: 'AuthPage',
})

const router = useRouter()
const userStore = useUserStore()

const email = ref('')
const password = ref('')
const error = ref('')
const isRegister = ref(false)
const showPassword = ref(false)

const rules = {
  required: (v: string) => !!v || 'Required',
  email: (v: string) => /.+@.+\..+/.test(v) || 'Invalid email',
  minLength: (v: string) => v.length >= 6 || 'Minimum 6 characters',
}

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
  <v-container class="fill-height" fluid>
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="pa-6" elevation="2">
          <v-card-title class="text-center text-h4 font-weight-bold pt-4">
            <span class="text-primary">Cortex</span>
          </v-card-title>

          <v-card-subtitle class="text-center">
            {{ isRegister ? 'Create your account' : 'Welcome back' }}
          </v-card-subtitle>

          <v-card-text>
            <v-form @submit.prevent="handleSubmit">
              <v-text-field
                v-model="email"
                label="Email"
                type="email"
                prepend-inner-icon="mdi-email-outline"
                :rules="[rules.required, rules.email]"
                autocomplete="email"
                class="mb-2"
              />

              <v-text-field
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                label="Password"
                prepend-inner-icon="mdi-lock-outline"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                :rules="[rules.required, rules.minLength]"
                autocomplete="current-password"
                class="mb-2"
              />

              <v-alert v-if="error" type="error" variant="tonal" class="mb-4" closable @click:close="error = ''">
                {{ error }}
              </v-alert>

              <v-btn type="submit" color="primary" size="large" block :loading="userStore.isLoading">
                {{ isRegister ? 'Create account' : 'Sign in' }}
              </v-btn>
            </v-form>
          </v-card-text>

          <v-card-actions class="justify-center pb-4">
            <v-btn variant="text" color="secondary" @click="isRegister = !isRegister">
              {{ isRegister ? 'Already have an account? Sign in' : "Don't have an account? Sign up" }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped></style>
