/**
 * 页面容器组件
 * 统一处理状态栏占位、导航栏占位、TabBar 占位
 */
Component({
  options: {
    multipleSlots: true,
    styleIsolation: 'apply-shared',
  },

  properties: {
    // 是否显示状态栏占位
    showStatusBar: {
      type: Boolean,
      value: true,
    },
    // 是否是 TabBar 页面
    isTabBar: {
      type: Boolean,
      value: false,
    },
    // 自定义背景色
    bgColor: {
      type: String,
      value: '',
    },
    // 是否有自定义头部（如 Banner）
    customHeader: {
      type: Boolean,
      value: false,
    },
    // 主题模式: 'light' | 'dark'
    themeMode: {
      type: String,
      value: '', // 默认为空，会尝试从全局读取
    },
  },

  data: {
    statusBarHeight: 44,
    navBarHeight: 44, // 导航栏高度（胶囊按钮区域）
    topPadding: 88, // 总顶部高度 = 状态栏 + 导航栏
  },

  lifetimes: {
    attached() {
      const systemInfo = wx.getSystemInfoSync()
      const statusBarHeight = systemInfo.statusBarHeight || 44

      // 如果 props 没有传入 themeMode，尝试从全局读取
      if (!this.data.themeMode) {
        const app = getApp<IAppOption>()
        this.setData({
          themeMode: app.globalData.themeMode || 'dark'
        })
      }

      // 获取胶囊按钮位置来计算导航栏高度
      let navBarHeight = 44
      try {
        const menuButton = wx.getMenuButtonBoundingClientRect()
        // 导航栏高度 = 胶囊底部位置 - 状态栏高度 + 底部间距
        navBarHeight = menuButton.bottom - statusBarHeight + (menuButton.top - statusBarHeight)
      } catch (e) {
        navBarHeight = 44
      }

      this.setData({
        statusBarHeight,
        navBarHeight,
        topPadding: statusBarHeight + navBarHeight,
      })
    },
  },
})
