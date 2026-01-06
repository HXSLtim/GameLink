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
      const value = row[col.key];
      newRow[col.key as string] = col.format ? col.format(value, row) : value;
    });
    return newRow as T;
  });

  const simpleColumns = columns.map((col) => ({ key: col.key, title: col.title }));

  if (format === 'excel') {
    return exportToExcel(formattedData, simpleColumns, filename);
  }
  return exportToCSV(formattedData, simpleColumns, filename);
}
