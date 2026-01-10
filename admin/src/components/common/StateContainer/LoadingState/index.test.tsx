/**
 * LoadingState Component Unit Tests
 * 测试加载状态的渲染和配置
 */
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import LoadingState from './index';

describe('LoadingState Component', () => {
    describe('基本渲染', () => {
        it('should render skeleton when loading is true', () => {
            const { container } = render(<LoadingState loading={true} />);
            const skeletons = container.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });

        it('should render children when loading is false', () => {
            render(
                <LoadingState loading={false}>
                    <div>加载完成的内容</div>
                </LoadingState>
            );
            expect(screen.getByText('加载完成的内容')).toBeInTheDocument();
        });

        it('should render skeleton by default (loading prop defaults to true)', () => {
            const { container } = render(<LoadingState />);
            const skeletons = container.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });
    });

    describe('卡片包裹', () => {
        it('should render with Card wrapper by default', () => {
            const { container } = render(<LoadingState loading={true} />);
            const card = container.querySelector('.ant-card');
            expect(card).toBeInTheDocument();
        });

        it('should render without Card wrapper when card is false', () => {
            const { container } = render(<LoadingState loading={true} card={false} />);
            const card = container.querySelector('.ant-card');
            expect(card).not.toBeInTheDocument();
        });

        it('should render title when provided', () => {
            render(<LoadingState loading={true} title="加载中..." />);
            expect(screen.getByText('加载中...')).toBeInTheDocument();
        });

        it('should not render title when not provided', () => {
            const { container } = render(<LoadingState loading={true} />);
            const cardTitle = container.querySelector('.ant-card-head-title');
            expect(cardTitle).not.toBeInTheDocument();
        });
    });

    describe('骨架屏配置', () => {
        it('should render skeleton columns', () => {
            const { container } = render(<LoadingState loading={true} />);
            // Should render skeletons
            const skeletons = container.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });

        it('should render with active animation', () => {
            const { container } = render(<LoadingState loading={true} />);
            const activeSkeleton = container.querySelector('.ant-skeleton-active');
            expect(activeSkeleton).toBeInTheDocument();
        });
    });

    describe('子内容渲染', () => {
        it('should render multiple children', () => {
            render(
                <LoadingState loading={false}>
                    <div>子内容 1</div>
                    <div>子内容 2</div>
                    <div>子内容 3</div>
                </LoadingState>
            );
            expect(screen.getByText('子内容 1')).toBeInTheDocument();
            expect(screen.getByText('子内容 2')).toBeInTheDocument();
            expect(screen.getByText('子内容 3')).toBeInTheDocument();
        });

        it('should render null children without crashing', () => {
            const { container } = render(
                <LoadingState loading={false}>
                    {null}
                </LoadingState>
            );
            expect(container.firstChild).toBeInTheDocument();
        });

        it('should render complex children', () => {
            render(
                <LoadingState loading={false}>
                    <div>
                        <span>嵌套内容</span>
                        <button>按钮</button>
                    </div>
                </LoadingState>
            );
            expect(screen.getByText('嵌套内容')).toBeInTheDocument();
            expect(screen.getByText('按钮')).toBeInTheDocument();
        });
    });

    describe('组合场景', () => {
        it('should render card with title and skeleton', () => {
            render(<LoadingState loading={true} title="数据加载中" />);
            expect(screen.getByText('数据加载中')).toBeInTheDocument();
            const skeletons = document.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });

        it('should render card with title and children when not loading', () => {
            render(
                <LoadingState loading={false} title="用户列表">
                    <div>用户数据</div>
                </LoadingState>
            );
            expect(screen.getByText('用户列表')).toBeInTheDocument();
            expect(screen.getByText('用户数据')).toBeInTheDocument();
        });

        it('should render without card when card is false and loading', () => {
            const { container } = render(
                <LoadingState loading={true} card={false} />
            );
            const card = container.querySelector('.ant-card');
            expect(card).not.toBeInTheDocument();
            const skeletons = container.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });
    });

    describe('边界情况', () => {
        it('should render without crashing', () => {
            expect(() => {
                render(<LoadingState loading={true} />);
            }).not.toThrow();
        });

        it('should handle empty string title gracefully', () => {
            render(<LoadingState loading={true} title="" />);
            // Card doesn't render header when title is empty
            const cardTitle = document.querySelector('.ant-card-head-title');
            expect(cardTitle).not.toBeInTheDocument();
        });
    });
});
