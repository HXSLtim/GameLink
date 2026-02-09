/**
 * 聊天室管理页面
 * 管理所有聊天会话，支持按类型、状态、用户筛选
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Row,
    Col,
    Statistic,
    Tag,
    Button,
    Space,
    Modal,
    Typography,
    Avatar,
    Popconfirm,
    Descriptions,
    Divider,
    App,
} from 'antd';
import {
    MessageOutlined,
    ReloadOutlined,
    EyeOutlined,
    ShoppingCartOutlined,
    UserOutlined,
    TeamOutlined,
    CloseCircleOutlined,
    DownloadOutlined,
} from '@ant-design/icons';
import { PageContainer, SearchTable, type ToolbarButton, type SearchField } from '@/components';
import { useAuthStore } from '@/stores/modules/authStore';
import { chatConversationApi } from '@/api/chat';
import type {
  ChatConversation,
  ChatConversationQueryParams,
  ChatConversationStatus,
  ChatConversationType,
  CloseConversationRequest,
  ChatStatsOverview,
} from '@/api/chat';
import {
  CONVERSATION_TYPE_TEXT,
  CONVERSATION_STATUS_TEXT,
  CONVERSATION_STATUS_COLOR,
} from '@/api/chat';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';

const { Text } = Typography;

/**
 * 聊天会话查询参数接口
 */
interface ConversationQueryParams extends ChatConversationQueryParams {
  type?: ChatConversationType;
  status?: ChatConversationStatus;
  userId?: number;
  playerId?: number;
  orderId?: number;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * 聊天室管理页面
 */
const ChatRoomsPage: React.FC = () => {
  const { message, modal } = App.useApp();
  const userInfo = useAuthStore((state) => state.userInfo);
  const [loading, setLoading] = useState(false);
  const [conversations, setConversations] = useState<ChatConversation[]>([]);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchParams, setSearchParams] = useState<ConversationQueryParams>({});
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Statistics
  const [stats, setStats] = useState<ChatStatsOverview | null>(null);

  // Detail modal
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentConversation, setCurrentConversation] = useState<ChatConversation | null>(null);

  /**
   * 加载统计数据
   */
  const loadStats = useCallback(async () => {
    try {
      const response = await chatConversationApi.getConversations({ pageSize: 1000 });
      if (response.data?.success && response.data.data) {
        const convs = response.data.data.items || [];
        const activeConvs = convs.filter(c => c.status === 'active');
        const userSet = new Set<number>();
        convs.forEach(c => {
          if (c.userId) userSet.add(c.userId);
        });
        setStats({
          totalConversations: convs.length,
          activeConversations: activeConvs.length,
          totalMessages: 0,
          todayMessages: 0,
          totalUsers: userSet.size,
          onlineUsers: 0,
        });
      }
    } catch (error) {
      logger.error('Load stats error:', error);
    }
  }, []);

  /**
   * 加载会话列表
   */
  const loadConversations = useCallback(async () => {
    setLoading(true);
    try {
      const response = await chatConversationApi.getConversations({
        page: current,
        pageSize,
        ...searchParams,
      });
      if (response.data?.success) {
        setConversations(response.data.data.items || []);
        setTotal(response.data.data.total || 0);
      } else {
        message.error(response.data?.message || '加载失败');
      }
    } catch (error) {
      logger.error('Load conversations error:', error);
      message.error('加载会话列表失败');
    } finally {
      setLoading(false);
    }
  }, [current, pageSize, searchParams, message]);

  useEffect(() => {
    loadStats();
  }, [loadStats]);

  useEffect(() => {
    loadConversations();
  }, [loadConversations]);

  /**
   * 搜索处理
   */
  const handleSearch = (values: Record<string, unknown>) => {
    setSearchParams(values as ConversationQueryParams);
    setCurrent(1);
  };

  /**
   * 查看详情
   */
  const handleViewDetail = (record: ChatConversation) => {
    setCurrentConversation(record);
    setDetailVisible(true);
  };

  /**
   * 关闭会话
   */
  const handleCloseConversation = async (record: ChatConversation, reason?: string) => {
    try {
      const data: CloseConversationRequest = {
        reason: reason || '管理员关闭',
        closedBy: userInfo?.id || 0,
      };
      await chatConversationApi.closeConversation(record.id, data);
      message.success('会话已关闭');
      loadConversations();
      loadStats();
      setDetailVisible(false);
    } catch (error) {
      logger.error('Close conversation error:', error);
      message.error('关闭会话失败');
    }
  };

  /**
   * 重新打开会话
   */
  const handleReopenConversation = async (record: ChatConversation) => {
    modal.confirm({
      title: '重新打开会话',
      content: '确定要重新打开这个会话吗？',
      onOk: async () => {
        try {
          await chatConversationApi.reopenConversation(record.id);
          message.success('会话已重新打开');
          loadConversations();
          loadStats();
          setDetailVisible(false);
        } catch (error) {
          logger.error('Reopen conversation error:', error);
          message.error('重新打开会话失败');
        }
      },
    });
  };

  /**
   * 导出数据
   */
  const handleExport = async () => {
    try {
      message.loading({ content: '正在导出...', key: 'export' });
      const response = await chatConversationApi.getConversations({
        ...searchParams,
        pageSize: 10000,
      });
      if (response.data?.success && response.data.data) {
        const exportColumns: ExportColumn<ChatConversation>[] = [
          { key: 'id', title: 'ID' },
          { key: 'type', title: '类型', format: (_, r) => CONVERSATION_TYPE_TEXT[r.type] || r.type },
          { key: 'orderId', title: '订单ID', format: (_, r) => r.orderNo || '-' },
          { key: 'userName', title: '用户', format: (_, r) => r.userName || '-' },
          { key: 'playerName', title: '陪玩师', format: (_, r) => r.playerName || '-' },
          { key: 'messageCount', title: '消息数' },
          { key: 'lastMessageContent', title: '最后消息', format: (_, r) => r.lastMessageContent?.substring(0, 50) + '...' || '-' },
          { key: 'status', title: '状态', format: (_, r) => CONVERSATION_STATUS_TEXT[r.status] || r.status },
          { key: 'createdAt', title: '创建时间', format: (v) => dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') },
        ];
        exportToCSV(
          response.data.data.items as unknown as Record<string, unknown>[],
          exportColumns,
          'chat_conversations'
        );
        message.success({ content: '导出成功', key: 'export' });
      } else {
        message.error({ content: response.data?.message || '导出失败', key: 'export' });
      }
    } catch {
      message.error({ content: '导出失败', key: 'export' });
    }
  };

  /**
   * 搜索字段配置
   */
  const searchFields: SearchField[] = [
    {
      name: 'keyword',
      label: '关键词',
      type: 'input',
      placeholder: '搜索用户名、陪玩师名、订单号',
    },
    {
      name: 'type',
      label: '会话类型',
      type: 'select',
      options: Object.entries(CONVERSATION_TYPE_TEXT).map(([key, val]) => ({ label: val, value: key })),
    },
    {
      name: 'status',
      label: '会话状态',
      type: 'select',
      options: Object.entries(CONVERSATION_STATUS_TEXT).map(([key, val]) => ({ label: val, value: key })),
    },
    {
      name: 'orderId',
      label: '订单ID',
      type: 'input',
      placeholder: '输入订单ID',
    },
  ];

  /**
   * 工具栏按钮
   */
  const toolbarButtons: ToolbarButton[] = [
    {
      text: '导出数据',
      icon: <DownloadOutlined />,
      needSelection: false,
      onClick: () => handleExport(),
    },
  ];

  /**
   * 表格列配置
   */
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
    },
    {
      title: '会话类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: ChatConversationType) => (
        <Tag
          color={type === 'user_order' ? 'blue' : 'green'}
          icon={type === 'user_order' ? <ShoppingCartOutlined /> : <TeamOutlined />}
        >
          {CONVERSATION_TYPE_TEXT[type] || type}
        </Tag>
      ),
    },
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      width: 140,
      render: (orderNo: string) => orderNo || '-',
    },
    {
      title: '用户',
      key: 'user',
      width: 150,
      render: (_: unknown, record: ChatConversation) => (
        <Space>
          <Avatar
            size={24}
            src={record.userAvatar}
            icon={<UserOutlined />}
          />
          <div>
            <div>{record.userName || `ID:${record.userId}`}</div>
            {record.userId && (
              <Text type="secondary" style={{ fontSize: 11 }}>
                ID:{record.userId}
              </Text>
            )}
          </div>
        </Space>
      ),
    },
    {
      title: '陪玩师',
      key: 'player',
      width: 150,
      render: (_: unknown, record: ChatConversation) => (
        <Space>
          <Avatar
            size={24}
            src={record.playerAvatar}
            icon={<UserOutlined />}
          />
          <div>
            <div>{record.playerName || `ID:${record.playerId}`}</div>
            {record.playerId && (
              <Text type="secondary" style={{ fontSize: 11 }}>
                ID:{record.playerId}
              </Text>
            )}
          </div>
        </Space>
      ),
    },
    {
      title: '消息数',
      dataIndex: 'messageCount',
      key: 'messageCount',
      width: 80,
    },
    {
      title: '最后消息',
      key: 'lastMessage',
      width: 200,
      render: (_: unknown, record: ChatConversation) => (
        <div>
          <Text ellipsis={{ tooltip: record.lastMessageContent }}>
            {record.lastMessageContent || '-'}
          </Text>
          {record.lastMessageAt && (
            <Text type="secondary" style={{ fontSize: 11 }}>
              {dayjs(record.lastMessageAt).format('MM-DD HH:mm')}
            </Text>
          )}
        </div>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 90,
      render: (status: ChatConversationStatus) => (
        <Tag color={CONVERSATION_STATUS_COLOR[status]}>
          {CONVERSATION_STATUS_TEXT[status] || status}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right' as const,
      render: (_: unknown, record: ChatConversation) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          {record.status === 'active' ? (
            <Popconfirm
              title="确定要关闭这个会话吗？"
              description="关闭后用户将无法发送新消息"
              onConfirm={() => handleCloseConversation(record)}
            >
              <Button type="link" size="small" danger icon={<CloseCircleOutlined />}>
                关闭
              </Button>
            </Popconfirm>
          ) : (
            <Popconfirm
              title="确定要重新打开这个会话吗？"
              onConfirm={() => handleReopenConversation(record)}
            >
              <Button type="link" size="small" icon={<ReloadOutlined />}>
                打开
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <PageContainer title="聊天室管理" subTitle="管理所有聊天会话">
      {/* Statistics */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总会话数"
              value={stats?.totalConversations || 0}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="活跃会话"
              value={stats?.activeConversations || 0}
              />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={stats?.totalUsers || 0}
              prefix={<UserOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="总消息数"
              value={stats?.totalMessages || 0}
              prefix={<MessageOutlined />}
            />
          </Card>
        </Col>
      </Row>

      {/* Quick Actions */}
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <Text>快捷操作：</Text>
          <Button size="small" onClick={() => setSearchParams({ status: 'active' as ChatConversationStatus })}>
            仅活跃会话
          </Button>
          <Button size="small" onClick={() => setSearchParams({ status: 'closed' as ChatConversationStatus })}>
            已关闭会话
          </Button>
          <Button
            size="small"
            onClick={() => {
              const today = dayjs().format('YYYY-MM-DD');
              setSearchParams({ dateFrom: today, dateTo: today });
            }}
          >
            今天创建的会话
          </Button>
          <Button
            size="small"
            onClick={() => {
              const today = dayjs().format('YYYY-MM-DD');
              setSearchParams({ status: 'active' as ChatConversationStatus, dateFrom: today, dateTo: today });
            }}
          >
            今天的活跃会话
          </Button>
          <Button size="small" icon={<ReloadOutlined />} onClick={() => loadConversations()}>
            刷新
          </Button>
        </Space>
      </Card>

      {/* Conversation Table */}
      <SearchTable
        columns={columns}
        dataSource={conversations}
        rowKey="id"
        searchFields={searchFields}
        onSearch={handleSearch}
        onRefresh={() => loadConversations()}
        loading={loading}
        toolbarButtons={toolbarButtons}
        pagination={{
          current,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (page, size) => {
            setCurrent(page);
            setPageSize(size);
          },
        }}
        scroll={{ x: 1600 }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
      />

      {/* Detail Modal */}
      <Modal
        title="会话详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={600}
        destroyOnHidden
      >
        {currentConversation && (
          <>
            <Descriptions column={2} size="small" bordered style={{ marginBottom: 16 }}>
              <Descriptions.Item label="会话ID">{currentConversation.id}</Descriptions.Item>
              <Descriptions.Item label="会话类型">
                <Tag
                  color={currentConversation.type === 'user_order' ? 'blue' : 'green'}
                >
                  {CONVERSATION_TYPE_TEXT[currentConversation.type] || currentConversation.type}
                </Tag>
              </Descriptions.Item>
              {currentConversation.orderNo && (
                <Descriptions.Item label="订单号">{currentConversation.orderNo}</Descriptions.Item>
              )}
              <Descriptions.Item label="用户">
                <Space>
                  <Avatar
                    size={32}
                    src={currentConversation.userAvatar}
                    icon={<UserOutlined />}
                  />
                  <div>
                    <div>{currentConversation.userName || `ID:${currentConversation.userId}`}</div>
                    {currentConversation.userId && (
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        ID:{currentConversation.userId}
                      </Text>
                    )}
                  </div>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="陪玩师">
                <Space>
                  <Avatar
                    size={32}
                    src={currentConversation.playerAvatar}
                    icon={<UserOutlined />}
                  />
                  <div>
                    <div>{currentConversation.playerName || `ID:${currentConversation.playerId}`}</div>
                    {currentConversation.playerId && (
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        ID:{currentConversation.playerId}
                      </Text>
                    )}
                  </div>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="消息数">
                <Statistic
                  value={currentConversation.messageCount}
                  suffix="条"
                  />
              </Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color={CONVERSATION_STATUS_COLOR[currentConversation.status]}>
                  {CONVERSATION_STATUS_TEXT[currentConversation.status] || currentConversation.status}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">
                {dayjs(currentConversation.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              {currentConversation.lastMessageAt && (
                <Descriptions.Item label="最后消息时间">
                  {dayjs(currentConversation.lastMessageAt).format('YYYY-MM-DD HH:mm:ss')}
                </Descriptions.Item>
              )}
            </Descriptions>

            {/* Last Message Preview */}
            {currentConversation.lastMessageContent && (
              <>
                <Divider />
                <Card size="small" title="最后消息预览">
                  <Text type="secondary">最后一条消息：</Text>
                  <Text code style={{ wordBreak: 'break-word' }}>
                    {currentConversation.lastMessageContent}
                  </Text>
                </Card>
              </>
            )}

            {/* Action Buttons */}
            <Divider />
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              {currentConversation.status === 'active' ? (
                <Button
                  danger
                  icon={<CloseCircleOutlined />}
                  onClick={() => handleCloseConversation(currentConversation)}
                >
                  关闭会话
                </Button>
              ) : (
                <Button
                  type="primary"
                  icon={<ReloadOutlined />}
                  onClick={() => handleReopenConversation(currentConversation)}
                >
                  重新打开会话
                </Button>
              )}
            </Space>
          </>
        )}
      </Modal>
    </PageContainer>
  );
};

export default ChatRoomsPage;
