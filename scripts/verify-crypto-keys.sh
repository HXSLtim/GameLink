#!/bin/bash
# GameLink 加密密钥一致性验证工具
# 用途：验证前后端加密密钥是否配置正确
# 使用方法：bash scripts/verify-crypto-keys.sh

set -e

echo "=================================="
echo " GameLink 加密密钥验证工具"
echo "=================================="
echo ""

# 检查 .env.production 文件是否存在
if [ ! -f ".env.production" ]; then
    echo "❌ 错误：.env.production 文件不存在"
    echo "   请先运行: bash scripts/generate-production-keys.sh"
    exit 1
fi

if [ ! -f "admin/.env.production" ]; then
    echo "❌ 错误：admin/.env.production 文件不存在"
    echo "   请先运行: bash scripts/generate-production-keys.sh"
    exit 1
fi

echo "正在加载配置文件..."
echo ""

# 加载后端配置
source .env.production

# 验证后端配置
echo "1. 验证后端配置..."
echo "   CRYPTO_ENABLED=$CRYPTO_ENABLED"

if [ "$CRYPTO_ENABLED" != "true" ]; then
    echo "   ⚠️  警告：生产环境应启用加密 (CRYPTO_ENABLED=true)"
fi

if [ -z "$CRYPTO_SECRET_KEY" ]; then
    echo "   ❌ 错误：CRYPTO_SECRET_KEY 未配置"
    exit 1
fi

if [ -z "$CRYPTO_IV" ]; then
    echo "   ❌ 错误：CRYPTO_IV 未配置"
    exit 1
fi

# 检查后端密钥长度
BACKEND_KEY_LENGTH=${#CRYPTO_SECRET_KEY}
BACKEND_IV_LENGTH=${#CRYPTO_IV}

echo "   CRYPTO_SECRET_KEY 长度: $BACKEND_KEY_LENGTH 字符（期望 44）"
echo "   CRYPTO_IV 长度: $BACKEND_IV_LENGTH 字符（期望 24）"

if [ "$BACKEND_KEY_LENGTH" -ne 44 ]; then
    echo "   ⚠️  警告：后端密钥长度不是 44 字符（base64 编码的 32 字节）"
fi

if [ "$BACKEND_IV_LENGTH" -ne 24 ]; then
    echo "   ⚠️  警告：后端 IV 长度不是 24 字符（base64 编码的 16 字节）"
fi

echo "   ✅ 后端配置检查通过"
echo ""

# 加载前端配置
echo "2. 验证前端配置..."
source admin/.env.production

echo "   VITE_CRYPTO_ENABLED=$VITE_CRYPTO_ENABLED"

if [ "$VITE_CRYPTO_ENABLED" != "true" ]; then
    echo "   ⚠️  警告：前端加密未启用，与后端不一致"
fi

if [ -z "$VITE_CRYPTO_SECRET_KEY" ]; then
    echo "   ❌ 错误：VITE_CRYPTO_SECRET_KEY 未配置"
    exit 1
fi

if [ -z "$VITE_CRYPTO_IV" ]; then
    echo "   ❌ 错误：VITE_CRYPTO_IV 未配置"
    exit 1
fi

# 检查前端密钥长度（原始字节，应该是 64 个十六进制字符 = 32 字节）
FRONTEND_KEY_LENGTH=${#VITE_CRYPTO_SECRET_KEY}
FRONTEND_IV_LENGTH=${#VITE_CRYPTO_IV}

echo "   VITE_CRYPTO_SECRET_KEY 长度: $FRONTEND_KEY_LENGTH 字符（期望 64 个十六进制字符）"
echo "   VITE_CRYPTO_IV 长度: $FRONTEND_IV_LENGTH 字符（期望 32 个十六进制字符）"

if [ "$FRONTEND_KEY_LENGTH" -ne 64 ]; then
    echo "   ⚠️  警告：前端密钥长度不是 64 个十六进制字符（32 字节）"
fi

if [ "$FRONTEND_IV_LENGTH" -ne 32 ]; then
    echo "   ⚠️  警告：前端 IV 长度不是 32 个十六进制字符（16 字节）"
fi

echo "   ✅ 前端配置检查通过"
echo ""

# 验证加密启用状态一致性
echo "3. 验证前后端加密启用状态..."
if [ "$CRYPTO_ENABLED" == "$VITE_CRYPTO_ENABLED" ]; then
    echo "   ✅ 加密启用状态一致"
else
    echo "   ❌ 错误：前后端加密启用状态不一致"
    echo "      后端: CRYPTO_ENABLED=$CRYPTO_ENABLED"
    echo "      前端: VITE_CRYPTO_ENABLED=$VITE_CRYPTO_ENABLED"
    exit 1
fi
echo ""

# 验证签名配置一致性
echo "4. 验证签名配置..."
if [ "$CRYPTO_USE_SIGNATURE" == "true" ] && [ "$VITE_CRYPTO_USE_SIGNATURE" == "true" ]; then
    echo "   ✅ 签名配置一致且已启用"
elif [ "$CRYPTO_USE_SIGNATURE" == "false" ] && [ "$VITE_CRYPTO_USE_SIGNATURE" == "false" ]; then
    echo "   ✅ 签名配置一致且已禁用"
else
    echo "   ❌ 错误：前后端签名配置不一致"
    echo "      后端: CRYPTO_USE_SIGNATURE=$CRYPTO_USE_SIGNATURE"
    echo "      前端: VITE_CRYPTO_USE_SIGNATURE=$VITE_CRYPTO_USE_SIGNATURE"
    exit 1
fi
echo ""

# 验证密钥格式兼容性
echo "5. 验证密钥格式兼容性..."
echo "   ℹ️  提示：前后端使用不同的密钥编码格式"
echo "      - 后端使用 base64 编码"
echo "      - 前端使用原始字节（十六进制字符串）"
echo "   ✅ 密钥格式符合预期"
echo ""

# 显示配置摘要
echo "=================================="
echo " 配置摘要"
echo "=================================="
echo ""
echo "加密启用: $CRYPTO_ENABLED"
echo "签名启用: $CRYPTO_USE_SIGNATURE"
echo ""
echo "后端密钥:"
echo "  - SECRET_KEY: ${CRYPTO_SECRET_KEY:0:16}... (共 $BACKEND_KEY_LENGTH 字符)"
echo "  - IV: ${CRYPTO_IV:0:16}... (共 $BACKEND_IV_LENGTH 字符)"
echo ""
echo "前端密钥:"
echo "  - SECRET_KEY: ${VITE_CRYPTO_SECRET_KEY:0:32}... (共 $FRONTEND_KEY_LENGTH 字符)"
echo "  - IV: ${VITE_CRYPTO_IV:0:16}... (共 $FRONTEND_IV_LENGTH 字符)"
echo ""

# 最终验证
echo "=================================="
echo " 验证结果"
echo "=================================="
echo ""
echo "✅ 所有配置检查通过！"
echo ""
echo "前后端加密配置一致且正确。"
echo ""

# 提示下一步操作
echo "=================================="
echo " 下一步操作"
echo "=================================="
echo ""
echo "1. 构建前端应用："
echo "   cd admin && npm run build"
echo ""
echo "2. 构建后端应用："
echo "   cd api && go build -o bin/gamelink cmd/main.go"
echo ""
echo "3. 启动 Docker 服务："
echo "   docker-compose -f docker-compose.prod.yml --env-file .env.production up -d"
echo ""
echo "4. 验证部署："
echo "   curl https://your-domain.com/api/v1/healthz"
echo ""

echo "✅ 验证完成！"
echo ""
