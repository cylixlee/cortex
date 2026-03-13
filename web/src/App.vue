<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import ConversationList from '@/views/ConversationList.vue'

const userStore = useUserStore()
const route = useRoute()
const showSidebar = ref(true)
const conversationListRef = ref<InstanceType<typeof ConversationList> | null>(null)

onMounted(async () => {
  await userStore.initialize()
})

function toggleSidebar() {
  showSidebar.value = !showSidebar.value
}

function handleRefreshConversations() {
  conversationListRef.value?.loadConversations()
}
</script>

<template>
  <v-app>
    <template v-if="userStore.isLoggedIn && route.path !== '/login' && route.path !== '/register'">
      <v-row no-gutters class="fill-height">
        <v-col v-if="showSidebar" cols="12" md="3" lg="2" class="conversation-sidebar fill-height bg-sidebar">
          <ConversationList ref="conversationListRef" class="fill-height" />
        </v-col>

        <v-col cols="12" :md="showSidebar ? 9 : 12" :lg="showSidebar ? 10 : 12" class="chat-main fill-height bg-chat">
          <RouterView @toggle-sidebar="toggleSidebar" @refresh-conversations="handleRefreshConversations" />
        </v-col>
      </v-row>
    </template>
    <template v-else>
      <RouterView />
    </template>
  </v-app>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600&display=swap');

html,
body {
  height: 100%;
}

body {
  font-family:
    'IBM Plex Sans',
    -apple-system,
    BlinkMacSystemFont,
    'Segoe UI',
    Roboto,
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
