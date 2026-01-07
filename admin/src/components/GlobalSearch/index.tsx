/**
 * 全局搜索组件
 * 支持搜索订单、用户、陪玩师等
 */
import React, { useState, useCallback, useRef, useEffect } from 'react';
import {
    Input,
    Modal,
    List,
    Avatar,
    Tag,
    Space,
    Typography,
    Spin,
    Empty,
    Tabs,
} from 'antd';
import type { InputRef } from 'antd';
import {
    SearchOutlined,
    UserOutlined,
    ShoppingOutlined,
    TeamOutlined,
    FileTextOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { adminApi } from '@/api/admin';
import { logger } from '@/utils/logger';

const { Text } = Typography;

interface SearchResult {
    type: 'order' | 'user' | 'player' | 'dispute';
    id: number;
    title: string;
    subtitle: string;
    avatar?: string;
    status?: string;
    statusColor?: string;
}

interface GlobalSearchProps {
    open: boolean;
    onClose: () => void;
}

const GlobalSearch: React.FC<GlobalSearchProps> = ({ open, onClose }) => {
    const navigate = useNavigate();
    const [keyword, setKeyword] = useState('');
    const [loading, setLoading] = useState(false);
    const [results, setResults] = useState<SearchResult[]>([]);
    const [activeTab, setActiveTab] = useState('all');
    const inputRef = useRef<InputRef>(null);
    const searchTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    // 聚焦输入框
    useEffect(() => {
        if (open) {
            setTimeout(() => {
                inputRef.current?.focus();
            }, 100);
        } else {
            setKeyword('');
            setResults([]);
        }
    }, [open]);

    // 搜索函数
    const doSearch = useCallback(async (searchKeyword: string) => {
        if (!searchKeyword.trim()) {
            setResults([]);
            return;
        }

        setLoading(true);
        const allResults: SearchResult[] = [];

        try {
            // 并行搜索多个类型
            const [ordersRes, usersRes, playersRes] = await Promise.allSettled([
                adminApi.getOrders({ orderNumber: searchKeyword, page_size: 5 }),
                adminApi.getUsers({ keyword: searchKeyword, page_size: 5 }),
                adminApi.getPlayers({ keyword: searchKeyword, page_size: 5 }),
            ]);

            // 处理订单结果
            if (ordersRes.status === 'fulfilled' && ordersRes.value.data.success) {
                const orders = ordersRes.value.data.data || [];
                orders.forEach((order: { id: number; orderNo: string; status: string; totalPriceCents: number; user?: { name?: string } }) => {
                    allResults.push({
                        type: 'order',
                        id: order.id,
                        title: `订单 ${order.orderNo}`,
                        subtitle: `¥${(order.totalPriceCents / 100).toFixed(2)} - ${order.user?.name || '未知用户'}`,
                        status: order.status,
                        statusColor: getOrderStatusColor(order.status),
                    });
                });
            }

            // 处理用户结果
            if (usersRes.status === 'fulfilled' && usersRes.value.data.success) {
                const users = usersRes.value.data.data || [];
                users.forEach((user: { id: number; name?: string; phone?: string; avatarUrl?: string; status?: string }) => {
                    allResults.push({
                        type: 'user',
                        id: user.id,
                        title: user.name || `用户${user.id}`,
                        subtitle: user.phone || '',
                        avatar: user.avatarUrl,
                        status: user.status,
                        statusColor: user.status === 'active' ? 'success' : 'default',
                    });
                });
            }

            // 处理陪玩师结果
            if (playersRes.status === 'fulfilled' && playersRes.value.data.success) {
                const players = playersRes.value.data.data || [];
                players.forEach((player: { id: number; nickname?: string; user?: { name?: string; avatarUrl?: string }; verificationStatus?: string }) => {
                    allResults.push({
                        type: 'player',
                        id: player.id,
                        title: player.nickname || player.user?.name || `陪玩师${player.id}`,
                        subtitle: `ID: ${player.id}`,
                        avatar: player.user?.avatarUrl,
                        status: player.verificationStatus,
                        statusColor: getPlayerStatusColor(player.verificationStatus),
                    });
                });
            }

            setResults(allResults);
        } catch (error) {
            logger.error('Global search error:', error);
        } finally {
            setLoading(false);
        }
    }, []);

    // 处理输入变化（带防抖）
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setKeyword(value);
        
        // 清除之前的定时器
        if (searchTimeoutRef.current) {
            clearTimeout(searchTimeoutRef.current);
        }
        
        // 设置新的防抖定时器
        searchTimeoutRef.current = setTimeout(() => {
            doSearch(value);
        }, 300);
    };

    // 清理定时器
    useEffect(() => {
        return () => {
            if (searchTimeoutRef.current) {
                clearTimeout(searchTimeoutRef.current);
            }
        };
    }, []);

    // 处理结果点击
    const handleResultClick = (result: SearchResult) => {
        onClose();
        switch (result.type) {
            case 'order':
                navigate('/admin/orders');
                break;
            case 'user':
                navigate('/admin/users');
                break;
            case 'player':
                navigate('/admin/players');
                break;
            case 'dispute':
                navigate('/admin/disputes');
                break;
        }
    };

    // 过滤结果
    const filteredResults = activeTab === 'all' 
        ? results 
        : results.filter(r => r.type === activeTab);

    // 获取类型图标
    const getTypeIcon = (type: string) => {
        switch (type) {
            case 'order': return <ShoppingOutlined />;
            case 'user': return <UserOutlined />;
            case 'player': return <TeamOutlined />;
            case 'dispute': return <FileTextOutlined />;
            default: return <SearchOutlined />;
        }
    };

    // 获取类型标签
    const getTypeLabel = (type: string) => {
        switch (type) {
            case 'order': return '订单';
            case 'user': return '用户';
            case 'player': return '陪玩师';
            case 'dispute': return '纠纷';
            default: return type;
        }
    };

    return (
        <Modal
            title={null}
            open={open}
            onCancel={onClose}
            footer={null}
            width={600}
            styles={{ body: { padding: 0 } }}
            closable={false}
        >
            <div style={{ padding: '16px 16px 0' }}>
                <Input
                    ref={inputRef}
                    size="large"
                    placeholder="搜索订单号、用户名、陪玩师..."
                    prefix={<SearchOutlined />}
                    value={keyword}
                    onChange={handleInputChange}
                    allowClear
                />
            </div>

            {keyword && (
                <Tabs
                    activeKey={activeTab}
                    onChange={setActiveTab}
                    style={{ padding: '0 16px' }}
                    items={[
                        { key: 'all', label: `全部 (${results.length})` },
                        { key: 'order', label: `订单 (${results.filter(r => r.type === 'order').length})` },
                        { key: 'user', label: `用户 (${results.filter(r => r.type === 'user').length})` },
                        { key: 'player', label: `陪玩师 (${results.filter(r => r.type === 'player').length})` },
                    ]}
                />
            )}

            <div style={{ maxHeight: 400, overflow: 'auto', padding: '0 16px 16px' }}>
                {loading ? (
                    <div style={{ textAlign: 'center', padding: 40 }}>
                        <Spin />
                    </div>
                ) : filteredResults.length > 0 ? (
                    <List
                        dataSource={filteredResults}
                        renderItem={(item) => (
                            <List.Item
                                style={{ cursor: 'pointer', padding: '12px 8px' }}
                                onClick={() => handleResultClick(item)}
                            >
                                <List.Item.Meta
                                    avatar={
                                        item.avatar ? (
                                            <Avatar src={item.avatar} />
                                        ) : (
                                            <Avatar icon={getTypeIcon(item.type)} />
                                        )
                                    }
                                    title={
                                        <Space>
                                            <span>{item.title}</span>
                                            <Tag>{getTypeLabel(item.type)}</Tag>
                                            {item.status && (
                                                <Tag color={item.statusColor}>{item.status}</Tag>
                                            )}
                                        </Space>
                                    }
                                    description={<Text type="secondary">{item.subtitle}</Text>}
                                />
                            </List.Item>
                        )}
                    />
                ) : keyword ? (
                    <Empty description="未找到相关结果" style={{ padding: 40 }} />
                ) : (
                    <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>
                        输入关键词开始搜索
                    </div>
                )}
            </div>

            <div style={{ 
                borderTop: '1px solid #f0f0f0', 
                padding: '8px 16px',
                color: '#999',
                fontSize: 12 
            }}>
                提示：按 ESC 关闭，按 ↑↓ 选择，按 Enter 确认
            </div>
        </Modal>
    );
};

// 辅助函数
function getOrderStatusColor(status?: string): string {
    const colors: Record<string, string> = {
        pending: 'gold',
        confirmed: 'blue',
        in_progress: 'processing',
        completed: 'success',
        canceled: 'default',
        refunded: 'error',
    };
    return colors[status || ''] || 'default';
}

function getPlayerStatusColor(status?: string): string {
    const colors: Record<string, string> = {
        pending: 'gold',
        verified: 'success',
        rejected: 'error',
    };
    return colors[status || ''] || 'default';
}

export default GlobalSearch;
