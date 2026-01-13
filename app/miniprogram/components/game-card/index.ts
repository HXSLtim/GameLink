/**
 * 游戏分类卡片组件
 */
Component({
  options: {
    styleIsolation: 'apply-shared',
  },

  properties: {
    // 游戏数据
    game: {
      type: Object,
      value: {},
    },
    // 是否显示在线人数
    showCount: {
      type: Boolean,
      value: true,
    },
  },

  methods: {
    onTap() {
      this.triggerEvent('tap', { game: this.data.game })
    },
  },
})
