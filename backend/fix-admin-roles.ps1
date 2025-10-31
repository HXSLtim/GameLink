# 修复管理员角色分配

Write-Host "`n=== Fixing Admin Roles ===" -ForegroundColor Cyan

# 查找数据库文件
$dbPath = Get-ChildItem -Path . -Recurse -Filter "*.db" | Select-Object -First 1 -ExpandProperty FullName

if (-not $dbPath) {
    Write-Host "ERROR: SQLite database not found" -ForegroundColor Red
    Write-Host "Please make sure the service has been run at least once" -ForegroundColor Yellow
    exit 1
}

Write-Host "Database found: $dbPath" -ForegroundColor Green

# 执行 SQL 修复
Write-Host "`nExecuting fix..." -ForegroundColor Yellow

$sql = @"
-- 给所有 role='admin' 的用户分配 super_admin 角色
INSERT OR IGNORE INTO user_roles (user_id, role_id)
SELECT 
    u.id as user_id,
    (SELECT id FROM roles WHERE slug = 'super_admin') as role_id
FROM users u
WHERE u.role = 'admin' 
AND u.deleted_at IS NULL
AND NOT EXISTS (
    SELECT 1 FROM user_roles ur 
    WHERE ur.user_id = u.id 
    AND ur.role_id = (SELECT id FROM roles WHERE slug = 'super_admin')
);

-- 显示结果
SELECT 
    u.id,
    u.email,
    u.role as old_role,
    GROUP_CONCAT(r.slug, ', ') as rbac_roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.role = 'admin' AND u.deleted_at IS NULL
GROUP BY u.id, u.email, u.role;
"@

# 使用 sqlite3 执行
try {
    $result = & sqlite3 $dbPath $sql 2>&1
    
    Write-Host "`nFixed successfully!" -ForegroundColor Green
    Write-Host "`nAdmin users with super_admin role:" -ForegroundColor Cyan
    Write-Host $result
    
    Write-Host "`n=== Next Steps ===" -ForegroundColor Cyan
    Write-Host "1. Restart the backend service" -ForegroundColor Yellow
    Write-Host "2. Re-login with your admin account" -ForegroundColor Yellow
    Write-Host "3. The 403 error should be resolved" -ForegroundColor Yellow
    
} catch {
    Write-Host "`nERROR: Failed to execute SQL" -ForegroundColor Red
    Write-Host "Manual fix required. Run this SQL:" -ForegroundColor Yellow
    Write-Host $sql -ForegroundColor Gray
}

