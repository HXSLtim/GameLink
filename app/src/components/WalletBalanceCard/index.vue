<template>
  <SectionCard class="balance-card" margin="var(--spacing-md)" padding="var(--spacing-md)">
    <view class="balance-content">
      <text class="balance-label">账户余额（元）</text>
      <view class="balance-row">
        <PriceTag
          v-if="showBalance"
          class="balance-value"
          :amount="balance"
          amount-unit="cents"
          size="large"
          :show-currency="false"
        />
        <text v-else class="balance-value">****</text>
        <view class="eye-btn" @tap="$emit('toggle-visibility')">
          <uv-icon :name="showBalance ? 'eye' : 'eye-off'" size="20" color="var(--color-text-secondary)"></uv-icon>
        </view>
      </view>
      
      <view class="balance-stats">
        <HeaderStatsRow :items="stats" size="md" item-padding="0">
          <template #value="{ item }">
            <PriceTag
              v-if="showBalance"
              :amount="item.value"
              amount-unit="cents"
              size="small"
              :show-currency="false"
            />
            <text v-else>****</text>
          </template>
        </HeaderStatsRow>
      </view>
      
      <view class="balance-actions">
        <GlButton class="balance-action-btn" type="default" size="small" round plain @click="$emit('recharge')">
          <uv-icon name="plus-circle" size="18" color="var(--color-text)"></uv-icon>
          <text class="action-text">充值</text>
        </GlButton>
        <GlButton class="balance-action-btn" type="default" size="small" round plain @click="$emit('withdraw')">
          <uv-icon name="red-packet" size="18" color="var(--color-text)"></uv-icon>
          <text class="action-text">提现</text>
        </GlButton>
      </view>
    </view>
  </SectionCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SectionCard from '@/components/SectionCard/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import PriceTag from '@/components/PriceTag/index.vue'
import HeaderStatsRow from '@/components/HeaderStatsRow/index.vue'
import type { HeaderStatItem } from '@/types/ui'

interface Props {
  balance: number // 分
  totalRecharge: number // 分
  totalSpent: number // 分
  showBalance: boolean
}

const props = defineProps<Props>()

defineEmits<{
  'toggle-visibility': []
  recharge: []
  withdraw: []
}>()

const stats = computed<HeaderStatItem[]>(() => [
  { key: 'recharge', label: '累计充值', value: props.totalRecharge },
  { key: 'spent', label: '累计消费', value: props.totalSpent },
])

</script>

<style lang="scss" scoped>
.balance-content {
  position: relative;
  z-index: 1;
}

.balance-label {
  display: block;
  font-size: var(--font-sm);
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-sm);
}

.balance-row {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.balance-value {
  font-size: var(--font-xl);
  font-weight: 700;
  color: var(--color-text);
}

.balance-value :deep(.amount) {
  font-weight: 700;
  font-size: var(--font-xl);
}

.eye-btn {
  padding: var(--spacing-xs);
  opacity: 0.8;
  @include press-effect;
}

.balance-stats {
  padding: var(--spacing-md) 0;
  border-top: 1rpx solid var(--color-border);
  border-bottom: 1rpx solid var(--color-border);
  margin-bottom: var(--spacing-lg);
}

.balance-actions {
  display: flex;
  gap: var(--spacing-md);
}

.balance-action-btn {
  background: var(--color-bg-secondary) !important;
  border-color: var(--color-border) !important;
  color: var(--color-text) !important;
}

.action-text {
  margin-left: var(--spacing-xs);
}
</style>
