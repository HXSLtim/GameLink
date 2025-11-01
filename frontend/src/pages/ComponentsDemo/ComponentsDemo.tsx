import React, { useMemo, useState } from 'react';
import {
  Button,
  Row,
  Col,
  Space,
  Tabs,
  Tree,
  Select,
  notify,
  message,
  Breadcrumb,
  Card,
  Input,
  PasswordInput,
  Pagination,
  Tag,
  Badge,
  Modal,
  DeleteConfirmModal,
  Rating,
  SimpleRating,
  Table,
  CardSkeleton,
  TableSkeleton,
  ActionButtons,
  BulkActions,
  Form,
  FormItem,
  FormField,
  DataTable,
  StatCardSkeleton,
  ListItemSkeleton,
  PageSkeleton,
  FEATURE_ICONS,
  PAYMENT_METHOD_ICONS,
  ReviewModal,
  Header,
  Sidebar,
} from 'components';
import styles from './ComponentsDemo.module.less';
import { FEATURE_FLAGS } from '../../config';

const tabs = [
  { key: 'tab-1', title: 'Tab 1', disabled: false, content: '内容一：基本信息' },
  { key: 'tab-2', title: 'Tab 2', disabled: false, content: '内容二：扩展信息' },
  { key: 'tab-3', title: 'Tab 3', disabled: true, content: '内容三：禁用' },
];

const treeData = [
  {
    key: 'root',
    title: '根节点',
    children: [
      { key: 'leaf-1', title: '叶子 1' },
      {
        key: 'branch-1',
        title: '分支 1',
        children: [
          { key: 'leaf-2', title: '叶子 2' },
          { key: 'leaf-3', title: '叶子 3', disabled: true },
        ],
      },
    ],
  },
];

type RowType = { id: number; name: string; status: 'active' | 'inactive'; createdAt: string };

const DataTableDemoSection: React.FC = () => {
  const dtMock = useMemo<RowType[]>(
    () =>
      Array.from({ length: 24 }).map((_, i) => ({
        id: i + 1,
        name: `用户-${i + 1}`,
        status: i % 3 === 0 ? 'inactive' : 'active',
        createdAt: new Date(2024, 0, 1 + i).toISOString().slice(0, 10),
      })),
    [],
  );
  const [dtQuery, setDtQuery] = useState('');
  const [dtStatus, setDtStatus] = useState<string | undefined>(undefined);
  const [dtPage, setDtPage] = useState(1);
  const dtPageSize = 5;
  const dtFiltered = useMemo(
    () =>
      dtMock.filter(
        (r) => (!dtQuery || r.name.includes(dtQuery)) && (!dtStatus || r.status === dtStatus),
      ),
    [dtMock, dtQuery, dtStatus],
  );
  const dtTotal = dtFiltered.length;
  const dtPaged = useMemo(
    () => dtFiltered.slice((dtPage - 1) * dtPageSize, dtPage * dtPageSize),
    [dtFiltered, dtPage],
  );
  const dtColumns = [
    { key: 'name', title: '姓名', dataIndex: 'name' as const },
    {
      key: 'status',
      title: '状态',
      render: (_: any, r: RowType) => (
        <Tag color={r.status === 'active' ? 'success' : 'default'}>{r.status}</Tag>
      ),
    },
    { key: 'created', title: '创建日期', dataIndex: 'createdAt' as const },
    {
      key: 'op',
      title: '操作',
      render: (_: any, r: RowType) => (
        <ActionButtons
          onView={() => message({ type: 'info', content: `查看 ${r.name}` })}
          onEdit={() => message({ type: 'success', content: `编辑 ${r.name}` })}
          onDelete={() => message({ type: 'warning', content: `删除 ${r.name}` })}
        />
      ),
    },
  ];

  return (
    <section className={styles.section}>
      <h2 className={styles.sectionTitle}>DataTable</h2>
      <DataTable<RowType>
        title="用户列表"
        headerActions={
          <Button
            onClick={() =>
              notify({ type: 'success', title: '新增', description: '创建成功' })
            }
          >
            新增
          </Button>
        }
        filters={[
          {
            label: '关键词',
            key: 'query',
            element: (
              <Input
                placeholder="按姓名搜索"
                value={dtQuery}
                onChange={(e) => setDtQuery(e.target.value)}
                allowClear
              />
            ),
          },
          {
            label: '状态',
            key: 'status',
            element: (
              <Select
                value={dtStatus}
                options={[
                  { label: '全部', value: undefined as any },
                  { label: '启用', value: 'active' },
                  { label: '停用', value: 'inactive' },
                ]}
                onChange={(val) =>
                  setDtStatus(Array.isArray(val) ? undefined : (val as string))
                }
                placeholder="选择状态"
              />
            ),
          },
        ]}
        filterActions={(
          <Space>
            <Button
              variant="outlined"
              onClick={() => {
                setDtQuery('');
                setDtStatus(undefined);
                setDtPage(1);
              }}
            >
              重置
            </Button>
            <Button
              onClick={() =>
                notify({ type: 'info', title: '筛选', description: `结果：${dtTotal} 条` })
              }
            >
              筛选
            </Button>
          </Space>
        )}
        columns={dtColumns}
        dataSource={dtPaged}
        rowKey="id"
        pagination={{
          current: dtPage,
          pageSize: dtPageSize,
          total: dtTotal,
          onChange: (p) => setDtPage(p),
        }}
      />
    </section>
  );
};

export const ComponentsDemo: React.FC = () => {
  const gutter = useMemo<[number, number]>(() => [16, 16], []);
  const [singleValue, setSingleValue] = useState<string | number | undefined>(undefined);
  const [multiValue, setMultiValue] = useState<Array<string | number>>([]);
  const [inputVal, setInputVal] = useState('');
  const [passwordVal, setPasswordVal] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deleteLoading, setDeleteLoading] = useState(false);
  const [ratingValue] = useState(4.2);
  const [selectedCount, setSelectedCount] = useState(0);
  const [totalCount] = useState(20);

  // Form demo state
  const [formUsername, setFormUsername] = useState('');
  const [formPassword, setFormPassword] = useState('');
  const [formGender, setFormGender] = useState<string | undefined>(undefined);
  const [formError, setFormError] = useState<string | undefined>(undefined);

  // ReviewModal demo state
  const [showReview, setShowReview] = useState(false);

  // Header/Sidebar demo state
  const [demoSidebarCollapsed, setDemoSidebarCollapsed] = useState(false);

  // 示例展开/收起，受配置默认值控制
  const [examplesExpanded, setExamplesExpanded] = useState(
    FEATURE_FLAGS.showcase.expandExamplesByDefault,
  );

  // Table & Pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);
  const total = 23;
  const tableData = useMemo(() => {
    return Array.from({ length: total }).map((_, i) => ({
      id: i + 1,
      name: `用户 ${i + 1}`,
      age: 18 + ((i + 1) % 10),
    }));
  }, []);
  const pagedData = useMemo(() => {
    const start = (currentPage - 1) * pageSize;
    return tableData.slice(start, start + pageSize);
  }, [currentPage, pageSize, tableData]);

  const columns = useMemo(
    () => [
      { key: 'name', title: '姓名', dataIndex: 'name' as const },
      { key: 'age', title: '年龄', dataIndex: 'age' as const },
      {
        key: 'action',
        title: '操作',
        render: (_: any, record: { id: number; name: string }) => (
          <Button onClick={() => notify({ type: 'info', title: '查看', description: record.name })}>
            查看
          </Button>
        ),
      },
    ],
    [],
  );

  const baseOptions = useMemo(
    () => [
      { label: '选项 A', value: 'A' },
      { label: '选项 B', value: 'B' },
      { label: '选项 C', value: 'C' },
      { label: '选项 D', value: 'D' },
      { label: '选项 E', value: 'E' },
    ],
    [],
  );

  const mockAsyncSearch = async (q: string): Promise<{ label: string; value: string }[]> => {
    const lowered = q.toLowerCase();
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve(
          baseOptions.filter((opt) => String(opt.label).toLowerCase().includes(lowered)),
        );
      }, 300);
    });
  };

  const handleNotify = () => {
    notify({ type: 'success', title: '操作成功', description: '已保存更改。' });
  };

  const handleMessage = () => {
    message({ type: 'info', content: '这是一条轻提示' });
  };

  const handleDeleteConfirm = async () => {
    setDeleteLoading(true);
    await new Promise((r) => setTimeout(r, 800));
    setDeleteLoading(false);
    setShowDeleteModal(false);
    notify({ type: 'success', title: '删除成功', description: '记录已删除。' });
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.title}>组件演示</h1>

      <section className={styles.section}>
        <Space>
          <Button onClick={() => setExamplesExpanded((v) => !v)}>
            {examplesExpanded ? '收起组件使用示例' : '展开组件使用示例'}
          </Button>
          <Tag color={examplesExpanded ? 'success' : 'default'}>
            {examplesExpanded ? '已展开' : '已收起'}
          </Tag>
        </Space>
      </section>

      {examplesExpanded && (
        <>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Breadcrumb</h2>
        <Breadcrumb items={[{ label: '首页', path: '/' }, { label: '展示' }, { label: '组件' }]} />
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Grid & Space</h2>
        <Row gutter={gutter}>
          <Col span={6}>
            <div className={styles.box}>span=6</div>
          </Col>
          <Col span={6}>
            <div className={styles.box}>span=6</div>
          </Col>
          <Col span={6}>
            <div className={styles.box}>span=6</div>
          </Col>
          <Col span={6}>
            <div className={styles.box}>span=6</div>
          </Col>
        </Row>

        <Space size="medium">
          <Button>默认按钮</Button>
          <Button variant="outlined">Outlined</Button>
          <Button variant="text">Text</Button>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Input</h2>
        <Space>
          <Input
            placeholder="请输入文本"
            value={inputVal}
            onChange={(e) => setInputVal(e.target.value)}
            allowClear
          />
          <PasswordInput
            placeholder="请输入密码"
            value={passwordVal}
            onChange={(e) => setPasswordVal(e.target.value)}
          />
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Button（更多示例）</h2>
        <Space size="medium" wrap>
          <Button>Primary</Button>
          <Button variant="secondary">Secondary</Button>
          <Button variant="outlined">Outlined</Button>
          <Button variant="text">Text</Button>
          <Button size="small">Small</Button>
          <Button size="large">Large</Button>
          <Button loading>Loading</Button>
          <Button icon={<span>🔍</span>}>With Icon</Button>
          <Button block style={{ maxWidth: 240 }}>Block Button</Button>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Icons</h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))', gap: 12 }}>
          {Object.entries(FEATURE_ICONS).slice(0, 18).map(([key, Icon]) => (
            <div key={key} style={{ display: 'flex', alignItems: 'center', gap: 8, padding: 8, border: '1px solid var(--border-color)' }}>
              <Icon size={20} />
              <span style={{ color: 'var(--text-secondary)' }}>{key}</span>
            </div>
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Payment Icons</h2>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(160px, 1fr))', gap: 12 }}>
          {Object.entries(PAYMENT_METHOD_ICONS).map(([key, Icon]) => (
            <div key={key} style={{ display: 'flex', alignItems: 'center', gap: 10, padding: 10, border: '1px dashed var(--border-color)' }}>
              <Icon size={24} />
              <span style={{ fontWeight: 500 }}>{key}</span>
            </div>
          ))}
        </div>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Card / Badge / Tag</h2>
        <Card title="卡片标题" extra={<Tag color="processing">进行中</Tag>} hoverable>
          <Space>
            <Badge count={8}>
              <Button>带徽章按钮</Button>
            </Badge>
            <Badge dot />
            <Tag color="success">成功</Tag>
            <Tag color="warning" closable>
              可关闭标签
            </Tag>
          </Space>
        </Card>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Tabs</h2>
        <Tabs
          items={tabs.map((t) => ({
            key: t.key,
            label: t.title,
            disabled: t.disabled,
            children: <div className={styles.panel}>{t.content}</div>,
          }))}
          defaultActiveKey="tab-1"
        />
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Tree</h2>
        <Tree data={treeData} />
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Select（搜索 + 异步）</h2>
        <Space>
          <div style={{ minWidth: 240 }}>
            <Select
              options={baseOptions}
              value={singleValue}
              onChange={(val) => setSingleValue(Array.isArray(val) ? undefined : val)}
              placeholder="选择一个"
              searchable
            />
          </div>

          <div style={{ minWidth: 300 }}>
            <Select
              options={baseOptions}
              value={multiValue}
              onChange={(val) => setMultiValue(Array.isArray(val) ? val : [])}
              placeholder="选择多个"
              searchable
              multiple
              asyncSearch={mockAsyncSearch}
            />
          </div>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Form / FormItem / FormField</h2>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Card title="垂直布局">
            <Form layout="vertical" onSubmit={(e) => { e.preventDefault(); setFormError(undefined); if (!formUsername) setFormError('用户名不能为空'); }}>
              <FormItem label="用户名" required error={formError} help={!formError ? '请输入 3-16 位用户名' : undefined}>
                <Input value={formUsername} onChange={(e) => setFormUsername(e.target.value)} allowClear />
              </FormItem>
              <FormItem label="密码" required>
                <PasswordInput value={formPassword} onChange={(e) => setFormPassword(e.target.value)} />
              </FormItem>
              <FormItem label="性别">
                <Select
                  value={formGender}
                  options={[{ label: '男', value: 'male' }, { label: '女', value: 'female' }]}
                  onChange={(val) => setFormGender(Array.isArray(val) ? undefined : (val as string))}
                  placeholder="请选择"
                />
              </FormItem>
              <Space>
                <Button type="submit">提交</Button>
                <Button variant="outlined" onClick={() => { setFormUsername(''); setFormPassword(''); setFormGender(undefined); setFormError(undefined); }}>重置</Button>
              </Space>
            </Form>
          </Card>

          <Card title="水平布局 + FormField">
            <Form layout="horizontal">
              <FormField label="邮箱" required>
                <Input placeholder="example@domain.com" allowClear />
              </FormField>
              <FormField label="手机号">
                <Input placeholder="请输入手机号" />
              </FormField>
            </Form>
          </Card>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Table & Pagination</h2>
        <Card>
          <Table columns={columns} dataSource={pagedData} rowKey="id" />
          <div style={{ marginTop: 12 }}>
            <Pagination
              current={currentPage}
              total={total}
              pageSize={pageSize}
              onChange={(p) => setCurrentPage(p)}
              showSizeChanger
              onSizeChange={(s) => {
                setPageSize(s);
                setCurrentPage(1);
              }}
              pageSizeOptions={[5, 10, 20]}
            />
          </div>
        </Card>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>BulkActions / ActionButtons</h2>
        <Space direction="vertical" size="medium" style={{ width: '100%' }}>
          <BulkActions
            selectedCount={selectedCount}
            totalCount={totalCount}
            actions={[
              { key: 'export', label: '批量导出', variant: 'outlined', onClick: () => notify({ type: 'info', title: '批量导出', description: `导出 ${selectedCount} 项` }) },
              { key: 'delete', label: '批量删除', danger: true, onClick: () => notify({ type: 'warning', title: '批量删除', description: `准备删除 ${selectedCount} 项` }) },
            ]}
            onClearSelection={() => setSelectedCount(0)}
          />
          <Space>
            <Button onClick={() => setSelectedCount((c) => Math.min(totalCount, c + 1))}>选择 +1</Button>
            <Button variant="outlined" onClick={() => setSelectedCount((c) => Math.max(0, c - 1))}>选择 -1</Button>
            <Button variant="text" onClick={() => setSelectedCount(0)}>清空选择</Button>
          </Space>

          <Card title="操作按钮（单条）">
            <ActionButtons
              onView={() => notify({ type: 'info', title: '查看详情', description: '打开详情页...' })}
              onEdit={() => notify({ type: 'success', title: '编辑', description: '进入编辑模式' })}
              onDelete={() => notify({ type: 'warning', title: '删除', description: '请确认删除' })}
            />
          </Card>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Modal / DeleteConfirmModal</h2>
        <Space>
          <Button onClick={() => setShowModal(true)}>打开Modal</Button>
          <Button variant="secondary" onClick={() => setShowDeleteModal(true)}>
            删除确认
          </Button>
        </Space>
        <Modal visible={showModal} title="示例弹窗" onClose={() => setShowModal(false)}>
          <Space>
            <span>这里是弹窗内容</span>
            <Button onClick={() => setShowModal(false)}>关闭</Button>
          </Space>
        </Modal>
        <DeleteConfirmModal
          visible={showDeleteModal}
          content="确定要删除该记录吗？此操作不可撤销。"
          onConfirm={handleDeleteConfirm}
          onCancel={() => setShowDeleteModal(false)}
          loading={deleteLoading}
        />
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>ReviewModal 审核弹窗</h2>
        <Space>
          <Button onClick={() => setShowReview(true)}>打开审核弹窗</Button>
        </Space>
        <ReviewModal
          visible={showReview}
          orderNo="ORD-2024-0001"
          onClose={() => setShowReview(false)}
          onSubmit={async (data) => {
            await new Promise((r) => setTimeout(r, 800));
            notify({ type: 'success', title: data.approved ? '审核通过' : '审核拒绝', description: `订单 ${'ORD-2024-0001'}：${data.reason || '无原因'}` });
            setShowReview(false);
          }}
        />
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Rating</h2>
        <Space>
          <Rating value={ratingValue} showNumber showText />
          <SimpleRating value={ratingValue} />
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Skeleton</h2>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Space>
            <CardSkeleton hasImage lines={3} />
            <TableSkeleton rows={3} columns={4} />
          </Space>
          <Space>
            <StatCardSkeleton />
            <ListItemSkeleton />
          </Space>
          <PageSkeleton />
        </Space>
      </section>

      <DataTableDemoSection />

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Notification & Message</h2>
        <Space>
          <Button onClick={handleNotify}>显示通知</Button>
          <Button variant="secondary" onClick={handleMessage}>
            显示消息
          </Button>
        </Space>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Header & Sidebar</h2>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Card title="Header 示例">
            <div style={{ border: '1px solid var(--border-color)' }}>
              <Header
                user={{ username: 'Admin', role: '超级管理员' }}
                onLogout={() => message({ type: 'info', content: '已退出登录' })}
                onToggleSidebar={() => message({ type: 'info', content: '切换侧边栏' })}
                breadcrumbs={[{ label: '首页', path: '/' }, { label: '展示', path: '/showcase' }, { label: 'Header' }]}
              />
            </div>
          </Card>

          <Card
            title="Sidebar 示例"
            extra={<Button size="small" onClick={() => setDemoSidebarCollapsed((c) => !c)}>{demoSidebarCollapsed ? '展开' : '收起'}</Button>}
          >
            <div style={{ height: 260, border: '1px solid var(--border-color)' }}>
              <Sidebar
                menuItems={[
                  { key: 'dashboard', label: '仪表盘', path: '/dashboard', icon: React.createElement(FEATURE_ICONS.chart) },
                  { key: 'orders', label: '订单', path: '/orders', icon: React.createElement(FEATURE_ICONS.clipboard) },
                  { key: 'payments', label: '支付', path: '/payments', icon: React.createElement(FEATURE_ICONS.creditCard) },
                  { key: 'showcase', label: '组件演示', path: '/showcase', icon: React.createElement(FEATURE_ICONS.layers) },
                ]}
                collapsed={demoSidebarCollapsed}
              />
            </div>
          </Card>
        </Space>
      </section>
      </>
      )}
    </div>
  );
};

export default ComponentsDemo;