import client from './client'

export interface LoginParams {
  phone: string
  password: string
}

export interface RegisterParams {
  phone: string
  password: string
  role: 'user' | 'player'
}

export interface LoginResponse {
  token: string
  user: {
    id: string
    nickname: string
    avatar?: string
    phone: string
    role: 'user' | 'player'
  }
}

export const authApi = {
  login: (params: LoginParams) => 
    client.post<unknown, LoginResponse>('/auth/login', params),
  
  register: (params: RegisterParams) => 
    client.post('/auth/register', params),
  
  logout: () => 
    client.post('/auth/logout'),
  
  getProfile: () => 
    client.get('/user/profile'),
}
