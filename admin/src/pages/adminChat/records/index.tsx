/**
 * 聊天记录管理页面
 * 管理所有聊天消息记录，支持按会话、发送者、类型筛选
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Card,
  Row,
  Col,
  Statistic,
  Tag,
  Space,
  Button,
  App,
  Popconfirm,
  Descriptions,
  Tag as AntTag,
  Typography,
  Avatar,
  Drawer,
  Alert,
} from 'antd';
import {
  ReloadOutlined,
  DeleteOutlined,
  EyeOutlined,
  UserOutlined,
  DownloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { PageContainer, SearchTable, type ToolbarButton, type SearchField } from '@/components';
import { chatMessageApi } from '@/api/chat';
import type {
  ChatMessage,
  ChatMessageQueryParams,
  MessageType,
  MessageSenderType,
  ChatStatsOverview,
} from '@/api/chat';
import { MESSAGE_TYPE_TEXT, SENDER_TYPE_TEXT, SENDER_TYPE_COLOR } from '@/api/chat';
import { exportToCSV, type ExportColumn } from '@/utils/export';
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';

const { Text } = Typography;

/**
 * 聊天消息查询参数接口
 */
interface MessageQueryParams extends ChatMessageQueryParams {
  conversationId?: number;
  senderId?: number;
  senderType?: MessageSenderType;
  messageType?: MessageType;
  keyword?: string;
  dateFrom?: string;
  dateTo?: string;
}

/**
 * 聊天记录管理页面
 */
const ChatRecordsPage: React.FC = () => {
  const { message, modal } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [total, setTotal] = useState(0);
  const [current, setCurrent] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchParams, setSearchParams] = useState<MessageQueryParams>({});
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);

  // Statistics
  const [stats, setStats] = useState<ChatStatsOverview | null>(null);

  // Detail drawer
  const [detailVisible, setDetailVisible] = useState(false);
  const [currentMessage, setCurrentMessage] = useState<ChatMessage | null>(null);

  /**
   * 加载统计数据
   */
  const loadStats = useCallback(async () => {
    try {
      const response = await chatMessageApi.getAllMessages({ pageSize: 1 });
      if (response.data?.success) {
        const msgTotal = response.data.data?.total || 0;
        setStats({
          totalConversations: 0,
          activeConversations: 0,
          totalMessages: msgTotal,
          todayMessages: 0,
          totalUsers: 0,
          onlineUsers: 0,
        });
      }
    } catch (error) {
      logger.error('Load stats error:', error);
    }
  }, []);

  /**
   * 加载消息列表
   */
  const loadMessages = useCallback(async () => {
    setLoading(true);
    try {
      const response = await chatMessageApi.getAllMessages({
        page: current,
        pageSize,
        ...searchParams,
      });
      if (response.data?.success) {
        setMessages(response.data.data.items || []);
        setTotal(response.data.data.total || 0);
      } else {
        message.error(response.data?.message || '加载失败');
      }
    } catch (error) {
      logger.error('Load messages error:', error);
      message.error('加载消息列表失败');
    } finally {
      setLoading(false);
    }
  }, [current, pageSize, searchParams, message]);

  useEffect(() => {
    loadStats();
  }, [loadStats]);

  useEffect(() => {
    loadMessages();
  }, [loadMessages]);

  /**
   * 搜索处理
   */
  const handleSearch = (values: Record<string, unknown>) => {
    setSearchParams(values as MessageQueryParams);
    setCurrent(1);
  };

  /**
   * 查看详情
   */
  const handleViewDetail = (record: ChatMessage) => {
    setCurrentMessage(record);
    setDetailVisible(true);
  };

  /**
   * 删除消息
   */
  const handleDelete = (record: ChatMessage) => {
    modal.confirm({
      title: '确认删除',
      content: (
        <div>
          <p>确定要删除这条消息吗？</p>
          <p style={{ color: '#ff4d4f' }}>此操作不可恢复，删除后用户端将无法查看该消息。</p>
        </div>
      ),
      onOk: async () => {
        try {
          await chatMessageApi.deleteMessage(record.id, '管理员删除');
          message.success('删除成功');
          loadMessages();
          loadStats();
        } catch (error) {
          logger.error('Delete message error:', error);
          message.error('删除失败');
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
      const response = await chatMessageApi.getAllMessages({
        ...searchParams,
        pageSize: 10000,
      });
      if (response.data?.success && response.data.data) {
        const exportColumns: ExportColumn<ChatMessage>[] = [
          { key: 'id', title: 'ID' },
          { key: 'conversationId', title: '会话ID', format: (_, r) => r.conversationId || '-' },
          { key: 'senderName', title: '发送者' },
          { key: 'senderType', title: '发送者类型', format: (_, r) => SENDER_TYPE_TEXT[r.senderType] || '-' },
          { key: 'content', title: '消息内容' },
          { key: 'messageType', title: '消息类型', format: (_, r) => MESSAGE_TYPE_TEXT[r.messageType] || '-' },
          { key: 'createdAt', title: '发送时间', format: (v) => dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') },
        ];
        exportToCSV(
          response.data.data.items as unknown as Record<string, unknown>[],
          exportColumns,
          'chat_messages'
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
   * 快速筛选
   */
  const handleQuickFilter = (filterType: string) => {
    if (filterType === 'today') {
      const today = dayjs().format('YYYY-MM-DD');
      setSearchParams({ dateFrom: today, dateTo: today });
    } else if (filterType === 'week') {
      const today = dayjs().format('YYYY-MM-DD');
      const weekAgo = dayjs().subtract(7, 'day').format('YYYY-MM-DD');
      setSearchParams({ dateFrom: weekAgo, dateTo: today });
    } else if (filterType === 'text') {
      setSearchParams({ messageType: 'text' as MessageType });
    } else if (filterType === 'image') {
      setSearchParams({ messageType: 'image' as MessageType });
    } else if (filterType === 'system') {
      setSearchParams({ messageType: 'system' as MessageType });
    } else {
      setSearchParams({});
    }
    setCurrent(1);
  };

  /**
   * 搜索字段配置
   */
  const searchFields: SearchField[] = [
    {
      name: 'keyword',
      label: '关键词',
      type: 'input',
      placeholder: '搜索消息内容、发送者',
    },
    {
      name: 'messageType',
      label: '消息类型',
      type: 'select',
      options: Object.entries(MESSAGE_TYPE_TEXT).map(([key, val]) => ({ label: val, value: key })),
    },
    {
      name: 'senderType',
      label: '发送者类型',
      type: 'select',
      options: Object.entries(SENDER_TYPE_TEXT).map(([key, val]) => ({ label: val, value: key })),
    },
    {
      name: 'conversationId',
      label: '会话ID',
      type: 'input',
      placeholder: '输入会话ID',
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
  const columns: ColumnsType<ChatMessage> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 70,
    },
    {
      title: '会话ID',
      dataIndex: 'conversationId',
      key: 'conversationId',
      width: 90,
    },
    {
      title: '发送者',
      key: 'sender',
      width: 150,
      render: (_: ChatMessage, record: ChatMessage) => (
        <Space>
          <Avatar
            size={24}
            src={record.senderAvatar}
            icon={<UserOutlined />}
          />
          <span>{record.senderName}</span>
          <AntTag color={SENDER_TYPE_COLOR[record.senderType]} style={{ fontSize: 11 }}>
            {SENDER_TYPE_TEXT[record.senderType]}
          </AntTag>
        </Space>
      ),
    },
    {
      title: '消息类型',
      dataIndex: 'messageType',
      key: 'messageType',
      width: 90,
      render: (type: MessageType) => (
        <Tag color={type === 'system' ? 'orange' : type === 'image' ? 'blue' : 'default'}>
          {MESSAGE_TYPE_TEXT[type] || type}
        </Tag>
      ),
    },
    {
      title: '消息内容',
      dataIndex: 'content',
      key: 'content',
      width: 300,
      ellipsis: { showTitle: true },
      render: (content: string, record: ChatMessage) => {
        if (record.messageType === 'image' && record.imageUrl) {
          return (
            <Space>
              <img
                src={record.imageUrl}
                alt={content}
                style={{ maxWidth: 200, maxHeight: 100 }}
              />
            </Space>
          );
        }
        return content || '-';
      },
    },
    {
      title: '发送时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 160,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right',
      render: (_: ChatMessage, record: ChatMessage) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
          <Popconfirm
            title="确定要删除这条消息吗？"
            onConfirm={() => handleDelete(record)}
          >
            <Button
              type="link"
              size="small"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <PageContainer title="聊天记录管理" subTitle="查看和管理所有聊天消息记录">
      {/* Statistics */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总消息数"
              value={stats?.totalMessages || 0}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="今日消息"
              value={stats?.todayMessages || 0}
              />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="活跃会话数"
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
      </Row>

      {/* Quick Filters */}
      <Card style={{ marginBottom: 16 }}>
        <Space wrap>
          <Text>快速筛选：</Text>
          <Button size="small" onClick={() => handleQuickFilter('all')}>
            全部
          </Button>
          <Button size="small" onClick={() => handleQuickFilter('today')}>
            今天
          </Button>
          <Button size="small" onClick={() => handleQuickFilter('week')}>
            最近7天
          </Button>
          <Button size="small" onClick={() => handleQuickFilter('text')}>
            文本消息
          </Button>
          <Button size="small" onClick={() => handleQuickFilter('image')}>
            图片消息
          </Button>
          <Button size="small" onClick={() => handleQuickFilter('system')}>
            系统消息
          </Button>
          <Button size="small" icon={<ReloadOutlined />} onClick={() => loadMessages()}>
            刷新
          </Button>
        </Space>
      </Card>

      {/* Message Table */}
      <SearchTable
        columns={columns}
        dataSource={messages}
        rowKey="id"
        searchFields={searchFields}
        onSearch={handleSearch}
        onRefresh={() => loadMessages()}
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
        scroll={{ x: 1400 }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
      />

      {/* Detail Drawer */}
      <Drawer
        title="消息详情"
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
        width={500}
      >
        {currentMessage && (
          <>
            <Descriptions column={1} size="small" bordered style={{ marginBottom: 16 }}>
              <Descriptions.Item label="消息ID">{currentMessage.id}</Descriptions.Item>
              <Descriptions.Item label="会话ID">{currentMessage.conversationId}</Descriptions.Item>
              <Descriptions.Item label="发送者">
                <Space>
                  <Avatar
                    size={32}
                    src={currentMessage.senderAvatar}
                    icon={<UserOutlined />}
                  />
                  <span>{currentMessage.senderName}</span>
                  <AntTag color={SENDER_TYPE_COLOR[currentMessage.senderType]} style={{ fontSize: 12 }}>
                    {SENDER_TYPE_TEXT[currentMessage.senderType]}
                  </AntTag>
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="消息类型">
                <Tag color={currentMessage.messageType === 'system' ? 'orange' : currentMessage.messageType === 'image' ? 'blue' : 'default'}>
                  {MESSAGE_TYPE_TEXT[currentMessage.messageType] || currentMessage.messageType}
                </Tag>
              </Descriptions.Item>
              <Descriptions.Item label="发送时间">
                {dayjs(currentMessage.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              <Descriptions.Item label="是否删除">
                <Tag color={currentMessage.isDeleted ? 'red' : 'green'}>
                  {currentMessage.isDeleted ? '已删除' : '正常'}
                </Tag>
              </Descriptions.Item>
              {currentMessage.replyToId && (
                <Descriptions.Item label="回复消息ID">{currentMessage.replyToId}</Descriptions.Item>
              )}
              <Descriptions.Item label="消息内容" span={2}>
                <Text style={{ wordBreak: 'break-all' }}>
                  {currentMessage.content}
                </Text>
              </Descriptions.Item>
              {currentMessage.imageUrl && (
                <Descriptions.Item label="图片链接" span={2}>
                  <a
                    href={currentMessage.imageUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {currentMessage.imageUrl}
                  </a>
                </Descriptions.Item>
              )}
            </Descriptions>

            {currentMessage.isDeleted && (
              <Alert
                message="此消息已被删除"
                description="用户端无法查看该消息"
                type="error"
                showIcon
                style={{ marginTop: 16 }}
              />
            )}
          </>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default ChatRecordsPage;
