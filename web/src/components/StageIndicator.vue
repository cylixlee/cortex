<script setup lang="ts">
defineProps<{
  stage: number
}>()

const getStageInfo = (stage: number): { text: string; class: string; isSpinning: boolean } => {
  switch (stage) {
    case 1:
      return { text: '上传中', class: 'stage-uploading', isSpinning: true }
    case 2:
      return { text: '解析文件中', class: 'stage-extracting', isSpinning: true }
    case 3:
      return { text: 'AI 分析中', class: 'stage-analyzing', isSpinning: true }
    case 4:
      return { text: '生成文档中', class: 'stage-generating', isSpinning: true }
    case 5:
      return { text: '已完成', class: 'stage-completed', isSpinning: false }
    case 6:
      return { text: '处理失败', class: 'stage-failed', isSpinning: false }
    default:
      return { text: '未知状态', class: 'stage-unknown', isSpinning: false }
  }
}
</script>

<template>
  <div :class="['stage-indicator', getStageInfo(stage).class]">
    <span v-if="getStageInfo(stage).isSpinning" class="spinner"></span>
    <span class="stage-text">{{ getStageInfo(stage).text }}</span>
  </div>
</template>

<style scoped>
.stage-indicator {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 500;
}

.stage-text {
  white-space: nowrap;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.stage-uploading {
  background-color: #dbeafe;
  color: #1e40af;
}

.stage-extracting {
  background-color: #fef3c7;
  color: #92400e;
}

.stage-analyzing {
  background-color: #e0e7ff;
  color: #4338ca;
}

.stage-generating {
  background-color: #fce7f3;
  color: #9d174d;
}

.stage-completed {
  background-color: #d1fae5;
  color: #065f46;
}

.stage-failed {
  background-color: #fee2e2;
  color: #991b1b;
}

.stage-unknown {
  background-color: #f3f4f6;
  color: #374151;
}
</style>
