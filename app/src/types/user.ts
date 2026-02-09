export type UserRole = 'user' | 'player' | 'admin'

export type AppUserRole = Exclude<UserRole, 'admin'>
