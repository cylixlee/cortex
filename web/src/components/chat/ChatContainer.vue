<script setup lang="ts">
import { layout } from '@/plugins/vuetify'

withDefaults(
  defineProps<{
    maxMessageWidth?: string
  }>(),
  {
    maxMessageWidth: '70%',
  },
)
</script>

<template>
  <v-main class="chat-container-main">
    <v-sheet class="chat-container-sheet fill-height" color="chat">
      <div class="messages-scroll">
        <div class="messages-wrapper w-100 h-100" :style="{ maxWidth: maxMessageWidth || layout.chatMaxWidth }">
          <slot />
        </div>
      </div>
    </v-sheet>
  </v-main>
</template>

<style scoped>
.chat-container-main {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.chat-container-sheet {
  flex: 1;
  overflow: hidden;
}

.messages-scroll {
  height: 100%;
  overflow-y: auto;
}

.messages-wrapper {
  margin: 0 auto;
}
</style>
