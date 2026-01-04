#!/bin/bash
export TEST_DB_PORT=5433
cd "D:/Desktop/Code/GameLink/api"
go test ./internal/handler/admin/... -v -run "$@"
