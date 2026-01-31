<template>
  <view class="balance-card">
    <view class="balance-content">
      <text class="balance-label">账户余额（元）</text>
      <view class="balance-row">
        <text class="balance-value">{{ showBalance ? formattedBalance : '****' }}</text>
        <view class="eye-btn" @tap="$emit('toggle-visibility')">
          <uv-icon :name="showBalance ? 'eye' : 'eye-off'" size="20" color="#fff"></uv-icon>
        </view>
      </view>
      
      <view class="balance-stats">
        <view class="stat-item">
          <text class="stat-label">累计充值</text>
          <text class="stat-value">{{ showBalance ? formattedRecharge : '****' }}</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item">
          <text class="stat-label">累计消费</text>
          <text class="stat-value">{{ showBalance ? formattedSpent : '****' }}</text>
        </view>
      </view>
      
      <view class="balance-actions">
        <GlButton type="default" size="small" round plain custom-style="background: rgba(255,255,255,0.2); border-color: rgba(255,255,255,0.3); color: #fff;" @click="$emit('recharge')">
          <uv-icon name="plus-circle" size="18" color="#fff"></uv-icon>
          <text style="margin-left: 8rpx;">充值</text>
        </GlButton>
        <GlButton type="default" size="small" round plain custom-style="background: rgba(255,255,255,0.2); border-color: rgba(255,255,255,0.3); color: #fff;" @click="$emit('withdraw')">
          <uv-icon name="red-packet" size="18" color="#fff"></uv-icon>
          <text style="margin-left: 8rpx;">提现</text>
        </GlButton>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'

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

const formatMoney = (cents: number) => (cents / 100).toFixed(2)

const formattedBalance = computed(() => formatMoney(props.balance))
const formattedRecharge = computed(() => formatMoney(props.totalRecharge))
const formattedSpent = computed(() => formatMoney(props.totalSpent))
</script>

<style lang="scss" scoped>
.balance-card {
  margin: 24rpx;
  padding: 40rpx 32rpx;
  border-radius: 24rpx;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-primary-light, #4ADE80) 100%);
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -30%;
    width: 400rpx;
    height: 400rpx;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.15) 0%, transparent 70%);
    border-radius: 50%;
  }
}

.balance-content {
  position: relative;
  z-index: 1;
}

.balance-label {
  display: block;
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.85);
  margin-bottom: 16rpx;
}

.balance-row {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 32rpx;
}

.balance-value {
  font-size: 64rpx;
  font-weight: 800;
  color: #FFFFFF;
}

.eye-btn {
  padding: 8rpx;
  opacity: 0.8;
}

.balance-stats {
  display: flex;
  padding: 24rpx 0;
  border-top: 1rpx solid rgba(255, 255, 255, 0.2);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.2);
  margin-bottom: 32rpx;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 8rpx;
}

.stat-value {
  font-size: 32rpx;
  font-weight: 600;
  color: #FFFFFF;
}

.stat-divider {
  width: 1rpx;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 32rpx;
}

.balance-actions {
  display: flex;
  gap: 24rpx;
}
</style>
