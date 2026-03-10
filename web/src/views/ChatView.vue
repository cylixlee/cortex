<script setup lang="ts">
import { ref, nextTick, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { sendMessage, getConversation } from '@/api/conversation'

const route = useRoute()
const router = useRouter()
const store = useChatStore()

const input = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

const conversationId = computed(() => route.params.id as string | undefined)

const emit = defineEmits<{
  (e: 'toggle-sidebar'): void
  (e: 'refresh-conversations'): void
}>()

onMounted(async () => {
  if (conversationId.value) {
    await loadConversation(conversationId.value)
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
    const newConversationId = await sendMessage(content, conversationId.value || undefined, (chunk) => {
      store.addAssistantMessage(chunk)
      nextTick(() => scrollToBottom())
    })

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
    <div class="chat-main">
      <div class="chat-header">
        <button class="sidebar-toggle" @click="$emit('toggle-sidebar')">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="12" x2="21" y2="12"></line>
            <line x1="3" y1="6" x2="21" y2="6"></line>
            <line x1="3" y1="18" x2="21" y2="18"></line>
          </svg>
        </button>
        <h1 class="chat-title">{{ store.conversationTitle }}</h1>
        <button v-if="!conversationId" class="new-chat-btn" @click="newChat">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
      </div>

      <div class="messages" ref="messagesContainer">
        <div v-if="store.messages.length === 0 && !store.isLoading" class="empty-state">
          <div class="empty-icon">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
          </div>
          <h2>How can I help you today?</h2>
        </div>

        <div v-for="msg in store.messages" :key="msg.id" :class="['message', msg.role]">
          <div class="message-avatar">
            <svg
              v-if="msg.role === 'user'"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z"></path>
              <path d="M12 6v6l4 2"></path>
            </svg>
          </div>
          <div class="message-content">{{ msg.content }}</div>
        </div>

        <div v-if="store.isLoading" class="message assistant loading">
          <div class="message-avatar">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 18a8 8 0 1 1 8-8 8 8 0 0 1-8 8z"></path>
              <path d="M12 6v6l4 2"></path>
            </svg>
          </div>
          <div class="message-content">
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
        </div>
      </div>

      <div class="input-area">
        <div class="input-container">
          <textarea
            v-model="input"
            @keydown="handleKeydown"
            placeholder="Message Cortex..."
            :disabled="store.isLoading"
            rows="1"
          ></textarea>
          <button @click="handleSend" :disabled="store.isLoading || !input.trim()" class="send-btn">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="22" y1="2" x2="11" y2="13"></line>
              <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
            </svg>
          </button>
        </div>
        <p class="disclaimer">AI can make mistakes. Please verify important information.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-layout {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  max-width: 768px;
  width: 100%;
  margin: 0 auto;
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  border-bottom: 1px solid #e5e5e5;
  background: #fafafa;
  flex-shrink: 0;
}

.sidebar-toggle {
  display: none;
  padding: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: #374151;
  border-radius: 6px;
}

.sidebar-toggle:hover {
  background: #f3f4f6;
}

.chat-title {
  flex: 1;
  font-size: 18px;
  font-weight: 500;
  color: #374151;
  margin: 0;
}

.new-chat-btn {
  padding: 8px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: #6b7280;
  border-radius: 6px;
  transition: all 0.15s;
}

.new-chat-btn:hover {
  background: #f3f4f6;
  color: #10a37f;
}

.messages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 24px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  color: #6b7280;
}

.empty-icon {
  margin-bottom: 16px;
  color: #9ca3af;
}

.empty-state h2 {
  font-size: 20px;
  font-weight: 500;
  margin: 0;
}

.message {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user .message-avatar {
  background: #10a37f;
  color: white;
}

.assistant .message-avatar {
  background: #fff;
  border: 1px solid #e5e5e5;
  color: #10a37f;
}

.message-content {
  padding: 14px 16px;
  border-radius: 12px;
  font-size: 16px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  max-width: calc(100% - 60px);
}

.user .message-content {
  background: #10a37f;
  color: white;
  border-top-right-radius: 4px;
}

.assistant .message-content {
  background: #f7f7f8;
  color: #374151;
  border-top-left-radius: 4px;
}

.loading .message-content {
  display: flex;
  gap: 4px;
  align-items: center;
  padding: 14px 16px;
}

.dot {
  width: 6px;
  height: 6px;
  background: #9ca3af;
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out both;
}

.dot:nth-child(1) {
  animation-delay: -0.32s;
}
.dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%,
  80%,
  100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.input-area {
  padding: 16px 24px 24px;
  flex-shrink: 0;
  background: #fff;
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 12px 16px;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition:
    border-color 0.2s,
    box-shadow 0.2s;
}

.input-container:focus-within {
  border-color: #10a37f;
  box-shadow: 0 2px 12px rgba(16, 163, 127, 0.1);
}

textarea {
  flex: 1;
  border: none;
  resize: none;
  font-family: inherit;
  font-size: 16px;
  line-height: 1.5;
  outline: none;
  background: transparent;
  max-height: 200px;
  min-height: 24px;
}

textarea::placeholder {
  color: #9ca3af;
}

.send-btn {
  padding: 10px;
  background: #10a37f;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  color: white;
  transition: background 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-btn:hover:not(:disabled) {
  background: #0d8c6d;
}

.send-btn:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.disclaimer {
  text-align: center;
  font-size: 12px;
  color: #9ca3af;
  margin: 12px 0 0;
}

@media (max-width: 768px) {
  .sidebar-toggle {
    display: flex;
  }
}
</style>
