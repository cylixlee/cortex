<script setup lang="ts">
defineProps<{
  size?: 'small' | 'medium' | 'large'
}>()
</script>

<template>
  <div :class="['loading-spinner', `loading-spinner-${size || 'medium'}`]">
    <div class="spinner"></div>
  </div>
</template>

<style scoped>
.loading-spinner {
  display: flex;
  justify-content: center;
  align-items: center;
}

.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-spinner-small .spinner {
  width: 16px;
  height: 16px;
}

.loading-spinner-medium .spinner {
  width: 32px;
  height: 32px;
}

.loading-spinner-large .spinner {
  width: 48px;
  height: 48px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
