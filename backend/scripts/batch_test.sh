#!/bin/bash

# 批量测试脚本
# 用于快速为多个包生成测试并运行

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "================================"
echo "  GameLink 测试覆盖率工具"
echo "================================"
echo ""

# 显示使用说明
show_usage() {
    echo "使用方法:"
    echo "  $0 generate <package>   - 为指定包生成测试骨架"
    echo "  $0 test <package>       - 测试指定包并显示覆盖率"
    echo "  $0 report <package>     - 生成指定包的HTML覆盖率报告"
    echo "  $0 all                  - 测试所有包"
    echo "  $0 progress             - 显示测试进度"
    echo ""
    echo "示例:"
    echo "  $0 generate internal/cache"
    echo "  $0 test internal/cache"
    echo "  $0 all"
}

# 生成测试骨架
generate_tests() {
    local package=$1
    echo -e "${YELLOW}为 $package 生成测试骨架...${NC}"
    go run scripts/generate_test_skeleton.go "$package"
    echo -e "${GREEN}✓ 完成${NC}"
}

# 测试单个包
test_package() {
    local package=$1
    echo -e "${YELLOW}测试 $package...${NC}"

    cd "$package" || exit 1

    # 运行测试
    if go test -cover -v; then
        echo -e "${GREEN}✓ $package 测试通过${NC}"

        # 生成覆盖率报告
        go test -coverprofile=coverage.out -covermode=atomic > /dev/null 2>&1
        coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
        echo -e "${GREEN}  覆盖率: $coverage${NC}"
    else
        echo -e "${RED}✗ $package 测试失败${NC}"
        cd - > /dev/null
        return 1
    fi

    cd - > /dev/null
}

# 生成HTML报告
generate_report() {
    local package=$1
    echo -e "${YELLOW}生成 $package 的HTML覆盖率报告...${NC}"

    cd "$package" || exit 1
    go test -coverprofile=coverage.out -covermode=atomic > /dev/null 2>&1
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ 报告已生成: $package/coverage.html${NC}"
    cd - > /dev/null
}

# 测试所有包
test_all() {
    echo -e "${YELLOW}测试所有包...${NC}"
    echo ""

    # 基础设施层
    echo "=== 基础设施层 ==="
    test_package "internal/apierr"
    test_package "internal/cache"
    test_package "internal/config"
    test_package "internal/db"
    echo ""

    # 数据访问层
    echo "=== 数据访问层 ==="
    for dir in internal/repository/*/; do
        if [ -f "$dir"/*.go ]; then
            test_package "$dir"
        fi
    done
    echo ""

    # 业务逻辑层
    echo "=== 业务逻辑层 ==="
    for dir in internal/service/*/; do
        if [ -f "$dir"/*.go ]; then
            test_package "$dir"
        fi
    done
    echo ""

    # HTTP处理层
    echo "=== HTTP处理层 ==="
    for dir in internal/handler/*/; do
        if [ -f "$dir"/*.go ]; then
            test_package "$dir"
        fi
    done

    echo ""
    echo "=== 总体覆盖率 ==="
    go test ./... -coverprofile=coverage.out -covermode=atomic > /dev/null 2>&1
    go tool cover -func=coverage.out | grep total
}

# 显示测试进度
show_progress() {
    echo "=== 测试覆盖率进度 ==="
    echo ""

    total_files=0
    tested_files=0

    # 统计源文件
    while IFS= read -r file; do
        ((total_files++))
        test_file="${file%.go}_test.go"
        if [ -f "$test_file" ]; then
            ((tested_files++))
            echo -e "${GREEN}✓${NC} $file"
        else
            echo -e "${RED}✗${NC} $file"
        fi
    done < <(find internal -name "*.go" -not -name "*_test.go" -not -name "doc.go" -not -name "wire_gen.go")

    echo ""
    echo "进度: $tested_files / $total_files 文件"
    percentage=$((tested_files * 100 / total_files))
    echo "完成度: ${percentage}%"

    # 显示覆盖率
    if [ -f "coverage.out" ]; then
        echo ""
        echo "总体覆盖率:"
        go tool cover -func=coverage.out | grep total
    fi
}

# 主逻辑
case "${1:-}" in
    generate)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请指定包路径${NC}"
            show_usage
            exit 1
        fi
        generate_tests "$2"
        ;;
    test)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请指定包路径${NC}"
            show_usage
            exit 1
        fi
        test_package "$2"
        ;;
    report)
        if [ -z "$2" ]; then
            echo -e "${RED}错误: 请指定包路径${NC}"
            show_usage
            exit 1
        fi
        generate_report "$2"
        ;;
    all)
        test_all
        ;;
    progress)
        show_progress
        ;;
    *)
        show_usage
        exit 1
        ;;
esac
