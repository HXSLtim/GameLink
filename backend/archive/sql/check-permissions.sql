-- 检查超级管理员账户和角色分配

-- 1. 查看所有角色
SELECT '=== 角色列表 ===' as info;
SELECT id, name, slug, description FROM roles WHERE deleted_at IS NULL;

-- 2. 查看超级管理员用户
SELECT '=== 超级管理员用户 ===' as info;
SELECT id, name, email, role, status FROM users WHERE email = 'admin@gamelink.local' AND deleted_at IS NULL;

-- 3. 查看用户-角色关系
SELECT '=== 用户角色关系 ===' as info;
SELECT 
    ur.user_id,
    u.email,
    ur.role_id,
    r.slug as role_slug,
    r.name as role_name
FROM user_roles ur
LEFT JOIN users u ON ur.user_id = u.id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.email = 'admin@gamelink.local';

-- 4. 检查 super_admin 角色的权限数量
SELECT '=== super_admin 角色权限数量 ===' as info;
SELECT 
    r.id as role_id,
    r.slug,
    COUNT(rp.permission_id) as permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
WHERE r.slug = 'super_admin'
GROUP BY r.id, r.slug;

-- 5. 查看前 10 个权限
SELECT '=== 权限示例（前10个）===' as info;
SELECT id, code, method, path, description 
FROM permissions 
WHERE deleted_at IS NULL 
ORDER BY id 
LIMIT 10;


