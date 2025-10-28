#!/bin/bash
# 完整文档归档脚本

echo "📁 开始归档所有文档..."

# 创建完整的目录结构
mkdir -p docs/{features,refactoring,archive/{reports,deprecated}}

# ========== 功能文档 ==========
echo -e "\n📦 归档功能文档..."
[ -f "THEME_RIPPLE_EFFECT.md" ] && mv THEME_RIPPLE_EFFECT.md docs/features/ && echo "  ✓ THEME_RIPPLE_EFFECT.md"
[ -f "THEME_TOGGLE_GUIDE.md" ] && mv THEME_TOGGLE_GUIDE.md docs/features/ && echo "  ✓ THEME_TOGGLE_GUIDE.md"
[ -f "NAVIGATION_SYSTEM.md" ] && mv NAVIGATION_SYSTEM.md docs/features/ && echo "  ✓ NAVIGATION_SYSTEM.md"
[ -f "I18N_IMPLEMENTATION_GUIDE.md" ] && mv I18N_IMPLEMENTATION_GUIDE.md docs/features/ && echo "  ✓ I18N_IMPLEMENTATION_GUIDE.md"
[ -f "FIGMA_TO_CODE_GUIDE.md" ] && mv FIGMA_TO_CODE_GUIDE.md docs/features/ && echo "  ✓ FIGMA_TO_CODE_GUIDE.md"
[ -f "MOCK_LOGIN_GUIDE.md" ] && mv MOCK_LOGIN_GUIDE.md docs/features/ && echo "  ✓ MOCK_LOGIN_GUIDE.md"

# ========== 重构文档 ==========
echo -e "\n🔧 归档重构文档..."
[ -f "COLOR_VARIABLES_REFACTOR.md" ] && mv COLOR_VARIABLES_REFACTOR.md docs/refactoring/ && echo "  ✓ COLOR_VARIABLES_REFACTOR.md"
[ -f "IMPORT_PATH_GUIDE.md" ] && mv IMPORT_PATH_GUIDE.md docs/refactoring/ && echo "  ✓ IMPORT_PATH_GUIDE.md"
[ -f "PATH_ALIAS_FIX.md" ] && mv PATH_ALIAS_FIX.md docs/refactoring/ && echo "  ✓ PATH_ALIAS_FIX.md"
[ -f "HMR_AND_THEME_FIX.md" ] && mv HMR_AND_THEME_FIX.md docs/refactoring/ && echo "  ✓ HMR_AND_THEME_FIX.md"

# ========== 完成报告/总结（归档） ==========
echo -e "\n📊 归档报告文档..."
[ -f "API_INTEGRATION_SUCCESS.md" ] && mv API_INTEGRATION_SUCCESS.md docs/archive/reports/ && echo "  ✓ API_INTEGRATION_SUCCESS.md"
[ -f "README_IMPROVEMENTS.md" ] && mv README_IMPROVEMENTS.md docs/archive/reports/ && echo "  ✓ README_IMPROVEMENTS.md"
[ -f "ROUTES_REGISTERED.md" ] && mv ROUTES_REGISTERED.md docs/archive/reports/ && echo "  ✓ ROUTES_REGISTERED.md"
[ -f "前端代码整洁度评估报告.md" ] && mv 前端代码整洁度评估报告.md docs/archive/reports/ && echo "  ✓ 前端代码整洁度评估报告.md"

# ========== MVP 计划（移到 guides） ==========
echo -e "\n📋 归档指南文档..."
[ -f "MVP_DEVELOPMENT_PLAN.md" ] && mv MVP_DEVELOPMENT_PLAN.md docs/guides/ && echo "  ✓ MVP_DEVELOPMENT_PLAN.md"

# ========== 设计文档 V2（移到 design） ==========
echo -e "\n🎨 归档设计文档..."
[ -f "DESIGN_SYSTEM_V2.md" ] && mv DESIGN_SYSTEM_V2.md docs/design/ && echo "  ✓ DESIGN_SYSTEM_V2.md"

# ========== 其他整理文档 ==========
echo -e "\n📄 归档其他文档..."
[ -f "DOCS_ORGANIZATION.md" ] && mv DOCS_ORGANIZATION.md docs/archive/ && echo "  ✓ DOCS_ORGANIZATION.md"

# ========== 统计结果 ==========
echo -e "\n✅ 归档完成！\n"
echo "📊 归档统计："
echo "  • 功能文档: $(ls -1 docs/features/*.md 2>/dev/null | wc -l) 个"
echo "  • 重构文档: $(ls -1 docs/refactoring/*.md 2>/dev/null | wc -l) 个"
echo "  • 指南文档: $(ls -1 docs/guides/*.md 2>/dev/null | wc -l) 个"
echo "  • 设计文档: $(ls -1 docs/design/*.md 2>/dev/null | wc -l) 个"
echo "  • 归档报告: $(ls -1 docs/archive/reports/*.md 2>/dev/null | wc -l) 个"
echo ""
echo "📁 根目录剩余 MD 文件: $(ls -1 *.md 2>/dev/null | wc -l) 个"

# 列出剩余文件
remaining=$(ls -1 *.md 2>/dev/null)
if [ ! -z "$remaining" ]; then
  echo -e "\n📋 剩余文件列表:"
  ls -1 *.md 2>/dev/null | sed 's/^/  • /'
fi

echo -e "\n🎯 查看文档: cd docs && ls -R"
