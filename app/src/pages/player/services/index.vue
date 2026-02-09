<template>
  <BasePageLayout
    class="services-page"
    padding="0"
    title="服务管理"
    :show-back="true"
    :show-tab-bar="true"
    :show-mobile-tab-bar="false"
  >
    <template #nav>
      <!-- 顶部导航 -->
      <NavBar title="服务管理" @back="goBack">
        <template #right>
          <GlButton type="primary" size="mini" @click="addService">添加</GlButton>
        </template>
      </NavBar>
    </template>

    <!-- 服务统计 -->
    <ServiceStatsBar :items="statItems" />

    <!-- 服务列表 -->
    <ServiceList
      :loading="loading"
      :services="services"
      @toggle-status="toggleStatus"
      @edit="editService"
      @delete="handleDelete"
      @add="addService"
    />

    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- 添加/编辑弹窗 -->
      <ServiceEditorPanel
        :show="showEditor"
        :editing="!!editingService"
        :saving="saving"
        :form="form"
        @close="closeEditor"
        @save="saveService"
      />
    </template>
  </BasePageLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import BasePageLayout from '@/components/layout/BasePageLayout/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
// Business 组件
import ServiceStatsBar from '@/components/ServiceStatsBar/index.vue'
import ServiceList from '@/components/ServiceList/index.vue'
import ServiceEditorPanel from '@/components/ServiceEditorPanel/index.vue'
// Composables
import { usePlayerServices } from '@/composables/usePlayerServices'

const {
  loading,
  showEditor,
  saving,
  editingService,
  services,
  form,
  statItems,
  loadServices,
  addService,
  editService,
  closeEditor,
  saveService,
  toggleStatus,
  handleDelete,
  goBack,
} = usePlayerServices()

onMounted(loadServices)
onShow(loadServices)
</script>

<style lang="scss" scoped>
.bottom-placeholder {
  height: 100rpx;
}
</style>
