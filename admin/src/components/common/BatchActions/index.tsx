/**
 * 统一的批量操作组件
 * 支持行内按钮和弹窗两种模式
 */
import React, { ReactNode, useState } from 'react';
import { Space, Button, Modal, Form } from 'antd';

export interface BatchAction {
    key: string;
    label: string;
    icon?: ReactNode;
    type?: 'default' | 'primary' | 'danger';
    mode?: 'inline' | 'modal';
    modalTitle?: string;
    modalContent?: ReactNode;
    onConfirm: (selectedKeys: React.Key[]) => void | Promise<void>;
}

export interface BatchActionsProps {
    /** 选中的行数 */
    selectedCount: number;
    /** 批量操作配置 */
    actions: BatchAction[];
    /** 选中的行 key */
    selectedRowKeys: React.Key[];
    /** 操作完成后回调 */
    onActionComplete?: () => void;
}

const BatchActions: React.FC<BatchActionsProps> = ({
    selectedCount,
    actions,
    selectedRowKeys,
    onActionComplete,
}) => {
    const [modalVisible, setModalVisible] = useState(false);
    const [currentAction, setCurrentAction] = useState<BatchAction | null>(null);
    const [form] = Form.useForm();

    const handleActionClick = (action: BatchAction) => {
        if (action.mode === 'modal') {
            setCurrentAction(action);
            setModalVisible(true);
        } else {
            action.onConfirm(selectedRowKeys).then(() => {
                onActionComplete?.();
            });
        }
    };

    const handleModalOk = async () => {
        if (currentAction) {
            try {
                await form.validateFields();
                await currentAction.onConfirm(selectedRowKeys);
                setModalVisible(false);
                form.resetFields();
                onActionComplete?.();
            } catch (error) {
                // 表单验证失败
            }
        }
    };

    const handleModalCancel = () => {
        setModalVisible(false);
        setCurrentAction(null);
        form.resetFields();
    };

    if (selectedCount === 0) {
        return null;
    }

    return (
        <>
            <Space>
                {actions.map((action) => {
                    const button = (
                        <Button
                            key={action.key}
                            type={action.type || 'default'}
                            icon={action.icon}
                            danger={action.type === 'danger'}
                            onClick={() => handleActionClick(action)}
                        >
                            {action.label}
                        </Button>
                    );
                    return button;
                })}
            </Space>

            <Modal
                title={currentAction?.modalTitle}
                open={modalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                destroyOnClose
            >
                {currentAction?.modalContent}
            </Modal>
        </>
    );
};

export default BatchActions;
