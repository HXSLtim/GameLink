-- ============================================================================
-- GameLink Table Partitioning Migration - Part 2: operation_logs
-- ============================================================================
-- Purpose: Implement monthly range partitioning for operation_logs table
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
-- - Reduced table bloat: Smaller partitions are easier to maintain
--
-- Strategy:
-- - Range partitioning by created_at (monthly)
-- - Create partitions for current month + 2 future months
-- - Automate partition creation via cron or pg_cron
-- - Recommended retention: 6-12 months for audit logs
--
-- Prerequisites:
-- - PostgreSQL 10+ (native declarative partitioning)
-- - Backup operation_logs table before migration
-- - Schedule during low-traffic window
--
-- Migration Commands:
--   psql -U gamelink -d gamelink -f 0008_partition_operation_logs.sql
--
-- Rollback:
--   See 0008_partition_operation_logs_rollback.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ============================================================================
-- SECTION 1: Backup Existing Table
-- ============================================================================

-- Create backup table (optional but recommended)
CREATE TABLE IF NOT EXISTS operation_logs_backup AS
SELECT * FROM operation_logs;

COMMENT ON TABLE operation_logs_backup IS
'Backup of operation_logs before partitioning migration (2026-03-02)';

-- ============================================================================
-- SECTION 2: Rename Existing Table
-- ============================================================================

-- Rename existing table to _old suffix
ALTER TABLE IF EXISTS operation_logs RENAME TO operation_logs_old;

-- Rename indexes to avoid conflicts
ALTER INDEX IF EXISTS idx_operation_logs_user_id
    RENAME TO idx_operation_logs_user_id_old;
ALTER INDEX IF EXISTS idx_operation_logs_action
    RENAME TO idx_operation_logs_action_old;
ALTER INDEX IF EXISTS idx_operation_logs_trace_id
    RENAME TO idx_operation_logs_trace_id_old;
ALTER INDEX IF EXISTS idx_operation_logs_created_at
    RENAME TO idx_operation_logs_created_at_old;
ALTER INDEX IF EXISTS idx_operation_logs_ext_json_gin
    RENAME TO idx_operation_logs_ext_json_gin_old;
ALTER INDEX IF EXISTS idx_operation_logs_metadata_gin
    RENAME TO idx_operation_logs_metadata_gin_old;

-- ============================================================================
-- SECTION 3: Create Partitioned Table
-- ============================================================================

CREATE TABLE operation_logs (
    id bigserial NOT NULL,
    user_id bigint,
    action varchar(128) NOT NULL,
    resource_type varchar(64),
    resource_id bigint,
    ip_address varchar(64),
    user_agent text,
    trace_id varchar(64),
    metadata_json jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    ext_json jsonb DEFAULT '{}'::jsonb,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

COMMENT ON TABLE operation_logs IS
'Partitioned operation logs table (monthly range partitioning by created_at)';

-- ============================================================================
-- SECTION 4: Create Partitions
-- ============================================================================

-- Create partitions for past 2 months, current month, and next 2 months
-- Adjust date ranges based on migration date

-- January 2026
CREATE TABLE IF NOT EXISTS operation_logs_2026_01 PARTITION OF operation_logs
FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- February 2026
CREATE TABLE IF NOT EXISTS operation_logs_2026_02 PARTITION OF operation_logs
FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- March 2026 (current month)
CREATE TABLE IF NOT EXISTS operation_logs_2026_03 PARTITION OF operation_logs
FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

-- April 2026
CREATE TABLE IF NOT EXISTS operation_logs_2026_04 PARTITION OF operation_logs
FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');

-- May 2026
CREATE TABLE IF NOT EXISTS operation_logs_2026_05 PARTITION OF operation_logs
FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');

-- Default partition for future data (catch-all)
CREATE TABLE IF NOT EXISTS operation_logs_default PARTITION OF operation_logs DEFAULT;

-- ============================================================================
-- SECTION 5: Recreate Indexes on Partitioned Table
-- ============================================================================

-- User activity index
CREATE INDEX idx_operation_logs_user_id ON operation_logs (user_id, created_at DESC);

-- Action type index
CREATE INDEX idx_operation_logs_action ON operation_logs (action, created_at DESC);

-- Trace ID for distributed tracing
CREATE INDEX idx_operation_logs_trace_id ON operation_logs (trace_id)
WHERE trace_id IS NOT NULL;

-- Resource lookup index
CREATE INDEX idx_operation_logs_resource ON operation_logs (resource_type, resource_id, created_at DESC)
WHERE resource_type IS NOT NULL AND resource_id IS NOT NULL;

-- Soft delete index
CREATE INDEX idx_operation_logs_deleted_at ON operation_logs (deleted_at);

-- JSONB GIN indexes
CREATE INDEX idx_operation_logs_ext_json_gin
ON operation_logs USING GIN (ext_json jsonb_path_ops);

CREATE INDEX idx_operation_logs_metadata_gin
ON operation_logs USING GIN (metadata_json jsonb_path_ops);

-- Composite index for audit queries
CREATE INDEX idx_operation_logs_audit ON operation_logs (user_id, action, created_at DESC);

-- ============================================================================
-- SECTION 6: Migrate Data from Old Table
-- ============================================================================

-- Insert data from old table to partitioned table
-- This may take time depending on data volume
-- Consider batching for very large tables

INSERT INTO operation_logs
SELECT * FROM operation_logs_old;

-- Verify row count matches
DO $$
DECLARE
    old_count bigint;
    new_count bigint;
BEGIN
    SELECT COUNT(*) INTO old_count FROM operation_logs_old;
    SELECT COUNT(*) INTO new_count FROM operation_logs;

    IF old_count = new_count THEN
        RAISE NOTICE 'Data migration successful: % rows migrated', new_count;
    ELSE
        RAISE EXCEPTION 'Data migration failed: old=%, new=%', old_count, new_count;
    END IF;
END $$;

-- ============================================================================
-- SECTION 7: Create Partition Management Function
-- ============================================================================

CREATE OR REPLACE FUNCTION create_operation_logs_partition(partition_date date)
RETURNS void AS $$
DECLARE
    partition_name text;
    start_date date;
    end_date date;
BEGIN
    -- Calculate partition boundaries
    start_date := date_trunc('month', partition_date);
    end_date := start_date + interval '1 month';
    partition_name := 'operation_logs_' || to_char(start_date, 'YYYY_MM');

    -- Create partition if not exists
    EXECUTE format(
        'CREATE TABLE IF NOT EXISTS %I PARTITION OF operation_logs FOR VALUES FROM (%L) TO (%L)',
        partition_name,
        start_date,
        end_date
    );

    RAISE NOTICE 'Created partition: % for range [%, %)', partition_name, start_date, end_date;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION create_operation_logs_partition(date) IS
'Create a monthly partition for operation_logs table. Call with first day of target month.';

-- ============================================================================
-- SECTION 8: Create Partition Cleanup Function
-- ============================================================================

CREATE OR REPLACE FUNCTION drop_old_operation_logs_partitions(retention_months int DEFAULT 6)
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
          AND tablename LIKE 'operation_logs_20%'
          AND tablename != 'operation_logs_default'
          AND tablename != 'operation_logs_old'
          AND tablename != 'operation_logs_backup'
        ORDER BY tablename
    LOOP
        -- Extract date from partition name (format: operation_logs_YYYY_MM)
        DECLARE
            partition_date date;
        BEGIN
            partition_date := to_date(substring(partition_record.tablename from 16), 'YYYY_MM');

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

COMMENT ON FUNCTION drop_old_operation_logs_partitions(int) IS
'Drop operation_logs partitions older than retention_months (default 6 months for audit compliance).';

-- ============================================================================
-- SECTION 9: Create Audit Query Helper Views
-- ============================================================================

-- View: Recent user activity (last 30 days)
CREATE OR REPLACE VIEW v_recent_user_activity AS
SELECT
    user_id,
    action,
    resource_type,
    resource_id,
    ip_address,
    created_at
FROM operation_logs
WHERE created_at >= CURRENT_DATE - interval '30 days'
  AND user_id IS NOT NULL
ORDER BY created_at DESC;

COMMENT ON VIEW v_recent_user_activity IS
'Recent user activity logs (last 30 days). Useful for user behavior analysis.';

-- View: Action frequency summary
CREATE OR REPLACE VIEW v_action_frequency AS
SELECT
    action,
    COUNT(*) AS total_count,
    COUNT(DISTINCT user_id) AS unique_users,
    MIN(created_at) AS first_seen,
    MAX(created_at) AS last_seen
FROM operation_logs
WHERE created_at >= CURRENT_DATE - interval '7 days'
GROUP BY action
ORDER BY total_count DESC;

COMMENT ON VIEW v_action_frequency IS
'Action frequency summary (last 7 days). Useful for monitoring system usage patterns.';

-- ============================================================================
-- SECTION 10: Verification
-- ============================================================================

-- List all partitions
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE tablename LIKE 'operation_logs%'
ORDER BY tablename;

-- Verify partition constraints
SELECT
    tablename,
    pg_get_expr(relpartbound, oid) AS partition_expression
FROM pg_class
JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
WHERE relname LIKE 'operation_logs_%'
  AND relkind = 'r'
ORDER BY relname;

-- Test query with partition pruning
EXPLAIN (ANALYZE, BUFFERS)
SELECT * FROM operation_logs
WHERE created_at >= '2026-03-01' AND created_at < '2026-04-01'
  AND action = 'user.login'
LIMIT 10;

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- Next Steps:
-- 1. Test application queries to ensure compatibility
-- 2. Monitor query performance with EXPLAIN ANALYZE
-- 3. Set up monthly cron job to create future partitions:
--    SELECT create_operation_logs_partition(CURRENT_DATE + interval '2 months');
-- 4. Set up monthly cron job to drop old partitions (6-12 month retention):
--    SELECT drop_old_operation_logs_partitions(6);
-- 5. After verification (1-2 weeks), drop old table:
--    DROP TABLE operation_logs_old;
-- 6. Update application code if needed (should be transparent)
--
-- Maintenance Schedule (Recommended):
-- - Weekly: ANALYZE operation_logs;
-- - Monthly: Create next 2 months' partitions
-- - Monthly: Drop partitions older than retention period (6 months)
-- - Quarterly: Review partition sizes and adjust strategy
-- - Annually: Review retention policy for compliance requirements
--
-- Compliance Note:
-- - Audit logs typically require 6-12 month retention
-- - Adjust retention_months parameter based on regulatory requirements
-- - Consider archiving to cold storage before dropping partitions
-- ============================================================================
