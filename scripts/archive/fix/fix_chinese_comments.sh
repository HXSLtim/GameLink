#!/bin/bash
# Fix common broken Chinese Swagger comments

cd "$(dirname "$0")"

# Function to fix file
fix_file() {
    local file=$1
    echo "Fixing: $file"

    # Common replacements
    sed -i 's|"开始时[^"]*"|"Start date (YYYY-MM-DD)"|g' "$file"
    sed -i 's|"结束时[^"]*"|"End date (YYYY-MM-DD)"|g' "$file"
    sed -i 's|"子分[^"]*"|"Sub category (solo/team/gift)"|g' "$file"
    sed -i 's|"是否激[^"]*"|"Is active"|g' "$file"
    sed -i 's|"订单状态，可多[^"]*"|"Order status (multiple allowed)"|g' "$file"
    sed -i 's|"退款信[^"]*"|"Refund information"|g' "$file"
    sed -i 's|"支付状[^"]*"|"Payment status"|g' "$file"
    sed -i 's|"导出列（逗号分隔[^"]*"|"Export columns (comma separated)"|g' "$file"
    sed -i 's|"月份筛[^"]*"|"Month filter (YYYY-MM)"|g' "$file"
    sed -i 's|"角色过滤，可多[^"]*"|"Role filter (multiple allowed)"|g' "$file"
    sed -i 's|"状态过滤，可多[^"]*"|"Status filter (multiple allowed)"|g' "$file"
}

# Fix all files with encoding issues
fix_file "internal/handler/admin/game.go"
fix_file "internal/handler/admin/item.go"
fix_file "internal/handler/admin/order.go"
fix_file "internal/handler/admin/player.go"
fix_file "internal/handler/admin/ranking.go"
fix_file "internal/handler/admin/review.go"
fix_file "internal/handler/admin/user.go"

echo "Done! All files fixed."
