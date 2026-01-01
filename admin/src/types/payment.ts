/**
 * Payment & Operations Types
 * 支付运营类型定义
 *
 * Contains types for:
 * - Settlement Company (结算公司)
 * - Routing Rule (路由规则)
 * - Recharge (充值系统)
 * - Dispute (纠纷系统)
 */

// ============================================================================
// Settlement Company Types (结算公司)
// ============================================================================

/**
 * Company Type
 * 公司类型
 */
export type CompanyType = 'individual' | 'company';

/**
 * Company Status
 * 公司状态
 */
export type CompanyStatus = 'active' | 'suspended';

/**
 * Settlement Company
 * 结算公司
 */
export interface SettlementCompany {
    id: number;
    name: string;                     // 公司名称
    type: CompanyType;                // 个人/企业
    businessLicense?: string;         // 营业执照
    taxNumber?: string;               // 税号
    bankName?: string;                // 开户银行
    bankAccount?: string;             // 银行账号
    contactPerson?: string;           // 联系人
    contactPhone?: string;            // 联系电话
    status: CompanyStatus;
    playerCount: number;              // 关联陪玩师数量
    createdAt: string;
    updatedAt: string;
}

/**
 * Settlement Company History
 * 结算公司变更历史
 */
export interface SettlementCompanyHistory {
    id: number;
    settlementCompanyId: number;
    fieldName: string;                // 变更字段
    oldValue: string;                 // 旧值
    newValue: string;                 // 新值
    changedBy: number;                // 操作人ID
    changedByAdmin?: {
        id: number;
        name: string;
        email: string;
    };
    changedAt: string;
}

/**
 * Player Company Assignment
 * 陪玩师结算公司分配记录
 */
export interface PlayerCompanyAssignment {
    id: number;
    playerId: number;
    settlementCompanyId: number;
    effectiveDate: string;            // 生效日期
    reason: string;                   // 分配/变更原因
    assignedBy: number;               // 操作人ID
    assignedByAdmin?: {
        id: number;
        name: string;
        email: string;
    };
    createdAt: string;
    player?: {
        id: number;
        nickname: string;
        user?: {
            name: string;
            phone: string;
        };
    };
    settlementCompany?: {
        id: number;
        name: string;
        type: CompanyType;
    };
}

// ============================================================================
// Routing Rule Types (支付路由规则)
// ============================================================================

/**
 * Condition Field
 * 条件字段
 */
export type ConditionField = 'game_type' | 'service_type' | 'order_amount' | 'region';

/**
 * Condition Operator
 * 条件操作符
 */
export type ConditionOperator = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'lt' | 'between';

/**
 * Rule Status
 * 规则状态
 */
export type RuleStatus = 'active' | 'inactive';

/**
 * Routing Condition
 * 路由条件
 */
export interface RoutingCondition {
    field: ConditionField;            // 条件字段
    operator: ConditionOperator;      // 操作符
    value: string | number | string[] | number[]; // 条件值
}

/**
 * Collection Entity (收款主体)
 * 支付收款主体
 */
export interface CollectionEntity {
    id: number;
    name: string;                     // 主体名称
    creditCode: string;               // 统一社会信用代码
    taxRegistrationNo?: string;       // 税务登记号
    status: 'active' | 'inactive';
    isDefault: boolean;               // 是否默认主体
    totalCollectionCents: number;     // 累计收款金额（分）
    transactionCount: number;         // 交易笔数
    createdAt: string;
    updatedAt: string;
}

/**
 * Routing Rule
 * 支付路由规则
 */
export interface RoutingRule {
    id: number;
    name: string;                     // 规则名称
    priority: number;                 // 优先级（数字越大越优先）
    conditions: RoutingCondition[];   // 匹配条件列表
    targetEntityId: number;           // 目标收款主体ID
    status: RuleStatus;
    description?: string;
    createdBy: number;                // 创建人ID
    updatedBy?: number;               // 更新人ID
    targetEntity?: CollectionEntity;  // 关联的收款主体
    createdAt: string;
    updatedAt: string;
}

/**
 * Routing Rule History
 * 路由规则历史记录
 */
export interface RoutingRuleHistory {
    id: number;
    routingRuleId: number;
    fieldName: string;
    oldValue: string;
    newValue: string;
    changedBy: number;
    createdAt: string;
    updatedAt: string;
}

/**
 * Routing Test Request
 * 路由测试请求
 */
export interface RoutingTestRequest {
    gameType?: string;
    serviceType?: string;
    amountCents?: number;
    region?: string;
}

/**
 * Routing Test Response
 * 路由测试响应
 */
export interface RoutingTestResponse {
    matchedRuleId?: number;           // 匹配的规则ID
    matchedRuleName?: string;         // 匹配的规则名称
    collectionEntityId: number;       // 收款主体ID
    entityName: string;               // 收款主体名称
    merchantNo: string;               // 商户号
    isDefault: boolean;               // 是否使用默认规则
    matchDetails?: RoutingCondition[]; // 匹配详情
}

// ============================================================================
// Recharge Types (充值系统)
// ============================================================================

/**
 * Recharge Status
 * 充值状态
 */
export type RechargeStatus = 'pending' | 'paid' | 'failed' | 'refunded' | 'expired';

/**
 * Recharge Option
 * 充值档位
 */
export interface RechargeOption {
    id: number;
    name: string;                     // 档位名称
    amountCents: number;              // 充值金额（分）
    bonusCents: number;               // 赠送金额（分）
    originalCents?: number;           // 原价（分）
    discountPercent?: number;         // 折扣百分比
    description?: string;
    tag?: string;                     // 标签（如"热销"）
    iconUrl?: string;                 // 图标
    sortOrder: number;
    isActive: boolean;                // 是否启用
    isRecommended: boolean;           // 是否推荐
    couponTemplateId?: number;        // 赠送优惠券模板ID
    couponCount: number;              // 赠送优惠券数量
    minVipLevel?: number;             // 最低VIP等级要求
    perUserLimit: number;             // 每用户限购次数
    totalLimit: number;               // 总限购次数
    createdAt: string;
    updatedAt: string;
}

/**
 * Recharge Record
 * 充值记录
 */
export interface RechargeRecord {
    id: number;
    orderNo: string;                  // 充值订单号
    userId: number;                   // 用户ID
    optionId: number;                 // 档位ID
    amountCents: number;              // 充值金额（分）
    bonusCents: number;               // 赠送金额（分）
    totalCents: number;               // 总到账金额（分）
    status: RechargeStatus;
    paymentChannel?: string;          // 支付渠道
    paymentNo?: string;               // 第三方支付流水号
    paidAt?: string;                  // 支付时间
    refundedAt?: string;              // 退款时间
    refundReason?: string;            // 退款原因
    createdAt: string;
    updatedAt: string;
    user?: {
        id: number;
        name: string;
        avatarUrl?: string;
    };
    option?: {
        id: number;
        name: string;
        amountCents: number;
        bonusCents: number;
    };
}

/**
 * Recharge Statistics
 * 充值统计
 */
export interface RechargeStats {
    totalOrders: number;
    totalAmountCents: number;
    totalBonusCents: number;
    paidOrders: number;
    pendingOrders: number;
    failedOrders: number;
    refundedOrders: number;
    todayOrders: number;
    todayAmountCents: number;
    monthOrders: number;
    monthAmountCents: number;
}

// ============================================================================
// Dispute Types (纠纷系统)
// ============================================================================

/**
 * Dispute Status
 * 纠纷状态
 */
export type DisputeStatus =
    | 'pending'    // 待处理
    | 'assigned'   // 已指派
    | 'mediating'  // 调解中
    | 'resolved'   // 已解决
    | 'rejected'   // 已驳回
    | 'canceled';  // 已取消

/**
 * Dispute Resolution
 * 纠纷处理结果
 */
export type DisputeResolution =
    | 'refund'     // 全额退款
    | 'partial'    // 部分退款
    | 'reassign'   // 重新指派
    | 'reject'     // 驳回
    | 'pending';   // 待决定

/**
 * Dispute Initiator Type
 * 纠纷发起人类型
 */
export type DisputeInitiatorType = 'user' | 'player';

/**
 * Dispute Type
 * 纠纷类型
 */
export type DisputeType =
    | 'service_quality'       // 服务质量问题
    | 'bad_attitude'          // 态度问题
    | 'incomplete_service'    // 未完成服务
    | 'user_not_cooperative'  // 用户不配合
    | 'user_harassment'       // 用户骚扰
    | 'other';                // 其他

/**
 * Assignment Source
 * 指派来源
 */
export type AssignmentSource = 'system' | 'manual' | 'team';

/**
 * Order Dispute
 * 订单纠纷
 */
export interface Dispute {
    id: number;
    orderId: number;
    orderNo?: string;
    initiatorId: number;             // 发起人ID
    initiatorName?: string;
    initiatorType: DisputeInitiatorType;
    type: DisputeType;
    status: DisputeStatus;
    reason: string;                  // 纠纷原因
    evidenceUrls?: string[];         // 证据图片URL列表
    evidenceText?: string;           // 文字证据
    chatSnapshotId?: number;         // 聊天快照ID

    // 双客服机制 (Dual-CS Mechanism)
    originalServiceId?: number;      // 原始客服ID
    originalServiceName?: string;
    assignedServiceId?: number;      // 独立客服ID
    assignedServiceName?: string;

    // SLA 信息 (SLA Information)
    slaDeadline?: string;            // 处理截止时间
    slaBreached: boolean;            // 是否超时
    slaBreachedAt?: string;          // 超时时间

    // 处理信息 (Resolution Information)
    resolution?: DisputeResolution;
    resolvedBy?: number;
    resolvedByName?: string;
    resolvedAt?: string;
    resolveRemark?: string;

    // 回退信息 (Rollback Information)
    rolledBackAt?: string;
    rolledBackByUserId?: number;
    rollbackReason?: string;

    // 追踪信息 (Trace Information)
    traceId: string;

    createdAt: string;
    updatedAt: string;
}

/**
 * Dispute Statistics
 * 纠纷统计
 */
export interface DisputeStats {
    total: number;
    pending: number;
    assigned: number;
    mediating: number;
    resolved: number;
    rejected: number;
    canceled: number;
    slaBreached: number;             // 超时数量
}

// ============================================================================
// Common Batch Operation Types
// ============================================================================

/**
 * Batch Operation Result
 * 批量操作结果
 */
export interface BatchOperationResult {
    successCount: number;
    failedCount: number;
    totalCount?: number;
    failedItems?: Array<{ id: number; message: string; error?: string }>;
    successItems?: number[];
    errors?: string[];
}

/**
 * Batch Operation Error
 * 批量操作错误详情
 */
export interface BatchOperationError {
    id: number;
    message: string;
    error?: string;
}
