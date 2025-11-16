#!/bin/bash

# 修复Swagger注释中的Go泛型语法
# 将 model.APIResponse[Type] 替换为 model.ErrorResponse 或保持原样

echo "修复Swagger注释中的泛型语法..."

# 修复所有handler文件中的Swagger注释
find internal/handler -name "*.go" -type f ! -name "*_test.go" | while read file; do
    echo "处理: $file"

    # 替换 Failure 注释中的泛型
    sed -i 's|// @Failure.*{object}  model\.APIResponse\[any\]|// @Failure      400            {object}  model.ErrorResponse|g' "$file"
    sed -i 's|// @Failure.*{object}  model\.APIResponse\[interface{}\]|// @Failure      400            {object}  model.ErrorResponse|g' "$file"
done

echo "完成！"
echo "请手动检查并修复 @Success 注释中的泛型类型，为每个响应创建具体的Response类型。"
