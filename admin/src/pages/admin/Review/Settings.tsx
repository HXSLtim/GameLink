/**
 * 评价管理 - 展示设置页面
 * 需求: 7.1, 7.2, 7.3, 7.4, 7.5
 */
import React, { useState, useEffect } from 'react';
import {
  Card,
  Form,
  Radio,
  InputNumber,
  Switch,
  Button,
  Space,
  message,
  Spin,
  Descriptions,
  Divider,
  Typography,
} from 'antd';
import { SaveOutlined, ReloadOutlined, UndoOutlined } from '@ant-design/icons';
import { reviewSettingsApi } from '@/api/review';
import type { ReviewDisplaySettings, UpdateSettingsFormData } from '@/types/review';
import { SORT_BY_TEXT } from '@/types/review';

const { Text } = Typography;

// 默认设置
const DEFAULT_SETTINGS: UpdateSettingsFormData = {
  sortBy: 'time',
  minScore: 1,
  showAnonymous: true,
  pageSize: 10,
};

const ReviewSettingsPage: React.FC = () => {
  const [form] = Form.useForm<UpdateSettingsFormData>();

  // 状态
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [currentSettings, setCurrentSettings] = useState<ReviewDisplaySettings | null>(null);

  // 加载当前设置
  const fetchSettings = async () => {
    setLoading(true);
    try {
      const response = await reviewSettingsApi.getSettings();
      if (response.data.success) {
        setCurrentSettings(response.data.data);
        form.setFieldsValue({
          sortBy: response.data.data.sortBy,
          minScore: response.data.data.minScore,
          showAnonymous: response.data.data.showAnonymous,
          pageSize: response.data.data.pageSize,
        });
      }
    } catch {
      message.error('获取设置失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSettings();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // 保存设置
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);
      const response = await reviewSettingsApi.updateSettings(values);
      if (response.data.success) {
        message.success('设置已保存');
        setCurrentSettings(response.data.data);
      }
    } catch {
      message.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  // 重置为默认值
  const handleResetToDefault = () => {
    form.setFieldsValue(DEFAULT_SETTINGS);
    message.info('已重置为默认值，请点击保存生效');
  };

  // 恢复当前设置
  const handleRevert = () => {
    if (currentSettings) {
      form.setFieldsValue({
        sortBy: currentSettings.sortBy,
        minScore: currentSettings.minScore,
        showAnonymous: currentSettings.showAnonymous,
        pageSize: currentSettings.pageSize,
      });
      message.info('已恢复为当前保存的设置');
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 100 }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      {/* 当前设置预览 */}
      <Card title="当前设置" style={{ marginBottom: 16 }}>
        <Descriptions column={2}>
          <Descriptions.Item label="排序方式">
            {currentSettings ? SORT_BY_TEXT[currentSettings.sortBy] : '-'}
          </Descriptions.Item>
          <Descriptions.Item label="最低评分">
            {currentSettings?.minScore || '-'} 分
          </Descriptions.Item>
          <Descriptions.Item label="显示匿名评价">
            {currentSettings?.showAnonymous ? '是' : '否'}
          </Descriptions.Item>
          <Descriptions.Item label="每页显示数量">
            {currentSettings?.pageSize || '-'} 条
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {/* 设置表单 */}
      <Card title="修改设置">
        <Form
          form={form}
          layout="vertical"
          initialValues={DEFAULT_SETTINGS}
          style={{ maxWidth: 600 }}
        >
          <Form.Item
            name="sortBy"
            label="评价排序规则"
            tooltip="设置前端评价列表的默认排序方式"
          >
            <Radio.Group>
              <Radio.Button value="time">按时间</Radio.Button>
              <Radio.Button value="score">按评分</Radio.Button>
              <Radio.Button value="likes">按点赞数</Radio.Button>
            </Radio.Group>
          </Form.Item>

          <Form.Item
            name="minScore"
            label="最低评分阈值"
            tooltip="低于此评分的评价将不在前端展示"
            rules={[
              { required: true, message: '请输入最低评分' },
              { type: 'number', min: 1, max: 5, message: '评分范围为1-5' },
            ]}
          >
            <InputNumber min={1} max={5} style={{ width: 120 }} addonAfter="分" />
          </Form.Item>

          <Form.Item
            name="showAnonymous"
            label="显示匿名评价"
            tooltip="是否在前端展示匿名评价"
            valuePropName="checked"
          >
            <Switch checkedChildren="显示" unCheckedChildren="隐藏" />
          </Form.Item>

          <Form.Item
            name="pageSize"
            label="每页显示数量"
            tooltip="前端评价列表每页显示的评价数量"
            rules={[
              { required: true, message: '请输入每页数量' },
              { type: 'number', min: 5, max: 50, message: '数量范围为5-50' },
            ]}
          >
            <InputNumber min={5} max={50} style={{ width: 120 }} addonAfter="条" />
          </Form.Item>

          <Divider />

          <Form.Item>
            <Space>
              <Button
                type="primary"
                icon={<SaveOutlined />}
                onClick={handleSave}
                loading={saving}
              >
                保存设置
              </Button>
              <Button icon={<UndoOutlined />} onClick={handleRevert}>
                恢复当前
              </Button>
              <Button icon={<ReloadOutlined />} onClick={handleResetToDefault}>
                重置为默认
              </Button>
            </Space>
          </Form.Item>
        </Form>

        <Divider />

        <div>
          <Text type="secondary">
            说明：这些设置将影响用户端评价的展示方式。修改后需要点击"保存设置"才能生效。
          </Text>
        </div>
      </Card>
    </div>
  );
};

export default ReviewSettingsPage;
