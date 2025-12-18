INSERT INTO menus (name, path, component, icon, parent_id, "order", hidden) VALUES
('纠纷管理', '/admin/biz/dispute', 'Dispute', 'warning', 7, 50, false),
('用户标签', '/admin/sys/user-tag', 'UserTag', 'tags', 2, 50, false),
('结算公司', '/admin/finance/settlement-company', 'SettlementCompany', 'bank', 12, 30, false),
('排行榜抽成', '/admin/finance/ranking-commission', 'RankingCommission', 'trophy', 12, 40, false),
('分流规则', '/admin/biz/routing-rule', 'RoutingRule', 'fork', 7, 60, false),
('用户行为', '/admin/monitor/user-behavior', 'UserBehavior', 'line-chart', 15, 40, false),
('提现分流', '/admin/finance/withdraw-routing', 'WithdrawRouting', 'swap', 12, 50, false);
