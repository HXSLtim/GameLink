import { useState } from 'react'
import { View } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { Button, Input, Form, Toast } from '@nutui/nutui-react-taro'
import { useUserStore } from '@/stores/user'
import './index.scss'

const Login = () => {
    const [loading, setLoading] = useState(false)
    const [showPassword, setShowPassword] = useState(false)
    const { setToken, setCurrentUser } = useUserStore()

    // Form instance
    const [form] = Form.useForm()

    const handleLogin = async (values: any) => {
        setLoading(true)
        console.log('Login values:', values)
        try {
            // Mock Login API Call
            // const res = await authApi.login(values)

            setTimeout(() => {
                const mockToken = 'mock-jwt-token-123'
                const mockUser = {
                    id: 1,
                    nickname: 'Gamer001',
                    avatar: '',
                    currentRole: 'user',
                    isPlayer: false
                }

                setToken(mockToken)
                setCurrentUser(mockUser as any)

                Toast.show('登录成功', { duration: 2 })

                setTimeout(() => {
                    Taro.switchTab({ url: '/pages/index/index' })
                }, 1000)
            }, 1000)

        } catch (error) {
            Toast.show('登录失败', { duration: 2 })
        } finally {
            setLoading(false)
        }
    }

    return (
        <View className="login-page">
            <View className="login-header">
                <View className="logo">GL</View>
                <View className="title">GameLink</View>
                <View className="subtitle">连接无限游戏乐趣</View>
            </View>

            <View className="login-form">
                <Form
                    form={form}
                    onFinish={handleLogin}
                    footer={
                        <View style={{ marginTop: '20px' }}>
                            <Button block type="primary" loading={loading} nativeType="submit">
                                立即登录
                            </Button>
                        </View>
                    }
                >
                    <Form.Item
                        name="username"
                        label="账号"
                        rules={[{ required: true, message: '请输入用户名' }]}
                    >
                        <Input
                            placeholder="请输入用户名/手机号"
                            clearable
                        />
                    </Form.Item>
                    <Form.Item
                        name="password"
                        label="密码"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input
                            type={showPassword ? 'text' : 'password'}
                            placeholder="请输入密码"
                        />
                    </Form.Item>
                </Form>
            </View>

            <Toast />
        </View>
    )
}

export default Login
