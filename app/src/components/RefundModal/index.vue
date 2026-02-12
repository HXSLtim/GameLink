<template>
  <uv-popup :show="show" mode="bottom" round="16" @close="$emit('close')">
    <view class="refund-modal">
      <view class="modal-header">
        <text class="modal-title">申请退款</text>
        <uv-icon name="close" size="24" @click="$emit('close')"></uv-icon>
      </view>

      <!-- 退款原因选择 -->
      <view class="reason-section">
        <text class="section-label">请选择退款原因</text>
        <view class="reason-list">
          <view
            v-for="item in presetReasons"
            :key="item"
            class="reason-item"
            :class="{ 'reason-item--active': selectedReason === item }"
            @tap="selectReason(item)"
          >
            <uv-icon
              :name="selectedReason === item ? 'checkmark-circle-fill' : 'circle'"
              size="18"
              :color="selectedReason === item ? 'var(--color-primary)' : 'var(--color-text-placeholder)'"
            />
            <text class="reason-text">{{ item }}</text>
          </view>
        </view>
      </view>

      <!-- 自定义输入（选择"其他"时展示） -->
      <view v-if="selectedReason === '其他'" class="custom-section">
        <GlInput
          class="custom-input"
          :model-value="customReason"
          type="textarea"
          size="small"
          placeholder="请描述退款原因（必填）"
          :maxlength="200"
          @update:modelValue="(value) => $emit('update:customReason', value)"
        />
      </view>

      <!-- 退款提示 -->
      <view class="refund-tip">
        <uv-icon name="info-circle" size="14" color="var(--color-text-placeholder)" />
        <text class="tip-text">退款将在 1-3 个工作日内原路退回</text>
      </view>

      <GlButton
        type="primary"
        block
        size="large"
        :loading="loading"
        :disabled="!canSubmit"
        @click="handleSubmit"
      >
        提交退款申请
      </GlButton>
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  show: boolean
  selectedReason: string
  customReason?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  customReason: '',
  loading: false,
})

const emit = defineEmits<{
  close: []
  submit: [reason: string]
  'update:selectedReason': [value: string]
  'update:customReason': [value: string]
}>()

const presetReasons = [
  '陪玩师未上线',
  '服务态度差',
  '陪玩师技术不符',
  '临时有事无法继续',
  '其他',
]

const canSubmit = computed(() => {
  if (!props.selectedReason) return false
  if (props.selectedReason === '其他' && !props.customReason.trim()) return false
  return true
})

const selectReason = (reason: string) => {
  emit('update:selectedReason', reason)
}

const handleSubmit = () => {
  if (!canSubmit.value) return
  const finalReason = props.selectedReason === '其他'
    ? props.customReason.trim()
    : props.selectedReason
  emit('submit', finalReason)
}
</script>

<style lang="scss" scoped>
.refund-modal {
  padding: var(--spacing-md);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--spacing-md);
}

.modal-title {
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
}

.section-label {
  display: block;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-sm);
}

.reason-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-xs);
}

.reason-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  transition: all 0.2s;

  &--active {
    border-color: var(--color-primary);
    background: var(--color-primary-tint);
  }

  &:active {
    opacity: 0.8;
  }
}

.reason-text {
  font-size: var(--font-sm);
  color: var(--color-text);
}

.custom-section {
  margin-top: var(--spacing-sm);
}

.custom-input {
  :deep(.gl-input__textarea) {
    min-height: 160rpx;
  }
}

.refund-tip {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin: var(--spacing-md) 0;
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
}

.tip-text {
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
}
</style>
