<template>
  <FormSection title="基本信息">
    <FormItem
      v-model="form.nickname"
      label="昵称"
      type="input"
      placeholder="请输入昵称"
      :maxlength="20"
    />
    <FormItem
      label="性别"
      :display-value="genderText"
      @click="$emit('pick-gender')"
    />
    <FormItem
      label="生日"
      :display-value="form.birthday"
      placeholder="请选择"
      @click="$emit('pick-birthday')"
    />
    <FormItem
      label="地区"
      :display-value="form.region"
      placeholder="请选择"
      @click="$emit('pick-region')"
    />
  </FormSection>
</template>

<script setup lang="ts">
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'
import type { ProfileForm } from '@/types/profile'

interface Props {
  form: ProfileForm
  genderText: string
}

defineProps<Props>()

defineEmits<{
  'pick-gender': []
  'pick-birthday': []
  'pick-region': []
}>()
</script>
