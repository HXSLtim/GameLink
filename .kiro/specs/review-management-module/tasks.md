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


- [x] 3. 实现评价审核功能













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

- [x] 4. 实现敏感词管理






- [x] 4.1 创建 SensitiveWord 模型


  - 定义敏感词模型（word, category, severity）
  - 实现数据库迁移
  - _需求: 5.1, 5.2_

- [x] 4.2 实现 SensitiveWord 仓储层


  - 实现 Create 方法（带唯一性验证）
  - 实现 List 方法（支持搜索和分类筛选）
  - 实现 Update 方法
  - 实现 Delete 方法
  - 实现 GetAll 方法（用于检测）
  - _需求: 5.1, 5.2, 5.3, 5.4_


- [x] 4.3 实现敏感词服务层

  - 实现 AddSensitiveWord 方法
  - 实现 UpdateSensitiveWord 方法
  - 实现 DeleteSensitiveWord 方法
  - 实现 DetectSensitiveWords 方法（检测并高亮）
  - 集成缓存机制
  - _需求: 5.2, 5.3, 5.4, 5.5_

- [x] 4.4 实现敏感词 API 端点



  - GET /admin/sensitive-words - 列出敏感词
  - POST /admin/sensitive-words - 添加敏感词
  - PUT /admin/sensitive-words/:id - 更新敏感词
  - DELETE /admin/sensitive-words/:id - 删除敏感词
  - POST /admin/reviews/detect-sensitive - 检测敏感词
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 4.5 编写敏感词唯一性属性测试






  - **属性 4: 敏感词唯一性**
  - **验证: 需求 5.3**


- [x] 4.6 编写敏感词检测准确性属性测试





  - **属性 7: 敏感词检测准确性**
  - **验证: 需求 2.4, 5.5**

- [x] 5. 实现评价统计分析



- [x] 5.1 实现统计服务层





  - 实现 GetReviewStats 方法（总数、平均分、分布）
  - 实现 GetReviewTrend 方法（30天趋势）
  - 实现 GetTopPlayers 方法（评价最多/评分最高）
  - 实现 GetGameStats 方法（按游戏统计）
  - _需求: 6.1, 6.2, 6.3, 6.4_

- [x] 5.2 实现统计 API 端点





  - GET /admin/reviews/stats - 获取统计概览
  - GET /admin/reviews/trend - 获取趋势数据
  - GET /admin/reviews/top-players - 获取陪玩师排行
  - GET /admin/reviews/game-stats - 获取游戏统计
  - GET /admin/reviews/export - 导出统计数据
  - _需求: 6.1, 6.2, 6.3, 6.4, 6.5_

- [x] 6. 实现评价展示设置





- [x] 6.1 创建 ReviewDisplaySettings 模型


  - 定义设置模型（sortBy, minScore, showAnonymous, pageSize）
  - 实现配置存储（可使用配置表或配置文件）
  - _需求: 7.1, 7.2, 7.3, 7.4_

- [x] 6.2 实现设置 API 端点


  - GET /admin/review-settings - 获取当前设置
  - PUT /admin/review-settings - 更新设置
  - _需求: 7.1, 7.2, 7.3, 7.4, 7.5_

- [x] 7. 完善评价回复功能





- [x] 7.1 扩展 ReviewReply 服务


  - 实现 UpdateReply 方法
  - 实现 DeleteReply 方法
  - 添加通知功能
  - _需求: 4.4, 4.5_

- [x] 7.2 添加回复 API 端点


  - PUT /admin/review-replies/:id - 更新回复
  - DELETE /admin/review-replies/:id - 删除回复
  - _需求: 4.4, 4.5_

- [x] 8. 实现操作日志记录




- [x] 8.1 扩展 OperationLog 使用


  - 在评价审核时记录日志
  - 在举报处理时记录日志
  - 在评价删除时记录日志
  - 在回复操作时记录日志
  - _需求: 9.1_

- [x] 8.2 实现日志查询 API


  - GET /admin/reviews/:id/logs - 获取评价操作日志
  - GET /admin/operation-logs - 搜索操作日志（支持筛选）
  - GET /admin/operation-logs/export - 导出日志
  - _需求: 9.2, 9.3, 9.4, 9.5_




- [x] 8.3 编写操作日志完整性属性测试




  - **属性 9: 操作日志完整性**
  - **验证: 需求 9.1**



- [x] 9. 完善权限控制








- [x] 9.1 添加评价管理权限


  - 定义 review:view 权限
  - 定义 review:approve 权限
  - 定义 review:delete 权限
  - 定义 review:manage 权限（敏感词管理）


  - 更新权限种子数据
  - _需求: 10.1, 10.2, 10.3_

- [x] 9.2 应用权限中间件


  - 在所有评价管理端点应用权限验证
  - 记录未授权访问日志
  - 确保超级管理员可访问所有功能
  - _需求: 10.1, 10.2, 10.3, 10.4, 10.5_

- [x] 9.3 编写权限验证属性测试






  - **属性 8: 权限验证一致性**
  - **验证: 需求 10.1, 10.2, 10.3, 10.4**


- [x] 10. 后端集成测试



- [x] 10.1 编写评价审核集成测试






  - 测试完整审核流程
  - 测试敏感词自动标记
  - 测试批量审核
  - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_

- [x] 10.2 编写举报处理集成测试
  - 测试举报创建
  - 测试举报处理流程
  - 测试举报状态更新
  - _需求: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 10.3 编写敏感词管理集成测试
  - 测试敏感词CRUD
  - 测试敏感词检测
  - 测试唯一性约束
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_





## 前端实现


- [ ] 11. 设置前端基础结构

- [x] 11.1 创建评价管理模块目录



  - 创建 pages/admin/Review 目录
  - 创建 components/review 目录（如需要）
  - 创建 api/review.ts
  - 创建 types/review.ts
  - _需求: 所有需求_

- [x] 11.2 定义 TypeScript 类型

  - 定义 Review 接口
  - 定义 ReviewReport 接口
  - 定义 SensitiveWord 接口
  - 定义 ReviewStats 接口
  - 定义 ReviewDisplaySettings 接口
  - 定义 API 请求和响应类型
  - _需求: 所有需求_




- [ ] 12. 实现评价列表页面

- [x] 12.1 创建 ReviewList 页面组件


  - 实现评价记录表格布局（使用 Ant Design Table）
  - 添加搜索和筛选功能（订单ID、评价者、被评价者、评分、状态、时间）
  - 实现分页功能（每页20条）
  - 添加举报标记显示（红色标记或 Badge）
  - 添加操作按钮（查看详情、审核、删除）
  - 集成后端 API 调用
  - _需求: 1.1, 1.2, 1.3, 1.5_

- [x] 13. 实现评价详情页面







- [x] 13.1 创建 ReviewDetail 页面组件


  - 显示评价完整信息（评分、内容、时间）
  - 显示评价图片（缩略图和大图预览，使用 Image 组件）
  - 显示订单信息（订单ID、用户、陪玩师）
  - 显示操作历史（使用 Timeline 组件）
  - 显示回复列表
  - 添加回复功能入口

  - _需求: 1.4, 8.1, 8.2, 9.2_

- [x] 14. 实现评价审核功能



- [x] 14.1 创建 ReviewModeration 页面组件

  - 显示待审核评价列表（调用 /admin/reviews/pending API）
  - 实现批准评价功能（单个批准）
  - 实现拒绝评价功能（带原因输入的 Modal）
  - 高亮显示敏感词（集成敏感词检测 API）
  - 实现批量选择和批量审核（使用 Table rowSelection）
  - 添加加载状态和错误处理
  - _需求: 2.1, 2.2, 2.3, 2.4, 2.5_



- [x] 14.2 创建 SensitiveWordHighlight 组件


  - 实现敏感词高亮显示（使用不同背景色）

  - 支持不同严重程度的颜色标记（低/中/高）
  - 可复用组件，接收内容和检测结果作为 props

  - _需求: 2.4_

- [ ] 15. 实现举报管理功能

- [x] 15.1 创建 ReviewReportList 页面组件


  - 显示所有被举报评价（Table 组件）
  - 显示举报原因、举报人、举报时间
  - 实现举报状态筛选（pending/approved/rejected）
  - 实现分页功能

  - 添加查看详情按钮
  - _需求: 3.1_




- [ ] 15.2 创建 ReviewReportDetail 页面或 Modal 组件
  - 显示举报详情（举报原因、证据、举报人信息）
  - 显示被举报评价完整内容

  - 实现删除评价功能（带确认对话框）

  - 实现警告评价者功能（带备注输入）
  - 实现驳回举报功能（带原因输入）
  - 处理成功后更新列表状态
  - _需求: 3.2, 3.3, 3.4, 3.5_

- [ ] 16. 实现评价回复功能

- [x] 16.1 创建 ReviewReply 组件


  - 实现回复输入框（TextArea）

  - 实现回复提交功能（调用 API）
  - 显示所有回复记录（List 或 Comment 组件）
  - 实现回复编辑和删除（带权限控制）
  - 添加加载和提交状态
  - _需求: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 17. 实现敏感词管理

- [x] 17.1 创建 SensitiveWordList 页面组件


  - 显示所有敏感词（Table 组件）
  - 显示敏感词分类和严重程度（使用 Tag 组件）

  - 实现敏感词搜索（按词内容搜索）


  - 实现分类筛选
  - 添加操作按钮（编辑、删除）
  - 添加新增敏感词按钮
  - _需求: 5.1_



- [ ] 17.2 创建 SensitiveWordForm 组件
  - 实现添加敏感词表单（Modal 或独立页面）
  - 实现编辑敏感词表单
  - 验证敏感词唯一性（前端和后端验证）
  - 支持分类选择（Select 组件）
  - 支持严重程度选择（Radio 或 Select）
  - 表单验证和错误提示
  - _需求: 5.2, 5.3, 5.4_

- [ ] 18. 实现评价统计分析

- [x] 18.1 创建 ReviewStats 页面组件





  - 显示统计概览卡片（总数、平均分，使用 Statistic 组件）
  - 显示评分段分布（使用 Ant Design Charts 或 ECharts）
  - 显示评价趋势（折线图，30天）
  - 实现时间范围选择（可选）
  - _需求: 6.1, 6.2_


- [ ] 18.2 创建 PlayerReviewRanking 组件
  - 显示评价最多的陪玩师排行（Table 或 List）
  - 显示评分最高的陪玩师排行
  - 支持切换排行类型
  - 显示排名、陪玩师信息、评价数/评分

  - _需求: 6.3_


- [x] 18.3 创建 GameReviewStats 组件

  - 按游戏展示评价数量（Table 或卡片）
  - 按游戏展示平均评分

  - 可视化展示（柱状图或饼图）
  - _需求: 6.4_


- [ ] 18.4 实现统计数据导出


  - 添加导出按钮（在统计页面）
  - 调用导出 API（/admin/reviews/export）
  - 处理文件下载（CSV 格式）
  - 添加导出类型选择（概览/趋势/陪玩师/游戏）
  - _需求: 6.5_


- [ ] 19. 实现评价展示设置

- [x] 19.1 创建 ReviewDisplaySettings 页面组件


  - 配置评价排序规则（时间/评分/点赞数，使用 Radio 或 Select）
  - 配置评价过滤规则（最低评分、是否显示匿名，使用 InputNumber 和 Switch）
  - 配置评价显示数量（每页条数）
  - 实现设置保存功能（调用 PUT API）
  - 显示当前设置（调用 GET API）
  - 添加重置为默认值功能
  - _需求: 7.1, 7.2, 7.3, 7.4, 7.5_

- [ ] 20. 实现 API 客户端

- [x] 20.1 创建评价 API 服务（api/review.ts）


  - 实现 getReviews 接口（GET /admin/reviews）
  - 实现 getReviewDetail 接口（GET /admin/reviews/:id）
  - 实现 getPendingReviews 接口（GET /admin/reviews/pending）
  - 实现 approveReview 接口（PUT /admin/reviews/:id/approve）
  - 实现 rejectReview 接口（PUT /admin/reviews/:id/reject）
  - 实现 batchApproveReviews 接口（PUT /admin/reviews/batch-approve）

  - 实现 batchRejectReviews 接口（PUT /admin/reviews/batch-reject）
  - 实现 deleteReview 接口（DELETE /admin/reviews/:id）
  - 实现 getReviewLogs 接口（GET /admin/reviews/:id/logs）
  - 使用统一的 axios 实例和错误处理
  - _需求: 1.1, 1.2, 1.4, 2.1, 2.2, 2.3, 2.5, 9.2_


- [x] 20.2 创建举报 API 服务（api/reviewReport.ts）

  - 实现 getReviewReports 接口（GET /admin/review-reports）
  - 实现 getReviewReportDetail 接口（GET /admin/review-reports/:id）
  - 实现 handleReviewReport 接口（PUT /admin/review-reports/:id/handle）
  - 实现 createReviewReport 接口（POST /admin/reviews/:id/reports）
  - _需求: 3.1, 3.2, 3.3, 3.4, 3.5_


- [x] 20.3 创建敏感词 API 服务（api/sensitiveWord.ts）

  - 实现 getSensitiveWords 接口（GET /admin/sensitive-words）
  - 实现 addSensitiveWord 接口（POST /admin/sensitive-words）
  - 实现 updateSensitiveWord 接口（PUT /admin/sensitive-words/:id）
  - 实现 deleteSensitiveWord 接口（DELETE /admin/sensitive-words/:id）
  - 实现 detectSensitiveWords 接口（POST /admin/reviews/detect-sensitive）
  - _需求: 5.1, 5.2, 5.3, 5.4, 5.5_



- [x] 20.4 创建统计 API 服务（api/reviewStats.ts）

  - 实现 getReviewStats 接口（GET /admin/reviews/stats）
  - 实现 getReviewTrend 接口（GET /admin/reviews/trend）
  - 实现 getTopPlayers 接口（GET /admin/reviews/top-players）
  - 实现 getGameStats 接口（GET /admin/reviews/game-stats）
  - 实现 exportReviewStats 接口（GET /admin/reviews/export）
  - _需求: 6.1, 6.2, 6.3, 6.4, 6.5_


- [ ] 20.5 创建回复 API 服务（api/reviewReply.ts）
  - 实现 createReply 接口（POST /admin/reviews/:id/reply）
  - 实现 updateReply 接口（PUT /admin/review-replies/:id）
  - 实现 deleteReply 接口（DELETE /admin/review-replies/:id）
  - _需求: 4.1, 4.2, 4.4, 4.5_


- [ ] 20.6 创建设置 API 服务（api/reviewSettings.ts）
  - 实现 getReviewSettings 接口（GET /admin/review-settings）


  - 实现 updateReviewSettings 接口（PUT /admin/review-settings）
  - _需求: 7.1, 7.5_

- [ ] 21. 实现工具函数

- [x] 21.1 创建格式化工具（utils/reviewFormat.ts）


  - 实现评分格式化（显示星级，可使用 Rate 组件或自定义）
  - 实现日期时间格式化（使用 dayjs）

  - 实现审核状态格式化（中文显示和颜色映射）
  - 实现举报状态格式化
  - 实现敏感词分类和严重程度格式化
  - _需求: 所有需求_

- [x] 22. 配置路由和权限






- [ ] 22.1 添加评价管理路由（router/index.tsx）
  - 添加评价列表路由（/admin/reviews）
  - 添加评价详情路由（/admin/reviews/:id）

  - 添加评价审核路由（/admin/reviews/moderation）
  - 添加举报管理路由（/admin/review-reports）


  - 添加敏感词管理路由（/admin/sensitive-words）
  - 添加统计分析路由（/admin/reviews/stats）
  - 添加展示设置路由（/admin/review-settings）
  - 配置路由懒加载

  - _需求: 所有需求_

- [x] 22.2 配置权限守卫

  - 应用 review:view 权限（评价列表、详情）
  - 应用 review:approve 权限（审核功能）
  - 应用 review:delete 权限（删除功能）
  - 应用 review:manage 权限（敏感词管理）
  - 在路由配置中添加权限元数据


  - 在组件中使用权限控制显示/隐藏功能


  - _需求: 10.1, 10.2, 10.3_

- [-] 22.3 更新侧边栏菜单

  - 在管理后台侧边栏添加评价管理菜单项
  - 添加子菜单（评价列表、待审核、举报管理、敏感词、统计、设置）
  - 配置菜单图标和权限
  - _需求: 所有需求_

- [ ] 23. 测试和验证

- [ ] 23.1 功能测试
  - 测试评价列表查询和筛选
  - 测试评价审核流程（批准、拒绝、批量）
  - 测试举报处理流程（删除、警告、驳回）
  - 测试敏感词管理（增删改查、检测）
  - 测试统计数据展示和导出
  - 测试回复功能
  - 测试展示设置
  - _需求: 所有需求_

- [ ] 23.2 权限测试
  - 测试不同角色的访问权限
  - 测试无权限用户的访问限制
  - 测试操作权限控制
  - _需求: 10.1, 10.2, 10.3_

- [ ] 23.3 UI/UX 测试
  - 测试响应式布局
  - 测试加载状态和错误提示
  - 测试表单验证
  - 测试用户交互流畅性
  - _需求: 所有需求_

- [ ] 24. 最终检查点
  - 确保所有后端测试通过
  - 确保所有前端功能正常
  - 验证所有需求已实现
  - 检查代码质量和规范（ESLint、Prettier）
  - 更新 API 文档（如有变更）
  - 编写用户使用文档（可选）
  - _需求: 所有需求_
