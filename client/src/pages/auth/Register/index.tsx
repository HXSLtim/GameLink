import { Form, Input, Button, Card, Radio, message } from 'antd'
import { UserOutlined, LockOutlined, PhoneOutlined } from '@ant-design/icons'
import { useNavigate, Link } from 'react-router-dom'
import '../Login/style.css'

interface RegisterForm {
  phone: string
  password: string
  confirmPassword: string
  role: 'user' | 'player'
}

export default function Register() {
  const navigate = useNavigate()
  const [form] = Form.useForm()

  const handleSubmit = async (values: RegisterForm) => {
    try {
      // TODO: 调用注册 API
      console.log('Register:', values)
      message.success('注册成功，请登录')
      navigate('/login')
    } catch {
      message.error('注册失败')
    }
  }

  return (
    <div className="login-container">
      <Card className="login-card">
        <h2 className="login-title">GameLink 注册</h2>
        <Form form={form} onFinish={handleSubmit} size="large" initialValues={{ role: 'user' }}>
          <Form.Item
            name="phone"
            rules={[{ required: true, message: '请输入手机号' }]}
          >
            <Input prefix={<PhoneOutlined />} placeholder="手机号" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>
          <Form.Item
            name="confirmPassword"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
          </Form.Item>
          <Form.Item name="role" label="注册身份">
            <Radio.Group>
              <Radio value="user">普通用户</Radio>
              <Radio value="player">陪玩师</Radio>
            </Radio.Group>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              注册
            </Button>
          </Form.Item>
        </Form>
        <div className="login-footer">
          已有账号？<Link to="/login">立即登录</Link>
        </div>
      </Card>
    </div>
  )
}
