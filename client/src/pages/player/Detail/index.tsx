import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Avatar, Tag, Rate, Button, Descriptions, Spin, message } from 'antd'
import { UserOutlined, ArrowLeftOutlined } from '@ant-design/icons'
import { useAuth } from '@/context/AuthContext'

interface PlayerDetail {
  id: string
  nickname: string
  avatar?: string
  games: string[]
  rating: number
  price: number
  status: 'online' | 'busy' | 'offline'
  bio: string
  orderCount: number
}

export default function PlayerDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const [player, setPlayer] = useState<PlayerDetail | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // TODO: 调用 API 获取陪玩师详情
    setTimeout(() => {
      setPlayer({
        id: id!,
        nickname: '王者小姐姐',
        games: ['王者荣耀', '和平精英'],
        rating: 4.8,
        price: 30,
        status: 'online',
        bio: '专业王者荣耀陪玩，擅长打野和中单，耐心温柔，欢迎下单~',
        orderCount: 128,
      })
      setLoading(false)
    }, 500)
  }, [id])

  const handleOrder = () => {
    if (!isAuthenticated) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    // TODO: 跳转到下单页面
    message.info('下单功能开发中')
  }

  if (loading) {
    return <Spin size="large" style={{ display: 'block', margin: '100px auto' }} />
  }

  if (!player) {
    return <div>陪玩师不存在</div>
  }

  const statusMap = {
    online: { color: 'green', text: '在线' },
    busy: { color: 'orange', text: '忙碌' },
    offline: { color: 'default', text: '离线' },
  }

  return (
    <div>
      <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>

      <Card>
        <div style={{ display: 'flex', gap: 24 }}>
          <Avatar size={120} icon={<UserOutlined />} src={player.avatar} />
          <div style={{ flex: 1 }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
              <h2 style={{ margin: 0 }}>{player.nickname}</h2>
              <Tag color={statusMap[player.status].color}>
                {statusMap[player.status].text}
              </Tag>
            </div>
            <div style={{ marginBottom: 8 }}>
              {player.games.map(game => (
                <Tag key={game} color="blue">{game}</Tag>
              ))}
            </div>
            <Rate disabled defaultValue={player.rating} allowHalf />
            <span style={{ marginLeft: 8 }}>{player.rating} 分</span>
          </div>
          <div style={{ textAlign: 'right' }}>
            <div style={{ fontSize: 24, color: '#f5222d', fontWeight: 'bold' }}>
              ¥{player.price}/小时
            </div>
            <Button
              type="primary"
              size="large"
              onClick={handleOrder}
              disabled={player.status === 'offline'}
              style={{ marginTop: 16 }}
            >
              立即下单
            </Button>
          </div>
        </div>

        <Descriptions style={{ marginTop: 24 }} column={2}>
          <Descriptions.Item label="个人简介" span={2}>{player.bio}</Descriptions.Item>
          <Descriptions.Item label="完成订单">{player.orderCount} 单</Descriptions.Item>
          <Descriptions.Item label="好评率">98%</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  )
}
