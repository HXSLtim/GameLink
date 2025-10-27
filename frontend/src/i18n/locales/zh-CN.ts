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
  };
  menu: {
    dashboard: string;
    users: string;
    orders: string;
    permissions: string;
  };
  error: {
    pageNotFound: string;
    serverError: string;
    networkError: string;
    unknownError: string;
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
  },
  menu: {
    dashboard: '仪表盘',
    users: '用户管理',
    orders: '订单管理',
    permissions: '权限管理',
  },
  error: {
    pageNotFound: '页面未找到',
    serverError: '服务器错误',
    networkError: '网络错误',
    unknownError: '未知错误',
  },
};

