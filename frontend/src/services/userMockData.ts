import {
  User,
  UserDetail,
  UserRole,
  UserStatus,
  VerificationStatus,
  PlayerInfo,
} from '../types/user.types';

/**
 * 生成 Mock 用户数据
 */
export const generateMockUsers = (): User[] => {
  const users: User[] = [];

  // 管理员
  users.push({
    id: 1,
    name: '系统管理员',
    email: 'admin@gamelink.com',
    phone: '13800138000',
    avatar_url: '',
    role: UserRole.ADMIN,
    status: UserStatus.ACTIVE,
    last_login_at: '2025-01-05 10:30:00',
    created_at: '2024-01-01 00:00:00',
    updated_at: '2025-01-05 10:30:00',
  });

  // 陪玩师
  const playerNames = [
    '王者荣耀大神',
    '绝地求生专家',
    '英雄联盟教练',
    '和平精英陪练',
    'CSGO老司机',
  ];
  for (let i = 0; i < 15; i++) {
    users.push({
      id: i + 2,
      name: playerNames[i % playerNames.length] + ` ${i + 1}`,
      email: `player${i + 1}@gamelink.com`,
      phone: `138${String(i + 1).padStart(8, '0')}`,
      avatar_url: `https://api.dicebear.com/7.x/avataaars/svg?seed=${i + 2}`,
      role: UserRole.PLAYER,
      status:
        i % 10 === 0 ? UserStatus.SUSPENDED : i % 15 === 0 ? UserStatus.BANNED : UserStatus.ACTIVE,
      last_login_at: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(
        2024,
        Math.floor(Math.random() * 12),
        Math.floor(Math.random() * 28),
      ).toISOString(),
      updated_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
    });
  }

  // 普通用户
  const userNames = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十'];
  for (let i = 0; i < 30; i++) {
    users.push({
      id: i + 17,
      name: userNames[i % userNames.length] + ` ${Math.floor(i / userNames.length) + 1}`,
      email: `user${i + 1}@example.com`,
      phone: `139${String(i + 1).padStart(8, '0')}`,
      avatar_url: `https://api.dicebear.com/7.x/avataaars/svg?seed=${i + 17}`,
      role: UserRole.USER,
      status: i % 20 === 0 ? UserStatus.SUSPENDED : UserStatus.ACTIVE,
      last_login_at: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      created_at: new Date(
        2024,
        Math.floor(Math.random() * 12),
        Math.floor(Math.random() * 28),
      ).toISOString(),
      updated_at: new Date(Date.now() - Math.random() * 60 * 24 * 60 * 60 * 1000).toISOString(),
    });
  }

  return users;
};

// 全局 Mock 数据
const mockUsers = generateMockUsers();

/**
 * 获取用户列表（支持筛选和分页）
 */
export const getMockUserList = (params: {
  page?: number;
  page_size?: number;
  keyword?: string;
  role?: UserRole;
  status?: UserStatus;
}) => {
  const { page = 1, page_size = 10, keyword, role, status } = params;

  let filtered = [...mockUsers];

  // 关键词搜索
  if (keyword) {
    const kw = keyword.toLowerCase();
    filtered = filtered.filter(
      (user) =>
        user.name.toLowerCase().includes(kw) ||
        user.email?.toLowerCase().includes(kw) ||
        user.phone?.includes(kw),
    );
  }

  // 角色筛选
  if (role) {
    filtered = filtered.filter((user) => user.role === role);
  }

  // 状态筛选
  if (status) {
    filtered = filtered.filter((user) => user.status === status);
  }

  // 分页
  const start = (page - 1) * page_size;
  const end = start + page_size;
  const list = filtered.slice(start, end);

  return {
    list,
    total: filtered.length,
    page,
    page_size,
  };
};

/**
 * 获取用户详情
 */
export const getMockUserDetail = (id: number): UserDetail | null => {
  const user = mockUsers.find((u) => u.id === id);
  if (!user) return null;

  const detail: UserDetail = {
    ...user,
    order_count: Math.floor(Math.random() * 100),
    total_spent: Math.floor(Math.random() * 100000),
    review_count: Math.floor(Math.random() * 50),
  };

  // 如果是陪玩师，添加陪玩师信息
  if (user.role === UserRole.PLAYER) {
    const playerInfo: PlayerInfo = {
      id: id * 10,
      user_id: id,
      nickname: user.name,
      bio: '热爱游戏，擅长多种游戏类型，陪玩经验丰富，态度认真负责。',
      rating_average: 4.5 + Math.random() * 0.5,
      rating_count: Math.floor(Math.random() * 200) + 10,
      hourly_rate_cents: Math.floor(Math.random() * 5000) + 2000, // 20-70元/小时
      main_game_id: Math.floor(Math.random() * 5) + 1,
      verification_status:
        id % 3 === 0
          ? VerificationStatus.PENDING
          : id % 5 === 0
            ? VerificationStatus.REJECTED
            : VerificationStatus.VERIFIED,
      created_at: user.created_at,
      updated_at: user.updated_at,
    };
    detail.player = playerInfo;
  }

  return detail;
};
