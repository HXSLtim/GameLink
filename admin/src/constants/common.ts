/**
 * 通用常量配置
 * 集中管理应用中的魔法数字
 */

// ============================================================================
// 金额相关
// ============================================================================

export const MONEY = {
    /** 快捷充值金额选项（元） */
    QUICK_AMOUNTS: [50, 100, 200, 500, 1000] as const,

    /** 最小充值金额（分） */
    MIN_RECHARGE: 100,

    /** 最大充值金额（分） */
    MAX_RECHARGE: 100000,

    /** 元与分转换比率 */
    YUAN_TO_FEN: 100,

    /** 最小自定义充值金额（元） */
    MIN_CUSTOM_AMOUNT: 1,

    /** 最大自定义充值金额（元） */
    MAX_CUSTOM_AMOUNT: 10000,

    /** 优惠券默认减免金额（元） */
    DEFAULT_DEDUCT_AMOUNT: 10,

    /** 优惠券默认最低消费（元） */
    DEFAULT_MIN_AMOUNT: 0,

    /** 优惠券默认最大优惠（元） */
    DEFAULT_MAX_DISCOUNT: 50,

    /** 优惠券默认有效期（天） */
    DEFAULT_VALIDITY_DAYS: 30,

    /** 优惠券默认每人限领数量 */
    DEFAULT_PER_USER_LIMIT: 1,

    /** 优惠券默认发放总数 */
    DEFAULT_TOTAL_COUNT: 100,

    /** 优惠券最大有效期（天） */
    MAX_VALIDITY_DAYS: 365,

    /** 优惠券最大发放总数 */
    MAX_TOTAL_COUNT: 1000000,

    /** 优惠券最大每人限领数量 */
    MAX_PER_USER_LIMIT: 100,

    /** 折扣率最小值 */
    MIN_DISCOUNT_RATE: 0.1,

    /** 折扣率最大值 */
    MAX_DISCOUNT_RATE: 1,

    /** 折扣率步进值 */
    DISCOUNT_RATE_STEP: 0.1,

    /** 折扣显示转换倍数（0.9 -> 9折） */
    DISCOUNT_DISPLAY_MULTIPLIER: 10,

    /** 优惠券名称最大长度 */
    COUPON_NAME_MAX_LENGTH: 50,

    /** 优惠券描述最大长度 */
    COUPON_DESC_MAX_LENGTH: 200,
} as const;

// ============================================================================
// 分页相关
// ============================================================================

export const PAGINATION = {
    /** 默认页码大小 */
    DEFAULT_PAGE_SIZE: 10,

    /** 页码大小选项 */
    PAGE_SIZE_OPTIONS: [10, 20, 50, 100] as const,

    /** 默认当前页 */
    DEFAULT_CURRENT: 1,

    /** 批量操作时的最大加载量 */
    BATCH_LOAD_SIZE: 1000,

    /** 导出数据时的最大加载量 */
    EXPORT_LOAD_SIZE: 10000,
} as const;

// ============================================================================
// 布局相关
// ============================================================================

export const LAYOUT = {
    /** 标准 Grid 间距 */
    GUTTER: 16,

    /** 大间距 */
    GUTTER_LARGE: 24,

    /** 页面内边距 */
    PADDING: 24,

    /** 卡片间距 */
    CARD_MARGIN: 16,

    /** 列宽配置 */
    COL_SPAN: {
        FULL: 24,
        HALF: 12,
        THIRD: 8,
        QUARTER: 6,
    },

    /** 最小卡片高度 */
    MIN_CARD_HEIGHT: 120,

    /** 弹窗内边距 */
    MODAL_PADDING: 16,
} as const;

// ============================================================================
// 尺寸相关
// ============================================================================

export const SIZES = {
    /** 头像尺寸 */
    AVATAR: {
        SMALL: 32,
        MEDIUM: 40,
        LARGE: 64,
        XLARGE: 80,
        XXLARGE: 128,
    },

    /** 图标尺寸 */
    ICON: {
        SMALL: 16,
        MEDIUM: 20,
        LARGE: 24,
        XLARGE: 32,
    },

    /** 图片圆角 */
    IMAGE_BORDER_RADIUS: 8,

    /** 标签内边距 */
    TAG_PADDING: '4px 16px',

    /** 标签字体大小 */
    TAG_FONT_SIZE: 16,

    /** 次要文本字体大小 */
    SECONDARY_FONT_SIZE: 12,
} as const;

// ============================================================================
// 时间相关
// ============================================================================

export const TIMING = {
    /** 防抖延迟（毫秒） */
    DEBOUNCE: 300,

    /** 节流延迟（毫秒） */
    THROTTLE: 1000,

    /** Toast 显示时长（毫秒） */
    TOAST_DURATION: 3000,

    /** API 超时时间（毫秒） */
    API_TIMEOUT: 10000,

    /** 模拟加载延迟（毫秒） */
    MOCK_LOAD_DELAY: 500,

    /** 充值模拟延迟（毫秒） */
    RECHARGE_MOCK_DELAY: 1000,
} as const;

// ============================================================================
// 业务规则
// ============================================================================

export const BUSINESS = {
    /** 订单超时时间（分钟） */
    ORDER_TIMEOUT: 30,

    /** 争议处理时限（小时） */
    DISPUTE_DEADLINE: 24,

    /** 佣金比例（%） */
    COMMISSION: {
        MIN: 10,
        MAX: 25,
        DEFAULT: 20,
    },

    /** 表单行数配置 */
    FORM_ROWS: {
        DEFAULT: 2,
        TEXTAREA: 3,
    },

    /** InputNumber 精度 */
    PRECISION: {
        AMOUNT: 2,
        RATE: 2,
        PERCENT: 0,
    },

    /** 步进值 */
    STEP: {
        DEFAULT: 0.1,
    },

    /** 折扣率步进值 */
    DISCOUNT_RATE_STEP: 0.1,
} as const;

// ============================================================================
// 表格相关
// ============================================================================

export const TABLE = {
    /** 列宽配置 */
    COLUMN_WIDTH: {
        ID: 60,
        SMALL: 70,
        MEDIUM: 80,
        LARGE: 100,
        XLARGE: 120,
        XXLARGE: 150,
        XXXLARGE: 180,
        ACTION: 200,
        ICON: 70,
    },

    /** 列跨度配置 */
    COL_SPAN: LAYOUT.COL_SPAN,

    /** 滚动宽度配置 */
    SCROLL_WIDTH: {
        SMALL: 1400,
        MEDIUM: 1500,
        LARGE: 1600,
    },

    /** 技能标签最大显示数量 */
    SKILL_TAGS_MAX_DISPLAY: 3,
} as const;

// ============================================================================
// 弹窗相关
// ============================================================================

export const MODAL = {
    /** 弹窗宽度配置 */
    WIDTH: {
        SMALL: 480,
        MEDIUM: 550,
        LARGE: 600,
        XLARGE: 700,
        XXLARGE: 800,
    },

    /** 抽屉尺寸 */
    DRAWER_SIZE: {
        DEFAULT: 'large',
        SMALL: 'small',
    },
} as const;

// ============================================================================
// 字符串相关
// ============================================================================

export const TEXT = {
    /** 最大长度 */
    MAX_LENGTH: {
        SHORT: 50,
        MEDIUM: 200,
        LONG: 500,
    },

    /** 默认占位符 */
    PLACEHOLDER: {
        SEARCH: '请输入关键词',
        SELECT: '请选择',
    },
} as const;

// ============================================================================
// 颜色相关
// ============================================================================

export const COLORS = {
    /** 陪玩师头像背景色 */
    PLAYER_AVATAR_BG: '#722ed1',

    /** 金额颜色 */
    MONEY: '#f5222d',

    /** 成功颜色 */
    SUCCESS: '#52c41a',

    /** 警告颜色 */
    WARNING: '#faad14',

    /** 错误颜色 */
    ERROR: '#ff4d4f',

    /** 信息颜色 */
    INFO: '#1890ff',
} as const;

// ============================================================================
// 其他常量
// ============================================================================

export const OTHER = {
    /** 最小评分 */
    MIN_RATING: 0,

    /** 最大评分 */
    MAX_RATING: 5,

    /** 评分精度 */
    RATING_PRECISION: 1,

    /** 折扣率精度 */
    DISCOUNT_PRECISION: 1,

    /** Space 组件尺寸 */
    SPACE_SIZE: {
        ZERO: 0,
        SMALL: 4,
        MEDIUM: 'small',
        LARGE: 'small',
    },
} as const;
