const { shared, timing, Easing } = wx.worklet

Component({
  data: {
    selected: 0,
    list: [
      {
        pagePath: '/pages/index/index',
        text: '首页',
        icon: '/assets/icons/home.svg',
        iconActive: '/assets/icons/home-active.svg',
      },
      {
        pagePath: '/pages/category/index',
        text: '分类',
        icon: '/assets/icons/category.svg',
        iconActive: '/assets/icons/category-active.svg',
      },
      {
        pagePath: '/pages/message/index',
        text: '消息',
        icon: '/assets/icons/message.svg',
        iconActive: '/assets/icons/message-active.svg',
      },
      {
        pagePath: '/pages/profile/index',
        text: '我的',
        icon: '/assets/icons/profile.svg',
        iconActive: '/assets/icons/profile-active.svg',
      },
    ],
    // 动画进度值（0-1）
    animProgress: [1, 0, 0, 0] as number[],
  },

  lifetimes: {
    attached() {
      this.initAnimations()
    },
  },

  methods: {
    initAnimations() {
      // 为每个 tab 创建共享变量
      const progress0 = shared(1) // 默认选中第一个
      const progress1 = shared(0)
      const progress2 = shared(0)
      const progress3 = shared(0)

      this._progressList = [progress0, progress1, progress2, progress3]

      // 为每个 tab 应用动画样式
      this._progressList.forEach((progress, index) => {
        // 空心图标的透明度（1 - progress）
        this.applyAnimatedStyle(`.tab-icon-outline-${index}`, () => {
          'worklet'
          return {
            opacity: 1 - progress.value,
          }
        })

        // 实心图标的动画效果（从左下到右上填充）
        this.applyAnimatedStyle(`.tab-icon-filled-${index}`, () => {
          'worklet'
          const p = progress.value
          // 使用 clip-path 模拟从左下到右上的填充效果
          // 由于 Skyline 不支持 clip-path，改用 transform + opacity 组合
          const scale = 0.5 + p * 0.5 // 0.5 -> 1
          const translateX = -26 * (1 - p) // -26rpx -> 0
          const translateY = 26 * (1 - p) // 26rpx -> 0
          return {
            opacity: p,
            transform: `translate(${translateX}rpx, ${translateY}rpx) scale(${scale})`,
          }
        })
      })
    },

    switchTab(e: WechatMiniprogram.TouchEvent) {
      const { index, path } = e.currentTarget.dataset
      const currentSelected = this.data.selected

      if (currentSelected === index) return

      // 动画：旧 tab 淡出，新 tab 填充进入
      const oldProgress = this._progressList[currentSelected]
      const newProgress = this._progressList[index]

      // 旧 tab 淡出
      oldProgress.value = timing(0, {
        duration: 150,
        easing: Easing.ease,
      })

      // 新 tab 填充动画（从左下到右上）
      newProgress.value = timing(1, {
        duration: 300,
        easing: Easing.out(Easing.cubic),
      })

      wx.switchTab({
        url: path,
      })
    },
  },
})
