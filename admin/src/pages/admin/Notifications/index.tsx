import React, { useEffect, useState } from 'react';
import { List, Typography, Card, Button, Badge, Space, Empty, App, Popconfirm, theme } from 'antd';
import { CheckOutlined, BellOutlined, DeleteOutlined } from '@ant-design/icons';
import { userApi, type Notification, type ApiResponse, type NotificationListResponse } from '@/api/user';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

const AdminNotificationsPage: React.FC = () => {
    const { message } = App.useApp();
    const { token } = theme.useToken();
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);

    const fetchNotifications = async (currentPage: number, append = false) => {
        setLoading(true);
        try {
            const res = await userApi.getNotifications({ page: currentPage, page_size: 20 }) as unknown as ApiResponse<NotificationListResponse>;
            if (res.success && res.data) {
                const newNotifications = res.data.items || [];
                setNotifications(prev => append ? [...prev, ...newNotifications] : newNotifications);
                setHasMore(newNotifications.length === 20);
            }
        } catch {
            message.error('加载通知失败');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchNotifications(1);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleLoadMore = () => {
        const nextPage = page + 1;
        setPage(nextPage);
        fetchNotifications(nextPage, true);
    };

    const handleMarkAsRead = async (id: number) => {
        try {
            await userApi.markAsRead(id);
            setNotifications(prev => prev.map(n => n.id === id ? { ...n, isRead: true } : n));
        } catch {
            message.error('操作失败');
        }
    };

    const handleMarkAllRead = async () => {
        try {
            await userApi.markAllAsRead();
            setNotifications(prev => prev.map(n => ({ ...n, isRead: true })));
            message.success('已全部标记为已读');
        } catch {
            message.error('操作失败');
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await userApi.deleteNotification(id);
            setNotifications(prev => prev.filter(n => n.id !== id));
            message.success('通知已删除');
        } catch {
            message.error('删除失败');
        }
    };

    return (
        <div style={{ padding: '24px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                <Title level={4} style={{ margin: 0 }}>
                    <BellOutlined /> 消息通知
                </Title>
                <Button onClick={handleMarkAllRead}>全部已读</Button>
            </div>

            <List
                loading={loading && page === 1}
                dataSource={notifications}
                loadMore={
                    hasMore && !loading ? (
                        <div style={{ textAlign: 'center', marginTop: 12, height: 32, lineHeight: '32px' }}>
                            <Button onClick={handleLoadMore}>加载更多</Button>
                        </div>
                    ) : null
                }
                locale={{ emptyText: <Empty description="暂无通知" /> }}
                renderItem={(item) => (
                    <Card
                        style={{
                            marginBottom: 16,
                            opacity: item.isRead ? 0.7 : 1,
                            borderLeft: item.isRead
                                ? `1px solid ${token.colorBorder}`
                                : `4px solid ${token.colorPrimary}`
                        }}
                        styles={{ body: { padding: '16px 24px' } }}
                    >
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                            <div style={{ flex: 1 }}>
                                <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
                                    {!item.isRead && <Badge status="processing" style={{ marginRight: 8 }} />}
                                    <Text strong={!item.isRead} style={{ fontSize: 16 }}>{item.title}</Text>
                                    <Text type="secondary" style={{ marginLeft: 12, fontSize: 12 }}>
                                        {dayjs(item.createdAt).format('YYYY-MM-DD HH:mm')}
                                    </Text>
                                </div>
                                <div style={{ color: token.colorText, fontSize: '14px', lineHeight: '1.5' }}>
                                    {item.message}
                                </div>
                            </div>
                            <Space orientation="vertical" align="end">
                                {!item.isRead && (
                                    <Button
                                        type="text"
                                        icon={<CheckOutlined />}
                                        onClick={() => handleMarkAsRead(item.id)}
                                        title="标记为已读"
                                    />
                                )}
                                <Popconfirm
                                    title="删除通知"
                                    description="确定要删除这条通知吗？"
                                    onConfirm={() => handleDelete(item.id)}
                                    okText="是"
                                    cancelText="否"
                                >
                                    <Button
                                        type="text"
                                        danger
                                        icon={<DeleteOutlined />}
                                        title="删除"
                                    />
                                </Popconfirm>
                            </Space>
                        </div>
                    </Card>
                )}
            />
        </div>
    );
};

export default AdminNotificationsPage;
