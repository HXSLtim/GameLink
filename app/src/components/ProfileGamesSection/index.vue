<template>
  <FormSection title="游戏信息">
    <FormItem
      label="常玩游戏"
      :display-value="displayValue"
      placeholder="请选择"
      @click="$emit('edit')"
    />
  </FormSection>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import FormSection from '@/components/FormSection/index.vue'
import FormItem from '@/components/FormItem/index.vue'

interface Props {
  gamesCount?: number
}

const props = defineProps<Props>()

defineEmits<{
  edit: []
}>()

const displayValue = computed(() => {
  return props.gamesCount ? `${props.gamesCount}款` : undefined
})
</script>
