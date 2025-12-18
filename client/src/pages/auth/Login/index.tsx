import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import './style.css'

interface LoginForm {
  phone: string
  password: string
}

export default function Login() {
  const navigate = useNavigate()
  const { login } = useAuth()
  const [form] = Form.useForm()

  const handleSubmit = async (values: LoginForm) => {
    try {
      // TODO: 调用登录 API
      console.log('Login:', values)
      message.success('登录成功')
      // 模拟登录
      login('mock-token', {
        id: '1',
        nickname: '测试用户',
        phone: values.phone,
        role: 'user',
      })
      navigate('/')
    } catch {
      message.error('登录失败')
    }
  }

  return (
    <div className="login-container">
      <Card className="login-card">
        <h2 className="login-title">GameLink 登录</h2>
        <Form form={form} onFinish={handleSubmit} size="large">
          <Form.Item
            name="phone"
            rules={[{ required: true, message: '请输入手机号' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="手机号" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              登录
            </Button>
          </Form.Item>
        </Form>
        <div className="login-footer">
          还没有账号？<Link to="/register">立即注册</Link>
        </div>
      </Card>
    </div>
  )
}
