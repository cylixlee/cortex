<template>
  <div class="upload-page">
    <div class="header">
      <router-link to="/skills" class="back-link">← Back to Skills</router-link>
      <h1>Upload New Skill</h1>
    </div>

    <div class="upload-form">
      <div class="form-group">
        <label>Skill Name</label>
        <input v-model="name" type="text" placeholder="e.g., my-awesome-lib" :disabled="uploading" />
      </div>

      <div class="form-group">
        <label>Source Code (ZIP)</label>
        <div
          class="dropzone"
          :class="{ dragging: isDragging, disabled: uploading }"
          @dragover.prevent="isDragging = true"
          @dragleave="isDragging = false"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <input ref="fileInput" type="file" accept=".zip" @change="handleFileSelect" :disabled="uploading" hidden />
          <div v-if="!selectedFile" class="dropzone-text">
            <p>Drag & drop a ZIP file here, or click to select</p>
            <p class="hint">Only .zip files are supported</p>
          </div>
          <div v-else class="selected-file">
            <span class="file-name">{{ selectedFile.name }}</span>
            <span class="file-size">{{ formatFileSize(selectedFile.size) }}</span>
          </div>
        </div>
      </div>

      <div v-if="error" class="error">{{ error }}</div>

      <button class="btn-upload" :disabled="!canUpload || uploading" @click="handleUpload">
        {{ uploading ? 'Uploading...' : 'Upload Skill' }}
      </button>
    </div>

    <div v-if="uploading" class="progress-section">
      <h3>Processing Status</h3>
      <StageIndicator :stage="stage" />
      <p class="progress-hint">You can check the status later in the skill list</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { uploadSkill, subscribeSkillStatus } from '@/api/skill'
import StageIndicator from '@/components/StageIndicator.vue'

const router = useRouter()

const name = ref('')
const selectedFile = ref<File | null>(null)
const isDragging = ref(false)
const uploading = ref(false)
const error = ref('')
const stage = ref(1) // 1 = pending/uploading

const fileInput = ref<HTMLInputElement | null>(null)

const canUpload = computed(() => {
  return name.value.trim() && selectedFile.value
})

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0]
  }
}

const handleDrop = (e: DragEvent) => {
  isDragging.value = false
  if (e.dataTransfer?.files && e.dataTransfer.files[0]) {
    const file = e.dataTransfer.files[0]
    if (file.name.endsWith('.zip')) {
      selectedFile.value = file
    } else {
      error.value = 'Only ZIP files are supported'
    }
  }
}

const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

const handleUpload = async () => {
  if (!name.value || !selectedFile.value) return

  error.value = ''
  uploading.value = true
  stage.value = 1

  try {
    const result = await uploadSkill(name.value, selectedFile.value, () => {})

    stage.value = 2

    const unsubscribe = await subscribeSkillStatus(result.skill_id, (newStage, newName) => {
      stage.value = newStage

      if (newName === 'completed' || newName === 'failed') {
        unsubscribe()
        router.push(`/skills/${result.skill_id}`)
      }
    })
  } catch {
    error.value = 'Failed to upload skill'
    uploading.value = false
  }
}
</script>

<style scoped>
.upload-page {
  padding: 24px;
  max-width: 600px;
  margin: 0 auto;
}

.header {
  margin-bottom: 24px;
}

.back-link {
  display: inline-block;
  color: #6b7280;
  text-decoration: none;
  font-size: 14px;
  margin-bottom: 8px;
}

.back-link:hover {
  color: #4f46e5;
}

h1 {
  font-size: 24px;
  font-weight: 600;
}

.upload-form {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
  color: #374151;
}

.form-group input[type='text'] {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input[type='text']:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.dropzone {
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s;
}

.dropzone:hover:not(.disabled) {
  border-color: #4f46e5;
}

.dropzone.dragging {
  border-color: #4f46e5;
  background: #f5f3ff;
}

.dropzone.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.dropzone-text p {
  color: #6b7280;
  margin: 0;
}

.dropzone-text .hint {
  font-size: 12px;
  margin-top: 4px;
}

.selected-file {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-name {
  font-weight: 500;
}

.file-size {
  color: #6b7280;
  font-size: 14px;
}

.error {
  color: #ef4444;
  margin-bottom: 16px;
}

.btn-upload {
  width: 100%;
  background: #4f46e5;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.btn-upload:hover:not(:disabled) {
  background: #4338ca;
}

.btn-upload:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.progress-section {
  margin-top: 24px;
  padding: 24px;
  background: #f9fafb;
  border-radius: 8px;
}

.progress-section h3 {
  margin: 0 0 16px;
  font-size: 16px;
}

.progress-bar {
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #4f46e5;
  transition: width 0.3s;
}

.progress-text {
  margin: 8px 0;
  font-size: 14px;
  color: #6b7280;
}

.progress-hint {
  font-size: 12px;
  color: #9ca3af;
  margin: 0;
}
</style>
