# 🎨 GameLink UI/UX 设计文档

## 📊 文档信息

| 项目 | 内容 |
|------|------|
| **文档名称** | UI/UX 设计规范文档 |
| **版本** | v1.0 |
| **创建日期** | 2025-12-05 |
| **作者** | Claude - 项目测试负责人 |
| **最后更新** | 2025-12-05 |

---

## 1. 设计原则

### 1.1 设计理念

**核心原则**：**简洁、专业、易用**

- **简洁**：界面清晰，信息层次明确，减少认知负担
- **专业**：适合游戏陪玩场景，体现游戏文化元素
- **易用**：操作流程顺畅，符合用户习惯，降低学习成本

### 1.2 设计目标

| 维度 | 目标 | 衡量标准 |
|------|------|----------|
| **用户体验** | 流畅自然的交互 | 用户满意度 > 4.5/5.0 |
| **视觉设计** | 美观一致的界面 | 符合品牌调性 |
| **可用性** | 直观易懂的界面 | 新用户上手时间 < 5分钟 |
| **响应式** | 多设备适配良好 | 移动端占比 > 50% |

### 1.3 色彩规范

**主色调**：
- **主色**：`#1890FF` (科技蓝) - 代表专业、可靠
- **辅助色**：`#722ED1` (紫罗兰) - 代表游戏、创新

**功能色彩**：
| 颜色 | 色值 | 用途 | 示例 |
|------|------|------|------|
| 成功色 | `#52C41A` | 成功状态、正向操作 | 支付成功、接单成功 |
| 警告色 | `#FAAD14` | 警告信息、注意操作 | 余额不足、即将到期 |
| 错误色 | `#F5222D` | 错误状态、危险操作 | 支付失败、封禁提示 |
| 信息色 | `#1890FF` | 友好提示、链接 | 查看详情、了解更多 |

**中性色**：
- 深色文字：`rgba(0, 0, 0, 0.85)`
- 次要文字：`rgba(0, 0, 0, 0.65)`
- 禁用文字：`rgba(0, 0, 0, 0.25)`
- 分割线：`rgba(0, 0, 0, 0.06)`
- 背景色：`#F0F2F5`

### 1.4 字体规范

**主字体**：
- 中文：`-apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif`
- 英文/数字：`'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif`

**字体层级**：

| 类型 | 大小 | 行高 | 字重 | 用途 |
|------|------|------|------|------|
| H1 | 38px | 1.23 | 600 | 页面主标题 |
| H2 | 30px | 1.35 | 600 | 区块标题 |
| H3 | 24px | 1.35 | 600 | 小标题 |
| 正文 | 14px | 1.5715 | 400 | 主要内容 |
| 次要 | 12px | 1.5 | 400 | 辅助信息 |

### 1.5 间距规范

**网格系统**：基于8px的网格系统

```css
/* 间距变量 */
:root {
  --space-xs: 4px;   /* 极小程序员 */
  --space-sm: 8px;   /* 小间距 */
  --space-md: 16px;  /* 中间距 */
  --space-lg: 24px;  /* 大间距 */
  --space-xl: 32px;  /* 极大间距 */
  --space-xxl: 48px; /* 超大间距 */
}
```

**组件间距**：
- 按钮内部：左右 16px，上下 8px
- 表单元素：垂直间距 24px
- 卡片内部：padding 24px
- 页面边距：左右 24px，上下 32px

---

## 2. 页面布局规范

### 2.1 布局框架

**管理后台布局**：

```
┌─────────────────────────────────────────────────────────────┐
│ 顶部导航栏 (Header)                                          │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Logo | 导航菜单 | 搜索 | 通知 | 用户头像 | 设置       │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ ┌──────────┬──────────────────────────────────────────────┐ │
│ │          │                                              │ │
│ │ 侧边菜单 │           内容区域 (Main)                    │ │
│ │          │                                              │ │
│ │  仪表盘  │                                              │ │
│ │  用户管理│                                              │ │
│ │  订单管理│           页面标题                           │ │
│ │  陪玩师  │                                              │ │
│ │  游戏管理│                                              │ │
│ │  系统设置│           表格/图表/表单                     │ │
│ │          │                                              │ │
│ │          │                                              │ │
│ │  折叠/展 │           分页                               │ │
│ │  开按钮  │                                              │ │
│ │          │                                              │ │
│ └──────────┴──────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ 底部版权信息 (Footer)                                        │
└─────────────────────────────────────────────────────────────┘
```

**用户端布局**：

```
┌─────────────────────────────────────────────────────────────┐
│ 顶部导航栏 (Header)                                          │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Logo | 游戏导航 | 搜索 | 登录/注册 | 个人中心        │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│            Hero区域 (大屏幕展示)                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                                                      │  │
│  │            主标题 + 副标题 + CTA按钮                │  │
│  │                                                      │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   特色功能区域 (Feature Section)                            │
│  ┌──────────┬──────────┬──────────┬──────────┐            │
│  │  特点1   │  特点2   │  特点3   │  特点4   │            │
│  └──────────┴──────────┴──────────┴──────────┘            │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   热门游戏区域 (Hot Games)                                  │
│  ┌─────────┬─────────┬─────────┬─────────┬─────────┐       │
│  │ 游戏卡片│ 游戏卡片│ 游戏卡片│ 游戏卡片│ 游戏卡片│       │
│  └─────────┴─────────┴─────────┴─────────┴─────────┘       │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│ 统计数据区域 (Statistics)                                   │
│  50,000+ 活跃用户 | 5,000+ 认证打手 | 50+ 支持游戏          │
├─────────────────────────────────────────────────────────────┤
│ 底部区域 (Footer)                                            │
│ 关于我们 | 服务条款 | 隐私政策 | 联系我们 | 帮助中心       │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 响应式断点

**断点设置**：

| 设备 | 断点 | 布局调整 |
|------|------|----------|
| 超小屏 | < 576px | 单列布局，简化导航 |
| 小屏 | ≥ 576px | 两列布局，抽屉式菜单 |
| 中屏 | ≥ 768px | 三列布局，侧边栏固定 |
| 大屏 | ≥ 992px | 四列布局，完整侧边栏 |
| 超大屏 | ≥ 1200px | 最大宽度限制，居中显示 |

**媒体查询示例**：

```less
// 响应式混合宏
.mobile-only {
  @media (max-width: 768px) {
    display: block;
  }
  @media (min-width: 769px) {
    display: none;
  }
}

.desktop-only {
  @media (max-width: 768px) {
    display: none;
  }
  @media (min-width: 769px) {
    display: block;
  }
}

// 容器响应式
.container {
  width: 100%;
  padding: 0 24px;
  margin: 0 auto;

  @media (min-width: 576px) {
    max-width: 540px;
  }

  @media (min-width: 768px) {
    max-width: 720px;
  }

  @media (min-width: 992px) {
    max-width: 960px;
  }

  @media (min-width: 1200px) {
    max-width: 1140px;
  }

  @media (min-width: 1600px) {
    max-width: 1320px;
  }
}
```

---

## 3. 组件设计规范

### 3.1 按钮 (Button)

**按钮类型**：

| 类型 | 样式 | 用途 | 示例 |
|------|------|------|------|
| 主要按钮 | 蓝色背景，白色文字 | 主要操作 | 立即下单、确认支付 |
| 次要按钮 | 白色背景，灰色边框 | 次要操作 | 取消、返回 |
| 危险按钮 | 红色背景 | 删除、警告操作 | 删除订单、封禁用户 |
| 链接按钮 | 无背景，蓝色文字 | 链接跳转 | 查看详情、忘记密码 |

**按钮尺寸**：

```less
.btn {
  // 小型按钮
  &-small {
    height: 24px;
    padding: 0 12px;
    font-size: 12px;
  }

  // 默认按钮
  &-default {
    height: 32px;
    padding: 0 16px;
    font-size: 14px;
  }

  // 大型按钮
  &-large {
    height: 40px;
    padding: 0 24px;
    font-size: 16px;
  }
}
```

**按钮状态**：

- **默认状态**：背景 `#1890FF`，文字白色
- **悬停状态**：背景 `#40A9FF`，轻微阴影
- **激活状态**：背景 `#096DD9`，无阴影
- **禁用状态**：背景 `#D9D9D9`，文字 `#BFBFBF`
- **加载状态**：显示加载动画

### 3.2 表单 (Form)

**输入框 (Input)**：

```less
.input {
  height: 32px;
  padding: 0 12px;
  border: 1px solid #D9D9D9;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.3s;

  &:hover {
    border-color: #40A9FF;
  }

  &:focus {
    border-color: #1890FF;
    box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
    outline: none;
  }

  &:disabled {
    background-color: #F5F5F5;
    color: #BFBFBF;
  }

  // 错误状态
  &.error {
    border-color: #F5222D;

    &:focus {
      box-shadow: 0 0 0 2px rgba(245, 34, 45, 0.2);
    }
  }

  // 成功状态
  &.success {
    border-color: #52C41A;

    &:focus {
      box-shadow: 0 0 0 2px rgba(82, 196, 26, 0.2);
    }
  }
}

// 标签
.input-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);

  .required {
    color: #F5222D;
    margin-left: 4px;
  }
}

// 错误提示
.input-error {
  margin-top: 4px;
  color: #F5222D;
  font-size: 12px;
}

// 帮助文本
.input-help {
  margin-top: 4px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
}
```

**下拉选择框 (Select)**：

```tsx
// Ant Design Select配置
<Select
  style={{ width: '100%' }}
  placeholder="请选择游戏类型"
  allowClear
  showSearch
  filterOption={(input, option) =>
    option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
  }
>
  <Option value="moba">MOBA</Option>
  <Option value="fps">FPS</Option>
  <Option value="rpg">RPG</Option>
</Select>
```

**多行文本框 (TextArea)**：

```tsx
<TextArea
  rows={4}
  maxLength={500}
  showCount
  placeholder="请输入服务要求描述"
/>
```

### 3.3 表格 (Table)

**标准表格**：

```tsx
// Ant Design Table配置
<Table
  columns={[
    {
      title: '订单编号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      width: 180,
      fixed: 'left',
    },
    {
      title: '用户信息',
      dataIndex: 'userName',
      key: 'userName',
      width: 120,
    },
    {
      title: '陪玩师',
      dataIndex: 'playerName',
      key: 'playerName',
      width: 120,
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      width: 100,
      align: 'right',
      render: (amount) => `¥${(amount / 100).toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (date) => formatDate(date),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (record) => (
        <Space>
          <Button type="link" onClick={() => viewDetail(record)}>
            查看
          </Button>
          <Button type="link" onClick={() => editOrder(record)}>
            编辑
          </Button>
        </Space>
      ),
    },
  ]}
  dataSource={orderList}
  loading={loading}
  rowKey="id"
  scroll={{ x: 1200 }}
  pagination={{
    total: total,
    pageSize: 20,
    showSizeChanger: true,
    showQuickJumper: true,
    showTotal: (total) => `共 ${total} 条记录`,
  }}
/>
```

**分页器**：

```tsx
<Pagination
  total={500}
  pageSize={20}
  showSizeChanger
  showQuickJumper
  showTotal={(total) => `共 ${total} 条`}
  onChange={(page, pageSize) => handlePageChange(page, pageSize)}
/>
```

### 3.4 卡片 (Card)

**信息卡片**：

```tsx
<Card
  hoverable
  cover={
    <img
      alt="游戏封面"
      src="https://example.com/game-cover.jpg"
    />
  }
  actions={[
    <Button type="primary">立即预约</Button>,
    <Button>查看详情</Button>,
  ]}
>
  <Card.Meta
    title="英雄联盟LOL代练"
    description="专业打手带你上分，快速晋升"
  />
  <div className="card-content">
    <div className="price">¥ 50/小时</div>
    <div className="stats">
      <span>⭐ 4.8分</span>
      <span>👥 1.2k单</span>
    </div>
  </div>
</Card>
```

**统计数据卡片**：

```tsx
<Card className="stat-card">
  <Statistic
    title="今日收入"
    value={12580}
    prefix="¥"
    valueStyle={{ color: '#52C41A' }}
    formatter={(value) => (value / 100).toFixed(2)}
  />
  <div className="trend up">
    <ArrowUpOutlined />
    <span>12.5%</span>
    <span>较昨日</span>
  </div>
</Card>
```

### 3.5 提示与反馈

**消息提示 (Message)**：

```tsx
// 成功提示
message.success('订单创建成功！')

// 错误提示
message.error('支付失败，请重试')

// 警告提示
message.warning('余额不足，请及时充值')

// 信息提示
message.info('订单已提交，等待接单')

// 加载提示
const hide = message.loading('正在处理中...', 0)
setTimeout(hide, 2500)
```

**通知提醒 (Notification)**：

```tsx
notification.success({
  message: '支付成功',
  description: '您的订单已支付成功，正在等待陪玩师接单',
  duration: 5,
  onClick: () => {
    navigate('/orders')
  },
})

notification.error({
  message: '接单失败',
  description: '订单已被其他陪玩师接走，请尝试其他订单',
  duration: 5,
  action: <Button onClick={() => navigate('/order-pool')}>查看其他订单</Button>,
})
```

**模态框 (Modal)**：

```tsx
<Modal
  title="确认接单"
  visible={visible}
  onOk={handleOk}
  onCancel={handleCancel}
  okText="确认接单"
  cancelText="取消"
  centered
>
  <div className="order-confirm-info">
    <p><strong>订单编号：</strong>202512050001</p>
    <p><strong>游戏类型：</strong>英雄联盟</p>
    <p><strong>服务时长：</strong>2小时</p>
    <p><strong>预估收入：</strong>¥80.00</p>
    <p><strong>客户要求：</strong>希望上分到钻石段位</p>
  </div>
</Modal>
```

---

## 4. 页面设计规范

### 4.1 登录注册页

**登录页布局**：

```tsx
// 页面结构
<LoginPage>
  <div className="login-container">
    {/* 左侧：品牌信息 */}
    <div className="login-brand">
      <Logo size="large" />
      <Title level={2}>欢迎来到 GameLink</Title>
      <Paragraph>
        专业游戏陪玩平台，连接每一位游戏玩家
      </Paragraph>
      <div className="features">
        <FeatureItem icon="safety" text="安全可靠" />
        <FeatureItem icon="team" text="专业打手" />
        <FeatureItem icon="clock" text="快速匹配" />
      </div>
    </div>

    {/* 右侧：登录表单 */}
    <div className="login-form-container">
      <Card className="login-card">
        <Tabs centered>
          <TabPane tab="账号登录" key="account">
            <Form onFinish={handleLogin}>
              <Form.Item name="username" rules={[{ required: true }]}>
                <Input prefix={<UserOutlined />} placeholder="用户名/手机号/邮箱" />
              </Form.Item>
              <Form.Item name="password" rules={[{ required: true }]}>
                <Input.Password prefix={<LockOutlined />} placeholder="密码" />
              </Form.Item>
              <Form.Item>
                <Checkbox>记住我</Checkbox>
                <a className="float-right" href="/forgot-password">忘记密码？</a>
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" block loading={loading}>
                  登录
                </Button>
              </Form.Item>
            </Form>
          </TabPane>
          <TabPane tab="短信登录" key="sms">
            {/* 短信登录表单 */}
          </TabPane>
        </Tabs>
        <Divider>或使用以下方式登录</Divider>
        <div className="social-login">
          <WechatOutlined className="wechat-icon" />
          <QqOutlined className="qq-icon" />
        </div>
        <div className="register-link">
          还没有账号？ <a href="/register">立即注册</a>
        </div>
      </Card>
    </div>
  </div>
</LoginPage>
```

**注册页布局**：

```tsx
<RegisterPage>
  <Card className="register-card">
    <Title level={3}>创建账户</Title>
    <Steps current={currentStep}>
      <Step title="填写信息" />
      <Step title="验证手机" />
      <Step title="完成注册" />
    </Steps>

    <Form form={form} onFinish={handleRegister}>
      {/* 第一步：基本信息 */}
      {currentStep === 0 && (
        <>
          <Form.Item name="username" rules={[{ required: true, min: 3, max: 20 }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名（3-20位）" />
          </Form.Item>
          <Form.Item name="phone" rules={[{ required: true, pattern: /^1[3-9]\d{9}$/ }]}>
            <Input prefix={<MobileOutlined />} placeholder="手机号" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, min: 8 }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码（8位以上）" />
          </Form.Item>
          <Form.Item name="confirmPassword" rules={[{ required: true }, ({ getFieldValue }) => ({
            validator(_, value) {
              if (!value || getFieldValue('password') === value) {
                return Promise.resolve();
              }
              return Promise.reject(new Error('两次密码输入不一致'));
            },
          })]}>
            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
          </Form.Item>
        </>
      )}

      {/* 第二步：短信验证 */}
      {currentStep === 1 && (
        <Form.Item name="smsCode" rules={[{ required: true, len: 6 }]}>
          <Input.Search
            placeholder="短信验证码"
            enterButton="获取验证码"
            onSearch={sendSmsCode}
          />
        </Form.Item>
      )}

      {/* 第三步：完成 */}
      {currentStep === 2 && (
        <Result
          status="success"
          title="注册成功"
          subTitle="恭喜您成功创建账户，即将跳转到登录页面"
        />
      )}
    </Form>

    <div className="steps-action">
      {currentStep > 0 && (
        <Button style={{ margin: '0 8px' }} onClick={() => prev()}>
          上一步
        </Button>
      )}
      {currentStep < 2 && (
        <Button type="primary" onClick={() => next()}>
          下一步
        </Button>
      )}
      {currentStep === 2 && (
        <Button type="primary" onClick={() => navigate('/login')}>
          立即登录
        </Button>
      )}
    </div>
  </Card>
</RegisterPage>
```

### 4.2 订单列表页

```tsx
<OrderListPage>
  {/* 搜索和筛选 */}
  <Card className="search-card">
    <Form layout="inline" onFinish={handleSearch}>
      <Form.Item name="orderNo">
        <Input placeholder="订单编号" />
      </Form.Item>
      <Form.Item name="status">
        <Select placeholder="订单状态" style={{ width: 120 }} allowClear>
          <Option value="pending">待支付</Option>
          <Option value="confirmed">已确认</Option>
          <Option value="in_progress">进行中</Option>
          <Option value="completed">已完成</Option>
        </Select>
      </Form.Item>
      <Form.Item name="dateRange">
        <RangePicker />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">搜索</Button>
        <Button onClick={() => form.resetFields()}>重置</Button>
      </Form.Item>
    </Form>
  </Card>

  {/* 统计卡片 */}
  <Row gutter={16} className="stats-row">
    <Col span={6}>
      <Card className="stat-card">
        <Statistic title="总订单" value={1250} />
      </Card>
    </Col>
    <Col span={6}>
      <Card className="stat-card">
        <Statistic title="进行中" value={45} valueStyle={{ color: '#1890FF' }} />
      </Card>
    </Col>
    <Col span={6}>
      <Card className="stat-card">
        <Statistic title="已完成" value={1150} valueStyle={{ color: '#52C41A' }} />
      </Card>
    </Col>
    <Col span={6}>
      <Card className="stat-card">
        <Statistic title="总收入" value={125800} prefix="¥" valueStyle={{ color: '#722ED1' }} />
      </Card>
    </Col>
  </Row>

  {/* 订单表格 */}
  <Card className="table-card">
    <div className="table-header">
      <Title level={4}>订单列表</Title>
      <Button type="primary" icon={<PlusOutlined />}>创建订单</Button>
    </div>
    <Table
      columns={columns}
      dataSource={orderList}
      loading={loading}
      rowKey="id"
      pagination={pagination}
    />
  </Card>
</OrderListPage>
```

### 4.3 陪玩师详情页

```tsx
<PlayerDetailPage>
  {/* 顶部信息区域 */}
  <Card className="player-header">
    <div className="player-info">
      <Avatar size={120} src={player.avatarUrl} />
      <div className="info-content">
        <div className="name-verify">
          <Title level={3}>{player.nickname}</Title>
          {player.verified && <CheckCircleFilled className="verify-icon" />}
        </div>
        <div className="stats">
          <div className="stat-item">
            <StarOutlined />
            <span>{player.rating}</span>
          </div>
          <div className="stat-item">
            <HistoryOutlined />
            <span>{player.orderCount} 单</span>
          </div>
          <div className="stat-item">
            <DollarOutlined />
            <span>{player.price}/小时</span>
          </div>
        </div>
        <Paragraph className="bio">{player.bio}</Paragraph>
        <div className="skills">
          {player.skills.map(skill => (
            <Tag key={skill.id} color="blue">{skill.name}</Tag>
          ))}
        </div>
        <div className="actions">
          <Button type="primary" size="large" onClick={createOrder}>
            立即预约
          </Button>
          <Button size="large" onClick={sendGift}>
            赠送礼物
          </Button>
          <Button size="large">关注</Button>
        </div>
      </div>
    </div>
  </Card>

  {/* 服务内容区域 */}
  <Card title="服务内容" className="services-section">
    <Tabs>
      <TabPane tab="护航服务" key="escort">
        <List
          dataSource={player.services}
          renderItem={service => (
            <List.Item
              actions={[
                <span className="price">¥{service.price}/小时</span>,
                <Button type="primary" onClick={() => selectService(service)}>选择</Button>
              ]}
            >
              <List.Item.Meta
                title={service.name}
                description={service.description}
              />
            </List.Item>
          )}
        />
      </TabPane>
      <TabPane tab="评价" key="reviews">
        <List
          dataSource={player.reviews}
          renderItem={review => (
            <List.Item>
              <Comment
                author={review.userName}
                avatar={<Avatar src={review.userAvatar} />}
                content={review.content}
                datetime={review.createdAt}
              />
            </List.Item>
          )}
        />
      </TabPane>
    </Tabs>
  </Card>
</PlayerDetailPage>
```

### 4.4 仪表盘 (Dashboard)

```tsx
<DashboardPage>
  {/* 统计卡片 */}
  <Row gutter={16}>
    <Col span={6}>
      <Card>
        <Statistic
          title="总用户"
          value={50000}
          valueStyle={{ color: '#1890FF' }}
          prefix={<UserOutlined />}
          suffix="人"
        />
      </Card>
    </Col>
    <Col span={6}>
      <Card>
        <Statistic
          title="总订单"
          value={12580}
          valueStyle={{ color: '#52C41A' }}
          prefix={<FileTextOutlined />}
          suffix="单"
        />
      </Card>
    </Col>
    <Col span={6}>
      <Card>
        <Statistic
          title="总收入"
          value={12580000}
          valueStyle={{ color: '#722ED1' }}
          prefix={<DollarOutlined />}
          precision={2}
          formatter={value => `¥${(value / 100).toFixed(2)}`}
        />
      </Card>
    </Col>
    <Col span={6}>
      <Card>
        <Statistic
          title="今日新增"
          value={125}
          valueStyle={{ color: '#FAAD14' }}
          prefix={<ArrowUpOutlined />}
          suffix="人"
        />
      </Card>
    </Col>
  </Row>

  {/* 图表区域 */}
  <Row gutter={16} style={{ marginTop: 24 }}>
    <Col span={16}>
      <Card title="收入趋势">
        <Line {...revenueChartConfig} />
      </Card>
    </Col>
    <Col span={8}>
      <Card title="订单状态分布">
        <Pie {...orderStatusChartConfig} />
      </Card>
    </Col>
  </Row>

  {/* 最新订单 */}
  <Row style={{ marginTop: 24 }}>
    <Col span={24}>
      <Card title="最新订单">
        <Table
          columns={columns}
          dataSource={recentOrders}
          pagination={false}
          rowKey="id"
          size="middle"
        />
      </Card>
    </Col>
  </Row>
</DashboardPage>
```

---

## 5. 动效与交互

### 5.1 加载动画

**全局加载**：

```tsx
// 页面加载
const [loading, setLoading] = useState(true)
useEffect(() => {
  fetchData().then(() => setLoading(false))
}, [])

return (
  <Spin spinning={loading} tip="加载中...">
    <PageContent>
      {/* 页面内容 */}
    </PageContent>
  </Spin>
)
```

**局部加载**：

```tsx
const LoadMoreList = () => {
  const [loading, setLoading] = useState(false)

  const onLoadMore = () => {
    setLoading(true)
    loadMore().then(() => setLoading(false))
  }

  return (
    <>
      <List
        loading={loading}
        loadMore={
          !loading && (
            <div style={{ textAlign: 'center', marginTop: 12, height: 32, lineHeight: '32px' }}>
              <Button onClick={onLoadMore}>加载更多</Button>
            </div>
          )
        }
        dataSource={list}
        renderItem={item => <List.Item>{item.name}</List.Item>}
      />
    </>
  )
}
```

**骨架屏**：

```tsx
<Skeleton active avatar paragraph={{ rows: 4 }} />
```

### 5.2 微交互

**悬停效果**：

```less
.hover-card {
  transition: all 0.3s ease;
  cursor: pointer;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }
}

.hover-button {
  transition: all 0.2s ease;

  &:hover {
    transform: scale(1.05);
  }

  &:active {
    transform: scale(0.95);
  }
}
```

**点击涟漪效果**：

```tsx
// Ant Design Button自带涟漪效果
<Button type="primary">点击我</Button>

// 自定义涟漪
const RippleButton: React.FC<ButtonProps> = (props) => {
  const [ripples, setRipples] = useState<Array<{ x: number; y: number; size: number }>>([])

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    const rect = e.currentTarget.getBoundingClientRect()
    const size = Math.max(rect.width, rect.height)
    const x = e.clientX - rect.left - size / 2
    const y = e.clientY - rect.top - size / 2

    setRipples([...ripples, { x, y, size }])

    setTimeout(() => {
      setRipples([])
    }, 600)

    props.onClick?.(e)
  }

  return (
    <button {...props} onClick={handleClick} className="ripple-button">
      {props.children}
      {ripples.map((ripple, index) => (
        <span key={index} className="ripple" style={{
          left: ripple.x,
          top: ripple.y,
          width: ripple.size,
          height: ripple.size
        }} />
      ))}
    </button>
  )
}
```

### 5.3 页面切换动画

**路由切换**：

```tsx
import { CSSTransition, TransitionGroup } from 'react-transition-group'
import { useLocation } from 'react-router-dom'

const AnimatedSwitch = () => {
  const location = useLocation()

  return (
    <TransitionGroup>
      <CSSTransition key={location.pathname} classNames="fade" timeout={300}>
        <Switch location={location}>
          <Route path="/home" component={HomePage} />
          <Route path="/orders" component={OrderPage} />
          <Route path="/profile" component={ProfilePage} />
        </Switch>
      </CSSTransition>
    </TransitionGroup>
  )
}
```

```less
.fade-enter {
  opacity: 0;
  transform: translateX(100%);
}

.fade-enter-active {
  opacity: 1;
  transform: translateX(0);
  transition: opacity 300ms, transform 300ms;
}

.fade-exit {
  opacity: 1;
  transform: translateX(0);
}

.fade-exit-active {
  opacity: 0;
  transform: translateX(-100%);
  transition: opacity 300ms, transform 300ms;
}
```

---

## 6. 移动端适配

### 6.1 移动端导航

**底部导航栏**：

```tsx
// 移动端底部导航
<nav className="mobile-nav">
  <NavLink to="/" exact activeClassName="active">
    <HomeOutlined />
    <span>首页</span>
  </NavLink>
  <NavLink to="/orders" activeClassName="active">
    <FileTextOutlined />
    <span>订单</span>
  </NavLink>
  <NavLink to="/messages" activeClassName="active">
    <MessageOutlined />
    <span>消息</span>
  </NavLink>
  <NavLink to="/profile" activeClassName="active">
    <UserOutlined />
    <span>我的</span>
  </NavLink>
</nav>
```

```less
.mobile-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  background: #fff;
  border-top: 1px solid #f0f0f0;
  display: flex;
  z-index: 1000;

  a {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    font-size: 10px;
    color: #999;
    text-decoration: none;

    &.active {
      color: #1890FF;
    }

    .anticon {
      font-size: 20px;
      margin-bottom: 2px;
    }
  }
}

// 为底部导航留出空间
.mobile-content {
  padding-bottom: 56px;
}
```

### 6.2 触摸优化

**触摸区域**：

```css
/* 最小触摸区域 44x44px */
.touchable {
  min-width: 44px;
  min-height: 44px;
  padding: 12px;
}

/* 防止双击缩放 */
.prevent-double-tap {
  touch-action: manipulation;
}

/* 滚动优化 */
.scroll-container {
  -webkit-overflow-scrolling: touch;
  overflow-y: auto;
}
```

**手势操作**：

```tsx
import { useSwipeable } from 'react-swipeable'

const SwipeableCard = ({ onSwipeLeft, onSwipeRight, children }) => {
  const handlers = useSwipeable({
    onSwipedLeft: () => onSwipeLeft && onSwipeLeft(),
    onSwipedRight: () => onSwipeRight && onSwipeRight(),
    trackMouse: true,
    delta: 10,
  })

  return <div {...handlers}>{children}</div>
}
```

---

## 7. 无障碍设计

### 7.1 语义化HTML

```tsx
// 正确的表单标签
<label htmlFor="username">用户名</label>
<input id="username" type="text" />

// ARIA属性
<button aria-label="关闭模态框" onClick={closeModal}>
  <CloseOutlined />
</button>

<div role="alert" aria-live="assertive">
  {errorMessage}
</div>

// 键盘导航
<div tabIndex={0} onKeyDown={handleKeyDown}>
  可聚焦元素
</div>
```

### 7.2 颜色对比

**对比度要求（WCAG 2.1 AA标准）**：

- 普通文本：≥ 4.5:1
- 大文本（18pt+）：≥ 3:1
- 非文本元素：≥ 3:1

```less
// 检查对比度
.text-primary {
  color: #1890FF; // 背景白色 #FFFFFF
  // 对比度: 3.18:1 (大文本可用)
}

.text-secondary {
  color: #8C8C8C; // 背景白色 #FFFFFF
  // 对比度: 4.54:1 (符合要求)
}
```

---

## 8. 设计资源

### 8.1 设计系统组件

**基础组件库**：
- [Ant Design](https://ant.design/) - UI组件库
- [Ant Design Icons](https://ant.design/components/icon/) - 图标库
- [Ant Design Charts](https://charts.ant.design/) - 图表库

### 8.2 图标规范

**使用Ant Design图标**：

```tsx
import {
  UserOutlined,
  LockOutlined,
  MobileOutlined,
  HomeOutlined,
  FileTextOutlined,
  MessageOutlined,
  DollarOutlined,
  StarOutlined,
  HistoryOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ExclamationCircleFilled,
  InfoCircleFilled,
} from '@ant-design/icons'
```

**图标映射关系**：

| 功能 | 图标 | 说明 |
|------|------|------|
| 用户 | UserOutlined | 个人中心、用户信息 |
| 密码 | LockOutlined | 密码输入、安全设置 |
| 手机 | MobileOutlined | 手机号、短信验证 |
| 订单 | FileTextOutlined | 订单列表、订单详情 |
| 消息 | MessageOutlined | 消息中心、通知 |
| 金钱 | DollarOutlined | 价格、收入、钱包 |
| 星级 | StarOutlined | 评分、收藏 |
| 历史 | HistoryOutlined | 历史记录 |
| 成功 | CheckCircleFilled | 成功状态 |
| 错误 | CloseCircleFilled | 错误状态 |
| 警告 | ExclamationCircleFilled | 警告状态 |
| 信息 | InfoCircleFilled | 信息提示 |

---

## 9. 开发规范

### 9.1 CSS命名规范

**BEM命名法**：

```tsx
// Block - Element - Modifier
<div className="order-card order-card--completed">
  <div className="order-card__header">
    <h3 className="order-card__title">订单标题</h3>
    <span className="order-card__status">已完成</span>
  </div>
  <div className="order-card__body">
    <p className="order-card__description">订单描述</p>
  </div>
  <div className="order-card__footer">
    <button className="order-card__button order-card__button--primary">查看详情</button>
  </div>
</div>
```

### 9.2 组件文件结构

```
Component/
├── index.tsx          # 主组件
├── index.less         # 样式（使用CSS Modules）
├── types.ts           # TypeScript类型定义
├── utils.ts           # 工具函数
├── hooks.ts           # 自定义Hook
├── constants.ts       # 常量定义
└── __tests__/         # 测试文件
    ├── index.test.tsx
    └── mock.ts
```

### 9.3 状态管理

**使用Context**：

```tsx
// UserContext.tsx
import { createContext, useContext, useReducer } from 'react'

const UserContext = createContext()

const userReducer = (state, action) => {
  switch (action.type) {
    case 'LOGIN':
      return { ...state, user: action.payload, isAuthenticated: true }
    case 'LOGOUT':
      return { ...state, user: null, isAuthenticated: false }
    default:
      return state
  }
}

export const UserProvider = ({ children }) => {
  const [state, dispatch] = useReducer(userReducer, {
    user: null,
    isAuthenticated: false,
  })

  return (
    <UserContext.Provider value={{ state, dispatch }}>
      {children}
    </UserContext.Provider>
  )
}

export const useUser = () => useContext(UserContext)
```

---

## 10. 设计检查清单

### 10.1 视觉设计检查

- [ ] 色彩使用是否符合规范
- [ ] 字体大小层级清晰
- [ ] 间距符合网格系统
- [ ] 图标风格统一
- [ ] 图片质量高清
- [ ] 加载状态完善
- [ ] 空状态友好

### 10.2 交互设计检查

- [ ] 按钮交互反馈明确
- [ ] 表单验证及时
- [ ] 错误提示清晰
- [ ] 页面切换流畅
- [ ] 手势操作自然
- [ ] 键盘导航完整
- [ ] 无障碍支持

### 10.3 响应式设计检查

- [ ] 移动端布局合理
- [ ] 触摸区域足够大
- [ ] 文字在小屏可读
- [ ] 横竖屏切换正常
- [ ] 加载速度优化
- [ ] 图片适配屏幕

---

## 11. 设计工具

### 11.1 推荐工具

- **Figma** - 界面设计、原型制作
- **MasterGo** - 团队协作设计
- **Iconfont** - 图标管理
- **Ant Design** - React组件库
- **Unsplash** - 高质量图片

### 11.2 设计资源

**Ant Design资源**：
- [官方组件库](https://ant.design/components/overview/)
- [设计规范](https://ant.design/docs/spec/introduce)
- [图标库](https://ant.design/components/icon/)

**配色工具**：
- [Ant Design Colors](https://ant.design/docs/spec/colors)

---

**文档版本历史**

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.0 | 2025-12-05 | Claude | 初始版本，包含完整UI/UX设计规范 |

---

<div align="center">

**🎨 好的设计让产品更有温度**

</div>
