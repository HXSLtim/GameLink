import apiClient from './client';
import type { ApiResponse, Pagination } from '../types/api';

// ============================================================================
// Type Definitions
// ============================================================================

/**
 * Game Category - 游戏分类模型
 * 对应后端 model.GameCategory
 */
export interface GameCategory {
    id: number;
    name: string;
    description?: string;
    iconUrl?: string;
    sortOrder: number;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
    // 关联数据（非持久化字段）
    games?: Game[];
    serviceItems?: any[];
}

/**
 * Game Rank - 游戏段位配置模型
 * 对应后端 model.GameRank
 */
export interface GameRank {
    id: number;
    gameId: number;
    name: string;
    level: number;
    priceCents: number;
    iconUrl?: string;
    color?: string;
    description?: string;
    sortOrder: number;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
    // 关联数据（非持久化字段）
    game?: Game;
    gameName?: string;
}

/**
 * Game - 游戏基础信息（段位关联时返回）
 */
export interface Game {
    id: number;
    key: string;
    name: string;
    category?: string;
    categoryId?: number;
    iconUrl?: string;
    coverUrl?: string;
    description?: string;
    isActive: boolean;
    sortOrder: number;
}

// ============================================================================
// Request/Response DTOs - Game Category
// ============================================================================

/**
 * 创建游戏分类请求
 */
export interface CreateGameCategoryDto {
    name: string;
    description?: string;
    iconUrl?: string;
    sortOrder?: number;
}

/**
 * 更新游戏分类请求
 */
export interface UpdateGameCategoryDto {
    name?: string;
    description?: string;
    iconUrl?: string;
    sortOrder?: number;
    isActive?: boolean;
}

/**
 * 批量更新分类状态请求
 */
export interface BatchUpdateCategoryStatusDto {
    categoryIds: number[];
    isActive: boolean;
}

/**
 * 批量删除分类请求
 */
export interface BatchDeleteCategoriesDto {
    categoryIds: number[];
}

/**
 * 游戏分类查询参数
 */
export interface GameCategoryQueryParams {
    page?: number;
    pageSize?: number;
    keyword?: string;
    isActive?: boolean;
}

/**
 * 批量操作响应
 */
export interface BatchOperationResponse {
    successCount: number;
    failureCount: number;
    errors?: string[];
}

// ============================================================================
// Request/Response DTOs - Game Rank
// ============================================================================

/**
 * 创建游戏段位请求
 */
export interface CreateGameRankDto {
    gameId: number;
    name: string;
    level?: number;
    priceCents?: number;
    iconUrl?: string;
    color?: string;
    description?: string;
    sortOrder?: number;
    isActive?: boolean;
}

/**
 * 更新游戏段位请求
 */
export interface UpdateGameRankDto {
    name?: string;
    level?: number;
    priceCents?: number;
    iconUrl?: string;
    color?: string;
    description?: string;
    sortOrder?: number;
    isActive?: boolean;
}

/**
 * 批量删除段位请求
 */
export interface BatchDeleteGameRanksDto {
    ids: string[];
}

/**
 * 批量更新段位状态请求
 */
export interface BatchUpdateRankStatusDto {
    ids: string[];
    isActive: boolean;
}

/**
 * 游戏段位查询参数
 */
export interface GameRankQueryParams {
    page?: number;
    pageSize?: number;
    gameId?: number;
    keyword?: string;
    isActive?: boolean;
}

// ============================================================================
// API Client - Game Category
// ============================================================================

/**
 * 获取游戏分类列表
 * GET /admin/game-categories
 */
const getGameCategories = (
    params?: GameCategoryQueryParams
): Promise<ApiResponse<GameCategory[]>> => {
    return apiClient.get('/admin/game-categories', { params }).then(res => res.data);
};

/**
 * 获取单个游戏分类
 * GET /admin/game-categories/{id}
 */
const getGameCategory = (
    id: number
): Promise<ApiResponse<GameCategory>> => {
    return apiClient.get(`/admin/game-categories/${id}`).then(res => res.data);
};

/**
 * 创建游戏分类
 * POST /admin/game-categories
 */
const createGameCategory = (
    data: CreateGameCategoryDto
): Promise<ApiResponse<GameCategory>> => {
    return apiClient.post('/admin/game-categories', data).then(res => res.data);
};

/**
 * 更新游戏分类
 * PUT /admin/game-categories/{id}
 */
const updateGameCategory = (
    id: number,
    data: UpdateGameCategoryDto
): Promise<ApiResponse<GameCategory>> => {
    return apiClient.put(`/admin/game-categories/${id}`, data).then(res => res.data);
};

/**
 * 删除游戏分类
 * DELETE /admin/game-categories/{id}
 */
const deleteGameCategory = (
    id: number
): Promise<ApiResponse<void>> => {
    return apiClient.delete(`/admin/game-categories/${id}`).then(res => res.data);
};

/**
 * 批量更新游戏分类状态
 * POST /admin/game-categories/batch/status
 */
const batchUpdateCategoryStatus = (
    data: BatchUpdateCategoryStatusDto
): Promise<ApiResponse<BatchOperationResponse>> => {
    return apiClient.post('/admin/game-categories/batch/status', data).then(res => res.data);
};

/**
 * 批量删除游戏分类
 * POST /admin/game-categories/batch/delete
 */
const batchDeleteCategories = (
    data: BatchDeleteCategoriesDto
): Promise<ApiResponse<BatchOperationResponse>> => {
    return apiClient.post('/admin/game-categories/batch/delete', data).then(res => res.data);
};

// ============================================================================
// API Client - Game Rank
// ============================================================================

/**
 * 获取游戏段位列表
 * GET /admin/game-ranks
 */
const getGameRanks = (
    params?: GameRankQueryParams
): Promise<ApiResponse<GameRank[]>> => {
    return apiClient.get('/admin/game-ranks', { params }).then(res => res.data);
};

/**
 * 根据游戏ID获取段位列表
 * GET /admin/games/{gameId}/ranks
 */
const getGameRanksByGame = (
    gameId: number
): Promise<ApiResponse<GameRank[]>> => {
    return apiClient.get(`/admin/games/${gameId}/ranks`).then(res => res.data);
};

/**
 * 获取单个游戏段位
 * GET /admin/game-ranks/{id}
 */
const getGameRank = (
    id: number
): Promise<ApiResponse<GameRank>> => {
    return apiClient.get(`/admin/game-ranks/${id}`).then(res => res.data);
};

/**
 * 创建游戏段位
 * POST /admin/game-ranks
 */
const createGameRank = (
    data: CreateGameRankDto
): Promise<ApiResponse<GameRank>> => {
    return apiClient.post('/admin/game-ranks', data).then(res => res.data);
};

/**
 * 更新游戏段位
 * PUT /admin/game-ranks/{id}
 */
const updateGameRank = (
    id: number,
    data: UpdateGameRankDto
): Promise<ApiResponse<GameRank>> => {
    return apiClient.put(`/admin/game-ranks/${id}`, data).then(res => res.data);
};

/**
 * 删除游戏段位
 * DELETE /admin/game-ranks/{id}
 */
const deleteGameRank = (
    id: number
): Promise<ApiResponse<void>> => {
    return apiClient.delete(`/admin/game-ranks/${id}`).then(res => res.data);
};

/**
 * 批量删除游戏段位
 * POST /admin/game-ranks/batch/delete
 */
const batchDeleteGameRanks = (
    data: BatchDeleteGameRanksDto
): Promise<ApiResponse<{ deleted: number }>> => {
    return apiClient.post('/admin/game-ranks/batch/delete', data).then(res => res.data);
};

/**
 * 批量更新游戏段位状态
 * POST /admin/game-ranks/batch/status
 */
const batchUpdateRankStatus = (
    data: BatchUpdateRankStatusDto
): Promise<ApiResponse<{ updated: number }>> => {
    return apiClient.post('/admin/game-ranks/batch/status', data).then(res => res.data);
};

// ============================================================================
// Combined API Export
// ============================================================================

/**
 * Game Rank & Category API
 * 统一导出所有游戏分类和段位相关接口
 */
export const gameRankApi = {
    // Game Category APIs
    getGameCategories,
    getGameCategory,
    createGameCategory,
    updateGameCategory,
    deleteGameCategory,
    batchUpdateCategoryStatus,
    batchDeleteCategories,

    // Game Rank APIs
    getGameRanks,
    getGameRanksByGame,
    getGameRank,
    createGameRank,
    updateGameRank,
    deleteGameRank,
    batchDeleteGameRanks,
    batchUpdateRankStatus,
};

// Re-export types for convenience
export type {
    GameCategory,
    GameRank,
    Game,
    CreateGameCategoryDto,
    UpdateGameCategoryDto,
    BatchUpdateCategoryStatusDto,
    BatchDeleteCategoriesDto,
    GameCategoryQueryParams,
    BatchOperationResponse,
    CreateGameRankDto,
    UpdateGameRankDto,
    BatchDeleteGameRanksDto,
    BatchUpdateRankStatusDto,
    GameRankQueryParams,
};

export default gameRankApi;
