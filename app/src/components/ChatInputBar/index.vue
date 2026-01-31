<template>
  <view class="input-bar">
    <view class="input-tools">
      <view class="tool-btn" @tap="showVoice = !showVoice">
        <text>{{ showVoice ? '⌨️' : '🎤' }}</text>
      </view>
    </view>
    
    <!-- 文字输入 -->
    <view v-if="!showVoice" class="input-wrap">
      <textarea 
        :value="modelValue"
        class="message-input"
        placeholder="输入消息..."
        :maxlength="500"
        :auto-height="true"
        :show-confirm-bar="false"
        :adjust-position="true"
        @input="(e: any) => $emit('update:modelValue', e.detail.value)"
        @confirm="$emit('send')"
        @focus="$emit('focus')"
      />
    </view>
    
    <!-- 语音输入 -->
    <view 
      v-else 
      class="voice-input-wrap"
      @touchstart="$emit('record-start')"
      @touchend="$emit('record-end')"
      @touchcancel="$emit('record-cancel')"
    >
      <text>{{ recording ? '松开发送' : '按住说话' }}</text>
    </view>
    
    <view class="input-actions">
      <GlButton v-if="modelValue?.trim()" type="primary" size="small" round @click="$emit('send')">
        发送
      </GlButton>
      <template v-else>
        <view class="tool-btn" @tap="$emit('image')">
          <text>🖼️</text>
        </view>
        <view class="tool-btn" @tap="$emit('more')">
          <text>➕</text>
        </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'

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
  gap: 16rpx;
  padding: 16rpx 24rpx;
  padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
  background: var(--color-bg-card);
  border-top: 1rpx solid var(--color-border);
}

.input-tools {
  display: flex;
  gap: 8rpx;
}

.tool-btn {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36rpx;
}

.input-wrap {
  flex: 1;
  background: var(--color-bg-secondary);
  border-radius: 32rpx;
  border: 2rpx solid var(--color-border);
  padding: 16rpx 24rpx;
}

.message-input {
  width: 100%;
  max-height: 200rpx;
  font-size: 30rpx;
  color: var(--color-text);
  line-height: 1.4;
}

.voice-input-wrap {
  flex: 1;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: 32rpx;
  border: 2rpx solid var(--color-border);
  font-size: 28rpx;
  color: var(--color-text-secondary);
  
  &:active {
    background: var(--color-primary);
    border-color: var(--color-primary);
    color: #fff;
  }
}

.input-actions {
  display: flex;
  align-items: center;
  gap: 8rpx;
}
</style>
