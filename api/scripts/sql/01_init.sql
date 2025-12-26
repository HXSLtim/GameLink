-- GameLink 数据库初始化脚本
-- 此脚本在 PostgreSQL 容器首次启动时自动执行

-- 设置数据库编码
SET client_encoding = 'UTF8';

-- 创建扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- 用于模糊搜索

-- 创建数据库元信息表
CREATE TABLE IF NOT EXISTS db_version (
    version VARCHAR(50) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- 记录初始化版本
INSERT INTO db_version (version, description) 
VALUES ('1.0.0', 'Initial database setup')
ON CONFLICT (version) DO NOTHING;

-- 创建索引优化查询性能的函数
CREATE OR REPLACE FUNCTION create_index_if_not_exists(
    index_name TEXT,
    table_name TEXT,
    column_name TEXT
) RETURNS VOID AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = index_name
    ) THEN
        EXECUTE format('CREATE INDEX %I ON %I (%I)', 
            index_name, table_name, column_name);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- 设置默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO gamelink;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT ON SEQUENCES TO gamelink;

-- 输出初始化完成信息
DO $$
BEGIN
    RAISE NOTICE 'GameLink database initialized successfully';
END $$;
