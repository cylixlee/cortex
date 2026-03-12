<template>
  <v-container>
    <div class="d-flex justify-space-between align-center mb-6">
      <h1 class="text-h4 font-weight-bold">My Skills</h1>

      <v-btn color="secondary" prepend-icon="mdi-upload" to="/skills/upload"> Upload New Skill </v-btn>
    </div>

    <v-alert v-if="error" type="error" variant="tonal" closable class="mb-4" @click:close="error = ''">
      {{ error }}
    </v-alert>

    <div v-if="loading" class="d-flex justify-center py-12">
      <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
    </div>

    <v-empty-state
      v-else-if="skills.length === 0"
      icon="mdi-brain"
      headline="No skills yet"
      text="Upload your first skill to get started."
    >
      <template #actions>
        <v-btn color="secondary" prepend-icon="mdi-upload" to="/skills/upload"> Upload Skill </v-btn>
      </template>
    </v-empty-state>

    <v-row v-else>
      <v-col v-for="skill in skills" :key="skill.id" cols="12" sm="6" md="4" lg="3">
        <v-card hover @click="goToDetail(skill.id)" class="skill-card h-100">
          <v-card-item>
            <template #prepend>
              <v-avatar color="secondary" variant="tonal">
                <v-icon icon="mdi-code-tags"></v-icon>
              </v-avatar>
            </template>

            <v-card-title>{{ skill.name }}</v-card-title>

            <v-card-subtitle>
              {{ formatDate(skill.created_at) }}
            </v-card-subtitle>
          </v-card-item>

          <v-card-actions>
            <StageIndicator :stage="skill.stage" />

            <v-spacer></v-spacer>

            <v-btn
              icon="mdi-delete"
              variant="text"
              color="error"
              size="small"
              @click.stop="handleDelete(skill.id)"
            ></v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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
.skill-card {
  transition:
    transform 0.2s,
    box-shadow 0.2s;
}

.skill-card:hover {
  transform: translateY(-4px);
}
</style>
