<script setup lang="ts">
import { layout } from '@/plugins/vuetify'

defineProps<{
  padding?: string
  elevation?: number
  color?: string
}>()
</script>

<template>
  <v-card
    :color="color"
    :style="{ padding: padding || layout.cardPadding }"
    :elevation="elevation"
    class="content-card"
  >
    <slot />
  </v-card>
</template>

<style scoped>
.content-card {
  width: 100%;
}
</style>
