@echo off
set TEST_DB_PORT=5433
cd api
go test ./internal/handler/admin/... -v -run "TestOrderHandler_Unit_CreateOrder_Success|TestGameHandler_UpdateGame_ValidationTests" %*
