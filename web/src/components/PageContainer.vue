<script setup lang="ts">
import { layout } from '@/plugins/vuetify'

defineProps<{
  maxWidth?: string
  padding?: string
  fluid?: boolean
}>()
</script>

<template>
  <v-container
    :fluid="fluid"
    class="page-container"
    :style="{
      maxWidth: maxWidth || layout.pageMaxWidth,
      padding: padding || layout.pagePadding,
    }"
  >
    <slot />
  </v-container>
</template>

<style scoped>
.page-container {
  margin-left: auto;
  margin-right: auto;
  width: 100%;
}
</style>
