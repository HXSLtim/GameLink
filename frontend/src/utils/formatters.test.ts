import { describe, it, expect, vi } from 'vitest'
import { formatOrderStatus, getOrderStatusColor, formatCurrency, formatDuration, formatDateTime, formatRelativeTime, formatPrice } from './formatters'
import { OrderStatus } from '../types/order'

describe('formatters', () => {
  it('formats order status text', () => {
    expect(formatOrderStatus(OrderStatus.PENDING)).toBe('待处理')
    expect(formatOrderStatus(OrderStatus.COMPLETED)).toBe('已完成')
  })

  it('returns tag color by order status', () => {
    expect(getOrderStatusColor(OrderStatus.PENDING)).toBe('warning')
    expect(getOrderStatusColor(OrderStatus.REFUNDED)).toBe('error')
  })

  it('formats currency cents to yuan', () => {
    expect(formatCurrency(0)).toBe('¥0.00')
    expect(formatCurrency(12345)).toBe('¥123.45')
    expect(formatPrice(999)).toBe('¥9.99')
  })

  it('formats duration hours', () => {
    expect(formatDuration(3)).toBe('3小时')
  })

  it('formats date time string', () => {
    const d = new Date('2024-01-02T03:04:00Z')
    const local = new Date(d.getTime() - d.getTimezoneOffset() * 60000)
    const str = local.toISOString()
    expect(formatDateTime(str)).toMatch(/\d{4}-\d{2}-\d{2} \d{2}:\d{2}/)
    expect(formatDateTime()).toBe('-')
  })

  it('formats relative time for recent dates', () => {
    const now = Date.now()
    vi.setSystemTime(new Date(now))
    expect(formatRelativeTime(new Date(now - 30 * 1000).toISOString())).toBe('刚刚')
    expect(formatRelativeTime(new Date(now - 5 * 60 * 1000).toISOString())).toBe('5分钟前')
    expect(formatRelativeTime(new Date(now - 2 * 60 * 60 * 1000).toISOString())).toBe('2小时前')
    expect(formatRelativeTime(new Date(now - 3 * 24 * 60 * 60 * 1000).toISOString())).toBe('3天前')
    vi.useRealTimers()
  })
})
