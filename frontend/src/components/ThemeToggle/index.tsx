import React from 'react';
import { Button } from 'antd';
import { SunOutlined, MoonOutlined } from '@ant-design/icons';
import { useTheme } from '@/context/ThemeContext';

export const ThemeToggle: React.FC = () => {
    const { mode, toggleTheme } = useTheme();

    return (
        <Button
            type="text"
            icon={mode === 'light' ? <SunOutlined /> : <MoonOutlined />}
            onClick={toggleTheme}
            title={mode === 'light' ? '切换到暗色模式' : '切换到亮色模式'}
        />
    );
};
