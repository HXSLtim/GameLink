import React, { useState, useEffect } from 'react';
import { ConfigProvider, theme as antTheme } from 'antd';
import { ThemeContext, type ThemeMode } from './themeTypes';

// Re-export types for convenience
export type { ThemeMode, ThemeContextType } from './themeTypes';

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [mode, setMode] = useState<ThemeMode>(() => {
        const savedMode = localStorage.getItem('theme_mode');
        return (savedMode as ThemeMode) || 'light';
    });

    useEffect(() => {
        localStorage.setItem('theme_mode', mode);
        // 设置 body 的 data-theme 属性，方便自定义 CSS 使用
        document.body.setAttribute('data-theme', mode);
    }, [mode]);

    const toggleTheme = () => {
        setMode((prev) => (prev === 'light' ? 'dark' : 'light'));
    };

    return (
        <ThemeContext.Provider value={{ mode, toggleTheme }}>
            <ConfigProvider
                theme={{
                    algorithm: mode === 'dark' ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm,
                    token: {
                        // GameLink 品牌色 - 与移动端保持一致
                        colorPrimary: '#7ACC35', // 品牌绿色
                        colorSuccess: '#52c41a',
                        colorWarning: '#faad14',
                        colorError: '#ff4d4f',
                        colorInfo: '#1677ff',

                        // 品牌色渐变
                        colorPrimaryBg: '#e6f7ff',
                        colorPrimaryBgHover: '#bae7ff',
                        colorPrimaryBorder: '#91d5ff',
                        colorPrimaryBorderHover: '#69c0ff',
                        colorPrimaryHover: '#87D149', // 浅绿色
                        colorPrimaryActive: '#6DB72F', // 深绿色
                        colorPrimaryText: 'rgba(0, 0, 0, 0.88)',

                        // 圆角 - 与移动端保持一致
                        borderRadius: 8,
                        borderRadiusLG: 12,
                        borderRadiusSM: 6,

                        // 间距 - 与移动端保持一致
                        marginXS: 8,
                        marginSM: 12,
                        margin: 16,
                        marginMD: 20,
                        marginLG: 24,
                        marginXL: 32,
                    },
                    components: {
                        // 按钮组件定制
                        Button: {
                            colorPrimary: '#7ACC35',
                            colorPrimaryHover: '#87D149',
                            colorPrimaryActive: '#6DB72F',
                            algorithm: true, // 启用算法
                        },
                        // 卡片组件定制
                        Card: {
                            borderRadiusLG: 12,
                        },
                        // 输入框组件定制
                        Input: {
                            borderRadius: 8,
                        },
                        // 选择器组件定制
                        Select: {
                            borderRadius: 8,
                        },
                    },
                }}
            >
                {children}
            </ConfigProvider>
        </ThemeContext.Provider>
    );
};

// Note: useTheme is available from './useTheme' for Fast Refresh compatibility
