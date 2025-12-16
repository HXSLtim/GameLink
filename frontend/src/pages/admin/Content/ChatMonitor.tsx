/**
 * 聊天监控页面
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table, Card, Button, Space, Tag, Input, Select, DatePicker,
  Modal, message, Typography, Tooltip, Popconfirm, Form, InputNumber,
} from 'antd';
import {
  SearchOutlined, DeleteOutlined, StopOutlined,
  ReloadOutlined, CheckCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { chatModerationApi } from '@/api/content';
import type { ChatMessage, ChatMessageAuditStatus } from '@/types/content';
import { CHAT_AUDIT_STATUS_TEXT, CHAT_AUDIT_STATUS_COLOR } from '@/types/content';

const { RangePicker } = DatePicker;
const { TextArea } = Input;
const { Paragraph } = Typography;

const ChatMonitorPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [muteVisible, setMuteVisible] = useState(false);
  const [muteForm] = Form.useForm();
  const [currentMessage, setCurrentMessage] = useState<ChatMessage | null>(null);

  // 筛选条件
  const [groupId, setGroupId] = useState<number | undefined>();
  const [auditStatus, setAuditStatus] = useState<ChatMessageAuditStatus | ''>('');
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);

  const fetchMessages = useCallback(async () => {
    setLoading(true);
    try {
      const params: Record<string, unknown> = { page, pageSize };
      if (groupId) params.groupId = groupId;
      if (auditStatus) params.auditStatus = auditStatus;
      if (dateRange) {
        params.dateFrom = dateRange[0].format('YYYY-MM-DD');
        params.dateTo = dateRange[1].format('YYYY-MM-DD');
      }
      const res = await chatModerationApi.getMessages(params) as unknown as { success: boolean; data: { items: ChatMessage[]; total: number } };
      if (res.success) {
        setMessages(res.data?.items || []);
        setTotal(res.data?.total || 0);
      }
    } catch {
      message.error('获取消息列表失败');
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, groupId, auditStatus, dateRange]);

  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  const handleDelete = async (id: number) => {
    try {
      await chatModerationApi.deleteMessage(id);
      message.success('已删除');
      fetchMessages();
    } catch {
      message.error('操作失败');
    }
  };

  const handleMute = async () => {
    if (!currentMessage) return;
    try {
      const values = await muteForm.validateFields();
      await chatModerationApi.muteUser({
        groupId: currentMessage.groupId,
        userId: currentMessage.senderId,
        duration: values.duration,
        reason: values.reason,
      });
      message.success('已禁言');
      setMuteVisible(false);
      muteForm.resetFields();
    } catch {
      message.error('操作失败');
    }
  };

  const handleUnmute = async (groupId: number, userId: number) => {
    try {
      await chatModerationApi.unmuteUser(groupId, userId);
      message.success('已解除禁言');
    } catch {
      message.error('操作失败');
    }
  };

  // 高亮敏感词
  const highlightContent = (content: string, flaggedWords?: string[]) => {
    if (!flaggedWords || flaggedWords.length === 0) return content;
    let result = content;
    flaggedWords.forEach((word) => {
      const regex = new RegExp(word, 'gi');
      result = result.replace(regex, `<span style="color: red; font-weight: bold;">${word}</span>`);
    });
    return <span dangerouslySetInnerHTML={{ __html: result }} />;
  };

  const columns: ColumnsType<ChatMessage> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '群组',
      dataIndex: 'groupName',
      width: 120,
      render: (name, record) => name || `群组${record.groupId}`,
    },
    {
      title: '发送者',
      dataIndex: 'senderName',
      width: 120,
      render: (name, record) => name || `用户${record.senderId}`,
    },
    {
      title: '消息内容',
      dataIndex: 'content',
      ellipsis: true,
      render: (content, record) => (
        <Tooltip title={content}>
          <Paragraph ellipsis={{ rows: 2 }} style={{ marginBottom: 0 }}>
            {highlightContent(content, record.flaggedWords)}
          </Paragraph>
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'auditStatus',
      width: 100,
      render: (status: ChatMessageAuditStatus) => (
        <Tag color={CHAT_AUDIT_STATUS_COLOR[status]}>
          {CHAT_AUDIT_STATUS_TEXT[status]}
        </Tag>
      ),
    },
    {
      title: '敏感词',
      dataIndex: 'flaggedWords',
      width: 150,
      render: (words: string[]) =>
        words?.length > 0 ? (
          <Space wrap size="small">
            {words.map((word, idx) => (
              <Tag key={idx} color="red">{word}</Tag>
            ))}
          </Space>
        ) : '-',
    },
    {
      title: '发送时间',
      dataIndex: 'createdAt',
      width: 160,
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="禁言用户">
            <Button
              type="link"
              size="small"
              icon={<StopOutlined />}
              onClick={() => { setCurrentMessage(record); setMuteVisible(true); }}
            />
          </Tooltip>
          <Tooltip title="解除禁言">
            <Popconfirm
              title="确定解除该用户在此群组的禁言？"
              onConfirm={() => handleUnmute(record.groupId, record.senderId)}
            >
              <Button type="link" size="small" icon={<CheckCircleOutlined />} />
            </Popconfirm>
          </Tooltip>
          {record.auditStatus !== 'deleted' && (
            <Popconfirm title="确定删除此消息？" onConfirm={() => handleDelete(record.id)}>
              <Tooltip title="删除">
                <Button type="link" size="small" danger icon={<DeleteOutlined />} />
              </Tooltip>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="聊天监控"
      extra={<Button icon={<ReloadOutlined />} onClick={fetchMessages}>刷新</Button>}
    >
      <Space style={{ marginBottom: 16 }} wrap>
        <InputNumber
          placeholder="群组ID"
          value={groupId}
          onChange={(v) => setGroupId(v || undefined)}
          style={{ width: 120 }}
        />
        <Select
          placeholder="审核状态"
          value={auditStatus}
          onChange={(v) => { setAuditStatus(v); setPage(1); }}
          style={{ width: 120 }}
          allowClear
        >
          <Select.Option value="pending">待审核</Select.Option>
          <Select.Option value="approved">已通过</Select.Option>
          <Select.Option value="rejected">已拒绝</Select.Option>
          <Select.Option value="deleted">已删除</Select.Option>
        </Select>
        <RangePicker
          value={dateRange}
          onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)}
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={() => { setPage(1); fetchMessages(); }}>
          搜索
        </Button>
      </Space>

      <Table
        rowKey="id"
        columns={columns}
        dataSource={messages}
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (t) => `共 ${t} 条`,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
        scroll={{ x: 1100 }}
      />

      {/* 禁言弹窗 - 显示用户和违规消息 */}
      <Modal
        title="禁言用户"
        open={muteVisible}
        onOk={handleMute}
        onCancel={() => { setMuteVisible(false); muteForm.resetFields(); }}
      >
        {currentMessage && (
          <div style={{ background: '#fff2f0', border: '1px solid #ffccc7', padding: 12, borderRadius: 8, marginBottom: 16 }}>
            <div style={{ marginBottom: 8 }}>
              <Tag color="blue">{currentMessage.groupName || `群组${currentMessage.groupId}`}</Tag>
              <span style={{ fontWeight: 500 }}>{currentMessage.senderName || `用户${currentMessage.senderId}`}</span>
            </div>
            <Paragraph style={{ marginBottom: 8 }}>
              {highlightContent(currentMessage.content, currentMessage.flaggedWords)}
            </Paragraph>
            {currentMessage.flaggedWords && currentMessage.flaggedWords.length > 0 && (
              <div>
                <span style={{ fontSize: 12, color: '#999' }}>触发敏感词：</span>
                {currentMessage.flaggedWords.map((word, idx) => (
                  <Tag key={idx} color="red" style={{ marginLeft: 4 }}>{word}</Tag>
                ))}
              </div>
            )}
            <div style={{ fontSize: 12, color: '#999', marginTop: 8 }}>
              发送时间：{dayjs(currentMessage.createdAt).format('YYYY-MM-DD HH:mm:ss')}
            </div>
          </div>
        )}
        <Form form={muteForm} layout="vertical">
          <Form.Item
            name="duration"
            label="禁言时长"
            rules={[{ required: true, message: '请选择禁言时长' }]}
          >
            <Select placeholder="选择禁言时长">
              <Select.Option value={10}>10分钟</Select.Option>
              <Select.Option value={30}>30分钟</Select.Option>
              <Select.Option value={60}>1小时</Select.Option>
              <Select.Option value={360}>6小时</Select.Option>
              <Select.Option value={1440}>1天</Select.Option>
              <Select.Option value={4320}>3天</Select.Option>
              <Select.Option value={10080}>7天</Select.Option>
              <Select.Option value={43200}>30天</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="reason" label="禁言原因">
            <TextArea rows={3} placeholder="请输入禁言原因（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};

export default ChatMonitorPage;
