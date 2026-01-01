/**
 * 加载状态组件
 * 提供统一的加载骨架屏
 */
import React from 'react';
import { Skeleton, Card, Row, Col } from 'antd';

export interface LoadingStateProps {
    /** 是否显示卡片包裹 */
    card?: boolean;
    /** 卡片标题 */
    title?: string;
    /** 骨架屏行数 */
    rows?: number;
    /** 是否加载中 */
    loading?: boolean;
    /** 子内容 */
    children?: React.ReactNode;
}

/**
 * LoadingState 组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 加载状态组件，props 相同时不需要重新渲染
 */
const LoadingState: React.FC<LoadingStateProps> = React.memo(({
    card = true,
    title,
    rows = 3,
    loading = true,
    children,
}) => {
    const content = loading ? (
        <Row gutter={[16, 16]}>
            <Col xs={24} sm={12} md={8}>
                <Skeleton active paragraph={{ rows }} />
            </Col>
            <Col xs={24} sm={12} md={8}>
                <Skeleton active paragraph={{ rows }} />
            </Col>
            <Col xs={24} sm={12} md={8}>
                <Skeleton active paragraph={{ rows }} />
            </Col>
        </Row>
    ) : (
        children
    );

    if (card) {
        return (
            <Card title={title}>
                {content}
            </Card>
        );
    }

    return <>{content}</>;
});

export default LoadingState;
