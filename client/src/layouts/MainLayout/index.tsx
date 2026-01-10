import { Outlet, useNavigate } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import { HomeOutlined, TeamOutlined, OrderedListOutlined, UserOutlined } from '@ant-design/icons'
import { useAuth } from '@/context/AuthContext'
import './style.css'

const { Header, Content, Footer } = Layout

export default function MainLayout() {
  const navigate = useNavigate()
  const { user, isAuthenticated, logout } = useAuth()

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/players', icon: <TeamOutlined />, label: '找陪玩' },
    { key: '/orders', icon: <OrderedListOutlined />, label: '我的订单' },
  ]

  const userMenuItems = [
    { key: 'profile', label: '个人中心' },
    { key: 'logout', label: '退出登录' },
  ]

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
  }

  const handleUserMenuClick = ({ key }: { key: string }) => {
    if (key === 'logout') {
      logout()
      navigate('/login')
    } else if (key === 'profile') {
      navigate('/profile')
    }
  }

  return (
    <Layout className="main-layout">
      <Header className="header">
        <div className="logo" onClick={() => navigate('/')}>
          GameLink
        </div>
        <Menu
          theme="dark"
          mode="horizontal"
          items={menuItems}
          onClick={handleMenuClick}
          className="nav-menu"
        />
        <div className="header-right">
          {isAuthenticated ? (
            <Dropdown menu={{ items: userMenuItems, onClick: handleUserMenuClick }}>
              <Avatar icon={<UserOutlined />} src={user?.avatar} />
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => navigate('/login')}>
              登录
            </Button>
          )}
        </div>
      </Header>
      <Content className="content">
        <Outlet />
      </Content>
      <Footer className="footer">
        © 2025 GameLink. All rights reserved.
      </Footer>
    </Layout>
  )
}
