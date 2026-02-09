<template>
  <GlCard :title="title" :shadow="false" bordered>
    <GlInput
      class="intro-input"
      :model-value="modelValue"
      type="textarea"
      :placeholder="placeholder"
      :maxlength="maxLength"
      @update:modelValue="(value) => emit('update:modelValue', value)"
    />
    <text class="char-count">{{ (modelValue || '').length }}/{{ maxLength }}</text>
  </GlCard>
</template>

<script setup lang="ts">
import GlCard from '@/components/gl/Card/index.vue'
import GlInput from '@/components/gl/Input/index.vue'

interface Props {
  modelValue: string
  title?: string
  placeholder?: string
  maxLength?: number
}

const props = withDefaults(defineProps<Props>(), {
  title: '个人介绍',
  placeholder: '介绍一下自己的陪玩特色和优势吧~',
  maxLength: 500,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style lang="scss" scoped>
.intro-input {
  :deep(.gl-input__textarea) {
    min-height: 200rpx;
  }
}

.char-count {
  display: block;
  text-align: right;
  font-size: var(--font-xs);
  color: var(--color-text-placeholder);
  margin-top: var(--spacing-xs);
}
</style>
