#!/bin/bash

# GameLink 文档清理脚本
# 整理和归档临时文档文件

echo "🧹 开始清理 GameLink 项目文档..."

# 创建归档目录
mkdir -p docs/archive/temp-reports
mkdir -p docs/archive/implementation
mkdir -p docs/archive/features

# 移动前端根目录的临时报告
echo "📁 整理前端临时报告..."

# 实现相关的报告
IMPLEMENTATION_REPORTS=(
    "AUTO_SEARCH_FIX.md"
    "CLEANUP_SUMMARY.md"
    "CODE_REUSABILITY_REFACTOR.md"
    "CRUD_COMPLETE_REPORT.md"
    "CRUD_IMPLEMENTATION_SUMMARY.md"
    "DASHBOARD_FIX.md"
    "DATA_TABLE_IMPLEMENTATION.md"
    "MODULES_COMPLETION_REPORT.md"
)

for report in "${IMPLEMENTATION_REPORTS[@]}"; do
    if [ -f "frontend/$report" ]; then
        echo "  移动: $report -> docs/archive/implementation/"
        mv "frontend/$report" "docs/archive/implementation/"
    fi
done

# 功能相关的报告
FEATURE_REPORTS=(
    "Emoji清理和功能增强报告.md"
    "颜色系统优化报告.md"
    "代码整洁度和规范性评估报告.md"
)

for report in "${FEATURE_REPORTS[@]}"; do
    if [ -f "frontend/$report" ]; then
        echo "  移动: $report -> docs/archive/features/"
        mv "frontend/$report" "docs/archive/features/"
    fi
done

# 移动后端临时报告
echo "📁 整理后端临时报告..."
if [ -f "backend/内部接口模型整理.md" ]; then
    echo "  移动: 内部接口模型整理.md -> docs/archive/implementation/"
    mv "backend/内部接口模型整理.md" "docs/archive/implementation/"
fi

# 创建归档索引
echo "📋 创建归档索引..."
cat > docs/archive/README.md << 'EOF'
# 📦 文档归档

此目录包含 GameLink 项目的临时文档、实现报告和功能记录。

## 📂 目录结构

### implementation/ - 实现报告
包含各种功能实现过程中的详细报告和技术记录。

### features/ - 功能报告
包含特定功能的开发、优化和测试报告。

### temp-reports/ - 临时报告
包含开发过程中的临时文档和草稿。

### reports/ - 重要报告
包含重要的项目报告和迁移记录。

## 🗂️ 文档查找

如果您在寻找特定的文档：

1. **当前活跃文档**: 查看 `../` 目录
2. **实现细节**: 查看 `implementation/` 目录
3. **功能记录**: 查看 `features/` 目录
4. **重要报告**: 查看 `reports/` 目录

## ⚠️ 注意事项

- 归档文档仅供参考，可能不是最新状态
- 当前有效的规范和指南在上级目录
- 如需更新文档，请创建新的 PR 而不是修改归档文件

---

最后更新: $(date +%Y-%m-%d)
EOF

echo "✅ 文档清理完成！"
echo "📊 统计信息:"
echo "  - 实现报告: $(find docs/archive/implementation -name "*.md" | wc -l) 个文件"
echo "  - 功能报告: $(find docs/archive/features -name "*.md" | wc -l) 个文件"
echo "  - 归档总计: $(find docs/archive -name "*.md" | wc -l) 个文件"
echo ""
echo "📚 查看归档索引: docs/archive/README.md"