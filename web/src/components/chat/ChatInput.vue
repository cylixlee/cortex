<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { layout } from '@/plugins/vuetify'

const props = withDefaults(
  defineProps<{
    disabled?: boolean
    loading?: boolean
    modelValue?: boolean
    maxWidth?: string
  }>(),
  {
    maxWidth: layout.chatMaxWidth,
  },
)

const emit = defineEmits<{
  (e: 'send', content: string): void
  (e: 'update:modelValue', value: boolean): void
}>()

const input = ref('')
const textareaRef = ref<InstanceType<typeof import('vuetify/components').VTextarea> | null>(null)

const enableRag = ref(props.modelValue ?? false)

watch(
  () => props.modelValue,
  (val) => {
    enableRag.value = val ?? false
  },
)

watch(enableRag, (val) => {
  emit('update:modelValue', val)
})

watch(input, async () => {
  await nextTick()
  adjustHeight()
})

function adjustHeight() {
  const el = document.querySelector('.message-textarea .v-field__input') as HTMLElement
  if (el) {
    el.style.height = 'auto'
    el.style.height = Math.min(el.scrollHeight, 200) + 'px'
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey && !e.ctrlKey && !e.metaKey) {
    e.preventDefault()
    handleSend()
  }
}

function handleSend() {
  const content = input.value.trim()
  if (!content || props.disabled || props.loading) return
  emit('send', content)
  input.value = ''
  setTimeout(() => {
    adjustHeight()
  }, 0)
}

defineExpose({
  focus: () => {
    textareaRef.value?.focus()
  },
})
</script>

<template>
  <div class="chat-input-container pa-4 bg-white">
    <div class="input-wrapper w-100">
      <div class="skill-toggle mb-2">
        <v-btn
          :color="enableRag ? 'primary' : 'button-secondary'"
          variant="flat"
          rounded="pill"
          @click="enableRag = !enableRag"
        >
          <v-icon start>mdi-brain</v-icon>
          Skill
        </v-btn>
      </div>

      <div class="input-box">
        <v-textarea
          ref="textareaRef"
          v-model="input"
          @keydown="handleKeydown"
          placeholder="Message Cortex..."
          :disabled="disabled"
          rows="1"
          auto-grow
          max-rows="8"
          hide-details
          variant="filled"
          elevation="0"
          class="message-textarea"
        ></v-textarea>

        <div class="send-btn-wrapper">
          <v-btn
            icon="mdi-arrow-up"
            color="primary"
            variant="flat"
            size="small"
            rounded="full"
            :disabled="disabled || loading || !input.trim()"
            @click="handleSend"
          ></v-btn>
        </div>
      </div>

      <div class="text-center mt-2">
        <span class="text-caption text-muted"> AI can make mistakes. Please verify important information. </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-input-container {
}

.input-wrapper {
  max-width: v-bind(maxWidth);
  margin: 0 auto;
}

.input-box {
  position: relative;
  border-radius: 24px;
  overflow: visible;
}

.message-textarea :deep(.v-field),
.message-textarea :deep(.v-field--focused) {
  border-radius: 24px;
  background: rgb(var(--v-theme-surface-variant));
  border: 2px solid rgb(var(--v-theme-surface-variant));
}

.message-textarea :deep(.v-field--focused) {
  border-color: rgb(var(--v-theme-primary));
}

.message-textarea :deep(.v-field::before),
.message-textarea :deep(.v-field::after) {
  display: none !important;
}

.message-textarea :deep(.v-field__outline) {
  display: none !important;
}

.message-textarea :deep(.v-field__input) {
  padding-top: 16px;
  padding-bottom: 16px;
  min-height: 56px !important;
}

.send-btn-wrapper {
  position: absolute;
  bottom: 8px;
  right: 12px;
}
</style>
