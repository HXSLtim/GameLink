/**
 * 区块标题组件
 * Discord 风格的大写标题
 */
Component({
  options: {
    multipleSlots: true,
    styleIsolation: 'apply-shared',
  },

  properties: {
    title: {
      type: String,
      value: '',
    },
    // 是否显示更多按钮
    showMore: {
      type: Boolean,
      value: false,
    },
    moreText: {
      type: String,
      value: 'See All',
    },
  },

  methods: {
    onMore() {
      this.triggerEvent('more')
    },
  },
})
