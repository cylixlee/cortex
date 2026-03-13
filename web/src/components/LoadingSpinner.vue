<script setup lang="ts">
withDefaults(
  defineProps<{
    size?: 'small' | 'medium' | 'large'
    color?: string
  }>(),
  {
    size: 'medium',
    color: 'primary',
  },
)

const sizeMap = {
  small: 16,
  medium: 32,
  large: 48,
}
</script>

<template>
  <div class="loading-spinner">
    <v-progress-circular :size="sizeMap[size]" :width="3" :color="color" indeterminate />
  </div>
</template>

<style scoped>
.loading-spinner {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
