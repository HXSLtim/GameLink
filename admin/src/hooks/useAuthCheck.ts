/**
 * 统一的认证检查 Hook
 *
 * 解决 Zustand persist 的异步恢复问题
 * 确保组件只在状态恢复后才渲染
 * 遵循 Super Dev 最佳实践：精确订阅，避免不必要的重渲染
 *
 * @module hooks/useAuthCheck
 *
 * @example
 * ```tsx
 * // 基础使用
 * const { userInfo, isAuthenticated, isReady } = useAuthCheck();
 *
 * if (!isReady) {
 *   return <Spin tip="加载中..." />;
 * }
 *
 * if (!isAuthenticated) {
 *   return <div>请先登录</div>;
 * }
 *
 * return <div>欢迎, {userInfo?.name}</div>;
 * ```
 *
 * @example
 * ```tsx
 * // 使用便捷 hooks（推荐）
 * const authUser = useUserInfo();
 * const isAuthenticated = useIsAuthenticated();
 * const isHydrated = useIsHydrated();
 * ```
 */

import { useMemo } from 'react';
import { useUserInfo, useIsAuthenticated, useIsHydrated } from '@/stores/modules/authStore';
import type { UserInfo } from '@/stores/modules/authStore';

/**
 * 认证检查结果接口
 */
export interface AuthCheckResult {
  /** 用户信息 */
  userInfo: UserInfo | null;
  /** 是否已认证 */
  isAuthenticated: boolean;
  /** 是否正在从 localStorage 恢复 */
  isHydrating: boolean;
  /** 是否准备就绪 (水合完成 或 无需认证) */
  isReady: boolean;
}

/**
 * 统一的认证检查 Hook
 *
 * **功能**:
 * - 提供统一的认证状态接口
 * - 自动处理 Zustand persist 的异步恢复
 * - 避免组件在水合完成前渲染
 *
 * **优势**:
 * - 精确订阅，避免不必要的重渲染
 * - 类型安全
 * - 统一的加载状态处理
 *
 * **性能优化**:
 * - 使用 useUserInfo 等选择器，只订阅需要的状态
 * - 避免订阅整个 store
 */
export function useAuthCheck(): AuthCheckResult {
  // Super Dev 最佳实践: 使用选择器精确订阅
  const userInfo = useUserInfo();
  const isAuthenticated = useIsAuthenticated();
  const isHydrated = useIsHydrated();

  // 计算派生状态
  const result = useMemo<AuthCheckResult>(
    () => ({
      userInfo,
      isAuthenticated,
      isHydrating: !isHydrated,
      isReady: isHydrated,
    }),
    [userInfo, isAuthenticated, isHydrated]
  );

  return result;
}

/**
 * 导出便捷的选择器 hooks
 * 这些 hooks 可以直接使用，提供更好的类型推导和性能
 *
 * @example
 * ```tsx
 * // 直接导入使用，无需解构
 * import { useUserInfo, useIsAuthenticated } from '@/hooks/useAuthCheck';
 *
 * const authUser = useUserInfo();
 * const isAuthenticated = useIsAuthenticated();
 * ```
 */
export { useUserInfo, useIsAuthenticated, useIsHydrated } from '@/stores/modules/authStore';

/**
 * 导出 token hook
 * 用于需要直接访问 token 的场景（如 API 拦截器）
 */
export { useAuthToken } from '@/stores/modules/authStore';

/**
 * 导出加载状态和错误 hooks
 */
export { useAuthLoading, useAuthError } from '@/stores/modules/authStore';

/**
 * 导出管理员检查 hook
 */
export { useIsAdmin } from '@/stores/modules/authStore';

export default useAuthCheck;
