package integration

import (
"encoding/csv"
"encoding/json"
"fmt"
"net/http"
"strings"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
"gorm.io/gorm"

adminhandler "gamelink/internal/handler/admin"
"gamelink/internal/handler/middleware"
"gamelink/internal/model"
adminrepo "gamelink/internal/repository/admin"
"gamelink/internal/repository/permission"
"gamelink/internal/repository/permissionauditlog"
"gamelink/internal/repository/user"
permissionservice "gamelink/internal/service/admin"
roleservice "gamelink/internal/service/admin"
"gamelink/internal/service/audit"
"gamelink/pkg/cache"
"gamelink/pkg/testutil"
)

func migrateAuditModels(t *testing.T, db *gorm.DB) {
t.Helper()
err := db.AutoMigrate(&model.User{}, &model.RoleModel{}, &model.Permission{},
&model.RolePermission{}, &model.UserRole{}, &model.PermissionAuditLog{})
if err != nil {
t.Fatalf("migrate audit models: %v", err)
}
}

func setupAuditTestEnv(t *testing.T) (*gorm.DB, *audit.Service, *roleservice.RoleService,
*permissionservice.PermissionService, *model.User, *model.User) {
t.Helper()
gin.SetMode(gin.TestMode)
db := testutil.NewMemoryDB(t)
migrateAuditModels(t, db)

userRepo := user.NewUserRepository(db)
roleRepo := adminrepo.NewRoleRepository(db)
permRepo := permission.NewPermissionRepository(db)
auditRepo := permissionauditlog.NewRepository(db)
memCache := cache.NewMemory()

permSvc := permissionservice.NewPermissionService(permRepo, memCache)
roleSvc := roleservice.NewRoleService(roleRepo, memCache)
auditSvc := audit.NewService(auditRepo, audit.Config{
BufferSize: 100, BatchSize: 5, FlushInterval: 100 * time.Millisecond,
})
auditSvc.Start()

superUser := &model.User{
Name: "Super Admin", Email: "super@example.com", Phone: "18800000001",
PasswordHash: "x", Role: model.RoleAdmin, Status: model.UserStatusActive,
}
_ = userRepo.Create(ctx(), superUser)

normalUser := &model.User{
Name: "Normal User", Email: "normal@example.com", Phone: "18800000002",
PasswordHash: "x", Role: model.RoleAdmin, Status: model.UserStatusActive,
}
_ = userRepo.Create(ctx(), normalUser)

superRole := &model.RoleModel{Slug: string(model.RoleSlugSuperAdmin), Name: "Super Admin", IsSystem: true}
_ = roleRepo.Create(ctx(), superRole)
_ = roleRepo.AssignToUser(ctx(), superUser.ID, []uint64{superRole.ID})

return db, auditSvc, roleSvc, permSvc, superUser, normalUser
}

func waitForAuditFlush()                                                { time.Sleep(200 * time.Millisecond) }
func ptrAuditTargetType(t model.AuditTargetType) *model.AuditTargetType { return &t }
func ptrAuditAction(a model.AuditAction) *model.AuditAction             { return &a }
func auditTestSetUserID(userID uint64) gin.HandlerFunc {
return func(c *gin.Context) { c.Set("user_id", userID); c.Set("userID", userID); c.Next() }
}
func formatUint64(id uint64) string { return fmt.Sprintf("%d", id) }

// Task 24.1: Permission Assignment Audit Log Tests
// Requirements: 6.1

func TestPermissionAssignmentAuditLog(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

roleRepo := adminrepo.NewRoleRepository(db)
permRepo := permission.NewPermissionRepository(db)

testRole := &model.RoleModel{Slug: "audit-test-role", Name: "Audit Test Role", IsSystem: false}
require.NoError(t, roleRepo.Create(ctx(), testRole))

perm1 := &model.Permission{
Method: model.HTTPMethodGET, Path: "/api/v1/admin/audit-test-1",
Code: "audit.test.read1", Group: "audit", Description: "Audit test permission 1",
}
require.NoError(t, permRepo.Create(ctx(), perm1))

beforeData := map[string]any{"permissionIds": []uint64{}}
afterData := map[string]any{"permissionIds": []uint64{perm1.ID}}

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssignPermission, testRole.ID, testRole.Name,
beforeData, afterData, "192.168.1.1", "TestAgent/1.0", "req-perm-001")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, TargetType: ptrAuditTargetType(model.AuditTargetTypeRole), TargetID: &testRole.ID,
})
require.NoError(t, err)
require.GreaterOrEqual(t, len(result.Logs), 1)

log := result.Logs[0]
assert.Equal(t, superUser.ID, log.OperatorID)
assert.Equal(t, superUser.Name, log.OperatorName)
assert.Equal(t, model.AuditTargetTypeRole, log.TargetType)
assert.Equal(t, testRole.ID, log.TargetID)
assert.Equal(t, testRole.Name, log.TargetName)
assert.Equal(t, model.AuditActionAssignPermission, log.Action)
assert.Equal(t, "192.168.1.1", log.IPAddress)
assert.Equal(t, "TestAgent/1.0", log.UserAgent)
assert.Equal(t, "req-perm-001", log.RequestID)
assert.NotEmpty(t, log.BeforeData)
assert.NotEmpty(t, log.AfterData)
assert.False(t, log.CreatedAt.IsZero())
}

func TestPermissionAssignmentAuditLogBeforeAfterData(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssignPermission, 1, "TestRole",
map[string]any{"permissionIds": []uint64{1}},
map[string]any{"permissionIds": []uint64{1, 2, 3}},
"192.168.1.1", "TestAgent/1.0", "req-before-after")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{Page: 1, PageSize: 10, RequestID: "req-before-after"})
require.NoError(t, err)
require.Len(t, result.Logs, 1)

log := result.Logs[0]
var beforeData map[string]any
err = json.Unmarshal([]byte(log.BeforeData), &beforeData)
require.NoError(t, err)
assert.Contains(t, beforeData, "permissionIds")

var afterData map[string]any
err = json.Unmarshal([]byte(log.AfterData), &afterData)
require.NoError(t, err)
assert.Contains(t, afterData, "permissionIds")
}

// Task 24.2: User Role Assignment Audit Log Tests
// Requirements: 6.2

func TestUserRoleAssignmentAuditLog(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

roleRepo := adminrepo.NewRoleRepository(db)
role1 := &model.RoleModel{Slug: "role-assign-test-1", Name: "Role Assign Test 1", IsSystem: false}
require.NoError(t, roleRepo.Create(ctx(), role1))

auditSvc.LogUserRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssign, normalUser.ID, normalUser.Name,
map[string]any{"roles": []string{}},
map[string]any{"roles": []string{role1.Slug}},
"192.168.1.10", "TestAgent/1.0", "req-role-assign")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, TargetType: ptrAuditTargetType(model.AuditTargetTypeUser), TargetID: &normalUser.ID,
})
require.NoError(t, err)
require.GreaterOrEqual(t, len(result.Logs), 1)

log := result.Logs[0]
assert.Equal(t, superUser.ID, log.OperatorID)
assert.Equal(t, superUser.Name, log.OperatorName)
assert.Equal(t, model.AuditTargetTypeUser, log.TargetType)
assert.Equal(t, normalUser.ID, log.TargetID)
assert.Equal(t, normalUser.Name, log.TargetName)
assert.Equal(t, model.AuditActionAssign, log.Action)
assert.NotEmpty(t, log.BeforeData)
assert.NotEmpty(t, log.AfterData)
}

func TestUserRoleChangeAuditLogBeforeAfterData(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogUserRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssign, normalUser.ID, normalUser.Name,
map[string]any{"roles": []string{"viewer"}},
map[string]any{"roles": []string{"viewer", "editor"}},
"192.168.1.11", "TestAgent/1.0", "req-role-before-after")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{Page: 1, PageSize: 10, RequestID: "req-role-before-after"})
require.NoError(t, err)
require.Len(t, result.Logs, 1)

log := result.Logs[0]
var beforeData map[string]any
err = json.Unmarshal([]byte(log.BeforeData), &beforeData)
require.NoError(t, err)
assert.Contains(t, beforeData, "roles")

var afterData map[string]any
err = json.Unmarshal([]byte(log.AfterData), &afterData)
require.NoError(t, err)
assert.Contains(t, afterData, "roles")
}

func TestUserRoleAssignmentAuditLogMultipleRoles(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

roleRepo := adminrepo.NewRoleRepository(db)
role1 := &model.RoleModel{Slug: "multi-role-1", Name: "Multi Role 1", IsSystem: false}
role2 := &model.RoleModel{Slug: "multi-role-2", Name: "Multi Role 2", IsSystem: false}
require.NoError(t, roleRepo.Create(ctx(), role1))
require.NoError(t, roleRepo.Create(ctx(), role2))

auditSvc.LogUserRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssign, normalUser.ID, normalUser.Name,
map[string]any{"roles": []string{}},
map[string]any{"roles": []string{role1.Slug, role2.Slug}},
"192.168.1.12", "TestAgent/1.0", "req-multi-role-assign")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, RequestID: "req-multi-role-assign",
})
require.NoError(t, err)
require.Len(t, result.Logs, 1)

log := result.Logs[0]
var afterData map[string]any
err = json.Unmarshal([]byte(log.AfterData), &afterData)
require.NoError(t, err)

roles, ok := afterData["roles"].([]any)
require.True(t, ok)
assert.Len(t, roles, 2)
}

func TestUserRoleRemovalAuditLog(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogUserRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssign, normalUser.ID, normalUser.Name,
map[string]any{"roles": []string{"admin", "editor"}},
map[string]any{"roles": []string{"admin"}},
"192.168.1.13", "TestAgent/1.0", "req-role-removal")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, RequestID: "req-role-removal",
})
require.NoError(t, err)
require.Len(t, result.Logs, 1)

log := result.Logs[0]

var beforeData map[string]any
err = json.Unmarshal([]byte(log.BeforeData), &beforeData)
require.NoError(t, err)
beforeRoles, ok := beforeData["roles"].([]any)
require.True(t, ok)
assert.Len(t, beforeRoles, 2)

var afterData map[string]any
err = json.Unmarshal([]byte(log.AfterData), &afterData)
require.NoError(t, err)
afterRoles, ok := afterData["roles"].([]any)
require.True(t, ok)
assert.Len(t, afterRoles, 1)
}

// Task 24.3: Audit Log Query Filter Tests
// Requirements: 6.3

func TestAuditLogQueryByTimeRange(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "TimeRangeTestRole",
nil, map[string]any{"name": "TimeRangeTestRole"},
"192.168.1.20", "TestAgent/1.0", "req-time-range")

waitForAuditFlush()

now := time.Now()
dateFrom := now.Add(-1 * time.Hour)
dateTo := now.Add(1 * time.Hour)

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, DateFrom: &dateFrom, DateTo: &dateTo,
})
require.NoError(t, err)
require.NotEmpty(t, result.Logs)

for _, log := range result.Logs {
assert.True(t, log.CreatedAt.After(dateFrom) || log.CreatedAt.Equal(dateFrom))
assert.True(t, log.CreatedAt.Before(dateTo) || log.CreatedAt.Equal(dateTo))
}
}

func TestAuditLogQueryByActionType(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "CreateRole",
nil, map[string]any{"name": "CreateRole"},
"192.168.1.21", "TestAgent/1.0", "req-action-create")

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionUpdate, 2, "UpdateRole",
map[string]any{"name": "OldName"},
map[string]any{"name": "NewName"},
"192.168.1.22", "TestAgent/1.0", "req-action-update")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, Action: ptrAuditAction(model.AuditActionCreate),
})
require.NoError(t, err)
for _, log := range result.Logs {
assert.Equal(t, model.AuditActionCreate, log.Action)
}
}

func TestAuditLogQueryByOperator(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "SuperRole",
nil, map[string]any{"name": "SuperRole"},
"192.168.1.30", "TestAgent/1.0", "req-op-super")

auditSvc.LogRoleChange(ctx(), normalUser.ID, normalUser.Name,
model.AuditActionCreate, 2, "NormalRole",
nil, map[string]any{"name": "NormalRole"},
"192.168.1.31", "TestAgent/1.0", "req-op-normal")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, OperatorID: &superUser.ID,
})
require.NoError(t, err)
for _, log := range result.Logs {
assert.Equal(t, superUser.ID, log.OperatorID)
}
}

func TestAuditLogQueryByTargetType(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "RoleTarget",
nil, map[string]any{"name": "RoleTarget"},
"192.168.1.40", "TestAgent/1.0", "req-target-role")

auditSvc.LogUserRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionAssign, normalUser.ID, normalUser.Name,
nil, map[string]any{"roles": []string{"admin"}},
"192.168.1.41", "TestAgent/1.0", "req-target-user")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, TargetType: ptrAuditTargetType(model.AuditTargetTypeUser),
})
require.NoError(t, err)
for _, log := range result.Logs {
assert.Equal(t, model.AuditTargetTypeUser, log.TargetType)
}
}

func TestAuditLogQueryCombinedFilters(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "CombinedRole1",
nil, map[string]any{"name": "CombinedRole1"},
"192.168.1.50", "TestAgent/1.0", "req-combined-1")

auditSvc.LogRoleChange(ctx(), normalUser.ID, normalUser.Name,
model.AuditActionCreate, 2, "CombinedRole2",
nil, map[string]any{"name": "CombinedRole2"},
"192.168.1.51", "TestAgent/1.0", "req-combined-2")

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10, OperatorID: &superUser.ID, Action: ptrAuditAction(model.AuditActionCreate),
})
require.NoError(t, err)
for _, log := range result.Logs {
assert.Equal(t, superUser.ID, log.OperatorID)
assert.Equal(t, model.AuditActionCreate, log.Action)
}
}

func TestAuditLogQueryPagination(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

for i := 0; i < 5; i++ {
auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, uint64(i+1), fmt.Sprintf("PaginationRole%d", i),
nil, map[string]any{"index": i},
"192.168.1.60", "TestAgent/1.0", fmt.Sprintf("req-page-%d", i))
}

waitForAuditFlush()

result1, err := auditSvc.Query(ctx(), audit.QueryOptions{Page: 1, PageSize: 2})
require.NoError(t, err)
assert.LessOrEqual(t, len(result1.Logs), 2)

if result1.Total > 2 {
result2, err := auditSvc.Query(ctx(), audit.QueryOptions{Page: 2, PageSize: 2})
require.NoError(t, err)
if len(result1.Logs) > 0 && len(result2.Logs) > 0 {
assert.NotEqual(t, result1.Logs[0].ID, result2.Logs[0].ID)
}
}
}

// Task 24.4: Audit Log Export Tests
// Requirements: 6.5

func TestAuditLogExportCSV(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "ExportRole1",
nil, map[string]any{"name": "ExportRole1"},
"192.168.1.80", "ExportAgent/1.0", "req-export-1")

waitForAuditFlush()

csvData, err := auditSvc.ExportCSV(ctx(), audit.ExportOptions{MaxRecords: 100})
require.NoError(t, err)
require.NotEmpty(t, csvData)

csvContent := strings.TrimPrefix(string(csvData), "\xef\xbb\xbf")
reader := csv.NewReader(strings.NewReader(csvContent))
records, err := reader.ReadAll()
require.NoError(t, err)
require.GreaterOrEqual(t, len(records), 2)
}

func TestAuditLogExportCSVContainsRequiredColumns(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "ColumnTestRole",
nil, map[string]any{"name": "ColumnTestRole"},
"192.168.1.82", "ColumnAgent/1.0", "req-column-test")

waitForAuditFlush()

csvData, err := auditSvc.ExportCSV(ctx(), audit.ExportOptions{MaxRecords: 100})
require.NoError(t, err)

csvContent := strings.TrimPrefix(string(csvData), "\xef\xbb\xbf")
reader := csv.NewReader(strings.NewReader(csvContent))
header, err := reader.Read()
require.NoError(t, err)

requiredColumns := []string{
"ID", "时间", "操作者ID", "操作者名称", "目标类型", "目标ID",
"目标名称", "操作类型", "操作前数据", "操作后数据", "IP地址", "用户代理", "请求ID",
}
for _, col := range requiredColumns {
assert.Contains(t, header, col, "Missing required column: %s", col)
}
}

func TestAuditLogExportCSVWithFilters(t *testing.T) {
db, auditSvc, _, _, superUser, normalUser := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, 1, "FilterExportRole1",
nil, map[string]any{"name": "FilterExportRole1"},
"192.168.1.83", "FilterAgent/1.0", "req-filter-export-1")

auditSvc.LogUserRoleChange(ctx(), normalUser.ID, normalUser.Name,
model.AuditActionAssign, 2, "FilterExportUser",
nil, map[string]any{"roles": []string{"admin"}},
"192.168.1.84", "FilterAgent/1.0", "req-filter-export-2")

waitForAuditFlush()

csvData, err := auditSvc.ExportCSV(ctx(), audit.ExportOptions{
Action: ptrAuditAction(model.AuditActionCreate), MaxRecords: 100,
})
require.NoError(t, err)

csvContent := strings.TrimPrefix(string(csvData), "\xef\xbb\xbf")
reader := csv.NewReader(strings.NewReader(csvContent))
records, err := reader.ReadAll()
require.NoError(t, err)

for i := 1; i < len(records); i++ {
assert.Equal(t, string(model.AuditActionCreate), records[i][7])
}
}

func TestAuditLogExportCSVMaxRecords(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

for i := 0; i < 5; i++ {
auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, uint64(i+1), fmt.Sprintf("MaxRecordRole%d", i),
nil, map[string]any{"index": i},
"192.168.1.86", "MaxRecordAgent/1.0", fmt.Sprintf("req-max-%d", i))
}

waitForAuditFlush()

csvData, err := auditSvc.ExportCSV(ctx(), audit.ExportOptions{MaxRecords: 2})
require.NoError(t, err)

csvContent := strings.TrimPrefix(string(csvData), "\xef\xbb\xbf")
reader := csv.NewReader(strings.NewReader(csvContent))
records, err := reader.ReadAll()
require.NoError(t, err)
assert.LessOrEqual(t, len(records), 3)
}

func TestAuditLogExportFilename(t *testing.T) {
filename := audit.GenerateExportFilename("audit_log")
assert.True(t, strings.HasPrefix(filename, "audit_log_"))
assert.True(t, strings.HasSuffix(filename, ".csv"))
}

// Additional Integration Tests

func TestAuditLogAsyncWriting(t *testing.T) {
db, auditSvc, _, _, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

start := time.Now()

for i := 0; i < 10; i++ {
auditSvc.LogRoleChange(ctx(), superUser.ID, superUser.Name,
model.AuditActionCreate, uint64(i+1), fmt.Sprintf("AsyncRole%d", i),
nil, map[string]any{"index": i},
"192.168.1.90", "AsyncAgent/1.0", fmt.Sprintf("req-async-%d", i))
}

elapsed := time.Since(start)
assert.Less(t, elapsed, 100*time.Millisecond)

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{Page: 1, PageSize: 20})
require.NoError(t, err)
assert.GreaterOrEqual(t, len(result.Logs), 1)
}

func TestAuditLogServiceStats(t *testing.T) {
db, auditSvc, _, _, _, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

stats := auditSvc.GetStats()
assert.True(t, stats.Running)
assert.Greater(t, stats.BufferSize, 0)
assert.GreaterOrEqual(t, stats.ProcessedCount, int64(0))
}

func TestAuditLogWithRoleHandler(t *testing.T) {
db, auditSvc, roleSvc, permSvc, superUser, _ := setupAuditTestEnv(t)
defer testutil.CleanDB(t, db)
defer auditSvc.Stop()

roleRepo := adminrepo.NewRoleRepository(db)
permRepo := permission.NewPermissionRepository(db)

testRole := &model.RoleModel{Slug: "handler-audit-role", Name: "Handler Audit Role", IsSystem: false}
require.NoError(t, roleRepo.Create(ctx(), testRole))

perm1 := &model.Permission{
Method: model.HTTPMethodGET, Path: "/api/v1/admin/handler-audit-1",
Code: "handler.audit.read1", Group: "handler", Description: "Handler audit permission 1",
}
require.NoError(t, permRepo.Create(ctx(), perm1))

roleHandler := adminhandler.NewRoleHandlerWithAudit(roleSvc, auditSvc)
pm := middleware.NewPermissionMiddleware(nil, permSvc, roleSvc)

router := gin.New()
router.PUT("/roles/:id/permissions/batch",
auditTestSetUserID(superUser.ID),
pm.RequirePermission(model.HTTPMethodPUT, "/api/v1/admin/roles/:id/permissions/batch"),
roleHandler.AssignPermissions,
)

reqBody := map[string]any{"permissionIds": []uint64{perm1.ID}}
resp := doJSON(router, http.MethodPut, "/roles/"+formatUint64(testRole.ID)+"/permissions/batch", reqBody, "")
assert.Equal(t, http.StatusOK, resp.Code)

waitForAuditFlush()

result, err := auditSvc.Query(ctx(), audit.QueryOptions{
Page: 1, PageSize: 10,
TargetType: ptrAuditTargetType(model.AuditTargetTypeRole),
TargetID:   &testRole.ID,
Action:     ptrAuditAction(model.AuditActionAssignPermission),
})
require.NoError(t, err)
require.GreaterOrEqual(t, len(result.Logs), 1)
}
