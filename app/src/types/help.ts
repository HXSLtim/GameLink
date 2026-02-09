import type { FaqItemData } from '@/types/faq'

export interface HelpCategory {
  id: string
  name: string
  icon: string
}

export interface HelpFaq extends FaqItemData {
  categoryId: string
}
