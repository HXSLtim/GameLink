-- ============================================================================
-- GameLink PostgreSQL Slow Query Monitoring Setup
-- ============================================================================
-- Purpose: Enable pg_stat_statements extension for slow query monitoring
--          and performance analysis.
--
-- Author: Database Team
-- Date: 2026-03-02
-- Version: 1.0
--
-- Features:
-- - Track query execution statistics
-- - Identify slow queries and performance bottlenecks
-- - Monitor query frequency and resource usage
-- - Enable query optimization decisions
--
-- Prerequisites:
-- - PostgreSQL 9.2+ (pg_stat_statements is built-in)
-- - Superuser privileges required for initial setup
-- - Requires postgresql.conf modification (see below)
--
-- Migration Commands:
--   psql -U postgres -d gamelink -f 0006_enable_slow_query_monitoring.sql
-- ============================================================================

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';

-- ============================================================================
-- SECTION 1: Enable pg_stat_statements Extension
-- ============================================================================

-- Create extension if not exists (requires superuser)
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

COMMENT ON EXTENSION pg_stat_statements IS
'Track execution statistics of all SQL statements executed by the server.';

-- Verify extension is installed
SELECT extname, extversion
FROM pg_extension
WHERE extname = 'pg_stat_statements';

-- ============================================================================
-- SECTION 2: Configuration Parameters
-- ============================================================================
-- NOTE: These settings require postgresql.conf modification and server restart
-- Add to postgresql.conf:
--
-- shared_preload_libraries = 'pg_stat_statements'
-- pg_stat_statements.max = 10000
-- pg_stat_statements.track = all
-- pg_stat_statements.track_utility = on
-- pg_stat_statements.save = on
--
-- After modifying postgresql.conf, restart PostgreSQL:
--   sudo systemctl restart postgresql
-- ============================================================================

-- Display current configuration
SELECT name, setting, unit, context
FROM pg_settings
WHERE name LIKE 'pg_stat_statements%'
ORDER BY name;

-- ============================================================================
-- SECTION 3: Monitoring Views and Queries
-- ============================================================================

-- ----------------------------------------------------------------------------
-- View 1: Top 20 Slowest Queries by Total Time
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_slow_queries_by_total_time AS
SELECT
    query,
    calls,
    total_exec_time / 1000 AS total_time_sec,
    mean_exec_time / 1000 AS mean_time_sec,
    max_exec_time / 1000 AS max_time_sec,
    min_exec_time / 1000 AS min_time_sec,
    stddev_exec_time / 1000 AS stddev_time_sec,
    rows,
    100.0 * shared_blks_hit / NULLIF(shared_blks_hit + shared_blks_read, 0) AS cache_hit_ratio
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY total_exec_time DESC
LIMIT 20;

COMMENT ON VIEW v_slow_queries_by_total_time IS
'Top 20 queries by total execution time. Use this to find queries consuming most database time.';

-- ----------------------------------------------------------------------------
-- View 2: Top 20 Slowest Queries by Average Time
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_slow_queries_by_avg_time AS
SELECT
    query,
    calls,
    mean_exec_time / 1000 AS mean_time_sec,
    max_exec_time / 1000 AS max_time_sec,
    min_exec_time / 1000 AS min_time_sec,
    total_exec_time / 1000 AS total_time_sec,
    stddev_exec_time / 1000 AS stddev_time_sec,
    rows,
    100.0 * shared_blks_hit / NULLIF(shared_blks_hit + shared_blks_read, 0) AS cache_hit_ratio
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
  AND calls > 10  -- Filter out one-off queries
ORDER BY mean_exec_time DESC
LIMIT 20;

COMMENT ON VIEW v_slow_queries_by_avg_time IS
'Top 20 queries by average execution time (min 10 calls). Use this to find consistently slow queries.';

-- ----------------------------------------------------------------------------
-- View 3: Most Frequently Executed Queries
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_most_frequent_queries AS
SELECT
    query,
    calls,
    total_exec_time / 1000 AS total_time_sec,
    mean_exec_time / 1000 AS mean_time_sec,
    rows,
    100.0 * shared_blks_hit / NULLIF(shared_blks_hit + shared_blks_read, 0) AS cache_hit_ratio
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
ORDER BY calls DESC
LIMIT 20;

COMMENT ON VIEW v_most_frequent_queries IS
'Top 20 most frequently executed queries. Use this to identify hot paths for optimization.';

-- ----------------------------------------------------------------------------
-- View 4: Queries with Low Cache Hit Ratio
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_queries_low_cache_hit AS
SELECT
    query,
    calls,
    mean_exec_time / 1000 AS mean_time_sec,
    shared_blks_hit,
    shared_blks_read,
    100.0 * shared_blks_hit / NULLIF(shared_blks_hit + shared_blks_read, 0) AS cache_hit_ratio
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
  AND (shared_blks_hit + shared_blks_read) > 0
  AND calls > 10
ORDER BY cache_hit_ratio ASC
LIMIT 20;

COMMENT ON VIEW v_queries_low_cache_hit IS
'Queries with low cache hit ratio (<90%). These may benefit from index optimization or increased shared_buffers.';

-- ----------------------------------------------------------------------------
-- View 5: Queries with High I/O
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW v_queries_high_io AS
SELECT
    query,
    calls,
    mean_exec_time / 1000 AS mean_time_sec,
    shared_blks_read,
    shared_blks_written,
    shared_blks_dirtied,
    temp_blks_read,
    temp_blks_written
FROM pg_stat_statements
WHERE query NOT LIKE '%pg_stat_statements%'
  AND (shared_blks_read + shared_blks_written + temp_blks_read + temp_blks_written) > 0
ORDER BY (shared_blks_read + shared_blks_written + temp_blks_read + temp_blks_written) DESC
LIMIT 20;

COMMENT ON VIEW v_queries_high_io IS
'Queries with highest I/O activity. These may benefit from index optimization or query rewriting.';

-- ============================================================================
-- SECTION 4: Utility Functions
-- ============================================================================

-- Function to reset statistics (use with caution in production)
CREATE OR REPLACE FUNCTION reset_query_stats()
RETURNS void AS $$
BEGIN
    PERFORM pg_stat_statements_reset();
    RAISE NOTICE 'Query statistics have been reset.';
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION reset_query_stats() IS
'Reset all pg_stat_statements statistics. Use with caution in production.';

-- Function to get query statistics summary
CREATE OR REPLACE FUNCTION get_query_stats_summary()
RETURNS TABLE(
    total_queries bigint,
    total_calls bigint,
    total_time_sec numeric,
    avg_time_sec numeric,
    cache_hit_ratio numeric
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        COUNT(*)::bigint AS total_queries,
        SUM(calls)::bigint AS total_calls,
        ROUND((SUM(total_exec_time) / 1000)::numeric, 2) AS total_time_sec,
        ROUND((AVG(mean_exec_time) / 1000)::numeric, 2) AS avg_time_sec,
        ROUND((100.0 * SUM(shared_blks_hit) / NULLIF(SUM(shared_blks_hit + shared_blks_read), 0))::numeric, 2) AS cache_hit_ratio
    FROM pg_stat_statements
    WHERE query NOT LIKE '%pg_stat_statements%';
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_query_stats_summary() IS
'Get overall query statistics summary.';

-- ============================================================================
-- SECTION 5: Monitoring Queries (Examples)
-- ============================================================================

-- Example 1: Find queries taking longer than 1 second on average
-- SELECT * FROM v_slow_queries_by_avg_time WHERE mean_time_sec > 1;

-- Example 2: Find queries with cache hit ratio below 90%
-- SELECT * FROM v_queries_low_cache_hit WHERE cache_hit_ratio < 90;

-- Example 3: Get overall statistics summary
-- SELECT * FROM get_query_stats_summary();

-- Example 4: Find specific query pattern
-- SELECT query, calls, mean_exec_time / 1000 AS mean_time_sec
-- FROM pg_stat_statements
-- WHERE query ILIKE '%orders%'
-- ORDER BY mean_exec_time DESC
-- LIMIT 10;

-- ============================================================================
-- SECTION 6: Alerting Thresholds (Recommended)
-- ============================================================================
-- Set up monitoring alerts for:
-- 1. Queries with mean_exec_time > 1000ms (1 second)
-- 2. Queries with cache_hit_ratio < 90%
-- 3. Queries with temp_blks_written > 1000 (using temp files)
-- 4. Total database time increasing >20% week-over-week
-- ============================================================================

-- Display initial statistics
SELECT 'Slow Query Monitoring Enabled' AS status;
SELECT * FROM get_query_stats_summary();

-- ============================================================================
-- Migration Complete
-- ============================================================================
-- Next Steps:
-- 1. Modify postgresql.conf with recommended settings (see SECTION 2)
-- 2. Restart PostgreSQL server
-- 3. Monitor views regularly: v_slow_queries_by_total_time, v_slow_queries_by_avg_time
-- 4. Set up automated alerts for slow queries
-- 5. Review and optimize identified slow queries
-- 6. Schedule weekly review of query performance trends
--
-- Useful Commands:
-- - View slow queries: SELECT * FROM v_slow_queries_by_total_time;
-- - View cache issues: SELECT * FROM v_queries_low_cache_hit;
-- - Get summary: SELECT * FROM get_query_stats_summary();
-- - Reset stats: SELECT reset_query_stats(); (use with caution)
-- ============================================================================
