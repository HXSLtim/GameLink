export type AgreementType = 'user' | 'privacy' | 'player' | 'recharge'

export interface AgreementSection {
  title: string
  text: string
}

export interface AgreementData {
  title: string
  updateTime: string
  sections: AgreementSection[]
}
