# Dashboard Layout

Dashboard layout with sidebar navigation and stats cards.

```vue
<template>
  <v-app>
    <!-- Top Navigation -->
    <v-app-bar color="primary" elevation="2">
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-app-bar-title>Dashboard</v-app-bar-title>
      <v-spacer></v-spacer>
      <v-btn icon="mdi-bell"></v-btn>
      <v-btn icon="mdi-account-circle"></v-btn>
    </v-app-bar>
    
    <!-- Sidebar -->
    <v-navigation-drawer v-model="drawer" permanent>
      <v-list density="compact" nav>
        <v-list-item
          v-for="item in navItems"
          :key="item.title"
          :prepend-icon="item.icon"
          :title="item.title"
          :value="item.value"
          :active="currentNav === item.value"
          @click="currentNav = item.value"
        ></v-list-item>
      </v-list>
    </v-navigation-drawer>
    
    <!-- Main Content -->
    <v-main>
      <v-container fluid>
        <!-- Stats Cards -->
        <v-row class="mb-4">
          <v-col
            v-for="stat in stats"
            :key="stat.title"
            cols="12"
            sm="6"
            md="3"
          >
            <v-card :color="stat.color" dark>
              <v-card-text>
                <div class="d-flex justify-space-between align-center">
                  <div>
                    <div class="text-h5">{{ stat.value }}</div>
                    <div class="text-body-2">{{ stat.title }}</div>
                  </div>
                  <v-icon size="48" color="white" style="opacity: 0.5">
                    {{ stat.icon }}
                  </v-icon>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
        
        <!-- Data Table -->
        <v-row>
          <v-col cols="12">
            <v-card>
              <v-card-title>Recent Orders</v-card-title>
              <v-data-table
                :headers="headers"
                :items="orders"
                :items-length="orders.length"
              >
                <template v-slot:item.status="{ item }">
                  <v-chip
                    :color="getStatusColor(item.status)"
                    size="small"
                  >
                    {{ item.status }}
                  </v-chip>
                </template>
              </v-data-table>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'

const drawer = ref(true)
const currentNav = ref('dashboard')

const navItems = [
  { title: 'Dashboard', icon: 'mdi-view-dashboard', value: 'dashboard' },
  { title: 'Orders', icon: 'mdi-cart', value: 'orders' },
  { title: 'Products', icon: 'mdi-package', value: 'products' },
  { title: 'Users', icon: 'mdi-account-group', value: 'users' },
  { title: 'Settings', icon: 'mdi-cog', value: 'settings' }
]

const stats = [
  { title: 'Total Orders', value: '1,234', icon: 'mdi-cart', color: 'primary' },
  { title: 'Total Revenue', value: '$56,789', icon: 'mdi-cash', color: 'success' },
  { title: 'Users', value: '5,678', icon: 'mdi-account-group', color: 'info' },
  { title: 'Pending', value: '12', icon: 'mdi-clock-alert', color: 'warning' }
]

const headers = [
  { title: 'Order ID', key: 'id' },
  { title: 'Customer', key: 'customer' },
  { title: 'Amount', key: 'amount' },
  { title: 'Status', key: 'status' },
  { title: 'Date', key: 'date' }
]

const orders = ref([
  { id: 'ORD-001', customer: 'John Doe', amount: 299, status: 'completed', date: '2024-01-15' },
  { id: 'ORD-002', customer: 'Jane Smith', amount: 599, status: 'processing', date: '2024-01-16' },
  { id: 'ORD-003', customer: 'Bob Wilson', amount: 199, status: 'pending', date: '2024-01-17' }
])

const getStatusColor = (status) => {
  const colors = {
    completed: 'success',
    processing: 'info',
    pending: 'warning',
    cancelled: 'error'
  }
  return colors[status] || 'grey'
}
</script>
```
