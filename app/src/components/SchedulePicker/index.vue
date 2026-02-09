<template>
  <GlCard title="预约时间" required :shadow="false" bordered>
    <view class="schedule-row">
      <view class="schedule-item" @tap="showDatePicker = true">
        <text class="schedule-label">日期</text>
        <text class="schedule-value" :class="{ placeholder: !date }">
          {{ date || '选择日期' }}
        </text>
        <uv-icon name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
      <view class="schedule-item" @tap="showTimePicker = true">
        <text class="schedule-label">时间</text>
        <text class="schedule-value" :class="{ placeholder: !time }">
          {{ time || '选择时间' }}
        </text>
        <uv-icon name="arrow-right" size="14" color="var(--color-text-secondary)"></uv-icon>
      </view>
    </view>
  </GlCard>

  <!-- 日期选择器 -->
  <uv-picker
    :show="showDatePicker"
    :columns="dateColumns"
    @confirm="onDateConfirm"
    @cancel="showDatePicker = false"
  ></uv-picker>

  <!-- 时间选择器 -->
  <uv-picker
    :show="showTimePicker"
    :columns="timeColumns"
    @confirm="onTimeConfirm"
    @cancel="showTimePicker = false"
  ></uv-picker>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import GlCard from '@/components/gl/Card/index.vue'

interface Props {
  date?: string
  time?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:date': [value: string]
  'update:time': [value: string]
}>()

const showDatePicker = ref(false)
const showTimePicker = ref(false)

// 生成日期选项（未来7天）
const dateColumns = computed(() => {
  const dates: string[] = []
  const today = new Date()
  for (let i = 0; i < 7; i++) {
    const date = new Date(today)
    date.setDate(today.getDate() + i)
    const month = date.getMonth() + 1
    const day = date.getDate()
    const weekDay = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'][date.getDay()]
    dates.push(i === 0 ? '今天' : i === 1 ? '明天' : `${month}月${day}日 ${weekDay}`)
  }
  return [dates]
})

// 生成时间选项
const timeColumns = computed(() => {
  const times: string[] = []
  for (let h = 8; h <= 23; h++) {
    times.push(`${String(h).padStart(2, '0')}:00`)
    times.push(`${String(h).padStart(2, '0')}:30`)
  }
  return [times]
})

const onDateConfirm = (e: any) => {
  emit('update:date', e.value[0])
  showDatePicker.value = false
}

const onTimeConfirm = (e: any) => {
  emit('update:time', e.value[0])
  showTimePicker.value = false
}
</script>

<style lang="scss" scoped>
.schedule-row {
  display: flex;
  gap: var(--spacing-md);
}

.schedule-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  cursor: pointer;
  @include press-effect;
}

.schedule-label {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.schedule-value {
  flex: 1;
  font-size: var(--font-sm);
  color: var(--color-text);
  text-align: right;

  &.placeholder {
    color: var(--color-text-placeholder);
  }
}
</style>
