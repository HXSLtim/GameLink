BEGIN;

-- 1) Remove reviews attached to non-completed orders.
--    Rule baseline: reviews are only valid for completed orders.
WITH invalid_reviews AS (
    SELECT r.id
    FROM reviews r
    JOIN orders o ON o.id = r.order_id
    WHERE o.status <> 'completed'
)
DELETE FROM review_replies rr
USING invalid_reviews ir
WHERE rr.review_id = ir.id;

WITH invalid_reviews AS (
    SELECT r.id
    FROM reviews r
    JOIN orders o ON o.id = r.order_id
    WHERE o.status <> 'completed'
)
DELETE FROM review_reports rr
USING invalid_reviews ir
WHERE rr.review_id = ir.id;

WITH invalid_reviews AS (
    SELECT r.id
    FROM reviews r
    JOIN orders o ON o.id = r.order_id
    WHERE o.status <> 'completed'
)
DELETE FROM review_appeals ra
USING invalid_reviews ir
WHERE ra.review_id = ir.id;

WITH invalid_reviews AS (
    SELECT r.id
    FROM reviews r
    JOIN orders o ON o.id = r.order_id
    WHERE o.status <> 'completed'
)
DELETE FROM reviews r
USING invalid_reviews ir
WHERE r.id = ir.id;

-- 2) For completed orders that have no paid/refunded payment record,
--    backfill a synthetic paid payment for dev integrity.
WITH missing AS (
    SELECT o.id AS order_id, o.user_id, o.total_price_cents, o.currency, o.created_at
    FROM orders o
    WHERE o.status = 'completed'
      AND NOT EXISTS (
          SELECT 1
          FROM payments p
          WHERE p.order_id = o.id
            AND p.status IN ('paid', 'refunded')
      )
)
INSERT INTO payments (
    order_id,
    user_id,
    amount_cents,
    currency,
    status,
    method,
    provider_trade_no,
    provider_raw,
    paid_at,
    created_at,
    updated_at
)
SELECT
    m.order_id,
    m.user_id,
    m.total_price_cents,
    COALESCE(m.currency, 'CNY'),
    'paid',
    'wechat',
    CONCAT('FIX-', m.order_id, '-', TO_CHAR(NOW(), 'YYYYMMDDHH24MISS')),
    '{"source":"data_integrity_fix","note":"auto backfill for completed order"}',
    NOW(),
    COALESCE(m.created_at, NOW()),
    NOW()
FROM missing m;

-- 3) Remove orphan player schedules.
DELETE FROM player_schedules s
WHERE NOT EXISTS (
    SELECT 1
    FROM players p
    WHERE p.id = s.player_id
);

-- 4) Reconcile paid payments on canceled orders.
--    For development integrity, canceled orders with paid records are marked refunded.
UPDATE payments p
SET
    status = 'refunded',
    refunded_amount_cents = p.amount_cents,
    refunded_at = COALESCE(p.refunded_at, NOW()),
    updated_at = NOW()
FROM orders o
WHERE p.order_id = o.id
  AND o.status = 'canceled'
  AND p.status = 'paid';

UPDATE orders o
SET
    refund_amount_cents = CASE
        WHEN agg.total_refund > o.refund_amount_cents THEN agg.total_refund
        ELSE o.refund_amount_cents
    END,
    refunded_at = COALESCE(o.refunded_at, NOW()),
    refund_reason = CASE
        WHEN COALESCE(NULLIF(TRIM(o.refund_reason), ''), '') = '' THEN 'order canceled'
        ELSE o.refund_reason
    END,
    updated_at = NOW()
FROM (
    SELECT order_id, SUM(amount_cents) AS total_refund
    FROM payments
    WHERE status = 'refunded'
    GROUP BY order_id
) agg
WHERE o.id = agg.order_id
  AND o.status = 'canceled';

COMMIT;
