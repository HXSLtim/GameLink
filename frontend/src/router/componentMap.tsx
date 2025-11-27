import React from 'react';

// Lazy load components
const Dashboard = React.lazy(() => import('@/pages/sys/dashboard'));
const Users = React.lazy(() => import('@/pages/sys/user'));
const Games = React.lazy(() => import('@/pages/biz/game'));
const Orders = React.lazy(() => import('@/pages/biz/order'));
const Audit = React.lazy(() => import('@/pages/sys/log'));
const Settings = React.lazy(() => import('@/pages/sys/setting'));
const ServiceItemList = React.lazy(() => import('@/pages/biz/service'));
const ServiceItemCreate = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemEdit = React.lazy(() => import('@/pages/biz/service/form'));
const ServiceItemDetail = React.lazy(() => import('@/pages/biz/service/detail'));
const Menu = React.lazy(() => import('@/pages/sys/menu'));

export const componentMap: Record<string, React.LazyExoticComponent<React.FC<any>>> = {
    'sys/dashboard': Dashboard,
    'sys/user': Users,
    'sys/menu': Menu,
    'biz/game': Games,
    'biz/order': Orders,
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
