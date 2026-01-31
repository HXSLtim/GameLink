<template>
  <view class="agreement-page page-container">
    <!-- 顶部导航 -->
    <NavBar :title="pageTitle" @back="goBack" />

    <!-- 内容区域 -->
    <scroll-view class="content-scroll" scroll-y>
      <AgreementContent :content="agreementContent" />
      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
// Business 组件
import AgreementContent from '@/components/AgreementContent/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// 数据
import { agreements, agreementTitles, type AgreementType } from '@/data/agreements'

const type = ref<AgreementType>('user')

const pageTitle = computed(() => agreementTitles[type.value])

const agreementContent = computed(() => agreements[type.value])

const goBack = () => uni.navigateBack()

onLoad((options) => {
  if (options?.type && options.type in agreements) {
    type.value = options.type as AgreementType
  }
})
</script>

<style lang="scss" scoped>
.agreement-page {
  height: 100vh;
  height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-card);
  box-sizing: border-box;
  overflow: hidden;
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

.bottom-placeholder {
  height: 60rpx;
}
</style>
