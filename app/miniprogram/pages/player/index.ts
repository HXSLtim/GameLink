import { http } from '../../utils/request'
import { checkLogin } from '../../utils/auth'

interface Service {
  id: number
  name: string
  description: string
  price: number
  unit: string
}

interface PlayerDetail {
  id: number
  nickname: string
  avatar: string
  level: number
  online: boolean
  verified: boolean
  price: number
  orderCount: number
  rating: string
  responseRate: number
  completionRate: number
  games: string[]
  intro?: string
  services: Service[]
}

Page({
  data: {
    loading: true,
    playerId: 0,
    player: {} as PlayerDetail,
  },
  
  onLoad(options) {
    const playerId = Number(options.id) || 0
    this.setData({ playerId })
    this.loadPlayer()
  },
  
  async loadPlayer() {
    this.setData({ loading: true })
    try {
      const player = await http.get<PlayerDetail>(`/players/${this.data.playerId}`)
      this.setData({ player, loading: false })
    } catch (err) {
      console.error('加载陪玩师详情失败:', err)
      // 模拟数据
      this.setData({
        loading: false,
        player: {
          id: this.data.playerId,
          nickname: '小甜甜',
          avatar: '',
          level: 5,
          online: true,
          verified: true,
          price: 30,
          orderCount: 128,
          rating: '4.9',
          responseRate: 98,
          completionRate: 99,
          games: ['王者荣耀', '和平精英'],
          intro: '王者荣耀国服韩信，带飞上分稳稳的~声音甜美，性格开朗，欢迎来撩~',
          services: [
            { id: 1, name: '王者荣耀-陪玩', description: '陪你上分，轻松愉快', price: 30, unit: '局' },
            { id: 2, name: '王者荣耀-代练', description: '快速上分，安全可靠', price: 50, unit: '星' },
            { id: 3, name: '和平精英-陪玩', description: '吃鸡带飞，稳稳的', price: 25, unit: '局' },
          ],
        },
      })
    }
  },
  
  handleOrder() {
    if (!checkLogin()) return
    
    const { playerId } = this.data
    wx.navigateTo({ url: `/pages/order/create/index?playerId=${playerId}` })
  },
})
