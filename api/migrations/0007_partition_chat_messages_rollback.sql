-- ============================================================================
-- GameLink Table Partitioning Rollback - chat_messages
-- ============================================================================
-- Purpose: Rollback chat_messages partitioning to original non-partitioned table
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Usage:
--   psql -U gamelink -d gamelink -f 0007_partition_chat_messages_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- Drop partitioned table and all partitions
DROP TABLE IF EXISTS chat_messages CASCADE;

-- Restore original table
ALTER TABLE IF EXISTS chat_messages_old RENAME TO chat_messages;

-- Restore original indexes
ALTER INDEX IF EXISTS idx_chat_messages_group_sent_covering_old
    RENAME TO idx_chat_messages_group_sent_covering;
ALTER INDEX IF EXISTS idx_chat_messages_sender_id_old
    RENAME TO idx_chat_messages_sender_id;
ALTER INDEX IF EXISTS idx_chat_messages_moderated_by_old
    RENAME TO idx_chat_messages_moderated_by;
ALTER INDEX IF EXISTS idx_chat_messages_ext_json_gin_old
    RENAME TO idx_chat_messages_ext_json_gin;
ALTER INDEX IF EXISTS idx_chat_messages_metadata_gin_old
    RENAME TO idx_chat_messages_metadata_gin;

-- Drop partition management functions
DROP FUNCTION IF EXISTS create_chat_messages_partition(date);
DROP FUNCTION IF EXISTS drop_old_chat_messages_partitions(int);

-- Drop backup table (optional - keep for safety)
-- DROP TABLE IF EXISTS chat_messages_backup;

SELECT 'Rollback complete: chat_messages restored to non-partitioned table' AS status;

-- ============================================================================
-- Rollback Complete
-- ============================================================================
