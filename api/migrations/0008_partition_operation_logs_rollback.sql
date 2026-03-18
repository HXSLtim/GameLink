-- ============================================================================
-- GameLink Table Partitioning Rollback - operation_logs
-- ============================================================================
-- Purpose: Rollback operation_logs partitioning to original non-partitioned table
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Usage:
--   psql -U gamelink -d gamelink -f 0008_partition_operation_logs_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- Drop partitioned table and all partitions
DROP TABLE IF EXISTS operation_logs CASCADE;

-- Restore original table
ALTER TABLE IF EXISTS operation_logs_old RENAME TO operation_logs;

-- Restore original indexes
ALTER INDEX IF EXISTS idx_operation_logs_user_id_old
    RENAME TO idx_operation_logs_user_id;
ALTER INDEX IF EXISTS idx_operation_logs_action_old
    RENAME TO idx_operation_logs_action;
ALTER INDEX IF EXISTS idx_operation_logs_trace_id_old
    RENAME TO idx_operation_logs_trace_id;
ALTER INDEX IF EXISTS idx_operation_logs_created_at_old
    RENAME TO idx_operation_logs_created_at;
ALTER INDEX IF EXISTS idx_operation_logs_ext_json_gin_old
    RENAME TO idx_operation_logs_ext_json_gin;
ALTER INDEX IF EXISTS idx_operation_logs_metadata_gin_old
    RENAME TO idx_operation_logs_metadata_gin;

-- Drop partition management functions
DROP FUNCTION IF EXISTS create_operation_logs_partition(date);
DROP FUNCTION IF EXISTS drop_old_operation_logs_partitions(int);

-- Drop views
DROP VIEW IF EXISTS v_recent_user_activity;
DROP VIEW IF EXISTS v_action_frequency;

-- Drop backup table (optional - keep for safety)
-- DROP TABLE IF EXISTS operation_logs_backup;

SELECT 'Rollback complete: operation_logs restored to non-partitioned table' AS status;

-- ============================================================================
-- Rollback Complete
-- ============================================================================
