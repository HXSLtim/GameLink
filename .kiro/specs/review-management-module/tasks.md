# 实施计划 - 评价管理模块

## 后端实现

- [x] 1. 扩展评价数据模型






- [x] 1.1 更新 Review 模型



  - 添加 status 字段（pending/approved/rejected/deleted）
  - 添加 isReported 字段
  - 添加 images 字段（JSON数组）
  - 添加 rejectionReason 字段
  - 更新数据库迁移
  - _需求: 1.1, 1.5, 2.1, 2.3, 8.1_

- [x] 1.2 编写评分范围约束属性测试







  - **属性 2: 评分范围约束**
  - **验证: 需求 1.1**

- [x] 2. 实现评价举报功能




- [x] 2.1 创建 ReviewReport 模型

  - 定义举报模型（reviewID, reporterID, reason, evidence, status, handledBy, handledAt）
  - 实现数据库迁移
  - _需求: 3.1, 3.2_

- [x] 2.2 实现 ReviewReport 仓储层


  - 实现 Create 方法
  - 实现 List 方法（支持状态筛选）
  - 实现 Get 方法
  - 实现 Update 方法（处理举报）
  - _需求: 3.1, 3.2_

- [x] 2.3 实现举报服务层


  - 实现 ReportReview 方法
  - 实现 ListReports 方法
  - 实现 HandleReport 方法（删除评价/警告/驳回）
  - _需求: 3.2, 3.3, 3.4, 3.5_

- [x] 2.4 实现举报 API 端点


  - POST /admin/reviews/:id/reports - 创建举报
  - GET /admin/review-reports - 列出举报
  - GET /admin/review-reports/:id - 获取举报详情
  - PUT /admin/review-reports/:id/handle - 处理举报
  - _需求: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 2.5 编写举报评价标记属性测试




  - **属性 5: 举报评价标记**
  - **验证: 需求 1.5, 3.1**


- [ ] 3. 实现评价审核功能





- [x] 3.1 扩展评价仓储层


  - 添加 ListPending 方法（获取待审核评价）
  - 添加 UpdateStatus 方法（更新审核状态）
  - 添加 BatchUpdateStatus 方法（批量审核）
  - _需求: 2.1, 2.2, 2.3, 2.5_

- [x] 3.2 实现审核服务层


  - 实现 ApproveReview 方法
  - 实现 RejectReview 方法（需要原因）
  - 实现 BatchApprove 方法
  - 实现 BatchReject 方法
  - 集成敏感词检测
  - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 3.3 实现审核 API 端点


  - GET /admin/reviews/pending - 获取待审核列表
  - PUT /admin/reviews/:id/approve - 批准评价
  - PUT /admin/reviews/:id/reject - 拒绝评价
  - PUT /admin/reviews/batch-approve - 批量批准
  - PUT /admin/reviews/batch-reject - 批量拒绝
  - _需求: 2.1, 2.2, 2.3, 2.5_

- [x] 3.4 编写审核状态转换属性测试






  - **属性 3: 审核状态转换合法性**
  - **验证: 需求 2.2, 2.3**

- [x] 3.5 编写批量操作原子性属性测试






  - **属性 10: 批量操作原子性**
  - **验证: 需求 2.5**

- [-] 4. 实现敏感词管理





- [-] 4.1 创建 SensitiveWord 模型

  - 定义敏感词模型（word, category, severity）
  - 实现数据库迁移
  - _需求: 5.1, 5.2_

- [ ] 4.2 实现 SensitiveWord 仓储层
  - 实现 Create 方法（带唯一性验证）
  - 实现 List 方法（支持搜索和分类筛选）
  - 实现 Update 方法
  - 实现 Delete 方法
  - 实现 GetAll 方法（用于检测）
  - _需求: 5.1, 5.2, 5.3, 5.4_

- [ ] 4.3 实现敏感词服务层
  - 实现 AddSensitiveWord 方法
  - 实现 UpdateSensitiveWord 方法
  - 实现 DeleteSensitiveWord 方法
  - 实现 DetectSensitiveWords 方法（检测并高亮）
  - 集成缓存机制
  - _需求: 5.2, 5.3, 5.4, 5.5_

- [ ] 4.4 实现敏感词 API 端点
  - GET /admin/sensitive-words - 列出敏感词
  - POST /admin/sensitive-words - 添加敏感词
  - PUT /admin/sensitive-words/:id - 更新敏感词
  - DELETE /admin/sensitive-words/:id - 删除敏感词
  - POST /admin/reviews/detect-sensitive - 检测敏感词
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ]* 4.5 编写敏感词唯一性属性测试
  - **属性 4: 敏感词唯一性**
  - **验证: 需求 5.3**

- [ ]* 4.6 编写敏感词检测准确性属性测试
  - **属性 7: 敏感词检测准确性**
  - **验证: 需求 2.4, 5.5**

- [ ] 5. 实现评价统计分析
- [ ] 5.1 实现统计服务层
  - 实现 GetReviewStats 方法（总数、平均分、分布）
  - 实现 GetReviewTrend 方法（30天趋势）
  - 实现 GetTopPlayers 方法（评价最多/评分最高）
  - 实现 GetGameStats 方法（按游戏统计）
  - _需求: 6.1, 6.2, 6.3, 6.4_

- [ ] 5.2 实现统计 API 端点
  - GET /admin/reviews/stats - 获取统计概览
  - GET /admin/reviews/trend - 获取趋势数据
  - GET /admin/reviews/top-players - 获取陪玩师排行
  - GET /admin/reviews/game-stats - 获取游戏统计
  - GET /admin/reviews/export - 导出统计数据
  - _需求: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ] 6. 实现评价展示设置
- [ ] 6.1 创建 ReviewDisplaySettings 模型
  - 定义设置模型（sortBy, minScore, showAnonymous, pageSize）
  - 实现配置存储（可使用配置表或配置文件）
  - _需求: 7.1, 7.2, 7.3, 7.4_

- [ ] 6.2 实现设置 API 端点
  - GET /admin/review-settings - 获取当前设置
  - PUT /admin/review-settings - 更新设置
  - _需求: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 7. 完善评价回复功能
- [ ] 7.1 扩展 ReviewReply 服务
  - 实现 UpdateReply 方法
  - 实现 DeleteReply 方法
  - 添加通知功能
  - _需求: 4.4, 4.5_

- [ ] 7.2 添加回复 API 端点
  - PUT /admin/review-replies/:id - 更新回复
  - DELETE /admin/review-replies/:id - 删除回复
  - _需求: 4.4, 4.5_

- [ ] 8. 实现操作日志记录
- [ ] 8.1 扩展 OperationLog 使用
  - 在评价审核时记录日志
  - 在举报处理时记录日志
  - 在评价删除时记录日志
  - 在回复操作时记录日志
  - _需求: 9.1_

- [ ] 8.2 实现日志查询 API
  - GET /admin/reviews/:id/logs - 获取评价操作日志
  - GET /admin/operation-logs - 搜索操作日志（支持筛选）
  - GET /admin/operation-logs/export - 导出日志
  - _需求: 9.2, 9.3, 9.4, 9.5_

- [ ]* 8.3 编写操作日志完整性属性测试
  - **属性 9: 操作日志完整性**
  - **验证: 需求 9.1**

- [ ] 9. 完善权限控制
- [ ] 9.1 添加评价管理权限
  - 定义 review:view 权限
  - 定义 review:approve 权限
  - 定义 review:delete 权限
  - 定义 review:manage 权限（敏感词管理）
  - 更新权限种子数据
  - _需求: 10.1, 10.2, 10.3_

- [ ] 9.2 应用权限中间件
  - 在所有评价管理端点应用权限验证
  - 记录未授权访问日志
  - 确保超级管理员可访问所有功能
  - _需求: 10.1, 10.2, 10.3, 10.4, 10.5_

- [ ]* 9.3 编写权限验证属性测试
  - **属性 8: 权限验证一致性**
  - **验证: 需求 10.1, 10.2, 10.3, 10.4**

- [ ] 10. 后端集成测试
- [ ]* 10.1 编写评价审核集成测试
  - 测试完整审核流程
  - 测试敏感词自动标记
  - 测试批量审核
  - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ]* 10.2 编写举报处理集成测试
  - 测试举报创建
  - 测试举报处理流程
  - 测试举报状态更新
  - _需求: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ]* 10.3 编写敏感词管理集成测试
  - 测试敏感词CRUD
  - 测试敏感词检测
  - 测试唯一性约束
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_

## 前端实现

- [ ] 11. 设置前端基础结构
- [ ] 11.1 创建评价管理模块目录
  - 创建 pages/admin/review 目录
  - 创建 components/review 目录
  - 创建 api/review.ts
  - 创建 types/review.ts
  - _需求: 所有需求_

- [ ] 11.2 定义 TypeScript 类型
  - 定义 Review 接口
  - 定义 ReviewReport 接口
  - 定义 SensitiveWord 接口
  - 定义 ReviewStats 接口
  - 定义 ReviewDisplaySettings 接口
  - _需求: 所有需求_

- [ ] 12. 实现评价列表页面
- [ ] 12.1 创建 ReviewList 页面组件
  - 实现评价记录表格布局
  - 添加搜索和筛选功能（订单ID、评价者、被评价者、评分、状态、时间）
  - 实现分页功能（每页20条）
  - 添加举报标记显示（红色标记）
  - 添加操作按钮（查看详情、审核、删除）
  - _需求: 1.1, 1.2, 1.3, 1.5_

- [ ]* 12.2 编写评价记录完整性属性测试
  - **属性 1: 评价记录完整性**
  - **验证: 需求 1.1**

- [ ] 13. 实现评价详情页面
- [ ] 13.1 创建 ReviewDetail 页面组件
  - 显示评价完整信息（评分、内容、时间）
  - 显示评价图片（缩略图和大图预览）
  - 显示订单信息
  - 显示操作历史（使用 Timeline 组件）
  - 显示回复列表
  - _需求: 1.4, 8.1, 8.2, 9.2_

- [ ] 14. 实现评价审核功能
- [ ] 14.1 创建 ReviewModeration 页面组件
  - 显示待审核评价列表
  - 实现批准评价功能
  - 实现拒绝评价功能（带原因输入）
  - 高亮显示敏感词
  - 实现批量选择和批量审核
  - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ] 14.2 创建 SensitiveWordHighlight 组件
  - 实现敏感词高亮显示
  - 支持不同严重程度的颜色标记
  - _需求: 2.4_

- [ ] 15. 实现举报管理功能
- [ ] 15.1 创建 ReviewReportList 页面组件
  - 显示所有被举报评价
  - 显示举报原因、举报人、举报时间
  - 实现举报状态筛选
  - _需求: 3.1_

- [ ] 15.2 创建 ReviewReportDetail 页面组件
  - 显示举报详情
  - 显示被举报评价完整内容
  - 实现删除评价功能
  - 实现警告评价者功能
  - 实现驳回举报功能
  - _需求: 3.2, 3.3, 3.4, 3.5_

- [ ] 16. 实现评价回复功能
- [ ] 16.1 创建 ReviewReply 组件
  - 实现回复输入框
  - 实现回复提交功能
  - 显示所有回复记录
  - 实现回复编辑和删除
  - _需求: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 17. 实现敏感词管理
- [ ] 17.1 创建 SensitiveWordList 页面组件
  - 显示所有敏感词
  - 显示敏感词分类和严重程度
  - 实现敏感词搜索
  - 添加操作按钮（编辑、删除）
  - _需求: 5.1_

- [ ] 17.2 创建 SensitiveWordForm 组件
  - 实现添加敏感词表单
  - 实现编辑敏感词表单
  - 验证敏感词唯一性
  - 支持分类和严重程度选择
  - _需求: 5.2, 5.3, 5.4_

- [ ] 18. 实现评价统计分析
- [ ] 18.1 创建 ReviewStats 页面组件
  - 显示统计概览卡片（总数、平均分）
  - 显示评分段分布（饼图或柱状图）
  - 显示评价趋势（折线图，30天）
  - _需求: 6.1, 6.2_

- [ ] 18.2 创建 PlayerReviewRanking 组件
  - 显示评价最多的陪玩师排行
  - 显示评分最高的陪玩师排行
  - _需求: 6.3_

- [ ] 18.3 创建 GameReviewStats 组件
  - 按游戏展示评价数量
  - 按游戏展示平均评分
  - _需求: 6.4_

- [ ] 18.4 实现统计数据导出
  - 添加导出按钮
  - 调用导出 API
  - 下载 Excel 文件
  - _需求: 6.5_

- [ ] 19. 实现评价展示设置
- [ ] 19.1 创建 ReviewDisplaySettings 页面组件
  - 配置评价排序规则（时间/评分/点赞数）
  - 配置评价过滤规则（最低评分、是否显示匿名）
  - 配置评价显示数量
  - 实现设置保存功能
  - _需求: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 20. 实现 API 客户端
- [ ] 20.1 创建评价 API 服务
  - 实现 getReviews 接口
  - 实现 getReviewDetail 接口
  - 实现 getPendingReviews 接口
  - 实现 approveReview 接口
  - 实现 rejectReview 接口
  - 实现 batchApproveReviews 接口
  - 实现 batchRejectReviews 接口
  - 实现 deleteReview 接口
  - _需求: 1.1, 1.2, 1.4, 2.1, 2.2, 2.3, 2.5_

- [ ] 20.2 创建举报 API 服务
  - 实现 getReviewReports 接口
  - 实现 getReviewReportDetail 接口
  - 实现 handleReviewReport 接口
  - _需求: 3.1, 3.2, 3.3, 3.4, 3.5_

- [ ] 20.3 创建敏感词 API 服务
  - 实现 getSensitiveWords 接口
  - 实现 addSensitiveWord 接口
  - 实现 updateSensitiveWord 接口
  - 实现 deleteSensitiveWord 接口
  - 实现 detectSensitiveWords 接口
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ] 20.4 创建统计 API 服务
  - 实现 getReviewStats 接口
  - 实现 getReviewTrend 接口
  - 实现 getTopPlayers 接口
  - 实现 getGameStats 接口
  - 实现 exportReviewStats 接口
  - _需求: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ] 20.5 创建回复 API 服务
  - 实现 replyReview 接口
  - 实现 updateReply 接口
  - 实现 deleteReply 接口
  - _需求: 4.1, 4.2, 4.4, 4.5_

- [ ] 20.6 创建设置 API 服务
  - 实现 getReviewSettings 接口
  - 实现 updateReviewSettings 接口
  - _需求: 7.1, 7.5_

- [ ]* 20.7 编写 API 客户端单元测试
  - 测试所有 API 接口
  - 测试错误处理
  - _需求: 所有需求_

- [ ] 21. 实现工具函数
- [ ] 21.1 创建格式化工具
  - 实现评分格式化（显示星级）
  - 实现日期时间格式化
  - 实现审核状态格式化（中文显示）
  - 实现举报状态格式化
  - _需求: 所有需求_

- [ ] 22. 配置路由和权限
- [ ] 22.1 添加评价管理路由
  - 添加评价列表路由
  - 添加评价详情路由
  - 添加评价审核路由
  - 添加举报管理路由
  - 添加敏感词管理路由
  - 添加统计分析路由
  - 添加展示设置路由
  - _需求: 所有需求_

- [ ] 22.2 配置权限守卫
  - 应用 review:view 权限
  - 应用 review:approve 权限
  - 应用 review:delete 权限
  - 应用 review:manage 权限
  - _需求: 10.1, 10.2, 10.3_

- [ ] 23. 前端集成测试
- [ ]* 23.1 编写评价管理 E2E 测试
  - 测试评价列表查询流程
  - 测试评价审核流程
  - 测试举报处理流程
  - 测试敏感词管理流程
  - _需求: 所有需求_

- [ ] 24. 最终检查点
  - 确保所有后端测试通过
  - 确保所有前端功能正常
  - 验证所有需求已实现
  - 检查代码质量和规范
  - 更新 API 文档
  - _需求: 所有需求_
