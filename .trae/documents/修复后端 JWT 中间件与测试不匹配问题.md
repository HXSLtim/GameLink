**问题定位**

* `middleware.JWTAuth`/`OptionalAuth`当前实现为有参函数：`func JWTAuth(secretKey string)` 与 `func OptionalAuth(secretKey string)`，文件 `backend/internal/handler/middleware/jwtAuth.go:21`、`:166`。

* 测试文件 `backend/internal/handler/middleware/jwt_test.go` 多处无参调用：`router.Use(JWTAuth())` 与 `router.Use(OptionalAuth())`（例如 `jwt_test.go:27,49,65,94,142,157,182,214,235,256` 以及 `284,301,328,349,374`）。这会导致编译错误：缺少必需的参数。

* 另有一处断言与实现不一致：过期 Token 的错误消息测试包含“无效”，而实现返回“Token已过期”（`jwtAuth.go:66-69`）。

**修改方案**

* 将所有无参中间件调用改为显式传参，使用环境变量 `JWT_SECRET_KEY`：

  * 替换 `router.Use(JWTAuth())` 为 `router.Use(JWTAuth(os.Getenv("JWT_SECRET_KEY")))`。

  * 替换 `router.Use(OptionalAuth())` 为 `router.Use(OptionalAuth(os.Getenv("JWT_SECRET_KEY")))`。

* 修正过期 Token 的断言：

  * 在 `ExpiredToken` 用例中，将 `assert.Contains(response["message"].(string), "无效")` 改为包含“过期”（例如：`assert.Contains(response["message"].(string), "过期")`），与实现的“Token已过期”保持一致。

* 其余断言与实现保持一致：

  * 缺少 Authorization 头断言“缺少Authorization”与实现“缺少Authorization头”相容。

  * `InvalidTokenFormat` 仍然返回 401，消息为“Authorization头格式错误，应为'Bearer <token>'”。

  * `OptionalAuth` 在密钥缺失或过短（`len(secretKey) < 32`）时返回 503，测试通过传参保留该行为。

**具体修改点（示例）**

* `backend/internal/handler/middleware/jwt_test.go:27`、`49` 等：

  * `router.Use(JWTAuth())` → `router.Use(JWTAuth(os.Getenv("JWT_SECRET_KEY")))`

* `backend/internal/handler/middleware/jwt_test.go:284`、`301` 等：

  * `router.Use(OptionalAuth())` → `router.Use(OptionalAuth(os.Getenv("JWT_SECRET_KEY")))`

* `backend/internal/handler/middleware/jwt_test.go:115-117`：

  * 过期消息断言改为包含“过期”。

**验证步骤**

* 设置环境变量：`JWT_SECRET_KEY=test-secret-key-that-is-32-characters-long`（测试用例已设置）。

* 运行后端测试：`go test ./...`。

* 重点检查：

  * `middleware/jwt_test.go` 所有子测试编译并通过（缺头、格式错误、有效、过期、密钥缺失/过短、自动刷新与刷新提示、RequireRole、OptionalAuth 分支）。

  * 与路由主干 `backend/internal/router/router.go:107-110` 的有参调用保持一致，无需改动生产代码。

**影响与兼容性**

* 仅修改测试文件，保持生产实现的“配置注入密钥”模式；不影响 `router.Setup()` 中的认证初始化。

* 断言文案调整为与当前实现一致，避免语义不匹配导致的用例失败。

**后续优化（可选）**

* 如需支持旧测试风格，可在中间件中新增包装函数：`JWTAuthFromEnv()`、`OptionalAuthFromEnv()`，内部读取 `JWT_SECRET_KEY` 并调用有参版本；测试与生产均可选择使用。但本次先以最小改动修复。

