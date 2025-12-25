-- 插入排行榜抽成配置
INSERT INTO ranking_commission_configs (name, ranking_type, period, month, rules_json, description, is_active, created_at, updated_at)
SELECT '月度收入排行抽成', 'income', 'monthly', TO_CHAR(NOW(), 'YYYY-MM'), 
       '[{"rankStart":1,"rankEnd":3,"commissionRate":10},{"rankStart":4,"rankEnd":10,"commissionRate":12},{"rankStart":11,"rankEnd":50,"commissionRate":15},{"rankStart":51,"rankEnd":100,"commissionRate":18}]',
       '本月收入排行榜抽成规则：TOP3享10%低抽成', true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM ranking_commission_configs WHERE ranking_type = 'income');

INSERT INTO ranking_commission_configs (name, ranking_type, period, month, rules_json, description, is_active, created_at, updated_at)
SELECT '月度订单量排行抽成', 'order_count', 'monthly', TO_CHAR(NOW(), 'YYYY-MM'),
       '[{"rankStart":1,"rankEnd":5,"commissionRate":8},{"rankStart":6,"rankEnd":20,"commissionRate":12},{"rankStart":21,"rankEnd":100,"commissionRate":16}]',
       '本月订单量排行榜抽成规则：TOP5享8%超低抽成', true, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM ranking_commission_configs WHERE ranking_type = 'order_count');

-- 插入订单纠纷数据
INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, sla_deadline, created_at, updated_at)
SELECT 1, 1, 'pending', '陪玩师迟到', '约定时间是晚上8点，但陪玩师9点才上线', 'pending', NOW() + INTERVAL '30 minutes', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 1);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, assigned_to_user_id, assignment_source, sla_deadline, created_at, updated_at)
SELECT 2, 3, 'assigned', '服务态度差', '陪玩师在游戏中频繁挂机，态度敷衍', 'pending', 16, 'manual', NOW() + INTERVAL '30 minutes', NOW() - INTERVAL '4 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 2);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, resolution_amount, resolution_notes, resolved_at, assigned_to_user_id, sla_deadline, created_at, updated_at)
SELECT 3, 4, 'resolved', '陪玩师中途离开', '游戏进行到一半，陪玩师突然下线', 'refund', 9900, '经核实，全额退款', NOW() - INTERVAL '2 hours', 16, NOW() - INTERVAL '24 hours', NOW() - INTERVAL '24 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 3);

-- 补充收款主体数据（如果不存在）
INSERT INTO collection_entities (name, credit_code, tax_registration_no, status, is_default, total_collection_cents, transaction_count, created_by, created_at, updated_at)
SELECT '游戏联盟科技有限公司', '91110108MA01ABCD1X', '110108MA01ABCD1X', 'active', true, 1580000, 156, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM collection_entities WHERE credit_code = '91110108MA01ABCD1X');

INSERT INTO collection_entities (name, credit_code, tax_registration_no, status, is_default, total_collection_cents, transaction_count, created_by, created_at, updated_at)
SELECT '星耀互娱网络科技公司', '91310115MA1KEFGH2Y', '310115MA1KEFGH2Y', 'active', false, 890000, 89, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM collection_entities WHERE credit_code = '91310115MA1KEFGH2Y');

-- 补充结算公司数据（如果不存在）
INSERT INTO settlement_companies (name, credit_code, bank_name, bank_account, contact_name, contact_phone, status, total_payout_cents, player_count, created_by, created_at, updated_at)
SELECT '游戏联盟结算中心', '91110108MA01WXYZ1A', '中国工商银行北京分行', '6222021234567890123', '张经理', '13912345678', 'active', 2580000, 12, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM settlement_companies WHERE credit_code = '91110108MA01WXYZ1A');

INSERT INTO settlement_companies (name, credit_code, bank_name, bank_account, contact_name, contact_phone, status, total_payout_cents, player_count, created_by, created_at, updated_at)
SELECT '星耀支付结算公司', '91310115MA1KMNOP2B', '招商银行上海分行', '6225881234567890456', '李总监', '13898765432', 'active', 1890000, 8, 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM settlement_companies WHERE credit_code = '91310115MA1KMNOP2B');

-- 补充分流规则数据（如果不存在）
INSERT INTO routing_rules (name, priority, conditions, target_entity_id, status, description, created_by, created_at, updated_at)
SELECT '大额订单分流', 1, '[{"field":"order_amount","operator":"gt","value":500}]'::json, 
       (SELECT id FROM collection_entities WHERE is_default = true LIMIT 1), 
       'active', '订单金额超过500元走主收款主体', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM routing_rules WHERE name = '大额订单分流');
