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
    showBack: {
      type: Boolean,
      value: true,
    },
    transparent: {
      type: Boolean,
      value: false,
    },
  },
  
  data: {
    statusBarHeight: 20,
    navBarHeight: 44,
  },
  
  lifetimes: {
    attached() {
      const systemInfo = wx.getSystemInfoSync()
      this.setData({
        statusBarHeight: systemInfo.statusBarHeight || 20,
        navBarHeight: 44,
      })
    },
  },
  
  methods: {
    handleBack() {
      const pages = getCurrentPages()
      if (pages.length > 1) {
        wx.navigateBack()
      } else {
        wx.switchTab({ url: '/pages/index/index' })
      }
    },
  },
})
