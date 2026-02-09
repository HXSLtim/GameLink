<template>
  <FormSection title="基本信息">
    <template #extra>
      <text class="required-tip">* 为必填项</text>
    </template>

    <FormItem
      v-model="form.realName"
      label="真实姓名"
      type="input"
      placeholder="请输入真实姓名"
      required
    />

    <FormItem
      v-model="form.idNumber"
      label="身份证号"
      type="input"
      placeholder="请输入身份证号"
      required
    />

    <FormItem
      label="性别"
      :display-value="genderText"
      required
      @click="$emit('pick-gender')"
    />
  </FormSection>
</template>

<script setup lang="ts">
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'
import type { CertificationForm } from '@/types/certification'

type CertificationBasicForm = Pick<CertificationForm, 'realName' | 'idNumber' | 'gender'>

interface Props {
  form: CertificationBasicForm
  genderText: string
}

defineProps<Props>()

defineEmits<{
  'pick-gender': []
}>()
</script>

<style lang="scss" scoped>
.required-tip {
  font-size: var(--font-sm);
  color: var(--color-error);
}
</style>
