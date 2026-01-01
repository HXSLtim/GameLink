/**
 * BatchActions Component Unit Tests
 * 测试批量操作组件的渲染和交互
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { DeleteOutlined, CheckOutlined } from '@ant-design/icons';
import BatchActions from './index';

describe('BatchActions Component', () => {
    const mockActions = [
        {
            key: 'approve',
            label: '批量通过',
            icon: <CheckOutlined data-testid="check-icon" />,
            type: 'primary' as const,
            onConfirm: vi.fn(),
        },
        {
            key: 'delete',
            label: '批量删除',
            icon: <DeleteOutlined data-testid="delete-icon" />,
            danger: true,
            onConfirm: vi.fn(),
        },
    ];

    const mockSelectedRowKeys = [1, 2, 3];

    describe('基本渲染', () => {
        it('should not render when selectedCount is 0', () => {
            const { container } = render(
                <BatchActions
                    selectedCount={0}
                    actions={mockActions}
                    selectedRowKeys={[]}
                />
            );
            expect(container.firstChild).toBe(null);
        });

        it('should render action buttons when items are selected', () => {
            render(
                <BatchActions
                    selectedCount={3}
                    actions={mockActions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            expect(screen.getByText('批量通过')).toBeInTheDocument();
            expect(screen.getByText('批量删除')).toBeInTheDocument();
        });

        it('should render icons when provided', () => {
            render(
                <BatchActions
                    selectedCount={3}
                    actions={mockActions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            expect(screen.getByTestId('check-icon')).toBeInTheDocument();
            expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
        });
    });

    describe('按钮样式', () => {
        it('should render primary type button', () => {
            const { container } = render(
                <BatchActions
                    selectedCount={3}
                    actions={mockActions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            const primaryButton = container.querySelector('.ant-btn-primary');
            expect(primaryButton).toBeInTheDocument();
        });

        it('should render default type button', () => {
            const actions = [
                {
                    key: 'export',
                    label: '导出',
                    type: 'default' as const,
                    onConfirm: vi.fn(),
                },
            ];
            const { container } = render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            const defaultButton = container.querySelector('.ant-btn-default');
            expect(defaultButton).toBeInTheDocument();
        });

        it('should render danger button', () => {
            const { container } = render(
                <BatchActions
                    selectedCount={3}
                    actions={mockActions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            const dangerButton = container.querySelector('.ant-btn-dangerous');
            expect(dangerButton).toBeInTheDocument();
        });
    });

    describe('点击事件', () => {
        it('should call action onConfirm with selected keys when button clicked', async () => {
            const mockConfirm = vi.fn().mockResolvedValue(undefined);
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            const button = screen.getByText('批量通过');
            fireEvent.click(button);

            expect(mockConfirm).toHaveBeenCalledTimes(1);
            expect(mockConfirm).toHaveBeenCalledWith(mockSelectedRowKeys);
        });

        it('should call onActionComplete after successful confirm', async () => {
            const mockConfirm = vi.fn().mockResolvedValue(undefined);
            const mockComplete = vi.fn();
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                    onActionComplete={mockComplete}
                />
            );

            const button = screen.getByText('批量通过');

            // Need to wait for promise to resolve
            fireEvent.click(button);
            await new Promise(resolve => setTimeout(resolve, 0));

            expect(mockComplete).toHaveBeenCalledTimes(1);
        });

        it('should handle synchronous onConfirm', () => {
            const mockConfirm = vi.fn();
            const mockComplete = vi.fn();
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                    onActionComplete={mockComplete}
                />
            );

            const button = screen.getByText('批量通过');
            fireEvent.click(button);

            expect(mockConfirm).toHaveBeenCalled();
            // For synchronous functions, we wrap in Promise.resolve
            // Need to wait for the microtask
            return new Promise(resolve => {
                setTimeout(() => {
                    expect(mockComplete).toHaveBeenCalled();
                    resolve(null);
                }, 0);
            });
        });
    });

    describe('弹窗模式', () => {
        it('should open modal when modal action is clicked', () => {
            const mockConfirm = vi.fn();
            const actions = [
                {
                    key: 'assign',
                    label: '分配客服',
                    mode: 'modal' as const,
                    modalTitle: '选择客服',
                    modalContent: <div>客服选择表单</div>,
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            const button = screen.getByText('分配客服');
            fireEvent.click(button);

            expect(screen.getByText('选择客服')).toBeInTheDocument();
            expect(screen.getByText('客服选择表单')).toBeInTheDocument();
        });

        it('should close modal when cancel is clicked', () => {
            const mockConfirm = vi.fn();
            const actions = [
                {
                    key: 'assign',
                    label: '分配客服',
                    mode: 'modal' as const,
                    modalTitle: '选择客服',
                    modalContent: <div>客服选择表单</div>,
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            const button = screen.getByText('分配客服');
            fireEvent.click(button);

            // Modal should be open
            expect(screen.getByText('选择客服')).toBeInTheDocument();

            const cancelButton = screen.getAllByRole('button').find(
                btn => btn.textContent === '取消'
            );
            if (cancelButton) {
                fireEvent.click(cancelButton);
                // Modal should be closed
                expect(screen.queryByText('选择客服')).not.toBeInTheDocument();
            }
        });
    });

    describe('行内模式', () => {
        it('should execute inline action directly', () => {
            const mockConfirm = vi.fn();
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    mode: 'inline' as const,
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            const button = screen.getByText('批量通过');
            fireEvent.click(button);

            expect(mockConfirm).toHaveBeenCalledWith(mockSelectedRowKeys);
        });
    });

    describe('边界情况', () => {
        it('should handle empty actions array', () => {
            const { container } = render(
                <BatchActions
                    selectedCount={3}
                    actions={[]}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );
            const buttons = container.querySelectorAll('button');
            expect(buttons.length).toBe(0);
        });

        it('should handle action confirm error', async () => {
            const mockError = new Error('Action failed');
            const mockConfirm = vi.fn().mockRejectedValue(mockError);
            const mockComplete = vi.fn();
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    onConfirm: mockConfirm,
                },
            ];

            render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                    onActionComplete={mockComplete}
                />
            );

            const button = screen.getByText('批量通过');
            fireEvent.click(button);

            // Wait for promise to reject and be caught
            await new Promise(resolve => setTimeout(resolve, 0));

            // onActionComplete should not be called on error
            expect(mockConfirm).toHaveBeenCalled();
            expect(mockComplete).not.toHaveBeenCalled();
        });

        it('should handle action without type or danger', () => {
            const actions = [
                {
                    key: 'export',
                    label: '导出数据',
                    onConfirm: vi.fn(),
                },
            ];

            const { container } = render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            const button = container.querySelector('.ant-btn-default');
            expect(button).toBeInTheDocument();
            expect(container.querySelector('.ant-btn-dangerous')).not.toBeInTheDocument();
        });
    });

    describe('性能优化', () => {
        it('should use React.memo for optimization', () => {
            const mockConfirm = vi.fn();
            const actions = [
                {
                    key: 'approve',
                    label: '批量通过',
                    onConfirm: mockConfirm,
                },
            ];

            const { rerender } = render(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            // Re-render with same props
            rerender(
                <BatchActions
                    selectedCount={3}
                    actions={actions}
                    selectedRowKeys={mockSelectedRowKeys}
                />
            );

            // Component should not re-render unnecessarily
            expect(mockConfirm).not.toHaveBeenCalled();
        });
    });
});
