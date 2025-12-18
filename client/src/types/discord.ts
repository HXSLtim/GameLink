// Discord Layout Type Definitions

/**
 * Server item in the server sidebar
 */
export interface ServerItem {
  id: string;
  name: string;
  icon?: string;
  iconColor?: string;
  unreadCount?: number;
  hasNotification?: boolean;
  isActive?: boolean;
}

/**
 * Channel category
 */
export interface ChannelCategory {
  id: string;
  name: string;
  collapsed?: boolean;
  channels: Channel[];
}

/**
 * Channel types
 */
export type ChannelType = 'text' | 'voice' | 'announcement' | 'stage';

/**
 * Channel item
 */
export interface Channel {
  id: string;
  name: string;
  type: ChannelType;
  unread?: boolean;
  mentionCount?: number;
  isNsfw?: boolean;
  isPrivate?: boolean;
  path?: string; // Route path for navigation
}

/**
 * Member status
 */
export type MemberStatus = 'online' | 'idle' | 'dnd' | 'offline';

/**
 * Member item in the member list
 */
export interface Member {
  id: string;
  username: string;
  nickname?: string;
  avatar?: string;
  status: MemberStatus;
  isBot?: boolean;
  roleColor?: string;
  activity?: string;
}

/**
 * Member group for categorizing members
 */
export interface MemberGroup {
  id: string;
  name: string;
  members: Member[];
}

/**
 * User panel info (bottom of channel sidebar)
 */
export interface UserPanelInfo {
  id: string;
  username: string;
  discriminator?: string;
  avatar?: string;
  status: MemberStatus;
  customStatus?: string;
}

/**
 * Server info header
 */
export interface ServerInfo {
  id: string;
  name: string;
  icon?: string;
  banner?: string;
  verified?: boolean;
  boostLevel?: number;
}
