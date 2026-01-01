-- ============================================================================
-- GameLink Covering Index Verification Script
-- ============================================================================
-- Purpose: Verify covering index performance and monitor index usage.
--
-- Usage:
--   psql -U gamelink -d gamelink -f verify_covering_indexes.sql
--
-- Schedule: Run daily or weekly to monitor index effectiveness
-- ============================================================================

\echo '============================================================================'
\echo 'Covering Index Verification Report'
\echo 'Generated: ' `date`
\echo '============================================================================'
\echo ''

-- ----------------------------------------------------------------------------
-- 1. Index Size Comparison
-- ----------------------------------------------------------------------------
\echo '1. INDEX SIZE COMPARISON'
\echo '--------------------------------------------------------------------'
\echo 'Comparing covering index size vs regular index size:'
\echo ''

SELECT
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size,
    indexdef
FROM pg_indexes
WHERE indexname LIKE '%_covering%'
ORDER BY tablename, indexname;

\echo ''

-- ----------------------------------------------------------------------------
-- 2. Index Usage Statistics
-- ----------------------------------------------------------------------------
\echo '2. INDEX USAGE STATISTICS'
\echo '--------------------------------------------------------------------'
\echo 'Monitor how often covering indexes are used vs heap fetches:'
\echo ''

SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS index_scans,
    idx_tup_read AS tuples_read,
    idx_tup_fetch AS tuples_fetched,
    ROUND(100.0 * idx_tup_fetch / NULLIF(idx_tup_read, 0), 2) AS fetch_ratio_pct,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE indexname LIKE '%_covering%'
   OR tablename IN ('orders', 'chat_messages', 'payments')
ORDER BY tablename, idx_scan DESC;

\echo ''
\echo 'Key Metrics:'
\echo '  - idx_scan: Total index scans (should be high for active indexes)'
\echo '  - idx_tup_fetch: Heap fetches (should be LOW for covering indexes)'
\echo '  - fetch_ratio_pct: Lower is better (covering indexes should be < 10%)'
\echo ''

-- ----------------------------------------------------------------------------
-- 3. Query Performance Test (Order List Query)
-- ----------------------------------------------------------------------------
\echo '3. QUERY PERFORMANCE TEST'
\echo '--------------------------------------------------------------------'
\echo 'Testing order list query with EXPLAIN ANALYZE:'
\echo ''
\echo 'Note: Run this on a production-like dataset for accurate results'
\echo ''

-- Prepare test function
CREATE OR REPLACE FUNCTION test_order_list_query(user_id_param BIGINT)
RETURNS TABLE(query_plan TEXT) AS $$
BEGIN
    RETURN QUERY
    EXECUTE format(
        'EXPLAIN (ANALYZE, BUFFERS, VERBOSE) SELECT id, player_id, total_price_cents
         FROM orders
         WHERE user_id = %L AND status IN (''pending'', ''confirmed'', ''completed'')
         ORDER BY created_at DESC LIMIT 20',
        user_id_param
    );
END;
$$ LANGUAGE plpgsql;

-- Run test (replace with actual user_id)
-- SELECT * FROM test_order_list_query(1);

\echo 'To test manually, run:'
\echo '  EXPLAIN (ANALYZE, BUFFERS, VERBOSE)'
\echo '  SELECT id, player_id, total_price_cents'
\echo '  FROM orders'
\echo '  WHERE user_id = <actual_user_id> AND status IN (''pending'', ''confirmed'', ''completed'')'
\echo '  ORDER BY created_at DESC LIMIT 20;'
\echo ''
\echo 'Look for: "Index Only Scan" (good) vs "Index Scan using..." (heap fetch)'
\echo ''

-- ----------------------------------------------------------------------------
-- 4. Chat Message Query Performance Test
-- ----------------------------------------------------------------------------
\echo '4. CHAT MESSAGE QUERY TEST'
\echo '--------------------------------------------------------------------'
\echo 'Testing chat message list query:'
\echo ''

\echo 'To test manually, run:'
\echo '  EXPLAIN (ANALYZE, BUFFERS, VERBOSE)'
\echo '  SELECT id, content, sender_id'
\echo '  FROM chat_messages'
\echo '  WHERE group_id = <actual_group_id>'
\echo '  ORDER BY created_at DESC LIMIT 50;'
\echo ''
\echo 'Look for: "Index Only Scan" (good)'
\echo ''

-- ----------------------------------------------------------------------------
-- 5. Index Bloat Check (requires pgstattuple extension)
-- ----------------------------------------------------------------------------
\echo '5. INDEX BLOAT CHECK'
\echo '--------------------------------------------------------------------'

-- Check if pgstattuple is available
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pgstattuple') THEN
        RAISE NOTICE 'pgstattuple extension is available';
    ELSE
        RAISE NOTICE 'pgstattuple extension not found. To install: CREATE EXTENSION pgstattuple;';
    END IF;
END $$;

\echo ''

-- If available, check bloat
-- SELECT * FROM pgstatindex('idx_orders_user_status_created_covering');

\echo 'To check index bloat (requires pgstattuple):'
\echo '  SELECT * FROM pgstatindex(''idx_orders_user_status_created_covering'');'
\echo '  SELECT * FROM pgstatindex(''idx_chat_messages_group_sent_covering'');'
\echo ''
\echo 'Bloat > 30%: Consider REINDEX CONCURRENTLY'
\echo ''

-- ----------------------------------------------------------------------------
-- 6. Missing Indexes Analysis
-- ----------------------------------------------------------------------------
\echo '6. POTENTIAL MISSING COVERING INDEXES'
\echo '--------------------------------------------------------------------'
\echo 'High-fetch queries that might benefit from covering indexes:'
\echo ''

SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_fetch,
    ROUND(100.0 * idx_tup_fetch / NULLIF(idx_scan, 0), 2) AS avg_fetches_per_scan
FROM pg_stat_user_indexes
WHERE idx_scan > 1000
  AND idx_tup_fetch > 0
  AND (100.0 * idx_tup_fetch / NULLIF(idx_scan, 0)) > 50
ORDER BY avg_fetches_per_scan DESC
LIMIT 10;

\echo ''
\echo 'Queries with high fetch ratios indicate potential covering index opportunities'
\echo ''

-- ----------------------------------------------------------------------------
-- 7. Recommendation Summary
-- ----------------------------------------------------------------------------
\echo '============================================================================'
\echo 'RECOMMENDATION SUMMARY'
\echo '============================================================================'
\echo ''
\echo 'Covering Index Best Practices:'
\echo '  ✓ Monitor index usage weekly'
\echo '  ✓ REINDEX when bloat > 30% or quarterly'
\echo '  ✓ Use EXPLAIN ANALYZE to verify Index Only Scan'
\echo '  ✓ Balance index size vs query performance'
\echo '  ✓ Remove unused indexes to reduce write overhead'
\echo ''
\echo 'Performance Targets:'
\echo '  - fetch_ratio_pct: < 10% (excellent), 10-30% (good), > 50% (needs review)'
\echo '  - Index size: < 50% of table size'
\echo '  - Bloat: < 30%'
\echo ''
\echo '============================================================================'

-- Cleanup
DROP FUNCTION IF EXISTS test_order_list_query(BIGINT);
