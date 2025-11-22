@echo off
echo ==========================================
echo GameLink Backend - 简单并发测试验证
echo ==========================================
echo.

echo [1] 运行订单并发接单测试
go test -v ./internal/service/order -run TestConcurrentAcceptOrder

echo.
echo [2] 运行支付并发测试
go test -v ./internal/service/payment -run TestConcurrentCreatePayment

echo.
echo [3] 运行佣金并发测试
go test -v ./internal/service/commission -run TestConcurrentRecordCommission

echo.
echo [4] 运行性能基准测试
go test -bench=BenchmarkCreateOrder -benchtime=1x ./internal/service/order -run=^TestConcurrentAcceptOrder$

echo.
echo ==========================================
echo 测试完成！
echo ==========================================
pause