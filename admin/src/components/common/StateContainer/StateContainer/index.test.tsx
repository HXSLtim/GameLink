/**
 * StateContainer Component Unit Tests
 * 测试状态容器的各种状态切换
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import StateContainer from './index';

describe('StateContainer Component', () => {
    describe('加载状态', () => {
        it('should render loading state when loading is true', () => {
            const { container } = render(
                <StateContainer loading={true} data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            const skeletons = container.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });

        it('should render loading with custom config', () => {
            render(
                <StateContainer
                    loading={true}
                    data={[]}
                    loadingConfig={{ card: false, rows: 5 }}
                >
                    <div>内容</div>
                </StateContainer>
            );
            // Should render skeleton without card
            const card = document.querySelector('.ant-card');
            expect(card).not.toBeInTheDocument();
        });

        it('should not render children when loading', () => {
            render(
                <StateContainer loading={true} data={[1, 2, 3]}>
                    <div>不应该显示的内容</div>
                </StateContainer>
            );
            expect(screen.queryByText('不应该显示的内容')).not.toBeInTheDocument();
        });
    });

    describe('错误状态', () => {
        it('should render error state when error is provided', () => {
            render(
                <StateContainer error="网络错误" data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('网络错误')).toBeInTheDocument();
        });

        it('should render error state with default message', () => {
            render(
                <StateContainer error="加载失败" data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            const errorMessages = screen.getAllByText('加载失败');
            expect(errorMessages.length).toBeGreaterThan(0);
        });

        it('should call onEmptyAction when retry button clicked', () => {
            const handleRetry = vi.fn();
            render(
                <StateContainer error="加载失败" data={[]} onEmptyAction={handleRetry}>
                    <div>内容</div>
                </StateContainer>
            );
            const button = screen.queryByRole('button', { name: '重新加载' });
            if (button) {
                fireEvent.click(button);
                expect(handleRetry).toHaveBeenCalledTimes(1);
            }
        });

        it('should prioritize error over loading', () => {
            render(
                <StateContainer loading={true} error="错误信息" data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('错误信息')).toBeInTheDocument();
            const skeletons = document.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBe(0);
        });
    });

    describe('空状态', () => {
        it('should render empty state when data array is empty', () => {
            render(
                <StateContainer data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });

        it('should render empty state when data is null/undefined', () => {
            const { rerender } = render(
                <StateContainer data={null}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();

            rerender(
                <StateContainer data={undefined}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });

        it('should render custom empty type', () => {
            render(
                <StateContainer data={[]} emptyType="no-search">
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('未找到相关结果')).toBeInTheDocument();
        });

        it('should render custom empty title and description', () => {
            render(
                <StateContainer
                    data={[]}
                    emptyTitle="自定义标题"
                    emptyDescription="自定义描述"
                >
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('自定义标题')).toBeInTheDocument();
            expect(screen.getByText('自定义描述')).toBeInTheDocument();
        });

        it('should render custom action button', () => {
            const handleAction = vi.fn();
            render(
                <StateContainer
                    data={[]}
                    emptyActionText="清除筛选"
                    onEmptyAction={handleAction}
                >
                    <div>内容</div>
                </StateContainer>
            );
            const button = screen.queryByRole('button', { name: '清除筛选' });
            if (button) {
                fireEvent.click(button);
                expect(handleAction).toHaveBeenCalledTimes(1);
            }
        });
    });

    describe('正常状态', () => {
        it('should render children when data has items', () => {
            render(
                <StateContainer data={[1, 2, 3]}>
                    <div>正常内容</div>
                </StateContainer>
            );
            expect(screen.getByText('正常内容')).toBeInTheDocument();
        });

        it('should render children when data is truthy non-array', () => {
            render(
                <StateContainer data={{ key: 'value' }}>
                    <div>对象数据内容</div>
                </StateContainer>
            );
            expect(screen.getByText('对象数据内容')).toBeInTheDocument();
        });

        it('should render children when data is true', () => {
            render(
                <StateContainer data={true}>
                    <div>布尔数据内容</div>
                </StateContainer>
            );
            expect(screen.getByText('布尔数据内容')).toBeInTheDocument();
        });

        it('should render children when data is non-zero number', () => {
            render(
                <StateContainer data={42}>
                    <div>数字数据内容</div>
                </StateContainer>
            );
            expect(screen.getByText('数字数据内容')).toBeInTheDocument();
        });

        it('should render children when data is non-empty string', () => {
            render(
                <StateContainer data="some text">
                    <div>字符串数据内容</div>
                </StateContainer>
            );
            expect(screen.getByText('字符串数据内容')).toBeInTheDocument();
        });
    });

    describe('边界情况', () => {
        it('should treat empty array as empty', () => {
            render(
                <StateContainer data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });

        it('should treat array with items as data', () => {
            render(
                <StateContainer data={[{}]}>
                    <div>有数据</div>
                </StateContainer>
            );
            expect(screen.getByText('有数据')).toBeInTheDocument();
        });

        it('should treat zero as empty', () => {
            render(
                <StateContainer data={0}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });

        it('should treat false as empty', () => {
            render(
                <StateContainer data={false}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });

        it('should treat empty string as empty', () => {
            render(
                <StateContainer data={''}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });
    });

    describe('状态优先级', () => {
        it('should prioritize error over loading and empty', () => {
            render(
                <StateContainer loading={true} error="错误" data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('错误')).toBeInTheDocument();
        });

        it('should prioritize loading over empty', () => {
            render(
                <StateContainer loading={true} data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            const skeletons = document.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);
        });

        it('should show empty when no loading and no error', () => {
            render(
                <StateContainer loading={false} error={null} data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();
        });
    });

    describe('复杂场景', () => {
        it('should handle complete loading flow', () => {
            const { rerender } = render(
                <StateContainer loading={true} data={[]}>
                    <div>内容</div>
                </StateContainer>
            );

            // Initially loading
            const skeletons = document.querySelectorAll('.ant-skeleton');
            expect(skeletons.length).toBeGreaterThan(0);

            // Then empty
            rerender(
                <StateContainer loading={false} data={[]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('暂无数据')).toBeInTheDocument();

            // Then data
            rerender(
                <StateContainer loading={false} data={[1, 2, 3]}>
                    <div>内容</div>
                </StateContainer>
            );
            expect(screen.getByText('内容')).toBeInTheDocument();
        });

        it('should render complex children with data', () => {
            render(
                <StateContainer data={[{ id: 1, name: 'Item 1' }]}>
                    <div>
                        <span>列表项</span>
                        <button>操作</button>
                    </div>
                </StateContainer>
            );
            expect(screen.getByText('列表项')).toBeInTheDocument();
            expect(screen.getByText('操作')).toBeInTheDocument();
        });
    });
});
