<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useChatStore } from '@/stores/chat'
import { listConversations, deleteConversation, type Conversation } from '@/api/conversation'

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
  <div class="sidebar">
    <div class="sidebar-header">
      <div class="logo">Cortex</div>
      <button class="new-chat-btn" @click="newChat">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"></line>
          <line x1="5" y1="12" x2="19" y2="12"></line>
        </svg>
        New Chat
      </button>
      <button class="skills-btn" @click="goToSkills">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
        </svg>
        My Skills
      </button>
    </div>

    <div class="sidebar-content">
      <div v-if="isLoading" class="loading">Loading...</div>
      <div v-else-if="conversations.length === 0" class="empty">No conversations yet</div>
      <div v-else class="conversation-list">
        <div v-for="conv in conversations" :key="conv.id" class="conversation-item" @click="openConversation(conv.id)">
          <div class="conv-info">
            <span class="conv-title">{{ conv.title || 'New Chat' }}</span>
            <span class="conv-date">{{ formatDate(conv.updated_at) }}</span>
          </div>
          <button class="delete-btn" @click="handleDelete(conv.id, $event)">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <div class="sidebar-footer">
      <div class="user-info">
        <div class="user-avatar">{{ userStore.user?.email?.[0]?.toUpperCase() }}</div>
        <span class="user-email">{{ userStore.user?.email }}</span>
      </div>
      <button class="logout-btn" @click="handleLogout">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.sidebar {
  width: 280px;
  height: 100vh;
  background: #f7f7f8;
  display: flex;
  flex-direction: column;
  color: #374151;
  border-right: 1px solid #e5e5e5;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e5e5e5;
}

.logo {
  font-size: 18px;
  font-weight: 600;
  color: #10a37f;
  margin-bottom: 16px;
  letter-spacing: -0.3px;
}

.new-chat-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.new-chat-btn:hover {
  background: #f3f4f6;
  border-color: #10a37f;
  color: #10a37f;
}

.skills-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  margin-top: 8px;
  background: #fff;
  border: 1px solid #e5e5e5;
  border-radius: 8px;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.skills-btn:hover {
  background: #f3f4f6;
  border-color: #6366f1;
  color: #6366f1;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.loading,
.empty {
  padding: 20px;
  text-align: center;
  color: #9ca3af;
  font-size: 14px;
}

.conversation-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.conversation-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.conversation-item:hover {
  background: #eee;
}

.conv-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
}

.conv-title {
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conv-date {
  font-size: 12px;
  color: #9ca3af;
}

.delete-btn {
  padding: 6px;
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 4px;
  opacity: 0;
  transition:
    opacity 0.15s,
    color 0.15s;
}

.conversation-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  color: #ef4444;
}

.sidebar-footer {
  padding: 12px 16px;
  border-top: 1px solid #e5e5e5;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  overflow: hidden;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #10a37f;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 500;
  color: white;
  flex-shrink: 0;
}

.user-email {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #374151;
}

.logout-btn {
  padding: 8px;
  background: transparent;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 6px;
  transition: color 0.15s;
}

.logout-btn:hover {
  color: #374151;
}
</style>
