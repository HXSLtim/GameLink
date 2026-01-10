import { useState, useEffect } from 'react'
import { Card, Row, Col, Avatar, Tag, Rate, Input, Select, Spin } from 'antd'
import { useNavigate } from 'react-router-dom'
import { UserOutlined } from '@ant-design/icons'

const { Search } = Input

interface Player {
  id: string
  nickname: string
  avatar?: string
  games: string[]
  rating: number
  price: number
  status: 'online' | 'busy' | 'offline'
}

export default function PlayerList() {
  const navigate = useNavigate()
  const [players, setPlayers] = useState<Player[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // TODO: 调用 API 获取陪玩师列表
    setTimeout(() => {
      setPlayers([
        { id: '1', nickname: '王者小姐姐', games: ['王者荣耀', '和平精英'], rating: 4.8, price: 30, status: 'online' },
        { id: '2', nickname: '吃鸡大神', games: ['和平精英', 'CSGO'], rating: 4.9, price: 50, status: 'online' },
        { id: '3', nickname: 'LOL高手', games: ['英雄联盟', '云顶之弈'], rating: 4.7, price: 40, status: 'busy' },
      ])
      setLoading(false)
    }, 500)
  }, [])

  const statusMap = {
    online: { color: 'green', text: '在线' },
    busy: { color: 'orange', text: '忙碌' },
    offline: { color: 'default', text: '离线' },
  }

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          <Col flex="auto">
            <Search placeholder="搜索陪玩师" allowClear />
          </Col>
          <Col>
            <Select placeholder="选择游戏" style={{ width: 150 }} allowClear>
              <Select.Option value="王者荣耀">王者荣耀</Select.Option>
              <Select.Option value="和平精英">和平精英</Select.Option>
              <Select.Option value="英雄联盟">英雄联盟</Select.Option>
            </Select>
          </Col>
        </Row>
      </Card>

      <Spin spinning={loading}>
        <Row gutter={[16, 16]}>
          {players.map(player => (
            <Col xs={24} sm={12} md={8} lg={6} key={player.id}>
              <Card
                hoverable
                onClick={() => navigate(`/players/${player.id}`)}
              >
                <Card.Meta
                  avatar={<Avatar size={64} icon={<UserOutlined />} src={player.avatar} />}
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      {player.nickname}
                      <Tag color={statusMap[player.status].color}>
                        {statusMap[player.status].text}
                      </Tag>
                    </div>
                  }
                  description={
                    <>
                      <div style={{ marginBottom: 8 }}>
                        {player.games.map(game => (
                          <Tag key={game}>{game}</Tag>
                        ))}
                      </div>
                      <Rate disabled defaultValue={player.rating} allowHalf style={{ fontSize: 14 }} />
                      <div style={{ marginTop: 8, color: '#f5222d', fontWeight: 'bold' }}>
                        ¥{player.price}/小时
                      </div>
                    </>
                  }
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Spin>
    </div>
  )
}
