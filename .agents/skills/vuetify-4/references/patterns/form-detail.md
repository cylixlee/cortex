# Form Detail Page

Form detail page for editing entities.

```vue
<template>
  <v-app>
    <v-app-bar>
      <v-btn icon="mdi-arrow-left" variant="text"></v-btn>
      <v-app-bar-title>Edit User</v-app-bar-title>
    </v-app-bar>
    
    <v-main>
      <v-container>
        <v-row justify="center">
          <v-col cols="12" md="8">
            <v-card>
              <v-card-title>Basic Information</v-card-title>
              <v-card-text>
                <v-form ref="formRef">
                  <v-row>
                    <v-col cols="12" sm="6">
                      <v-text-field
                        v-model="form.firstName"
                        label="First Name"
                        variant="outlined"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="6">
                      <v-text-field
                        v-model="form.lastName"
                        label="Last Name"
                        variant="outlined"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  
                  <v-text-field
                    v-model="form.email"
                    label="Email"
                    type="email"
                    variant="outlined"
                  ></v-text-field>
                  
                  <v-text-field
                    v-model="form.phone"
                    label="Phone"
                    variant="outlined"
                  ></v-text-field>
                  
                  <v-select
                    v-model="form.role"
                    label="Role"
                    :items="['Admin', 'Editor', 'User']"
                    variant="outlined"
                  ></v-select>
                  
                  <v-switch
                    v-model="form.active"
                    label="Enable Account"
                    color="success"
                  ></v-switch>
                </v-form>
              </v-card-text>
              
              <v-divider></v-divider>
              
              <v-card-actions class="pa-4">
                <v-btn variant="text">Cancel</v-btn>
                <v-spacer></v-spacer>
                <v-btn color="primary" variant="flat">
                  Save
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

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  role: null,
  active: true
})
</script>
```
