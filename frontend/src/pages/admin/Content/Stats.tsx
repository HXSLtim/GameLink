/**
 * 内容统计页面
 */
import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Statistic, Spin, message, Select, Table, Button, Space } from 'antd';
import {
  FileTextOutlined, MessageOutlined, WarningOutlined,
  CheckCircleOutlined, ClockCircleOutlined, CloseCircleOutlined,
  DownloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { contentStatsApi } from '@/api/content';
import type { ContentStatsDTO, ContentTrend } from '@/types/content';

const StatsPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [exporting, setExporting] = useState(false);
  const [stats, setStats] = useState<ContentStatsDTO | null>(null);
  const [days, setDays] = useState(30);

  const fetchStats = async () => {
    setLoading(true);
    try {
      const res = await contentStatsApi.getStats(days) as unknown as { success: boolean; data: ContentStatsDTO };
      if (res.success) {
        setStats(res.data);
      }
    } catch {
      message.error('获取统计数据失败');
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async () => {
    setExporting(true);
    try {
      const res = await contentStatsApi.exportStats(days);
      // 从响应头获取文件名
      const contentDisposition = res.headers['content-disposition'];
      let filename = `content_stats_${days}days.xlsx`;
      if (contentDisposition) {
        const match = contentDisposition.match(/filename="?([^"]+)"?/);
        if (match) {
          filename = match[1];
        }
      }
      // 创建下载链接
      const blob = new Blob([res.data], {
        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      message.success('导出成功');
    } catch {
      message.error('导出失败');
    } finally {
      setExporting(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, [days, fetchStats]);

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
              styles={{ content: { color: '#faad14' } }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="已通过动态"
              value={stats?.stats?.approvedFeeds || 0}
              prefix={<CheckCircleOutlined />}
              styles={{ content: { color: '#52c41a' } }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={8} lg={6}>
          <Card>
            <Statistic
              title="已拒绝动态"
              value={stats?.stats?.rejectedFeeds || 0}
              prefix={<CloseCircleOutlined />}
              styles={{ content: { color: '#ff4d4f' } }}
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
              styles={{ content: { color: '#faad14' } }}
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
              styles={{ content: { color: '#52c41a' } }}
            />
          </Card>
        </Col>

        {/* 趋势表格 */}
        <Col span={24}>
          <Card
            title="内容趋势"
            extra={
              <Space>
                <Select value={days} onChange={setDays} style={{ width: 120 }}>
                  <Select.Option value={7}>最近7天</Select.Option>
                  <Select.Option value={14}>最近14天</Select.Option>
                  <Select.Option value={30}>最近30天</Select.Option>
                </Select>
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  loading={exporting}
                  onClick={handleExport}
                >
                  导出Excel
                </Button>
              </Space>
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
