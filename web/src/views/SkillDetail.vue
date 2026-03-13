<template>
  <v-container style="max-width: 900px">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" to="/skills" class="mb-4"> Back to Skills </v-btn>

    <div v-if="loading" class="d-flex flex-column align-center justify-center py-12">
      <v-progress-circular indeterminate color="primary" size="64" class="mb-4"></v-progress-circular>
      <p class="text-muted">Loading skill...</p>
    </div>

    <v-alert v-else-if="error && !skill" type="error" variant="tonal" closable>
      {{ error }}
    </v-alert>

    <template v-else-if="skill">
      <div class="d-flex justify-space-between align-center mb-6">
        <h1 class="text-h4 font-weight-bold">{{ skill.name }}</h1>

        <StageIndicator :stage="skill.stage" />
      </div>

      <v-card v-if="skill.stage < 5" class="pa-6 mb-6" color="surface">
        <div class="d-flex align-center gap-3">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
          <div>
            <p class="text-body-1 mb-1">Your skill is being processed.</p>
            <p class="text-caption text-muted">This may take a few minutes.</p>
          </div>
        </div>
      </v-card>

      <v-card v-else-if="skill.stage === 5 && skill.skill" class="pa-6">
        <section class="mb-8">
          <h2 class="text-h5 font-weight-bold mb-4">Overview</h2>
          <div class="markdown-content" v-html="renderMarkdown(skill.skill.overview)"></div>
        </section>

        <section v-if="skill.skill.references?.length" class="mb-8">
          <h2 class="text-h5 font-weight-bold mb-4">References</h2>

          <v-expansion-panels variant="accordion">
            <v-expansion-panel v-for="ref in skill.skill.references" :key="ref.filename">
              <v-expansion-panel-title>
                <v-icon icon="mdi-file-document-outline" class="mr-2"></v-icon>
                {{ ref.filename }}
              </v-expansion-panel-title>

              <v-expansion-panel-text>
                <div class="markdown-content" v-html="renderMarkdown(ref.content)"></div>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </section>

        <v-divider class="mb-6"></v-divider>

        <div class="d-flex gap-3">
          <v-btn variant="tonal" color="secondary" size="large" prepend-icon="mdi-download" @click="handleDownload">
            Download Skill Package
          </v-btn>
        </div>
      </v-card>

      <v-card v-else-if="skill.stage === 6" class="pa-6 text-center" color="error">
        <v-icon icon="mdi-alert-circle" size="48" class="mb-4"></v-icon>
        <p class="text-h6 mb-4">Processing failed. Please try again.</p>
        <v-btn color="white" variant="outlined" @click="handleDelete"> Delete and Retry </v-btn>
      </v-card>
    </template>
  </v-container>
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

    if (skill.value!.stage < 5) {
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
.markdown-content {
  line-height: 1.7;
}

.markdown-content :deep(h1),
.markdown-content :deep(h2),
.markdown-content :deep(h3) {
  margin-top: 1.5em;
  margin-bottom: 0.5em;
}

.markdown-content :deep(pre) {
  background: rgb(var(--v-theme-surface));
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-content :deep(code) {
  font-family: monospace;
}
</style>
