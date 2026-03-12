# Login Page

```vue
<template>
  <v-app>
    <v-main>
      <v-container class="fill-height" fluid>
        <v-row justify="center" align="center">
          <v-col cols="12" sm="8" md="4">
            <v-card class="pa-4">
              <v-card-title class="text-center text-h5 mb-4">
                Login
              </v-card-title>
              
              <v-card-text>
                <v-form ref="formRef" v-model="valid">
                  <v-text-field
                    v-model="form.username"
                    label="Username"
                    prepend-inner-icon="mdi-account"
                    variant="outlined"
                    :rules="[rules.required]"
                  ></v-text-field>
                  
                  <v-text-field
                    v-model="form.password"
                    label="Password"
                    :type="showPassword ? 'text' : 'password'"
                    prepend-inner-icon="mdi-lock"
                    :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                    variant="outlined"
                    :rules="[rules.required]"
                    @click:append-inner="showPassword = !showPassword"
                  ></v-text-field>
                  
                  <v-checkbox
                    v-model="form.remember"
                    label="Remember me"
                    density="compact"
                  ></v-checkbox>
                </v-form>
              </v-card-text>
              
              <v-card-actions>
                <v-btn
                  block
                  color="primary"
                  size="large"
                  :disabled="!valid"
                  @click="handleLogin"
                >
                  Login
                </v-btn>
              </v-card-actions>
              
              <v-card-actions class="justify-center mt-2">
                <span class="text-body-2">Don't have an account?</span>
                <v-btn variant="text" color="primary" size="small">
                  Sign up
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { reactive, ref } from 'vue'

const formRef = ref(null)
const valid = ref(false)
const showPassword = ref(false)

const form = reactive({
  username: '',
  password: '',
  remember: false
})

const rules = {
  required: v => !!v || 'Required'
}

const handleLogin = () => {
  console.log('Login:', form)
}
</script>
```
