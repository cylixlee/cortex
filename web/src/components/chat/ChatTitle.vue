<script setup lang="ts">
defineProps<{
  title: string
}>()

const emit = defineEmits<{
  (e: 'new-chat'): void
  (e: 'toggle-sidebar'): void
}>()
</script>

<template>
  <div class="chat-title d-flex align-center pa-2">
    <v-btn icon="mdi-menu" variant="text" class="d-md-none" @click="emit('toggle-sidebar')"></v-btn>

    <v-spacer class="d-md-none"></v-spacer>

    <v-toolbar-title class="text-center">
      {{ title || 'New Chat' }}
    </v-toolbar-title>

    <v-spacer></v-spacer>

    <v-btn icon="mdi-plus" variant="text" @click="emit('new-chat')"></v-btn>
  </div>
</template>

<style scoped>
.chat-title {
  min-height: 56px;
}
</style>
