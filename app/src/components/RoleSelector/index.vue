<template>
  <view class="role-section">
    <text class="section-title">{{ title }}</text>
    <view class="role-cards">
      <view 
        v-for="role in roles"
        :key="role.value"
        class="role-card"
        :class="{ active: modelValue === role.value }"
        @click="$emit('update:modelValue', role.value)"
      >
        <view class="role-icon">{{ role.icon }}</view>
        <text class="role-name">{{ role.name }}</text>
        <text class="role-desc">{{ role.desc }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
export interface RoleOption {
  value: string
  icon: string
  name: string
  desc: string
}

interface Props {
  title?: string
  roles: RoleOption[]
  modelValue?: string
}

withDefaults(defineProps<Props>(), {
  title: '选择注册身份',
})

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style lang="scss" scoped>
.role-section {
  margin-bottom: 32rpx;
}

.section-title {
  display: block;
  font-size: 28rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 20rpx;
}

.role-cards {
  display: flex;
  gap: 24rpx;
}

.role-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32rpx 16rpx;
  background: var(--color-bg-secondary);
  border: 2rpx solid var(--color-border);
  border-radius: 20rpx;
  transition: all 0.2s;
  
  &.active {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.1);
  }
}

.role-icon {
  width: 80rpx;
  height: 80rpx;
  background: var(--color-bg);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: 16rpx;
}

.role-name {
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 8rpx;
}

.role-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
  text-align: center;
}
</style>
