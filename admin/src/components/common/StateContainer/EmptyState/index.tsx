/**
 * 空状态组件
 * 用于展示无数据、无搜索结果等场景
 */
import React from 'react';
import { Empty, Button, Typography } from 'antd';
import type { EmptyProps } from 'antd';

const { Text } = Typography;

export interface EmptyStateProps extends Partial<EmptyProps> {
    /** 空状态类型 */
    type?: 'no-data' | 'no-search' | 'no-permission' | 'error' | 'custom';
    /** 标题文本 */
    title?: string;
    /** 描述文本 */
    description?: string;
    /** 操作按钮文本 */
    actionText?: string;
    /** 操作按钮回调 */
    onAction?: () => void;
    /** 是否显示图片 */
    showImage?: boolean;
}

/**
 * EmptyState 组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 纯展示型空状态组件，props 相同时不需要重新渲染
 */
const EmptyState: React.FC<EmptyStateProps> = React.memo(({
    type = 'no-data',
    title,
    description,
    actionText,
    onAction,
    showImage = true,
    ...restProps
}) => {
    // 预设配置
    const presets = {
        'no-data': {
            image: Empty.PRESENTED_IMAGE_SIMPLE,
            title: title || '暂无数据',
            description: description || '当前页面还没有任何数据',
        },
        'no-search': {
            image: Empty.PRESENTED_IMAGE_SIMPLE,
            title: title || '未找到相关结果',
            description: description || '请尝试调整搜索条件',
        },
        'no-permission': {
            image: Empty.PRESENTED_IMAGE_SIMPLE,
            title: title || '暂无访问权限',
            description: description || '您没有权限查看此内容，请联系管理员',
        },
        'error': {
            image: Empty.PRESENTED_IMAGE_SIMPLE,
            title: title || '加载失败',
            description: description || '请稍后重试，如果问题持续存在请联系技术支持',
        },
        'custom': {
            image: restProps.image,
            title: title || '',
            description: description || '',
        },
    };

    const preset = presets[type];

    return (
        <div style={{
            padding: '60px 0',
            textAlign: 'center',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: 300,
        }}>
            <Empty
                image={showImage ? preset.image : undefined}
                description={
                    <div>
                        {preset.title && (
                            <Text strong style={{ fontSize: 16, display: 'block', marginBottom: 8 }}>
                                {preset.title}
                            </Text>
                        )}
                        {preset.description && (
                            <Text type="secondary" style={{ fontSize: 14 }}>
                                {preset.description}
                            </Text>
                        )}
                    </div>
                }
                {...restProps}
            />
            {actionText && onAction && (
                <Button type="primary" onClick={onAction} style={{ marginTop: 16 }}>
                    {actionText}
                </Button>
            )}
        </div>
    );
});

export default EmptyState;
