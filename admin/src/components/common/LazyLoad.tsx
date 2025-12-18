/**
 * 懒加载包装器组件
 */
import { Suspense } from 'react';
import type { ReactNode } from 'react';
import { Spin } from 'antd';

interface LazyLoadProps {
    children: ReactNode;
}

const LazyLoad = ({ children }: LazyLoadProps) => (
    <Suspense
        fallback={
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%', minHeight: 300 }}>
                <Spin size="large" />
            </div>
        }
    >
        {children}
    </Suspense>
);

export default LazyLoad;
