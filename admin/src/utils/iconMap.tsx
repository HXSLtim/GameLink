import * as Icons from '@ant-design/icons';

type IconsType = typeof Icons;

// 图标名称映射表（数据库中的简写 -> Ant Design 图标名）
const iconNameMap: Record<string, string> = {
    // 系统管理
    'setting': 'SettingOutlined',
    'user': 'UserOutlined',
    'team': 'TeamOutlined',
    'safety-certificate': 'SafetyCertificateOutlined',
    'menu': 'MenuOutlined',
    'tags': 'TagsOutlined',
    // 业务管理
    'appstore': 'AppstoreOutlined',
    'shopping': 'ShoppingOutlined',
    'customer-service': 'CustomerServiceOutlined',
    'warning': 'WarningOutlined',
    'branches': 'BranchesOutlined',
    'fork': 'BranchesOutlined',
    // 财务管理
    'dollar': 'DollarOutlined',
    'wallet': 'WalletOutlined',
    'percentage': 'PercentageOutlined',
    'bank': 'BankOutlined',
    'trophy': 'TrophyOutlined',
    'swap': 'SwapOutlined',
    // 监控中心
    'monitor': 'MonitorOutlined',
    'line-chart': 'LineChartOutlined',
    'fund': 'FundOutlined',
    'bar-chart': 'BarChartOutlined',
    // 内容管理
    'file-text': 'FileTextOutlined',
    'notification': 'NotificationOutlined',
    'picture': 'PictureOutlined',
    // 评价管理
    'star': 'StarOutlined',
    'comment': 'CommentOutlined',
    // 通用
    'dashboard': 'DashboardOutlined',
    'home': 'HomeOutlined',
    'profile': 'ProfileOutlined',
    'audit': 'AuditOutlined',
    'solution': 'SolutionOutlined',
    // 新增模块图标
    'dispute': 'ExclamationCircleOutlined',
    'coupon': 'GiftOutlined',
    'vip': 'CrownOutlined',
    'referral': 'UserSwitchOutlined',
    'routing': 'BranchesOutlined',
    'recharge': 'WalletOutlined',
    'activity': 'TrophyOutlined',
    'settlement': 'BankOutlined',
    'chat': 'MessageOutlined',
    'rank': 'RocketOutlined',
    'game': 'GameOutlined',
};

export const getIcon = (iconName?: string) => {
    if (!iconName) return null;

    // 先查找映射表
    const mappedName = iconNameMap[iconName];
    
    // 尝试多种方式查找图标
    const IconComponent = 
        (mappedName && (Icons as IconsType)[mappedName as keyof IconsType]) ||
        (Icons as IconsType)[iconName as keyof IconsType] || 
        (Icons as IconsType)[(iconName + 'Outlined') as keyof IconsType] ||
        (Icons as IconsType)[(iconName.charAt(0).toUpperCase() + iconName.slice(1) + 'Outlined') as keyof IconsType];

    if (IconComponent && typeof IconComponent === 'function') {
        const Icon = IconComponent as React.ComponentType;
        return <Icon />;
    }

    return null;
};
