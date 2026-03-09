import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
}

export const useChatStore = defineStore('chat', () => {
  const messages = ref<Message[]>([])
  const isLoading = ref(false)

  function addUserMessage(content: string) {
    messages.value.push({
      id: Date.now().toString(),
      role: 'user',
      content,
    })
  }

  function addAssistantMessage(content: string) {
    const existingMsg = messages.value.find(m => m.role === 'assistant' && m.id === 'streaming')
    if (existingMsg) {
      existingMsg.content += content
    } else {
      messages.value.push({
        id: 'streaming',
        role: 'assistant',
        content,
      })
    }
  }

  function setLoading(loading: boolean) {
    isLoading.value = loading
  }

  function clearStreamingMessage() {
    const idx = messages.value.findIndex(m => m.id === 'streaming')
    if (idx !== -1) {
      const msg = messages.value[idx]
      if (msg) {
        msg.id = Date.now().toString()
      }
    }
  }

  function clear() {
    messages.value = []
  }

  return {
    messages,
    isLoading,
    addUserMessage,
    addAssistantMessage,
    setLoading,
    clearStreamingMessage,
    clear,
  }
})
