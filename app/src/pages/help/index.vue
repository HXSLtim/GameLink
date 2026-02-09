<template>
  <PageShell class="help-page" padding="0">
    <template #header>
      <ListHeaderStack>
        <template #nav>
          <!-- 顶部导航 -->
          <NavBar title="帮助中心" @back="goBack" />
        </template>

        <template #search>
          <!-- 搜索栏 -->
          <SearchBar
            v-model="searchKeyword"
            placeholder="搜索问题"
            @search="handleSearch"
          />
        </template>
      </ListHeaderStack>
    </template>

    <!-- 常见问题分类 -->
    <HelpCategoryGrid
      :categories="categories"
      :selected-id="selectedCategory"
      @select="selectCategory"
    />

    <!-- FAQ 列表 -->
    <HelpFaqSection
      :faqs="displayFaqs"
      :expanded-id="expandedId"
      @toggle="toggleFaq"
    />

    <!-- 联系客服 -->
    <HelpContactCard @click="goToService" />

    <view class="bottom-placeholder"></view>

    <template #footer>
      <!-- PC 端侧边栏 -->
      <CustomTabBar :show-mobile-tab-bar="false" />
    </template>
  </PageShell>
</template>

<script setup lang="ts">
// Pattern 组件
import NavBar from '@/components/NavBar/index.vue'
import PageShell from '@/components/layout/PageShell/index.vue'
import ListHeaderStack from '@/components/layout/ListHeaderStack/index.vue'
import SearchBar from '@/components/SearchBar/index.vue'
// Business 组件
import HelpCategoryGrid from '@/components/HelpCategoryGrid/index.vue'
import HelpFaqSection from '@/components/HelpFaqSection/index.vue'
import HelpContactCard from '@/components/HelpContactCard/index.vue'
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
.bottom-placeholder {
  height: 100rpx;
}
</style>
