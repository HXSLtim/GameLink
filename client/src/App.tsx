import { Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from '@/context/AuthContext'
import { ThemeProvider } from '@/context/ThemeContext'
import DiscordLayout from '@/layouts/DiscordLayout'
import Home from '@/pages/Home'
import Login from '@/pages/auth/Login'
import Register from '@/pages/auth/Register'
import PlayerList from '@/pages/player/List'
import PlayerDetail from '@/pages/player/Detail'
import OrderList from '@/pages/order/List'
import Profile from '@/pages/user/Profile'

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Routes>
          {/* 公开路由 */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />

          {/* 主布局路由 */}
          <Route path="/" element={<DiscordLayout />}>
            <Route index element={<Home />} />
            <Route path="players" element={<PlayerList />} />
            <Route path="players/:id" element={<PlayerDetail />} />
            <Route path="orders" element={<OrderList />} />
            <Route path="profile" element={<Profile />} />
          </Route>

          {/* 404 重定向 */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </ThemeProvider>
  )
}

export default App
