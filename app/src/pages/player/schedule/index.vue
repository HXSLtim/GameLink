<template>
  <BasePageLayout
    class="schedule-page"
    padding="0"
    title="排班设置"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <NavBar title="排班设置" @back="goBack" />
    </template>

    <view v-if="loading" class="loading-wrap">
      <Skeleton :rows="6" />
    </view>

    <template v-else>
      <SectionCard title="工作设置">
        <view class="setting-row">
          <text class="setting-label">空闲自动下线</text>
          <GlSwitch
            :model-value="autoOffline"
            size="small"
            @update:modelValue="(value) => (autoOffline = value)"
          />
        </view>
        <view class="setting-row">
          <text class="setting-label">时区</text>
          <text class="setting-value">{{ timezone }}</text>
        </view>
      </SectionCard>

      <SectionCard title="每周排班">
        <view class="slot-list">
          <view v-for="(slot, index) in slots" :key="slot.dayOfWeek" class="slot-item">
            <view class="slot-head">
              <text class="slot-day">{{ getDayLabel(slot.dayOfWeek) }}</text>
              <GlSwitch
                :model-value="slot.isAvailable"
                size="small"
                @update:modelValue="(value) => updateSlotAvailability(index, value)"
              />
            </view>

            <view v-if="slot.isAvailable" class="slot-time">
              <picker mode="time" :value="slot.startTime" @change="onStartTimeChange(index, $event)">
                <view class="time-pill">{{ slot.startTime }}</view>
              </picker>
              <text class="time-sep">-</text>
              <picker mode="time" :value="slot.endTime" @change="onEndTimeChange(index, $event)">
                <view class="time-pill">{{ slot.endTime }}</view>
              </picker>
            </view>
          </view>
        </view>
      </SectionCard>

      <view class="save-wrap">
        <GlButton type="primary" block :loading="saving" @click="saveSchedule">保存排班</GlButton>
      </view>
      <view class="bottom-placeholder"></view>
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import SectionCard from '@/components/SectionCard/index.vue'
import GlSwitch from '@/components/gl/Switch/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import Skeleton from '@/components/Skeleton/index.vue'
import { usePlayerSchedule } from '@/composables/usePlayerSchedule'

const {
  loading,
  saving,
  autoOffline,
  timezone,
  slots,
  getDayLabel,
  loadSchedule,
  updateSlotAvailability,
  updateSlotTime,
  saveSchedule,
  goBack,
} = usePlayerSchedule()

const onTimeChange = (
  index: number,
  key: 'startTime' | 'endTime',
  event: { detail?: { value?: string } }
) => {
  const value = event?.detail?.value
  if (!value) return
  updateSlotTime(index, key, value)
}

const onStartTimeChange = (index: number, event: { detail?: { value?: string } }) => {
  onTimeChange(index, 'startTime', event)
}

const onEndTimeChange = (index: number, event: { detail?: { value?: string } }) => {
  onTimeChange(index, 'endTime', event)
}

onMounted(loadSchedule)
onShow(loadSchedule)
</script>

<style scoped lang="scss">
.loading-wrap {
  padding: var(--spacing-md);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-xs) 0;
}

.setting-label {
  font-size: var(--font-sm);
  color: var(--color-text);
}

.setting-value {
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
}

.slot-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.slot-item {
  padding-bottom: var(--spacing-sm);
  border-bottom: 1rpx solid var(--color-border);

  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.slot-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.slot-day {
  font-size: var(--font-sm);
  color: var(--color-text);
  font-weight: 500;
}

.slot-time {
  margin-top: var(--spacing-xs);
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.time-pill {
  min-width: 120rpx;
  padding: 10rpx 16rpx;
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  background: var(--color-bg-secondary);
  text-align: center;
  font-size: var(--font-xs);
  color: var(--color-text);
}

.time-sep {
  color: var(--color-text-secondary);
}

.save-wrap {
  margin: var(--spacing-sm) var(--spacing-md) 0;
}

.bottom-placeholder {
  height: 120rpx;
}
</style>
