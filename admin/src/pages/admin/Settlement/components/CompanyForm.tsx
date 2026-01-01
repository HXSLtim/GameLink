/**
 * 结算公司创建/编辑表单
 */
import React, { useEffect } from 'react';
import {
    Modal,
    Form,
    Input,
    Select,
    Row,
    Col,
    message,
} from 'antd';
import { settlementApi, type SettlementCompany, type CompanyType } from '@/api/settlement';

export interface CompanyFormProps {
    /** 是否可见 */
    open: boolean;
    /** 编辑的公司数据（新建时为空） */
    company?: SettlementCompany | null;
    /** 确认回调 */
    onOk: () => void;
    /** 取消回调 */
    onCancel: () => void;
    /** 加载状态 */
    loading?: boolean;
}

/**
 * 公司类型选项
 */
const COMPANY_TYPE_OPTIONS: { label: string; value: CompanyType }[] = [
    { label: '个人', value: 'individual' },
    { label: '企业', value: 'company' },
];

/**
 * 结算公司表单组件
 */
const CompanyForm: React.FC<CompanyFormProps> = ({
    open,
    company,
    onOk,
    onCancel,
    loading: outerLoading,
}) => {
    const [form] = Form.useForm();
    const [loading, setLoading] = React.useState(false);
    const isEdit = !!company;

    useEffect(() => {
        if (open) {
            if (company) {
                // 编辑模式：填充表单
                form.setFieldsValue({
                    name: company.name,
                    type: company.type,
                    businessLicense: company.businessLicense,
                    taxNumber: company.taxNumber,
                    bankName: company.bankName,
                    bankAccount: company.bankAccount,
                    contactPerson: company.contactPerson,
                    contactPhone: company.contactPhone,
                });
            } else {
                // 新建模式：重置表单
                form.resetFields();
            }
        }
    }, [open, company, form]);

    /**
     * 提交表单
     */
    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();

            setLoading(true);
            if (isEdit && company) {
                // 更新
                const updateData = {
                    name: values.name,
                    type: values.type,
                    businessLicense: values.businessLicense,
                    taxNumber: values.taxNumber,
                    bankName: values.bankName,
                    bankAccount: values.bankAccount,
                    contactPerson: values.contactPerson,
                    contactPhone: values.contactPhone,
                };
                await settlementApi.updateSettlementCompany(company.id, updateData);
                message.success('更新结算公司成功');
            } else {
                // 新建
                const createData = {
                    name: values.name,
                    type: values.type,
                    businessLicense: values.businessLicense,
                    taxNumber: values.taxNumber,
                    bankName: values.bankName,
                    bankAccount: values.bankAccount,
                    contactPerson: values.contactPerson,
                    contactPhone: values.contactPhone,
                };
                await settlementApi.createSettlementCompany(createData);
                message.success('创建结算公司成功');
            }

            onOk();
        } catch (error) {
            console.error('Submit company form error:', error);
            message.error(isEdit ? '更新结算公司失败' : '创建结算公司失败');
        } finally {
            setLoading(false);
        }
    };

    /**
     * 公司类型变化时的处理
     */
    const handleTypeChange = (type: CompanyType) => {
        // 个人类型时清空企业相关字段
        if (type === 'individual') {
            form.setFieldsValue({
                businessLicense: undefined,
                taxNumber: undefined,
            });
        }
    };

    return (
        <Modal
            title={isEdit ? '编辑结算公司' : '新增结算公司'}
            open={open}
            onOk={handleSubmit}
            onCancel={onCancel}
            confirmLoading={loading || outerLoading}
            width={700}
            destroyOnClose
        >
            <Form
                form={form}
                layout="vertical"
                autoComplete="off"
            >
                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="name"
                            label="公司名称"
                            rules={[
                                { required: true, message: '请输入公司名称' },
                                { min: 2, max: 100, message: '公司名称长度为2-100个字符' },
                            ]}
                        >
                            <Input placeholder="请输入公司名称" maxLength={100} showCount />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="type"
                            label="公司类型"
                            rules={[{ required: true, message: '请选择公司类型' }]}
                        >
                            <Select
                                placeholder="请选择公司类型"
                                options={COMPANY_TYPE_OPTIONS}
                                onChange={handleTypeChange}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => prevValues.type !== currentValues.type}>
                    {({ getFieldValue }) => {
                        const type = getFieldValue('type');
                        return type === 'company' ? (
                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        name="businessLicense"
                                        label="营业执照号"
                                        rules={[
                                            { required: true, message: '请输入营业执照号' },
                                            { len: 18, message: '营业执照号为18位' },
                                        ]}
                                    >
                                        <Input placeholder="请输入18位营业执照号" maxLength={18} />
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        name="taxNumber"
                                        label="税务登记号"
                                        rules={[
                                            { required: true, message: '请输入税务登记号' },
                                        ]}
                                    >
                                        <Input placeholder="请输入税务登记号" maxLength={50} />
                                    </Form.Item>
                                </Col>
                            </Row>
                        ) : null;
                    }}
                </Form.Item>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="bankName"
                            label="开户银行"
                        >
                            <Input placeholder="请输入开户银行" maxLength={50} />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="bankAccount"
                            label="银行账号"
                            rules={[
                                { pattern: /^[0-9]{10,30}$/, message: '请输入10-30位数字账号' },
                            ]}
                        >
                            <Input placeholder="请输入银行账号" maxLength={30} />
                        </Form.Item>
                    </Col>
                </Row>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="contactPerson"
                            label="联系人"
                        >
                            <Input placeholder="请输入联系人姓名" maxLength={50} />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="contactPhone"
                            label="联系电话"
                            rules={[
                                { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' },
                            ]}
                        >
                            <Input placeholder="请输入联系电话" maxLength={11} />
                        </Form.Item>
                    </Col>
                </Row>
            </Form>
        </Modal>
    );
};

export default CompanyForm;
