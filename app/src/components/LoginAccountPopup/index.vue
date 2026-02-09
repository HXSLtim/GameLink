<template>
  <uv-popup :show="show" mode="bottom" round="24" @close="$emit('update:show', false)">
    <view class="popup-content">
      <view class="popup-header">
        <text class="popup-title">账号登录</text>
        <uv-icon name="close" size="20" @click="$emit('update:show', false)"></uv-icon>
      </view>

      <LoginForm
        v-model:account="form.account"
        v-model:password="form.password"
        :loading="loading"
        @submit="$emit('submit')"
      />
    </view>
  </uv-popup>
</template>

<script setup lang="ts">
import LoginForm from '@/components/LoginForm/index.vue'
import type { LoginFormData } from '@/types/auth'

interface Props {
  show: boolean
  loading: boolean
  form: LoginFormData
}

defineProps<Props>()

defineEmits<{
  'update:show': [value: boolean]
  submit: []
}>()
</script>

<style lang="scss" scoped>
.popup-content {
  padding: 32rpx 0;
  padding-bottom: calc(32rpx + env(safe-area-inset-bottom));
}

.popup-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 48rpx 32rpx;
}

.popup-title {
  font-size: var(--font-md);
  font-weight: 700;
  color: var(--color-text);
}
</style>
