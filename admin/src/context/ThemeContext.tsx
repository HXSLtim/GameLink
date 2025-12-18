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
                        colorPrimary: '#1890ff',
                    },
                }}
            >
                {children}
            </ConfigProvider>
        </ThemeContext.Provider>
    );
};

// Note: useTheme is available from './useTheme' for Fast Refresh compatibility
