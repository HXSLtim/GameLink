/**
 * 统计卡片组件
 * 用于仪表盘展示关键指标
 */
import React, { useMemo } from 'react';
import type { ReactNode } from 'react';
import { Card, Statistic, Tooltip, Skeleton } from 'antd';
import type { StatisticProps } from 'antd';
import {
    ArrowUpOutlined,
    ArrowDownOutlined,
    InfoCircleOutlined,
} from '@ant-design/icons';
import { useCountUp } from '@/hooks/useCountUp';
import { ANIMATION_DURATION } from '@/constants/dashboard';
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
    /** 是否启用数字动画 */
    animated?: boolean;
    /** 动画持续时间（毫秒） */
    animationDuration?: number;
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
    animated = true,
    animationDuration = ANIMATION_DURATION.number,
    ...statisticProps
}) => {
    // 提取数值用于动画
    const numericValue = useMemo(() => {
        if (typeof statisticProps.value === 'number') {
            return statisticProps.value;
        }
        if (typeof statisticProps.value === 'string') {
            const parsed = parseFloat(statisticProps.value);
            return isNaN(parsed) ? 0 : parsed;
        }
        return 0;
    }, [statisticProps.value]);

    // 使用数字动画Hook
    const animatedValue = useCountUp(numericValue, {
        duration: animationDuration,
        enabled: animated && !loading,
    });

    // 确定最终显示的值
    const displayValue = useMemo(() => {
        if (!animated || loading) {
            return statisticProps.value;
        }
        
        // 如果原始值是数字，使用动画值
        if (typeof statisticProps.value === 'number') {
            // 保持原始值的小数位数
            const decimals = statisticProps.precision ?? 0;
            return Number(animatedValue.toFixed(decimals));
        }
        
        return statisticProps.value;
    }, [animated, loading, statisticProps.value, statisticProps.precision, animatedValue]);

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
                        <Statistic {...statisticProps} value={displayValue} />
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
