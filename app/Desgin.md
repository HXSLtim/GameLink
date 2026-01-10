🚀 GameLink 小程序技术栈与开发规划
Super Dev Team - 综合专家组

版本: v1.0

日期: 2025-01-10

状态: 🟢 规划完成，准备实施

📋 执行摘要

┌─────────────────────────────────────────────────────────────┐
│  技术栈决策矩阵                                                │
├─────────────────────────────────────────────────────────────┤
│  框架:     Taro 3.x         React/Vue3     京东开源        │
│  状态:     Pinia              Vue3        Vue官方          │
│  UI:       NutUI              Taro        京东组件库        │
│  语言:     TypeScript 5+     类型安全     微软开发          │
│  构建:     Vite 5.x          极速构建     Evan You          │
├─────────────────────────────────────────────────────────────┤
│  核心优势:                                                   │
│  ✅ 一套代码，多端运行（微信/支付宝/抖音/H5）              │
│  ✅ React/Vue3 双生态支持，团队技术栈匹配                  │
│  ✅ 开箱即用的组件库，开发效率提升 50%                     │
│  ✅ 类型安全，减少 80% 的运行时错误                       │
│  ✅ 完善的工具链，调试、热更新、性能分析                    │
└─────────────────────────────────────────────────────────────┘
1️⃣ 技术选型
框架选择：Taro 3.x ✅
框架 优势 劣势 评分 推荐度
Taro 3.x • 京东开源，生态成熟<br>• 支持 React/Vue3<br>• 多端编译<br>• 社区活跃 • 编译产物较大<br>• 部分原生 API 适配延迟 9.5/10 ⭐⭐⭐⭐⭐ 强烈推荐
uni-app • DCloud 出品<br>• Vue 技术栈<br>• 插件丰富 • Vue 3 支持不完善<br>• 性能略逊 8.5/10 ⭐⭐⭐⭐ 推荐
原生小程序 • 性能最优<br>• API 最新 • 开发效率低<br>• 维护成本高<br>• 无法跨平台 6.0/10 ⭐⭐⭐ 不推荐
最终决策: Taro 3.x + React（可切换至 Vue3）

选择理由:

✅ 生态成熟: 京东 40+ 团队维护，社区活跃
✅ 技术匹配: 后端使用 Go，前端团队熟悉 React
✅ 跨平台: 一套代码可编译到微信/支付宝/抖音/H5
✅ 性能优秀: 运行时性能接近原生（95%+）
✅ 开发体验: 支持 React Hooks，热更新快
✅ 类型安全: 完美支持 TypeScript
状态管理：Pinia ✅
方案 优势 劣势 评分 推荐度
Pinia • Vue 官方推荐<br>• API 简洁<br>• TypeScript 友好<br>• 无需 mutations • 生态较新 9.5/10 ⭐⭐⭐⭐⭐ 强烈推荐
Zustand • 轻量简洁<br>• 无需 Provider • React 生态更强 9.0/10 ⭐⭐⭐⭐⭐ 推荐
Redux Toolkit • 生态成熟<br>• 中间件丰富 • 学习曲线陡<br>• 样板代码多 8.0/10 ⭐⭐⭐⭐ 可选
最终决策: Pinia（Vue3）或 Zustand（React）

选择理由:

✅ API 简洁: 无需 mutations，代码减少 50%
✅ TypeScript: 完美类型推导
✅ DevTools: 出色的开发工具支持
✅ 模块化: 天然支持代码分割
UI 组件库：NutUI ✅
组件库 优势 劣势 评分 推荐度
NutUI • 京东出品<br>• Taro 官方推荐<br>• 组件丰富（80+）<br>• 设计精美 • 文档部分中文 9.0/10 ⭐⭐⭐⭐⭐ 强烈推荐
自定义组件 • 完全可控<br>• 轻量级 • 开发成本高<br>• 维护负担 7.0/10 ⭐⭐⭐ 不推荐
最终决策: NutUI + 少量自定义组件

选择理由:

✅ 官方支持: Taro 官方推荐和维护
✅ 设计统一: 京东设计语言，专业美观
✅ 组件完整: 涵盖 80+ 常用组件
✅ 持续更新: 跟随 Taro 版本更新
技术栈总览

┌─────────────────────────────────────────────────────────────┐
│  GameLink 小程序技术栈                                     │
├─────────────────────────────────────────────────────────────┤
│  框架层                                                      │
│  ├─ Taro 3.x + React              或  Taro 3.x + Vue3           │
│  ├─ TypeScript 5+                                           │
│  └─ Vite 5.x                                                │
├─────────────────────────────────────────────────────────────┤
│  状态管理                                                    │
│  ├─ Pinia (Vue3)                  或  Zustand (React)         │
│  └─ Taro.setStorageSync (持久化)                            │
├─────────────────────────────────────────────────────────────┤
│  UI 组件                                                     │
│  ├─ NutUI React (Taro)             或  NutUI Vue (Taro)        │
│  ├─ 样式系统：SCSS                                           │
│  └─ 响应式：Taro.rpx                                         │
├─────────────────────────────────────────────────────────────┤
│  网络请求                                                    │
│  ├─ Taro.request (封装)                                      │
│  ├─ Axios (备用)                                              │
│  └─ WebSocket (Taro Socket)                                   │
├─────────────────────────────────────────────────────────────┤
│  工具库                                                      │
│  ├─ Day.js (日期)                                            │
│  ├─ Lodash-es (工具)                                         │
│  ├─ Crypto-js (加密)                                         │
│  └─ Socket.io-client (WebSocket)                              │
├─────────────────────────────────────────────────────────────┤
│  开发工具                                                    │
│  ├─ VS Code + 插件                                           │
│  ├─ Taro CLI (开发工具)                                       │
│  ├─ 微信开发者工具                                             │
│  └─ ESLint + Prettier                                         │
├─────────────────────────────────────────────────────────────┤
│  测试工具                                                    │
│  ├─ Vitest (单元测试)                                        │
│  ├─ Taro Simulator (模拟器)                                  │
│  └─ Puppeteer (E2E 测试)                                     │
└─────────────────────────────────────────────────────────────┘
2️⃣ 项目架构设计
目录结构

miniprogram/
├── src/
│   ├── app.config.ts           # 全局配置
│   ├── app.ts                  # 应用入口
│   │
│   ├── assets/                 # 静态资源
│   │   ├── images/             # 图片
│   │   ├── icons/              # 图标
│   │   └── fonts/              # 字体
│   │
│   ├── components/             # 全局组件
│   │   ├── common/             # 通用组件
│   │   │   ├── Button/        # 按钮
│   │   │   ├── Card/          # 卡片
│   │   │   ├── List/          # 列表
│   │   │   ├── Input/         # 输入框
│   │   │   └── Modal/         # 弹窗
│   │   │
│   │   ├── business/           # 业务组件
│   │   │   ├── PlayerCard/    # 陪玩师卡片
│   │   │   ├── OrderCard/     # 订单卡片
│   │   │   ├── ChatBubble/    # 聊天气泡
│   │   │   └── PriceDisplay/  # 价格显示
│   │   │
│   │   └── layout/             # 布局组件
│   │       ├── CustomTabBar/  # 自定义 TabBar
│   │       ├── NavigationBar/  # 导航栏
│   │       └── SafeArea/      # 安全区域
│   │
│   ├── pages/                  # 页面
│   │   ├── index/             # 首页
│   │   │   └── index.tsx
│   │   │
│   │   ├── search/            # 搜索页
│   │   ├── player-detail/     # 陪玩师详情
│   │   ├── order-create/      # 下单页
│   │   ├── order-list/        # 订单列表
│   │   ├── chat/              # 聊天页
│   │   ├── wallet/            # 钱包页
│   │   ├── profile/           # 个人中心
│   │   │
│   │   └── player/            # 陪玩师端页面
│   │       ├── workbench/      # 工作台
│   │       ├── order-hall/     # 接单大厅
│   │       └── earnings/       # 收入页
│   │
│   ├── services/               # 服务层
│   │   ├── api/                # API 客户端
│   │   │   ├── request.ts     # 请求封装
│   │   │   ├── user.ts        # 用户 API
│   │   │   ├── player.ts      # 陪玩师 API
│   │   │   ├── order.ts       # 订单 API
│   │   │   ├── chat.ts        # 聊天 API
│   │   │   ├── wallet.ts      # 钱包 API
│   │   │   └── auth.ts        # 认证 API
│   │   │
│   │   ├── auth/              # 认证服务
│   │   │   ├── token.ts       # Token 管理
│   │   │   └── storage.ts     # 存储服务
│   │   │
│   │   └── websocket/         # WebSocket
│   │       └── socket.ts      # Socket 封装
│   │
│   ├── stores/                 # 状态管理
│   │   ├── user.ts             # 用户状态
│   │   ├── player.ts           # 陪玩师状态
│   │   ├── order.ts            # 订单状态
│   │   ├── chat.ts             # 聊天状态
│   │   └── app.ts              # 应用状态
│   │
│   ├── hooks/                  # 自定义 Hooks
│   │   ├── useAuth.ts          # 认证 Hook
│   │   ├── usePlayer.ts        # 陪玩师 Hook
│   │   ├── useOrder.ts         # 订单 Hook
│   │   ├── useChat.ts          # 聊天 Hook
│   │   └── useWebSocket.ts     # WebSocket Hook
│   │
│   ├── utils/                  # 工具函数
│   │   ├── date.ts             # 日期处理
│   │   ├── format.ts           # 格式化
│   │   ├── validate.ts         # 验证
│   │   ├── storage.ts          # 存储
│   │   └── crypto.ts           # 加密
│   │
│   ├── constants/              # 常量
│   │   ├── api.ts              # API 地址
│   │   ├── config.ts           # 配置
│   │   └── enum.ts             # 枚举
│   │
│   ├── types/                  # 类型定义
│   │   ├── user.ts             # 用户类型
│   │   ├── player.ts           # 陪玩师类型
│   │   ├── order.ts            # 订单类型
│   │   ├── chat.ts             # 聊天类型
│   │   └── api.ts              # API 响应类型
│   │
│   └── styles/                 # 全局样式
│       ├── variables.scss      # 样式变量
│       ├── mixins.scss         # 样式混入
│       └── global.scss         # 全局样式
│
├── project.config.json         # 项目配置
├── tsconfig.json                # TypeScript 配置
├── vite.config.ts               # Vite 配置
├── .eslintrc.js                  # ESLint 配置
└── package.json                 # 依赖配置
模块划分方案

graph TB
    subgraph "页面层 (Pages)"
        A[首页]
        B[搜索页]
        C[陪玩师详情]
        D[下单页]
        E[订单页]
        F[聊天页]
        G[钱包页]
        H[个人中心]
        I[工作台]
        J[接单大厅]
        K[收入页]
    end

    subgraph "业务组件 (Business Components)"
        BA[PlayerCard]
        BB[OrderCard]
        BC[ChatBubble]
        BD[PriceDisplay]
    end
    
    subgraph "通用组件 (Common Components)"
        CA[Button]
        CB[Card]
        CC[List]
        CD[Input]
        CE[Modal]
    end
    
    subgraph "服务层 (Services)"
        SA[API Client]
        SB[Auth Service]
        SC[WebSocket Service]
    end
    
    subgraph "状态管理 (Stores)"
        STA[User Store]
        STB[Player Store]
        STC[Order Store]
        STD[Chat Store]
    end
    
    A --> CA
    A --> CB
    C --> BA
    D --> BD
    E --> BB
    F --> BC
    
    A --> STA
    I --> STB
    D --> STC
    F --> STD
    
    A --> SA
    F --> SC
API 客户端封装方案

// services/api/request.ts

import Taro from '@tarojs/taro';
import { getToken, refreshToken, clearToken } from '@/services/auth/storage';

interface RequestConfig {
  url: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  data?: any;
  headers?: Record<string, string>;
  skipAuth?: boolean;
}

interface ApiResponse<T = any> {
  success: boolean;
  code: number;
  message: string;
  data: T;
  traceId?: string;
}

class ApiClient {
  private baseURL: string;
  private timeout: number = 10000;

  constructor() {
    this.baseURL = process.env.TARO_APP_API_BASE_URL || '<http://localhost:8080/api/v1>';
  }

  /**

* 统一请求方法
   */
  async request<T = any>(config: RequestConfig): Promise<ApiResponse<T>> {
    const { url, method = 'GET', data, headers = {}, skipAuth = false } = config;

    // 添加认证头
    if (!skipAuth) {
      const token = getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    try {
      const response = await Taro.request({
        url: `${this.baseURL}${url}`,
        method,
        data,
        header: {
          'Content-Type': 'application/json',
          ...headers,
        },
        timeout: this.timeout,
      });

      const result = response.data as ApiResponse<T>;

      // Token 自动刷新
      if (result.code === 401 && !skipAuth) {
        const refreshed = await refreshToken();
        if (refreshed) {
          // 重试原请求
          return this.request<T>(config);
        }
      }

      return result;
    } catch (error: any) {
      // 网络错误处理
      console.error('API request failed:', error);
      throw error;
    }
  }

  /**

* GET 请求
   */
  get<T = any>(url: string, params?: any, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'GET',
      data: params,
      ...config,
    });
  }

  /**

* POST 请求
   */
  post<T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'POST',
      data,
      ...config,
    });
  }

  /**

* PUT 请求
   */
  put<T = any>(url: string, data?: any, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'PUT',
      data,
      ...config,
    });
  }

  /**

* DELETE 请求
   */
  delete<T = any>(url: string, config?: Partial<RequestConfig>): Promise<ApiResponse<T>> {
    return this.request<T>({
      url,
      method: 'DELETE',
      ...config,
    });
  }
}

export default new ApiClient();
状态管理方案

// stores/user.ts

import { defineStore } from 'pinia';
import { User } from '@/types/user';
import { userApi } from '@/services/api/user';

interface UserState {
  currentUser: User | null;
  currentRole: 'user' | 'player';
  isPlayer: boolean;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    currentUser: null,
    currentRole: 'user',
    isPlayer: false,
  }),

  getters: {
    /**
     *获取用户 ID
     */
    userId(): number {
      return this.currentUser?.id || 0;
    },

    /**
     * 获取用户头像
     */
    avatar(): string {
      return this.currentUser?.avatar || '';
    },

    /**
     * 获取用户昵称
     */
    nickname(): string {
      return this.currentUser?.nickname || '';
    },

    /**
     * 是否已登录
     */
    isLoggedIn(): boolean {
      return !!this.currentUser;
    },
  },

  actions: {
    /**
     *设置当前用户
     */
    setCurrentUser(user: User | null) {
      this.currentUser = user;
      this.currentRole = user?.currentRole || 'user';
      this.isPlayer = user?.isPlayer || false;
    },

    /**
     * 更新用户信息
     */
    async updateUserInfo() {
      if (!this.isLoggedIn) return;

      try {
        const response = await userApi.getProfile();
        if (response.success) {
          this.setCurrentUser(response.data);
        }
      } catch (error) {
        console.error('Failed to update user info:', error);
      }
    },

    /**
     * 切换角色
     */
    async switchRole(role: 'user' | 'player') {
      try {
        const response = await userApi.switchRole({ role });
        if (response.success) {
          // 更新 Token
          const { accessToken, refreshToken } = response.data;
          // TODO: 保存新 Token
          
          // 更新状态
          this.currentRole = role;
        }
      } catch (error) {
        console.error('Failed to switch role:', error);
        throw error;
      }
    },

    /**
     * 登出
     */
    logout() {
      this.currentUser = null;
      this.currentRole = 'user';
      this.isPlayer = false;
      // TODO: 清除 Token
    },
  },
});
3️⃣ 开发规范
命名规范

文件命名:
├── 页面组件: PascalCase     例: PlayerDetail.tsx
├── 业务组件: PascalCase     例: PlayerCard.tsx
├── 通用组件: PascalCase     例: CustomButton.tsx
├── Hooks: camelCase + use  例: useAuth.ts
├── 工具函数: camelCase     例: formatDate.ts
└── 常量: UPPER_SNAKE_CASE   例: API_BASE_URL

变量命名:
├── 组件状态: camelCase       例: isLoading
├── 函数/方法: camelCase     例: handleSubmit
├── 常量: UPPER_SNAKE_CASE   例: MAX_RETRY
├── 类/接口: PascalCase       例: UserService
└── 类型别名: PascalCase     例: UserResponse

CSS 命名:
├── 类名: kebab-case          例: custom-button
├── 变量: kebab-case          例: --primary-color
└── ID: kebab-case            例: submit-button
代码风格

// 组件结构标准
import React, { useState, useEffect } from 'react';
import { View, Text } from '@tarojs/components';
import { useLoad } from '@tarojs/taro';
import CustomButton from '@/components/common/Button';
import { useUserStore } from '@/stores/user';
import { formatDate } from '@/utils/date';
import type { Player } from '@/types/player';

import './styles.scss';

/**

* 陪玩师卡片组件
*
* @component
* @example
* <PlayerCard
* player={player}
* onContact={() => handleContact(player.id)}
* />
 */
const PlayerCard: React.FC<PlayerCardProps> = ({ player, onContact }) => {
  // 1. 状态声明
  const [isFavorite, setIsFavorite] = useState(false);
  
  // 2. Hooks
  const { currentUser } = useUserStore();
  
  // 3. 副作用
  useEffect(() => {
    checkFavorite();
  }, [player.id]);

  // 4. 事件处理
  const handleContact = (playerId: number) => {
    if (onContact) {
      onContact(playerId);
    }
  };

  // 5. 渲染
  return (
    <View className="player-card">
      {/*组件内容*/}
    </View>
  );
};

export default PlayerCard;
组件开发规范

// 组件 Props 规范
interface ComponentProps {
  // 必填属性
  required: string;
  
  // 可选属性
  optional?: string;
  
  // 回调函数
  onAction?: (data: any) => void;
  
  // 渲染函数
  renderCustom?: (data: any) => JSX.Element;
}

// 组件注释规范
/**

* 组件描述
*
* @component
* @example
* <ComponentName
* required="value"
* optional="value"
* />
*
* @param {string} required - 必填参数描述
* @param {string} [optional] - 可选参数描述
* @returns {JSX.Element} 组件渲染结果
 */
API 调用规范

// services/api/order.ts

import apiClient from './request';
import type { Order, CreateOrderRequest, OrderResponse } from '@/types/order';

/**

* 订单 API
 */
export const orderApi = {
  /**
  * 创建订单
   */
  create: (data: CreateOrderRequest) => {
    return apiClient.post<OrderResponse>('/user/orders', data);
  },

  /**

* 获取订单列表
   */
  list: (params: {
    status?: string;
    page: number;
    pageSize: number;
  }) => {
    return apiClient.get<{
      items: Order[];
      pagination: Pagination;
    }>('/user/orders', params);
  },

  /**

* 获取订单详情
   */
  detail: (id: number) => {
    return apiClient.get<Order>(`/user/orders/${id}`);
  },

  /**

* 取消订单
   */
  cancel: (id: number, reason: string) => {
    return apiClient.post(`/user/orders/${id}/cancel`, { reason });
  },

  /**

* 确认订单完成
   */
  confirm: (id: number, data: {
    rating: number;
    tags: string[];
    comment: string;
  }) => {
    return apiClient.post(`/user/orders/${id}/confirm`, data);
  },
};
错误处理规范

// utils/error-handler.ts

/**

* 错误处理类
 */
class ErrorHandler {
  /**
  * 处理 API 错误
   */
  static handleApiError(error: any): void {
    console.error('API Error:', error);

    const { code, message } = error;

    switch (code) {
      case 401:
        Taro.showToast({
          title: '登录已过期',
          icon: 'none',
          duration: 2000,
        });
        // 跳转登录
        setTimeout(() => {
          Taro.navigateTo({
            url: '/pages/login/index',
          });
        }, 2000);
        break;

      case 403:
        Taro.showToast({
          title: '权限不足',
          icon: 'none',
        });
        break;

      case 404:
        Taro.showToast({
          title: '资源不存在',
          icon: 'none',
        });
        break;

      case 500:
        Taro.showToast({
          title: '服务器错误',
          icon: 'none',
        });
        break;

      default:
        Taro.showToast({
          title: message || '请求失败',
          icon: 'none',
        });
    }
  }

  /**

* 处理网络错误
   */
  static handleNetworkError(error: any): void {
    console.error('Network Error:', error);
    Taro.showToast({
      title: '网络连接失败',
      icon: 'none',
    });
  }

  /**

* 处理业务错误
   */
  static handleBusinessError(error: any): void {
    const { message } = error;
    Taro.showToast({
      title: message || '操作失败',
      icon: 'none',
    });
  }
}

export default ErrorHandler;
4️⃣ 开发计划
总体时间线

gantt
    title GameLink 小程序开发计划
    dateFormat  YYYY-MM-DD
    section Phase 1
    基础架构搭建           :done, p1, 2025-01-10, 7d
    section Phase 2
    用户端核心功能         :p2, 2025-01-17, 21d
    section Phase 3
    陪玩师端核心功能       :p3, 2025-02-07, 14d
    section Phase 4
    营销功能              :p4, 2025-02-21, 7d
    section Phase 5
    测试与优化             :p5, 2025-02-28, 7d
Phase 1: 基础架构搭建 (Week 1)

目标: 搭建项目基础架构，实现核心功能
┌─────────────────────────────────────────────────────────────┐
│  Week 1 (1.10 - 1.16)                                       │
├─────────────────────────────────────────────────────────────┤
│  Day 1-2: 项目初始化                                         │
│  ├─ 使用 Taro CLI 创建项目                                 │
│  ├─ 配置 TypeScript + ESLint + Prettier                       │
│  ├─ 安装依赖 (NutUI, Pinia, Day.js, Lodash)                │
│  └─ 配置 Vite 构建工具                                      │
│                                                              │
│  Day 3-4: 核心架构                                           │
│  ├─ 创建目录结构                                            │
│  ├─ 实现路由配置                                            │
│  ├─ 封装 API 客户端                                         │
│  ├─ 实现认证服务 (Token 管理)                              │
│  └─ 实现状态管理 (Pinia Stores)                            │
│                                                              │
│  Day 5-7: 通用组件                                           │
│  ├─ 实现布局组件 (NavigationBar, TabBar, SafeArea)          │
│  ├─ 实现通用组件 (Button, Card, Input, Modal, List)         │
│  ├─ 实现工具函数 (date, format, validate, storage)         │
│  └─ 实现全局样式 (variables, mixins, global)               │
└─────────────────────────────────────────────────────────────┘

✅ 里程碑: 项目可运行，基础架构完成
Phase 2: 用户端核心功能 (Week 2-5)

目标: 实现用户端核心下单流程
┌─────────────────────────────────────────────────────────────┐
│  Week 2 (1.17 - 1.23)                                       │
├─────────────────────────────────────────────────────────────┤
│  首页                                                        │
│  ├─ 游戏横向滚动列表                                        │
│  ├─ 精选陪玩师列表                                          │
│  ├─ 充值优惠横幅                                            │
│  └─ 功能入口 (优惠券, VIP, 活动)                            │
│                                                              │
│  Week 3 (1.24 - 1.30)                                       │
├─────────────────────────────────────────────────────────────┤
│  搜索 + 详情                                                  │
│  ├─ 搜索页 (搜索框 + 历史 + 热门)                          │
│  ├─ 游戏详情页                                              │
│  ├─ 陪玩师详情页 (头像 + 信息 + 游戏 + 评价)               │
│  └─ 收藏功能 (添加/取消/检查)                               │
│                                                              │
│  Week 4 (1.31 - 2.06)                                       │
├─────────────────────────────────────────────────────────────┤
│  下单流程                                                    │
│  ├─ 下单页 (选择游戏 + 时长 + 特殊要求)                    │
│  ├─ 价格计算 (VIP 折扣 + 优惠券)                           │
│  ├─ 支付集成 (微信支付 + 钱包)                             │
│  └─ 订单确认页                                              │
│                                                              │
│  Week 5 (2.07 - 2.13)                                       │
├─────────────────────────────────────────────────────────────┤
│  订单 + 聊天                                                 │
│  ├─ 订单列表页 (全部/待支付/服务中/已完成)                    │
│  ├─ 订单详情页                                              │
│  ├─ 聊天页 (发送消息 + 消息列表)                           │
│  └─ WebSocket 实时消息                                     │
└─────────────────────────────────────────────────────────────┘

✅ 里程碑: 用户端可完整下单流程
Phase 3: 陪玩师端核心功能 (Week 6-7)

目标: 实现陪玩师端接单和收入管理
┌─────────────────────────────────────────────────────────────┐
│  Week 6 (2.14 - 2.20)                                       │
├─────────────────────────────────────────────────────────────┤
│  工作台 + 接单大厅                                            │
│  ├─ 工作台 (今日数据 + 在线状态)                           │
│  ├─ 接单大厅 (可接订单列表)                                │
│  ├─ 接单功能 (接单/拒单)                                    │
│  └─ 在线状态管理 (上线/下线/忙碌)                          │
│                                                              │
│  Week 7 (2.21 - 2.27)                                       │
├─────────────────────────────────────────────────────────────┤
│  订单管理 + 收入                                             │
│  ├─ 订单管理 (服务中订单 + 完成服务)                        │
│  ├─ 收入页 (收入统计 + 收入明细)                           │
│  ├─ 提现功能 (申请提现 + 提现记录)                           │
│  └─ 认证功能 (成为陪玩师)                                   │
└─────────────────────────────────────────────────────────────┘

✅ 里程碑: 陪玩师端可完整接单流程
Phase 4: 营销功能 (Week 8)

目标: 实现营销功能模块
┌─────────────────────────────────────────────────────────────┐
│  Week 8 (2.28 - 3.05)                                       │
├─────────────────────────────────────────────────────────────┤
│  营销模块                                                    │
│  ├─ VIP 会员 (购买 + 使用 + 折扣显示)                        │
│  ├─ 优惠券 (领取 + 使用 + 我的优惠券)                         │
│  ├─ 充值功能 (档位 + 微信支付 + 赠送)                       │
│  ├─ 活动列表 (活动详情 + 参与活动)                           │
│  └─ 推荐功能 (推荐码 + 分享 + 奖励)                           │
└─────────────────────────────────────────────────────────────┘

✅ 里程碑: 营销功能完整
Phase 5: 测试与优化 (Week 9-10)

目标: 测试、优化、上线准备
┌─────────────────────────────────────────────────────────────┐
│  Week 9 (3.06 - 3.12)                                       │
├─────────────────────────────────────────────────────────────┤
│  功能测试                                                    │
│  ├─ 单元测试 (Vitest)                                       │
│  ├─ 组件测试                                                │
│  ├─ 集成测试                                                │
│  └─ Bug 修复                                                 │
│                                                              │
│  Week 10 (3.13 - 3.19)                                      │
├─────────────────────────────────────────────────────────────┤
│  性能优化 + 上线准备                                          │
│  ├─ 性能优化 (首屏加载、图片压缩、代码分包)                │
│  ├─ 安全检查 (权限检查、数据加密)                           │
│  ├─ 提交微信审核                                              │
│  └─ 生产环境部署                                              │
└─────────────────────────────────────────────────────────────┘

✅ 里程碑: 小程序正式上线
5️⃣ 性能优化策略
首屏加载优化

目标: 首屏加载时间 < 2 秒

优化策略:
┌─────────────────────────────────────────────────────────────┐
│  1. 分包加载                                                  │
│  ├─ 主包: 首页 + 登录               ~200KB                   │
│  ├─ 分包1: 搜索 + 详情             ~150KB                   │
│  ├─ 分包2: 订单 + 支付             ~180KB                   │
│  └─ 分包3: 聊天 + 钱包             ~120KB                   │
│                                                              │
│  2. 预加载资源                                                │
│  ├─ 首屏图片优先加载                                        │
│  ├─ 关键组件预加载                                          │
│  └─ API 数据预取                                            │
│                                                              │
│  3. 骨架优化                                                  │
│  ├─ 使用 Taro.setStorageSync 缓存 Token                     │
│  ├─ 使用 React.memo 避免重渲染                               │
│  └─ 虚拟列表处理长列表                                      │
└─────────────────────────────────────────────────────────────┘
图片资源优化

优化策略:
┌─────────────────────────────────────────────────────────────┐
│  1. 图片压缩                                                  │
│  ├─ JPG 压缩: 质量 80%                                     │
│  ├─ PNG 压缩: 使用 TinyPNG                                  │
│  └─ WebP 转换: 支持 WebP 的平台使用 WebP                    │
│                                                              │
│   2. 图片尺寸                                                  │
│  ├─ 头像: 200x200                                           │
│  ├─ 游戏图标: 150x150                                         │
│  ├─ 商品图: 750x400 (2:1)                                  │
│  └── Banner: 750x300 (2.5:1)                                 │
│                                                              │
│  3. CDN 加速                                                   │
│  ├─ 使用 CDN 存储图片                                         │
│  ├─ 配置图片缓存策略                                         │
│  └─ 支持图片懒加载                                            │
└─────────────────────────────────────────────────────────────┘
代码分包策略

// app.config.ts

export default {
  pages: [
    'pages/index/index',              // 首页
    'pages/search/index',            // 搜索
    'pages/player-detail/index',     // 陪玩师详情
    'pages/order-create/index',      // 下单
    'pages/order-list/index',        // 订单列表
    'pages/chat/index',              // 聊天
    'pages/wallet/index',            // 钱包
    'pages/profile/index',           // 个人中心
  ],
  
  // 分包配置
  subPackages: [
    {
      root: 'pages/order/',
      pages: [
        'pages/order-detail/index',
      ],
    },
    {
      root: 'pages/player/',
      pages: [
        'pages/player/workbench/index',      // 陪玩师工作台
        'pages/player/order-hall/index',      // 接单大厅
        'pages/player/earnings/index',        // 收入
      ],
    },
  ],
  
  // 预下载包
  preloadRule: [
    {
      network: 'all',
      packages: ['packages/player'],
    },
  ],
};
缓存策略

缓存层级:
┌─────────────────────────────────────────────────────────────┐
│  1. 本地缓存 (Taro.setStorageSync)                            │
│  ├─ Token (持久化)                                            │
│  ├─ 用户信息 (1小时)                                          │
│  ├─ 陪玩师列表 (30分钟)                                      │
│  └─ 游戏列表 (24小时)                                        │
│                                                              │
│   2. 内存缓存 (状态管理)                                      │
│  ├─ 当前用户 (状态管理)                                       │
│  ├─ 购物车/订单 (页面级)                                    │
│  └─ WebSocket 连接 (全局)                                   │
│                                                              │
│  3. CDN 缓存                                                  │
│  ├─ 静态资源 (JS, CSS, 图片)                               │
│  └─ 缓存策略: Cache-Control: max-age=2592000                 │
└─────────────────────────────────────────────────────────────┘
网络请求优化

// 请求优化策略

// 1. 请求合并
interface BatchRequest {
  [key: string]: Promise<any>;
}

async function batchRequest(requests: BatchRequest): Promise<void> {
  // 并行执行多个请求
  await Promise.all(Object.values(requests));
}

// 2. 请求取消
class RequestManager {
  private pendingRequests = new Map<string, AbortController>();

  async request(url: string, options: RequestInit) {
    // 取消相同的 pending 请求
    const key = `${options.method}:${url}`;
    if (this.pendingRequests.has(key)) {
      this.pendingRequests.get(key)!.abort();
    }

    const controller = new AbortController();
    this.pendingRequests.set(key, controller);

    options.signal = controller.signal;

    try {
      const response = await fetch(url, options);
      return response.json();
    } finally {
      this.pendingRequests.delete(key);
    }
  }
}

// 3. 请求重试
async function fetchWithRetry(
  url: string,
  options: RequestInit,
  maxRetries = 3
): Promise<Response> {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      if (response.ok) {
        return response;
      }
    } catch (error) {
      if (i === maxRetries - 1) throw error;
      await new Promise(resolve => setTimeout(resolve, 1000 * (i + 1)));
    }
  }
  throw new Error('Max retries exceeded');
}
6️⃣ 质量保证
单元测试方案

测试框架: Vitest + @testing-library/react

测试覆盖:
┌─────────────────────────────────────────────────────────────┐
│  组件测试: 80%+                                              │
│  ├─ 通用组件 (Button, Card, Input, Modal)                    │
│  ├─ 业务组件 (PlayerCard, OrderCard, ChatBubble)             │
│  └─ 布局组件 (NavigationBar, TabBar)                        │
│                                                              │
│  Hooks 测试: 70%+                                             │
│  ├─ useAuth (认证)                                            │
│  ├─ usePlayer (陪玩师)                                       │
│  ├─ useOrder (订单)                                          │
│  └─ useWebSocket (聊天)                                     │
│                                                              │
│  工具函数测试: 90%+                                           │
│  ├─ formatDate (日期格式化)                                  │
│  ├─ formatPrice (价格格式化)                                 │
│  ├─ validateEmail (邮箱验证)                                  │
│  └─ crypto-js (加密)                                        │
└─────────────────────────────────────────────────────────────┘
测试示例:

// components/common/Button/__tests__/Button.test.tsx

import { render, fireEvent } from '@testing-library/react';
import Button from '../index';

describe('Button Component', () => {
  it('should render correctly', () => {
    const { getByText } = render(<Button>Click Me</Button>);
    expect(getByText('Click Me')).toBeInTheDocument();
  });

  it('should call onClick when clicked', () => {
    const handleClick = jest.fn();
    const { getByText } = render(
      <Button onClick={handleClick}>Click Me</Button>
    );

    fireEvent.click(getByText('Click Me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should be disabled when loading is true', () => {
    const { getByText } = render(
      <Button loading={true}>Click Me</Button>
    );

    expect(getByText('Click Me')).toBeDisabled();
  });
});
E2E 测试方案

测试框架: Puppeteer + Taro Simulator

测试场景:
┌─────────────────────────────────────────────────────────────┐
│  1. 用户注册/登录                                            │
│  ├─ 微信授权登录                                            │
│  ├─ 手机号验证码登录                                        │
│  └─ Token 自动刷新                                          │
│                                                              │
│  2. 下单流程                                                  │
│  ├─ 搜索陪玩师                                              │
│  ├─ 查看陪玩师详情                                          │
│  ├─ 创建订单                                                │
│  ├─ 选择优惠券                                              │
│  ├─ 微信支付                                                │
│  └─ 查看订单状态                                          │
│                                                              │
│  3. 聊天流程                                                  │
│  ├─ 发送消息                                                │
│  ├─ 接收消息                                                │
│  └─ WebSocket 连接                                          │
│                                                              │
│   4. 角色切换                                                  │
│  ├─ 用户切换到陪玩师                                        │
│  ├─ 陪玩师切换到用户                                        │
│  └─ 切换后功能验证                                          │
└─────────────────────────────────────────────────────────────┘
E2E 测试示例:

// e2e/order-flow.test.ts

import { test, expect } from '@playwright/test';

test.describe('Order Flow', () => {
  test('complete order flow', async ({ page }) => {
    // 1. 登录
    await page.goto('<https://example.com>');
    await page.click('text=登录');
    // ... 登录流程

    // 2. 搜索陪玩师
    await page.fill('input[placeholder="搜索"]', '技术流');
    await page.press('Enter');

    // 3. 查看详情
    await page.click('.player-card:first-child');
    await page.waitForURL('**/player-detail/**');
    await expect(page.locator('text=立即下单')).toBeVisible();

    // 4. 下单
    await page.click('text=立即下单');
    await page.selectOption('#game', '王者荣耀');
    await page.click('text=立即支付');

    // 5. 验证
    await expect(page.locator('.order-status')).toContainText('待支付');
  });
});
代码审查流程

Pull Request 流程:
┌─────────────────────────────────────────────────────────────┐
│  1. 开发者提交 PR                                            │
│  ├─ 填写 PR 描述 (功能说明、改动文件、测试结果)             │
│  └─ 关联 Issue                                               │
│                                                              │
│  2. 自动检查                                                  │
│  ├─ ESLint 检查                                              │
│  ├─ Prettier 格式化                                          │
│  ├─ TypeScript 类型检查                                     │
│  ├─ 单元测试 (要求 80%+ 覆盖率)                            │
│  └─ 构建检查                                                │
│                                                              │
│  3. 代码审查 (至少 1 人 Approve)                              │
│  ├─ 代码质量检查                                            │
│  ├─ 性能影响检查                                            │
│  ├─ 安全漏洞检查                                            │
│  └─ 业务逻辑检查                                            │
│                                                              │
│  4. 合并主分支                                                  │
│  ├─ 解决冲突 (如有)                                          │
│  ├─ 再次运行测试                                            │
│  └─ 部署到测试环境                                          │
└─────────────────────────────────────────────────────────────┘
性能监控方案

监控指标:
┌─────────────────────────────────────────────────────────────┐
│  1. 前端性能                                                │
│  ├─ 首屏加载时间 (FCP)                                     │
│  ├─ 页面交互延迟 (FID)                                     │
│  ├─ 布局偏移 (CLS)                                          │
│  └─ 累积布局偏移 (LCP)                                   │
│                                                              │
│   2. 网络性能                                                │
│  ├─ API 响应时间                                              │
│  ├─ API 错误率                                              │
│  ├─ WebSocket 连接成功率                                    │
│  └─ 文件上传速度                                            │
│                                                              │
│  3. 用户行为                                                │
│  ├─ PV/UV                                                  │
│  ├─ 页面访问路径                                            │
│  ├– 功能使用频率                                            │
│  └─ 转化漏斗                                                │
│                                                              │
│  4. 错误监控                                                │
│  ├─ JavaScript 错误                                          │
│  ├─ API 错误                                                │
│  ├─ 白屏异常                                                │
│  └── 性能异常                                                │
└─────────────────────────────────────────────────────────────┘

监控工具:
├─ 微信小程序后台 (官方数据)
├─ 腾讯云前端性能监控 (如接入)
└─ 自建监控日志 (Sentry)
7️⃣ 实施计划总结
关键里程碑

Week 1: 项目初始化 ✅
├─ Taro 项目创建
├─ 技术栈配置
├─ 基础架构搭建
└─ CI/CD 配置

Week 3: 用户端核心功能 ✅
├─ 首页 + 搜索
├─ 陪玩师详情
├─ 下单流程
└─ 支付集成

Week 6: 陪玩师端核心功能 ✅
├─ 工作台 + 接单大厅
├─ 订单管理
└─ 收入管理

Week 8: 功能完整 ✅
├─ 营销模块
├─ 角色切换
└─ 完整测试

Week 10: 正式上线 ✅
├─ 性能优化
├─ 安全加固
└─ 提交审核
8️⃣ 总结
技术栈总结

推荐技术栈:
├── 框架: Taro 3.x + React/Vue3
├── 状态: Pinia / Zustand
├── UI: NutUI
├── 语言: TypeScript 5+
├── 构建: Vite 5.x
└── 测试: Vitest + Playwright

核心优势:
✅ 跨平台支持 (微信/支付宝/抖音/H5)
✅ 类型安全
✅ 开箱即用的组件库
✅ 完善的工具链
✅ 活跃的社区支持
开发规范总结

代码规范:
✅ 统一的命名规范
✅ 清晰的代码结构
✅ 完善的类型定义
✅ 详细的代码注释
✅ 严格的错误处理

质量保证:
✅ 单元测试 80%+
✅ E2E 测试覆盖核心流程
✅ 代码审查机制
✅ 性能监控方案
✅ 错误日志系统
时间规划总结

10周完整开发周期:
├── Week 1: 基础架构
├── Week 2-5: 用户端核心
├── Week 6-7: 陪玩师端核心
├── Week 8: 营销功能
└── Week 9-10: 测试与上线
📚 附录
推荐学习资源
官方文档:

Taro 官方文档
NutUI React 文档
Pinia 文档
社区资源:

Taro GitHub
Taro 社区论坛
Taro 插件市场
最佳实践:

Taro 最佳实践
小程序性能优化
文档版本: v1.0
最后更新: 2025-01-10
维护者: Super Dev Team

🚀 准备好开始小程序开发了吗？让我们开始吧！
