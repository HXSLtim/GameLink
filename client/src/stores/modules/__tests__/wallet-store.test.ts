/**
 * Wallet Store Tests
 * Tests for wallet balance, transactions, payments, and withdrawals
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useWalletStore, PaymentStatus, PaymentMethod, TransactionType, WithdrawStatus } from '../wallet-store';

// Mock http module
vi.mock('@/lib/http', () => ({
    http: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
    },
}));

import { http } from '@/lib/http';

const mockHttp = http as unknown as {
    get: ReturnType<typeof vi.fn>;
    post: ReturnType<typeof vi.fn>;
    put: ReturnType<typeof vi.fn>;
    delete: ReturnType<typeof vi.fn>;
};

describe('Wallet Store', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        // Reset store state
        useWalletStore.setState({
            wallet: null,
            transactions: [],
            rechargeOptions: [],
            currentPayment: null,
            paymentStatus: null,
            withdrawRecords: [],
            loading: false,
            error: null,
        });
    });

    afterEach(() => {
        vi.restoreAllMocks();
    });

    describe('Initial State', () => {
        it('should have correct initial state', () => {
            const state = useWalletStore.getState();

            expect(state.wallet).toBeNull();
            expect(state.transactions).toEqual([]);
            expect(state.rechargeOptions).toEqual([]);
            expect(state.currentPayment).toBeNull();
            expect(state.paymentStatus).toBeNull();
            expect(state.withdrawRecords).toEqual([]);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });
    });

    describe('fetchWallet', () => {
        it('should fetch wallet successfully', async () => {
            const mockWallet = {
                userId: 1,
                balanceCents: 100000,
                frozenCents: 5000,
                totalRechargeCents: 200000,
                totalConsumeCents: 100000,
                updatedAt: '2024-01-01T00:00:00Z',
            };

            mockHttp.get.mockResolvedValueOnce(mockWallet);

            await useWalletStore.getState().fetchWallet();

            const state = useWalletStore.getState();
            expect(state.wallet).toEqual(mockWallet);
            expect(state.wallet?.balanceCents).toBe(100000);
            expect(state.loading).toBe(false);
            expect(state.error).toBeNull();
        });

        it('should handle fetch wallet error', async () => {
            mockHttp.get.mockRejectedValueOnce(new Error('Network error'));

            await useWalletStore.getState().fetchWallet();

            const state = useWalletStore.getState();
            expect(state.wallet).toBeNull();
            expect(state.loading).toBe(false);
            expect(state.error).toBe('Network error');
        });
    });

    describe('fetchTransactions', () => {
        it('should fetch transactions successfully', async () => {
            const mockTransactions = {
                items: [
                    {
                        id: 1,
                        transactionNo: 'TXN001',
                        type: TransactionType.RECHARGE,
                        direction: 'in',
                        amountCents: 10000,
                        balanceAfterCents: 110000,
                        title: '充值',
                        createdAt: '2024-01-01T00:00:00Z',
                    },
                    {
                        id: 2,
                        transactionNo: 'TXN002',
                        type: TransactionType.CONSUME,
                        direction: 'out',
                        amountCents: 5000,
                        balanceAfterCents: 105000,
                        title: '订单消费',
                        orderId: 101,
                        createdAt: '2024-01-02T00:00:00Z',
                    },
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockTransactions);

            await useWalletStore.getState().fetchTransactions();

            const state = useWalletStore.getState();
            expect(state.transactions).toHaveLength(2);
            expect(state.transactions[0].type).toBe(TransactionType.RECHARGE);
            expect(state.loading).toBe(false);
        });

        it('should fetch transactions with filters', async () => {
            mockHttp.get.mockResolvedValueOnce({ items: [] });

            await useWalletStore.getState().fetchTransactions({
                type: [TransactionType.RECHARGE],
                direction: 'in',
                page: 2,
            });

            expect(mockHttp.get).toHaveBeenCalledWith('/user/wallet/transactions', {
                params: expect.objectContaining({
                    type: [TransactionType.RECHARGE],
                    direction: 'in',
                    page: 2,
                }),
            });
        });
    });

    describe('fetchRechargeOptions', () => {
        it('should fetch recharge options successfully', async () => {
            const mockOptions = [
                { id: 1, amountCents: 1000, bonusCents: 0, isPopular: false },
                { id: 2, amountCents: 5000, bonusCents: 500, label: '推荐', isPopular: true },
                { id: 3, amountCents: 10000, bonusCents: 1500, isPopular: false },
            ];

            mockHttp.get.mockResolvedValueOnce(mockOptions);

            await useWalletStore.getState().fetchRechargeOptions();

            const state = useWalletStore.getState();
            expect(state.rechargeOptions).toHaveLength(3);
            expect(state.rechargeOptions[1].isPopular).toBe(true);
        });
    });

    describe('recharge', () => {
        it('should create recharge successfully', async () => {
            const mockResponse = {
                rechargeId: 123,
                payParams: {
                    appId: 'wx123',
                    timeStamp: '1234567890',
                    nonceStr: 'abc123',
                    package: 'prepay_id=xxx',
                    signType: 'RSA',
                    paySign: 'signature',
                },
            };

            mockHttp.post.mockResolvedValueOnce(mockResponse);

            const result = await useWalletStore.getState().recharge(10000, 'wechat');

            expect(result.rechargeId).toBe(123);
            expect(result.payParams.appId).toBe('wx123');
            expect(mockHttp.post).toHaveBeenCalledWith('/user/recharge', {
                amountCents: 10000,
                method: 'wechat',
            });
        });

        it('should handle recharge error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Payment failed'));

            await expect(
                useWalletStore.getState().recharge(10000, 'wechat')
            ).rejects.toThrow('Payment failed');

            const state = useWalletStore.getState();
            expect(state.error).toBe('Payment failed');
        });
    });

    describe('createPayment', () => {
        it('should create payment successfully', async () => {
            const mockPayment = {
                paymentId: 456,
                paymentNo: 'PAY001',
                status: PaymentStatus.PENDING,
                breakdown: {
                    originalAmountCents: 10000,
                    couponDiscountCents: 1000,
                    vipDiscountCents: 500,
                    finalAmountCents: 8500,
                    walletAmountCents: 5000,
                    thirdPartyAmountCents: 3500,
                },
            };

            mockHttp.post.mockResolvedValueOnce(mockPayment);

            const result = await useWalletStore.getState().createPayment({
                orderId: 101,
                method: PaymentMethod.COMBINED,
                walletAmountCents: 5000,
                thirdPartyMethod: 'wechat',
                couponId: 1,
            });

            expect(result.paymentId).toBe(456);
            expect(result.breakdown.finalAmountCents).toBe(8500);

            const state = useWalletStore.getState();
            expect(state.currentPayment).toEqual(mockPayment);
            expect(state.paymentStatus).toBe(PaymentStatus.PENDING);
        });
    });

    describe('checkPaymentStatus', () => {
        it('should check payment status successfully', async () => {
            mockHttp.get.mockResolvedValueOnce({ status: PaymentStatus.COMPLETED });

            const status = await useWalletStore.getState().checkPaymentStatus(456);

            expect(status).toBe(PaymentStatus.COMPLETED);
            expect(useWalletStore.getState().paymentStatus).toBe(PaymentStatus.COMPLETED);
        });
    });

    describe('cancelPayment', () => {
        it('should cancel payment successfully', async () => {
            mockHttp.put.mockResolvedValueOnce({});

            await useWalletStore.getState().cancelPayment(456);

            expect(useWalletStore.getState().paymentStatus).toBe(PaymentStatus.CANCELED);
        });
    });

    describe('withdraw', () => {
        it('should withdraw successfully', async () => {
            // Set up wallet first
            useWalletStore.setState({
                wallet: {
                    userId: 1,
                    balanceCents: 100000,
                    frozenCents: 0,
                    totalRechargeCents: 200000,
                    totalConsumeCents: 100000,
                    updatedAt: '2024-01-01',
                },
            });

            mockHttp.post.mockResolvedValueOnce({});
            mockHttp.get.mockResolvedValueOnce({
                userId: 1,
                balanceCents: 50000,
                frozenCents: 0,
                totalRechargeCents: 200000,
                totalConsumeCents: 100000,
                updatedAt: '2024-01-02',
            });

            await useWalletStore.getState().withdraw(50000, 'bank_card', {
                bankName: 'Test Bank',
                cardNumber: '1234567890',
            });

            expect(mockHttp.post).toHaveBeenCalledWith('/player/withdraw', {
                amountCents: 50000,
                method: 'bank_card',
                accountInfo: {
                    bankName: 'Test Bank',
                    cardNumber: '1234567890',
                },
            });
        });

        it('should handle withdraw error', async () => {
            mockHttp.post.mockRejectedValueOnce(new Error('Insufficient balance'));

            await expect(
                useWalletStore.getState().withdraw(1000000, 'bank_card', {})
            ).rejects.toThrow('Insufficient balance');

            expect(useWalletStore.getState().error).toBe('Insufficient balance');
        });
    });

    describe('fetchWithdrawRecords', () => {
        it('should fetch withdraw records successfully', async () => {
            const mockRecords = {
                items: [
                    {
                        id: 1,
                        amountCents: 50000,
                        method: 'bank_card',
                        status: WithdrawStatus.COMPLETED,
                        createdAt: '2024-01-01',
                        completedAt: '2024-01-02',
                    },
                    {
                        id: 2,
                        amountCents: 30000,
                        method: 'wechat',
                        status: WithdrawStatus.PENDING,
                        createdAt: '2024-01-03',
                    },
                ],
            };

            mockHttp.get.mockResolvedValueOnce(mockRecords);

            await useWalletStore.getState().fetchWithdrawRecords();

            const state = useWalletStore.getState();
            expect(state.withdrawRecords).toHaveLength(2);
            expect(state.withdrawRecords[0].status).toBe(WithdrawStatus.COMPLETED);
        });
    });

    describe('Helper Methods', () => {
        it('getBalance should return correct balance', () => {
            useWalletStore.setState({
                wallet: {
                    userId: 1,
                    balanceCents: 50000,
                    frozenCents: 0,
                    totalRechargeCents: 100000,
                    totalConsumeCents: 50000,
                    updatedAt: '2024-01-01',
                },
            });

            expect(useWalletStore.getState().getBalance()).toBe(50000);
        });

        it('getBalance should return 0 when wallet is null', () => {
            expect(useWalletStore.getState().getBalance()).toBe(0);
        });

        it('canAfford should return true when balance is sufficient', () => {
            useWalletStore.setState({
                wallet: {
                    userId: 1,
                    balanceCents: 50000,
                    frozenCents: 0,
                    totalRechargeCents: 100000,
                    totalConsumeCents: 50000,
                    updatedAt: '2024-01-01',
                },
            });

            expect(useWalletStore.getState().canAfford(30000)).toBe(true);
            expect(useWalletStore.getState().canAfford(50000)).toBe(true);
        });

        it('canAfford should return false when balance is insufficient', () => {
            useWalletStore.setState({
                wallet: {
                    userId: 1,
                    balanceCents: 50000,
                    frozenCents: 0,
                    totalRechargeCents: 100000,
                    totalConsumeCents: 50000,
                    updatedAt: '2024-01-01',
                },
            });

            expect(useWalletStore.getState().canAfford(60000)).toBe(false);
        });

        it('canAfford should return false when wallet is null', () => {
            expect(useWalletStore.getState().canAfford(1000)).toBe(false);
        });
    });
});
