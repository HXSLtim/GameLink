/**
 * 翻译键接口
 */
export interface TranslationKeys {
  common: {
    confirm: string;
    cancel: string;
    save: string;
    delete: string;
    edit: string;
    add: string;
    search: string;
    reset: string;
    submit: string;
    back: string;
    loading: string;
    success: string;
    error: string;
    noData: string;
    viewAll: string;
    view: string;
    total: string;
    page: string;
    perPage: string;
  };
  auth: {
    login: string;
    logout: string;
    username: string;
    password: string;
    rememberMe: string;
    forgotPassword: string;
    loginSuccess: string;
    loginFailed: string;
    welcomeBack: string;
  };
  menu: {
    dashboard: string;
    users: string;
    orders: string;
    permissions: string;
    settings: string;
  };
  dashboard: {
    title: string;
    totalOrders: string;
    todayOrders: string;
    pendingReview: string;
    todayRevenue: string;
    quickActions: string;
    recentOrders: string;
    pendingReviews: string;
    allOrders: string;
    pendingReviewOrders: string;
    inProgressOrders: string;
    financialReports: string;
    viewAllOrders: string;
    processReviews: string;
    monitorProgress: string;
    viewReports: string;
    needsAction: string;
    vsLastMonth: string;
    vsYesterday: string;
  };
  orders: {
    title: string;
    orderNo: string;
    user: string;
    player: string;
    game: string;
    service: string;
    amount: string;
    duration: string;
    status: string;
    reviewStatus: string;
    createdAt: string;
    actions: string;
    searchPlaceholder: string;
    allStatus: string;
    allReviewStatus: string;
    allGames: string;
    allServices: string;
    waitingForPlayer: string;
    noReview: string;
    hours: string;
    rating: string;
    ordersCount: string;
  };
  orderStatus: {
    pending_payment: string;
    pending_accept: string;
    pending_review: string;
    in_review: string;
    review_approved: string;
    review_rejected: string;
    in_progress: string;
    completed: string;
    cancelled: string;
    refunded: string;
  };
  reviewStatus: {
    pending: string;
    in_review: string;
    approved: string;
    rejected: string;
  };
  gameType: {
    honor_of_kings: string;
    league_of_legends: string;
    peacekeeper_elite: string;
    genshin_impact: string;
    other: string;
  };
  serviceType: {
    accompany: string;
    boost: string;
    rank_up: string;
    entertainment: string;
  };
  theme: {
    light: string;
    dark: string;
    system: string;
    switchToLight: string;
    switchToDark: string;
  };
  error: {
    pageNotFound: string;
    serverError: string;
    networkError: string;
    unknownError: string;
    somethingWentWrong: string;
    backToHome: string;
    tryAgain: string;
  };
  time: {
    justNow: string;
    minutesAgo: string;
    hoursAgo: string;
    daysAgo: string;
  };
}

/**
 * 中文（简体）翻译
 */
export const zhCN: TranslationKeys = {
  common: {
    confirm: '确认',
    cancel: '取消',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    add: '添加',
    search: '搜索',
    reset: '重置',
    submit: '提交',
    back: '返回',
    loading: '加载中...',
    success: '操作成功',
    error: '操作失败',
    noData: '暂无数据',
    viewAll: '查看全部',
    view: '查看',
    total: '共',
    page: '页',
    perPage: '/ 页',
  },
  auth: {
    login: '登录',
    logout: '退出登录',
    username: '用户名',
    password: '密码',
    rememberMe: '记住我',
    forgotPassword: '忘记密码？',
    loginSuccess: '登录成功',
    loginFailed: '登录失败',
    welcomeBack: '欢迎回来',
  },
  menu: {
    dashboard: '仪表盘',
    users: '用户管理',
    orders: '订单管理',
    permissions: '权限管理',
    settings: '系统设置',
  },
  dashboard: {
    title: '仪表盘',
    totalOrders: '总订单数',
    todayOrders: '今日订单',
    pendingReview: '待审核',
    todayRevenue: '今日收入',
    quickActions: '快捷入口',
    recentOrders: '最近订单',
    pendingReviews: '待审核订单',
    allOrders: '所有订单',
    pendingReviewOrders: '待审核订单',
    inProgressOrders: '进行中订单',
    financialReports: '财务报表',
    viewAllOrders: '查看和管理所有订单',
    processReviews: '处理需要审核的订单',
    monitorProgress: '监控正在进行的订单',
    viewReports: '查看收入和统计数据',
    needsAction: '需要处理',
    vsLastMonth: '较上月',
    vsYesterday: '较昨日',
  },
  orders: {
    title: '订单管理',
    orderNo: '订单号',
    user: '用户',
    player: '陪玩者',
    game: '游戏',
    service: '服务',
    amount: '金额',
    duration: '时长',
    status: '订单状态',
    reviewStatus: '审核状态',
    createdAt: '创建时间',
    actions: '操作',
    searchPlaceholder: '搜索订单号、用户、陪玩者',
    allStatus: '全部状态',
    allReviewStatus: '全部审核状态',
    allGames: '全部游戏',
    allServices: '全部服务',
    waitingForPlayer: '待接单',
    noReview: '-',
    hours: '小时',
    rating: '评分',
    ordersCount: '单',
  },
  orderStatus: {
    pending_payment: '待支付',
    pending_accept: '待接单',
    pending_review: '待审核',
    in_review: '审核中',
    review_approved: '审核通过',
    review_rejected: '审核拒绝',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款',
  },
  reviewStatus: {
    pending: '待审核',
    in_review: '审核中',
    approved: '已通过',
    rejected: '已拒绝',
  },
  gameType: {
    honor_of_kings: '王者荣耀',
    league_of_legends: '英雄联盟',
    peacekeeper_elite: '和平精英',
    genshin_impact: '原神',
    other: '其他',
  },
  serviceType: {
    accompany: '陪玩',
    boost: '代练',
    rank_up: '上分',
    entertainment: '娱乐',
  },
  theme: {
    light: '亮色',
    dark: '暗色',
    system: '跟随系统',
    switchToLight: '切换到亮色模式',
    switchToDark: '切换到暗色模式',
  },
  error: {
    pageNotFound: '页面未找到',
    serverError: '服务器错误',
    networkError: '网络错误',
    unknownError: '未知错误',
    somethingWentWrong: '哎呀，出错了！',
    backToHome: '返回首页',
    tryAgain: '重试',
  },
  time: {
    justNow: '刚刚',
    minutesAgo: '分钟前',
    hoursAgo: '小时前',
    daysAgo: '天前',
  },
};
