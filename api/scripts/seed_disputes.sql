-- 插入订单纠纷数据（使用实际存在的 order_id 和 user_id）
INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, sla_deadline, created_at, updated_at)
SELECT 111, 482, 'pending', '陪玩师迟到', '约定时间是晚上8点，但陪玩师9点才上线，严重影响游戏体验', 'pending', NOW() + INTERVAL '30 minutes', NOW() - INTERVAL '2 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 111);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, assigned_to_user_id, assignment_source, sla_deadline, created_at, updated_at)
SELECT 112, 484, 'assigned', '服务态度差', '陪玩师在游戏中频繁挂机，态度敷衍', 'pending', 487, 'manual', NOW() + INTERVAL '30 minutes', NOW() - INTERVAL '4 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 112);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, assigned_to_user_id, assignment_source, sla_deadline, sla_breached, created_at, updated_at)
SELECT 113, 485, 'mediating', '技术水平不符', '陪玩师声称是王者段位，实际游戏表现只有黄金水平', 'pending', 487, 'manual', NOW() - INTERVAL '6 hours', true, NOW() - INTERVAL '8 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 113);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, resolution_amount, resolution_notes, resolved_at, resolved_by_user_id, assigned_to_user_id, sla_deadline, created_at, updated_at)
SELECT 114, 484, 'resolved', '陪玩师中途离开', '游戏进行到一半，陪玩师突然下线不再回来', 'refund', 9900, '经核实，陪玩师确实中途离开，全额退款', NOW() - INTERVAL '20 hours', 487, 487, NOW() - INTERVAL '24 hours', NOW() - INTERVAL '24 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 114);

INSERT INTO order_disputes (order_id, user_id, status, reason, description, resolution, resolution_amount, resolution_notes, resolved_at, resolved_by_user_id, assigned_to_user_id, sla_deadline, created_at, updated_at)
SELECT 115, 488, 'resolved', '服务时长不足', '购买了2小时服务，实际只陪玩了1.5小时', 'partial', 5000, '按实际服务时长计算，退还0.5小时费用', NOW() - INTERVAL '44 hours', 487, 487, NOW() - INTERVAL '48 hours', NOW() - INTERVAL '48 hours', NOW()
WHERE NOT EXISTS (SELECT 1 FROM order_disputes WHERE order_id = 115);
