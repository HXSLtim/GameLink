#!/bin/bash
# 文档整理脚本

echo "📁 开始整理文档..."

# 创建目录结构
mkdir -p docs/{crypto,api,design,guides,archive/reports}

# 移动加密文档
[ -f "CRYPTO_README.md" ] && mv CRYPTO_README.md docs/crypto/README.md && echo "✓ CRYPTO_README.md"
[ -f "CRYPTO_INTEGRATION.md" ] && mv CRYPTO_INTEGRATION.md docs/crypto/INTEGRATION.md && echo "✓ CRYPTO_INTEGRATION.md"
[ -f "CRYPTO_MIDDLEWARE.md" ] && mv CRYPTO_MIDDLEWARE.md docs/crypto/MIDDLEWARE.md && echo "✓ CRYPTO_MIDDLEWARE.md"
[ -f "CRYPTO_USAGE_EXAMPLES.md" ] && mv CRYPTO_USAGE_EXAMPLES.md docs/crypto/EXAMPLES.md && echo "✓ CRYPTO_USAGE_EXAMPLES.md"
[ -f "ENV_CONFIG_GUIDE.md" ] && mv ENV_CONFIG_GUIDE.md docs/crypto/ENV_CONFIG.md && echo "✓ ENV_CONFIG_GUIDE.md"

# 移动 API 文档
[ -f "API_REQUIREMENTS.md" ] && mv API_REQUIREMENTS.md docs/api/ && echo "✓ API_REQUIREMENTS.md"
[ -f "API_INTEGRATION_GUIDE.md" ] && mv API_INTEGRATION_GUIDE.md docs/api/INTEGRATION.md && echo "✓ API_INTEGRATION_GUIDE.md"
[ -f "内部接口模型整理.md" ] && mv 内部接口模型整理.md docs/api/backend-models.md && echo "✓ 内部接口模型整理.md"

# 移动设计文档
[ -f "DESIGN_SYSTEM.md" ] && mv DESIGN_SYSTEM.md docs/design/ && echo "✓ DESIGN_SYSTEM.md"
[ -f "CODING_STANDARDS.md" ] && mv CODING_STANDARDS.md docs/design/ && echo "✓ CODING_STANDARDS.md"

# 移动指南文档
[ -f "QUICK_START.md" ] && mv QUICK_START.md docs/guides/ && echo "✓ QUICK_START.md"
[ -f "MIGRATION_GUIDE.md" ] && mv MIGRATION_GUIDE.md docs/guides/ && echo "✓ MIGRATION_GUIDE.md"

# 归档报告
for file in *_COMPLETE.md *_SUMMARY.md *_REPORT.md; do
  [ -f "$file" ] && mv "$file" docs/archive/reports/ && echo "✓ 归档: $file"
done

# 创建文档索引
[ -f "DOCS_INDEX.md" ] && mv DOCS_INDEX.md docs/README.md && echo "✓ DOCS_INDEX.md → docs/README.md"

echo ""
echo "✅ 文档整理完成！"
echo "📂 文档位置: docs/"
echo ""
echo "📋 整理结果:"
echo "  • 加密文档: docs/crypto/ (5个)"
echo "  • API文档: docs/api/ (3个)"
echo "  • 设计文档: docs/design/ (2个)"
echo "  • 指南文档: docs/guides/ (2个)"
echo "  • 归档报告: docs/archive/reports/"
echo ""
echo "🔍 查看文档索引: cat docs/README.md"
