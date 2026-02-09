<template>
  <PageShell class="agreement-page" padding="0">
    <template #header>
      <!-- 顶部导航 -->
      <NavBar :title="pageTitle" @back="goBack" />
    </template>

    <!-- 内容区域 -->
    <AgreementContent :content="agreementContent" />
    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- PC 端侧边栏 -->
      <CustomTabBar :show-mobile-tab-bar="false" />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
// Business 组件
import AgreementContent from '@/components/AgreementContent/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// 数据
import { agreements, agreementTitles } from '@/data/agreements'
import type { AgreementType } from '@/types/agreement'

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
  background: var(--color-bg-card);
}

.bottom-placeholder {
  height: 60rpx;
}
</style>
