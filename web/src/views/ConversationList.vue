<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { listConversations, deleteConversation, type Conversation } from '@/api/conversation'
import PrimaryButton from '@/components/common/PrimaryButton.vue'
import SecondaryButton from '@/components/common/SecondaryButton.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const chatStore = useChatStore()

const conversations = ref<Conversation[]>([])
const isLoading = ref(true)

onMounted(async () => {
  await loadConversations()
})

async function loadConversations() {
  isLoading.value = true
  try {
    conversations.value = await listConversations()
  } catch (e) {
    console.error('Failed to load conversations:', e)
    if ((e as Error).message.includes('401') || (e as Error).message.includes('Not authenticated')) {
      localStorage.clear()
      window.location.href = '/login'
    }
  } finally {
    isLoading.value = false
  }
}

function openConversation(id: string) {
  router.push(`/chat/${id}`)
}

function newChat() {
  router.push('/chat')
}

function goToSkills() {
  router.push('/skills')
}

async function handleDelete(id: string, event: Event) {
  event.stopPropagation()
  if (!confirm('Delete this conversation?')) return
  try {
    await deleteConversation(id)
    conversations.value = conversations.value.filter((c) => c.id !== id)
    if (route.params.id === id) {
      chatStore.clear()
      router.push('/chat')
    }
  } catch (e) {
    console.error('Failed to delete:', e)
  }
}

function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function handleLogout() {
  userStore.logout()
  chatStore.clear()
  router.push('/login')
}

defineExpose({
  loadConversations,
})
</script>

<template>
  <div class="conversation-list-container d-flex flex-column fill-height bg-sidebar">
    <div class="pa-4">
      <div class="text-h6 text-grey-darken-3 font-weight-bold">Cortex</div>
    </div>

    <div class="pa-3">
      <SecondaryButton block icon="mdi-plus" class="mb-2" @click="newChat"> New Chat </SecondaryButton>

      <PrimaryButton block icon="mdi-brain" @click="goToSkills"> My Skills </PrimaryButton>
    </div>

    <v-list nav density="compact" class="pa-2 flex-grow-1 overflow-y-auto bg-sidebar">
      <v-list-subheader v-if="conversations.length">Recent</v-list-subheader>

      <v-list-item
        v-for="conv in conversations"
        :key="conv.id"
        :active="route.params.id === conv.id"
        @click="openConversation(conv.id)"
        rounded="lg"
        class="mb-1"
      >
        <template #prepend>
          <v-icon :icon="route.params.id === conv.id ? 'mdi-chat' : 'mdi-chat-outline'"></v-icon>
        </template>

        <v-list-item-title class="text-truncate">
          {{ conv.title || 'New Chat' }}
        </v-list-item-title>

        <v-list-item-subtitle>
          {{ formatDate(conv.updated_at) }}
        </v-list-item-subtitle>

        <template #append>
          <v-btn
            icon="mdi-delete-outline"
            variant="text"
            size="small"
            color="error"
            @click.stop="handleDelete(conv.id, $event)"
          ></v-btn>
        </template>
      </v-list-item>

      <v-list-item v-if="!isLoading && conversations.length === 0" class="text-center bg-transparent">
        <v-list-item-title class="text-muted">No conversations yet</v-list-item-title>
      </v-list-item>

      <v-list-item v-if="isLoading" class="text-center">
        <v-progress-circular indeterminate size="24" color="secondary"></v-progress-circular>
      </v-list-item>
    </v-list>

    <div class="pa-3 d-flex align-center">
      <v-avatar color="primary" size="36" class="mr-3">
        <span class="text-white">{{ userStore.user?.email?.[0]?.toUpperCase() }}</span>
      </v-avatar>

      <div class="flex-grow-1 text-truncate">
        <div class="text-body-2 text-truncate">{{ userStore.user?.email }}</div>
      </div>

      <v-btn icon="mdi-logout" variant="text" size="small" @click="handleLogout"></v-btn>
    </div>
  </div>
</template>

<style scoped>
.conversation-list-container {
  height: 100%;
}
</style>
