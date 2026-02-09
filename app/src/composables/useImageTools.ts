import { previewImage } from '@/utils'
import type { ImagePickOptions } from '@/types/media'

const defaultPickOptions: Required<ImagePickOptions> = {
  count: 1,
  sizeType: ['compressed'],
  sourceType: ['album', 'camera'],
}

export function useImageTools() {
  const pickImages = (options: ImagePickOptions = {}) => {
    return new Promise<string[]>((resolve, reject) => {
      uni.chooseImage({
        ...defaultPickOptions,
        ...options,
        success: (res) => {
          resolve(res.tempFilePaths || [])
        },
        fail: (error) => {
          reject(error)
        },
      })
    })
  }

  const previewImages = (urls: string[], current?: number | string) => {
    previewImage(urls, current ?? 0)
  }

  return {
    pickImages,
    previewImages,
  }
}
