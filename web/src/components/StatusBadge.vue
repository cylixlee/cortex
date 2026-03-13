<script setup lang="ts">
import { statusColors } from '@/plugins/vuetify'

defineProps<{
  status: string
}>()

const getStatusText = (status: string) => {
  switch (status) {
    case 'pending':
      return 'Pending'
    case 'processing':
      return 'Processing'
    case 'completed':
      return 'Completed'
    case 'failed':
      return 'Failed'
    default:
      return status
  }
}

const getStatusColor = (status: string) => {
  return statusColors[status as keyof typeof statusColors] || statusColors.default
}
</script>

<template>
  <v-chip :color="getStatusColor(status).bg" :class="['status-badge']" label size="small">
    <span :style="{ color: `rgb(var(--v-theme-${getStatusColor(status).text}))` }">
      {{ getStatusText(status) }}
    </span>
  </v-chip>
</template>

<style scoped>
.status-badge {
  font-weight: 500;
  font-size: 12px;
}
</style>
