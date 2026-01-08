/**
 * ImportModal Component
 * Modal for importing data from Excel/CSV files
 *
 * Features:
 * - File upload with drag-and-drop
 * - Template download button
 * - Preview table for parsed data
 * - Error display with row numbers
 *
 * @module components/ImportModal
 */
import React, { useState, useCallback } from 'react';
import {
  Modal,
  Upload,
  Button,
  Table,
  Alert,
  Space,
  Typography,
  Progress,
  message,
  Tabs,
  Tag,
  Tooltip,
  Result,
  Radio,
} from 'antd';
import type { UploadProps, TabsProps } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  InboxOutlined,
  DownloadOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  WarningOutlined,
  FileExcelOutlined,
} from '@ant-design/icons';
import type { ImportType } from '@/services/import/templates/types';
import type { ImportPreview, ParsedRow, ImportResult, DuplicateKeyHandling } from '@/services/import';
import { importService } from '@/services/import';

const { Dragger } = Upload;
const { Text, Title } = Typography;

/**
 * Import step enum
 */
type ImportStep = 'upload' | 'preview' | 'importing' | 'result';

/**
 * ImportModal props
 */
export interface ImportModalProps {
  /** Whether the modal is visible */
  open: boolean;
  /** Import type (user, player, game) */
  type: ImportType;
  /** Modal title */
  title?: string;
  /** Callback when modal is closed */
  onClose: () => void;
  /** Callback when import is successful */
  onSuccess?: (result: ImportResult) => void;
}

/**
 * Get type display name in Chinese
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
 * ImportModal Component
 */
export const ImportModal: React.FC<ImportModalProps> = ({
  open,
  type,
  title,
  onClose,
  onSuccess,
}) => {
  // State
  const [step, setStep] = useState<ImportStep>('upload');
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<ImportPreview | null>(null);
  const [_importing, setImporting] = useState(false);
  const [progress, setProgress] = useState(0);
  const [result, setResult] = useState<ImportResult | null>(null);
  const [duplicateHandling, setDuplicateHandling] = useState<DuplicateKeyHandling>('fail');

  // Reset state when modal closes
  const handleClose = useCallback(() => {
    setStep('upload');
    setFile(null);
    setPreview(null);
    setImporting(false);
    setProgress(0);
    setResult(null);
    setDuplicateHandling('fail');
    onClose();
  }, [onClose]);

  // Download template
  const handleDownloadTemplate = useCallback(() => {
    try {
      const blob = importService.downloadTemplate(type);
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `${getTypeDisplayName(type)}导入模板.csv`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
      message.success('模板下载成功');
    } catch {
      message.error('模板下载失败');
    }
  }, [type]);

  // Handle file upload
  const handleFileUpload = useCallback(async (uploadedFile: File) => {
    setFile(uploadedFile);
    
    try {
      const previewResult = await importService.parseFile(uploadedFile, type);
      setPreview(previewResult);
      setStep('preview');
    } catch (error) {
      message.error('文件解析失败');
      console.error('File parse error:', error);
    }
  }, [type]);

  // Upload props
  const uploadProps: UploadProps = {
    name: 'file',
    multiple: false,
    accept: '.xlsx,.xls,.csv',
    showUploadList: false,
    beforeUpload: (uploadedFile) => {
      // Validate file size (10MB max)
      const maxSize = 10 * 1024 * 1024;
      if (uploadedFile.size > maxSize) {
        message.error('文件大小不能超过 10MB');
        return false;
      }
      
      handleFileUpload(uploadedFile);
      return false; // Prevent auto upload
    },
  };

  // Execute import
  const handleImport = useCallback(async () => {
    if (!preview) return;

    setStep('importing');
    setImporting(true);
    setProgress(0);

    try {
      const validRows = preview.validRows;
      let importResult: ImportResult;

      const options = {
        duplicateKeyHandling: duplicateHandling,
        onProgress: (completed: number, total: number) => {
          setProgress(Math.round((completed / total) * 100));
        },
      };

      switch (type) {
        case 'user':
          importResult = await importService.importUsers(validRows, options);
          break;
        case 'player':
          importResult = await importService.importPlayers(validRows, options);
          break;
        case 'game':
          importResult = await importService.importGames(validRows, options);
          break;
        default:
          throw new Error(`Unsupported import type: ${type}`);
      }

      setResult(importResult);
      setStep('result');

      if (importResult.success) {
        message.success(`成功导入 ${importResult.importedCount} 条数据`);
        onSuccess?.(importResult);
      } else if (importResult.importedCount > 0) {
        message.warning(`部分导入成功：${importResult.importedCount} 条成功，${importResult.skippedCount} 条失败`);
        onSuccess?.(importResult);
      } else {
        message.error('导入失败');
      }
    } catch (error) {
      message.error('导入过程中发生错误');
      console.error('Import error:', error);
      setStep('preview');
    } finally {
      setImporting(false);
    }
  }, [preview, type, duplicateHandling, onSuccess]);

  // Get preview columns
  const getPreviewColumns = useCallback((): ColumnsType<ParsedRow> => {
    const template = importService.getTemplate(type);
    
    const columns: ColumnsType<ParsedRow> = [
      {
        title: '行号',
        dataIndex: 'rowNumber',
        key: 'rowNumber',
        width: 70,
        fixed: 'left',
      },
      ...template.columns.map((col) => ({
        title: col.labelZh,
        dataIndex: ['data', col.key],
        key: col.key,
        width: 120,
        ellipsis: true,
        render: (value: unknown) => {
          if (value === null || value === undefined || value === '') {
            return <Text type="secondary">-</Text>;
          }
          if (typeof value === 'boolean') {
            return value ? '是' : '否';
          }
          return String(value);
        },
      })),
      {
        title: '状态',
        key: 'status',
        width: 80,
        fixed: 'right',
        render: (_, record) => (
          record.isValid ? (
            <Tag color="success" icon={<CheckCircleOutlined />}>有效</Tag>
          ) : (
            <Tooltip title={record.errors.map(e => `${e.field}: ${e.message}`).join('\n')}>
              <Tag color="error" icon={<CloseCircleOutlined />}>无效</Tag>
            </Tooltip>
          )
        ),
      },
    ];

    return columns;
  }, [type]);

  // Render upload step
  const renderUploadStep = () => (
    <div style={{ padding: '20px 0' }}>
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        {/* Template download */}
        <Alert
          message="下载导入模板"
          description={
            <Space direction="vertical">
              <Text>请先下载模板文件，按照模板格式填写数据后上传。</Text>
              <Button
                type="primary"
                icon={<DownloadOutlined />}
                onClick={handleDownloadTemplate}
              >
                下载{getTypeDisplayName(type)}导入模板
              </Button>
            </Space>
          }
          type="info"
          showIcon
          icon={<FileExcelOutlined />}
        />

        {/* File upload */}
        <Dragger {...uploadProps}>
          <p className="ant-upload-drag-icon">
            <InboxOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p className="ant-upload-hint">
            支持 Excel (.xlsx, .xls) 和 CSV (.csv) 格式，文件大小不超过 10MB
          </p>
        </Dragger>
      </Space>
    </div>
  );

  // Render preview step
  const renderPreviewStep = () => {
    if (!preview) return null;

    const hasStructureErrors = preview.structureErrors.length > 0;
    const hasValidRows = preview.validRows.length > 0;
    const hasInvalidRows = preview.invalidRows.length > 0;

    // Tab items
    const tabItems: TabsProps['items'] = [];

    if (hasValidRows) {
      tabItems.push({
        key: 'valid',
        label: (
          <span>
            <CheckCircleOutlined style={{ color: '#52c41a' }} />
            有效数据 ({preview.validRows.length})
          </span>
        ),
        children: (
          <Table
            columns={getPreviewColumns()}
            dataSource={preview.validRows}
            rowKey="rowNumber"
            size="small"
            scroll={{ x: 'max-content', y: 300 }}
            pagination={{ pageSize: 10, showSizeChanger: false }}
          />
        ),
      });
    }

    if (hasInvalidRows) {
      tabItems.push({
        key: 'invalid',
        label: (
          <span>
            <CloseCircleOutlined style={{ color: '#ff4d4f' }} />
            无效数据 ({preview.invalidRows.length})
          </span>
        ),
        children: (
          <Table
            columns={getPreviewColumns()}
            dataSource={preview.invalidRows}
            rowKey="rowNumber"
            size="small"
            scroll={{ x: 'max-content', y: 300 }}
            pagination={{ pageSize: 10, showSizeChanger: false }}
            expandable={{
              expandedRowRender: (record) => (
                <div style={{ padding: '8px 0' }}>
                  {record.errors.map((error, index) => (
                    <Alert
                      key={index}
                      message={`${error.field}: ${error.message}`}
                      type="error"
                      showIcon
                      style={{ marginBottom: index < record.errors.length - 1 ? 8 : 0 }}
                    />
                  ))}
                </div>
              ),
              rowExpandable: (record) => record.errors.length > 0,
            }}
          />
        ),
      });
    }

    return (
      <div style={{ padding: '20px 0' }}>
        <Space direction="vertical" size="middle" style={{ width: '100%' }}>
          {/* File info */}
          <Alert
            message={`已选择文件: ${file?.name}`}
            type="info"
            showIcon
            action={
              <Button size="small" onClick={() => setStep('upload')}>
                重新选择
              </Button>
            }
          />

          {/* Structure errors */}
          {hasStructureErrors && (
            <Alert
              message="文件结构错误"
              description={
                <ul style={{ margin: 0, paddingLeft: 20 }}>
                  {preview.structureErrors.map((error, index) => (
                    <li key={index}>{error}</li>
                  ))}
                </ul>
              }
              type="error"
              showIcon
            />
          )}

          {/* Summary */}
          {!hasStructureErrors && (
            <Alert
              message={`共 ${preview.totalRows} 行数据，${preview.validRows.length} 行有效，${preview.invalidRows.length} 行无效`}
              type={hasInvalidRows ? 'warning' : 'success'}
              showIcon
              icon={hasInvalidRows ? <WarningOutlined /> : <CheckCircleOutlined />}
            />
          )}

          {/* Duplicate handling for game import */}
          {type === 'game' && hasValidRows && (
            <Alert
              message="重复游戏标识处理方式"
              description={
                <Radio.Group
                  value={duplicateHandling}
                  onChange={(e) => setDuplicateHandling(e.target.value)}
                  style={{ marginTop: 8 }}
                >
                  <Space direction="vertical">
                    <Radio value="fail">报错 - 遇到重复标识时停止导入该行</Radio>
                    <Radio value="skip">跳过 - 遇到重复标识时跳过该行</Radio>
                    <Radio value="update">更新 - 遇到重复标识时更新已有数据</Radio>
                  </Space>
                </Radio.Group>
              }
              type="info"
              showIcon
            />
          )}

          {/* Preview tabs */}
          {tabItems.length > 0 && (
            <Tabs items={tabItems} defaultActiveKey="valid" />
          )}
        </Space>
      </div>
    );
  };

  // Render importing step
  const renderImportingStep = () => (
    <div style={{ padding: '40px 0', textAlign: 'center' }}>
      <Space direction="vertical" size="large" align="center">
        <Title level={4}>正在导入数据...</Title>
        <Progress
          type="circle"
          percent={progress}
          status="active"
        />
        <Text type="secondary">请勿关闭此窗口</Text>
      </Space>
    </div>
  );

  // Render result step
  const renderResultStep = () => {
    if (!result) return null;

    const isSuccess = result.success;
    const isPartial = !result.success && result.importedCount > 0;

    return (
      <div style={{ padding: '20px 0' }}>
        <Result
          status={isSuccess ? 'success' : isPartial ? 'warning' : 'error'}
          title={isSuccess ? '导入成功' : isPartial ? '部分导入成功' : '导入失败'}
          subTitle={
            <Space direction="vertical">
              <Text>总计 {result.totalRows} 行</Text>
              <Text type="success">成功导入 {result.importedCount} 行</Text>
              {result.skippedCount > 0 && (
                <Text type="danger">跳过 {result.skippedCount} 行</Text>
              )}
            </Space>
          }
          extra={[
            <Button key="close" type="primary" onClick={handleClose}>
              完成
            </Button>,
            result.errors.length > 0 && (
              <Button key="retry" onClick={() => setStep('preview')}>
                查看详情
              </Button>
            ),
          ].filter(Boolean)}
        />

        {/* Error details */}
        {result.errors.length > 0 && (
          <div style={{ marginTop: 24 }}>
            <Title level={5}>错误详情</Title>
            <Table
              columns={[
                { title: '行号', dataIndex: 'rowNumber', key: 'rowNumber', width: 80 },
                { title: '字段', dataIndex: 'field', key: 'field', width: 120 },
                { title: '错误信息', dataIndex: 'message', key: 'message' },
              ]}
              dataSource={result.errors.map((e, i) => ({ ...e, key: i }))}
              size="small"
              pagination={{ pageSize: 5 }}
              scroll={{ y: 200 }}
            />
          </div>
        )}
      </div>
    );
  };

  // Render step content
  const renderStepContent = () => {
    switch (step) {
      case 'upload':
        return renderUploadStep();
      case 'preview':
        return renderPreviewStep();
      case 'importing':
        return renderImportingStep();
      case 'result':
        return renderResultStep();
      default:
        return null;
    }
  };

  // Modal footer
  const getFooter = () => {
    switch (step) {
      case 'upload':
        return [
          <Button key="cancel" onClick={handleClose}>
            取消
          </Button>,
        ];
      case 'preview':
        return [
          <Button key="cancel" onClick={handleClose}>
            取消
          </Button>,
          <Button key="back" onClick={() => setStep('upload')}>
            上一步
          </Button>,
          <Button
            key="import"
            type="primary"
            onClick={handleImport}
            disabled={!preview || preview.validRows.length === 0}
          >
            开始导入 ({preview?.validRows.length || 0} 条)
          </Button>,
        ];
      case 'importing':
        return null;
      case 'result':
        return null;
      default:
        return null;
    }
  };

  return (
    <Modal
      title={title || `导入${getTypeDisplayName(type)}数据`}
      open={open}
      onCancel={step === 'importing' ? undefined : handleClose}
      footer={getFooter()}
      width={800}
      maskClosable={step !== 'importing'}
      closable={step !== 'importing'}
      destroyOnHidden
    >
      {renderStepContent()}
    </Modal>
  );
};

export default ImportModal;
