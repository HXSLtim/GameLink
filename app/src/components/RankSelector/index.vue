<template>
  <view class="rank-selector">
    <!-- 已选段位展示 -->
    <view v-if="selectedRank" class="selected-rank">
      <image v-if="selectedRank.icon" :src="selectedRank.icon" class="rank-icon" mode="aspectFit" />
      <text class="rank-name">{{ selectedRank.name }}</text>
    </view>
    
    <!-- 选择按钮 -->
    <view class="select-trigger" @tap="showPicker">
      <text class="trigger-text">{{ selectedRank ? selectedRank.name : placeholder }}</text>
      <text class="trigger-arrow">›</text>
    </view>
    
    <!-- 段位选择弹窗 -->
    <view v-if="visible" class="picker-mask" @tap="hidePicker">
      <view class="picker-container" @tap.stop>
        <!-- 头部 -->
        <view class="picker-header">
          <text class="picker-title">选择段位</text>
          <text class="picker-close" @tap="hidePicker">×</text>
        </view>
        
        <!-- 游戏切换 -->
        <scroll-view v-if="games.length > 1" class="game-tabs" scroll-x>
          <view
            v-for="game in games"
            :key="game.id"
            class="game-tab"
            :class="{ 'game-tab-active': currentGameId === game.id }"
            @tap="switchGame(game.id)"
          >
            {{ game.name }}
          </view>
        </scroll-view>
        
        <!-- 段位列表 -->
        <scroll-view class="picker-content" scroll-y>
          <view class="rank-grid">
            <view
              v-for="rank in currentRanks"
              :key="rank.id"
              class="rank-item"
              :class="{ 'rank-item-selected': tempSelected?.id === rank.id }"
              @tap="selectRank(rank)"
            >
              <image v-if="rank.icon" :src="rank.icon" class="rank-icon" mode="aspectFit" />
              <text class="rank-name">{{ rank.name }}</text>
              <view v-if="tempSelected?.id === rank.id" class="rank-check">✓</view>
            </view>
          </view>
          
          <view v-if="currentRanks.length === 0" class="picker-empty">
            <GlEmpty title="暂无段位数据" compact />
          </view>
        </scroll-view>
        
        <!-- 底部按钮 -->
        <view class="picker-footer">
          <view class="picker-btn cancel" @tap="hidePicker">取消</view>
          <view class="picker-btn confirm" @tap="confirmSelect">确定</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import GlEmpty from '@/components/gl/Empty/index.vue'
import type { GameOption, GameRankOption } from '@/types/game'

const props = withDefaults(defineProps<{
  modelValue?: GameRankOption | null
  ranks?: GameRankOption[]
  games?: GameOption[]
  gameId?: number | string
  placeholder?: string
}>(), {
  modelValue: null,
  ranks: () => [],
  games: () => [],
  placeholder: '请选择段位',
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: GameRankOption | null): void
  (e: 'change', value: GameRankOption | null): void
}>()

const visible = ref(false)
const tempSelected = ref<GameRankOption | null>(null)
const currentGameId = ref<number | string | undefined>(props.gameId)

// 已选段位
const selectedRank = computed(() => props.modelValue)

// 当前游戏的段位列表
const currentRanks = computed(() => {
  if (!currentGameId.value) return props.ranks
  return props.ranks.filter(rank => rank.gameId === currentGameId.value)
})

// 显示选择器
const showPicker = () => {
  tempSelected.value = selectedRank.value
  if (props.games.length > 0 && !currentGameId.value) {
    currentGameId.value = props.games[0].id
  }
  visible.value = true
}

// 隐藏选择器
const hidePicker = () => {
  visible.value = false
}

// 切换游戏
const switchGame = (gameId: number | string) => {
  currentGameId.value = gameId
  tempSelected.value = null
}

// 选择段位
const selectRank = (rank: GameRankOption) => {
  tempSelected.value = rank
}

// 确认选择
const confirmSelect = () => {
  emit('update:modelValue', tempSelected.value)
  emit('change', tempSelected.value)
  hidePicker()
}
</script>

<style lang="scss" scoped>
.rank-selector {
  width: 100%;
}

.selected-rank {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-sm);
  
  .rank-icon {
    width: 48rpx;
    height: 48rpx;
  }
  
  .rank-name {
    font-size: var(--font-sm);
    color: var(--color-text);
  }
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) var(--spacing-md);
  background: var(--color-bg-card);
  border: 1rpx solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  @include press-effect;
  
  .trigger-text {
    font-size: var(--font-sm);
    color: var(--color-text-secondary);
  }
  
  .trigger-arrow {
    font-size: var(--font-md);
    color: var(--color-text-placeholder);
  }
}

.picker-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: flex-end;
}

.picker-container {
  width: 100%;
  max-height: 80vh;
  background: var(--color-bg-card);
  border-radius: var(--radius-md) var(--radius-md) 0 0;
  border-top: 1rpx solid var(--color-border);
  display: flex;
  flex-direction: column;
}

.picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
  
  .picker-title {
    font-size: var(--font-md);
    font-weight: 600;
    color: var(--color-text);
  }
  
  .picker-close {
    font-size: var(--font-lg);
    color: var(--color-text-secondary);
  }
}

.game-tabs {
  white-space: nowrap;
  padding: var(--spacing-xs) var(--spacing-md);
  border-bottom: 1rpx solid var(--color-border);
  @include hide-scrollbar;
}

.game-tab {
  display: inline-block;
  padding: var(--spacing-xs) var(--spacing-md);
  margin-right: var(--spacing-xs);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  font-size: var(--font-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  @include press-effect;
  
  &-active {
    background: var(--color-bg-secondary);
    color: var(--color-text);
  }
}

.picker-content {
  flex: 1;
  max-height: 50vh;
  padding: var(--spacing-md) var(--spacing-md);
}

.rank-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-sm);
}

.rank-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-sm);
  border: 1rpx solid var(--color-border);
  position: relative;
  cursor: pointer;
  @include press-effect;
  
  &-selected {
    border-color: var(--color-border);
    background: var(--color-bg-secondary);
  }
  
  .rank-icon {
    width: 64rpx;
    height: 64rpx;
    margin-bottom: var(--spacing-xs);
  }
  
  .rank-name {
    font-size: var(--font-sm);
    color: var(--color-text);
    text-align: center;
  }
  
  .rank-check {
    position: absolute;
    top: 8rpx;
    right: 8rpx;
    width: 28rpx;
    height: 28rpx;
    background: var(--color-primary);
    color: #FFFFFF;
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--font-sm);
  }
}

.picker-empty {
  padding: var(--spacing-xl) 0;
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--font-sm);
  grid-column: 1 / -1;
}

.picker-footer {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md);
  padding-bottom: calc(var(--spacing-md) + env(safe-area-inset-bottom));
  border-top: 1rpx solid var(--color-border);
  
  .picker-btn {
    flex: 1;
    height: 72rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    font-size: var(--font-sm);
    font-weight: 500;
    @include press-effect;
    
    &.cancel {
      background: var(--color-bg-secondary);
      color: var(--color-text);
      border: 1rpx solid var(--color-border);
    }
    
    &.confirm {
      background: var(--color-primary);
      color: #FFFFFF;
      box-shadow: none;
    }
  }
}
</style>
