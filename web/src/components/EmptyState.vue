<script setup lang="ts">
import { spacing } from '@/plugins/vuetify'

withDefaults(
  defineProps<{
    title?: string
    description?: string
    icon?: string
  }>(),
  {
    title: 'No data',
    icon: 'mdi-folder-open-outline',
  },
)
</script>

<template>
  <div class="empty-state" :style="{ padding: `${spacing.lg} ${spacing.md}` }">
    <v-icon :icon="icon" :size="64" color="grey-lighten-1" class="empty-icon mb-4" />
    <h3 class="empty-title">{{ title }}</h3>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <slot></slot>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  width: 100%;
}

.empty-title {
  font-size: 18px;
  font-weight: 600;
  color: rgb(var(--v-theme-grey-darken-3));
  margin: 0 0 8px;
}

.empty-description {
  font-size: 14px;
  color: rgb(var(--v-theme-muted));
  margin: 0;
}
</style>
