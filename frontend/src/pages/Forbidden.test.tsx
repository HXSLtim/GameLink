/**
 * Forbidden Page Tests
 * 
 * Tests for 403 forbidden page functionality
 * Requirements: 8.3 - 无权限重定向到 403 页面
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Forbidden from './Forbidden';

// Mock navigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

// Mock antd theme
vi.mock('antd', async () => {
    const actual = await vi.importActual('antd');
    return {
        ...actual,
        theme: {
            useToken: () => ({
                token: {
                    colorBgLayout: '#f5f5f5',
                    colorError: '#ff4d4f',
                    colorBgContainerDisabled: '#f5f5f5',
                    borderRadius: 6,
                    boxShadowSecondary: '0 2px 8px rgba(0,0,0,0.15)',
                },
            }),
        },
    };
});

describe('Forbidden Page', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should render 403 page with default content', () => {
        render(
            <MemoryRouter>
                <Forbidden />
            </MemoryRouter>
        );

        expect(screen.getByText('403')).toBeInTheDocument();
        expect(screen.getByText('抱歉，您没有权限访问此页面')).toBeInTheDocument();
    });

    it('should display from path when provided in state', () => {
        render(
            <MemoryRouter
                initialEntries={[
                    {
                        pathname: '/403',
                        state: { from: '/admin/users', requiredPermission: 'admin.users.list' },
                    },
                ]}
            >
                <Forbidden />
            </MemoryRouter>
        );

        expect(screen.getByText('/admin/users')).toBeInTheDocument();
        expect(screen.getByText('admin.users.list')).toBeInTheDocument();
    });

    it('should navigate to admin home when clicking home button', () => {
        render(
            <MemoryRouter>
                <Forbidden />
            </MemoryRouter>
        );

        const homeButton = screen.getByText('返回首页');
        fireEvent.click(homeButton);

        expect(mockNavigate).toHaveBeenCalledWith('/admin');
    });

    it('should navigate back when clicking back button with history', () => {
        // Mock history length > 2
        Object.defineProperty(window, 'history', {
            value: { length: 5 },
            writable: true,
        });

        render(
            <MemoryRouter>
                <Forbidden />
            </MemoryRouter>
        );

        const backButton = screen.getByText('返回上页');
        fireEvent.click(backButton);

        expect(mockNavigate).toHaveBeenCalledWith(-1);
    });

    it('should navigate to admin when clicking back button without history', () => {
        // Mock history length <= 2
        Object.defineProperty(window, 'history', {
            value: { length: 1 },
            writable: true,
        });

        render(
            <MemoryRouter>
                <Forbidden />
            </MemoryRouter>
        );

        const backButton = screen.getByText('返回上页');
        fireEvent.click(backButton);

        expect(mockNavigate).toHaveBeenCalledWith('/admin');
    });

    it('should have request permission button', () => {
        render(
            <MemoryRouter>
                <Forbidden />
            </MemoryRouter>
        );

        expect(screen.getByText('申请权限')).toBeInTheDocument();
    });
});
