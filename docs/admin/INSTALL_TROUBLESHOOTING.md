# GameLink 安装问题解决方案

## ❌ 遇到安装错误？

别担心，这里有多种解决方案！

---

## 🚀 方案 1: 使用淘宝镜像 (推荐)

**最常见的原因是网络问题，使用国内镜像可解决90%的问题**

```bash
cd admin
npm config set registry https://registry.npmmirror.com
npm install -D vite-plugin-pwa picocolors workbox-window
```

**永久配置镜像**:
```bash
npm config set registry https://registry.npmmirror.com
```

**临时使用镜像**:
```bash
npm install -D vite-plugin-pwa picocolors workbox-window --registry=https://registry.npmmirror.com
```

---

## 🧹 方案 2: 清除缓存重试

```bash
cd admin

# 清除npm缓存
npm cache clean --force

# 删除node_modules和lock文件
Remove-Item -Recurse -Force node_modules
Remove-Item package-lock.json

# 重新安装
npm install -D vite-plugin-pwa picocolors workbox-window
```

---

## 🔧 方案 3: 修复脚本 (自动化)

```powershell
cd admin
.\fix-install.ps1
```

脚本会：
- ✅ 检查Node.js和npm版本
- ✅ 测试网络连接
- ✅ 提供5种修复方案
- ✅ 自动应用修复

---

## 📦 方案 4: 分步安装

**先安装基础依赖**:
```bash
cd admin
npm install
```

**然后单独安装PWA依赖**:
```bash
npm install -D vite-plugin-pwa
npm install -D picocolors
npm install -D workbox-window
```

---

## 🎯 方案 5: 跳过PWA (临时方案)

如果PWA依赖一直安装失败，可以先跳过：

```bash
cd admin

# 仅安装基础依赖
npm install

# 修改vite.config.ts，注释掉PWA插件
# 或使用 --no-verify 标志启动
```

稍后网络改善时再安装PWA：
```bash
npm install -D vite-plugin-pwa picocolors workbox-window
```

---

## 🔍 方案 6: 详细诊断

**查看详细错误信息**:
```bash
npm install -D vite-plugin-pwa picocolors workbox-window --verbose
```

**检查npm配置**:
```bash
npm config list
```

**检查网络**:
```bash
# 测试npm registry连接
ping registry.npmjs.org

# 查看当前使用的registry
npm config get registry
```

---

## 🛠️ 常见错误及解决方案

### 错误 1: ECONNREFUSED / ETIMEDOUT

**原因**: 无法连接到npm registry

**解决方案**:
```bash
# 使用淘宝镜像
npm config set registry https://registry.npmmirror.com
npm install -D vite-plugin-pwa picocolors workbox-window
```

### 错误 2: peer dependency warning

**原因**: 依赖版本冲突

**解决方案**:
```bash
# 使用legacy模式
npm install -D vite-plugin-pwa picocolors workbox-window --legacy-peer-deps
```

### 错误 3: permission denied

**原因**: 权限问题

**解决方案**:
```bash
# 以管理员身份运行PowerShell，然后
npm install -D vite-plugin-pwa picocolors workbox-window
```

### 错误 4: disk space / out of memory

**原因**: 磁盘空间不足或内存不足

**解决方案**:
```bash
# 清理npm缓存
npm cache clean --force

# 清理全局缓存
npm cache verify
```

### 错误 5: package.json not found

**原因**: 不在正确的目录

**解决方案**:
```bash
# 确保在admin目录
cd admin
dir package.json  # 确认文件存在
npm install -D vite-plugin-pwa picocolors workbox-window
```

---

## 🌐 其他可用的镜像源

除了淘宝镜像，还可以尝试：

### 华为云镜像
```bash
npm config set registry https://mirrors.huaweicloud.com/repository/npm/
```

### 腾讯云镜像
```bash
npm config set registry https://mirrors.cloud.tencent.com/npm/
```

### 官方源 (如果网络良好)
```bash
npm config set registry https://registry.npmjs.org
```

**恢复默认源**:
```bash
npm config set registry https://registry.npmjs.org
```

---

## 🔧 高级解决方案

### 使用代理

如果需要通过代理访问：
```bash
npm config set proxy http://proxy-server:port
npm config set https-proxy http://proxy-server:port
```

**取消代理**:
```bash
npm config delete proxy
npm config delete https-proxy
```

### 使用yarn代替npm

```bash
# 安装yarn
npm install -g yarn

# 使用yarn安装
cd admin
yarn add -D vite-plugin-pwa picocolors workbox-window
```

### 使用pnpm (更快)

```bash
# 安装pnpm
npm install -g pnpm

# 使用pnpm安装
cd admin
pnpm add -D vite-plugin-pwa picocolors workbox-window
```

---

## 📋 检查清单

在安装前，请确保：

- [ ] Node.js版本 >= 18.0.0
  - 检查: `node --version`
  - 下载: https://nodejs.org/

- [ ] npm版本 >= 9.0.0
  - 检查: `npm --version`

- [ ] 在正确的目录 (admin/)
  - 检查: `dir package.json`

- [ ] 网络连接正常
  - 测试: `ping registry.npmjs.org`

- [ ] 磁盘空间充足 (>2GB)

- [ ] 没有防火墙/杀毒软件阻止

---

## 🎯 快速修复命令

复制粘贴运行：

```powershell
# 进入admin目录
cd admin

# 配置淘宝镜像
npm config set registry https://registry.npmmirror.com

# 清除缓存
npm cache clean --force

# 安装依赖
npm install -D vite-plugin-pwa picocolors workbox-window
```

---

## 💬 仍无法解决？

1. **查看详细日志**:
   ```bash
   npm install -D vite-plugin-pwa picocolors workbox-window --verbose > install-error.log
   ```
   然后查看 `install-error.log` 文件

2. **运行诊断脚本**:
   ```powershell
   .\fix-install.ps1
   ```

3. **查看npm文档**:
   https://docs.npmjs.com/

4. **在GitHub搜索类似问题**:
   https://github.com/npm/cli/issues

---

## ✅ 安装成功后的下一步

```bash
# 启动开发服务器
npm run dev
```

然后访问显示的地址（通常是 http://localhost:5173）

---

**最后更新**: 2026-01-04
**维护者**: GameLink Team
