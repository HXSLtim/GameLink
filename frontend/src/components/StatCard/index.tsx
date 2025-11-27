/**
 * 统计卡片组件
 * 用于仪表盘展示关键指标
 */
import React, { ReactNode } from 'react';
import { Card, Statistic, Tooltip, Skeleton } from 'antd';
import type { StatisticProps } from 'antd';
import {
    ArrowUpOutlined,
    ArrowDownOutlined,
    InfoCircleOutlined,
} from '@ant-design/icons';
import styles from './index.module.css';

export interface StatCardProps extends Omit<StatisticProps, 'title'> {
    /** 卡片标题 */
    title: ReactNode;
    /** 提示信息 */
    tooltip?: string;
    /** 图标 */
    icon?: ReactNode;
    /** 图标背景色 */
    iconBgColor?: string;
    /** 趋势（正数表示上升，负数表示下降） */
    trend?: number;
    /** 趋势描述 */
    trendLabel?: string;
    /** 底部描述 */
    footer?: ReactNode;
    /** 是否加载中 */
    loading?: boolean;
    /** 点击回调 */
    onClick?: () => void;
}

/**
 * StatCard组件
 */
const StatCard: React.FC<StatCardProps> = ({
    title,
    tooltip,
    icon,
    iconBgColor = '#1890ff',
    trend,
    trendLabel = '较昨日',
    footer,
    loading = false,
    onClick,
    ...statisticProps
}) => {
    const renderTrend = () => {
        if (trend === undefined) return null;

        const isUp = trend > 0;
        const trendValue = Math.abs(trend);
        const TrendIcon = isUp ? ArrowUpOutlined : ArrowDownOutlined;
        const trendColor = isUp ? '#52c41a' : '#ff4d4f';

        return (
            <div className={styles.trend} style={{ color: trendColor }}>
                <TrendIcon />
                <span>{trendValue}%</span>
                <span className={styles.trendLabel}>{trendLabel}</span>
            </div>
        );
    };

    return (
        <Card
            className={`${styles.card} ${onClick ? styles.clickable : ''}`}
            onClick={onClick}
            hoverable={!!onClick}
        >
            {loading ? (
                <Skeleton active paragraph={{ rows: 2 }} />
            ) : (
                <>
                    <div className={styles.header}>
                        <div className={styles.titleWrapper}>
                            <span className={styles.title}>{title}</span>
                            {tooltip && (
                                <Tooltip title={tooltip}>
                                    <InfoCircleOutlined className={styles.infoIcon} />
                                </Tooltip>
                            )}
                        </div>
                        {icon && (
                            <div
                                className={styles.iconWrapper}
                                style={{ backgroundColor: iconBgColor }}
                            >
                                {icon}
                            </div>
                        )}
                    </div>

                    <div className={styles.body}>
                        <Statistic {...statisticProps} />
                    </div>

                    {(trend !== undefined || footer) && (
                        <div className={styles.footer}>
                            {renderTrend()}
                            {footer && <div className={styles.footerContent}>{footer}</div>}
                        </div>
                    )}
                </>
            )}
        </Card>
    );
};

export default StatCard;
