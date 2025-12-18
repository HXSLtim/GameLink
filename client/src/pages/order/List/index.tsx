import { useState, useEffect } from 'react'
import { Card, Table, Tag, Button, Empty } from 'antd'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'

interface Order {
  id: string
  playerName: string
  game: string
  hours: number
  amount: number
  status: 'pending' | 'accepted' | 'in_progress' | 'completed' | 'cancelled'
  createdAt: string
}

export default function OrderList() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login')
      return
    }
    // TODO: 调用 API 获取订单列表
    setTimeout(() => {
      setOrders([
        { id: '1', playerName: '王者小姐姐', game: '王者荣耀', hours: 2, amount: 60, status: 'completed', createdAt: '2025-01-15 14:30' },
        { id: '2', playerName: '吃鸡大神', game: '和平精英', hours: 1, amount: 50, status: 'in_progress', createdAt: '2025-01-16 10:00' },
      ])
      setLoading(false)
    }, 500)
  }, [isAuthenticated, navigate])

  const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'default', text: '待接单' },
    accepted: { color: 'blue', text: '已接单' },
    in_progress: { color: 'processing', text: '进行中' },
    completed: { color: 'success', text: '已完成' },
    cancelled: { color: 'error', text: '已取消' },
  }

  const columns = [
    { title: '订单号', dataIndex: 'id', key: 'id' },
    { title: '陪玩师', dataIndex: 'playerName', key: 'playerName' },
    { title: '游戏', dataIndex: 'game', key: 'game' },
    { title: '时长', dataIndex: 'hours', key: 'hours', render: (h: number) => `${h}小时` },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (a: number) => `¥${a}` },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={statusMap[status]?.color}>{statusMap[status]?.text}</Tag>
      ),
    },
    { title: '下单时间', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: '操作',
      key: 'action',
      render: (_: unknown, record: Order) => (
        <Button type="link" size="small">
          {record.status === 'completed' ? '评价' : '详情'}
        </Button>
      ),
    },
  ]

  if (!isAuthenticated) {
    return null
  }

  return (
    <Card title="我的订单">
      {orders.length === 0 && !loading ? (
        <Empty description="暂无订单">
          <Button type="primary" onClick={() => navigate('/players')}>
            去找陪玩
          </Button>
        </Empty>
      ) : (
        <Table
          columns={columns}
          dataSource={orders}
          rowKey="id"
          loading={loading}
        />
      )}
    </Card>
  )
}
