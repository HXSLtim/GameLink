-- ============================================================================
-- GameLink JSONB GIN Index Rollback
-- ============================================================================
-- Purpose: Remove GIN indexes created by 0005_add_jsonb_gin_indexes.sql
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Usage:
--   psql -U gamelink -d gamelink -f 0005_add_jsonb_gin_indexes_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ============================================================================
-- Drop all GIN indexes created in migration 0005
-- ============================================================================

-- Base model ext_json indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_users_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_players_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_payments_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_groups_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_messages_ext_json_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_reviews_ext_json_gin;

-- Business-specific JSON column indexes
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_config_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_groups_settings_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_messages_metadata_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_payments_provider_raw_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_operation_logs_metadata_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_financial_reports_data_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_settings_notifications_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_settings_privacy_gin;
DROP INDEX CONCURRENTLY IF EXISTS idx_vip_levels_benefits_gin;

-- ============================================================================
-- Rollback Complete
-- ============================================================================
