/**
 * Discord风格三栏布局组件
 */
import React, { type ReactNode } from 'react';
import './DiscordLayout.less';

interface DiscordLayoutProps {
  children: ReactNode;
  className?: string;
}

interface SidebarProps {
  children: ReactNode;
  className?: string;
}

interface MainProps {
  children: ReactNode;
  className?: string;
}

interface HeaderProps {
  children: ReactNode;
  className?: string;
}

interface ContentProps {
  children: ReactNode;
  className?: string;
}

interface PanelProps {
  children: ReactNode;
  className?: string;
  visible?: boolean;
}

const DiscordLayout: React.FC<DiscordLayoutProps> & {
  Sidebar: React.FC<SidebarProps>;
  Main: React.FC<MainProps>;
  Header: React.FC<HeaderProps>;
  Content: React.FC<ContentProps>;
  Panel: React.FC<PanelProps>;
} = ({ children, className = '' }) => {
  return (
    <div className={`discord-layout ${className}`}>
      {children}
    </div>
  );
};

// 左侧导航栏
const Sidebar: React.FC<SidebarProps> = ({ children, className = '' }) => {
  return (
    <aside className={`discord-layout__sidebar ${className}`}>
      {children}
    </aside>
  );
};

// 主要内容区域
const Main: React.FC<MainProps> = ({ children, className = '' }) => {
  return (
    <main className={`discord-layout__main ${className}`}>
      {children}
    </main>
  );
};

// 头部
const Header: React.FC<HeaderProps> = ({ children, className = '' }) => {
  return (
    <header className={`discord-layout__header ${className}`}>
      {children}
    </header>
  );
};

// 内容区域
const Content: React.FC<ContentProps> = ({ children, className = '' }) => {
  return (
    <div className={`discord-layout__content ${className}`}>
      {children}
    </div>
  );
};

// 右侧面板
const Panel: React.FC<PanelProps> = ({ children, className = '', visible = true }) => {
  return (
    <aside className={`discord-layout__panel ${visible ? 'visible' : ''} ${className}`}>
      {children}
    </aside>
  );
};

DiscordLayout.Sidebar = Sidebar;
DiscordLayout.Main = Main;
DiscordLayout.Header = Header;
DiscordLayout.Content = Content;
DiscordLayout.Panel = Panel;

export default DiscordLayout;
