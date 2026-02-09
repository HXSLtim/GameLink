<template>
  <FormSection title="联系方式">
    <FormItem
      label="手机号"
      :display-value="phone"
      :clickable="false"
    >
      <template #extra>
        <text class="change-btn" @tap="$emit('change-phone')">更换</text>
      </template>
    </FormItem>
    <FormItem
      label="微信"
      :display-value="wechatBound ? '已绑定' : '未绑定'"
      :clickable="false"
    >
      <template #extra>
        <text v-if="!wechatBound" class="bind-btn" @tap="$emit('bind-wechat')">绑定</text>
      </template>
    </FormItem>
  </FormSection>
</template>

<script setup lang="ts">
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'
import type { ProfileContactInfo } from '@/types/profile'

type Props = ProfileContactInfo

defineProps<Props>()

defineEmits<{
  'change-phone': []
  'bind-wechat': []
}>()
</script>

<style lang="scss" scoped>
.change-btn,
.bind-btn {
  font-size: var(--font-sm);
  color: var(--color-primary);
  font-weight: 500;
}
</style>
