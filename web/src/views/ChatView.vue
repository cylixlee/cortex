<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useChatStore } from '@/stores/chat'
import { sendMessage } from '@/api/chat'

const store = useChatStore()
const input = ref('')
const messagesContainer = ref<HTMLElement | null>(null)

async function handleSend() {
  const content = input.value.trim()
  if (!content || store.isLoading) return

  input.value = ''
  store.addUserMessage(content)
  store.setLoading(true)

  await nextTick()
  scrollToBottom()

  try {
    await sendMessage(content, (chunk) => {
      store.addAssistantMessage(chunk)
      nextTick(() => scrollToBottom())
    })
  } catch (error) {
    store.addAssistantMessage(`Error: ${error}`)
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
</script>

<template>
  <div class="chat-container">
    <div class="messages" ref="messagesContainer">
      <div
        v-for="msg in store.messages"
        :key="msg.id"
        :class="['message', msg.role]"
      >
        <div class="message-content">{{ msg.content }}</div>
      </div>
      <div v-if="store.isLoading" class="message assistant">
        <div class="message-content loading">Thinking...</div>
      </div>
    </div>
    <div class="input-area">
      <textarea
        v-model="input"
        @keydown="handleKeydown"
        placeholder="Type your message..."
        :disabled="store.isLoading"
      ></textarea>
      <button @click="handleSend" :disabled="store.isLoading || !input.trim()">
        Send
      </button>
    </div>
  </div>
</template>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px 0;
}

.message {
  margin-bottom: 16px;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  white-space: pre-wrap;
  word-break: break-word;
}

.user .message-content {
  background-color: #007AFF;
  color: white;
}

.assistant .message-content {
  background-color: #f0f0f0;
  color: #333;
}

.loading {
  color: #888;
  font-style: italic;
}

.input-area {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

textarea {
  flex: 1;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  resize: none;
  font-family: inherit;
  font-size: 14px;
}

textarea:focus {
  outline: none;
  border-color: #007AFF;
}

button {
  padding: 12px 24px;
  background-color: #007AFF;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
}

button:hover:not(:disabled) {
  background-color: #0056b3;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
</style>
