import React from 'react';
import { Outlet } from 'react-router-dom';

/**
 * Layout 组件 - 用于父级菜单，渲染子路由
 */
const LayoutOutlet: React.FC = () => <Outlet />;

export default LayoutOutlet;
