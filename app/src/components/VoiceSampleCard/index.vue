<template>
  <GlCard title="语音样本（选填）" :shadow="false" bordered>
    <view class="voice-upload">
      <view v-if="sample" class="voice-item">
        <view class="voice-play" @tap="$emit('play')">
          <uv-icon :name="isPlaying ? 'pause-circle' : 'play-circle'" size="32" color="var(--color-primary)"></uv-icon>
        </view>
        <text class="voice-duration">{{ duration }}s</text>
        <view class="voice-delete" @tap="$emit('delete')">
          <uv-icon name="close" size="16" color="var(--color-text-secondary)"></uv-icon>
        </view>
      </view>

      <view
        v-else
        class="voice-record"
        @touchstart="$emit('record-start')"
        @touchend="$emit('record-end')"
      >
        <uv-icon name="mic" size="32" :color="recording ? 'var(--color-icon-accent)' : 'var(--color-text-secondary)'"></uv-icon>
        <text class="record-text">{{ recording ? '松开结束' : '按住录制语音' }}</text>
      </view>
    </view>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  sample?: string
  duration?: number
  recording?: boolean
  isPlaying?: boolean
}

withDefaults(defineProps<Props>(), {
  sample: '',
  duration: 0,
  recording: false,
  isPlaying: false,
})

defineEmits<{
  play: []
  delete: []
  'record-start': []
  'record-end': []
}>()
</script>

<style lang="scss" scoped>
.voice-upload {
  padding: var(--spacing-sm) 0;
}

.voice-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
}

.voice-duration {
  flex: 1;
  font-size: var(--font-md);
  color: var(--color-text);
}

.voice-record {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx dashed var(--color-border);
  @include press-effect;
}

.record-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}
</style>
