<template>
  <view class="other-login">
    <!-- #ifdef MP-WEIXIN -->
    <text class="link" @click="$emit('open-account')">使用账号密码登录</text>
    <!-- #endif -->
    <text class="link" @click="$emit('register')">还没有账号？立即注册</text>
  </view>
</template>

<script setup lang="ts">
defineEmits<{
  'open-account': []
  register: []
}>()
</script>

<style lang="scss" scoped>
.other-login {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-md);
}

.link {
  font-size: var(--font-md);
  color: var(--color-primary);
  cursor: pointer;
  @include press-effect;
}
</style>
