<script setup lang="ts">
withDefaults(
  defineProps<{
    to?: string
    icon?: string
    loading?: boolean
    block?: boolean
    disabled?: boolean
  }>(),
  {
    loading: false,
    block: false,
    disabled: false,
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
    color="primary"
    :loading="loading"
    :disabled="disabled"
    :block="block"
    class="primary-button"
  >
    <v-icon v-if="icon" start>{{ icon }}</v-icon>
    <slot />
  </v-btn>
  <v-btn
    v-else
    color="primary"
    :loading="loading"
    :disabled="disabled"
    :block="block"
    class="primary-button"
    @click="$emit('click', $event)"
  >
    <v-icon v-if="icon" start>{{ icon }}</v-icon>
    <slot />
  </v-btn>
</template>

<style scoped>
.primary-button {
  text-transform: none;
}
</style>
