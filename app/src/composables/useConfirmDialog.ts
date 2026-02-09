import { reactive } from 'vue'
import type { ConfirmOptions } from '@/types/ui'

const dialogState = reactive({
  show: false,
  title: '提示',
  content: '',
  confirmText: '确定',
  cancelText: '取消',
})

let resolver: ((value: boolean) => void) | null = null

const resetResolver = (value: boolean) => {
  if (resolver) {
    resolver(value)
    resolver = null
  }
}

const closeDialog = (value = false) => {
  dialogState.show = false
  resetResolver(value)
}

export function useConfirmDialog() {
  const open = (options: ConfirmOptions = {}) => {
    dialogState.title = options.title || '提示'
    dialogState.content = options.content || ''
    dialogState.confirmText = options.confirmText || '确定'
    dialogState.cancelText = options.cancelText || '取消'
    dialogState.show = true
    return new Promise<boolean>((resolve) => {
      resolver = resolve
    })
  }

  const confirm = () => closeDialog(true)
  const cancel = () => closeDialog(false)
  const close = () => closeDialog(false)

  return {
    dialogState,
    open,
    confirm,
    cancel,
    close,
  }
}

export function confirmDialog(options: ConfirmOptions) {
  return useConfirmDialog().open(options)
}
