-- ============================================================================
-- GameLink Covering Index Migration - Rollback Script
-- ============================================================================
-- Purpose: Rollback covering index migration and restore original indexes.
--
-- WARNING: This will remove optimized indexes and may impact query performance.
--
-- Rollback Commands:
--   psql -U gamelink -d gamelink -f 0001_add_covering_indexes_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ----------------------------------------------------------------------------
-- Rollback: Order Table
-- ----------------------------------------------------------------------------

-- Drop covering index
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_user_status_created_covering;

-- Recreate original partial index (without INCLUDE)
CREATE INDEX CONCURRENTLY idx_orders_user_status_created
ON orders (user_id, status, created_at DESC);

-- ----------------------------------------------------------------------------
-- Rollback: Chat Messages Table
-- ----------------------------------------------------------------------------

-- Drop covering index
DROP INDEX CONCURRENTLY IF EXISTS idx_chat_messages_group_sent_covering;

-- Recreate original index (without INCLUDE)
CREATE INDEX CONCURRENTLY idx_chat_messages_group_created
ON chat_messages (group_id, created_at DESC);

-- ----------------------------------------------------------------------------
-- Rollback: Payment Table (Bonus)
-- ----------------------------------------------------------------------------

-- Drop covering index
DROP INDEX CONCURRENTLY IF EXISTS idx_payments_user_status_created_covering;

-- ============================================================================
-- Rollback Complete
-- ============================================================================
