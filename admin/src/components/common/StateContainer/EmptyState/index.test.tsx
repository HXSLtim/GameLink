/**
 * EmptyState Component Unit Tests
 * 测试各种空状态类型的渲染和交互
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import EmptyState from './index';

describe('EmptyState Component', () => {
    describe('基本渲染', () => {
        it('should render default no-data state', () => {
            render(<EmptyState />);
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
            expect(screen.getByText('当前页面还没有任何数据')).toBeInTheDocument();
        });

        it('should render no-search state', () => {
            render(<EmptyState type="no-search" />);
            expect(screen.getByText('未找到相关结果')).toBeInTheDocument();
            expect(screen.getByText('请尝试调整搜索条件')).toBeInTheDocument();
        });

        it('should render no-permission state', () => {
            render(<EmptyState type="no-permission" />);
            expect(screen.getByText('暂无访问权限')).toBeInTheDocument();
            expect(screen.getByText('您没有权限查看此内容，请联系管理员')).toBeInTheDocument();
        });

        it('should render error state', () => {
            render(<EmptyState type="error" />);
            expect(screen.getByText('加载失败')).toBeInTheDocument();
            expect(screen.getByText('请稍后重试，如果问题持续存在请联系技术支持')).toBeInTheDocument();
        });
    });

    describe('自定义内容', () => {
        it('should render custom title and description', () => {
            render(
                <EmptyState
                    type="no-data"
                    title="自定义标题"
                    description="自定义描述"
                />
            );
            expect(screen.getByText('自定义标题')).toBeInTheDocument();
            expect(screen.getByText('自定义描述')).toBeInTheDocument();
        });

        it('should hide image when showImage is false', () => {
            const { container } = render(
                <EmptyState showImage={false} />
            );
            // When showImage is false, image should not be in the Empty component
            const emptyElement = container.querySelector('.ant-empty');
            expect(emptyElement).toBeInTheDocument();
            // The empty element should not have image class
            const imageElement = container.querySelector('.ant-empty-img');
            expect(imageElement).not.toBeInTheDocument();
        });

        it('should show image by default', () => {
            const { container } = render(
                <EmptyState />
            );
            const emptyElement = container.querySelector('.ant-empty');
            expect(emptyElement).toBeInTheDocument();
        });
    });

    describe('操作按钮', () => {
        it('should render action button when actionText and onAction provided', () => {
            const handleAction = vi.fn();
            render(
                <EmptyState
                    actionText="重新加载"
                    onAction={handleAction}
                />
            );
            const button = screen.getByRole('button', { name: '重新加载' });
            expect(button).toBeInTheDocument();
        });

        it('should call onAction when button is clicked', () => {
            const handleAction = vi.fn();
            render(
                <EmptyState
                    actionText="重新加载"
                    onAction={handleAction}
                />
            );
            const button = screen.getByRole('button', { name: '重新加载' });
            fireEvent.click(button);
            expect(handleAction).toHaveBeenCalledTimes(1);
        });

        it('should not render button when only actionText provided', () => {
            render(
                <EmptyState actionText="重新加载" />
            );
            const button = screen.queryByRole('button', { name: '重新加载' });
            expect(button).not.toBeInTheDocument();
        });

        it('should not render button when only onAction provided', () => {
            const handleAction = vi.fn();
            render(
                <EmptyState onAction={handleAction} />
            );
            const button = screen.queryByRole('button');
            expect(button).not.toBeInTheDocument();
        });
    });

    describe('custom 类型', () => {
        it('should render custom type with empty title/description when not provided', () => {
            const { container } = render(
                <EmptyState type="custom" />
            );
            // Should render but without title/description
            const emptyElement = container.querySelector('.ant-empty');
            expect(emptyElement).toBeInTheDocument();
        });

        it('should render custom type with provided title and description', () => {
            render(
                <EmptyState
                    type="custom"
                    title="自定义空状态"
                    description="这是自定义的空状态描述"
                />
            );
            expect(screen.getByText('自定义空状态')).toBeInTheDocument();
            expect(screen.getByText('这是自定义的空状态描述')).toBeInTheDocument();
        });
    });

    describe('边界情况', () => {
        it('should handle missing description when title is provided', () => {
            render(
                <EmptyState
                    type="no-data"
                    title="仅标题"
                />
            );
            expect(screen.getByText('仅标题')).toBeInTheDocument();
        });

        it('should render without crashing when no props provided', () => {
            expect(() => {
                render(<EmptyState />);
            }).not.toThrow();
        });

        it('should handle action callback gracefully', () => {
            const handleAction = vi.fn();
            render(
                <EmptyState
                    actionText="操作按钮"
                    onAction={handleAction}
                />
            );
            const button = screen.getByRole('button', { name: '操作按钮' });
            expect(() => {
                fireEvent.click(button);
            }).not.toThrow();
            expect(handleAction).toHaveBeenCalled();
        });
    });
});
