import React, { useState, useEffect } from 'react';
import { Modal, Input, Select } from '../../components';
import type {
  Player,
  CreatePlayerRequest,
  UpdatePlayerRequest,
  VerificationStatus,
} from '../../types/user';
import styles from './PlayerFormModal.module.less';

interface PlayerFormModalProps {
  visible: boolean;
  player?: Player | null;
  onClose: () => void;
  onSubmit: (data: CreatePlayerRequest | UpdatePlayerRequest) => Promise<void>;
}

export const PlayerFormModal: React.FC<PlayerFormModalProps> = ({
  visible,
  player,
  onClose,
  onSubmit,
}) => {
  const isEdit = !!player;
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState<CreatePlayerRequest | UpdatePlayerRequest>({
    userId: 0,
    nickname: '',
    bio: '',
    hourlyRateCents: 0,
    mainGameId: undefined,
  });

  useEffect(() => {
    if (player) {
      setFormData({
        nickname: player.nickname || '',
        bio: player.bio || '',
        hourlyRateCents: player.hourlyRateCents,
        mainGameId: player.mainGameId,
        verificationStatus: player.verificationStatus,
      });
    } else {
      setFormData({
        userId: 0,
        nickname: '',
        bio: '',
        hourlyRateCents: 0,
        mainGameId: undefined,
      });
    }
  }, [player]);

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
      title={isEdit ? '编辑陪玩师' : '新增陪玩师'}
      onClose={onClose}
      onOk={handleSubmit}
      onCancel={onClose}
      okText={loading ? '提交中...' : '确定'}
      cancelText="取消"
      width={600}
    >
      <div className={styles.form}>
        {!isEdit && (
          <div className={styles.formItem}>
            <label className={styles.label}>
              用户ID <span className={styles.required}>*</span>
            </label>
            <Input
              value={(formData as CreatePlayerRequest).userId?.toString() || ''}
              onChange={(e) => setFormData({ ...formData, userId: parseInt(e.target.value) || 0 })}
              placeholder="请输入用户ID"
              type="number"
            />
          </div>
        )}

        <div className={styles.formItem}>
          <label className={styles.label}>昵称</label>
          <Input
            value={formData.nickname || ''}
            onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
            placeholder="请输入昵称"
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>个人简介</label>
          <textarea
            className={styles.textarea}
            value={formData.bio || ''}
            onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
            placeholder="请输入个人简介"
            rows={4}
          />
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>
            时薪（分） <span className={styles.required}>*</span>
          </label>
          <Input
            value={formData.hourlyRateCents?.toString() || '0'}
            onChange={(e) =>
              setFormData({ ...formData, hourlyRateCents: parseInt(e.target.value) || 0 })
            }
            placeholder="请输入时薪（单位：分）"
            type="number"
          />
          <div className={styles.hint}>例如：5000 表示 50.00 元/小时</div>
        </div>

        <div className={styles.formItem}>
          <label className={styles.label}>主游戏ID</label>
          <Input
            value={formData.mainGameId?.toString() || ''}
            onChange={(e) =>
              setFormData({ ...formData, mainGameId: parseInt(e.target.value) || undefined })
            }
            placeholder="请输入主游戏ID"
            type="number"
          />
        </div>

        {isEdit && (
        <div className={styles.formItem}>
          <label className={styles.label}>认证状态</label>
          <Select
            value={(formData as UpdatePlayerRequest).verificationStatus || 'pending'}
            onChange={(value) =>
              setFormData({ ...formData, verificationStatus: value as VerificationStatus })
            }
            options={[
              { label: '待认证', value: 'pending' },
              { label: '已认证', value: 'verified' },
              { label: '已拒绝', value: 'rejected' },
            ]}
          />
        </div>
      )}
      </div>
    </Modal>
  );
};
