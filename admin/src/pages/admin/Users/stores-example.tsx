import { useEffect } from 'react';
import { Table, Button, Space, Tag, message } from 'antd';
import { useUserStore } from '@/stores/modules/userStore';
import { useAuthStore } from '@/stores/modules/authStore';

/**
 * 示例：使用 Zustand stores 替换 Context API
 *
 * 迁移前: 每个 page 独立管理状态
 * 迁移后: 使用全局 stores，数据自动共享
 */
const UsersStoresExample = () => {
  // 使用 userStore
  const {
    users,
    loading,
    error,
    pagination,
    fetchUsers,
    deleteUser,
    setFilters,
    reset: resetUserStore,
  } = useUserStore();

  // 使用 authStore
  const { hasPermission, isAdmin } = useAuthStore();

  // 初始化加载数据
  useEffect(() => {
    fetchUsers();
  }, []);

  // 权限检查
  if (!hasPermission('admin.users.read')) {
    return <div>无权限访问</div>;
  }

  // 删除用户
  const handleDelete = async (userId: number) => {
    try {
      await deleteUser(userId);
      message.success('删除成功');
      await fetchUsers(); // 刷新列表
    } catch (error) {
      message.error('删除失败');
    }
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '邮箱', dataIndex: 'email', key: 'email' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '启用' : '禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link">编辑</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <h1>用户管理 (Zustand Stores 示例)</h1>

      {/* 权限提示 */}
      {isAdmin() && <Tag color="blue">管理员</Tag>}

      {/* 错误提示 */}
      {error && <div style={{ color: 'red' }}>{error}</div>}

      {/* 数据表格 */}
      <Table
        dataSource={users}
        columns={columns}
        loading={loading}
        rowKey="id"
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          onChange: (page) => fetchUsers(page),
        }}
      />

      {/* 清理按钮 */}
      <Button onClick={resetUserStore} style={{ marginTop: 16 }}>
        重置状态
      </Button>
    </div>
  );
};

export default UsersStoresExample;
