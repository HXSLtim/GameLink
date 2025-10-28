# 📁 文档整理建议

## 当前状态

项目根目录有 **40+** 个 Markdown 文档，建议整理归档。

## 🎯 整理方案

### 方案 A: 创建 docs 目录（推荐）

```bash
# 创建文档目录结构
mkdir -p docs/{crypto,api,design,guides,reports,archive}

# 移动加密文档
mv CRYPTO_*.md docs/crypto/
mv ENV_CONFIG_GUIDE.md docs/crypto/

# 移动 API 文档
mv API_*.md docs/api/
mv 内部接口模型整理.md docs/api/

# 移动设计文档
mv DESIGN_*.md docs/design/
mv CODING_STANDARDS.md docs/design/

# 移动指南文档
mv QUICK_START.md docs/guides/
mv MIGRATION_GUIDE.md docs/guides/
mv MVP_*.md docs/guides/

# 归档完成报告
mv *_COMPLETE.md docs/archive/
mv *_SUMMARY.md docs/archive/
mv *_REPORT.md docs/archive/
```

### 整理后的结构

```
frontend/
├── README.md                 # 项目主文档（保留在根目录）
├── docs/
│   ├── README.md            # 文档索引（DOCS_INDEX.md）
│   ├── crypto/              # 加密相关
│   │   ├── README.md       (CRYPTO_README.md)
│   │   ├── INTEGRATION.md  (CRYPTO_INTEGRATION.md)
│   │   ├── MIDDLEWARE.md   (CRYPTO_MIDDLEWARE.md)
│   │   ├── EXAMPLES.md     (CRYPTO_USAGE_EXAMPLES.md)
│   │   └── ENV_CONFIG.md   (ENV_CONFIG_GUIDE.md)
│   ├── api/                 # API 文档
│   │   ├── REQUIREMENTS.md
│   │   ├── INTEGRATION.md
│   │   └── backend-models.md (内部接口模型整理.md)
│   ├── design/              # 设计规范
│   │   ├── DESIGN_SYSTEM.md
│   │   └── CODING_STANDARDS.md
│   ├── guides/              # 开发指南
│   │   ├── QUICK_START.md
│   │   ├── MIGRATION_GUIDE.md
│   │   └── MVP_PLAN.md
│   └── archive/             # 历史文档
│       ├── reports/
│       └── deprecated/
└── src/
```

## 🚀 执行整理

### 一键整理脚本

```bash
#!/bin/bash
# organize-docs.sh

echo "开始整理文档..."

# 创建目录
mkdir -p docs/{crypto,api,design,guides,archive/reports}

# 加密文档
mv CRYPTO_README.md docs/crypto/README.md 2>/dev/null
mv CRYPTO_INTEGRATION.md docs/crypto/INTEGRATION.md 2>/dev/null
mv CRYPTO_MIDDLEWARE.md docs/crypto/MIDDLEWARE.md 2>/dev/null
mv CRYPTO_USAGE_EXAMPLES.md docs/crypto/EXAMPLES.md 2>/dev/null
mv ENV_CONFIG_GUIDE.md docs/crypto/ENV_CONFIG.md 2>/dev/null

# API 文档
mv API_REQUIREMENTS.md docs/api/ 2>/dev/null
mv API_INTEGRATION_GUIDE.md docs/api/INTEGRATION.md 2>/dev/null
mv 内部接口模型整理.md docs/api/backend-models.md 2>/dev/null

# 设计文档
mv DESIGN_SYSTEM.md docs/design/ 2>/dev/null
mv CODING_STANDARDS.md docs/design/ 2>/dev/null

# 指南文档
mv QUICK_START.md docs/guides/ 2>/dev/null
mv MIGRATION_GUIDE.md docs/guides/ 2>/dev/null
mv MVP_DEVELOPMENT_PLAN.md docs/guides/ 2>/dev/null

# 归档报告
mv *_COMPLETE.md docs/archive/reports/ 2>/dev/null
mv *_SUMMARY.md docs/archive/reports/ 2>/dev/null
mv *_REPORT.md docs/archive/reports/ 2>/dev/null

# 创建文档索引
mv DOCS_INDEX.md docs/README.md 2>/dev/null

echo "✅ 文档整理完成！"
echo "📁 文档位置: docs/"
```

### 手动执行

```bash
# 1. 保存脚本
cat > organize-docs.sh << 'EOF'
# ... (上面的脚本内容)
EOF

# 2. 添加执行权限
chmod +x organize-docs.sh

# 3. 执行
./organize-docs.sh
```

## 📝 更新引用

整理后需要更新以下文件中的文档链接：

### 1. README.md

```markdown
# 旧

- [加密文档](./CRYPTO_README.md)

# 新

- [加密文档](./docs/crypto/README.md)
```

### 2. package.json 中的文档链接

```json
{
  "homepage": "docs/README.md"
}
```

### 3. 各文档间的交叉引用

使用相对路径：

```markdown
# 在 docs/crypto/README.md 中

- 详细文档: [INTEGRATION.md](./INTEGRATION.md)
- API 文档: [../api/REQUIREMENTS.md](../api/REQUIREMENTS.md)
```

## 🔄 Git 操作建议

```bash
# 使用 git mv 保留历史
git mv CRYPTO_README.md docs/crypto/README.md
git mv CRYPTO_INTEGRATION.md docs/crypto/INTEGRATION.md
# ... 其他文件

# 或使用脚本后提交
./organize-docs.sh
git add docs/
git commit -m "docs: 整理项目文档结构"
```

## 📊 整理前后对比

### 整理前

```
frontend/
├── API_INTEGRATION_COMPLETE.md
├── API_INTEGRATION_GUIDE.md
├── API_REQUIREMENTS.md
├── CODE_QUALITY_IMPROVEMENTS_COMPLETE.md
├── CODING_STANDARDS.md
├── CRYPTO_INTEGRATION.md
├── CRYPTO_MIDDLEWARE.md
├── CRYPTO_README.md
├── CRYPTO_USAGE_EXAMPLES.md
├── DESIGN_SYSTEM.md
├── ENV_CONFIG_GUIDE.md
├── QUICK_START.md
├── ... (30+ 个文档)
└── src/
```

### 整理后

```
frontend/
├── README.md                 # 项目主文档
├── docs/
│   ├── README.md            # 文档索引
│   ├── crypto/              # 5个文件
│   ├── api/                 # 3个文件
│   ├── design/              # 2个文件
│   ├── guides/              # 3个文件
│   └── archive/             # 20+个文件
└── src/
```

## ✅ 整理收益

1. **结构清晰** - 文档分类明确
2. **易于查找** - 按功能组织
3. **便于维护** - 相关文档集中
4. **减少混乱** - 根目录整洁
5. **历史保留** - 归档而非删除

## 🎯 下一步

1. **立即整理** - 执行整理脚本
2. **更新链接** - 修改文档引用
3. **提交代码** - git commit
4. **通知团队** - 告知文档新位置

---

**需要帮助？** 查看 [DOCS_INDEX.md](./DOCS_INDEX.md)
