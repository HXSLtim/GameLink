/**
 * 数据导出工具
 * 支持 CSV 和 Excel 格式导出
 */

/**
 * 将数据导出为 CSV 文件
 *
 * @param data 数据数组
 * @param columns 列配置 { key, title }
 * @param filename 文件名（不含扩展名）
 */
export function exportToCSV<T extends Record<string, unknown>>(
  data: T[],
  columns: { key: keyof T; title: string }[],
  filename: string
): void {
  if (data.length === 0) {
    console.warn('No data to export');
    return;
  }

  // 构建 CSV 内容
  const headers = columns.map((col) => `"${col.title}"`).join(',');
  const rows = data.map((row) =>
    columns
      .map((col) => {
        const value = row[col.key];
        // 处理特殊字符
        if (value === null || value === undefined) return '""';
        const str = String(value).replace(/"/g, '""');
        return `"${str}"`;
      })
      .join(',')
  );

  const csvContent = '\uFEFF' + [headers, ...rows].join('\n'); // BOM for Excel UTF-8

  // 下载文件
  downloadFile(csvContent, `${filename}.csv`, 'text/csv;charset=utf-8');
}

/**
 * 将数据导出为 Excel 文件
 * 使用 xlsx 库（需要已安装）
 *
 * @param data 数据数组
 * @param columns 列配置 { key, title }
 * @param filename 文件名（不含扩展名）
 * @param sheetName 工作表名称
 */
export async function exportToExcel<T extends Record<string, unknown>>(
  data: T[],
  columns: { key: keyof T; title: string }[],
  filename: string,
  sheetName = 'Sheet1'
): Promise<void> {
  if (data.length === 0) {
    console.warn('No data to export');
    return;
  }

  try {
    // 动态导入 xlsx 库
    const XLSX = await import('xlsx');

    // 转换数据格式
    const headers = columns.map((col) => col.title);
    const rows = data.map((row) => columns.map((col) => row[col.key] ?? ''));

    // 创建工作簿
    const wb = XLSX.utils.book_new();
    const ws = XLSX.utils.aoa_to_sheet([headers, ...rows]);

    // 设置列宽
    ws['!cols'] = columns.map(() => ({ wch: 15 }));

    XLSX.utils.book_append_sheet(wb, ws, sheetName);

    // 导出文件
    XLSX.writeFile(wb, `${filename}.xlsx`);
  } catch (error) {
    console.error('Failed to export Excel:', error);
    // 降级到 CSV
    exportToCSV(data, columns, filename);
  }
}

/**
 * 下载文件
 */
function downloadFile(content: string | Blob, filename: string, mimeType: string): void {
  const blob = content instanceof Blob ? content : new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);

  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  link.style.display = 'none';

  document.body.appendChild(link);
  link.click();

  // 清理
  setTimeout(() => {
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }, 100);
}

/**
 * 格式化日期为导出友好格式
 */
export function formatDateForExport(date: Date | string | number): string {
  const d = new Date(date);
  if (isNaN(d.getTime())) return '';
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

/**
 * 格式化金额为导出友好格式（分转元）
 */
export function formatAmountForExport(cents: number): string {
  return (cents / 100).toFixed(2);
}

/**
 * 导出配置类型
 */
export interface ExportColumn<T> {
  key: keyof T;
  title: string;
  /** 自定义格式化函数 */
  format?: (value: T[keyof T], row: T) => string | number;
}

/**
 * 获取嵌套对象的值
 */
function getNestedValue<T extends Record<string, unknown>>(obj: T, key: string): unknown {
  const keys = String(key).split('.');
  let value: unknown = obj;
  for (const k of keys) {
    if (value === null || value === undefined) return undefined;
    value = (value as Record<string, unknown>)[k];
  }
  return value;
}

/**
 * 带格式化的导出
 */
export function exportWithFormat<T extends Record<string, unknown>>(
  data: T[],
  columns: ExportColumn<T>[],
  filename: string,
  format: 'csv' | 'excel' = 'csv'
): void | Promise<void> {
  // 应用格式化
  const formattedData = data.map((row) => {
    const newRow: Record<string, unknown> = {};
    columns.forEach((col) => {
      const value = getNestedValue(row, String(col.key));
      newRow[col.key as string] = col.format ? col.format(value as T[keyof T], row) : value;
    });
    return newRow as T;
  });

  const simpleColumns = columns.map((col) => ({ key: col.key, title: col.title }));

  if (format === 'excel') {
    return exportToExcel(formattedData, simpleColumns, filename);
  }
  return exportToCSV(formattedData, simpleColumns, filename);
}

/**
 * 预定义导出列配置
 */

// 角色映射
const roleMap: Record<string, string> = {
  user: '普通用户',
  player: '陪玩师',
  admin: '管理员',
};

// 状态映射
const statusMap: Record<string, string> = {
  active: '正常',
  inactive: '未激活',
  banned: '已封禁',
  pending: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
};

// 订单状态映射
const orderStatusMap: Record<string, string> = {
  pending: '待确认',
  confirmed: '已确认',
  in_progress: '进行中',
  completed: '已完成',
  canceled: '已取消',
  refunded: '已退款',
};

/**
 * 用户导出列配置（完整版）
 */
export const userExportColumns: ExportColumn<Record<string, unknown>>[] = [
  { key: 'id', title: 'ID' },
  { key: 'name', title: '用户名' },
  { key: 'phone', title: '手机号' },
  { key: 'email', title: '邮箱' },
  {
    key: 'role',
    title: '角色',
    format: (value) => roleMap[String(value)] || String(value),
  },
  {
    key: 'status',
    title: '状态',
    format: (value) => statusMap[String(value)] || String(value),
  },
  { key: 'level', title: '等级' },
  {
    key: 'wallet',
    title: '积分/余额',
    format: (value, row) => {
      // 支持嵌套对象 wallet.balanceCents
      const wallet = row.wallet as { balanceCents?: number } | undefined;
      return wallet?.balanceCents ?? 0;
    },
  },
  {
    key: 'tags',
    title: '标签',
    format: (value) => {
      const tags = value as string[] | undefined;
      return tags?.join(', ') || '';
    },
  },
  {
    key: 'vipLevelId',
    title: 'VIP等级ID',
    format: (value) => (value ? String(value) : '无'),
  },
  {
    key: 'vipExp',
    title: 'VIP经验',
    format: (value) => (value ? String(value) : '0'),
  },
  {
    key: 'totalRechargeCents',
    title: '累计充值(元)',
    format: (value) => formatAmountForExport(Number(value) || 0),
  },
  {
    key: 'vipExpiry',
    title: 'VIP到期时间',
    format: (value) => (value ? formatDateForExport(value as string) : '永久'),
  },
  {
    key: 'lastLoginAt',
    title: '最后登录',
    format: (value) => (value ? formatDateForExport(value as string) : ''),
  },
  {
    key: 'createdAt',
    title: '注册时间',
    format: (value) => (value ? formatDateForExport(value as string) : ''),
  },
];

/**
 * 订单导出列配置
 */
export const orderExportColumns: ExportColumn<Record<string, unknown>>[] = [
  { key: 'id', title: 'ID' },
  { key: 'orderNo', title: '订单号' },
  { key: 'userId', title: '用户ID' },
  { key: 'playerId', title: '陪玩ID' },
  {
    key: 'totalPriceCents',
    title: '金额(元)',
    format: (value) => formatAmountForExport(Number(value) || 0),
  },
  {
    key: 'status',
    title: '状态',
    format: (value) => orderStatusMap[String(value)] || String(value),
  },
  {
    key: 'createdAt',
    title: '创建时间',
    format: (value) => (value ? formatDateForExport(value as string) : ''),
  },
];

/**
 * 陪玩师导出列配置
 */
export const playerExportColumns: ExportColumn<Record<string, unknown>>[] = [
  { key: 'id', title: 'ID' },
  { key: 'userId', title: '用户ID' },
  { key: 'nickname', title: '昵称' },
  { key: 'realName', title: '真实姓名' },
  {
    key: 'status',
    title: '状态',
    format: (value) => statusMap[String(value)] || String(value),
  },
  { key: 'ratingAverage', title: '平均评分' },
  { key: 'ratingCount', title: '评价数' },
  { key: 'orderCount', title: '订单数' },
  {
    key: 'createdAt',
    title: '注册时间',
    format: (value) => (value ? formatDateForExport(value as string) : ''),
  },
];
