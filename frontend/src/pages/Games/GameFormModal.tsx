import React, { useState, useEffect } from 'react';
import { Modal, Input, Select } from '../../components';
import type { Game, CreateGameRequest, UpdateGameRequest, GameCategory } from '../../types/game';
import styles from './GameFormModal.module.less';

interface GameFormModalProps {
  visible: boolean;
  game?: Game | null;
  onClose: () => void;
  onSubmit: (data: CreateGameRequest | UpdateGameRequest) => Promise<void>;
}

export const GameFormModal: React.FC<GameFormModalProps> = ({
  visible,
  game,
  onClose,
  onSubmit,
}) => {
  const isEdit = !!game;
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreateGameRequest | UpdateGameRequest>({
    key: '',
    name: '',
    category: 'other',
    iconUrl: '',
    description: '',
  });

  useEffect(() => {
    if (game) {
      setFormData({
        key: game.key,
        name: game.name,
        category: game.category,
        iconUrl: game.iconUrl || '',
        description: game.description || '',
      });
    } else {
      setFormData({
        key: '',
        name: '',
        category: 'other',
        iconUrl: '',
        description: '',
      });
    }
  }, [game]);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await onSubmit(formData);
      onClose();
    } catch (err) {
      console.error('提交失败:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      visible={visible}
      title={isEdit ? '编辑游戏' : '新增游戏'}
      onClose={onClose}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={loading ? '提交中...' : '确定'}
      cancelText="取消"
      width={600}
    >
      <div className={styles.form}>
        <div className={styles.formItem}>
          <label className={styles.label}>
            游戏KEY <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.key || ''}
            onChange={(e) => setFormData({ ...formData, key: e.target.value })}
            placeholder="请输入游戏KEY（英文字母、数字、下划线）"
            disabled={isEdit}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>
            游戏名称 <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.name || ''}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            placeholder="请输入游戏名称"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>分类</label>
          <Select
            value={formData.category || 'other'}
            onChange={(value) => setFormData({ ...formData, category: value as GameCategory })}
            options={[
              { label: 'MOBA', value: 'moba' },
              { label: '射击', value: 'fps' },
              { label: '角色扮演', value: 'rpg' },
              { label: '策略', value: 'strategy' },
              { label: '体育', value: 'sports' },
              { label: '竞速', value: 'racing' },
              { label: '益智', value: 'puzzle' },
              { label: '其他', value: 'other' },
            ]}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>图标URL</label>
          <Input
            value={formData.iconUrl || ''}
            onChange={(e) => setFormData({ ...formData, iconUrl: e.target.value })}
            placeholder="请输入图标URL"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>描述</label>
          <textarea
            className={styles.textarea}
            value={formData.description || ''}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="请输入游戏描述"
            rows={4}
          />
        </div>
      </div>
    </Modal>
  );
};
