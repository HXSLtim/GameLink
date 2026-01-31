<template>
  <view class="help-page page-container">
    <!-- 顶部导航 -->
    <NavBar title="帮助中心" @back="goBack" />

    <!-- 搜索栏 -->
    <SearchBar
      v-model="searchKeyword"
      placeholder="搜索问题"
      @search="handleSearch"
    />

    <scroll-view class="content-scroll" scroll-y>
      <!-- 常见问题分类 -->
      <view class="category-section">
        <view class="category-grid">
          <view
            v-for="cat in categories"
            :key="cat.id"
            class="category-item"
            :class="{ active: selectedCategory === cat.id }"
            @tap="selectCategory(cat.id)"
          >
            <view class="category-icon">{{ cat.icon }}</view>
            <text class="category-name">{{ cat.name }}</text>
          </view>
        </view>
      </view>

      <!-- FAQ 列表 -->
      <view class="faq-section">
        <SectionHeader title="常见问题" :show-more="false" />
        
        <view class="faq-list">
          <FaqItem
            v-for="faq in displayFaqs"
            :key="faq.id"
            :question="faq.question"
            :answer="faq.answer"
            :expanded="expandedId === faq.id"
            @toggle="toggleFaq(faq.id)"
          />
          
          <GlEmpty v-if="displayFaqs.length === 0" title="未找到相关问题" compact />
        </view>
      </view>

      <!-- 联系客服 -->
      <view class="contact-section">
        <SectionHeader title="没有找到答案？" :show-more="false" />
        <GlCard :shadow="false" bordered @click="goToService">
          <view class="contact-card">
            <view class="contact-icon">💬</view>
            <view class="contact-info">
              <text class="contact-title">在线客服</text>
              <text class="contact-desc">工作时间 9:00-22:00</text>
            </view>
            <uv-icon name="arrow-right" size="16" color="var(--color-text-secondary)"></uv-icon>
          </view>
        </GlCard>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- PC 端侧边栏 -->
    <CustomTabBar :show-mobile-tab-bar="false" />
  </view>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
import SectionHeader from '@/components/SectionHeader/index.vue'
import GlCard from '@/components/gl/Card/index.vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
// Business 组件
import FaqItem from '@/components/FaqItem/index.vue'
import CustomTabBar from '@/components/CustomTabBar/index.vue'
// Composables
import { useHelp } from '@/composables/useHelp'

const {
  searchKeyword,
  selectedCategory,
  expandedId,
  categories,
  displayFaqs,
  selectCategory,
  toggleFaq,
  handleSearch,
  goBack,
  goToService,
} = useHelp()
</script>

<style lang="scss" scoped>
.help-page {
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

.category-section {
  padding: 24rpx;
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20rpx;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  padding: 28rpx 16rpx;
  background: var(--color-bg-card);
  border-radius: 16rpx;
  border: 2rpx solid var(--color-border);
  transition: all 0.2s;
  
  &:active,
  &.active {
    border-color: var(--color-primary);
    background: rgba(0, 210, 106, 0.05);
  }
}

.category-icon {
  font-size: 40rpx;
}

.category-name {
  font-size: 26rpx;
  color: var(--color-text);
  font-weight: 500;
}

.faq-section {
  padding: 0 24rpx;
}

.faq-list {
  margin-top: 16rpx;
}

.contact-section {
  padding: 24rpx;
}

.contact-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.contact-icon {
  font-size: 48rpx;
}

.contact-info {
  flex: 1;
}

.contact-title {
  display: block;
  font-size: 30rpx;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 4rpx;
}

.contact-desc {
  font-size: 24rpx;
  color: var(--color-text-secondary);
}

.bottom-placeholder {
  height: 100rpx;
}
</style>
