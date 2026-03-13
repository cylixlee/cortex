<script setup lang="ts">
defineProps<{
  disabled?: boolean
  loading?: boolean
  canSend?: boolean
}>()

const emit = defineEmits<{
  (e: 'send'): void
}>()

function handleSend() {
  emit('send')
}
</script>

<template>
  <v-btn
    icon="mdi-arrow-up"
    color="primary"
    variant="flat"
    size="small"
    rounded="full"
    :disabled="disabled || loading || !canSend"
    @click="handleSend"
  ></v-btn>
</template>
