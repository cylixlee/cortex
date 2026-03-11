<script setup lang="ts">
defineProps<{
  status: string
}>()

const getStatusClass = (status: string) => {
  switch (status) {
    case 'pending':
      return 'status-badge-pending'
    case 'processing':
      return 'status-badge-processing'
    case 'completed':
      return 'status-badge-completed'
    case 'failed':
      return 'status-badge-failed'
    default:
      return 'status-badge-default'
  }
}

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
</script>

<template>
  <span :class="['status-badge', getStatusClass(status)]">
    {{ getStatusText(status) }}
  </span>
</template>

<style scoped>
.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge-pending {
  background-color: #fef3c7;
  color: #92400e;
}

.status-badge-processing {
  background-color: #dbeafe;
  color: #1e40af;
}

.status-badge-completed {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge-failed {
  background-color: #fee2e2;
  color: #991b1b;
}

.status-badge-default {
  background-color: #f3f4f6;
  color: #374151;
}
</style>
