import React, { useState } from 'react';
import { Card, Button, Tag, ActionButtons, DeleteConfirmModal } from '../../components';
import { formatDateTime } from '../../utils/formatters';
import { AnnouncementFormModal } from './AnnouncementFormModal';
import type { Announcement } from '../../types/settings';
import styles from './SettingsDashboard.module.less';

/**
 * 公告管理组件
 */
export const AnnouncementSettings: React.FC = () => {
  // 模拟数据
  const [announcements, setAnnouncements] = useState<Announcement[]>([
    {
      id: 1,
      title: '系统维护通知',
      content: '系统将于明天凌晨 2:00-4:00 进行维护，期间服务将暂停',
      type: 'warning',
      enabled: true,
      createdAt: '2024-01-15T10:00:00Z',
      updatedAt: '2024-01-15T10:00:00Z',
    },
    {
      id: 2,
      title: '新功能上线',
      content: '我们推出了全新的陪玩匹配功能，快来体验吧！',
      type: 'success',
      enabled: true,
      createdAt: '2024-01-10T10:00:00Z',
      updatedAt: '2024-01-10T10:00:00Z',
    },
  ]);

  const [showModal, setShowModal] = useState(false);
  const [editingAnnouncement, setEditingAnnouncement] = useState<Announcement | null>(null);
  const [deletingAnnouncement, setDeletingAnnouncement] = useState<Announcement | null>(null);

  // 创建公告
  const handleCreate = () => {
    setEditingAnnouncement(null);
    setShowModal(true);
  };

  // 编辑公告
  const handleEdit = (announcement: Announcement) => {
    setEditingAnnouncement(announcement);
    setShowModal(true);
  };

  // 删除公告
  const handleDelete = (announcement: Announcement) => {
    setDeletingAnnouncement(announcement);
  };

  // 确认删除
  const handleConfirmDelete = () => {
    if (deletingAnnouncement) {
      setAnnouncements(announcements.filter((a) => a.id !== deletingAnnouncement.id));
      setDeletingAnnouncement(null);
    }
  };

  // 切换启用状态
  const handleToggleEnabled = (announcement: Announcement) => {
    setAnnouncements(
      announcements.map((a) =>
        a.id === announcement.id ? { ...a, enabled: !a.enabled } : a
      )
    );
  };

  // 获取公告类型颜色
  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      info: 'blue',
      warning: 'orange',
      error: 'red',
      success: 'green',
    };
    return colors[type] || 'default';
  };

  // 获取公告类型文本
  const getTypeText = (type: string) => {
    const texts: Record<string, string> = {
      info: '信息',
      warning: '警告',
      error: '错误',
      success: '成功',
    };
    return texts[type] || type;
  };

  return (
    <>
      <Card className={styles.settingsCard}>
        <div className={styles.announcementHeader}>
          <h3 className={styles.announcementTitle}>平台公告列表</h3>
          <Button type="primary" onClick={handleCreate}>
            创建公告
          </Button>
        </div>

        <div className={styles.announcementList}>
          {announcements.length === 0 ? (
            <div className={styles.emptyState}>暂无公告</div>
          ) : (
            announcements.map((announcement) => (
              <div key={announcement.id} className={styles.announcementItem}>
                <div className={styles.announcementContent}>
                  <div className={styles.announcementTop}>
                    <h4 className={styles.announcementItemTitle}>{announcement.title}</h4>
                    <div className={styles.announcementTags}>
                      <Tag color={getTypeColor(announcement.type)}>
                        {getTypeText(announcement.type)}
                      </Tag>
                      <Tag color={announcement.enabled ? 'green' : 'default'}>
                        {announcement.enabled ? '已启用' : '已禁用'}
                      </Tag>
                    </div>
                  </div>

                  <p className={styles.announcementText}>{announcement.content}</p>

                  <div className={styles.announcementMeta}>
                    <span className={styles.announcementDate}>
                      创建时间：{formatDateTime(announcement.createdAt)}
                    </span>
                    {announcement.startDate && announcement.endDate && (
                      <span className={styles.announcementDate}>
                        有效期：{announcement.startDate} 至 {announcement.endDate}
                      </span>
                    )}
                  </div>
                </div>

                <div className={styles.announcementActions}>
                  <ActionButtons
                    actions={[
                      {
                        label: announcement.enabled ? '禁用' : '启用',
                        onClick: () => handleToggleEnabled(announcement),
                      },
                      {
                        label: '编辑',
                        onClick: () => handleEdit(announcement),
                      },
                      {
                        label: '删除',
                        onClick: () => handleDelete(announcement),
                        type: 'danger',
                      },
                    ]}
                  />
                </div>
              </div>
            ))
          )}
        </div>
      </Card>

      {/* 公告表单模态框 */}
      {showModal && (
        <AnnouncementFormModal
          visible={showModal}
          announcement={editingAnnouncement}
          onClose={() => {
            setShowModal(false);
            setEditingAnnouncement(null);
          }}
          onSuccess={(data) => {
            if (editingAnnouncement) {
              // 更新
              setAnnouncements(
                announcements.map((a) =>
                  a.id === editingAnnouncement.id
                    ? { ...editingAnnouncement, ...data, updatedAt: new Date().toISOString() }
                    : a
                )
              );
            } else {
              // 创建
              const newAnnouncement: Announcement = {
                id: Math.max(...announcements.map((a) => a.id), 0) + 1,
                ...data,
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
              };
              setAnnouncements([newAnnouncement, ...announcements]);
            }
            setShowModal(false);
            setEditingAnnouncement(null);
          }}
        />
      )}

      {/* 删除确认模态框 */}
      {deletingAnnouncement && (
        <DeleteConfirmModal
          visible={!!deletingAnnouncement}
          title="删除公告"
          content={
            <>
              确定要删除公告 <strong>{deletingAnnouncement.title}</strong> 吗？
            </>
          }
          onConfirm={handleConfirmDelete}
          onCancel={() => setDeletingAnnouncement(null)}
        />
      )}
    </>
  );
};

