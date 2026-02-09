/**
 * Banner 相关 API
 */

import { get } from './request'

// Banner 数据
export interface BannerItem {
  id: number
  title?: string
  description?: string
  imageUrl: string
  type: 'link' | 'preview'
  link?: string
  actionText?: string
}

// Banner 列表响应
export interface BannerListResponse {
  banners: BannerItem[]
}

/**
 * 获取首页 banner 列表（公开接口，无需登录）
 */
export function getBanners() {
  return get<BannerListResponse>('/public/banners')
}
