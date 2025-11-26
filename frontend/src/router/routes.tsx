import { RouteConfig } from './types';
import ClientLayout from '../layouts/ClientLayout';
import CompanionLayout from '../layouts/CompanionLayout';
import AdminLayout from '../layouts/AdminLayout';
import ClientHome from '../pages/client/Home';
import CompanionDashboard from '../pages/companion/Dashboard';
import AdminDashboard from '../pages/admin/Dashboard';
import Login from '../pages/auth/Login';

export const routes: RouteConfig[] = [
    {
        path: '/login',
        element: <Login />,
        meta: { title: 'Login' }
    },
    {
        path: '/',
        element: <ClientLayout />,
        children: [
            {
                path: '',
                element: <ClientHome />,
                meta: { title: 'Home' }
            }
        ]
    },
    {
        path: '/companion',
        element: <CompanionLayout />,
        meta: { roles: ['COMPANION'], requiresAuth: true, title: 'Companion Dashboard' },
        children: [
            {
                path: '',
                element: <CompanionDashboard />,
                meta: { title: 'Dashboard' }
            }
        ]
    },
    {
        path: '/admin',
        element: <AdminLayout />,
        meta: { roles: ['ADMIN', 'CS', 'FINANCE'], requiresAuth: true, title: 'Admin Console' },
        children: [
            {
                path: '',
                element: <AdminDashboard />,
                meta: { title: 'Dashboard' }
            }
        ]
    }
];
