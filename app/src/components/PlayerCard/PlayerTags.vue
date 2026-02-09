<template>
  <view class="player-tags">
    <GlTag v-if="showRank && rank" size="mini" type="warning">{{ rank }}</GlTag>
    <GlTag v-if="showVerified && isVerified" size="mini" plain>认证</GlTag>
    <GlTag
      v-if="showStatusTag"
      size="mini"
      :type="isOnline ? 'success' : 'default'"
      plain
    >
      {{ isOnline ? '在线' : '离线' }}
    </GlTag>
  </view>
</template>

<script setup lang="ts">
import GlTag from '@/components/gl/Tag/index.vue'

interface Props {
  showRank: boolean
  showVerified: boolean
  showStatusTag: boolean
  rank?: string
  isVerified?: boolean
  isOnline?: boolean
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.player-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
}
</style>
