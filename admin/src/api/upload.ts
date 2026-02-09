/**
 * 文件上传 API
 * 对接后端 /admin/upload/image 接口
 */
import apiClient from './client';
import type { ApiResponse } from './admin';

// ============================================================================
// Upload Types (上传类型)
// ============================================================================

/**
 * 上传响应
 */
export interface UploadResponse {
    filePath: string;
    hash: string;
}

/**
 * 上传进度回调
 */
export type UploadProgressCallback = (percent: number) => void;

/**
 * 上传选项
 */
export interface UploadOptions {
    onProgress?: UploadProgressCallback;
    maxSize?: number; // 最大文件大小（字节），默认 5MB
}

// ============================================================================
// Upload API
// ============================================================================

export const uploadApi = {
    /**
     * 上传图片
     * POST /api/v1/admin/upload/image
     *
     * @param file - 图片文件
     * @param options - 上传选项
     * @returns 上传结果，包含文件路径和哈希值
     *
     * 支持的图片格式：image/jpeg, image/png, image/gif, image/webp
     * 最大文件大小：5MB
     */
    uploadImage: async (file: File, options?: UploadOptions): Promise<ApiResponse<UploadResponse>> => {
        const maxSize = options?.maxSize || 5 * 1024 * 1024; // 默认 5MB

        // 前端预校验文件大小
        if (file.size > maxSize) {
            return {
                success: false,
                code: 400,
                message: `文件大小超过限制（最大 ${Math.round(maxSize / 1024 / 1024)}MB）`,
                data: { filePath: '', hash: '' },
            };
        }

        // 前端预校验文件类型
        const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
        if (!allowedTypes.includes(file.type)) {
            return {
                success: false,
                code: 400,
                message: '不支持的文件类型，仅支持 JPG、PNG、GIF、WebP 格式',
                data: { filePath: '', hash: '' },
            };
        }

        const formData = new FormData();
        formData.append('file', file);

        return apiClient.post<ApiResponse<UploadResponse>>('/admin/upload/image', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
            onUploadProgress: (progressEvent) => {
                if (options?.onProgress && progressEvent.total) {
                    const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total);
                    options.onProgress(percent);
                }
            },
        }).then(res => res.data);
    },

    /**
     * 批量上传图片
     * 依次调用单图上传接口
     *
     * @param files - 图片文件数组
     * @param options - 上传选项
     * @returns 上传结果数组
     */
    uploadImages: async (
        files: File[],
        options?: UploadOptions
    ): Promise<{ success: boolean; results: Array<{ file: File; response: ApiResponse<UploadResponse> }> }> => {
        const results: Array<{ file: File; response: ApiResponse<UploadResponse> }> = [];
        let successCount = 0;

        for (const file of files) {
            try {
                const response = await uploadApi.uploadImage(file, options);
                results.push({ file, response });
                if (response.success) {
                    successCount++;
                }
            } catch (error) {
                results.push({
                    file,
                    response: {
                        success: false,
                        code: 500,
                        message: error instanceof Error ? error.message : '上传失败',
                        data: { filePath: '', hash: '' },
                    },
                });
            }
        }

        return {
            success: successCount === files.length,
            results,
        };
    },
};

// ============================================================================
// Helper Functions (辅助函数)
// ============================================================================

/**
 * 从上传响应获取完整的图片 URL
 * @param filePath - 上传返回的文件路径
 * @param baseUrl - 基础 URL，默认为空（使用相对路径）
 */
export const getImageUrl = (filePath: string, baseUrl = ''): string => {
    if (!filePath) return '';
    if (filePath.startsWith('http://') || filePath.startsWith('https://')) {
        return filePath;
    }
    return `${baseUrl}${filePath}`;
};

/**
 * 验证文件是否为图片
 * @param file - 文件对象
 */
export const isImageFile = (file: File): boolean => {
    const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    return allowedTypes.includes(file.type);
};

/**
 * 格式化文件大小显示
 * @param bytes - 字节数
 */
export const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export default uploadApi;
