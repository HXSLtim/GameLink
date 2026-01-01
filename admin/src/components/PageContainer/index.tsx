/**
 * 页面容器组件
 * 提供统一的页面头部和内容布局
 */
import React, { type ReactNode } from 'react';
import { Button, Space, Tabs } from 'antd';
import type { TabsProps } from 'antd';
import { ReloadOutlined } from '@ant-design/icons';
import styles from './index.module.css';

export interface PageContainerProps {
    /** 页面标题 */
    title?: ReactNode;
    /** 页面副标题 */
    subTitle?: ReactNode;
    /** 额外操作区域 */
    extra?: ReactNode;
    /** 标签页配置 */
    tabList?: TabsProps['items'];
    /** 当前激活的标签 */
    activeTab?: string;
    /** 标签切换回调 */
    onTabChange?: (key: string) => void;
    /** 是否显示刷新按钮 */
    showRefresh?: boolean;
    /** 刷新回调 */
    onRefresh?: () => void;
    /** 是否加载中 */
    loading?: boolean;
    /** 页面内容 */
    children?: ReactNode;
    /** 自定义类名 */
    className?: string;
}

/**
 * PageContainer组件
 * 优化: 使用 React.memo 避免不必要的重新渲染
 * 适用场景: 纯展示型布局组件，props 相同时不需要重新渲染
 */
const PageContainer: React.FC<PageContainerProps> = React.memo(({
    title,
    subTitle,
    extra,
    tabList,
    activeTab,
    onTabChange,
    showRefresh = false,
    onRefresh,
    loading = false,
    children,
    className,
}) => {
    return (
        <div className={`${styles.container} ${className || ''}`}>
            {/* 页面头部 */}
            {(title || extra) && (
                <div className={styles.header}>
                    <div className={styles.headerMain}>
                        <div className={styles.headerTitle}>
                            {title && <h1 className={styles.title}>{title}</h1>}
                            {subTitle && <span className={styles.subTitle}>{subTitle}</span>}
                        </div>
                        <div className={styles.headerExtra}>
                            <Space>
                                {showRefresh && (
                                    <Button
                                        icon={<ReloadOutlined spin={loading} />}
                                        onClick={onRefresh}
                                        loading={loading}
                                    >
                                        刷新
                                    </Button>
                                )}
                                {extra}
                            </Space>
                        </div>
                    </div>

                    {/* 标签页 */}
                    {tabList && tabList.length > 0 && (
                        <Tabs
                            activeKey={activeTab}
                            items={tabList}
                            onChange={onTabChange}
                            className={styles.tabs}
                        />
                    )}
                </div>
            )}

            {/* 页面内容 */}
            <div className={styles.content}>{children}</div>
        </div>
    );
});

export default PageContainer;
