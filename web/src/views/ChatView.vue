<script setup lang="ts">
import { ref, nextTick, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { sendMessage, getConversation } from '@/api/conversation'

const route = useRoute()
const router = useRouter()
const store = useChatStore()

const input = ref('')
const enableRag = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

const conversationId = computed(() => route.params.id as string | undefined)

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'refresh-conversations'): void
}>()

onMounted(async () => {
  if (conversationId.value) {
    await loadConversation(conversationId.value)
  } else {
    store.clear()
  }
})

watch(conversationId, async (newId) => {
  if (newId) {
    await loadConversation(newId)
  } else {
    store.clear()
  }
})

async function loadConversation(id: string) {
  try {
    const conv = await getConversation(id)
    store.setConversation(conv.id, conv.title)
    store.setMessages(conv.messages)
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error('Failed to load conversation:', e)
    router.push('/chat')
  }
}

async function handleSend() {
  const content = input.value.trim()
  if (!content || store.isLoading) return

  input.value = ''
  store.addUserMessage(content)
  store.setLoading(true)

  await nextTick()
  scrollToBottom()

  try {
    const newConversationId = await sendMessage(
      content,
      conversationId.value || undefined,
      (chunk) => {
        store.addAssistantMessage(chunk)
        nextTick(() => scrollToBottom())
      },
      enableRag.value,
    )

    if (newConversationId) {
      if (!conversationId.value) {
        router.push(`/chat/${newConversationId}`)
      }
      emit('refresh-conversations')
    }
  } catch (error) {
    store.addAssistantMessage(`Error: ${(error as Error).message}`)
  } finally {
    store.setLoading(false)
    store.clearStreamingMessage()
  }
}

function scrollToBottom() {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

function newChat() {
  router.push('/chat')
}
</script>

<template>
  <div class="chat-layout">
    <v-app-bar flat density="comfortable" border="b">
      <v-app-bar-nav-icon class="d-md-none" @click="$emit('toggle-sidebar')"></v-app-bar-nav-icon>

      <v-app-bar-title>
        {{ store.conversationTitle || 'New Chat' }}
      </v-app-bar-title>

      <v-btn icon="mdi-plus" variant="text" @click="newChat"></v-btn>
    </v-app-bar>

    <v-main class="chat-main">
      <v-sheet class="messages-container" color="surface">
        <div ref="messagesContainer" class="messages-scroll">
          <div class="messages-wrapper">
            <div
              v-if="store.messages.length === 0 && !store.isLoading"
              class="empty-state d-flex flex-column align-center justify-center fill-height"
            >
              <v-icon icon="mdi-robot" size="64" color="grey-lighten-1" class="mb-4"></v-icon>
              <h2 class="text-h5 text-grey-darken-1">How can I help you today?</h2>
            </div>

            <div v-else>
              <div v-for="msg in store.messages" :key="msg.id" class="message-wrapper mb-4" :class="msg.role">
                <v-avatar
                  :color="msg.role === 'user' ? 'primary' : 'surface'"
                  :variant="msg.role === 'assistant' ? 'tonal' : 'flat'"
                  size="36"
                  class="message-avatar"
                >
                  <v-icon v-if="msg.role === 'user'" icon="mdi-account"></v-icon>
                  <v-icon v-else icon="mdi-robot"></v-icon>
                </v-avatar>

                <v-card
                  :color="msg.role === 'user' ? 'primary' : 'surface'"
                  :class="msg.role === 'user' ? 'text-white' : ''"
                  flat
                  class="message-card pa-3"
                >
                  <div class="message-text">{{ msg.content }}</div>
                </v-card>
              </div>

              <div v-if="store.isLoading" class="message-wrapper assistant mb-4">
                <v-avatar color="surface" variant="tonal" size="36" class="message-avatar">
                  <v-icon icon="mdi-robot"></v-icon>
                </v-avatar>
                <v-card color="surface" flat class="message-card pa-3 d-flex align-center">
                  <v-progress-circular indeterminate size="20" width="2"></v-progress-circular>
                </v-card>
              </div>
            </div>
          </div>
        </div>
      </v-sheet>
    </v-main>

    <div class="chat-footer pa-4">
      <div class="input-wrapper">
        <div class="d-flex align-center mb-3">
          <v-switch v-model="enableRag" label="Skill 检索" color="primary" density="compact" hide-details></v-switch>
        </div>

        <v-textarea
          v-model="input"
          @keydown="handleKeydown"
          placeholder="Message Cortex..."
          :disabled="store.isLoading"
          rows="1"
          auto-grow
          max-rows="6"
          hide-details
          class="message-input"
        >
          <template #append-inner>
            <v-btn
              icon="mdi-send"
              color="primary"
              variant="flat"
              :disabled="store.isLoading || !input.trim()"
              @click="handleSend"
            ></v-btn>
          </template>
        </v-textarea>

        <div class="text-center mt-2">
          <span class="text-caption text-grey"> AI can make mistakes. Please verify important information. </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-main {
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
  max-width: 768px;
  margin: 0 auto;
  padding: 24px;
}

.message-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.message-wrapper.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-card {
  max-width: 70%;
}

.message-input :deep(.v-field__input) {
  padding-top: 12px;
  padding-bottom: 12px;
}

.chat-footer {
  background: rgb(var(--v-theme-surface));
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.input-wrapper {
  width: 100%;
  max-width: 768px;
  margin: 0 auto;
}
</style>
