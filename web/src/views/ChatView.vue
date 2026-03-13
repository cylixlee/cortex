<script setup lang="ts">
import { ref, nextTick, onMounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { sendMessage, getConversation } from '@/api/conversation'
import ChatTitle from '@/components/chat/ChatTitle.vue'
import ChatHistory from '@/components/chat/ChatHistory.vue'
import ChatInput from '@/components/chat/ChatInput.vue'

const route = useRoute()
const router = useRouter()
const store = useChatStore()

const enableRag = ref(false)
const historyRef = ref<InstanceType<typeof ChatHistory> | null>(null)

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

async function handleSend(content: string) {
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
  historyRef.value?.scrollToBottom()
}

function newChat() {
  router.push('/chat')
}
</script>

<template>
  <div class="chat-view-container d-flex flex-column fill-height">
    <ChatTitle :title="store.conversationTitle" @new-chat="newChat" @toggle-sidebar="$emit('toggle-sidebar')" />

    <ChatHistory ref="historyRef" :messages="store.messages" :is-loading="store.isLoading" />

    <ChatInput v-model="enableRag" :disabled="store.isLoading" :loading="store.isLoading" @send="handleSend" />
  </div>
</template>

<style scoped>
.chat-view-container {
  height: 100%;
}
</style>
