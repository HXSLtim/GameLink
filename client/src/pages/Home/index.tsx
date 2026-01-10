import { Card, Row, Col, Typography, Button } from 'antd'
import { useNavigate } from 'react-router-dom'
import { TeamOutlined, SafetyOutlined, ThunderboltOutlined } from '@ant-design/icons'

const { Title, Paragraph } = Typography

export default function Home() {
  const navigate = useNavigate()

  const features = [
    {
      icon: <TeamOutlined style={{ fontSize: 48, color: '#1890ff' }} />,
      title: '专业陪玩',
      desc: '海量认证陪玩师，技术过硬，服务贴心',
    },
    {
      icon: <SafetyOutlined style={{ fontSize: 48, color: '#52c41a' }} />,
      title: '安全保障',
      desc: '平台担保交易，资金安全有保障',
    },
    {
      icon: <ThunderboltOutlined style={{ fontSize: 48, color: '#faad14' }} />,
      title: '快速匹配',
      desc: '智能推荐系统，快速找到合适的陪玩',
    },
  ]

  return (
    <div>
      {/* Hero Section */}
      <Card style={{ textAlign: 'center', marginBottom: 24 }}>
        <Title>欢迎来到 GameLink</Title>
        <Paragraph style={{ fontSize: 16, color: '#666' }}>
          专业的游戏陪玩服务平台，让游戏更有趣
        </Paragraph>
        <Button type="primary" size="large" onClick={() => navigate('/players')}>
          立即找陪玩
        </Button>
      </Card>

      {/* Features */}
      <Row gutter={24}>
        {features.map((feature, index) => (
          <Col xs={24} md={8} key={index}>
            <Card style={{ textAlign: 'center', height: '100%' }}>
              {feature.icon}
              <Title level={4} style={{ marginTop: 16 }}>{feature.title}</Title>
              <Paragraph style={{ color: '#666' }}>{feature.desc}</Paragraph>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}
