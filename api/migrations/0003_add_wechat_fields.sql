-- Migration: Add WeChat fields to users table
-- Description: Support WeChat mini-program login
-- Date: 2026-01-10

-- Add WeChat OpenID field (unique per mini-program)
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_open_id VARCHAR(64);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_wechat_open_id ON users(wechat_open_id) WHERE wechat_open_id IS NOT NULL;

-- Add WeChat UnionID field (shared across apps under same WeChat Open Platform account)
ALTER TABLE users ADD COLUMN IF NOT EXISTS wechat_union_id VARCHAR(64);
CREATE INDEX IF NOT EXISTS idx_users_wechat_union_id ON users(wechat_union_id) WHERE wechat_union_id IS NOT NULL;

-- Comment
COMMENT ON COLUMN users.wechat_open_id IS '微信小程序 OpenID，每个小程序唯一';
COMMENT ON COLUMN users.wechat_union_id IS '微信 UnionID，同一开放平台下的应用共享';
