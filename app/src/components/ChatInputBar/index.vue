<template>
  <view
    class="input-bar"
    :style="{ paddingBottom: resolvedPaddingBottom + 'px', transitionDuration: duration + 's' }"
  >
    <view class="input-tools">
      <view class="tool-btn" :class="{ active: showVoice }" @tap="toggleVoice">
        <uv-icon :name="showVoice ? 'edit-pen' : 'mic'" size="24" color="var(--color-text-secondary)" />
      </view>
    </view>

    <!-- 文字输入 -->
    <view v-if="!showVoice" class="input-wrap">
      <textarea
        class="chat-textarea"
        :value="modelValue"
        :auto-height="true"
        :maxlength="500"
        :show-confirm-bar="false"
        :adjust-position="false"
        :cursor-spacing="20"
        fixed
        disable-default-padding
        placeholder="输入消息..."
        @input="handleInput"
        @confirm="$emit('send')"
        @focus="handleFocus"
        @blur="handleBlur"
      />
    </view>

    <!-- 语音输入 -->
    <view
      v-else
      class="voice-input-wrap"
      :class="{ recording: recording }"
      @touchstart="$emit('record-start')"
      @touchend="$emit('record-end')"
      @touchcancel="$emit('record-cancel')"
    >
      <uv-icon :name="recording ? 'mic-filled' : 'mic'" size="20" :color="recording ? '#fff' : 'var(--color-text-secondary)'" />
      <text>{{ recording ? '松开发送' : '按住说话' }}</text>
    </view>

    <view class="input-actions">
      <!-- 有文字时显示发送按钮 -->
      <view v-if="modelValue?.trim()" class="send-btn" @tap.stop="$emit('send')">
        <uv-icon name="arrow-up" size="20" color="#fff" />
      </view>
      <!-- 无文字时显示功能按钮 -->
      <template v-else>
        <view class="tool-btn" @tap="$emit('image')">
          <uv-icon name="photo" size="24" color="var(--color-text-secondary)" />
        </view>
        <view class="tool-btn" @tap="$emit('more')">
          <uv-icon name="plus-circle" size="24" color="var(--color-text-secondary)" />
        </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  modelValue?: string
  recording?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  recording: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  send: []
  focus: []
  blur: []
  image: []
  more: []
  'record-start': []
  'record-end': []
  'record-cancel': []
  'keyboard-height-change': [height: number]
}>()

const showVoice = ref(false)
const keyboardHeight = ref(0)
const duration = ref(0.25)
const isFocus = ref(false)

// 计算底部 padding：键盘高度 + 安全区
// 当键盘弹起时，键盘高度已包含安全区，所以不需要额外加 safe-area
// 当键盘收起时，需要保留 safe-area
const resolvedPaddingBottom = computed(() => {
  const safeArea = uni.getSystemInfoSync().safeAreaInsets?.bottom || 0
  if (keyboardHeight.value > 0) {
    return keyboardHeight.value
  }
  return 12 + safeArea // 默认 12px 间距 + 安全区
})

const toggleVoice = () => {
  showVoice.value = !showVoice.value
  if (showVoice.value) {
    uni.hideKeyboard()
  }
}

const handleInput = (e: any) => {
  emit('update:modelValue', e.detail.value)
}

const handleFocus = (e: any) => {
  isFocus.value = true
  emit('focus')
}

const handleBlur = () => {
  isFocus.value = false
  emit('blur')
}

const onKeyboardHeightChange = (res: any) => {
  // iOS 下 duration 为 0.25s，Android 可能不同，这里做统一处理
  if (res.duration > 0) {
    duration.value = res.duration
  }
  keyboardHeight.value = res.height
  emit('keyboard-height-change', res.height)
}

onMounted(() => {
  uni.onKeyboardHeightChange(onKeyboardHeightChange)
})

onUnmounted(() => {
  uni.offKeyboardHeightChange(onKeyboardHeightChange)
})
</script>

<style lang="scss" scoped>
.input-bar {
  position: relative;
  display: flex;
  align-items: flex-end;
  gap: $gl-spacing-sm;
  padding: $gl-spacing-md $gl-spacing-md;
  background: rgba(255, 255, 255, 0.98); // 微透明背景
  backdrop-filter: blur(20px);
  border-top: 1rpx solid rgba(0, 0, 0, 0.05);
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.02);
  transition-property: padding-bottom;
  transition-timing-function: cubic-bezier(0.25, 0.1, 0.25, 1);
  z-index: 100;
}

.input-tools, .input-actions {
  display: flex;
  align-items: center;
  gap: $gl-spacing-sm;
  height: 72rpx; // 对齐输入框高度
}

.tool-btn {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: $gl-radius-circle;
  transition: all 0.2s;

  &:active {
    background: rgba(0, 0, 0, 0.05);
    transform: scale(0.95);
  }

  &.active {
    color: $gl-color-primary;
    background: rgba($gl-color-primary, 0.1);
  }
}

.input-wrap {
  flex: 1;
  min-height: 72rpx;
  background: #F7F8FA;
  border-radius: $gl-radius-lg;
  padding: 16rpx 24rpx;
  display: flex;
  align-items: center;
}

.chat-textarea {
  width: 100%;
  max-height: 200rpx;
  font-size: 30rpx;
  color: $gl-text-main;
  line-height: 1.5;
  padding: 0;
  margin: 0;
}

.voice-input-wrap {
  flex: 1;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $gl-spacing-xs;
  background: #F7F8FA;
  border-radius: $gl-radius-lg;
  font-size: 30rpx;
  font-weight: 500;
  color: $gl-text-main;
  transition: all 0.2s;

  &:active, &.recording {
    background: $gl-color-primary;
    color: #fff;
    transform: scale(0.98);
  }
}

.send-btn {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $gl-color-primary;
  border-radius: $gl-radius-circle;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); // 弹性动效

  &:active {
    transform: scale(0.9);
    background: $gl-color-primary-dark;
  }
}
</style>
