/**
 * 数据导出组件
 * 支持多种格式导出和进度显示
 */
import React, { useState, useCallback } from 'react';
import {
  Button,
  Dropdown,
  Modal,
  Progress,
  Space,
  Typography,
  message,
  Radio,
} from 'antd';
import type { MenuProps } from 'antd';
import {
  DownloadOutlined,
  FileExcelOutlined,
  FileTextOutlined,
  FilePdfOutlined,
  LoadingOutlined,
} from '@ant-design/icons';
import styles from './index.module.css';

const { Text } = Typography;

type ExportFormat = 'xlsx' | 'csv' | 'pdf' | 'json';

interface ExportColumn {
  title: string;
  dataIndex: string;
  render?: (value: unknown, record: Record<string, unknown>) => string;
}

interface DataExportProps {
  /**
   * 导出数据（同步数据）
   */
  data?: Record<string, unknown>[];
  /**
   * 异步获取数据函数
   */
  fetchData?: () => Promise<Record<string, unknown>[]>;
  /**
   * 导出列配置
   */
  columns: ExportColumn[];
  /**
   * 文件名（不含扩展名）
   * @default 'export'
   */
  filename?: string;
  /**
   * 支持的导出格式
   * @default ['xlsx', 'csv']
   */
  formats?: ExportFormat[];
  /**
   * 按钮文本
   * @default '导出'
   */
  buttonText?: string;
  /**
   * 是否禁用
   */
  disabled?: boolean;
  /**
   * 按钮类型
   */
  buttonType?: 'default' | 'primary' | 'dashed' | 'link' | 'text';
  /**
   * 尺寸
   */
  size?: 'small' | 'middle' | 'large';
  /**
   * 导出前的确认
   */
  confirmBeforeExport?: boolean;
  /**
   * 最大导出条数
   */
  maxRecords?: number;
  /**
   * 导出成功回调
   */
  onSuccess?: (format: ExportFormat, count: number) => void;
  /**
   * 导出失败回调
   */
  onError?: (error: Error) => void;
}

// 格式配置
const FORMAT_CONFIG: Record<ExportFormat, { label: string; icon: React.ReactNode; mimeType: string }> = {
  xlsx: { label: 'Excel (.xlsx)', icon: <FileExcelOutlined />, mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' },
  csv: { label: 'CSV (.csv)', icon: <FileTextOutlined />, mimeType: 'text/csv' },
  pdf: { label: 'PDF (.pdf)', icon: <FilePdfOutlined />, mimeType: 'application/pdf' },
  json: { label: 'JSON (.json)', icon: <FileTextOutlined />, mimeType: 'application/json' },
};

export const DataExport: React.FC<DataExportProps> = ({
  data,
  fetchData,
  columns,
  filename = 'export',
  formats = ['xlsx', 'csv'],
  buttonText = '导出',
  disabled = false,
  buttonType = 'default',
  size = 'middle',
  confirmBeforeExport = false,
  maxRecords,
  onSuccess,
  onError,
}) => {
  const [exporting, setExporting] = useState(false);
  const [progress, setProgress] = useState(0);
  const [modalVisible, setModalVisible] = useState(false);
  const [selectedFormat, setSelectedFormat] = useState<ExportFormat>(formats[0]);

  // 转换数据为导出格式
  const transformData = useCallback((records: Record<string, unknown>[]) => {
    return records.map((record) => {
      const row: Record<string, unknown> = {};
      columns.forEach((col) => {
        const value = record[col.dataIndex];
        row[col.title] = col.render ? col.render(value, record) : value;
      });
      return row;
    });
  }, [columns]);

  // 导出为 CSV
  const exportCSV = useCallback((records: Record<string, unknown>[]) => {
    const headers = columns.map((col) => col.title);
    const rows = transformData(records);

    const csvContent = [
      headers.join(','),
      ...rows.map((row) =>
        headers.map((header) => {
          const value = row[header];
          // 处理包含逗号、引号或换行的值
          if (typeof value === 'string' && (value.includes(',') || value.includes('"') || value.includes('\n'))) {
            return `"${value.replace(/"/g, '""')}"`;
          }
          return value ?? '';
        }).join(',')
      ),
    ].join('\n');

    return new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8' });
  }, [columns, transformData]);

  // 导出为 JSON
  const exportJSON = useCallback((records: Record<string, unknown>[]) => {
    const rows = transformData(records);
    return new Blob([JSON.stringify(rows, null, 2)], { type: 'application/json' });
  }, [transformData]);

  // 执行导出
  const handleExport = useCallback(async (format: ExportFormat) => {
    setExporting(true);
    setProgress(0);

    try {
      // 获取数据
      let exportData: Record<string, unknown>[];
      if (fetchData) {
        setProgress(10);
        exportData = await fetchData();
        setProgress(50);
      } else if (data) {
        exportData = data;
        setProgress(50);
      } else {
        throw new Error('没有可导出的数据');
      }

      // 检查最大条数
      if (maxRecords && exportData.length > maxRecords) {
        exportData = exportData.slice(0, maxRecords);
        message.warning(`数据量超过限制，仅导出前 ${maxRecords} 条`);
      }

      setProgress(70);

      // 根据格式生成文件
      let blob: Blob;
      switch (format) {
        case 'csv':
          blob = exportCSV(exportData);
          break;
        case 'json':
          blob = exportJSON(exportData);
          break;
        case 'xlsx':
          // XLSX 需要额外的库支持，这里简化处理为 CSV
          blob = exportCSV(exportData);
          message.info('XLSX 格式暂不支持，已导出为 CSV');
          format = 'csv';
          break;
        case 'pdf':
          // PDF 需要额外的库支持
          throw new Error('PDF 导出功能暂未实现');
        default:
          throw new Error(`不支持的导出格式: ${format}`);
      }

      setProgress(90);

      // 下载文件
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${filename}_${new Date().toISOString().split('T')[0]}.${format}`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);

      setProgress(100);
      message.success(`导出成功，共 ${exportData.length} 条数据`);
      onSuccess?.(format, exportData.length);
    } catch (error) {
      const err = error instanceof Error ? error : new Error('导出失败');
      message.error(err.message);
      onError?.(err);
    } finally {
      setTimeout(() => {
        setExporting(false);
        setProgress(0);
        setModalVisible(false);
      }, 500);
    }
  }, [data, fetchData, filename, maxRecords, exportCSV, exportJSON, onSuccess, onError]);

  // 确认导出
  const confirmExport = useCallback((format: ExportFormat) => {
    if (confirmBeforeExport) {
      setSelectedFormat(format);
      setModalVisible(true);
    } else {
      handleExport(format);
    }
  }, [confirmBeforeExport, handleExport]);

  // 下拉菜单
  const menuItems: MenuProps['items'] = formats.map((format) => ({
    key: format,
    icon: FORMAT_CONFIG[format].icon,
    label: FORMAT_CONFIG[format].label,
    onClick: () => confirmExport(format),
  }));

  // 单一格式时直接显示按钮
  if (formats.length === 1) {
    return (
      <>
        <Button
          type={buttonType}
          size={size}
          icon={exporting ? <LoadingOutlined /> : <DownloadOutlined />}
          disabled={disabled || exporting}
          onClick={() => confirmExport(formats[0])}
        >
          {buttonText}
        </Button>

        <Modal
          title="确认导出"
          open={modalVisible}
          onOk={() => handleExport(selectedFormat)}
          onCancel={() => setModalVisible(false)}
          confirmLoading={exporting}
        >
          {exporting ? (
            <div className={styles.progress}>
              <Progress percent={progress} status={progress === 100 ? 'success' : 'active'} />
              <Text type="secondary">
                {progress < 50 ? '正在获取数据...' : progress < 90 ? '正在生成文件...' : '即将完成...'}
              </Text>
            </div>
          ) : (
            <Text>确定要导出数据吗？</Text>
          )}
        </Modal>
      </>
    );
  }

  // 多格式时显示下拉菜单
  return (
    <>
      <Dropdown menu={{ items: menuItems }} disabled={disabled || exporting}>
        <Button type={buttonType} size={size}>
          <Space>
            {exporting ? <LoadingOutlined /> : <DownloadOutlined />}
            {buttonText}
          </Space>
        </Button>
      </Dropdown>

      <Modal
        title="确认导出"
        open={modalVisible}
        onOk={() => handleExport(selectedFormat)}
        onCancel={() => setModalVisible(false)}
        confirmLoading={exporting}
      >
        {exporting ? (
          <div className={styles.progress}>
            <Progress percent={progress} status={progress === 100 ? 'success' : 'active'} />
            <Text type="secondary">
              {progress < 50 ? '正在获取数据...' : progress < 90 ? '正在生成文件...' : '即将完成...'}
            </Text>
          </div>
        ) : (
          <Space direction="vertical">
            <Text>选择导出格式:</Text>
            <Radio.Group value={selectedFormat} onChange={(e) => setSelectedFormat(e.target.value)}>
              <Space direction="vertical">
                {formats.map((format) => (
                  <Radio key={format} value={format}>
                    {FORMAT_CONFIG[format].icon} {FORMAT_CONFIG[format].label}
                  </Radio>
                ))}
              </Space>
            </Radio.Group>
          </Space>
        )}
      </Modal>
    </>
  );
};

export default DataExport;
