-- 给 admin@gamelink.com 用户分配 super_admin 角色

-- 1. 查看当前用户
SELECT '=== 检查用户 ===' as step;
SELECT id, email, name, role, status FROM users WHERE email = 'admin@gamelink.com';

-- 2. 查看 super_admin 角色 ID
SELECT '=== 检查 super_admin 角色 ===' as step;
SELECT id, slug, name FROM roles WHERE slug = 'super_admin';

-- 3. 给用户分配 super_admin 角色（假设 user_id=7, role_id=1）
-- 如果已存在会跳过
INSERT OR IGNORE INTO user_roles (user_id, role_id)
SELECT 
    (SELECT id FROM users WHERE email = 'admin@gamelink.com') as user_id,
    (SELECT id FROM roles WHERE slug = 'super_admin') as role_id;

-- 4. 验证分配结果
SELECT '=== 验证角色分配 ===' as step;
SELECT 
    u.id as user_id,
    u.email,
    r.id as role_id,
    r.slug,
    r.name
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
JOIN roles r ON ur.role_id = r.id
WHERE u.email = 'admin@gamelink.com';


