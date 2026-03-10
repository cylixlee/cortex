import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message } from '@/api/conversation'

export const useChatStore = defineStore('chat', () => {
  const messages = ref<Message[]>([])
  const isLoading = ref(false)
  const conversationId = ref<string | null>(null)
  const conversationTitle = ref('New Chat')

  function setConversation(id: string | null, title: string = 'New Chat') {
    conversationId.value = id
    conversationTitle.value = title
  }

  function setMessages(msgs: Message[]) {
    messages.value = msgs
  }

  function addUserMessage(content: string) {
    messages.value.push({
      id: Date.now().toString(),
      role: 'user',
      content,
      created_at: new Date().toISOString(),
    })
  }

  function addAssistantMessage(content: string) {
    const existingMsg = messages.value.find((m) => m.role === 'assistant' && m.id === 'streaming')
    if (existingMsg) {
      existingMsg.content += content
    } else {
      messages.value.push({
        id: 'streaming',
        role: 'assistant',
        content,
        created_at: new Date().toISOString(),
      })
    }
  }

  function setLoading(loading: boolean) {
    isLoading.value = loading
  }

  function clearStreamingMessage() {
    const idx = messages.value.findIndex((m) => m.id === 'streaming')
    if (idx !== -1) {
      const msg = messages.value[idx]
      if (msg) {
        msg.id = Date.now().toString()
      }
    }
  }

  function clear() {
    messages.value = []
    conversationId.value = null
    conversationTitle.value = 'New Chat'
  }

  return {
    messages,
    isLoading,
    conversationId,
    conversationTitle,
    setConversation,
    setMessages,
    addUserMessage,
    addAssistantMessage,
    setLoading,
    clearStreamingMessage,
    clear,
  }
})
