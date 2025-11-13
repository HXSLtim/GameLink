import React, { useState, useEffect } from 'react';
import { Modal, Button, Checkbox } from '../../components';
import { roleApi, permissionApi } from '../../services/api/rbac';
import type { Role, Permission, PermissionGroup, HTTPMethod } from '../../types/rbac';
import styles from './PermissionAssignModal.module.less';

interface PermissionAssignModalProps {
  visible: boolean;
  role: Role;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 权限分配模态框
 * 显示权限树形结构，按分组展示
 */
export const PermissionAssignModal: React.FC<PermissionAssignModalProps> = ({
  visible,
  role,
  onClose,
  onSuccess,
}) => {
  const [permissionGroups, setPermissionGroups] = useState<PermissionGroup[]>([]);
  const [selectedPermissionIds, setSelectedPermissionIds] = useState<Set<number>>(new Set());
  const [isLoading, setIsLoading] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 加载权限分组和角色已有权限
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // 获取所有权限列表（循环分页获取所有数据）
        // 先获取第一页以了解总数
        const firstPage = await permissionApi.getList({
          page: 1,
          pageSize: 100, // 尝试每页100条
        });
        
        const total = firstPage.total || 0;
        const pageSize = (firstPage.list || []).length; // 实际返回的每页数量
        
        console.log(`📊 权限总数: ${total}, 每页实际返回: ${pageSize} 条`);
        
        // 收集所有权限
        let allPermissions: Permission[] = [...(firstPage.list || [])];
        
        // 如果还有更多页，继续获取
        if (pageSize > 0 && total > allPermissions.length) {
          const totalPages = Math.ceil(total / pageSize);
          console.log(`📄 需要获取 ${totalPages} 页数据`);
          
          // 从第2页开始获取剩余数据
          const remainingPages = [];
          for (let page = 2; page <= totalPages; page++) {
            remainingPages.push(
              permissionApi.getList({
                page,
                pageSize: 100,
              })
            );
          }
          
          // 并行获取所有剩余页面
          const results = await Promise.all(remainingPages);
          results.forEach((result) => {
            allPermissions = [...allPermissions, ...(result.list || [])];
          });
        }
        
        console.log(`✅ 成功加载了 ${allPermissions.length} 个权限（总数：${total}）`);
        
        // 按分组聚合权限
        const groupMap = new Map<string, Permission[]>();
        allPermissions.forEach((perm) => {
          const group = perm.group || '其他';
          if (!groupMap.has(group)) {
            groupMap.set(group, []);
          }
          groupMap.get(group)!.push(perm);
        });
        
        // 转换为 PermissionGroup 数组
        const groups: PermissionGroup[] = Array.from(groupMap.entries())
          .map(([groupName, permissions]) => ({
            group: groupName,
            permissions: permissions.sort((a, b) => {
              // 按 HTTP 方法和路径排序
              if (a.method !== b.method) {
                const methodOrder: Record<HTTPMethod, number> = { 
                  GET: 1, 
                  POST: 2, 
                  PUT: 3, 
                  PATCH: 4, 
                  DELETE: 5 
                };
                return (methodOrder[a.method] || 99) - (methodOrder[b.method] || 99);
              }
              return a.path.localeCompare(b.path);
            }),
          }))
          .sort((a, b) => a.group.localeCompare(b.group));
        
        console.log(`📦 生成了 ${groups.length} 个权限分组`);
        setPermissionGroups(groups);

        // 尝试获取角色已有权限（如果失败则跳过）
        try {
          const rolePermissions = await roleApi.getPermissions(role.id);
          setSelectedPermissionIds(new Set(rolePermissions.map((p) => p.id)));
        } catch (permErr) {
          console.warn('无法加载角色现有权限，将从空白开始:', permErr);
          // 如果获取角色权限失败，从空白开始（不报错）
          setSelectedPermissionIds(new Set());
        }
      } catch (err) {
        console.error('加载权限数据失败:', err);
        setError('加载权限数据失败，请刷新重试');
      } finally {
        setIsLoading(false);
      }
    };

    if (visible) {
      loadData();
    }
  }, [visible, role.id]);

  // 切换权限选择
  const togglePermission = (permissionId: number) => {
    const newSet = new Set(selectedPermissionIds);
    if (newSet.has(permissionId)) {
      newSet.delete(permissionId);
    } else {
      newSet.add(permissionId);
    }
    setSelectedPermissionIds(newSet);
  };

  // 切换分组所有权限
  const toggleGroup = (permissions: Permission[]) => {
    const groupIds = permissions.map((p) => p.id);
    const allSelected = groupIds.every((id) => selectedPermissionIds.has(id));
    
    const newSet = new Set(selectedPermissionIds);
    if (allSelected) {
      groupIds.forEach((id) => newSet.delete(id));
    } else {
      groupIds.forEach((id) => newSet.add(id));
    }
    setSelectedPermissionIds(newSet);
  };

  // 提交保存
  const handleSubmit = async () => {
    try {
      setIsSubmitting(true);
      setError(null);

      await roleApi.assignPermissions(role.id, {
        permissionIds: Array.from(selectedPermissionIds),
      });

      onSuccess();
    } catch (err: any) {
      console.error('分配权限失败:', err);
      setError(err.message || '分配权限失败，请重试');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal
      visible={visible}
      title={`分配权限 - ${role.name}`}
      onCancel={onClose}
      width={900}
      footer={
        <div className={styles.actions}>
          <Button variant="secondary" onClick={onClose} disabled={isSubmitting}>
            取消
          </Button>
          <Button variant="primary" onClick={handleSubmit} loading={isSubmitting}>
            保存
          </Button>
        </div>
      }
    >
      <div className={styles.container}>
        {error && <div className={styles.error}>{error}</div>}

        {isLoading ? (
          <div className={styles.loading}>加载中...</div>
        ) : (
          <>
            <div className={styles.treeContainer}>
              {permissionGroups.map((group) => {
                const allSelected = group.permissions.every((p) =>
                  selectedPermissionIds.has(p.id)
                );
                const someSelected = group.permissions.some((p) =>
                  selectedPermissionIds.has(p.id)
                );

                return (
                  <div key={group.group} className={styles.permissionGroup}>
                    <div className={styles.groupHeader}>
                      <div className={styles.groupCheckbox}>
                        <Checkbox
                          checked={allSelected}
                          indeterminate={someSelected && !allSelected}
                          onChange={() => toggleGroup(group.permissions)}
                        />
                      </div>
                      <span className={styles.groupName}>{group.group}</span>
                      <span className={styles.groupCount}>
                        {group.permissions.length} 项
                      </span>
                    </div>

                    <div className={styles.permissionList}>
                      {group.permissions.map((permission) => (
                        <div key={permission.id} className={styles.permissionItem}>
                          <div className={styles.permissionCheckbox}>
                            <Checkbox
                              checked={selectedPermissionIds.has(permission.id)}
                              onChange={() => togglePermission(permission.id)}
                            />
                          </div>
                          <div className={styles.permissionDetails}>
                            <div className={styles.permissionMain}>
                              <span
                                className={`${styles.methodBadge} ${styles[permission.method.toLowerCase()]}`}
                              >
                                {permission.method}
                              </span>
                              <code className={styles.permissionPath}>{permission.path}</code>
                            </div>
                            <div className={styles.permissionCode}>{permission.code}</div>
                            {permission.description && (
                              <div className={styles.permissionDescription}>
                                {permission.description}
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>

            <div className={styles.statsBar}>
              <span className={styles.statsText}>
                已选择 <strong>{selectedPermissionIds.size}</strong> 个权限
              </span>
            </div>
          </>
        )}
      </div>
    </Modal>
  );
};

