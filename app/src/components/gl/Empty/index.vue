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
  padding: var(--spacing-xl) var(--spacing-lg);
}

.gl-empty__icon {
  margin-bottom: var(--spacing-md);
  opacity: 0.6;
}

.gl-empty__title {
  font-size: var(--font-lg);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--spacing-xs);
}

.gl-empty__desc {
  font-size: var(--font-sm);
  color: var(--color-text-placeholder);
  text-align: center;
  max-width: 400rpx;
}

.gl-empty__action {
  margin-top: var(--spacing-lg);
}

.gl-empty--compact {
  padding: var(--spacing-md);
  
  .gl-empty__icon {
    margin-bottom: var(--spacing-sm);
  }
  
  .gl-empty__title {
    font-size: var(--font-base);
  }
  
  .gl-empty__desc {
    font-size: var(--font-xs);
  }
  
  .gl-empty__action {
    margin-top: var(--spacing-md);
  }
}
</style>
