-- Migration: Enable pg_stat_statements extension for slow query monitoring
-- Date: 2026-01-10
-- Description: Enables PostgreSQL's pg_stat_statements extension for query performance monitoring

-- Enable the extension (requires superuser privileges)
-- Note: This extension must be added to shared_preload_libraries in postgresql.conf
-- and requires a server restart to take effect.
-- 
-- In postgresql.conf:
--   shared_preload_libraries = 'pg_stat_statements'
--   pg_stat_statements.track = all
--   pg_stat_statements.max = 10000

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Useful queries for monitoring:

-- Top 20 slowest queries by mean execution time
-- SELECT query, calls, total_exec_time, mean_exec_time, max_exec_time, rows
-- FROM pg_stat_statements
-- ORDER BY mean_exec_time DESC
-- LIMIT 20;

-- Most frequently called queries
-- SELECT query, calls, total_exec_time, mean_exec_time
-- FROM pg_stat_statements
-- ORDER BY calls DESC
-- LIMIT 20;

-- Queries with highest total execution time
-- SELECT query, calls, total_exec_time, mean_exec_time
-- FROM pg_stat_statements
-- ORDER BY total_exec_time DESC
-- LIMIT 20;

-- Reset statistics (use with caution)
-- SELECT pg_stat_statements_reset();
