<template>
  <view class="form-section">
    <view class="input-group">
      <text class="input-label">手机号/邮箱</text>
      <GlInput
        :model-value="account"
        type="text"
        placeholder="请输入手机号或邮箱"
        size="medium"
        clearable
        @update:modelValue="(value) => $emit('update:account', value)"
      />
    </view>

    <view class="input-group">
      <text class="input-label">密码</text>
      <GlInput
        :model-value="password"
        :type="showPassword ? 'text' : 'password'"
        placeholder="请输入密码"
        size="medium"
        :suffix-icon="showPassword ? 'eye' : 'eye-off'"
        :suffix-clickable="true"
        @suffix-click="togglePassword"
        @update:modelValue="(value) => $emit('update:password', value)"
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
import { computed, ref } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  account: string
  password: string
  loading?: boolean
}

const props = defineProps<Props>()

defineEmits<{
  'update:account': [value: string]
  'update:password': [value: string]
  submit: []
}>()

const canSubmit = computed(() => props.account && props.password)

const showPassword = ref(false)
const togglePassword = () => {
  showPassword.value = !showPassword.value
}
</script>

<style lang="scss" scoped>
.form-section {
  padding: 0 var(--spacing-lg);
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

</style>
