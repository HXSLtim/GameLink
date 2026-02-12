# GameLink 第三方接入验收模板

> 适用范围：OSS、微信支付、支付宝、短信、TRTC、微信登录等外部依赖接入后的联调与上线前验收。

## 1. 基础信息

| 字段 | 内容 |
| --- | --- |
| 验收日期 |  |
| 验收环境 | `local` / `staging` / `production` |
| 验收分支/提交 |  |
| 验收负责人 |  |
| 参与角色 | 后端 / 前端 / 测试 / 运维 |

## 2. 外部服务配置清单（脱敏）

| 服务 | 必填配置 | 已配置 | 备注 |
| --- | --- | --- | --- |
| OSS | `OSS_PROVIDER/ENDPOINT/ACCESS_KEY/SECRET_KEY/BUCKET/REGION/ENABLED` |  |  |
| 微信支付 | `WECHAT_PAY_*` |  |  |
| 支付宝 | `ALIPAY_*` |  |  |
| 短信 | `SMS_PROVIDER/ACCESS_KEY/SECRET_KEY/SIGN_NAME/ENABLED` |  |  |
| TRTC | `TRTC_*` |  |  |
| 微信登录 | `WECHAT_APP_ID/WECHAT_SECRET` |  |  |

## 3. OSS 验收

### 3.1 配置核对
- `OSS_ENABLED=true`
- `OSS_PROVIDER=qcloud`（若走腾讯 COS）
- Endpoint/Region/Bucket 与控制台一致

### 3.2 自动化验证
- 执行命令：
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_oss_smoke.ps1 -BaseUrl "http://127.0.0.1:8080/api/v1"`
- 通过标准：
  - 返回 `PASS`
  - `upload_url` 非本地路径（不是 `/uploads/...`）
  - URL Host 与预期 OSS 域名一致

### 3.3 手工验证
- 前端上传头像/聊天图片/认证图片
- 页面回显正常（含签名 URL 过期前可访问）
- 记录截图与请求日志

## 4. 支付验收（微信/支付宝）

### 4.1 正向流程
- 下单 -> 拉起支付 -> 回调成功 -> 订单状态流转正确
- 钱包/流水/订单支付记录三方一致

### 4.2 异常流程
- 取消支付
- 超时未支付
- 重复回调幂等
- 回调签名错误拦截

### 4.3 对账与补偿
- 支付平台订单号与本地 `payments` 可对齐
- 异常单可重试补偿，不产生重复扣款

## 5. 短信验收

### 5.1 发送链路
- 验证码发送成功率
- 模板变量渲染正确
- 限流生效（手机号/IP/设备）

### 5.2 安全
- 验证码有效期与错误次数限制
- 同一验证码不可重复使用

## 6. TRTC/实时能力验收

### 6.1 基础链路
- 房间创建/加入/离开
- Token 获取与过期处理

### 6.2 异常链路
- 弱网重连
- 断线重进
- 房间人数上限

## 7. 统一数据可信度校验

### 7.1 自动化
- 执行：
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_full_service_flow_acceptance.ps1`
  - `powershell -ExecutionPolicy Bypass -File api/scripts/run_data_integrity.ps1`

### 7.2 通过标准
- 全链路验收脚本 `PASS`
- 数据完整性违规项 `violations=0`

## 8. 发布阻断项（Go/No-Go）

| 编号 | 条件 | 结果 | 说明 |
| --- | --- | --- | --- |
| G1 | 核心支付链路通过 |  |  |
| G2 | OSS 上传链路通过 |  |  |
| G3 | 数据一致性为 0 |  |  |
| G4 | 关键告警可触发 |  |  |
| G5 | 回滚方案已演练 |  |  |

> 只有所有 G 项通过，才允许进入上线窗口。

## 9. 验收结论

- 总结：  
- 遗留风险：  
- 上线建议：`Go` / `No-Go`
- 签字确认：  

