import { useUserStore } from '@/store'
import { getStorage, removeStorage, setStorage } from '@/utils/storage'

const REDIRECT_KEY = 'post_login_redirect'
let guardInitialized = false

const TAB_PAGES = new Set<string>([
  '/pages/index/index',
  '/pages/player/list/index',
  '/pages/message/list/index',
  '/pages/profile/index/index',
])

const GUEST_ONLY_PAGES = new Set<string>([
  '/pages/auth/login/index',
  '/pages/auth/register/index',
])

const AUTH_ONLY_PAGES = new Set<string>([
  '/pages/profile/index/index',
  '/pages/profile/edit/index',
  '/pages/settings/index/index',
  '/pages/favorite/list/index',
  '/pages/review/list/index',
  '/pages/player/certification/index',
  '/pages/payment/result/index',
])

const PLAYER_ONLY_PAGES = new Set<string>([
  '/pages/player/dashboard/index',
  '/pages/player/orders/index',
  '/pages/player/earnings/index',
  '/pages/player/services/index',
  '/pages/player/schedule/index',
])

const AUTH_PREFIXES = [
  '/pages/order/',
  '/pages/message/',
  '/pages/wallet/',
]

const CHAT_PAGE = '/pages/message/chat/index'

function normalizePath(url: string): string {
  return url.split('?')[0]
}

function isAuthRequired(path: string): boolean {
  if (GUEST_ONLY_PAGES.has(path)) return false
  if (PLAYER_ONLY_PAGES.has(path)) return true
  if (AUTH_ONLY_PAGES.has(path)) return true
  return AUTH_PREFIXES.some(prefix => path.startsWith(prefix))
}

function isPlayerOnly(path: string): boolean {
  return PLAYER_ONLY_PAGES.has(path)
}

function buildLoginUrl(redirectUrl: string): string {
  const encoded = encodeURIComponent(redirectUrl)
  return `/pages/auth/login/index?redirect=${encoded}`
}

function safeRedirectToLogin(redirectUrl: string) {
  if (!redirectUrl) {
    uni.navigateTo({ url: '/pages/auth/login/index' })
    return
  }
  uni.navigateTo({ url: buildLoginUrl(redirectUrl) })
}

function parseQuery(url: string): Record<string, string> {
  const queryIndex = url.indexOf('?')
  if (queryIndex < 0) return {}
  const query = url.slice(queryIndex + 1)
  if (!query) return {}
  return query.split('&').reduce<Record<string, string>>((acc, item) => {
    const [rawKey, rawValue] = item.split('=')
    if (!rawKey) return acc
    acc[decodeURIComponent(rawKey)] = rawValue ? decodeURIComponent(rawValue) : ''
    return acc
  }, {})
}

export function isTabPage(url: string): boolean {
  return TAB_PAGES.has(normalizePath(url))
}

export function setRedirectPath(url: string): void {
  if (!url) return
  setStorage(REDIRECT_KEY, url)
}

export function consumeRedirectPath(): string | null {
  const value = getStorage<string>(REDIRECT_KEY)
  if (!value) return null
  removeStorage(REDIRECT_KEY)
  return value
}

export function redirectToUrl(url: string): void {
  const path = normalizePath(url)
  if (TAB_PAGES.has(path)) {
    uni.switchTab({ url: path })
    return
  }
  uni.redirectTo({ url })
}

export function setupRouteGuard(): void {
  if (guardInitialized) return
  guardInitialized = true

  const guard = (url?: string) => {
    if (!url) return true
    const path = normalizePath(url)
    const userStore = useUserStore()
    const query = parseQuery(url)

    if (GUEST_ONLY_PAGES.has(path) && userStore.isLoggedIn) {
      redirectToUrl('/pages/index/index')
      return false
    }

    if (path === CHAT_PAGE && query.groupType === 'public' && !userStore.isLoggedIn) {
      setRedirectPath(url)
      safeRedirectToLogin(url)
      return false
    }

    if (isAuthRequired(path) && !userStore.isLoggedIn) {
      setRedirectPath(url)
      safeRedirectToLogin(url)
      return false
    }

    if (isPlayerOnly(path) && userStore.isLoggedIn && !userStore.isPlayer) {
      uni.showToast({ title: '请先完成陪玩师认证', icon: 'none' })
      redirectToUrl('/pages/player/certification/index')
      return false
    }

    return true
  }

  ;['navigateTo', 'redirectTo', 'reLaunch', 'switchTab'].forEach((method) => {
    uni.addInterceptor(method as 'navigateTo', {
      invoke(args: { url: string }) {
        return guard(args.url)
      },
    })
  })
}
