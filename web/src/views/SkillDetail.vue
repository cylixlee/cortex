<template>
  <div class="skill-detail">
    <div class="header">
      <router-link to="/skills" class="back-link">← Back to Skills</router-link>
      <h1>{{ skill?.name || 'Skill Details' }}</h1>
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Loading skill...</p>
    </div>
    <div v-else-if="error && !skill" class="error">{{ error }}</div>

    <template v-else-if="skill">
      <div class="status-section">
        <StageIndicator :stage="skill.stage" />
      </div>

      <div v-if="skill.stage < 5" class="processing-notice">
        <p>Your skill is being processed. This may take a few minutes.</p>
      </div>

      <div v-else-if="skill.stage === 5 && skill.skill" class="content-section">
        <div class="section">
          <h2>Overview</h2>
          <div class="markdown-content" v-html="renderMarkdown(skill.skill.overview)"></div>
        </div>

        <div v-if="skill.skill.references?.length" class="section">
          <h2>References</h2>
          <div v-for="ref in skill.skill.references" :key="ref.filename" class="reference-item">
            <h3>{{ ref.filename }}</h3>
            <div class="markdown-content" v-html="renderMarkdown(ref.content)"></div>
          </div>
        </div>

        <div class="actions">
          <button class="btn-download" @click="handleDownload">Download Skill Package</button>
        </div>
      </div>

      <div v-else-if="skill.stage === 6" class="error-section">
        <p>Processing failed. Please try again.</p>
        <button class="btn-retry" @click="handleDelete">Delete and Retry</button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSkill, deleteSkill, subscribeSkillStatus, downloadSkill } from '@/api/skill'
import StageIndicator from '@/components/StageIndicator.vue'

const route = useRoute()
const router = useRouter()

interface SkillData {
  id: string
  name: string
  status: string
  stage: number
  skill?: {
    overview: string
    references: Array<{
      filename: string
      content: string
    }>
  }
}

const skill = ref<SkillData | null>(null)
const loading = ref(true)
const error = ref('')
let unsubscribe: (() => void) | null = null

const loadSkill = async () => {
  loading.value = true
  error.value = ''
  try {
    skill.value = await getSkill(route.params.id as string)

    if (skill.value.stage < 5) {
      unsubscribe = await subscribeSkillStatus(route.params.id as string, (stage, name) => {
        if (skill.value) {
          skill.value.stage = stage
          skill.value.status = name
        }
      })
    }
  } catch {
    error.value = 'Failed to load skill'
  } finally {
    loading.value = false
  }
}

const renderMarkdown = (content: string) => {
  return content
    .replace(/^### (.+)$/gm, '<h3>$1</h3>')
    .replace(/^## (.+)$/gm, '<h2>$1</h2>')
    .replace(/^# (.+)$/gm, '<h1>$1</h1>')
    .replace(/```(\w+)?\n([\s\S]*?)```/g, '<pre><code>$2</code></pre>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\n\n/g, '</p><p>')
    .replace(/^(.+)$/gm, (match) => {
      if (match.startsWith('<h') || match.startsWith('<pre') || match.startsWith('<p')) {
        return match
      }
      return `<p>${match}</p>`
    })
}

const handleDownload = async () => {
  if (!skill.value) return
  try {
    await downloadSkill(skill.value.id, skill.value.name)
  } catch {
    alert('Failed to download skill')
  }
}

const handleDelete = async () => {
  if (!confirm('Are you sure you want to delete this skill?')) return
  try {
    await deleteSkill(route.params.id as string)
    router.push('/skills')
  } catch {
    alert('Failed to delete skill')
  }
}

onMounted(loadSkill)

onUnmounted(() => {
  if (unsubscribe) {
    unsubscribe()
  }
})
</script>

<style scoped>
.skill-detail {
  padding: 24px;
  max-width: 900px;
  margin: 0 auto;
  height: 100%;
  overflow-y: auto;
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
  margin: 0;
}

.loading,
.error {
  text-align: center;
  padding: 48px;
  color: #6b7280;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.mini-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid #e5e7eb;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error {
  color: #ef4444;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  text-transform: uppercase;
  font-weight: 500;
}

.status-badge.pending {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.processing {
  background: #dbeafe;
  color: #1e40af;
}

.status-badge.completed {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.failed {
  background: #fee2e2;
  color: #991b1b;
}

.progress {
  color: #6b7280;
}

.processing-notice {
  background: #f9fafb;
  padding: 16px;
  border-radius: 8px;
  text-align: center;
}

.content-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 24px;
}

.section {
  margin-bottom: 32px;
}

.section h2 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e5e7eb;
}

.reference-item {
  margin-bottom: 24px;
}

.reference-item h3 {
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 12px;
}

.markdown-content {
  line-height: 1.6;
  color: #374151;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3) {
  margin-top: 24px;
  margin-bottom: 12px;
}

.markdown-content :deep(pre) {
  background: #f3f4f6;
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
}

.markdown-content :deep(code) {
  font-family: monospace;
  font-size: 13px;
}

.actions {
  margin-top: 32px;
  display: flex;
  gap: 12px;
}

.btn-download {
  background: #4f46e5;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.btn-download:hover {
  background: #4338ca;
}

.error-section {
  text-align: center;
  padding: 24px;
}

.btn-retry {
  background: #ef4444;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
</style>
