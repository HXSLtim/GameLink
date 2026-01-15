import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { http } from '@/lib/http';

// ============ Enums ============

export const PaymentStatus = {
    PENDING: 'pending',           // 待支付
    PROCESSING: 'processing',     // 处理中
    COMPLETED: 'completed',       // 已完成
    FAILED: 'failed',             // 失败
    REFUNDED: 'refunded',         // 已退款
    PARTIAL_REFUNDED: 'partial_refunded',  // 部分退款
    CANCELED: 'canceled'          // 已取消
} as const;

export type PaymentStatus = typeof PaymentStatus[keyof typeof PaymentStatus];

export const PaymentMethod = {
    WECHAT: 'wechat',             // 微信支付
    ALIPAY: 'alipay',             // 支付宝
    WALLET: 'wallet',             // 余额支付
    COMBINED: 'combined'          // 组合支付
} as const;

export type PaymentMethod = typeof PaymentMethod[keyof typeof PaymentMethod];

export const TransactionType = {
    RECHARGE: 'recharge',         // 充值
    CONSUME: 'consume',           // 消费
    REFUND: 'refund',             // 退款
    INCOME: 'income',             // 收入 (陪玩师)
    WITHDRAW: 'withdraw',         // 提现
    BONUS: 'bonus'                // 奖励
} as const;

export type TransactionType = typeof TransactionType[keyof typeof TransactionType];

export const WithdrawStatus = {
    PENDING: 'pending',           // 待处理
    PROCESSING: 'processing',     // 处理中
    COMPLETED: 'completed',       // 已完成
    FAILED: 'failed',             // 失败
    REJECTED: 'rejected'          // 已拒绝
} as const;

export type WithdrawStatus = typeof WithdrawStatus[keyof typeof WithdrawStatus];

// ============ Interfaces ============

export interface Wallet {
    userId: number;
    balanceCents: number;          // 可用余额
    frozenCents: number;           // 冻结金额
    incomeCents?: number;          // 累计收入 (陪玩师)
    withdrawableCents?: number;    // 可提现金额
    totalRechargeCents: number;    // 累计充值
    totalConsumeCents: number;     // 累计消费
    updatedAt: string;
}

export interface Transaction {
    id: number;
    transactionNo: string;
    type: TransactionType;
    direction: 'in' | 'out';       // 收入/支出
    amountCents: number;
    balanceAfterCents: number;     // 交易后余额
    orderId?: number;
    orderNo?: string;
    rechargeId?: number;
    withdrawId?: number;
    title: string;
    description?: string;
    createdAt: string;
}

export interface RechargeOption {
    id: number;
    amountCents: number;           // 充值金额
    bonusCents: number;            // 赠送金额
    label?: string;                // 标签 (如 "推荐")
    isPopular: boolean;            // 是否热门
}

export interface PaymentBreakdown {
    originalAmountCents: number;   // 原价
    couponDiscountCents: number;   // 优惠券抵扣
    vipDiscountCents: number;      // VIP 折扣
    finalAmountCents: number;      // 实付金额
    walletAmountCents: number;     // 余额支付
    thirdPartyAmountCents: number; // 第三方支付
}

export interface WeChatPayParams {
    appId: string;
    timeStamp: string;
    nonceStr: string;
    package: string;
    signType: string;
    paySign: string;
}

export interface CreatePaymentRequest {
    orderId: number;
    method: PaymentMethod;
    walletAmountCents?: number;    // 组合支付时余额部分
    thirdPartyMethod?: 'wechat' | 'alipay';
    couponId?: number;
}

export interface CreatePaymentResponse {
    paymentId: number;
    paymentNo: string;
    status: PaymentStatus;
    payParams?: WeChatPayParams;
    breakdown: PaymentBreakdown;
}

export interface WithdrawRecord {
    id: number;
    amountCents: number;
    method: 'bank_card' | 'wechat' | 'alipay';
    status: WithdrawStatus;
    createdAt: string;
    processedAt?: string;
    completedAt?: string;
    failureReason?: string;
}

// ============ State & Actions ============

export interface WalletState {
    wallet: Wallet | null;
    transactions: Transaction[];
    rechargeOptions: RechargeOption[];
    currentPayment: CreatePaymentResponse | null;
    paymentStatus: PaymentStatus | null;
    withdrawRecords: WithdrawRecord[];
    loading: boolean;
    error: string | null;
}

export interface WalletActions {
    fetchWallet: () => Promise<void>;
    fetchTransactions: (params?: { type?: TransactionType[]; direction?: 'in' | 'out'; page?: number }) => Promise<void>;
    fetchRechargeOptions: () => Promise<void>;
    recharge: (amountCents: number, method: 'wechat' | 'alipay') => Promise<{ rechargeId: number; payParams: WeChatPayParams }>;
    createPayment: (request: CreatePaymentRequest) => Promise<CreatePaymentResponse>;
    checkPaymentStatus: (paymentId: number) => Promise<PaymentStatus>;
    cancelPayment: (paymentId: number) => Promise<void>;
    withdraw: (amountCents: number, method: 'bank_card' | 'wechat' | 'alipay', accountInfo: Record<string, string>) => Promise<void>;
    fetchWithdrawRecords: () => Promise<void>;
    getBalance: () => number;
    canAfford: (amountCents: number) => boolean;
}

// ============ Store ============

export const useWalletStore = create<WalletState & WalletActions>()(
    persist(
        (set, get) => ({
            wallet: null,
            transactions: [],
            rechargeOptions: [],
            currentPayment: null,
            paymentStatus: null,
            withdrawRecords: [],
            loading: false,
            error: null,

            fetchWallet: async () => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<Wallet>('/user/wallet');
                    set({ wallet: data, loading: false });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch wallet' });
                }
            },

            fetchTransactions: async (params = {}) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.get<{ items: Transaction[] }>('/user/wallet/transactions', {
                        params: { page: params.page || 1, pageSize: 20, ...params }
                    });
                    set({ transactions: data.items || [], loading: false });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to fetch transactions' });
                }
            },

            fetchRechargeOptions: async () => {
                try {
                    const data = await http.get<RechargeOption[]>('/user/recharge/options');
                    set({ rechargeOptions: data });
                } catch (err) {
                    console.error('Failed to fetch recharge options:', err);
                }
            },

            recharge: async (amountCents, method) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.post<{ rechargeId: number; payParams: WeChatPayParams }>('/user/recharge', {
                        amountCents,
                        method
                    });
                    set({ loading: false });
                    return data;
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Recharge failed' });
                    throw err;
                }
            },

            createPayment: async (request) => {
                set({ loading: true, error: null });
                try {
                    const data = await http.post<CreatePaymentResponse>('/payments', request);
                    set({ currentPayment: data, paymentStatus: PaymentStatus.PENDING, loading: false });
                    return data;
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Failed to create payment' });
                    throw err;
                }
            },

            checkPaymentStatus: async (paymentId) => {
                try {
                    const data = await http.get<{ status: PaymentStatus }>(`/payments/${paymentId}/status`);
                    set({ paymentStatus: data.status });
                    return data.status;
                } catch (err) {
                    console.error('Failed to check payment status:', err);
                    throw err;
                }
            },

            cancelPayment: async (paymentId) => {
                try {
                    await http.put(`/payments/${paymentId}/cancel`);
                    set({ paymentStatus: PaymentStatus.CANCELED });
                } catch (err) {
                    console.error('Failed to cancel payment:', err);
                    throw err;
                }
            },

            withdraw: async (amountCents, method, accountInfo) => {
                set({ loading: true, error: null });
                try {
                    await http.post('/player/withdraw', { amountCents, method, accountInfo });
                    await get().fetchWallet();
                    set({ loading: false });
                } catch (err) {
                    set({ loading: false, error: err instanceof Error ? err.message : 'Withdrawal failed' });
                    throw err;
                }
            },

            fetchWithdrawRecords: async () => {
                try {
                    const data = await http.get<{ items: WithdrawRecord[] }>('/player/withdraw/records');
                    set({ withdrawRecords: data.items || [] });
                } catch (err) {
                    console.error('Failed to fetch withdraw records:', err);
                }
            },

            getBalance: () => {
                const { wallet } = get();
                return wallet?.balanceCents || 0;
            },

            canAfford: (amountCents) => {
                const { wallet } = get();
                return (wallet?.balanceCents || 0) >= amountCents;
            }
        }),
        {
            name: 'wallet-storage',
            partialize: (state) => ({
                wallet: state.wallet
            })
        }
    )
);
