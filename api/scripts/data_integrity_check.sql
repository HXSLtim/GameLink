\pset format aligned
\pset tuples_only off

\echo '=== Core Counts ==='
SELECT 'users' AS metric, COUNT(*) AS value FROM users
UNION ALL SELECT 'players', COUNT(*) FROM players
UNION ALL SELECT 'orders', COUNT(*) FROM orders
UNION ALL SELECT 'payments', COUNT(*) FROM payments
UNION ALL SELECT 'reviews', COUNT(*) FROM reviews
ORDER BY metric;

\echo ''
\echo '=== Orphan Integrity ==='
SELECT 'orders.user_id -> users.id' AS check_item, COUNT(*) AS violations
FROM orders o LEFT JOIN users u ON u.id = o.user_id
WHERE u.id IS NULL
UNION ALL
SELECT 'orders.player_id -> players.id', COUNT(*)
FROM orders o LEFT JOIN players p ON p.id = o.player_id
WHERE o.player_id IS NOT NULL AND p.id IS NULL
UNION ALL
SELECT 'payments.order_id -> orders.id', COUNT(*)
FROM payments p LEFT JOIN orders o ON o.id = p.order_id
WHERE o.id IS NULL
UNION ALL
SELECT 'payments.user_id -> users.id', COUNT(*)
FROM payments p LEFT JOIN users u ON u.id = p.user_id
WHERE u.id IS NULL
UNION ALL
SELECT 'reviews.order_id -> orders.id', COUNT(*)
FROM reviews r LEFT JOIN orders o ON o.id = r.order_id
WHERE o.id IS NULL
UNION ALL
SELECT 'reviews.user_id -> users.id', COUNT(*)
FROM reviews r LEFT JOIN users u ON u.id = r.user_id
WHERE u.id IS NULL
UNION ALL
SELECT 'reviews.player_id -> players.id', COUNT(*)
FROM reviews r LEFT JOIN players p ON p.id = r.player_id
WHERE p.id IS NULL
UNION ALL
SELECT 'player_schedules.player_id -> players.id', COUNT(*)
FROM player_schedules s LEFT JOIN players p ON p.id = s.player_id
WHERE p.id IS NULL;

\echo ''
\echo '=== Business Consistency ==='
SELECT 'reviews on non-completed orders' AS check_item, COUNT(*) AS violations
FROM reviews r
JOIN orders o ON o.id = r.order_id
WHERE o.status <> 'completed'
UNION ALL
SELECT 'multi pending payments on same order', COUNT(*)
FROM (
    SELECT order_id
    FROM payments
    WHERE status = 'pending'
    GROUP BY order_id
    HAVING COUNT(*) > 1
) x
UNION ALL
SELECT 'completed orders without paid/refunded payment', COUNT(*)
FROM orders o
WHERE o.status = 'completed'
  AND NOT EXISTS (
      SELECT 1
      FROM payments p
      WHERE p.order_id = o.id
        AND p.status IN ('paid', 'refunded')
  )
UNION ALL
SELECT 'paid payments on canceled orders', COUNT(*)
FROM payments p
JOIN orders o ON o.id = p.order_id
WHERE p.status = 'paid'
  AND o.status = 'canceled';

\echo ''
\echo '=== Sample Violations (Top 20) ==='
SELECT r.id AS review_id, r.order_id, o.status AS order_status, r.score
FROM reviews r
JOIN orders o ON o.id = r.order_id
WHERE o.status <> 'completed'
ORDER BY r.id
LIMIT 20;

SELECT o.id AS order_id, o.status, o.total_price_cents
FROM orders o
WHERE o.status = 'completed'
  AND NOT EXISTS (
      SELECT 1
      FROM payments p
      WHERE p.order_id = o.id
        AND p.status IN ('paid', 'refunded')
  )
ORDER BY o.id
LIMIT 20;
