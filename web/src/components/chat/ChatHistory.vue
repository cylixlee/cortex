<script setup lang="ts">
import type { Message } from '@/api/conversation'

defineProps<{
  messages: Message[]
  isLoading: boolean
}>()

defineExpose({
  scrollToBottom: () => {
    const container = document.querySelector('.messages-scroll') as HTMLElement
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  },
})
</script>

<template>
  <v-main class="chat-history-main">
    <v-sheet class="messages-container fill-height" color="chat">
      <div class="messages-scroll">
        <div class="messages-wrapper w-100 h-100">
          <div
            v-if="messages.length === 0 && !isLoading"
            class="empty-state d-flex flex-column align-center justify-center fill-height"
          >
            <v-icon icon="mdi-robot" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
            <h2 class="text-h5 text-grey-darken-1">How can I help you today?</h2>
          </div>

          <div v-else class="pa-4">
            <div v-for="msg in messages" :key="msg.id" class="message-wrapper mb-3" :class="msg.role">
              <v-avatar
                :color="msg.role === 'user' ? 'primary' : 'grey-lighten-3'"
                :variant="msg.role === 'assistant' ? 'flat' : 'flat'"
                size="36"
                class="message-avatar flex-shrink-0"
              >
                <v-icon v-if="msg.role === 'user'" icon="mdi-account"></v-icon>
                <v-icon v-else icon="mdi-robot"></v-icon>
              </v-avatar>

              <v-card
                :color="msg.role === 'user' ? 'primary' : 'grey-lighten-4'"
                :class="msg.role === 'user' ? 'text-white' : ''"
                flat
                class="message-card pa-3"
                max-width="70%"
              >
                <div class="message-text">{{ msg.content }}</div>
              </v-card>
            </div>

            <div v-if="isLoading" class="message-wrapper assistant mb-3">
              <v-avatar color="grey-lighten-3" variant="flat" size="36" class="message-avatar flex-shrink-0">
                <v-icon icon="mdi-robot"></v-icon>
              </v-avatar>
              <v-card color="surface" flat class="message-card pa-3 d-flex align-center" min-width="60">
                <v-progress-circular indeterminate size="20" width="2"></v-progress-circular>
              </v-card>
            </div>
          </div>
        </div>
      </div>
    </v-sheet>
  </v-main>
</template>

<style scoped>
.chat-history-main {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.messages-container {
  flex: 1;
  overflow: hidden;
}

.messages-scroll {
  height: 100%;
  overflow-y: auto;
}

.messages-wrapper {
  max-width: 900px;
  margin: 0 auto;
}

.message-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.message-wrapper.user {
  flex-direction: row-reverse;
}
</style>
