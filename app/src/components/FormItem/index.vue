<template>
  <view class="form-item" :class="{ clickable }" @tap="handleClick">
    <text class="form-label">{{ label }}</text>
    
    <!-- 输入框 -->
    <template v-if="type === 'input'">
      <input 
        :value="modelValue"
        class="form-input"
        :placeholder="placeholder"
        :maxlength="maxlength"
        @input="(e: any) => $emit('update:modelValue', e.detail.value)"
      />
    </template>
    
    <!-- 选择器/只读 -->
    <template v-else>
      <view class="form-value-wrap">
        <text class="form-value" :class="{ placeholder: !displayValue }">
          {{ displayValue || placeholder }}
        </text>
        <slot name="extra"></slot>
        <uv-icon v-if="clickable && !$slots.extra" name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  label: string
  type?: 'input' | 'picker' | 'readonly'
  modelValue?: string | number
  displayValue?: string
  placeholder?: string
  maxlength?: number
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'picker',
  placeholder: '请选择',
  clickable: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  click: []
}>()

const handleClick = () => {
  if (props.type !== 'input' && props.clickable) {
    emit('click')
  }
}
</script>

<style lang="scss" scoped>
.form-item {
  display: flex;
  align-items: center;
  padding: 28rpx 0;
  border-bottom: 1rpx solid var(--color-border);
  
  &:last-child {
    border-bottom: none;
  }
  
  &.clickable {
    cursor: pointer;
  }
}

.form-label {
  width: 160rpx;
  font-size: 30rpx;
  color: var(--color-text);
  flex-shrink: 0;
}

.form-input {
  flex: 1;
  font-size: 30rpx;
  color: var(--color-text);
  text-align: right;
  background: transparent;
}

.form-value-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8rpx;
}

.form-value {
  font-size: 30rpx;
  color: var(--color-text);
  
  &.placeholder {
    color: var(--color-text-placeholder);
  }
}
</style>
