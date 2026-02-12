<template>
  <view class="gl-empty" :class="{ 'gl-empty--compact': compact }" :style="customStyle">
    <view class="gl-empty__icon">
      <slot name="icon">
        <uv-icon :name="icon" :size="resolvedIconSize" color="var(--color-text-placeholder)" />
      </slot>
    </view>

    <text v-if="title" class="gl-empty__title">{{ title }}</text>
    <text v-if="description" class="gl-empty__desc">{{ description }}</text>

    <view v-if="showActionArea" class="gl-empty__action">
      <slot v-if="$slots.default"></slot>
      <GlButton
        v-else
        type="primary"
        size="small"
        round
        @click="handleAction"
      >
        {{ resolvedActionText }}
      </GlButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'
import GlButton from '@/components/gl/Button/index.vue'

interface Props {
  icon?: string
  iconSize?: number
  title?: string
  description?: string
  showAction?: boolean
  actionText?: string
  compact?: boolean
  customStyle?: string | Record<string, string>
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'empty-data',
  title: '暂无数据',
  description: '',
  showAction: false,
  actionText: '',
  compact: false,
})

const emit = defineEmits<{
  action: []
}>()

const slots = useSlots()

const resolvedIconSize = computed(() => props.iconSize ?? (props.compact ? 80 : 120))
const resolvedActionText = computed(() => props.actionText || '去看看')
const showActionArea = computed(
  () => Boolean(slots.default) || props.showAction || Boolean(props.actionText)
)

const handleAction = () => {
  emit('action')
}
</script>

<style lang="scss" scoped>
.gl-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $gl-spacing-xxl $gl-spacing-lg;
  min-height: 400rpx;
}

.gl-empty__icon {
  margin-bottom: $gl-spacing-lg;
  opacity: 0.8;
  // 增加轻微浮动动画，增加生动感
  animation: float 3s ease-in-out infinite;
}

.gl-empty__title {
  font-size: 32rpx;
  font-weight: 600;
  color: $gl-text-secondary;
  margin-bottom: $gl-spacing-sm;
}

.gl-empty__desc {
  font-size: 26rpx;
  color: $gl-text-placeholder;
  text-align: center;
  max-width: 480rpx;
  line-height: 1.5;
}

.gl-empty__action {
  margin-top: $gl-spacing-xl;
}

.gl-empty--compact {
  padding: $gl-spacing-lg;
  min-height: auto;

  .gl-empty__icon {
    margin-bottom: $gl-spacing-sm;
  }

  .gl-empty__title {
    font-size: 28rpx;
  }

  .gl-empty__desc {
    font-size: 24rpx;
  }

  .gl-empty__action {
    margin-top: $gl-spacing-md;
  }
}

@keyframes float {
  0% { transform: translateY(0); }
  50% { transform: translateY(-10rpx); }
  100% { transform: translateY(0); }
}
</style>