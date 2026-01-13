import { http } from '../../utils/request'

interface Category {
  id: number
  name: string
  icon: string
  playerCount: number
}

Page({
  data: {
    navPadding: 88,
    loading: true,
    categories: [] as Category[],
  },

  onLoad() {
    // 计算导航栏占位高度
    const systemInfo = wx.getSystemInfoSync()
    const statusBarHeight = systemInfo.statusBarHeight || 44
    let navBarHeight = 44

    try {
      const menuButton = wx.getMenuButtonBoundingClientRect()
      navBarHeight = menuButton.height + (menuButton.top - statusBarHeight) * 2
    } catch (e) {
      navBarHeight = 44
    }

    this.setData({
      navPadding: statusBarHeight + navBarHeight,
    })

    this.loadCategories()
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 1 })
    }
  },

  onPullDownRefresh() {
    this.loadCategories().finally(() => wx.stopPullDownRefresh())
  },

  getMockCategories(): Category[] {
    return [
      { id: 1, name: '王者荣耀', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=wzry', playerCount: 1280 },
      { id: 2, name: '和平精英', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=hpjy', playerCount: 856 },
      { id: 3, name: '英雄联盟', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=lol', playerCount: 2100 },
      { id: 4, name: '永劫无间', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=yjwj', playerCount: 320 },
      { id: 5, name: '原神', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=ys', playerCount: 450 },
      { id: 6, name: '第五人格', icon: 'https://api.dicebear.com/7.x/shapes/svg?seed=d5rg', playerCount: 280 },
    ]
  },

  async loadCategories() {
    this.setData({ loading: true })
    try {
      const categories = await http.get<Category[]>('/games')
      this.setData({ categories, loading: false })
    } catch (err) {
      console.error('加载分类失败:', err)
      this.setData({ loading: false, categories: this.getMockCategories() })
    }
  },

  onSearch(e: WechatMiniprogram.CustomEvent) {
    const { value } = e.detail
    console.log('Search games:', value)
  },

  handleCategoryTap(e: WechatMiniprogram.CustomEvent) {
    const { game } = e.detail
    wx.navigateTo({ url: `/pages/player-list/index?gameId=${game.id}` })
  },

  onTagTap(e: WechatMiniprogram.TouchEvent) {
    const tag = e.currentTarget.dataset.tag
    console.log('Tag tapped:', tag)
  },
})
