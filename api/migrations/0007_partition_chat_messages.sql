-- ============================================================================
-- GameLink Table Partitioning Migration - Part 1: chat_messages
-- ============================================================================
-- Purpose: Implement monthly range partitioning for chat_messages table
--          to improve query performance and enable efficient data archival.
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Benefits:
-- - Faster queries: Partition pruning reduces scan size
-- - Efficient archival: Drop old partitions instead of DELETE
-- - Better maintenance: VACUUM/ANALYZE per partition
-- - Aligned with retention policy: 30-day message retention
--
-- Strategy:
-- - Range partitioning by created_at (monthly)
-- - Create partitions for current month + 2 future months
-- - Automate partition creation via cron or pg_cron
--
-- Prerequisites:
-- - PostgreSQL 10+ (native declarative partitioning)
-- - Backup chat_messages table before migration
-- - Schedule during low-traffic window
--
-- Migration Commands:
--   psql -U gamelink -d gamelink -f 0007_partition_chat_messages.sql
--
-- Rollback:
--   See 0007_partition_chat_messages_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ============================================================================
-- SECTION 1: Backup Existing Table
-- ============================================================================

-- Create backup table (optional but recommended)
CREATE TABLE IF NOT EXISTS chat_messages_backup AS
SELECT * FROM chat_messages;

COMMENT ON TABLE chat_messages_backup IS
'Backup of chat_messages before partitioning migration (2026-03-02)';

-- ============================================================================
-- SECTION 2: Rename Existing Table
-- ============================================================================

-- Rename existing table to _old suffix
ALTER TABLE IF EXISTS chat_messages RENAME TO chat_messages_old;

-- Rename indexes to avoid conflicts
ALTER INDEX IF EXISTS idx_chat_messages_group_sent_covering
    RENAME TO idx_chat_messages_group_sent_covering_old;
ALTER INDEX IF EXISTS idx_chat_messages_sender_id
    RENAME TO idx_chat_messages_sender_id_old;
ALTER INDEX IF EXISTS idx_chat_messages_moderated_by
    RENAME TO idx_chat_messages_moderated_by_old;
ALTER INDEX IF EXISTS idx_chat_messages_ext_json_gin
    RENAME TO idx_chat_messages_ext_json_gin_old;
ALTER INDEX IF EXISTS idx_chat_messages_metadata_gin
    RENAME TO idx_chat_messages_metadata_gin_old;

-- ============================================================================
-- SECTION 3: Create Partitioned Table
-- ============================================================================

CREATE TABLE chat_messages (
    id bigserial NOT NULL,
    group_id bigint NOT NULL,
    sender_id bigint NOT NULL,
    content text NOT NULL,
    message_type varchar(16) DEFAULT 'text',
    reply_to_id bigint,
    image_url varchar(255),
    metadata json DEFAULT '{}',
    is_deleted boolean DEFAULT false,
    audit_status varchar(16) DEFAULT 'pending',
    moderated_by bigint,
    moderated_at timestamp with time zone,
    reject_reason text,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    ext_json jsonb DEFAULT '{}'::jsonb,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

COMMENT ON TABLE chat_messages IS
'Partitioned chat messages table (monthly range partitioning by created_at)';

-- ============================================================================
-- SECTION 4: Create Partitions
-- ============================================================================

-- Create partitions for past 2 months, current month, and next 2 months
-- Adjust date ranges based on migration date

-- January 2026
CREATE TABLE IF NOT EXISTS chat_messages_2026_01 PARTITION OF chat_messages
FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- February 2026
CREATE TABLE IF NOT EXISTS chat_messages_2026_02 PARTITION OF chat_messages
FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- March 2026 (current month)
CREATE TABLE IF NOT EXISTS chat_messages_2026_03 PARTITION OF chat_messages
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

-- April 2026
CREATE TABLE IF NOT EXISTS chat_messages_2026_04 PARTITION OF chat_messages
FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

-- May 2026
CREATE TABLE IF NOT EXISTS chat_messages_2026_05 PARTITION OF chat_messages
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');

-- Default partition for future data (catch-all)
CREATE TABLE IF NOT EXISTS chat_messages_default PARTITION OF chat_messages DEFAULT;

-- ============================================================================
-- SECTION 5: Recreate Indexes on Partitioned Table
-- ============================================================================

-- Covering index for message list queries
CREATE INDEX idx_chat_messages_group_sent_covering
ON chat_messages (group_id, created_at DESC)
INCLUDE (id, content, sender_id, message_type, audit_status);

-- Foreign key indexes
CREATE INDEX idx_chat_messages_sender_id ON chat_messages (sender_id);
CREATE INDEX idx_chat_messages_moderated_by ON chat_messages (moderated_by)
WHERE moderated_by IS NOT NULL;

-- Audit status index
CREATE INDEX idx_chat_messages_audit_status ON chat_messages (audit_status);

-- Soft delete index
CREATE INDEX idx_chat_messages_deleted_at ON chat_messages (deleted_at);

-- JSONB GIN indexes
CREATE INDEX idx_chat_messages_ext_json_gin
ON chat_messages USING GIN (ext_json jsonb_path_ops);

CREATE INDEX idx_chat_messages_metadata_gin
ON chat_messages USING GIN (CAST(metadata AS jsonb) jsonb_path_ops);

-- ============================================================================
-- SECTION 6: Migrate Data from Old Table
-- ============================================================================

-- Insert data from old table to partitioned table
-- This may take time depending on data volume
-- Consider using pg_bulkload or COPY for large datasets

INSERT INTO chat_messages
SELECT * FROM chat_messages_old;

-- Verify row count matches
DO $$
DECLARE
    old_count bigint;
    new_count bigint;
BEGIN
    SELECT COUNT(*) INTO old_count FROM chat_messages_old;
    SELECT COUNT(*) INTO new_count FROM chat_messages;

    IF old_count = new_count THEN
        RAISE NOTICE 'Data migration successful: % rows migrated', new_count;
    ELSE
        RAISE EXCEPTION 'Data migration failed: old=%, new=%', old_count, new_count;
    END IF;
END $$;

-- ============================================================================
-- SECTION 7: Recreate Foreign Key Constraints
-- ============================================================================

-- Add foreign key to chat_groups
ALTER TABLE chat_messages
ADD CONSTRAINT fk_chat_messages_group
FOREIGN KEY (group_id) REFERENCES chat_groups(id)
ON UPDATE CASCADE ON DELETE CASCADE;

-- Note: Foreign key to users table for sender_id
-- Uncomment if users table exists and constraint is needed:
-- ALTER TABLE chat_messages
-- ADD CONSTRAINT fk_chat_messages_sender
-- FOREIGN KEY (sender_id) REFERENCES users(id)
-- ON UPDATE CASCADE ON DELETE RESTRICT;

-- ============================================================================
-- SECTION 8: Create Partition Management Function
-- ============================================================================

CREATE OR REPLACE FUNCTION create_chat_messages_partition(partition_date date)
RETURNS void AS $$
DECLARE
    partition_name text;
    start_date date;
    end_date date;
BEGIN
    -- Calculate partition boundaries
    start_date := date_trunc('month', partition_date);
    end_date := start_date + interval '1 month';
    partition_name := 'chat_messages_' || to_char(start_date, 'YYYY_MM');

    -- Create partition if not exists
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF chat_messages FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        start_date,
        end_date
    );

    RAISE NOTICE 'Created partition: % for range [%, %)', partition_name, start_date, end_date;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION create_chat_messages_partition(date) IS
'Create a monthly partition for chat_messages table. Call with first day of target month.';

-- ============================================================================
-- SECTION 9: Create Partition Cleanup Function
-- ============================================================================

CREATE OR REPLACE FUNCTION drop_old_chat_messages_partitions(retention_months int DEFAULT 2)
RETURNS void AS $$
DECLARE
    partition_record record;
    cutoff_date date;
BEGIN
    cutoff_date := date_trunc('month', CURRENT_DATE) - (retention_months || ' months')::interval;

    FOR partition_record IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
          AND tablename LIKE 'chat_messages_20%'
          AND tablename != 'chat_messages_default'
          AND tablename != 'chat_messages_old'
          AND tablename != 'chat_messages_backup'
        ORDER BY tablename
    LOOP
        -- Extract date from partition name (format: chat_messages_YYYY_MM)
        DECLARE
            partition_date date;
        BEGIN
            partition_date := to_date(substring(partition_record.tablename from 15), 'YYYY_MM');

            IF partition_date < cutoff_date THEN
                EXECUTE format('DROP TABLE IF EXISTS %I', partition_record.tablename);
                RAISE NOTICE 'Dropped old partition: %', partition_record.tablename;
            END IF;
        EXCEPTION WHEN OTHERS THEN
            RAISE WARNING 'Failed to process partition %: %', partition_record.tablename, SQLERRM;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION drop_old_chat_messages_partitions(int) IS
'Drop chat_messages partitions older than retention_months (default 2 months).';

-- ============================================================================
-- SECTION 10: Verification
-- ============================================================================

-- List all partitions
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'chat_messages%'
ORDER BY tablename;

-- Verify partition constraints
SELECT
    tablename,
    pg_get_expr(relpartbound, oid) AS partition_expression
FROM pg_class
JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
WHERE relname LIKE 'chat_messages_%'
  AND relkind = 'r'
ORDER BY relname;

-- Test query with partition pruning
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM chat_messages
WHERE created_at >= '2026-03-01' AND created_at < '2026-04-01'
LIMIT 10;

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- Next Steps:
-- 1. Test application queries to ensure compatibility
-- 2. Monitor query performance with EXPLAIN ANALYZE
-- 3. Set up monthly cron job to create future partitions:
--    SELECT create_chat_messages_partition(CURRENT_DATE + interval '2 months');
-- 4. Set up monthly cron job to drop old partitions:
--    SELECT drop_old_chat_messages_partitions(2);
-- 5. After verification (1-2 weeks), drop old table:
--    DROP TABLE chat_messages_old;
-- 6. Update application code if needed (should be transparent)
--
-- Maintenance Schedule (Recommended):
-- - Weekly: ANALYZE chat_messages;
-- - Monthly: Create next 2 months' partitions
-- - Monthly: Drop partitions older than retention period
-- - Quarterly: Review partition sizes and adjust strategy
-- ============================================================================
