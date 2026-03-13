<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'

const props = defineProps<{
  disabled?: boolean
  loading?: boolean
  modelValue?: boolean
}>()

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
          :color="enableRag ? 'primary' : 'surface-variant'"
          :class="!enableRag ? 'text-black' : 'text-white'"
          :variant="enableRag ? 'flat' : 'flat'"
          size="default"
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
  border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.input-wrapper {
  max-width: 900px;
  margin: 0 auto;
}

.input-box {
  position: relative;
  border-radius: 24px;
  overflow: hidden;
}

.message-textarea :deep(.v-field) {
  border-radius: 24px;
  padding-bottom: 40px;
}

.message-textarea :deep(.v-field__field) {
  padding-bottom: 0;
}

.message-textarea :deep(.v-field__input) {
  padding-top: 16px;
  padding-bottom: 16px;
  min-height: 56px !important;
}

.message-textarea :deep(.v-field::after) {
  display: none;
}

.message-textarea :deep(.v-field__wrapper::before) {
  display: none !important;
}

.send-btn-wrapper {
  position: absolute;
  bottom: 8px;
  right: 12px;
}
</style>
