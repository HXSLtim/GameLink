import React, { useEffect, useState } from 'react';
import { List, Typography, Card, Button, Badge, Space, Empty, message } from 'antd';
import { CheckOutlined, DeleteOutlined, BellOutlined } from '@ant-design/icons';
import { userApi, type Notification, type ApiResponse } from '@/api/user';
import dayjs from 'dayjs';

const { Title, Text, Paragraph } = Typography;

const NotificationsPage: React.FC = () => {
    const [notifications, setNotifications] = useState<Notification[]>([]);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [hasMore, setHasMore] = useState(true);

    const fetchNotifications = async (currentPage: number, append = false) => {
        setLoading(true);
        try {
            const res = await userApi.getNotifications({ page: currentPage, page_size: 20 }) as unknown as ApiResponse<Notification[]>;
            if (res.success) {
                const newNotifications = res.data || [];
                setNotifications(prev => append ? [...prev, ...newNotifications] : newNotifications);
                setHasMore(newNotifications.length === 20);
            }
        } catch (error) {
            message.error('Failed to load notifications');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchNotifications(1);
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
        } catch (error) {
            message.error('Operation failed');
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await userApi.deleteNotification(id);
            setNotifications(prev => prev.filter(n => n.id !== id));
            message.success('Deleted');
        } catch (error) {
            message.error('Delete failed');
        }
    };

    const handleMarkAllRead = async () => {
        try {
            await userApi.markAllAsRead();
            setNotifications(prev => prev.map(n => ({ ...n, isRead: true })));
            message.success('All marked as read');
        } catch (error) {
            message.error('Operation failed');
        }
    };

    return (
        <div style={{ maxWidth: 800, margin: '24px auto', padding: '0 24px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
                <Title level={2} style={{ margin: 0 }}>
                    <BellOutlined /> Notifications
                </Title>
                <Button onClick={handleMarkAllRead}>Mark all as read</Button>
            </div>

            <List
                loading={loading && page === 1}
                dataSource={notifications}
                loadMore={
                    hasMore && !loading ? (
                        <div style={{ textAlign: 'center', marginTop: 12, height: 32, lineHeight: '32px' }}>
                            <Button onClick={handleLoadMore}>Load More</Button>
                        </div>
                    ) : null
                }
                locale={{ emptyText: <Empty description="No notifications" /> }}
                renderItem={(item) => (
                    <Card
                        style={{
                            marginBottom: 16,
                            opacity: item.isRead ? 0.7 : 1,
                            borderLeft: item.isRead ? '1px solid #f0f0f0' : '4px solid #1890ff'
                        }}
                        bodyStyle={{ padding: '16px 24px' }}
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
                                <Paragraph style={{ margin: 0, color: 'var(--text-normal)' }}>
                                    {item.content}
                                </Paragraph>
                            </div>
                            <Space direction="vertical" align="end">
                                {!item.isRead && (
                                    <Button
                                        type="text"
                                        icon={<CheckOutlined />}
                                        onClick={() => handleMarkAsRead(item.id)}
                                        title="Mark as read"
                                    />
                                )}
                                <Button
                                    type="text"
                                    danger
                                    icon={<DeleteOutlined />}
                                    onClick={() => handleDelete(item.id)}
                                    title="Delete"
                                />
                            </Space>
                        </div>
                    </Card>
                )}
            />
        </div>
    );
};

export default NotificationsPage;
