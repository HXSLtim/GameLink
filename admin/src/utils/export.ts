/**
 * 通用导出工具
 * 支持 CSV 和 Excel 格式导出
 */
import dayjs from 'dayjs';

import { logger } from '@/utils/logger';
export interface ExportColumn {
    key: string;
    title: string;
    render?: (value: unknown, record: Record<string, unknown>) => string;
}

/**
 * 导出为 CSV 文件
 */
export function exportToCSV<T extends Record<string, unknown>>(
    data: T[],
    columns: ExportColumn[],
    filename: string
): void {
    if (!data || data.length === 0) {
        logger.warn('No data to export');
        return;
    }

    // 生成表头
    const headers = columns.map(col => `"${col.title}"`).join(',');

    // 生成数据行
    const rows = data.map(record => {
        return columns.map(col => {
            let value: unknown;
            if (col.render) {
                value = col.render(record[col.key], record);
            } else {
                // 支持嵌套属性，如 'user.name'
                value = col.key.split('.').reduce((obj, key) => {
                    return obj && typeof obj === 'object' ? (obj as Record<string, unknown>)[key] : undefined;
                }, record as unknown);
            }
            // 处理特殊字符
            const strValue = value === null || value === undefined ? '' : String(value);
            return `"${strValue.replace(/"/g, '""')}"`;
        }).join(',');
    });

    // 组合 CSV 内容（添加 BOM 以支持中文）
    const csvContent = '\uFEFF' + [headers, ...rows].join('\n');

    // 创建下载链接
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${filename}_${dayjs().format('YYYYMMDD_HHmmss')}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
}

/**
 * 导出为 Excel 文件（使用 xlsx 库，如果可用）
 * 如果 xlsx 库不可用，回退到 CSV 格式
 */
export async function exportToExcel<T extends Record<string, unknown>>(
    data: T[],
    columns: ExportColumn[],
    filename: string,
    sheetName = 'Sheet1'
): Promise<void> {
    try {
        // 动态导入 xlsx 库
        const XLSX = await import('xlsx');
        
        // 准备数据
        const headers = columns.map(col => col.title);
        const rows = data.map(record => {
            return columns.map(col => {
                if (col.render) {
                    return col.render(record[col.key], record);
                }
                // 支持嵌套属性
                return col.key.split('.').reduce((obj, key) => {
                    return obj && typeof obj === 'object' ? (obj as Record<string, unknown>)[key] : undefined;
                }, record as unknown) ?? '';
            });
        });

        // 创建工作表
        const ws = XLSX.utils.aoa_to_sheet([headers, ...rows]);
        
        // 设置列宽
        ws['!cols'] = columns.map(() => ({ wch: 20 }));

        // 创建工作簿
        const wb = XLSX.utils.book_new();
        XLSX.utils.book_append_sheet(wb, ws, sheetName);

        // 导出文件
        XLSX.writeFile(wb, `${filename}_${dayjs().format('YYYYMMDD_HHmmss')}.xlsx`);
    } catch {
        // xlsx 库不可用，回退到 CSV
        logger.warn('xlsx library not available, falling back to CSV export');
        exportToCSV(data, columns, filename);
    }
}

/**
 * 用户导出列配置
 */
export const userExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '用户名' },
    { key: 'email', title: '邮箱' },
    { key: 'phone', title: '手机号' },
    { key: 'role', title: '角色', render: (v) => {
        const map: Record<string, string> = { user: '普通用户', player: '陪玩师', admin: '管理员' };
        return map[v as string] || String(v);
    }},
    { key: 'status', title: '状态', render: (v) => {
        const map: Record<string, string> = { active: '正常', banned: '已封禁', suspended: '已停用' };
        return map[v as string] || String(v);
    }},
    { key: 'level', title: '等级' },
    { key: 'createdAt', title: '注册时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
    { key: 'lastLoginAt', title: '最后登录', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

/**
 * 订单导出列配置
 */
export const orderExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'orderNo', title: '订单号' },
    { key: 'user.name', title: '用户' },
    { key: 'player.nickname', title: '陪玩师' },
    { key: 'game.name', title: '游戏' },
    { key: 'title', title: '标题' },
    { key: 'totalPriceCents', title: '金额(元)', render: (v) => ((v as number) / 100).toFixed(2) },
    { key: 'status', title: '状态', render: (v) => {
        const map: Record<string, string> = {
            pending: '待确认', confirmed: '已确认', in_progress: '进行中',
            completed: '已完成', canceled: '已取消', refunded: '已退款'
        };
        return map[v as string] || String(v);
    }},
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
    { key: 'completedAt', title: '完成时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

/**
 * 陪玩师导出列配置
 */
export const playerExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'nickname', title: '昵称' },
    { key: 'user.name', title: '用户名' },
    { key: 'user.email', title: '邮箱' },
    { key: 'user.phone', title: '手机号' },
    { key: 'level', title: '等级' },
    { key: 'rating', title: '评分' },
    { key: 'orderCount', title: '订单数' },
    { key: 'status', title: '状态', render: (v) => {
        const map: Record<string, string> = { active: '正常', pending: '待审核', rejected: '已拒绝', suspended: '已停用' };
        return map[v as string] || String(v);
    }},
    { key: 'createdAt', title: '注册时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

/**
 * 游戏导出列配置
 */
export const gameExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: '游戏名称' },
    { key: 'category', title: '分类' },
    { key: 'playerCount', title: '陪玩师数' },
    { key: 'orderCount', title: '订单数' },
    { key: 'isActive', title: '状态', render: (v) => v ? '启用' : '禁用' },
    { key: 'sortOrder', title: '排序' },
    { key: 'createdAt', title: '创建时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];

/**
 * 提现导出列配置
 */
export const withdrawExportColumns: ExportColumn[] = [
    { key: 'id', title: 'ID' },
    { key: 'withdrawNo', title: '提现单号' },
    { key: 'user.name', title: '用户' },
    { key: 'amountCents', title: '金额(元)', render: (v) => ((v as number) / 100).toFixed(2) },
    { key: 'feeCents', title: '手续费(元)', render: (v) => ((v as number) / 100).toFixed(2) },
    { key: 'actualAmountCents', title: '实际到账(元)', render: (v) => ((v as number) / 100).toFixed(2) },
    { key: 'bankName', title: '银行' },
    { key: 'bankAccount', title: '账号' },
    { key: 'status', title: '状态', render: (v) => {
        const map: Record<string, string> = { pending: '待审核', approved: '已通过', rejected: '已拒绝', completed: '已完成' };
        return map[v as string] || String(v);
    }},
    { key: 'createdAt', title: '申请时间', render: (v) => v ? dayjs(v as string).format('YYYY-MM-DD HH:mm:ss') : '' },
];
