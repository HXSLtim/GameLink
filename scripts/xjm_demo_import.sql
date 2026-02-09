BEGIN;

-- Helper mapping tables for idempotent imports.
CREATE TABLE IF NOT EXISTS xjm_import.user_map (
    source_type TEXT NOT NULL,
    source_key TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (source_type, source_key)
);

CREATE TABLE IF NOT EXISTS xjm_import.player_map (
    source_key TEXT PRIMARY KEY,
    player_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS xjm_import.item_map (
    projectname TEXT PRIMARY KEY,
    item_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS xjm_import.order_map (
    source_order_id BIGINT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    order_item_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS xjm_import_ordersettlements_orderid_idx
    ON xjm_import.ordersettlements (orderid);

-- Demo game category.
INSERT INTO game_categories (name, description, is_active, sort_order)
SELECT 'xjm_demo', 'XJM demo category', TRUE, 0
WHERE NOT EXISTS (SELECT 1 FROM game_categories WHERE name = 'xjm_demo');

-- Games derived from project types.
WITH game_src AS (
    SELECT DISTINCT TRIM(SPLIT_PART(projecttype, '/', 1)) AS game_name
    FROM xjm_import.theprojectdisk
    WHERE projecttype IS NOT NULL AND TRIM(projecttype) <> ''
)
INSERT INTO games (created_at, updated_at, key, name, category, category_id, is_active, sort_order, ext_json)
SELECT
    NOW(),
    NOW(),
    'xjm_' || SUBSTR(MD5(game_name), 1, 24),
    LEFT(game_name, 128),
    'xjm_demo',
    (SELECT id FROM game_categories WHERE name = 'xjm_demo'),
    TRUE,
    0,
    JSONB_BUILD_OBJECT('source', 'xjm')
FROM game_src
WHERE NOT EXISTS (SELECT 1 FROM games g WHERE g.name = game_name);

-- Service items from project disk.
WITH src AS (
    SELECT
        p.*,
        COALESCE(NULLIF(TRIM(SPLIT_PART(projecttype, '/', 2)), ''),
                 NULLIF(TRIM(projecttype), ''),
                 'custom') AS sub_category,
        COALESCE(NULLIF(TRIM(SPLIT_PART(projecttype, '/', 1)), ''),
                 NULLIF(TRIM(projecttype), ''),
                 NULL) AS game_name,
        COALESCE(NULLIF(TRIM(projectname), ''),
                 NULLIF(TRIM(projecttype), ''),
                 'Unnamed') AS name_key
    FROM xjm_import.theprojectdisk p
)
INSERT INTO service_items (
    item_code,
    name,
    description,
    category,
    sub_category,
    game_id,
    category_id,
    base_price_cents,
    service_hours,
    icon_url,
    is_active,
    created_at,
    updated_at
)
SELECT
    'xjm_' || p.id::text,
    LEFT(p.name_key, 128),
    p.projecttext,
    'escort',
    LEFT(p.sub_category, 32),
    g.id,
    (SELECT id FROM game_categories WHERE name = 'xjm_demo'),
    COALESCE(ROUND(COALESCE(p.projectamount, p.singleamount, 0) * 100)::bigint, 0),
    1,
    p.projectpictureurl,
    TRUE,
    p.createtime,
    COALESCE(p.updatetime, p.createtime)
FROM src p
LEFT JOIN games g ON g.name = p.game_name
WHERE NOT EXISTS (
    SELECT 1 FROM service_items s WHERE s.item_code = 'xjm_' || p.id::text
);

-- Map project names to service items.
WITH src AS (
    SELECT
        p.id,
        COALESCE(NULLIF(TRIM(projectname), ''),
                 NULLIF(TRIM(projecttype), ''),
                 'Unnamed') AS name_key
    FROM xjm_import.theprojectdisk p
)
INSERT INTO xjm_import.item_map (projectname, item_id)
SELECT DISTINCT ON (s.name_key)
    s.name_key,
    si.id
FROM src s
JOIN service_items si ON si.item_code = 'xjm_' || s.id::text
WHERE s.name_key IS NOT NULL
ON CONFLICT DO NOTHING;

-- Service items for order-only project names.
WITH proj AS (
    SELECT
        COALESCE(NULLIF(TRIM(projectname), ''), 'Unnamed') AS name_key,
        AVG(orderamount) AS avg_amount,
        MIN(createtime) AS min_time
    FROM xjm_import.orders
    GROUP BY COALESCE(NULLIF(TRIM(projectname), ''), 'Unnamed')
),
missing AS (
    SELECT p.*
    FROM proj p
    WHERE NOT EXISTS (SELECT 1 FROM xjm_import.item_map m WHERE m.projectname = p.name_key)
)
INSERT INTO service_items (
    item_code,
    name,
    description,
    category,
    sub_category,
    game_id,
    category_id,
    base_price_cents,
    service_hours,
    icon_url,
    is_active,
    created_at,
    updated_at
)
SELECT
    'xjm_auto_' || SUBSTR(MD5(name_key), 1, 20),
    LEFT(name_key, 128),
    NULL,
    'escort',
    'custom',
    NULL,
    (SELECT id FROM game_categories WHERE name = 'xjm_demo'),
    COALESCE(ROUND(avg_amount * 100)::bigint, 0),
    1,
    NULL,
    TRUE,
    COALESCE(min_time, NOW()),
    COALESCE(min_time, NOW())
FROM missing
WHERE NOT EXISTS (
    SELECT 1 FROM service_items s WHERE s.item_code = 'xjm_auto_' || SUBSTR(MD5(name_key), 1, 20)
);

INSERT INTO xjm_import.item_map (projectname, item_id)
SELECT
    p.name_key,
    si.id
FROM (
    SELECT DISTINCT COALESCE(NULLIF(TRIM(projectname), ''), 'Unnamed') AS name_key
    FROM xjm_import.orders
) p
JOIN service_items si ON si.item_code = 'xjm_auto_' || SUBSTR(MD5(p.name_key), 1, 20)
ON CONFLICT DO NOTHING;

-- Boss users.
WITH boss_src AS (
    SELECT DISTINCT COALESCE(NULLIF(TRIM(bossname), ''), 'anonymous') AS boss_name
    FROM xjm_import.orders
    UNION
    SELECT DISTINCT COALESCE(NULLIF(TRIM(bossname), ''), 'anonymous') AS boss_name
    FROM xjm_import.rechargerecord
)
INSERT INTO users (created_at, updated_at, name, nickname, email, status, ext_json)
SELECT
    NOW(),
    NOW(),
    LEFT(b.boss_name, 64),
    LEFT(b.boss_name, 64),
    'xjm_boss_' || SUBSTR(MD5(b.boss_name), 1, 20) || '@demo.local',
    'active',
    JSONB_BUILD_OBJECT('source', 'xjm', 'type', 'boss', 'key', b.boss_name)
FROM boss_src b
WHERE NOT EXISTS (
    SELECT 1 FROM users u
    WHERE u.email = 'xjm_boss_' || SUBSTR(MD5(b.boss_name), 1, 20) || '@demo.local'
);

INSERT INTO xjm_import.user_map (source_type, source_key, user_id)
SELECT
    'boss',
    b.boss_name,
    u.id
FROM (
    SELECT DISTINCT COALESCE(NULLIF(TRIM(bossname), ''), 'anonymous') AS boss_name
    FROM xjm_import.orders
    UNION
    SELECT DISTINCT COALESCE(NULLIF(TRIM(bossname), ''), 'anonymous') AS boss_name
    FROM xjm_import.rechargerecord
) b
JOIN users u ON u.email = 'xjm_boss_' || SUBSTR(MD5(b.boss_name), 1, 20) || '@demo.local'
ON CONFLICT DO NOTHING;

-- Handler users.
WITH handler_src AS (
    SELECT DISTINCT ON (handlerid)
        handlerid::text AS handler_key,
        handlergamename
    FROM xjm_import.ordersettlements
    WHERE handlerid IS NOT NULL
    ORDER BY handlerid, handlergamename
)
INSERT INTO users (created_at, updated_at, name, nickname, email, status, ext_json)
SELECT
    NOW(),
    NOW(),
    LEFT(COALESCE(NULLIF(TRIM(handlergamename), ''), 'handler_' || handler_key), 64),
    LEFT(COALESCE(NULLIF(TRIM(handlergamename), ''), 'handler_' || handler_key), 64),
    'xjm_handler_' || handler_key || '@demo.local',
    'active',
    JSONB_BUILD_OBJECT('source', 'xjm', 'type', 'handler', 'key', handler_key)
FROM handler_src h
WHERE NOT EXISTS (
    SELECT 1 FROM users u WHERE u.email = 'xjm_handler_' || h.handler_key || '@demo.local'
);

INSERT INTO xjm_import.user_map (source_type, source_key, user_id)
SELECT
    'handler',
    h.handler_key,
    u.id
FROM (
    SELECT DISTINCT handlerid::text AS handler_key
    FROM xjm_import.ordersettlements
    WHERE handlerid IS NOT NULL
) h
JOIN users u ON u.email = 'xjm_handler_' || h.handler_key || '@demo.local'
ON CONFLICT DO NOTHING;

-- Players for handlers.
INSERT INTO players (
    created_at,
    updated_at,
    user_id,
    nickname,
    verification_status,
    accepting_orders,
    online_status,
    ext_json
)
SELECT
    NOW(),
    NOW(),
    u.id,
    LEFT(u.name, 64),
    'verified',
    TRUE,
    'online',
    JSONB_BUILD_OBJECT('source', 'xjm', 'handler_key', m.source_key)
FROM xjm_import.user_map m
JOIN users u ON u.id = m.user_id
WHERE m.source_type = 'handler'
  AND NOT EXISTS (SELECT 1 FROM players p WHERE p.user_id = u.id);

INSERT INTO xjm_import.player_map (source_key, player_id)
SELECT
    m.source_key,
    p.id
FROM xjm_import.user_map m
JOIN players p ON p.user_id = m.user_id
WHERE m.source_type = 'handler'
ON CONFLICT DO NOTHING;

-- Orders.
WITH order_src AS (
    SELECT
        o.*,
        COALESCE(NULLIF(TRIM(o.bossname), ''), 'anonymous') AS boss_key,
        COALESCE(NULLIF(TRIM(o.projectname), ''), 'Unnamed') AS project_key
    FROM xjm_import.orders o
),
mapped AS (
    SELECT
        o.*,
        u.user_id,
        im.item_id
    FROM order_src o
    JOIN xjm_import.user_map u
      ON u.source_type = 'boss' AND u.source_key = o.boss_key
    JOIN xjm_import.item_map im
      ON im.projectname = o.project_key
),
settlement_stats AS (
    SELECT
        orderid,
        MIN(starttime) AS starttime,
        MAX(endtime) FILTER (WHERE osstatus IN ('APPROVED', 'CSAPPROVED')) AS endtime,
        BOOL_OR(osstatus IN ('APPROVED', 'CSAPPROVED')) AS has_approved,
        COUNT(*) > 0 AS has_any
    FROM xjm_import.ordersettlements
    GROUP BY orderid
),
status_calc AS (
    SELECT
        m.*,
        CASE
            WHEN ss.has_approved THEN 'completed'
            WHEN ss.has_any THEN 'processing'
            ELSE 'pending'
        END AS final_status,
        ss.starttime,
        ss.endtime
    FROM mapped m
    LEFT JOIN settlement_stats ss ON ss.orderid = m.id
)
INSERT INTO orders (
    created_at,
    updated_at,
    ext_json,
    order_no,
    user_id,
    item_id,
    player_id,
    quantity,
    unit_price_cents,
    total_price_cents,
    commission_cents,
    player_income_cents,
    status,
    title,
    description,
    game_id,
    scheduled_start,
    scheduled_end,
    started_at,
    completed_at
)
SELECT
    sc.createtime,
    COALESCE(sc.updatetime, sc.createtime),
    JSONB_BUILD_OBJECT(
        'source', 'xjm',
        'source_id', sc.id,
        'outside_order_id', sc.outsideorderid,
        'spare_flag', sc.spareflag,
        'is_deposits', sc.isdeposits
    ),
    'xjm_' || sc.id::text,
    sc.user_id,
    sc.item_id,
    pm.player_id,
    1,
    COALESCE(ROUND(sc.orderamount * 100)::bigint, 0),
    COALESCE(ROUND(sc.orderamount * 100)::bigint, 0),
    0,
    0,
    sc.final_status,
    LEFT(sc.projectname, 128),
    sc.remake,
    NULL,
    sc.starttime,
    sc.endtime,
    sc.starttime,
    sc.endtime
FROM status_calc sc
LEFT JOIN xjm_import.player_map pm ON pm.source_key = sc.handler1::text
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.order_no = 'xjm_' || sc.id::text
);

INSERT INTO xjm_import.order_map (source_order_id, order_id)
SELECT
    (SUBSTR(o.order_no, 5))::bigint AS source_order_id,
    o.id
FROM orders o
WHERE o.order_no LIKE 'xjm_%'
ON CONFLICT DO NOTHING;

-- Order items.
INSERT INTO order_items (
    created_at,
    updated_at,
    order_id,
    item_id,
    slot,
    unit_price_cents,
    quantity,
    total_cents,
    commission_rate,
    status,
    player_id,
    ext_json
)
SELECT
    o.created_at,
    o.updated_at,
    o.id,
    o.item_id,
    1,
    o.unit_price_cents,
    o.quantity,
    o.total_price_cents,
    0.20,
    o.status,
    o.player_id,
    o.ext_json
FROM orders o
WHERE o.order_no LIKE 'xjm_%'
  AND NOT EXISTS (
    SELECT 1 FROM order_items oi WHERE oi.order_id = o.id AND oi.slot = 1
  );

UPDATE xjm_import.order_map m
SET order_item_id = oi.id
FROM orders o
JOIN order_items oi ON oi.order_id = o.id AND oi.slot = 1
WHERE o.order_no = 'xjm_' || m.source_order_id::text
  AND m.order_item_id IS NULL;

-- Order players from settlements.
WITH settle AS (
    SELECT
        s.*,
        o.id AS order_id,
        oi.id AS order_item_id,
        pm.player_id,
        COALESCE(s.starttime, o.created_at) AS joined_at,
        CASE
            WHEN s.osstatus IN ('APPROVED', 'CSAPPROVED') THEN 'completed'
            WHEN s.osstatus IN ('REJECTED') THEN 'rejected'
            ELSE 'joined'
        END AS final_status
    FROM xjm_import.ordersettlements s
    JOIN orders o ON o.order_no = 'xjm_' || s.orderid::text
    JOIN order_items oi ON oi.order_id = o.id AND oi.slot = 1
    JOIN xjm_import.player_map pm ON pm.source_key = s.handlerid::text
)
INSERT INTO order_players (
    created_at,
    updated_at,
    order_id,
    order_item_id,
    player_id,
    joined_at,
    income_cents,
    commission_cents,
    status,
    ext_json
)
SELECT
    NOW(),
    NOW(),
    order_id,
    order_item_id,
    player_id,
    joined_at,
    COALESCE(ROUND(handleramount * 100)::bigint, 0),
    COALESCE(ROUND(commission * 100)::bigint, 0),
    final_status,
    JSONB_BUILD_OBJECT('source', 'xjm', 'settlement_id', id)
FROM settle
ON CONFLICT (order_item_id, player_id) DO NOTHING;

-- Order players from assigned handlers.
WITH assigned AS (
    SELECT id AS source_order_id, handler1 AS handler_id
    FROM xjm_import.orders
    WHERE handler1 IS NOT NULL
    UNION
    SELECT id AS source_order_id, handler2 AS handler_id
    FROM xjm_import.orders
    WHERE handler2 IS NOT NULL
),
mapped AS (
    SELECT
        a.source_order_id,
        o.id AS order_id,
        oi.id AS order_item_id,
        pm.player_id
    FROM assigned a
    JOIN orders o ON o.order_no = 'xjm_' || a.source_order_id::text
    JOIN order_items oi ON oi.order_id = o.id AND oi.slot = 1
    JOIN xjm_import.player_map pm ON pm.source_key = a.handler_id::text
)
INSERT INTO order_players (
    created_at,
    updated_at,
    order_id,
    order_item_id,
    player_id,
    joined_at,
    income_cents,
    commission_cents,
    status,
    ext_json
)
SELECT
    NOW(),
    NOW(),
    order_id,
    order_item_id,
    player_id,
    NOW(),
    0,
    0,
    'joined',
    JSONB_BUILD_OBJECT('source', 'xjm', 'assigned', true)
FROM mapped
ON CONFLICT (order_item_id, player_id) DO NOTHING;

-- Backfill order/player links.
UPDATE orders o
SET player_id = op.player_id
FROM (
    SELECT o.id, op.player_id,
           ROW_NUMBER() OVER (PARTITION BY o.id ORDER BY op.id) AS rn
    FROM orders o
    JOIN order_players op ON op.order_id = o.id
    WHERE o.order_no LIKE 'xjm_%'
) op
WHERE o.id = op.id AND o.player_id IS NULL AND op.rn = 1;

UPDATE order_items oi
SET player_id = o.player_id
FROM orders o
WHERE oi.order_id = o.id AND oi.player_id IS NULL AND o.player_id IS NOT NULL;

-- Recharge records.
WITH src AS (
    SELECT
        r.*,
        COALESCE(NULLIF(TRIM(r.bossname), ''), 'anonymous') AS boss_key
    FROM xjm_import.rechargerecord r
),
mapped AS (
    SELECT s.*, u.user_id
    FROM src s
    JOIN xjm_import.user_map u
      ON u.source_type = 'boss' AND u.source_key = s.boss_key
)
INSERT INTO recharge_records (
    created_at,
    updated_at,
    user_id,
    amount_cents,
    bonus_cents,
    total_cents,
    status,
    order_no,
    payment_channel,
    payment_method,
    paid_at,
    remark,
    ext_json
)
SELECT
    s.createtime,
    COALESCE(s.updatetime, s.createtime),
    s.user_id,
    COALESCE(ROUND(s.rechargenum * 100)::bigint, 0),
    CASE WHEN s.isgift THEN COALESCE(ROUND(s.rechargenum * 100)::bigint, 0) ELSE 0 END,
    COALESCE(ROUND(s.rechargenum * 100)::bigint, 0)
      + CASE WHEN s.isgift THEN COALESCE(ROUND(s.rechargenum * 100)::bigint, 0) ELSE 0 END,
    'success',
    'xjm_recharge_' || s.id::text,
    'manual',
    'manual',
    s.rechargetime,
    s.remark,
    JSONB_BUILD_OBJECT('source', 'xjm', 'picurl', s.picurl, 'reporter', s.reporter)
FROM mapped s
WHERE NOT EXISTS (
    SELECT 1 FROM recharge_records rr WHERE rr.order_no = 'xjm_recharge_' || s.id::text
);

-- Update demo users total recharge and wallets.
WITH totals AS (
    SELECT user_id, SUM(total_cents) AS total_cents
    FROM recharge_records
    GROUP BY user_id
)
UPDATE users u
SET total_recharge_cents = t.total_cents
FROM totals t
WHERE u.id = t.user_id;

INSERT INTO wallets (created_at, updated_at, user_id, balance_cents)
SELECT
    NOW(),
    NOW(),
    u.id,
    COALESCE(t.total_cents, 0)
FROM users u
LEFT JOIN (
    SELECT user_id, SUM(total_cents) AS total_cents
    FROM recharge_records
    GROUP BY user_id
) t ON t.user_id = u.id
WHERE u.email LIKE 'xjm_%'
  AND NOT EXISTS (SELECT 1 FROM wallets w WHERE w.user_id = u.id);

-- Reviews from bad reviews.
WITH src AS (
    SELECT
        br.*,
        o.id AS order_id,
        oi.id AS order_item_id,
        o.user_id,
        o.player_id
    FROM xjm_import.badreviews br
    JOIN orders o ON o.order_no = 'xjm_' || br.relatedorderid::text
    JOIN order_items oi ON oi.order_id = o.id AND oi.slot = 1
)
INSERT INTO reviews (
    created_at,
    updated_at,
    order_id,
    order_item_id,
    user_id,
    player_id,
    score,
    content,
    status,
    is_reported,
    images,
    is_public,
    is_anonymous,
    edit_count,
    last_edit_at,
    ext_json
)
SELECT
    s.createtime,
    COALESCE(s.updatetime, s.createtime),
    s.order_id,
    s.order_item_id,
    s.user_id,
    s.player_id,
    1,
    s.content,
    'approved',
    FALSE,
    '[]'::json,
    TRUE,
    FALSE,
    0,
    s.updatetime,
    JSONB_BUILD_OBJECT('source', 'xjm', 'source_id', s.id)
FROM src s
WHERE s.content IS NOT NULL AND s.content <> ''
  AND NOT EXISTS (
    SELECT 1 FROM reviews r
    WHERE r.order_id = s.order_id AND r.content = s.content
  );

UPDATE order_items oi
SET review_id = r.id
FROM (
    SELECT DISTINCT ON (order_item_id) id, order_item_id
    FROM reviews
    ORDER BY order_item_id, created_at DESC
) r
WHERE oi.id = r.order_item_id AND oi.review_id IS NULL;

COMMIT;
