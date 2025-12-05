import React from 'react';

// Lazy load components
// Lazy load components
const Audit = React.lazy(() => import('@/pages/sys/log'));
const Settings = React.lazy(() => import('@/pages/sys/setting'));
const ServiceItemList = React.lazy(() => import('@/pages/biz/service'));
const ServiceItemCreate = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemEdit = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = React.lazy(() => import('@/pages/biz/service/detail'));
const Menu = React.lazy(() => import('@/pages/sys/menu'));

export const componentMap: Record<string, React.LazyExoticComponent<React.FC<any>>> = {
    'sys/dashboard': React.lazy(() => import('@/pages/admin/Dashboard')),
    'sys/user': React.lazy(() => import('@/pages/admin/User')),
    'sys/menu': Menu,
    'biz/game': React.lazy(() => import('@/pages/admin/Game')),
    'biz/order': React.lazy(() => import('@/pages/admin/Order')),
    'biz/player': React.lazy(() => import('@/pages/admin/Player')),
    'sys/log': Audit,
    'sys/setting': Settings,
    'biz/service/list': ServiceItemList,
    'biz/service/create': ServiceItemCreate,
    'biz/service/edit': ServiceItemEdit,
    'biz/service/detail': ServiceItemDetail,
};

export const getComponent = (componentKey: string) => {
    return componentMap[componentKey] || null;
};
