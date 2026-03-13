<script setup lang="ts">
import { statusColors } from '@/plugins/vuetify'

defineProps<{
  stage: number
}>()

type StageKey = 'uploading' | 'extracting' | 'analyzing' | 'generating' | 'completed' | 'failed' | 'unknown'

const stageMap: Record<number, { key: StageKey; text: string }> = {
  1: { key: 'uploading', text: '上传中' },
  2: { key: 'extracting', text: '解析文件中' },
  3: { key: 'analyzing', text: 'AI 分析中' },
  4: { key: 'generating', text: '生成文档中' },
  5: { key: 'completed', text: '已完成' },
  6: { key: 'failed', text: '处理失败' },
}

const getStageInfo = (stage: number): { key: StageKey; text: string } => {
  return stageMap[stage] || { key: 'unknown', text: '未知状态' }
}

const getStatusColor = (key: StageKey) => {
  return statusColors[key] || statusColors.unknown
}

const isSpinning = (stage: number) => {
  return stage >= 1 && stage <= 4
}
</script>

<template>
  <v-chip :color="getStatusColor(getStageInfo(stage).key).bg" :class="['stage-indicator']" label size="small">
    <v-progress-circular
      v-if="isSpinning(stage)"
      indeterminate
      size="12"
      width="2"
      class="mr-1"
      :color="getStatusColor(getStageInfo(stage).key).text"
    />
    <span
      :style="{
        color: `rgb(var(--v-theme-${getStatusColor(getStageInfo(stage).key).text}))`,
      }"
    >
      {{ getStageInfo(stage).text }}
    </span>
  </v-chip>
</template>

<style scoped>
.stage-indicator {
  font-weight: 500;
  font-size: 14px;
}
</style>
