import { Card, Descriptions, Avatar, Button, Form, Input, message } from 'antd'
import { UserOutlined, EditOutlined } from '@ant-design/icons'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'

export default function Profile() {
  const navigate = useNavigate()
  const { user, isAuthenticated, logout } = useAuth()
  const [editing, setEditing] = useState(false)
  const [form] = Form.useForm()

  if (!isAuthenticated) {
    navigate('/login')
    return null
  }

  const handleSave = async (values: { nickname: string }) => {
    try {
      // TODO: 调用 API 更新用户信息
      console.log('Update profile:', values)
      message.success('保存成功')
      setEditing(false)
    } catch {
      message.error('保存失败')
    }
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <Card
      title="个人中心"
      extra={
        !editing && (
          <Button icon={<EditOutlined />} onClick={() => setEditing(true)}>
            编辑
          </Button>
        )
      }
    >
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Avatar size={100} icon={<UserOutlined />} src={user?.avatar} />
      </div>

      {editing ? (
        <Form
          form={form}
          initialValues={{ nickname: user?.nickname }}
          onFinish={handleSave}
          layout="vertical"
          style={{ maxWidth: 400, margin: '0 auto' }}
        >
          <Form.Item
            name="nickname"
            label="昵称"
            rules={[{ required: true, message: '请输入昵称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ marginRight: 8 }}>
              保存
            </Button>
            <Button onClick={() => setEditing(false)}>取消</Button>
          </Form.Item>
        </Form>
      ) : (
        <Descriptions column={1} style={{ maxWidth: 400, margin: '0 auto' }}>
          <Descriptions.Item label="昵称">{user?.nickname}</Descriptions.Item>
          <Descriptions.Item label="手机号">{user?.phone || '未绑定'}</Descriptions.Item>
          <Descriptions.Item label="身份">
            {user?.role === 'player' ? '陪玩师' : '普通用户'}
          </Descriptions.Item>
        </Descriptions>
      )}

      <div style={{ textAlign: 'center', marginTop: 24 }}>
        <Button danger onClick={handleLogout}>
          退出登录
        </Button>
      </div>
    </Card>
  )
}
