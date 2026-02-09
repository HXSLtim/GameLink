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
        <view class="role-icon"><uv-icon :name="role.icon" size="28" color="var(--color-text-secondary)" /></view>
        <text class="role-name">{{ role.name }}</text>
        <text class="role-desc">{{ role.desc }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { RoleOption } from '@/types/auth'

interface Props {
  title?: string
  roles: RoleOption[]
  modelValue?: RoleOption['value']
}

withDefaults(defineProps<Props>(), {
  title: '选择注册身份',
})

defineEmits<{
  'update:modelValue': [value: RoleOption['value']]
}>()
</script>

<style lang="scss" scoped>
.role-section {
  margin-bottom: var(--spacing-lg);
}

.section-title {
  display: block;
  font-size: var(--font-md);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-sm);
}

.role-cards {
  display: flex;
  gap: var(--spacing-md);
}

.role-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-sm);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all 0.2s;
  cursor: pointer;
  @include press-effect;
  
  &.active {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
}

.role-icon {
  width: 64rpx;
  height: 64rpx;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--spacing-sm);
}

.role-name {
  font-size: var(--font-base);
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: var(--spacing-xs);
}

.role-desc {
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  text-align: center;
}
</style>
