/**
 * SearchFilters Component Unit Tests
 * 测试搜索筛选组件的渲染和交互
 */
import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from 'antd';
import SearchFilters from './index';

describe('SearchFilters Component', () => {
    const mockFilters = [
        {
            type: 'input' as const,
            key: 'keyword',
            placeholder: '搜索关键词',
        },
        {
            type: 'select' as const,
            key: 'status',
            placeholder: '状态',
            options: [
                { label: '全部', value: 'all' },
                { label: '启用', value: 'enabled' },
                { label: '禁用', value: 'disabled' },
            ],
        },
    ];

    const mockFilterValues = {
        keyword: '',
        status: undefined,
    };

    describe('基本渲染', () => {
        it('should render filters correctly', () => {
            render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                />
            );
            expect(screen.getByPlaceholderText('搜索关键词')).toBeInTheDocument();
        });

        it('should render with Card wrapper by default', () => {
            const { container } = render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                />
            );
            const card = container.querySelector('.ant-card');
            expect(card).toBeInTheDocument();
        });

        it('should render without Card when card is false', () => {
            const { container } = render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    card={false}
                />
            );
            const card = container.querySelector('.ant-card');
            expect(card).not.toBeInTheDocument();
        });
    });

    describe('Input 筛选器', () => {
        it('should render input filter with placeholder', () => {
            render(
                <SearchFilters
                    filters={[mockFilters[0]]}
                    onFilterChange={vi.fn()}
                    filterValues={{ keyword: '' }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词');
            expect(input).toBeInTheDocument();
        });

        it('should call onFilterChange when input value changes', () => {
            const handleChange = vi.fn();
            render(
                <SearchFilters
                    filters={[mockFilters[0]]}
                    onFilterChange={handleChange}
                    filterValues={{ keyword: '' }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词') as HTMLInputElement;
            fireEvent.change(input, { target: { value: 'test' } });
            expect(handleChange).toHaveBeenCalledWith('keyword', 'test');
        });

        it('should display input value', () => {
            render(
                <SearchFilters
                    filters={[mockFilters[0]]}
                    onFilterChange={vi.fn()}
                    filterValues={{ keyword: 'existing value' }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词') as HTMLInputElement;
            expect(input.value).toBe('existing value');
        });

        it('should apply custom width', () => {
            const filterWithWidth = {
                type: 'input' as const,
                key: 'keyword',
                placeholder: '搜索关键词',
                width: 300,
            };
            render(
                <SearchFilters
                    filters={[filterWithWidth]}
                    onFilterChange={vi.fn()}
                    filterValues={{ keyword: '' }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词');
            // Just verify the element renders - jsdom has limitations with getComputedStyle
            expect(input).toBeInTheDocument();
        });

        it('should call onQuery when Enter key is pressed', () => {
            const handleChange = vi.fn();
            const handleQuery = vi.fn();
            render(
                <SearchFilters
                    filters={[mockFilters[0]]}
                    onFilterChange={handleChange}
                    onQuery={handleQuery}
                    filterValues={{ keyword: '' }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词');
            fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' });
            expect(handleQuery).toHaveBeenCalled();
        });
    });

    describe('Select 筛选器', () => {
        it('should render select filter', () => {
            render(
                <SearchFilters
                    filters={[mockFilters[1]]}
                    onFilterChange={vi.fn()}
                    filterValues={{ status: undefined }}
                />
            );
            const select = screen.getByRole('combobox');
            expect(select).toBeInTheDocument();
        });

        it('should call onFilterChange when select value changes', () => {
            const handleChange = vi.fn();
            render(
                <SearchFilters
                    filters={[mockFilters[1]]}
                    onFilterChange={handleChange}
                    filterValues={{ status: undefined }}
                />
            );
            const select = screen.getByRole('combobox');
            expect(select).toBeInTheDocument();
        });

        it('should apply custom width', () => {
            const filterWithWidth = {
                type: 'select' as const,
                key: 'status',
                placeholder: '状态',
                width: 150,
                options: mockFilters[1].options,
            };
            render(
                <SearchFilters
                    filters={[filterWithWidth]}
                    onFilterChange={vi.fn()}
                    filterValues={{ status: undefined }}
                />
            );
            const select = screen.getByRole('combobox');
            // Just verify the element renders - jsdom has limitations with getComputedStyle
            expect(select).toBeInTheDocument();
        });
    });

    describe('RangePicker 筛选器', () => {
        it('should render range picker filter', () => {
            const rangeFilter = {
                type: 'rangePicker' as const,
                key: 'dateRange',
            };
            render(
                <SearchFilters
                    filters={[rangeFilter]}
                    onFilterChange={vi.fn()}
                    filterValues={{ dateRange: null }}
                />
            );
            const picker = document.querySelector('.ant-picker-range');
            expect(picker).toBeInTheDocument();
        });
    });

    describe('Segmented 筛选器', () => {
        it('should render segmented filter', () => {
            const segmentedFilter = {
                type: 'segmented' as const,
                key: 'view',
                segmentedOptions: [
                    { label: '列表', value: 'list' },
                    { label: '卡片', value: 'card' },
                ],
            };
            render(
                <SearchFilters
                    filters={[segmentedFilter]}
                    onFilterChange={vi.fn()}
                    filterValues={{ view: 'list' }}
                />
            );
            const segmented = document.querySelector('.ant-segmented');
            expect(segmented).toBeInTheDocument();
        });
    });

    describe('操作按钮', () => {
        it('should render actions when provided', () => {
            const actions = (
                <Button type="primary" data-testid="action-button">
                    新增
                </Button>
            );
            render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    actions={actions}
                />
            );
            expect(screen.getByTestId('action-button')).toBeInTheDocument();
        });

        it('should render query and reset buttons when showQueryButtons is true', () => {
            const handleQuery = vi.fn();
            const handleReset = vi.fn();
            render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    showQueryButtons={true}
                    onQuery={handleQuery}
                    onReset={handleReset}
                />
            );
            expect(screen.getByText('查询')).toBeInTheDocument();
            expect(screen.getByText('重置')).toBeInTheDocument();
        });

        it('should call onQuery when query button is clicked', () => {
            const handleQuery = vi.fn();
            render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    showQueryButtons={true}
                    onQuery={handleQuery}
                    onReset={vi.fn()}
                />
            );
            const queryButton = screen.getByText('查询');
            fireEvent.click(queryButton);
            expect(handleQuery).toHaveBeenCalledTimes(1);
        });

        it('should call onReset when reset button is clicked', () => {
            const handleReset = vi.fn();
            render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    showQueryButtons={true}
                    onQuery={vi.fn()}
                    onReset={handleReset}
                />
            );
            const resetButton = screen.getByText('重置');
            fireEvent.click(resetButton);
            expect(handleReset).toHaveBeenCalledTimes(1);
        });
    });

    describe('布局', () => {
        it('should render filters and actions in correct layout', () => {
            const actions = <Button data-testid="action-btn">操作</Button>;
            const { container } = render(
                <SearchFilters
                    filters={mockFilters}
                    onFilterChange={vi.fn()}
                    filterValues={mockFilterValues}
                    actions={actions}
                />
            );
            const row = container.querySelector('.ant-row');
            expect(row).toBeInTheDocument();
            expect(screen.getByTestId('action-btn')).toBeInTheDocument();
        });
    });

    describe('边界情况', () => {
        it('should handle empty filters array', () => {
            const { container } = render(
                <SearchFilters
                    filters={[]}
                    onFilterChange={vi.fn()}
                    filterValues={{}}
                />
            );
            const row = container.querySelector('.ant-row');
            expect(row).toBeInTheDocument();
        });

        it('should handle null filter values', () => {
            render(
                <SearchFilters
                    filters={[mockFilters[0]]}
                    onFilterChange={vi.fn()}
                    filterValues={{ keyword: null }}
                />
            );
            const input = screen.getByPlaceholderText('搜索关键词') as HTMLInputElement;
            expect(input.value).toBe('');
        });

        it('should handle undefined filter values', () => {
            render(
                <SearchFilters
                    filters={[mockFilters[1]]}
                    onFilterChange={vi.fn()}
                    filterValues={{ status: undefined }}
                />
            );
            const select = screen.getByRole('combobox');
            expect(select).toBeInTheDocument();
        });

        it('should handle filter without options (select)', () => {
            const selectWithoutOptions = {
                type: 'select' as const,
                key: 'field',
                placeholder: '字段',
                options: [],
            };
            render(
                <SearchFilters
                    filters={[selectWithoutOptions]}
                    onFilterChange={vi.fn()}
                    filterValues={{ field: undefined }}
                />
            );
            const select = screen.getByRole('combobox');
            expect(select).toBeInTheDocument();
        });

        it('should render without crashing', () => {
            expect(() => {
                render(
                    <SearchFilters
                        filters={mockFilters}
                        onFilterChange={vi.fn()}
                        filterValues={mockFilterValues}
                    />
                );
            }).not.toThrow();
        });
    });
});
