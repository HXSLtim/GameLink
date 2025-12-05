/**
 * 内容统计页面
 */
import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, message, Select, Table } from 'antd';
import {
  FileTextOutlined, MessageOutlined, WarningOutlined,
  CheckCircleOutlined, ClockCircleOutlined, CloseCircleOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { contentStatsApi } from '@/api/content';
import type { ContentStatsDTO, ContentTrend } from '@/types/content';

const StatsPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<ContentStatsDTO | null>(null);
  const [days, setDays] = useState(30);

  const fetchStats = async () => {
    setLoading(true);
    try {
      const res = await contentStatsApi.getStats(days);
      if (res.data.success) {
        setStats(res.data.data);
      }
    } catch {
      message.error('获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, [days]);

  const trendColumns: ColumnsType<ContentTrend> = [
    { title: '日期', dataIndex: 'date', key: 'date' },
    { title: '动态数', dataIndex: 'feedCount', key: 'feedCount' },
    { title: '消息数', dataIndex: 'messageCount', key: 'messageCount' },
    { title: '举报数', dataIndex: 'reportCount', key: 'reportCount' },
  ];

  return (
    <Spin spinning={loading}>
      <Row gutter={[16, 16]}>
        {/* 统计卡片 */}
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="总动态数"
              value={stats?.stats?.totalFeeds || 0}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="待审核动态"
              value={stats?.stats?.pendingFeeds || 0}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="已通过动态"
              value={stats?.stats?.approvedFeeds || 0}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="已拒绝动态"
              value={stats?.stats?.rejectedFeeds || 0}
              prefix={<CloseCircleOutlined />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="总消息数"
              value={stats?.stats?.totalMessages || 0}
              prefix={<MessageOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="标记消息数"
              value={stats?.stats?.flaggedMessages || 0}
              prefix={<WarningOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="总举报数"
              value={stats?.stats?.totalReports || 0}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="举报处理率"
              value={stats?.stats?.reportHandleRate || 0}
              suffix="%"
              precision={1}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>

        {/* 趋势表格 */}
        <Col span={24}>
          <Card
            title="内容趋势"
            extra={
              <Select value={days} onChange={setDays} style={{ width: 120 }}>
                <Select.Option value={7}>最近7天</Select.Option>
                <Select.Option value={14}>最近14天</Select.Option>
                <Select.Option value={30}>最近30天</Select.Option>
              </Select>
            }
          >
            <Table
              rowKey="date"
              columns={trendColumns}
              dataSource={stats?.trend || []}
              pagination={false}
              size="small"
            />
          </Card>
        </Col>
      </Row>
    </Spin>
  );
};

export default StatsPage;
