import { isLoggedIn } from '../../utils/auth'
import { http } from '../../utils/request'

interface Conversation {
  id: number
  name: string
  avatar: string
  lastMessage: string
  time: string
  unread: number
  online?: boolean
}

Page({
  data: {
    navPadding: 88,
    isLoggedIn: false,
    loading: true,
    conversations: [] as Conversation[],
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
      isLoggedIn: isLoggedIn(),
    })

    if (this.data.isLoggedIn) {
      this.loadConversations()
    } else {
      this.setData({ loading: false, conversations: this.getMockData() })
    }
  },

  onShow() {
    if (typeof this.getTabBar === 'function' && this.getTabBar()) {
      this.getTabBar().setData({ selected: 2 })
    }

    const loggedIn = isLoggedIn()
    if (loggedIn !== this.data.isLoggedIn) {
      this.setData({ isLoggedIn: loggedIn })
      if (loggedIn) this.loadConversations()
    }
  },

  onPullDownRefresh() {
    if (this.data.isLoggedIn) {
      this.loadConversations().finally(() => wx.stopPullDownRefresh())
    } else {
      wx.stopPullDownRefresh()
    }
  },

  getMockData(): Conversation[] {
    return [
      {
        id: 1,
        name: '小明',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=1',
        lastMessage: '好的，那我们开始吧！',
        time: '2分钟前',
        unread: 2,
        online: true,
      },
      {
        id: 2,
        name: '王者大神',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=2',
        lastMessage: '上次的对局打得不错',
        time: '1小时前',
        unread: 0,
        online: true,
      },
      {
        id: 3,
        name: '和平精英小队',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=3',
        lastMessage: '[图片]',
        time: '昨天',
        unread: 5,
        online: false,
      },
      {
        id: 4,
        name: '客服小助手',
        avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=4',
        lastMessage: '您的订单已完成，感谢使用！',
        time: '3天前',
        unread: 0,
        online: false,
      },
    ]
  },

  async loadConversations() {
    this.setData({ loading: true })
    try {
      const conversations = await http.get<Conversation[]>('/chat/conversations')
      this.setData({ conversations, loading: false })
    } catch (err) {
      console.error('加载会话失败:', err)
      this.setData({ loading: false, conversations: this.getMockData() })
    }
  },

  handleLogin() {
    wx.navigateTo({ url: '/pages/login/index' })
  },

  onNewMessage() {
    wx.navigateTo({ url: '/pages/chat/new' })
  },

  handleConversationTap(e: WechatMiniprogram.CustomEvent) {
    const { message } = e.detail
    wx.navigateTo({ url: `/pages/chat/index?id=${message.id}` })
  },
})
