<template>
  <view class="game-grid-container">
    <uv-grid :col="cols" :border="false">
      <uv-grid-item
        v-for="game in displayGames"
        :key="game.id"
        @click="handleSelect(game.id)"
      >
        <view class="grid-game-item" :class="{ active: modelValue === game.id }">
          <image v-if="game.icon" :src="game.icon" class="grid-game-icon" mode="aspectFit" />
          <view v-else class="grid-game-placeholder">
            <uv-icon name="grid" size="20" color="var(--color-text-secondary)"></uv-icon>
          </view>
          <text class="grid-game-name">{{ game.name }}</text>
        </view>
      </uv-grid-item>
      
      <!-- 展开/收起按钮 -->
      <uv-grid-item v-if="expandable && hasMore" @click="toggleExpand">
        <view class="grid-game-item expand-item">
          <view class="grid-game-placeholder expand-icon">
            <uv-icon 
              :name="expanded ? 'arrow-up' : 'arrow-down'" 
              size="20" 
              color="var(--color-primary)"
            ></uv-icon>
          </view>
          <text class="grid-game-name expand-text">
            {{ expanded ? collapseText : expandText }}
          </text>
        </view>
      </uv-grid-item>
    </uv-grid>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

export interface GameItem {
  id: number | string
  name: string
  icon?: string
}

interface Props {
  games: GameItem[]
  modelValue: number | string
  cols?: number
  expandable?: boolean
  maxVisible?: number
  expandText?: string
  collapseText?: string
}

const props = withDefaults(defineProps<Props>(), {
  cols: 4,
  expandable: true,
  maxVisible: 3, // 默认显示 3 个 + 1 个展开按钮
  expandText: '更多',
  collapseText: '收起',
})

const emit = defineEmits<{
  'update:modelValue': [id: number | string]
  select: [game: GameItem]
}>()

const expanded = ref(false)

const hasMore = computed(() => {
  return props.games.length > props.maxVisible
})

const displayGames = computed(() => {
  if (!props.expandable || expanded.value) {
    return props.games
  }
  return props.games.slice(0, props.maxVisible)
})

const toggleExpand = () => {
  expanded.value = !expanded.value
}

const handleSelect = (id: number | string) => {
  emit('update:modelValue', id)
  const game = props.games.find(g => g.id === id)
  if (game) {
    emit('select', game)
  }
}
</script>

<style lang="scss" scoped>
.game-grid-container {
  background: var(--color-bg-card);
  border-bottom: 1rpx solid var(--color-border);
  padding: 16rpx 8rpx;
}

.grid-game-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
  padding: 12rpx 4rpx;
  border-radius: 12rpx;
  transition: all 0.2s;
  width: 100%;
  
  &.active {
    background: rgba(0, 210, 106, 0.15);
    
    .grid-game-name {
      color: var(--color-primary);
      font-weight: 600;
    }
    
    .grid-game-icon {
      transform: scale(1.05);
    }
  }
  
  &.expand-item {
    .grid-game-placeholder {
      background: transparent;
      border: 2rpx dashed var(--color-primary);
    }
    
    .expand-text {
      color: var(--color-primary);
    }
  }
  
  &:active {
    transform: scale(0.95);
  }
}

.grid-game-icon {
  width: 56rpx;
  height: 56rpx;
  border-radius: 10rpx;
  transition: transform 0.2s;
}

.grid-game-placeholder {
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: 10rpx;
  
  &.expand-icon {
    background: transparent;
    border: 2rpx dashed var(--color-primary);
  }
}

.grid-game-name {
  font-size: 22rpx;
  color: var(--color-text);
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
  line-height: 1.2;
}
</style>
