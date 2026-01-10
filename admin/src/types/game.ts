/**
 * Game Module Types
 * 游戏模块类型定义
 *
 * Contains types for:
 * - Game Category (游戏分类)
 * - Game Rank (游戏段位)
 * - Player Rank (陪玩师段位)
 * - Game (游戏基础信息)
 */

// ============================================================================
// Game Category Types (游戏分类)
// ============================================================================

/**
 * Game Category
 * 游戏分类
 * 对应后端 model.GameCategory
 */
export interface GameCategory {
    id: number;
    name: string;                     // 分类名称
    description?: string;             // 分类描述
    iconUrl?: string;                 // 分类图标
    sortOrder: number;                // 排序
    isActive: boolean;                // 是否启用
    createdAt: string;
    updatedAt: string;
    // 关联数据（非持久化字段）
    games?: Game[];                   // 该分类下的游戏列表
    serviceItems?: GameServiceItem[]; // 该分类下的服务项
}

// ============================================================================
// Game Types (游戏基础信息)
// ============================================================================

/**
 * Game
 * 游戏基础信息
 */
export interface Game {
    id: number;
    key: string;                      // 游戏唯一标识符
    name: string;                     // 游戏名称
    category?: string;                // 分类名称
    categoryId?: number;              // 分类ID
    iconUrl?: string;                 // 游戏图标
    coverUrl?: string;                // 游戏封面图
    description?: string;             // 游戏描述
    isActive: boolean;                // 是否启用
    sortOrder: number;                // 排序
}

// ============================================================================
// Game Rank Types (游戏段位)
// ============================================================================

/**
 * Game Rank
 * 游戏段位配置
 * 对应后端 model.GameRank
 */
export interface GameRank {
    id: number;
    gameId: number;                   // 关联游戏ID
    name: string;                     // 段位名称
    level: number;                    // 段位等级
    priceCents: number;               // 价格（分）- 该段位的陪玩价格
    iconUrl?: string;                 // 段位图标
    color?: string;                   // 段位颜色
    description?: string;             // 段位描述
    sortOrder: number;                // 排序
    isActive: boolean;                // 是否启用
    createdAt: string;
    updatedAt: string;
    // 关联数据（非持久化字段）
    game?: Game;                      // 关联的游戏
    gameName?: string;                // 游戏名称（冗余字段）
}

/**
 * Game Rank with Price Range
 * 带价格区间的游戏段位（用于展示）
 */
export interface GameRankWithPriceRange extends GameRank {
    minPriceCents?: number;           // 最低价格
    maxPriceCents?: number;           // 最高价格
}

// ============================================================================
// Player Rank Types (陪玩师段位)
// ============================================================================

/**
 * Player Rank Level
 * 陪玩师段位等级
 */
export type PlayerRankLevel =
    | 'bronze'    // 青铜
    | 'silver'    // 白银
    | 'gold'      // 黄金
    | 'platinum'  // 铂金
    | 'diamond'   // 钻石
    | 'master'    // 大师
    | 'champion'; // 王者

/**
 * Player Rank Configuration
 * 陪玩师段位配置
 */
export interface PlayerRank {
    id: number;
    level: PlayerRankLevel;           // 段位等级
    title: string;                    // 段位称号
    iconUrl?: string;                 // 段位图标
    color?: string;                   // 段位颜色
    minOrders: number;                // 最低订单数要求
    minRating: number;                // 最低评分要求
    minCompletionRate: number;        // 最低完成率要求 (0-1)
    commissionDiscountRate: number;   // 佣金折扣率 (0-1) - 段位越高折扣越少
    priceBonusRate: number;           // 价格加成比例 (0-1) - 段位越高加成越多
    benefits: string;                 // 段位权益描述
    sortOrder: number;                // 排序
    isActive: boolean;                // 是否启用
    createdAt: string;
    updatedAt: string;
}

/**
 * Player Rank Progress
 * 陪玩师段位进度
 */
export interface PlayerRankProgress {
    currentRankId: number;            // 当前段位ID
    currentRank: PlayerRank;          // 当前段位详情
    nextRankId?: number;              // 下一级段位ID
    nextRank?: PlayerRank;            // 下一级段位详情
    progress: number;                 // 当前进度 (0-1)
    currentOrders: number;            // 当前订单数
    requiredOrders: number;           // 晋升所需订单数
    currentRating: number;            // 当前评分
    requiredRating: number;           // 晋升所需评分
    currentCompletionRate: number;    // 当前完成率
    requiredCompletionRate: number;   // 晋升所需完成率
}

/**
 * Player Rank Statistics
 * 陪玩师段位统计
 */
export interface PlayerRankStats {
    rankId: number;
    rankLevel: PlayerRankLevel;
    rankTitle: string;
    playerCount: number;              // 该段位陪玩师数量
    totalOrders: number;              // 总订单数
    totalIncomeCents: number;         // 总收益（分）
    averageRating: number;            // 平均评分
    averageCompletionRate: number;    // 平均完成率
}

// ============================================================================
// Game Service Item Types (游戏服务项)
// ============================================================================

/**
 * Service Item Type
 * 服务项类型
 */
export type ServiceItemType =
    | 'solo'      // 单人陪玩
    | 'team'      // 组队陪玩
    | 'coach'     // 教学指导
    | 'boost'     // 代练
    | 'companion';// 陪伴聊天

/**
 * Game Service Item
 * 游戏服务项
 */
export interface GameServiceItem {
    id: number;
    gameId: number;                   // 关联游戏ID
    gameRankId?: number;              // 关联游戏段位ID
    name: string;                     // 服务项名称
    type: ServiceItemType;            // 服务项类型
    description?: string;             // 描述
    basePriceCents: number;           // 基础价格（分）/小时
    commissionRate: number;           // 佣金比例 (0-1)
    durationMinutes: number;          // 标准时长（分钟）
    minPlayers: number;               // 最少人数
    maxPlayers: number;               // 最多人数
    iconUrl?: string;                 // 图标
    sortOrder: number;
    isActive: boolean;
    createdAt: string;
    updatedAt: string;
    // 关联数据
    game?: Game;
    gameRank?: GameRank;
}

// ============================================================================
// Helper Types
// ============================================================================

/**
 * Game with Ranks
 * 带段位的游戏信息
 */
export interface GameWithRanks extends Omit<Game, 'category'> {
    ranks?: GameRank[];               // 该游戏的段位列表
    serviceItems?: GameServiceItem[]; // 该游戏的服务项列表
    category?: GameCategory;          // 游戏分类对象（覆盖基类的category: string）
}

/**
 * Game Category with Games
 * 带游戏的分类信息
 */
export interface GameCategoryWithGames extends Omit<GameCategory, 'games'> {
    games?: GameWithRanks[];          // 该分类下的游戏列表（含段位）
}

// ============================================================================
// Display Constants (UI展示用常量)
// ============================================================================

/**
 * Player Rank Level Labels
 * 陪玩师段位等级显示标签
 */
export const PLAYER_RANK_LEVEL_LABELS: Record<PlayerRankLevel, string> = {
    bronze: '青铜',
    silver: '白银',
    gold: '黄金',
    platinum: '铂金',
    diamond: '钻石',
    master: '大师',
    champion: '王者',
};

/**
 * Player Rank Level Colors
 * 陪玩师段位等级颜色（用于Ant Design Tag）
 */
export const PLAYER_RANK_LEVEL_COLORS: Record<PlayerRankLevel, string> = {
    bronze: '#CD7F32',    // 青铜色
    silver: '#C0C0C0',    // 银色
    gold: '#FFD700',      // 金色
    platinum: '#E5E4E2',  // 铂金
    diamond: '#B9F2FF',   // 钻石蓝
    master: '#9370DB',    // 紫色
    champion: '#FF4500',  // 橙红色
};

/**
 * Service Item Type Labels
 * 服务项类型显示标签
 */
export const SERVICE_ITEM_TYPE_LABELS: Record<ServiceItemType, string> = {
    solo: '单人陪玩',
    team: '组队陪玩',
    coach: '教学指导',
    boost: '代练服务',
    companion: '陪伴聊天',
};

/**
 * Service Item Type Icons
 * 服务项类型图标
 */
export const SERVICE_ITEM_TYPE_ICONS: Record<ServiceItemType, string> = {
    solo: 'user',
    team: 'team',
    coach: 'book',
    boost: 'trophy',
    companion: 'message',
};
