# 快速提升到80%的执行计划

## 当前状态
- 覆盖率: 35.5%
- 目标: 80.0%
- 差距: 44.5%
- 策略: 批量快速添加测试

## 执行策略

### 阶段1: 补充0%函数的简单测试 (预计+5%)
- auth包的0%函数
- cache/redis的0%函数
- config的低覆盖函数

### 阶段2: Service层核心方法 (预计+15%)
- Admin Service: 40% → 75%
- Role Service: 55% → 80%
- Player Service: 66% → 80%
- Order Service: 67% → 80%

### 阶段3: Repository层补充 (预计+8%)
- Commission: 77% → 90%
- ServiceItem: 78% → 90%
- Permission: 63% → 85%

### 阶段4: Handler层关键接口 (预计+12%)
- User Handler: 39% → 70%
- Player Handler: 39% → 70%
- Admin Handler: 0% → 50%

### 阶段5: 小模块批量提升 (预计+5%)
- db: 30% → 60%
- logging: 29% → 60%
- metrics: 19% → 50%

预计总提升: 45% → 达到80%+

