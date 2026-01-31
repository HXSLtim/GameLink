<template>
  <view class="services-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="服务管理" @back="goBack">
      <template #right>
        <GlButton type="primary" size="mini" @click="addService">添加</GlButton>
      </template>
    </NavBar>

    <scroll-view class="content-scroll" scroll-y>
      <!-- 服务统计 -->
      <view class="stats-bar">
        <view v-for="item in statItems" :key="item.label" class="stat-item">
          <text class="stat-value" :class="{ highlight: item.highlight }">{{ item.value }}</text>
          <text class="stat-label">{{ item.label }}</text>
        </view>
      </view>

      <!-- 服务列表 -->
      <view class="services-list">
        <template v-if="loading">
          <uv-skeleton rows="3" title loading v-for="i in 2" :key="i"></uv-skeleton>
        </template>
        
        <template v-else-if="services.length > 0">
          <ServiceCard
            v-for="service in services"
            :key="service.id"
            :service="service"
            @toggle-status="toggleStatus(service)"
            @edit="editService(service)"
            @delete="handleDelete(service)"
          />
        </template>
        
        <GlEmpty 
          v-else
          title="暂无服务"
          description="添加服务开始接单吧"
          action-text="添加服务"
          @action="addService"
        />
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 添加/编辑弹窗 -->
    <uv-popup :show="showEditor" mode="bottom" round="20" @close="closeEditor">
      <view class="editor-panel">
        <view class="editor-header">
          <text class="editor-title">{{ editingService ? '编辑服务' : '添加服务' }}</text>
          <uv-icon name="close" size="20" @click="closeEditor"></uv-icon>
        </view>
        
        <scroll-view class="editor-content" scroll-y>
          <FormItem label="游戏" :display-value="form.gameName" required @click="() => {}" />
          <FormItem label="服务类型" :display-value="form.serviceName" required @click="() => {}" />
          <FormItem label="段位" :display-value="form.rankName" required @click="() => {}" />
          
          <view class="form-row">
            <text class="form-label">价格 <text class="required">*</text></text>
            <view class="price-input-wrap">
              <text class="currency">¥</text>
              <input v-model.number="form.price" type="digit" class="price-input" placeholder="0.00" />
              <text class="unit">/{{ form.unit }}</text>
            </view>
          </view>
          
          <view class="form-row">
            <text class="form-label">服务介绍</text>
            <textarea 
              v-model="form.description"
              class="desc-textarea"
              placeholder="介绍您的服务特点..."
              :maxlength="200"
            />
          </view>
        </scroll-view>
        
        <view class="editor-footer">
          <GlButton type="primary" block round :loading="saving" @click="saveService">
            {{ editingService ? '保存修改' : '添加服务' }}
          </GlButton>
        </view>
      </view>
    </uv-popup>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import GlButton from '@/components/gl/Button/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import FormItem from '@/components/FormItem/index.vue'
// Business 组件
import ServiceCard from '@/components/ServiceCard/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
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
.services-page {
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  box-sizing: border-box;
  
  @include desktop {
    height: 100vh;
    min-height: auto;
    overflow: hidden;
  }
}

.content-scroll {
  flex: 1;
  overflow-y: auto;
}

.stats-bar {
  display: flex;
  padding: 24rpx;
  margin: 24rpx;
  background: var(--color-bg-card);
  border-radius: 20rpx;
  border: 2rpx solid var(--color-border);
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-value {
  font-size: 40rpx;
  font-weight: 700;
  color: var(--color-text);
  
  &.highlight {
    color: var(--color-primary);
  }
}

.stat-label {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.services-list {
  padding: 0 24rpx;
}

.editor-panel {
  padding: 24rpx;
  padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;
}

.editor-title {
  font-size: 34rpx;
  font-weight: 700;
  color: var(--color-text);
}

.editor-content {
  flex: 1;
  overflow-y: auto;
}

.form-row {
  margin-bottom: 24rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: var(--color-text);
  margin-bottom: 12rpx;
  
  .required {
    color: var(--color-error);
  }
}

.price-input-wrap {
  display: flex;
  align-items: center;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  border: 2rpx solid var(--color-border);
}

.currency {
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-right: 8rpx;
}

.price-input {
  flex: 1;
  font-size: 32rpx;
  font-weight: 600;
  color: var(--color-text);
}

.unit {
  font-size: 26rpx;
  color: var(--color-text-secondary);
}

.desc-textarea {
  width: 100%;
  height: 160rpx;
  padding: 20rpx;
  background: var(--color-bg-secondary);
  border-radius: 12rpx;
  border: 2rpx solid var(--color-border);
  font-size: 28rpx;
  color: var(--color-text);
}

.editor-footer {
  padding-top: 24rpx;
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
