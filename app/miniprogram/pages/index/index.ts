// Mock 数据
const CATEGORIES = [
  { id: '1', name: 'LoL', icon: '⚔️' },
  { id: '2', name: 'Valorant', icon: '🎯' },
  { id: '3', name: '王者荣耀', icon: '👑' },
  { id: '4', name: 'Apex', icon: '🛡️' },
  { id: '5', name: '原神', icon: '✨' },
]

const PLAYERS = [
  {
    id: '1',
    name: 'Nana Chan',
    avatar: 'https://picsum.photos/100/100?random=1',
    imageUrl: 'https://picsum.photos/400/500?random=1',
    tags: ['Sweet Voice', 'Carry', 'Humorous'],
    price: 35,
    rating: 5.0,
    orders: 1208,
    game: '王者荣耀',
    audioDuration: 12,
    isOnline: true,
  },
  {
    id: '2',
    name: 'Pro Zed',
    avatar: 'https://picsum.photos/100/100?random=2',
    imageUrl: 'https://picsum.photos/400/500?random=2',
    tags: ['Ex-Pro', 'Coaching', 'Challenger'],
    price: 60,
    rating: 4.9,
    orders: 450,
    game: 'League of Legends',
    audioDuration: 8,
    isOnline: false,
  },
  {
    id: '3',
    name: 'Viper Mist',
    avatar: 'https://picsum.photos/100/100?random=3',
    imageUrl: 'https://picsum.photos/400/500?random=3',
    tags: ['Tactical', 'Chill', 'Immortal'],
    price: 45,
    rating: 4.8,
    orders: 890,
    game: 'Valorant',
    audioDuration: 15,
    isOnline: true,
  },
]

Page({
  data: {
    navStyle: 'height: 87px;',
    categories: CATEGORIES,
    players: PLAYERS,
    hasNotification: true,
  },

  onLoad() {
    // 计算导航栏占位高度（状态栏 + 胶囊按钮区域）
    const systemInfo = wx.getSystemInfoSync()
    const statusBarHeight = systemInfo.statusBarHeight || 44

    let navPadding = statusBarHeight + 44 // 默认值

    try {
      const menuButton = wx.getMenuButtonBoundingClientRect()
      // 胶囊按钮底部位置 + 底部间距
      navPadding = menuButton.bottom + (menuButton.top - statusBarHeight)
    } catch (e) {
      console.log('获取胶囊按钮位置失败', e)
    }

    console.log('navPadding:', navPadding)

    this.setData({
      navStyle: `height: ${navPadding}px;`,
    })
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 0 })
    }
  },

  onSearch(e: WechatMiniprogram.CustomEvent) {
    const { value } = e.detail
    console.log('Search:', value)
  },

  onBellTap() {
    wx.navigateTo({ url: '/pages/notification/index' })
  },

  onCategoryTap(e: WechatMiniprogram.TouchEvent) {
    const { id } = e.currentTarget.dataset
    wx.navigateTo({ url: `/pages/player-list/index?gameId=${id}` })
  },

  onSeeAll() {
    wx.switchTab({ url: '/pages/category/index' })
  },

  onPlayerTap(e: WechatMiniprogram.CustomEvent) {
    const { player } = e.detail
    wx.navigateTo({ url: `/pages/player/index?id=${player.id}` })
  },
})
