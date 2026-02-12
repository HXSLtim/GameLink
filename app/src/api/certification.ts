/**
 * 陪玩认证相关 API
 */

import { get, post, put, ApiError } from './request'
import { uploadFile } from './request'
import type { Gender } from '@/types/common'
import type { CertStatus } from '@/types/status'

// 认证状态
export type CertificationStatus = CertStatus

// 游戏段位认证
export interface GameCertification {
  gameId: number
  gameName: string
  rankId?: number
  rankName?: string
  screenshotUrl: string
}

// 认证信息
export interface CertificationInfo {
  id?: number
  status: CertificationStatus
  realName?: string
  idNumber?: string
  gender?: Gender
  idCardFront?: string
  idCardBack?: string
  games?: GameCertification[]
  introduction?: string
  voiceSample?: string
  rejectedReason?: string
  createdAt?: string
  reviewedAt?: string
}

// 提交认证参数
export interface SubmitCertificationParams {
  realName: string
  idNumber: string
  gender: Gender
  idCardFront: string
  idCardBack: string
  games: GameCertification[]
  introduction: string
  voiceSample?: string
}

/**
 * 获取实名认证状态
 */
export function getCertificationStatus() {
  return get<CertificationInfo>('/player/certification/identity', undefined, { showError: false })
    .catch((error: unknown) => {
      if (error instanceof ApiError && error.code === 401 && error.message.includes('missing player')) {
        return {
          success: true,
          code: 200,
          message: 'OK',
          data: { status: 'none' } as CertificationInfo,
        }
      }
      throw error
    })
}

/**
 * 提交实名认证申请
 */
export function submitCertification(data: SubmitCertificationParams) {
  return post<CertificationInfo>('/player/certification/identity', data)
}

/**
 * 更新认证信息
 */
export function updateCertification(data: Partial<SubmitCertificationParams>) {
  return put<CertificationInfo>('/player/certification/identity', data)
}

/**
 * 上传身份证照片
 */
export function uploadIdCardImage(filePath: string, type: 'front' | 'back') {
  return uploadFile('/certification/upload/image', filePath, 'image', { type: 'id-card', side: type })
}

/**
 * 上传段位截图
 */
export function uploadRankScreenshot(filePath: string) {
  return uploadFile('/certification/upload/image', filePath, 'image', { type: 'skill-proof' })
}

/**
 * 上传语音样本
 */
export function uploadVoiceSample(filePath: string) {
  return uploadFile('/certification/upload/voice', filePath, 'voice')
}

export default {
  getCertificationStatus,
  submitCertification,
  updateCertification,
  uploadIdCardImage,
  uploadRankScreenshot,
  uploadVoiceSample,
}
