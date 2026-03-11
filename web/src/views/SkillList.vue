<template>
  <div class="skill-list">
    <div class="header">
      <h1>My Skills</h1>
      <router-link to="/skills/upload" class="btn-primary"> Upload New Skill </router-link>
    </div>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="skills.length === 0" class="empty">
      <p>No skills yet. Upload your first skill to get started.</p>
    </div>
    <div v-else class="skills-grid">
      <div v-for="skill in skills" :key="skill.id" class="skill-card" @click="goToDetail(skill.id)">
        <div class="skill-name">{{ skill.name }}</div>
        <div class="skill-status">
          <StageIndicator :stage="skill.stage" />
        </div>
        <div class="skill-date">{{ formatDate(skill.created_at) }}</div>
        <button class="btn-delete" @click.stop="handleDelete(skill.id)">Delete</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { listSkills, deleteSkill } from '@/api/skill'
import StageIndicator from '@/components/StageIndicator.vue'

const router = useRouter()
const skills = ref<
  Array<{
    id: string
    name: string
    status: string
    stage: number
    created_at: string
  }>
>([])
const loading = ref(true)
const error = ref('')

const loadSkills = async () => {
  loading.value = true
  error.value = ''
  try {
    skills.value = await listSkills()
  } catch {
    error.value = 'Failed to load skills'
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: string) => {
  router.push(`/skills/${id}`)
}

const handleDelete = async (id: string) => {
  if (!confirm('Are you sure you want to delete this skill?')) return
  try {
    await deleteSkill(id)
    await loadSkills()
  } catch {
    alert('Failed to delete skill')
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString()
}

onMounted(loadSkills)
</script>

<style scoped>
.skill-list {
  padding: 24px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  font-size: 24px;
  font-weight: 600;
}

.btn-primary {
  background: #4f46e5;
  color: white;
  padding: 8px 16px;
  border-radius: 6px;
  text-decoration: none;
  font-size: 14px;
}

.btn-primary:hover {
  background: #4338ca;
}

.loading,
.error,
.empty {
  text-align: center;
  padding: 48px;
  color: #6b7280;
}

.error {
  color: #ef4444;
}

.skills-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.skill-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  cursor: pointer;
  transition: box-shadow 0.2s;
}

.skill-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.skill-name {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.skill-status {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.status-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  text-transform: uppercase;
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
  font-size: 12px;
  color: #6b7280;
}

.skill-date {
  font-size: 12px;
  color: #9ca3af;
}

.btn-delete {
  margin-top: 8px;
  background: #ef4444;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.btn-delete:hover {
  background: #dc2626;
}
</style>
