import { UserRole, UserStatus, VerificationStatus } from '../types/user.types';

/**
 * 格式化用户角色
 */
export const formatUserRole = (role: UserRole): string => {
  const roleMap: Record<UserRole, string> = {
    [UserRole.USER]: '普通用户',
    [UserRole.PLAYER]: '陪玩师',
    [UserRole.ADMIN]: '管理员',
  };
  return roleMap[role] || role;
};

/**
 * 获取用户角色颜色
 */
export const getUserRoleColor = (role: UserRole): 'default' | 'info' | 'success' | 'warning' => {
  const colorMap: Record<UserRole, 'default' | 'info' | 'success' | 'warning'> = {
    [UserRole.USER]: 'default',
    [UserRole.PLAYER]: 'info',
    [UserRole.ADMIN]: 'warning',
  };
  return colorMap[role] || 'default';
};

/**
 * 格式化用户状态
 */
export const formatUserStatus = (status: UserStatus): string => {
  const statusMap: Record<UserStatus, string> = {
    [UserStatus.ACTIVE]: '正常',
    [UserStatus.SUSPENDED]: '暂停',
    [UserStatus.BANNED]: '封禁',
  };
  return statusMap[status] || status;
};

/**
 * 获取用户状态颜色
 */
export const getUserStatusColor = (
  status: UserStatus,
): 'default' | 'info' | 'success' | 'warning' | 'error' => {
  const colorMap: Record<UserStatus, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
    [UserStatus.ACTIVE]: 'success',
    [UserStatus.SUSPENDED]: 'warning',
    [UserStatus.BANNED]: 'error',
  };
  return colorMap[status] || 'default';
};

/**
 * 格式化认证状态
 */
export const formatVerificationStatus = (status: VerificationStatus): string => {
  const statusMap: Record<VerificationStatus, string> = {
    [VerificationStatus.PENDING]: '待认证',
    [VerificationStatus.VERIFIED]: '已认证',
    [VerificationStatus.REJECTED]: '已拒绝',
  };
  return statusMap[status] || status;
};

/**
 * 获取认证状态颜色
 */
export const getVerificationStatusColor = (
  status: VerificationStatus,
): 'default' | 'info' | 'success' | 'warning' | 'error' => {
  const colorMap: Record<VerificationStatus, 'default' | 'info' | 'success' | 'warning' | 'error'> =
    {
      [VerificationStatus.PENDING]: 'warning',
      [VerificationStatus.VERIFIED]: 'success',
      [VerificationStatus.REJECTED]: 'error',
    };
  return colorMap[status] || 'default';
};

/**
 * 格式化金额（分 -> 元）
 */
export const formatPrice = (cents: number): string => {
  return `¥${(cents / 100).toFixed(2)}`;
};

/**
 * 格式化时薪
 */
export const formatHourlyRate = (cents: number): string => {
  return `${formatPrice(cents)}/小时`;
};

/**
 * 格式化手机号（脱敏）
 */
export const formatPhone = (phone?: string): string => {
  if (!phone) return '-';
  if (phone.length === 11) {
    return `${phone.slice(0, 3)}****${phone.slice(7)}`;
  }
  return phone;
};

/**
 * 格式化邮箱（脱敏）
 */
export const formatEmail = (email?: string): string => {
  if (!email) return '-';
  const [name, domain] = email.split('@');
  if (name && domain) {
    const maskedName = name.length > 2 ? `${name[0]}***${name[name.length - 1]}` : name;
    return `${maskedName}@${domain}`;
  }
  return email;
};

/**
 * 格式化评分
 */
export const formatRating = (rating: number): string => {
  return `${rating.toFixed(1)} 分`;
};
