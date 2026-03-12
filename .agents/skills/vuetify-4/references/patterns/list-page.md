# List Page

List page with search and actions.

```vue
<template>
  <v-app>
    <v-app-bar elevation="1">
      <v-app-bar-title>User Management</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        label="Search"
        prepend-inner-icon="mdi-magnify"
        hide-details
        single-line
        density="compact"
        class="mr-4"
        style="max-width: 300px"
      ></v-text-field>
      <v-btn color="primary" prepend-icon="mdi-plus">
        Add User
      </v-btn>
    </v-app-bar>
    
    <v-main>
      <v-container fluid>
        <v-card>
          <v-data-table
            v-model="selected"
            :headers="headers"
            :items="items"
            :search="search"
            :loading="loading"
            show-select
            hover
          >
            <template v-slot:item.name="{ item }">
              <div class="d-flex align-center py-2">
                <v-avatar color="primary" size="32" class="mr-2">
                  <span class="text-white">{{ item.name[0] }}</span>
                </v-avatar>
                {{ item.name }}
              </div>
            </template>
            
            <template v-slot:item.status="{ item }">
              <v-chip
                :color="item.status === 'active' ? 'success' : 'grey'"
                size="small"
              >
                {{ item.status === 'active' ? 'Active' : 'Inactive' }}
              </v-chip>
            </template>
            
            <template v-slot:item.createdAt="{ item }">
              {{ formatDate(item.createdAt) }}
            </template>
            
            <template v-slot:item.actions="{ item }">
              <v-btn
                icon="mdi-pencil"
                variant="text"
                size="small"
              ></v-btn>
              <v-btn
                icon="mdi-delete"
                variant="text"
                size="small"
                color="error"
              ></v-btn>
            </template>
          </v-data-table>
        </v-card>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'

const search = ref('')
const selected = ref([])
const loading = ref(false)

const headers = [
  { title: 'Username', key: 'name' },
  { title: 'Email', key: 'email' },
  { title: 'Role', key: 'role' },
  { title: 'Status', key: 'status' },
  { title: 'Created At', key: 'createdAt' },
  { title: 'Actions', key: 'actions', sortable: false, align: 'end' }
]

const items = ref([
  { name: 'John Doe', email: 'john@example.com', role: 'Admin', status: 'active', createdAt: '2024-01-01' },
  { name: 'Jane Smith', email: 'jane@example.com', role: 'User', status: 'active', createdAt: '2024-01-02' },
  { name: 'Bob Wilson', email: 'bob@example.com', role: 'User', status: 'inactive', createdAt: '2024-01-03' }
])

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}
</script>
```
