<template>
  <v-container style="max-width: 600px">
    <v-btn variant="text" prepend-icon="mdi-arrow-left" to="/skills" class="mb-4"> Back to Skills </v-btn>

    <h1 class="text-h4 font-weight-bold mb-6">Upload New Skill</h1>

    <v-card class="pa-6">
      <v-form @submit.prevent="handleUpload">
        <v-text-field
          v-model="name"
          label="Skill Name"
          placeholder="e.g., my-awesome-lib"
          prepend-inner-icon="mdi-tag"
          :disabled="uploading"
          class="mb-4"
        />

        <v-file-input
          v-model="selectedFile"
          label="Source Code (ZIP)"
          accept=".zip"
          prepend-icon="mdi-folder-zip"
          chips
          show-size
          :disabled="uploading"
          class="mb-4"
        >
          <template #selection="{ fileNames }">
            <v-chip v-for="fileName in fileNames" :key="fileName" variant="tonal" color="secondary" size="small">
              {{ fileName }}
            </v-chip>
          </template>
        </v-file-input>

        <v-alert v-if="error" type="error" variant="tonal" closable class="mb-4" @click:close="error = ''">
          {{ error }}
        </v-alert>

        <v-btn
          variant="tonal"
          color="secondary"
          size="large"
          block
          :disabled="!canUpload"
          :loading="uploading"
          type="submit"
        >
          <v-icon icon="mdi-upload" class="mr-2"></v-icon>
          Upload Skill
        </v-btn>
      </v-form>
    </v-card>

    <v-card v-if="uploading" class="mt-6 pa-6">
      <h3 class="text-h6 mb-4">Processing Status</h3>

      <div class="d-flex align-center gap-3">
        <v-progress-circular
          :indeterminate="stage < 5"
          :model-value="stageProgress"
          variant="tonal"
          color="secondary"
        ></v-progress-circular>

        <StageIndicator :stage="stage" />
      </div>

      <p class="text-caption text-muted mt-4">You can check the status later in the skill list</p>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { uploadSkill, subscribeSkillStatus } from '@/api/skill'
import StageIndicator from '@/components/StageIndicator.vue'

const router = useRouter()

const name = ref('')
const selectedFile = ref<File[] | null>(null)
const uploading = ref(false)
const error = ref('')
const stage = ref(1)

const canUpload = computed(() => {
  return name.value.trim() && selectedFile.value && selectedFile.value.length > 0
})

const stageProgress = computed(() => {
  return ((stage.value - 1) / 4) * 100
})

const handleUpload = async () => {
  if (!name.value || !selectedFile.value || selectedFile.value.length === 0) return

  error.value = ''
  uploading.value = true
  stage.value = 1

  const file = selectedFile.value[0]
  if (!file) return

  try {
    const result = await uploadSkill(name.value, file, () => {})

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
