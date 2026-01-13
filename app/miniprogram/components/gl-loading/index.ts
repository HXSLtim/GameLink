/**
 * 加载组件
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    // 加载文字
    text: {
      type: String,
      value: 'Loading...',
    },
    // 是否显示文字
    showText: {
      type: Boolean,
      value: true,
    },
    // 尺寸: small | medium | large
    size: {
      type: String,
      value: 'medium',
    },
    // 是否全屏居中
    fullscreen: {
      type: Boolean,
      value: false,
    },
  },
})
