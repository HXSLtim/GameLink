export interface HomeBannerItem {
  id: number
  type: 'link' | 'preview'
  /** 图片地址 */
  image: string
  /** 跳转链接（type=link 时有效） */
  link?: string
  /** 预览大图列表（type=preview 时有效） */
  previewImages?: string[]
  /** 标题（PC 端 Hero 区域展示） */
  title?: string
  /** 描述（PC 端 Hero 区域展示） */
  description?: string
  /** CTA 按钮文字 */
  actionText?: string
}
