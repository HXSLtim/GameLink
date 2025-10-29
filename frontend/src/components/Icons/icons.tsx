/**
 * SVG 图标库
 * 替换 emoji 使用的 SVG 图标组件
 */

import React from 'react';

export interface IconProps {
  size?: number;
  color?: string;
  className?: string;
}

// 支付方式图标
export const WechatPayIcon: React.FC<IconProps> = ({ size = 20, color = '#09BB07' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M9.5 8.5C9.5 9.88071 8.38071 11 7 11C5.61929 11 4.5 9.88071 4.5 8.5C4.5 7.11929 5.61929 6 7 6C8.38071 6 9.5 7.11929 9.5 8.5Z"
      fill={color}
    />
    <path
      d="M19.5 8.5C19.5 9.88071 18.3807 11 17 11C15.6193 11 14.5 9.88071 14.5 8.5C14.5 7.11929 15.6193 6 17 6C18.3807 6 19.5 7.11929 19.5 8.5Z"
      fill={color}
    />
    <path
      d="M12 16C15 16 18 14 19 12C19 12 18 11 12 11C6 11 5 12 5 12C6 14 9 16 12 16Z"
      fill={color}
    />
    <circle cx="12" cy="12" r="10" stroke={color} strokeWidth="2" />
  </svg>
);

export const AlipayIcon: React.FC<IconProps> = ({ size = 20, color = '#1677FF' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="3" y="3" width="18" height="18" rx="2" stroke={color} strokeWidth="2" />
    <path d="M7 12H17" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M12 7V17" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <circle cx="12" cy="12" r="6" stroke={color} strokeWidth="2" />
  </svg>
);

export const BalanceIcon: React.FC<IconProps> = ({ size = 20, color = '#FFB800' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <circle cx="12" cy="12" r="9" stroke={color} strokeWidth="2" />
    <path
      d="M12 6V7M12 17V18M15 9.5C15 8.67157 14.3284 8 13.5 8H10.5C9.67157 8 9 8.67157 9 9.5C9 10.3284 9.67157 11 10.5 11H13.5C14.3284 11 15 11.6716 15 12.5C15 13.3284 14.3284 14 13.5 14H10.5C9.67157 14 9 13.3284 9 12.5"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
    />
  </svg>
);

// 功能图标
export const CheckIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M20 6L9 17L4 12"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const CrossIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M18 6L6 18M6 6L18 18" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

export const SearchIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <circle cx="11" cy="11" r="6" stroke={color} strokeWidth="2" />
    <path d="M20 20L16 16" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

export const LocationIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M12 21C12 21 5 15 5 10C5 6.13401 8.13401 3 12 3C15.866 3 19 6.13401 19 10C19 15 12 21 12 21Z"
      stroke={color}
      strokeWidth="2"
    />
    <circle cx="12" cy="10" r="3" stroke={color} strokeWidth="2" />
  </svg>
);

// 实心星星
export const StarFilledIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor', className }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" className={className}>
    <path
      d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"
      fill={color}
    />
  </svg>
);

// 空心星星
export const StarOutlinedIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor', className }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" className={className}>
    <path
      d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"
      stroke={color}
      strokeWidth="2"
      fill="none"
      strokeLinejoin="round"
    />
  </svg>
);

// 兼容旧代码
export const StarIcon = StarFilledIcon;

export const DatabaseIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <ellipse cx="12" cy="5" rx="9" ry="3" stroke={color} strokeWidth="2" />
    <path d="M3 5V19C3 20.66 7.03 22 12 22C16.97 22 21 20.66 21 19V5" stroke={color} strokeWidth="2" />
    <path d="M3 12C3 13.66 7.03 15 12 15C16.97 15 21 13.66 21 12" stroke={color} strokeWidth="2" />
  </svg>
);

export const RefreshIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C14.8273 3 17.35 4.26284 19 6.23077"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
    />
    <path d="M21 3V7H17" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

export const ListIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M8 6H21" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M8 12H21" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M8 18H21" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <circle cx="3.5" cy="6" r="1.5" fill={color} />
    <circle cx="3.5" cy="12" r="1.5" fill={color} />
    <circle cx="3.5" cy="18" r="1.5" fill={color} />
  </svg>
);

export const BlockIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <circle cx="12" cy="12" r="10" stroke={color} strokeWidth="2" />
    <path d="M4 4L20 20" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

export const ChartIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M3 3V18C3 19.1046 3.89543 20 5 20H21" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M7 14L12 9L16 13L21 8" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

export const LightbulbIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M9 18H15M12 3V4M18.364 5.636L17.657 6.343M21 12H20M4 12H3M6.343 6.343L5.636 5.636M9 21H15M12 7C14.2091 7 16 8.79086 16 11C16 12.8638 14.7252 14.4299 13 14.874V17H11V14.874C9.27477 14.4299 8 12.8638 8 11C8 8.79086 9.79086 7 12 7Z"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const BoltIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" fill={color} stroke={color} strokeWidth="2" strokeLinejoin="round" />
  </svg>
);

export const TargetIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <circle cx="12" cy="12" r="10" stroke={color} strokeWidth="2" />
    <circle cx="12" cy="12" r="6" stroke={color} strokeWidth="2" />
    <circle cx="12" cy="12" r="2" fill={color} />
  </svg>
);

export const ToolIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

export const SparklesIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M12 3L13.5 8.5L19 10L13.5 11.5L12 17L10.5 11.5L5 10L10.5 8.5L12 3Z"
      fill={color}
      stroke={color}
      strokeWidth="2"
      strokeLinejoin="round"
    />
    <path d="M19 3L19.5 5L21.5 5.5L19.5 6L19 8L18.5 6L16.5 5.5L18.5 5L19 3Z" fill={color} />
    <path d="M19 16L19.5 18L21.5 18.5L19.5 19L19 21L18.5 19L16.5 18.5L18.5 18L19 16Z" fill={color} />
  </svg>
);

export const PinIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M16 3L18 5L12 11L14 13L12 15L8 11L10 9L4 3L5 2L11 8L17 2L16 3Z"
      stroke={color}
      strokeWidth="2"
      strokeLinejoin="round"
    />
  </svg>
);

export const LayersIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke={color} strokeWidth="2" strokeLinejoin="round" />
    <path d="M2 17L12 22L22 17" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M2 12L12 17L22 12" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

export const RulerIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="3" y="6" width="18" height="12" rx="2" stroke={color} strokeWidth="2" />
    <path d="M7 10V14" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M12 10V14" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M17 10V14" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 文档图标
export const FolderIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M3 7C3 5.89543 3.89543 5 5 5H9L11 7H19C20.1046 7 21 7.89543 21 9V17C21 18.1046 20.1046 19 19 19H5C3.89543 19 3 18.1046 3 17V7Z"
      stroke={color}
      strokeWidth="2"
    />
  </svg>
);

export const BookIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path
      d="M6.5 2H20V22H6.5A2.5 2.5 0 0 1 4 19.5V4.5A2.5 2.5 0 0 1 6.5 2Z"
      stroke={color}
      strokeWidth="2"
    />
  </svg>
);

export const TrendingUpIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path d="M23 6L13.5 15.5L8.5 10.5L1 18" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
    <path d="M17 6H23V12" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

export const LockIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="3" y="11" width="18" height="10" rx="2" stroke={color} strokeWidth="2" />
    <path
      d="M7 11V7C7 4.23858 9.23858 2 12 2C14.7614 2 17 4.23858 17 7V11"
      stroke={color}
      strokeWidth="2"
    />
    <circle cx="12" cy="16" r="1.5" fill={color} />
  </svg>
);

export const ArrowRightIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <path
      d="M5 12H19M19 12L12 5M19 12L12 19"
      stroke={color}
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
);

// 手机图标
export const PhoneIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="5" y="2" width="14" height="20" rx="2" stroke={color} strokeWidth="2" />
    <path d="M12 18H12.01" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 邮件图标
export const EmailIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="2" y="4" width="20" height="16" rx="2" stroke={color} strokeWidth="2" />
    <path d="M2 7L12 13L22 7" stroke={color} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
  </svg>
);

// 禁止/封禁图标
export const BanIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <circle cx="12" cy="12" r="10" stroke={color} strokeWidth="2" />
    <path d="M4.93 4.93L19.07 19.07" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 信用卡/支付图标
export const CreditCardIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="2" y="5" width="20" height="14" rx="2" stroke={color} strokeWidth="2" />
    <path d="M2 10H22" stroke={color} strokeWidth="2" />
    <path d="M6 15H10" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 剪贴板/订单列表图标
export const ClipboardIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="4" y="4" width="16" height="18" rx="2" stroke={color} strokeWidth="2" />
    <path d="M9 2H15V4H9V2Z" stroke={color} strokeWidth="2" strokeLinejoin="round" />
    <path d="M8 12H16" stroke={color} strokeWidth="2" strokeLinecap="round" />
    <path d="M8 16H16" stroke={color} strokeWidth="2" strokeLinecap="round" />
  </svg>
);

// 暂停图标
export const PauseIcon: React.FC<IconProps> = ({ size = 20, color = 'currentColor' }) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none">
    <rect x="6" y="4" width="4" height="16" rx="1" fill={color} />
    <rect x="14" y="4" width="4" height="16" rx="1" fill={color} />
  </svg>
);

// 导出所有图标的映射
export const PAYMENT_METHOD_ICONS = {
  wechat: WechatPayIcon,
  alipay: AlipayIcon,
  balance: BalanceIcon,
} as const;

export const FEATURE_ICONS = {
  check: CheckIcon,
  cross: CrossIcon,
  search: SearchIcon,
  location: LocationIcon,
  star: StarIcon,
  starFilled: StarFilledIcon,
  starOutlined: StarOutlinedIcon,
  database: DatabaseIcon,
  refresh: RefreshIcon,
  list: ListIcon,
  block: BlockIcon,
  chart: ChartIcon,
  lightbulb: LightbulbIcon,
  bolt: BoltIcon,
  target: TargetIcon,
  tool: ToolIcon,
  sparkles: SparklesIcon,
  pin: PinIcon,
  layers: LayersIcon,
  ruler: RulerIcon,
  folder: FolderIcon,
  book: BookIcon,
  trendingUp: TrendingUpIcon,
  lock: LockIcon,
  arrowRight: ArrowRightIcon,
  phone: PhoneIcon,
  email: EmailIcon,
  ban: BanIcon,
  creditCard: CreditCardIcon,
  clipboard: ClipboardIcon,
  pause: PauseIcon,
} as const;

