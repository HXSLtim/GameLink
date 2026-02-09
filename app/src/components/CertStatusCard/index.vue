<template>
  <StatusCard :status-class="status" border-width="2rpx" :border-color="borderColor">
    <StatusInfo :icon="icon" :title="title" :description="description" size="md" />
  </StatusCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import StatusCard from '@/components/StatusCard/index.vue'
import StatusInfo from '@/components/StatusInfo/index.vue'
import { getCertStatusPreset, type CertStatus } from '@/components/StatusCard/presets'

interface Props {
  status: CertStatus
}

const props = defineProps<Props>()

const config = computed(() => getCertStatusPreset(props.status))

const icon = computed(() => config.value.icon)
const title = computed(() => config.value.title)
const description = computed(() => config.value.desc)
const borderColor = computed(() => config.value.borderColor)
</script>

<style lang="scss" scoped>
</style>
