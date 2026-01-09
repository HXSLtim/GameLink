import React, { useCallback } from 'react';
import { Button } from 'antd';
import { SunOutlined, MoonOutlined } from '@ant-design/icons';
import { useTheme } from '@/context/useTheme';

/**
 * 主题切换按钮
 * 优化: 使用 React.memo + useCallback 减少重渲染
 */
export const ThemeToggle: React.FC = React.memo(() => {
    const { mode, toggleTheme } = useTheme();

    const handleClick = useCallback(() => {
        toggleTheme();
    }, [toggleTheme]);

    return (
        <Button
            type="text"
            icon={mode === 'light' ? <SunOutlined /> : <MoonOutlined />}
            onClick={handleClick}
            title={mode === 'light' ? '切换到暗色模式' : '切换到亮色模式'}
        />
    );
});
