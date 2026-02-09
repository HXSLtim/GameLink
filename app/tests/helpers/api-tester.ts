/**
 * API 测试辅助工具
 * 用于快速验证后端接口是否正常工作
 */

import { request } from '@/api/request'

// 测试结果类型
export interface TestResult {
  name: string
  passed: boolean
  error?: string
  duration: number
}

export class ApiTester {
  private results: TestResult[] = []
  private baseUrl: string

  constructor(baseUrl: string = 'http://localhost:8080/api/v1') {
    this.baseUrl = baseUrl
  }

  /**
   * 运行单个测试
   */
  async runTest(
    name: string,
    testFn: () => Promise<void>
  ): Promise<TestResult> {
    const startTime = Date.now()
    try {
      await testFn()
      const duration = Date.now() - startTime
      const result: TestResult = { name, passed: true, duration }
      this.results.push(result)
      console.log(`✅ ${name} (${duration}ms)`)
      return result
    } catch (error: any) {
      const duration = Date.now() - startTime
      const result: TestResult = {
        name,
        passed: false,
        error: error.message || String(error),
        duration,
      }
      this.results.push(result)
      console.error(`❌ ${name}: ${error.message || String(error)}`)
      return result
    }
  }

  /**
   * 打印测试结果汇总
   */
  printSummary() {
    const passed = this.results.filter((r) => r.passed).length
    const failed = this.results.length - passed
    const passRate = ((passed / this.results.length) * 100).toFixed(1)

    console.log('\n========== 测试结果汇总 ==========')
    console.log(`总计: ${this.results.length} 个测试`)
    console.log(`通过: ${passed} 个`)
    console.log(`失败: ${failed} 个`)
    console.log(`通过率: ${passRate}%`)

    if (failed > 0) {
      console.log('\n失败的测试:')
      this.results
        .filter((r) => !r.passed)
        .forEach((r) => {
          console.log(`  - ${r.name}: ${r.error}`)
        })
    }
    console.log('====================================\n')
  }

  /**
   * 获取测试结果
   */
  getResults(): TestResult[] {
    return [...this.results]
  }

  /**
   * 清空测试结果
   */
  clear() {
    this.results = []
  }
}

/**
 * 认证模块测试
 */
export async function testAuth(tester: ApiTester) {
  console.log('\n📝 测试认证模块...')

  await tester.runTest('发送验证码', async () => {
    await request({
      url: '/public/verification/send',
      method: 'POST',
      data: { phone: '13800138000' },
      showLoading: false,
      showError: false,
    })
  })

  // 注意：登录测试需要先注册或使用已存在的账号
  await tester.runTest('账号密码登录', async () => {
    await request({
      url: '/auth/login',
      method: 'POST',
      data: {
        phone: '13800138000',
        password: 'password123',
      },
      showLoading: false,
      showError: false,
    })
  })
}

/**
 * 陪玩师模块测试
 */
export async function testPlayers(tester: ApiTester) {
  console.log('\n🎮 测试陪玩师模块...')

  await tester.runTest('获取陪玩师列表', async () => {
    await request({
      url: '/public/players',
      method: 'GET',
      showLoading: false,
      showError: false,
    })
  })

  await tester.runTest('获取陪玩师详情', async () => {
    await request({
      url: '/public/players/1',
      method: 'GET',
      showLoading: false,
      showError: false,
    })
  })

  await tester.runTest('获取陪玩师评价', async () => {
    await request({
      url: '/public/players/1/reviews',
      method: 'GET',
      showLoading: false,
      showError: false,
    })
  })
}

/**
 * 订单模块测试 (需要登录)
 */
export async function testOrders(tester: ApiTester, token: string) {
  console.log('\n📦 测试订单模块...')

  await tester.runTest('获取订单列表', async () => {
    await request({
      url: '/orders',
      method: 'GET',
      header: { Authorization: `Bearer ${token}` },
      showLoading: false,
      showError: false,
    })
  })

  await tester.runTest('获取订单详情', async () => {
    await request({
      url: '/orders/1',
      method: 'GET',
      header: { Authorization: `Bearer ${token}` },
      showLoading: false,
      showError: false,
    })
  })
}

/**
 * 游戏模块测试
 */
export async function testGames(tester: ApiTester) {
  console.log('\n🎯 测试游戏模块...')

  await tester.runTest('获取游戏列表', async () => {
    await request({
      url: '/public/games',
      method: 'GET',
      showLoading: false,
      showError: false,
    })
  })

  await tester.runTest('获取游戏详情', async () => {
    await request({
      url: '/public/games/1',
      method: 'GET',
      showLoading: false,
      showError: false,
    })
  })
}

/**
 * 钱包模块测试 (需要登录)
 */
export async function testWallet(tester: ApiTester, token: string) {
  console.log('\n💰 测试钱包模块...')

  await tester.runTest('获取钱包余额', async () => {
    await request({
      url: '/wallet/balance',
      method: 'GET',
      header: { Authorization: `Bearer ${token}` },
      showLoading: false,
      showError: false,
    })
  })

  await tester.runTest('获取交易记录', async () => {
    await request({
      url: '/wallet/transactions',
      method: 'GET',
      header: { Authorization: `Bearer ${token}` },
      showLoading: false,
      showError: false,
    })
  })
}

/**
 * 运行所有公开接口测试 (不需要登录)
 */
export async function runPublicApiTests() {
  const tester = new ApiTester()

  console.log('🚀 开始运行公开 API 测试...\n')

  await testAuth(tester)
  await testPlayers(tester)
  await testGames(tester)

  tester.printSummary()

  return tester.getResults()
}

/**
 * 运行需要认证的测试
 */
export async function runAuthenticatedApiTests(token: string) {
  const tester = new ApiTester()

  console.log('🚀 开始运行认证 API 测试...\n')

  await testOrders(tester, token)
  await testWallet(tester, token)

  tester.printSummary()

  return tester.getResults()
}

/**
 * 运行完整测试套件
 */
export async function runAllTests(token?: string) {
  const tester = new ApiTester()

  console.log('🚀 开始运行完整 API 测试套件...\n')

  // 公开接口
  await testAuth(tester)
  await testPlayers(tester)
  await testGames(tester)

  // 认证接口
  if (token) {
    await testOrders(tester, token)
    await testWallet(tester, token)
  }

  tester.printSummary()

  return tester.getResults()
}

// 导出测试器类
export default ApiTester
