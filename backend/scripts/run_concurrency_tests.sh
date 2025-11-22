#!/bin/bash

# GameLink Backend - 并发测试和性能测试运行器
# ==========================================

set -e

# 设置项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/.."
cd "$PROJECT_ROOT"

# 彩色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}==========================================${NC}"
    echo -e "${BLUE}GameLink Backend - 并发测试和性能测试运行器${NC}"
    echo -e "${BLUE}==========================================${NC}"
    echo
}

print_section() {
    echo -e "${GREEN}===== $1 =====${NC}"
    echo "$2"
    echo
}

print_error() {
    echo -e "${RED}错误: $1${NC}"
}

print_success() {
    echo -e "${GREEN}成功: $1${NC}"
}

run_race_tests() {
    print_section "运行竞态条件检测测试" "检测关键业务模块的竞态条件..."
    
    echo -e "${YELLOW}[订单服务 - 并发接单测试]${NC}"
    go test -race -v ./internal/service/order -run TestConcurrentAcceptOrder || true
    go test -race -v ./internal/service/order -run TestRaceConditionAcceptOrder || true
    
    echo
    echo -e "${YELLOW}[支付服务 - 并发支付测试]${NC}"
    go test -race -v ./internal/service/payment -run TestConcurrentCreatePayment || true
    go test -race -v ./internal/service/payment -run TestRaceConditionPaymentStatus || true
    
    echo
    echo -e "${YELLOW}[佣金服务 - 并发结算测试]${NC}"
    go test -race -v ./internal/service/commission -run TestConcurrentRecordCommission || true
    go test -race -v ./internal/service/commission -run TestRaceConditionCommissionSettlement || true
}

run_concurrency_tests() {
    print_section "运行并发测试" "测试高并发场景下的业务逻辑..."
    
    echo -e "${YELLOW}[订单服务 - 并发测试]${NC}"
    go test -v ./internal/service/order -run TestConcurrent || true
    go test -v ./internal/service/order -run TestRaceCondition || true
    
    echo
    echo -e "${YELLOW}[支付服务 - 并发测试]${NC}"
    go test -v ./internal/service/payment -run TestConcurrent || true
    go test -v ./internal/service/payment -run TestRaceCondition || true
    
    echo
    echo -e "${YELLOW}[佣金服务 - 并发测试]${NC}"
    go test -v ./internal/service/commission -run TestConcurrent || true
    go test -v ./internal/service/commission -run TestRaceCondition || true
}

run_benchmark_tests() {
    print_section "运行性能基准测试" "测试关键操作的性能表现..."
    
    echo -e "${YELLOW}[订单服务 - 性能测试]${NC}"
    go test -bench=. -benchmem ./internal/service/order -run=^$ || true
    
    echo
    echo -e "${YELLOW}[支付服务 - 性能测试]${NC}"
    go test -bench=. -benchmem ./internal/service/payment -run=^$ || true
    
    echo
    echo -e "${YELLOW}[佣金服务 - 性能测试]${NC}"
    go test -bench=. -benchmem ./internal/service/commission -run=^$ || true
    
    echo
    echo -e "${YELLOW}[认证服务 - 性能测试]${NC}"
    go test -bench=. -benchmem ./internal/service/auth -run=^$ || true
}

generate_benchmark_report() {
    print_section "生成性能测试报告" "生成详细的性能分析报告..."
    
    # 创建报告目录
    mkdir -p "var/benchmark_reports"
    
    echo -e "${YELLOW}[订单服务性能报告]${NC}"
    go test -bench=. -benchmem ./internal/service/order -run=^$ > var/benchmark_reports/order_benchmark.txt 2>&1 || true
    
    echo -e "${YELLOW}[支付服务性能报告]${NC}"
    go test -bench=. -benchmem ./internal/service/payment -run=^$ > var/benchmark_reports/payment_benchmark.txt 2>&1 || true
    
    echo -e "${YELLOW}[佣金服务性能报告]${NC}"
    go test -bench=. -benchmem ./internal/service/commission -run=^$ > var/benchmark_reports/commission_benchmark.txt 2>&1 || true
    
    echo -e "${YELLOW}[认证服务性能报告]${NC}"
    go test -bench=. -benchmem ./internal/service/auth -run=^$ > var/benchmark_reports/auth_benchmark.txt 2>&1 || true
    
    echo
    print_success "性能测试报告已生成到 var/benchmark_reports 目录"
    echo "查看报告文件了解详细的性能数据"
}

run_all_tests() {
    print_section "运行所有并发和性能测试" "全面检测并发问题和性能表现..."
    
    echo -e "${YELLOW}[阶段1] 竞态条件检测${NC}"
    run_race_tests
    
    echo
    echo -e "${YELLOW}[阶段2] 并发测试${NC}"
    run_concurrency_tests
    
    echo
    echo -e "${YELLOW}[阶段3] 性能基准测试${NC}"
    run_benchmark_tests
}

show_menu() {
    echo
    echo "请选择要运行的测试类型："
    echo "1) 运行竞态条件检测测试 (-race)"
    echo "2) 运行并发测试"
    echo "3) 运行性能基准测试 (Benchmark)"
    echo "4) 运行所有测试"
    echo "5) 生成性能测试报告"
    echo "q) 退出"
    echo
}

main() {
    print_header
    
    # 检查是否安装了Go
    if ! command -v go &> /dev/null; then
        print_error "未安装Go，请先安装Go语言环境"
        exit 1
    fi
    
    while true; do
        show_menu
        read -p "请输入选择 (1-5 或 q): " choice
        
        case $choice in
            1)
                run_race_tests
                ;;
            2)
                run_concurrency_tests
                ;;
            3)
                run_benchmark_tests
                ;;
            4)
                run_all_tests
                ;;
            5)
                generate_benchmark_report
                ;;
            q|Q)
                echo
                echo -e "${GREEN}感谢使用，再见！${NC}"
                exit 0
                ;;
            *)
                print_error "无效选择，请重新输入"
                ;;
        esac
        
        echo
        read -p "按回车键继续..."
        clear
        print_header
    done
}

# 运行主函数
main "$@"