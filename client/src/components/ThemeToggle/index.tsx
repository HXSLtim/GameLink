import { useTheme } from '@/context/ThemeContext';
import { Tooltip } from 'antd';
import './style.css';

const ThemeToggle = () => {
    const { theme, toggleTheme } = useTheme();

    return (
        <Tooltip title={`切换到${theme === 'dark' ? '浅色' : '深色'}主题`} placement="right">
            <button
                className={`theme-toggle ${theme}`}
                onClick={toggleTheme}
                aria-label="Toggle Theme"
            >
                <div className="toggle-track">
                    <div className="toggle-thumb">
                        {theme === 'dark' ? (
                            // Moon Icon
                            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="toggle-icon">
                                <path d="M12.3 22c-5 0-9.2-3.6-10-8.4C2 10 5.4 6.4 9.6 5.7c.3-.1.6.2.5.5C10 7.3 10.9 8 12 8c4.4 0 8 3.6 8 8 0 1.1.7 2 1.7 1.9.4.1.6.4.6.7-.7 2-2.3 3.3-4.4 3.4-1.8.1-3.6 0-5.6 0z" fill="currentColor" />
                            </svg>
                        ) : (
                            // Sun Icon 
                            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="toggle-icon">
                                <path d="M12 7c-2.8 0-5 2.2-5 5s2.2 5 5 5 5-2.2 5-5-2.2-5-5-5zm0-5c.6 0 1 .4 1 1v2c0 .6-.4 1-1 1s-1-.4-1-1V3c0-.6.4-1 1-1zm0 18c.6 0 1 .4 1 1v2c0 .6-.4 1-1 1s-1-.4-1-1v-2c0-.6.4-1 1-1zm9-9c0 .6-.4 1-1 1h-2c-.6 0-1-.4-1-1s.4-1 1-1h2c.6 0 1 .4 1 1zM4 12c0 .6-.4 1-1 1H1c-.6 0-1-.4-1-1s.4-1 1-1h2c.6 0 1 .4 1 1zm15.1-6.3c.4.4.4 1 0 1.4l-1.4 1.4c-.4.4-1 .4-1.4 0-.4-.4-.4-1 0-1.4l1.4-1.4c.4-.4 1-.4 1.4 0zM6.3 19.1c.4.4.4 1 0 1.4l-1.4 1.4c-.4.4-1 .4-1.4 0-.4-.4-.4-1 0-1.4l1.4-1.4c.4-.4 1-.4 1.4 0zm12.8 12.8c.4.4.4 1 0 1.4l-1.4 1.4c-.4.4-1 .4-1.4 0-.4-.4-.4-1 0-1.4l1.4-1.4c.4-.4 1-.4 1.4 0zM6.3 4.9c.4.4.4 1 0 1.4L4.9 7.7c-.4.4-1 .4-1.4 0-.4-.4-.4-1 0-1.4l1.4-1.4c.4-.4 1-.4 1.4 0z" fill="currentColor" />
                            </svg>
                        )}
                    </div>
                </div>
            </button>
        </Tooltip>
    );
};

export default ThemeToggle;
