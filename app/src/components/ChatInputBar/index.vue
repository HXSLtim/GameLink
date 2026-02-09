<template>
  <view class="input-bar">
    <view class="input-tools">
      <view class="tool-btn" :class="{ active: showVoice }" @tap="showVoice = !showVoice">
        <uv-icon :name="showVoice ? 'edit-pen' : 'mic'" size="20" color="var(--color-text-secondary)" />
      </view>
    </view>
    
    <!-- 文字输入 -->
    <GlInput
      v-if="!showVoice"
      class="input-wrap"
      :model-value="modelValue"
      type="textarea"
      size="small"
      placeholder="输入消息..."
      :maxlength="500"
      :auto-height="true"
      :show-confirm-bar="false"
      :adjust-position="true"
      @update:modelValue="(value) => $emit('update:modelValue', value)"
      @confirm="$emit('send')"
      @focus="$emit('focus')"
    />
    
    <!-- 语音输入 -->
    <view 
      v-else 
      class="voice-input-wrap"
      :class="{ recording: recording }"
      @touchstart="$emit('record-start')"
      @touchend="$emit('record-end')"
      @touchcancel="$emit('record-cancel')"
    >
      <uv-icon :name="recording ? 'mic' : 'mic'" size="16" :color="recording ? 'var(--color-primary)' : 'var(--color-text-secondary)'" />
      <text>{{ recording ? '松开发送' : '按住说话' }}</text>
    </view>
    
    <view class="input-actions">
      <!-- 有文字时显示发送按钮 -->
      <view v-if="modelValue?.trim()" class="send-btn" @click="$emit('send')">
        <uv-icon name="play-right-fill" size="20" color="#fff" />
      </view>
      <!-- 无文字时显示功能按钮 -->
      <template v-else>
        <view class="tool-btn" @tap="$emit('image')">
          <uv-icon name="photo" size="20" color="var(--color-text-secondary)" />
        </view>
        <view class="tool-btn" @tap="$emit('more')">
          <uv-icon name="plus-circle" size="20" color="var(--color-text-secondary)" />
        </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  modelValue?: string
  recording?: boolean
}

withDefaults(defineProps<Props>(), {
  modelValue: '',
  recording: false,
})

defineEmits<{
  'update:modelValue': [value: string]
  send: []
  focus: []
  image: []
  more: []
  'record-start': []
  'record-end': []
  'record-cancel': []
}>()

const showVoice = ref(false)
</script>

<style lang="scss" scoped>
.input-bar {
  display: flex;
  align-items: flex-end;
  gap: var(--spacing-md);
  padding: var(--spacing-md) var(--spacing-lg);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.input-tools {
  display: flex;
  gap: var(--spacing-sm);
  padding-bottom: var(--spacing-xs);
}

.tool-btn {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: var(--color-bg-secondary);
  }

  &:active {
    transform: scale(0.9);
  }

  &.active {
    background: rgba(var(--color-primary-rgb), 0.1);
  }
}

.input-wrap {
  flex: 1;
  border-radius: var(--radius-lg);

  :deep(.gl-input__textarea) {
    max-height: 200rpx;
    font-size: var(--font-sm);
    color: var(--color-text);
    line-height: 1.5;
  }
}

.voice-input-wrap {
  flex: 1;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-lg);
  border: 1rpx solid var(--color-border);
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all 0.2s;

  &:active,
  &.recording {
    background: rgba(var(--color-primary-rgb), 0.08);
    border-color: var(--color-primary);
    color: var(--color-primary);
  }
}

.input-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding-bottom: var(--spacing-xs);
}

.send-btn {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4rpx 12rpx rgba(var(--color-primary-rgb), 0.3);

  &:hover {
    transform: scale(1.05);
  }

  &:active {
    transform: scale(0.92);
    box-shadow: 0 2rpx 6rpx rgba(var(--color-primary-rgb), 0.2);
  }
}
</style>
