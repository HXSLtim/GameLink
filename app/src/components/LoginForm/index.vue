<template>
  <view class="form-section">
    <view class="input-group">
      <text class="input-label">手机号/邮箱</text>
      <input
        :value="username"
        class="input"
        type="text"
        placeholder="请输入手机号或邮箱"
        @input="(e: any) => $emit('update:username', e.detail.value)"
      />
    </view>

    <view class="input-group">
      <text class="input-label">密码</text>
      <input
        :value="password"
        class="input"
        type="password"
        placeholder="请输入密码"
        @input="(e: any) => $emit('update:password', e.detail.value)"
      />
    </view>

    <GlButton 
      type="primary"
      block
      round
      size="large"
      :loading="loading"
      :disabled="!canSubmit"
      @click="$emit('submit')"
    >
      登录
    </GlButton>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  username: string
  password: string
  loading?: boolean
}

const props = defineProps<Props>()

defineEmits<{
  'update:username': [value: string]
  'update:password': [value: string]
  submit: []
}>()

const canSubmit = computed(() => props.username && props.password)
</script>

<style lang="scss" scoped>
.form-section {
  padding: 0 48rpx;
}

.input-group {
  margin-bottom: 32rpx;
}

.input-label {
  display: block;
  font-size: 28rpx;
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: 16rpx;
}

.input {
  width: 100%;
  height: 88rpx;
  padding: 0 24rpx;
  background: var(--color-bg-secondary);
  border: 2rpx solid var(--color-border);
  border-radius: 16rpx;
  font-size: 30rpx;
  color: var(--color-text);
}
</style>
