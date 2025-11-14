import { describe, it, expect } from 'vitest'
import { mergeUrlParams } from './urlParams'

describe('mergeUrlParams', () => {
  it('merges numeric values from URL into initial params', () => {
    const search = new URLSearchParams({ page: '2', pageSize: '50' })
    const initial = { page: 1, pageSize: 20, q: '' }
    const result = mergeUrlParams(search, initial, ['page', 'pageSize', 'q'])
    expect(result).toEqual({ page: 2, pageSize: 50, q: '' })
  })

  it('merges string values and ignores empty values', () => {
    const search = new URLSearchParams({ q: 'abc', page: '' })
    const initial = { page: 1, q: '' }
    const result = mergeUrlParams(search, initial, ['page', 'q'])
    expect(result).toEqual({ page: 1, q: 'abc' })
  })

  it('ignores missing keys and keeps initial params', () => {
    const search = new URLSearchParams({})
    const initial = { page: 1, pageSize: 20 }
    const result = mergeUrlParams(search, initial, ['page', 'pageSize'])
    expect(result).toEqual({ page: 1, pageSize: 20 })
  })
})
