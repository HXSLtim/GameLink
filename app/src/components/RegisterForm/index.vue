<template>
  <view class="form-section">
    <view class="input-group">
      <text class="input-label">手机号</text>
      <GlInput
        v-model="form.phone"
        type="number"
        maxlength="11"
        placeholder="请输入手机号"
        size="medium"
        clearable
      />
    </view>

    <view class="input-group">
      <text class="input-label">昵称</text>
      <GlInput
        v-model="form.nickname"
        type="text"
        maxlength="20"
        placeholder="请输入昵称"
        size="medium"
        clearable
      />
    </view>

    <view class="input-group">
      <text class="input-label">密码</text>
      <GlInput
        v-model="form.password"
        type="password"
        placeholder="请设置密码（6-20位）"
        size="medium"
      />
    </view>

    <view class="input-group">
      <text class="input-label">确认密码</text>
      <GlInput
        v-model="form.confirmPassword"
        type="password"
        placeholder="请再次输入密码"
        size="medium"
      />
    </view>

    <!-- 协议勾选 -->
    <view class="agreement" @click="toggleAgreed">
      <view class="checkbox" :class="{ checked: agreed }">
        <uv-icon v-if="agreed" name="checkbox-mark" size="14" color="#fff"></uv-icon>
      </view>
      <text class="agreement-text">
        我已阅读并同意
        <text class="link-text" @tap.stop="$emit('agreement', 'user')">《用户协议》</text>
        和
        <text class="link-text" @tap.stop="$emit('agreement', 'privacy')">《隐私政策》</text>
      </text>
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
      注册
    </GlButton>

    <view class="login-link" @click="$emit('go-login')">
      <text>已有账号？去登录</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import GlButton from '@/components/gl/Button/index.vue'
import GlInput from '@/components/gl/Input/index.vue'
import type { RegisterFormData } from '@/types/auth'
import type { AgreementType } from '@/types/agreement'

type AgreementLinkType = Extract<AgreementType, 'user' | 'privacy'>

interface Props {
  form: RegisterFormData
  agreed: boolean
  loading: boolean
  canSubmit: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:agreed': [value: boolean]
  submit: []
  'go-login': []
  agreement: [type: AgreementLinkType]
}>()

const toggleAgreed = () => {
  emit('update:agreed', !props.agreed)
}
</script>

<style lang="scss" scoped>
.form-section {
  padding-top: var(--spacing-sm);
}

.input-group {
  margin-bottom: var(--spacing-md);
}

.input-label {
  display: block;
  font-size: var(--font-sm);
  font-weight: 500;
  color: var(--color-text);
  margin-bottom: var(--spacing-sm);
}


.agreement {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-md);
  cursor: pointer;
  @include press-effect;
}

.checkbox {
  width: 32rpx;
  height: 32rpx;
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  &.checked {
    background: var(--color-primary);
    border-color: var(--color-primary);
  }
}

.agreement-text {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
}

.link-text {
  color: var(--color-primary);
}

.login-link {
  margin-top: var(--spacing-md);
  text-align: center;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  @include press-effect;
}
</style>
