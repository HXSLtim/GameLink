/**
 * API 统一导出
 * 三层架构：interface（类型）-> request（HTTP）-> modules（业务）
 */

// 导出 API 客户端
export { apiClient as default } from './client';

// 导出类型定义层（interface）
export * from './interface/auth';

// 导出业务逻辑层（modules）
export * as authApi from './modules/auth';
export * as playerApi from './modules/player';

// 导出请求工具层（request）- 一般不直接使用，由 modules 内部调用
export { createRequest } from './request/http';
