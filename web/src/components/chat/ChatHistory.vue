<script setup lang="ts">
import { computed, toRefs } from 'vue'
import type { Message } from '@/api/conversation'
import { layout, spacing } from '@/plugins/vuetify'

const props = withDefaults(
  defineProps<{
    messages: Message[]
    isLoading: boolean
    maxWidth?: string
  }>(),
  {
    maxWidth: layout.chatMaxWidth,
  },
)

const { messages } = toRefs(props)

const hasStreamingMessage = computed(() => messages.value.some((m: Message) => m.id === 'streaming'))

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
            <template v-for="msg in messages" :key="msg.id">
              <div
                v-if="msg.role === 'user'"
                class="message-wrapper user"
                :style="{ marginBottom: spacing.chatMessageGap }"
              >
                <v-card color="primary" class="message-card pa-3" max-width="70%" elevation="0">
                  <div class="message-text text-white">{{ msg.content }}</div>
                </v-card>
              </div>
              <div v-else class="message-wrapper assistant" :style="{ marginBottom: spacing.chatMessageGap }">
                <div class="message-text assistant-text">
                  {{ msg.content }}
                </div>
              </div>
            </template>

            <div
              v-if="isLoading && !hasStreamingMessage"
              class="message-wrapper assistant"
              :style="{ marginBottom: spacing.chatMessageGap }"
            >
              <div class="message-text assistant-text">
                <v-progress-circular indeterminate size="16" width="2" color="primary"></v-progress-circular>
              </div>
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
  max-width: v-bind(maxWidth);
  margin: 0 auto;
}

.message-wrapper {
  display: flex;
}

.message-wrapper.user {
  justify-content: flex-end;
}

.message-wrapper.assistant {
  justify-content: flex-start;
}

.message-card {
  border-radius: 12px;
}

.message-text {
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

.assistant-text {
  width: 100%;
  font-size: 14px;
  line-height: 1.6;
  color: rgb(var(--v-theme-on-surface));
}
</style>
