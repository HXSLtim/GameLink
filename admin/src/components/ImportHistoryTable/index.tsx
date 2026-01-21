/**
 * ImportHistoryTable Component
 * Displays paginated list of past imports with filtering
 *
 * Features:
 * - Paginated list of past imports
 * - Filter by type and date range
 * - Link to import details
 *
 * @module components/ImportHistoryTable
 */
import React, { useState, useEffect, useCallback } from 'react';
import {
  Table,
  Card,
  Space,
  Tag,
  Button,
  Select,
  DatePicker,
  Typography,
  Tooltip,
  Modal,
  Descriptions,
  Progress,
  Alert,
  Empty,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  EyeOutlined,
  DownloadOutlined,
  ReloadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  LoadingOutlined,
  WarningOutlined,
  ClockCircleOutlined,
} from '@ant-design/icons';
import dayjs from 'dayjs';
import type { ImportType } from '@/services/import/templates/types';
import type {
  ImportHistoryRecord,
  ImportStatus,
  ImportHistoryPage,
} from '@/services/import/history/types';
import type { ImportHistoryDetails } from '@/services/import/history/historyService';
import { importHistoryService } from '@/services/import/history';
import { downloadErrorReport, hasErrorDetails } from '@/services/import/history/errorReport';

const { RangePicker } = DatePicker;
const { Text, Title } = Typography;

/**
 * ImportHistoryTable props
 */
export interface ImportHistoryTableProps {
  /** Filter by specific import type */
  filterType?: ImportType;
  /** Page size (default: 10) */
  pageSize?: number;
  /** Show type filter (default: true) */
  showTypeFilter?: boolean;
  /** Card title */
  title?: string;
  /** Callback when viewing details */
  onViewDetails?: (record: ImportHistoryRecord) => void;
}

/**
 * Get status display config
 */
function getStatusConfig(status: ImportStatus): {
  color: string;
  icon: React.ReactNode;
  text: string;
} {
  const configs: Record<ImportStatus, { color: string; icon: React.ReactNode; text: string }> = {
    pending: { color: 'default', icon: <ClockCircleOutlined />, text: '等待中' },
    processing: { color: 'processing', icon: <LoadingOutlined />, text: '处理中' },
    completed: { color: 'success', icon: <CheckCircleOutlined />, text: '已完成' },
    failed: { color: 'error', icon: <CloseCircleOutlined />, text: '失败' },
    partial: { color: 'warning', icon: <WarningOutlined />, text: '部分成功' },
  };
  return configs[status] || configs.pending;
}

/**
 * Get type display name
 */
function getTypeDisplayName(type: ImportType): string {
  const names: Record<ImportType, string> = {
    user: '用户',
    player: '陪玩师',
    game: '游戏',
  };
  return names[type] || type;
}

/**
 * ImportHistoryTable Component
 */
export const ImportHistoryTable: React.FC<ImportHistoryTableProps> = ({
  filterType,
  pageSize: defaultPageSize = 10,
  showTypeFilter = true,
  title = '导入历史',
  onViewDetails,
}) => {
  // State
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<ImportHistoryPage | null>(null);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(defaultPageSize);
  const [typeFilter, setTypeFilter] = useState<ImportType | undefined>(filterType);
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs] | null>(null);
  
  // Detail modal state
  const [detailModalVisible, setDetailModalVisible] = useState(false);
  const [selectedRecord, setSelectedRecord] = useState<ImportHistoryDetails | null>(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // Load data
  const loadData = useCallback(async () => {
    setLoading(true);
    try {
      const result = await importHistoryService.getImportHistory(
        {
          type: typeFilter,
          startDate: dateRange?.[0]?.toISOString(),
          endDate: dateRange?.[1]?.toISOString(),
        },
        { page, pageSize }
      );
      setData(result);
    } catch (error) {
      console.error('Failed to load import history:', error);
    } finally {
      setLoading(false);
    }
  }, [typeFilter, dateRange, page, pageSize]);

  // Load data on mount and filter change
  useEffect(() => {
    loadData();
  }, [loadData]);

  // Handle view details
  const handleViewDetails = useCallback(async (record: ImportHistoryRecord) => {
    if (onViewDetails) {
      onViewDetails(record);
      return;
    }

    setDetailLoading(true);
    setDetailModalVisible(true);
    
    try {
      const details = await importHistoryService.getImportDetails(record.id);
      setSelectedRecord(details);
    } catch (error) {
      console.error('Failed to load import details:', error);
    } finally {
      setDetailLoading(false);
    }
  }, [onViewDetails]);

  // Handle download error report
  const handleDownloadReport = useCallback((record: ImportHistoryRecord) => {
    try {
      downloadErrorReport(record);
    } catch (error) {
      console.error('Failed to download error report:', error);
    }
  }, []);

  // Table columns
  const columns: ColumnsType<ImportHistoryRecord> = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 100,
      ellipsis: true,
      render: (id: string) => (
        <Tooltip title={id}>
          <Text copyable={{ text: id }}>{id.substring(0, 8)}...</Text>
        </Tooltip>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
      render: (type: ImportType) => (
        <Tag>{getTypeDisplayName(type)}</Tag>
      ),
    },
    {
      title: '文件名',
      dataIndex: 'fileName',
      key: 'fileName',
      width: 200,
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: ImportStatus) => {
        const config = getStatusConfig(status);
        return (
          <Tag color={config.color} icon={config.icon}>
            {config.text}
          </Tag>
        );
      },
    },
    {
      title: '导入结果',
      key: 'result',
      width: 150,
      render: (_, record) => {
        const successRate = record.totalRows > 0
          ? Math.round((record.importedCount / record.totalRows) * 100)
          : 0;
        
        return (
          <Space orientation="vertical" size={0}>
            <Text>
              <Text type="success">{record.importedCount}</Text>
              {' / '}
              <Text>{record.totalRows}</Text>
              {record.skippedCount > 0 && (
                <Text type="danger"> ({record.skippedCount} 失败)</Text>
              )}
            </Text>
            <Progress
              percent={successRate}
              size="small"
              status={successRate === 100 ? 'success' : successRate > 0 ? 'normal' : 'exception'}
              showInfo={false}
            />
          </Space>
        );
      },
    },
    {
      title: '操作人',
      dataIndex: 'uploadedByName',
      key: 'uploadedByName',
      width: 100,
      ellipsis: true,
      render: (name: string | undefined, record) => name || `用户 ${record.uploadedBy}`,
    },
    {
      title: '导入时间',
      dataIndex: 'uploadedAt',
      key: 'uploadedAt',
      width: 160,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Tooltip title="查看详情">
            <Button
              type="link"
              size="small"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetails(record)}
            >
              详情
            </Button>
          </Tooltip>
          {hasErrorDetails(record) && (
            <Tooltip title="下载错误报告">
              <Button
                type="link"
                size="small"
                icon={<DownloadOutlined />}
                onClick={() => handleDownloadReport(record)}
              >
                报告
              </Button>
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  // Render detail modal content
  const renderDetailContent = () => {
    if (detailLoading) {
      return <div style={{ textAlign: 'center', padding: 40 }}><LoadingOutlined style={{ fontSize: 24 }} /></div>;
    }

    if (!selectedRecord) {
      return <Empty description="无法加载详情" />;
    }

    const statusConfig = getStatusConfig(selectedRecord.status);

    return (
      <Space orientation="vertical" size="middle" style={{ width: '100%' }}>
        <Descriptions column={2} bordered size="small">
          <Descriptions.Item label="导入ID">{selectedRecord.id}</Descriptions.Item>
          <Descriptions.Item label="类型">
            <Tag>{getTypeDisplayName(selectedRecord.type)}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="文件名">{selectedRecord.fileName}</Descriptions.Item>
          <Descriptions.Item label="文件大小">
            {(selectedRecord.fileSize / 1024).toFixed(2)} KB
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={statusConfig.color} icon={statusConfig.icon}>
              {statusConfig.text}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="成功率">
            <Progress
              percent={selectedRecord.successRate}
              size="small"
              status={selectedRecord.successRate === 100 ? 'success' : 'normal'}
            />
          </Descriptions.Item>
          <Descriptions.Item label="总行数">{selectedRecord.totalRows}</Descriptions.Item>
          <Descriptions.Item label="成功导入">{selectedRecord.importedCount}</Descriptions.Item>
          <Descriptions.Item label="跳过/失败">{selectedRecord.skippedCount}</Descriptions.Item>
          <Descriptions.Item label="错误数">{selectedRecord.errorCount}</Descriptions.Item>
          <Descriptions.Item label="导入时间">
            {dayjs(selectedRecord.uploadedAt).format('YYYY-MM-DD HH:mm:ss')}
          </Descriptions.Item>
          <Descriptions.Item label="耗时">
            {selectedRecord.durationFormatted || '-'}
          </Descriptions.Item>
        </Descriptions>

        {selectedRecord.errorSummary && (
          <Alert
            message="错误摘要"
            description={selectedRecord.errorSummary}
            type="error"
            showIcon
          />
        )}

        {selectedRecord.rowResults && selectedRecord.rowResults.length > 0 && (
          <>
            <Title level={5}>行级结果</Title>
            <Table
              columns={[
                { title: '行号', dataIndex: 'rowNumber', key: 'rowNumber', width: 80 },
                {
                  title: '状态',
                  dataIndex: 'success',
                  key: 'success',
                  width: 80,
                  render: (success: boolean) => (
                    success ? (
                      <Tag color="success" icon={<CheckCircleOutlined />}>成功</Tag>
                    ) : (
                      <Tag color="error" icon={<CloseCircleOutlined />}>失败</Tag>
                    )
                  ),
                },
                { title: '错误字段', dataIndex: 'errorField', key: 'errorField', width: 100 },
                { title: '错误信息', dataIndex: 'errorMessage', key: 'errorMessage', ellipsis: true },
              ]}
              dataSource={selectedRecord.rowResults.filter(r => !r.success)}
              rowKey="rowNumber"
              size="small"
              pagination={{ pageSize: 5 }}
              scroll={{ y: 200 }}
            />
          </>
        )}
      </Space>
    );
  };

  return (
    <>
      <Card
        title={title}
        extra={
          <Space>
            {showTypeFilter && !filterType && (
              <Select
                placeholder="选择类型"
                allowClear
                style={{ width: 120 }}
                value={typeFilter}
                onChange={setTypeFilter}
                options={[
                  { label: '用户', value: 'user' },
                  { label: '陪玩师', value: 'player' },
                  { label: '游戏', value: 'game' },
                ]}
              />
            )}
            <RangePicker
              value={dateRange}
              onChange={(dates) => setDateRange(dates as [dayjs.Dayjs, dayjs.Dayjs] | null)}
              placeholder={['开始日期', '结束日期']}
            />
            <Tooltip title="刷新">
              <Button
                icon={<ReloadOutlined spin={loading} />}
                onClick={loadData}
              />
            </Tooltip>
          </Space>
        }
      >
        <Table
          columns={columns}
          dataSource={data?.records || []}
          rowKey="id"
          loading={loading}
          pagination={{
            current: page,
            pageSize,
            total: data?.total || 0,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条`,
            onChange: (p, ps) => {
              setPage(p);
              setPageSize(ps);
            },
          }}
          scroll={{ x: 'max-content' }}
          size="middle"
        />
      </Card>

      {/* Detail Modal */}
      <Modal
        title="导入详情"
        open={detailModalVisible}
        onCancel={() => {
          setDetailModalVisible(false);
          setSelectedRecord(null);
        }}
        footer={[
          <Button key="close" onClick={() => setDetailModalVisible(false)}>
            关闭
          </Button>,
          selectedRecord && hasErrorDetails(selectedRecord) && (
            <Button
              key="download"
              type="primary"
              icon={<DownloadOutlined />}
              onClick={() => handleDownloadReport(selectedRecord)}
            >
              下载错误报告
            </Button>
          ),
        ].filter(Boolean)}
        width={800}
      >
        {renderDetailContent()}
      </Modal>
    </>
  );
};

export default ImportHistoryTable;
