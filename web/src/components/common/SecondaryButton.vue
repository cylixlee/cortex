<script setup lang="ts">
withDefaults(
  defineProps<{
    to?: string
    icon?: string
    loading?: boolean
    block?: boolean
    disabled?: boolean
    variant?: 'tonal' | 'outlined' | 'text'
  }>(),
  {
    loading: false,
    block: false,
    disabled: false,
    variant: 'tonal',
  },
)

defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()
</script>

<template>
  <v-btn
    v-if="to"
    :to="to"
    color="secondary"
    :variant="variant"
    :loading="loading"
    :disabled="disabled"
    :block="block"
    class="secondary-button"
  >
    <v-icon v-if="icon" start>{{ icon }}</v-icon>
    <slot />
  </v-btn>
  <v-btn
    v-else
    color="secondary"
    :variant="variant"
    :loading="loading"
    :disabled="disabled"
    :block="block"
    class="secondary-button"
    @click="$emit('click', $event)"
  >
    <v-icon v-if="icon" start>{{ icon }}</v-icon>
    <slot />
  </v-btn>
</template>

<style scoped>
.secondary-button {
  text-transform: none;
}
</style>
