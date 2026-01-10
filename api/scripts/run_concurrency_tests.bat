@echo off
echo ==========================================
echo GameLink Backend - 并发测试和性能测试运行器
echo ==========================================
echo.

REM 设置项目根目录
set PROJECT_ROOT=%~dp0..
cd /d "%PROJECT_ROOT%"

echo [1] 运行竞态条件检测测试 (-race)
echo [2] 运行并发测试
echo [3] 运行性能基准测试 (Benchmark)
echo [4] 运行所有测试
echo [5] 生成性能测试报告
echo.
set /p choice=请选择要运行的测试类型 (1-5): 

if "%choice%"=="1" goto RACE_TEST
if "%choice%"=="2" goto CONCURRENCY_TEST
if "%choice%"=="3" goto BENCHMARK_TEST
if "%choice%"=="4" goto ALL_TESTS
if "%choice%"=="5" goto BENCHMARK_REPORT
echo 无效选择，退出脚本。
goto END

:RACE_TEST
echo.
echo ===== 运行竞态条件检测测试 =====
echo 检测关键业务模块的竞态条件...
echo.

echo [订单服务 - 并发接单测试]
go test -race -v ./internal/service/order -run TestConcurrentAcceptOrder
go test -race -v ./internal/service/order -run TestRaceConditionAcceptOrder

echo.
echo [支付服务 - 并发支付测试]
go test -race -v ./internal/service/payment -run TestConcurrentCreatePayment
go test -race -v ./internal/service/payment -run TestRaceConditionPaymentStatus

echo.
echo [佣金服务 - 并发结算测试]
go test -race -v ./internal/service/commission -run TestConcurrentRecordCommission
go test -race -v ./internal/service/commission -run TestRaceConditionCommissionSettlement

goto END

:CONCURRENCY_TEST
echo.
echo ===== 运行并发测试 =====
echo 测试高并发场景下的业务逻辑...
echo.

echo [订单服务 - 并发测试]
go test -v ./internal/service/order -run TestConcurrent
go test -v ./internal/service/order -run TestRaceCondition

echo.
echo [支付服务 - 并发测试]
go test -v ./internal/service/payment -run TestConcurrent
go test -v ./internal/service/payment -run TestRaceCondition

echo.
echo [佣金服务 - 并发测试]
go test -v ./internal/service/commission -run TestConcurrent
go test -v ./internal/service/commission -run TestRaceCondition

goto END

:BENCHMARK_TEST
echo.
echo ===== 运行性能基准测试 =====
echo 测试关键操作的性能表现...
echo.

echo [订单服务 - 性能测试]
go test -bench=. -benchmem ./internal/service/order -run=^$

echo.
echo [支付服务 - 性能测试]
go test -bench=. -benchmem ./internal/service/payment -run=^$

echo.
echo [佣金服务 - 性能测试]
go test -bench=. -benchmem ./internal/service/commission -run=^$

echo.
echo [认证服务 - 性能测试]
go test -bench=. -benchmem ./internal/service/auth -run=^$

goto END

:ALL_TESTS
echo.
echo ===== 运行所有并发和性能测试 =====
echo 全面检测并发问题和性能表现...
echo.

echo [阶段1] 竞态条件检测
call :RACE_TEST

echo.
echo [阶段2] 并发测试
call :CONCURRENCY_TEST

echo.
echo [阶段3] 性能基准测试
call :BENCHMARK_TEST

goto END

:BENCHMARK_REPORT
echo.
echo ===== 生成性能测试报告 =====
echo 生成详细的性能分析报告...
echo.

REM 创建报告目录
if not exist "var\benchmark_reports" mkdir "var\benchmark_reports"

echo [订单服务性能报告]
go test -bench=. -benchmem ./internal/service/order -run=^$ > var\benchmark_reports\order_benchmark.txt 2>&1

echo [支付服务性能报告]
go test -bench=. -benchmem ./internal/service/payment -run=^$ > var\benchmark_reports\payment_benchmark.txt 2>&1

echo [佣金服务性能报告]
go test -bench=. -benchmem ./internal/service/commission -run=^$ > var\benchmark_reports\commission_benchmark.txt 2>&1

echo [认证服务性能报告]
go test -bench=. -benchmem ./internal/service/auth -run=^$ > var\benchmark_reports\auth_benchmark.txt 2>&1

echo.
echo 性能测试报告已生成到 var\benchmark_reports 目录
echo 查看报告文件了解详细的性能数据

goto END

:END
echo.
echo ==========================================
echo 测试运行完成！
echo ==========================================
echo.
echo 提示：
echo - 使用 'go test -race' 检测竞态条件
echo - 使用 'go test -bench=.' 运行性能测试
echo - 使用 'go test -benchmem' 查看内存分配
echo.
pause