/**
 * 按钮组件
 * 支持多种变体和尺寸
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    // 按钮类型（variant 的别名，优先使用 type）
    type: {
      type: String,
      value: '',
    },
    // 按钮变体: primary | secondary | outline | ghost | gold
    variant: {
      type: String,
      value: 'primary',
    },
    // 按钮尺寸: small | medium | large (或 sm | md | lg)
    size: {
      type: String,
      value: 'medium',
    },
    // 是否占满宽度
    block: {
      type: Boolean,
      value: false,
    },
    // 是否禁用
    disabled: {
      type: Boolean,
      value: false,
    },
    // 是否加载中
    loading: {
      type: Boolean,
      value: false,
    },
  },

  data: {
    computedVariant: 'primary',
    computedSize: 'md',
  },

  observers: {
    'type, variant': function (type: string, variant: string) {
      // type 优先于 variant
      this.setData({ computedVariant: type || variant || 'primary' })
    },
    size: function (size: string) {
      // 支持 small/medium/large 和 sm/md/lg 两种写法
      const sizeMap: Record<string, string> = {
        small: 'sm',
        medium: 'md',
        large: 'lg',
        sm: 'sm',
        md: 'md',
        lg: 'lg',
      }
      this.setData({ computedSize: sizeMap[size] || 'md' })
    },
  },

  methods: {
    handleTap() {
      if (!this.data.disabled && !this.data.loading) {
        this.triggerEvent('tap')
      }
    },
  },
})
