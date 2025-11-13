@echo off
setlocal enabledelayedexpansion

REM GameLink 快速启动脚本 (Windows)
REM 一键搭建开发环境

title GameLink 快速启动工具

REM 项目信息
set PROJECT_NAME=GameLink
set VERSION=v2.1.0
set BACKEND_PORT=8080
set FRONTEND_PORT=5173

REM 颜色定义 (Windows 10+ 支持 ANSI 转义序列)
set "RED=[91m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "BLUE=[94m"
set "NC=[0m"

REM 显示横幅
echo %BLUE%
echo ╔══════════════════════════════════════════════════════════════╗
echo ║                    🎮 GameLink 快速启动工具                  ║
echo ║                                                              ║
echo ║    现代化游戏陪玩管理平台 - Go + React 全栈项目               ║
echo ║                       %VERSION%                          ║
echo ╚══════════════════════════════════════════════════════════════╝
echo %NC%

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% == 0 (
    echo %GREEN%[✓] 检测到管理员权限%NC%
) else (
    echo %YELLOW%[!] 建议使用管理员权限运行%NC%
)

echo.
echo %GREEN%开始检查系统要求...%NC%

REM 检查必要命令
echo [%time%] 检查必要命令...

call :CheckCommand git Git
call :CheckCommand curl Curl
call :CheckCommand wget Wget

REM 检查 Docker
echo [%time%] 检查 Docker...
docker --version >nul 2>&1
if %errorLevel% == 0 (
    docker info >nul 2>&1
    if !errorLevel! == 0 (
        echo %GREEN%[✓] Docker 已安装并运行%NC%
        set DEPLOY_MODE=docker
    ) else (
        echo %RED%[✗] Docker 未运行%NC%
        set DEPLOY_MODE=local
    )
) else (
    echo %YELLOW%[!] Docker 未安装，将使用本地部署模式%NC%
    set DEPLOY_MODE=local
)

REM 检查 Go
echo [%time%] 检查 Go...
go version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
    echo %GREEN%[✓] Go 已安装: !GO_VERSION!%NC%
) else (
    echo %RED%[✗] Go 未安装%NC%
    if "!DEPLOY_MODE!"=="local" (
        echo.
        echo %RED%请先安装 Go 1.25.3+:%NC%
        echo 1. 访问 https://golang.org/dl/
        echo 2. 下载并安装 Go
        echo 3. 重启命令提示符
        pause
        exit /b 1
    )
)

REM 检查 Node.js
echo [%time%] 检查 Node.js...
node --version >nul 2>&1
if %errorLevel% == 0 (
    for /f "tokens=*" %%i in ('node --version') do set NODE_VERSION=%%i
    echo !NODE_VERSION! | findstr /R "^v1[8-9]\|^v2[0-9]" >nul
    if !errorLevel! == 0 (
        echo %GREEN%[✓] Node.js 已安装: !NODE_VERSION!%NC%
    ) else (
        echo %YELLOW%[!] Node.js 版本过低: !NODE_VERSION!，需要 18+ %NC%
        echo 请访问 https://nodejs.org/ 安装最新版本
        pause
        exit /b 1
    )
) else (
    echo %RED%[✗] Node.js 未安装%NC%
    echo.
    echo %RED%请先安装 Node.js 18+:%NC%
    echo 1. 访问 https://nodejs.org/
    echo 2. 下载并安装 Node.js
    echo 3. 重启命令提示符
    pause
    exit /b 1
)

REM 检查端口占用
echo [%time%] 检查端口占用...
call :CheckPort %BACKEND_PORT% 后端API
call :CheckPort %FRONTEND_PORT% 前端应用

REM 设置环境变量
echo [%time%] 设置环境变量...
if not exist .env (
    if exist .env.example (
        copy .env.example .env >nul
        echo %GREEN%[✓] 已创建 .env 配置文件%NC%
    ) else (
        echo %YELLOW%[!] 未找到 .env.example，创建默认配置%NC%
        call :CreateDefaultEnv
    )
) else (
    echo %GREEN%[✓] .env 配置文件已存在%NC%
)

if "!DEPLOY_MODE!"=="docker" (
    if not exist docker-compose.yml (
        if exist docker-compose.example.yml (
            copy docker-compose.example.yml docker-compose.yml >nul
            echo %GREEN%[✓] 已创建 docker-compose.yml 配置文件%NC%
        )
    )
)

REM 生成随机密钥
powershell -Command "if ((Get-Content .env) -match 'JWT_SECRET=change_me') { $secret = -join ((48..57) + (65..90) + (97..122) | Get-Random -Count 32 | ForEach-Object {[char]$_}); (Get-Content .env) -replace 'JWT_SECRET=.*', \"JWT_SECRET=$secret\" | Set-Content .env; Write-Host '[✓] 已生成新的 JWT 密钥' }"

REM 创建必要目录
echo [%time%] 创建必要目录...
if not exist scripts mkdir scripts
if not exist logs mkdir logs
if not exist uploads mkdir uploads

REM 创建管理脚本
call :CreateManagementScripts

REM 部署服务
echo.
echo %GREEN%开始部署 %PROJECT_NAME%...%NC%
echo.

if "!DEPLOY_MODE!"=="docker" (
    call :DeployDocker
) else (
    call :DeployLocal
)

REM 验证部署
echo.
echo [%time%] 验证部署...
call :VerifyDeployment

REM 显示访问信息
call :ShowAccessInfo

echo.
echo %GREEN%🎉 快速启动完成！%NC%
echo.
pause
goto :eof

:CheckCommand
where %1 >nul 2>&1
if %errorLevel% == 0 (
    echo %GREEN%[✓] %2 已安装%NC%
) else (
    echo %RED%[✗] %2 未找到，请先安装%NC%
    pause
    exit /b 1
)
goto :eof

:CheckPort
netstat -an | findstr ":%1" >nul 2>&1
if %errorLevel% == 0 (
    echo %YELLOW%[!] 端口 %1 (%2) 已被占用%NC%
    set /p continue="是否继续？(y/N): "
    if /i not "!continue!"=="y" exit /b 1
) else (
    echo %GREEN%[✓] 端口 %1 (%2) 可用%NC%
)
goto :eof

:CreateDefaultEnv
(
echo # 应用配置
echo APP_ENV=development
echo DEBUG=true
echo.
echo # 数据库配置
echo DB_HOST=localhost
echo DB_PORT=3306
echo DB_NAME=gamelink_dev
echo DB_USER=gamelink
echo DB_PASSWORD=dev_password_123
echo.
echo # Redis 配置
echo REDIS_HOST=localhost
echo REDIS_PORT=6379
echo.
echo # JWT 配置
echo JWT_SECRET=change_me_please_update_this_key
echo JWT_EXPIRE_HOURS=24
echo.
echo # 服务端口
echo API_PORT=%BACKEND_PORT%
echo WEB_PORT=%FRONTEND_PORT%
) > .env
goto :eof

:CreateManagementScripts
REM 状态检查脚本
(
echo @echo off
echo echo === GameLink 服务状态 ===
echo.
echo if exist docker-compose.yml ^(
echo     echo Docker 模式:
echo     docker-compose ps
echo ^) else ^(
echo     echo 本地模式:
echo     tasklist ^| findstr user-service ^>nul
echo     if !errorLevel! == 0 ^(
echo         echo ✓ 后端服务运行中
echo     ^) else ^(
echo         ✗ 后端服务未运行
echo     ^)
echo.
echo     tasklist ^| findstr "node.exe" ^| findstr "vite" ^>nul
echo     if !errorLevel! == 0 ^(
echo         echo ✓ 前端服务运行中
echo     ^) else ^(
echo         ✗ 前端服务未运行
echo     ^)
echo ^)
echo.
echo 端口检查:
echo netstat -an ^| findstr ":%BACKEND_PORT%"
echo netstat -an ^| findstr ":%FRONTEND_PORT%"
echo pause
) > scripts\status.bat

REM 停止服务脚本
(
echo @echo off
echo echo 停止 GameLink 服务...
echo.
echo if exist docker-compose.yml ^(
echo     echo 停止 Docker 服务...
echo     docker-compose down
echo ^) else ^(
echo     echo 停止本地服务...
echo     taskkill /f /im user-service.exe 2^>nul
echo     taskkill /f /im node.exe 2^>nul
echo ^)
echo.
echo echo 服务已停止
echo pause
) > scripts\stop.bat

REM 重启服务脚本
(
echo @echo off
echo echo 重启 GameLink 服务...
echo.
echo call scripts\stop.bat
echo timeout /t 5 /nobreak ^>nul
echo call quick-start.bat
) > scripts\restart.bat

echo %GREEN%[✓] 管理脚本创建完成%NC%
goto :eof

:DeployDocker
echo [%time%] 使用 Docker 部署...

echo 构建 Docker 镜像...
docker-compose build
if %errorLevel% neq 0 (
    echo %RED%[✗] Docker 镜像构建失败%NC%
    pause
    exit /b 1
)

echo 启动服务...
docker-compose up -d
if %errorLevel% neq 0 (
    echo %RED%[✗] Docker 服务启动失败%NC%
    pause
    exit /b 1
)

echo 等待服务启动...
timeout /t 30 /nobreak >nul

echo 运行数据库迁移...
docker-compose exec -T api make migrate 2>nul

echo %GREEN%[✓] Docker 部署完成%NC%
goto :eof

:DeployLocal
echo [%time%] 使用本地模式部署...

REM 构建后端
echo 构建后端服务...
cd backend
go mod download
if %errorLevel% neq 0 (
    echo %RED%[✗] Go 依赖下载失败%NC%
    pause
    exit /b 1
)

go build -o bin\user-service.exe .\cmd\user-service
if %errorLevel% neq 0 (
    echo %RED%[✗] 后端构建失败%NC%
    pause
    exit /b 1
)

REM 构建前端
echo 构建前端应用...
cd ..\frontend
npm install
if %errorLevel% neq 0 (
    echo %RED%[✗] 前端依赖安装失败%NC%
    pause
    exit /b 1
)

REM 启动服务
echo 启动服务...
cd ..\backend
start "GameLink API" /min cmd /c "bin\user-service.exe > ..\logs\api.log 2>&1"

cd ..\frontend
start "GameLink Frontend" /min cmd /c "npm run dev > ..\logs\frontend.log 2>&1"

cd ..

echo %GREEN%[✓] 本地部署完成%NC%
goto :eof

:VerifyDeployment
echo 验证后端服务...
set /a attempt=1
:verify_loop
powershell -Command "try { (Invoke-WebRequest -Uri 'http://localhost:%BACKEND_PORT%/health' -TimeoutSec 5).StatusCode } catch { '999' }" > temp_status.txt
set /p status=<temp_status.txt
del temp_status.txt

if "%status%"=="200" (
    echo %GREEN%[✓] 后端服务验证成功%NC%
) else (
    if %attempt% geq 10 (
        echo %RED%[✗] 后端服务验证失败%NC%
    ) else (
        echo 等待后端服务启动... (%attempt%/10)
        timeout /t 5 /nobreak >nul
        set /a attempt+=1
        goto verify_loop
    )
)

echo 验证前端服务...
powershell -Command "try { (Invoke-WebRequest -Uri 'http://localhost:%FRONTEND_PORT%' -TimeoutSec 5).StatusCode } catch { '999' }" > temp_status.txt
set /p status=<temp_status.txt
del temp_status.txt

if "%status%"=="200" (
    echo %GREEN%[✓] 前端服务验证成功%NC%
) else (
    echo %YELLOW%[!] 前端服务验证失败，请手动检查%NC%
)
goto :eof

:ShowAccessInfo
echo.
echo %GREEN%
echo ╔══════════════════════════════════════════════════════════════╗
echo ║                    🎉 部署成功！                          ║
echo ║                                                              ║
echo ║    访问地址:                                                   ║
echo ║    🌐 前端应用: http://localhost:%FRONTEND_PORT%                ║
echo ║    🔌 后端API: http://localhost:%BACKEND_PORT%                 ║
echo ║    📚 API文档: http://localhost:%BACKEND_PORT%/swagger         ║
echo ║                                                              ║
echo ║    管理命令:                                                   ║
echo ║    📋 查看状态: scripts\status.bat                          ║
echo ║    🛑 停止服务: scripts\stop.bat                            ║
echo ║    🔄 重启服务: scripts\restart.bat                         ║
echo ╚══════════════════════════════════════════════════════════════╝
echo %NC%
goto :eof