# Database Migration Guide

**Version:** 1.0.0
**Author:** Database-Architect
**Last Updated:** 2026-02-09

---

## Table of Contents

1. [Migration Architecture](#migration-architecture)
2. [Current Migration Status](#current-migration-status)
3. [Migration Execution Guide](#migration-execution-guide)
4. [Backup & Recovery](#backup--recovery)
5. [Rollback Procedures](#rollback-procedures)
6. [Validation & Testing](#validation--testing)
7. [Troubleshooting](#troubleshooting)

---

## Migration Architecture

### Migration System Overview

GameLink uses **GORM AutoMigrate** with version tracking:

```
┌─────────────────────────────────────────────────────────────┐
│                     Application Startup                      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 1: Check Migration Version (isMigrateUpToDate)        │
│  - Reads 'migrate_version' from seed_metadata table          │
│  - If version matches → SKIP migration                       │
│  - If version differs → PROCEED to migration                 │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 2: Special Orders Migration (prepareOrdersMigration) │
│  - Checks and adds item_id, order_no, unit_price_cents       │
│  - PostgreSQL-specific schema evolution                      │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 3: AutoMigrate - Base Tables                         │
│  - Game, User, Player, Order, Payment                        │
│  - Tables without foreign key dependencies                   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 4: AutoMigrate - Dependent Tables                    │
│  - All 80+ tables with foreign keys                          │
│  - Order items, wallets, reviews, chat, etc.                 │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 5: Data Fixups (runDataFixups)                        │
│  - Normalize order status spelling                           │
│  - Generate order numbers                                    │
│  - Ensure default roles and commission rules                 │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 6: Index Creation (ensureIndexes)                     │
│  - Composite indexes for performance                         │
│  - Covering indexes for common queries                       │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Phase 7: Seed Data (applySeeds)                             │
│  - Initial games, categories, demo data                      │
│  - Only if SEED_ENABLED=true                                 │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Mark Migration Version (markMigrateVersion)                 │
│  - Writes current version to seed_metadata                   │
│  - Prevents re-running on next startup                       │
└─────────────────────────────────────────────────────────────┘
```

### Current Version

**Migration Version:** `2026-02-07-v1`
**Location:** `api/pkg/db/migrate.go:115`

### Version Tracking

```sql
-- Version is stored in seed_metadata table
CREATE TABLE IF NOT EXISTS seed_metadata (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- Current version record
INSERT INTO seed_metadata (key, value) VALUES ('migrate_version', '2026-02-07-v1');
```

---

## Current Migration Status

### ✅ Implemented Migrations

| Phase | Status | Tables | Notes |
|-------|--------|--------|-------|
| **Base Tables** | ✅ Complete | 5 | Game, User, Player, Order, Payment |
| **Dependent Tables** | ✅ Complete | 75+ | All related tables |
| **Data Fixups** | ✅ Complete | - | Order status, order numbers |
| **Indexes** | ✅ Complete | 40+ | Performance indexes |
| **Seed Data** | ✅ Complete | - | v4 version (2026-02-07) |

### 📊 Database Statistics

- **Total Tables:** 80+
- **Total Indexes:** 40+ composite indexes
- **Foreign Keys:** 50+ relationships
- **Database Size (Development):** ~50MB (with seed data)

### ⚠️ Known Limitations

1. **No SQL Migration Files**
   - Current system uses GORM AutoMigrate only
   - No manual SQL migration files
   - All schema changes through Go structs

2. **No Automatic Rollback**
   - GORM doesn't support down migrations
   - Rollback requires manual SQL scripts
   - Backups are critical before migration

3. **Idempotent Safe**
   - `isMigrateUpToDate()` prevents re-running
   - Safe to run multiple times
   - No duplicate schema changes

---

## Migration Execution Guide

### Pre-Migration Checklist

#### ✅ Environment Verification

```bash
# 1. Check current migration version
docker-compose exec postgres psql -U gamelink -d gamelink -c \
  "SELECT value FROM seed_metadata WHERE key = 'migrate_version';"

# 2. Verify database connection
docker-compose exec postgres pg_isready -U gamelink

# 3. Check table count
docker-compose exec postgres psql -U gamelink -d gamellink -c \
  "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"

# 4. Verify indexes
docker-compose exec postgres psql -U gamelink -d gamelink -c \
  "SELECT COUNT(*) FROM pg_indexes WHERE schemaname = 'public';"
```

#### ✅ Backup Verification

```bash
# Ensure backup directory exists
ls -lh ./backups/

# Check last backup timestamp
ls -lt ./backups/ | head -5
```

#### ✅ Configuration Check

```bash
# Verify environment variables
docker-compose config | grep -E "(POSTGRES_|DB_|SEED_)"

# Check database settings
docker-compose exec postgres psql -U gamelink -d gamelink -c "SHOW max_connections;"
docker-compose exec postgres psql -U gamelink -d gamelink -c "SHOW shared_buffers;"
```

### Migration Execution Methods

#### Method 1: Application Startup (Automatic)

**Best for:** Development, Staging

```bash
# Migration runs automatically on startup
docker-compose up -d backend

# Watch migration logs
docker-compose logs -f backend | grep -E "(auto-migrate|Phase)"

# Expected output:
# [startup] auto-migrate outdated, running migration (version=2026-02-07-v1)...
# [startup] Phase 1: Creating base tables...
# [startup] Phase 1: Base tables created successfully
# [startup] Phase 2: Creating dependent tables...
# [startup] Phase 2: Dependent tables created successfully
# [startup] indexes: 234ms
# [startup] seed data: 1.2s
```

#### Method 2: Manual Migration Command

**Best for:** Production, Controlled Rollout

```bash
# Run migration explicitly
docker-compose exec backend /app/api migrate

# Or via Go command
cd api
go run cmd/main.go migrate
```

### Post-Migration Validation

```sql
-- 1. Verify migration version
SELECT value FROM seed_metadata WHERE key = 'migrate_version';
-- Expected: 2026-02-07-v1

-- 2. Check all tables exist
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;

-- 3. Verify critical indexes
SELECT indexname, tablename
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname LIKE 'idx_%'
ORDER BY tablename, indexname;

-- 4. Check foreign keys
SELECT
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
ORDER BY tc.table_name;

-- 5. Verify data integrity
SELECT COUNT(*) FROM users;           -- Should be > 0 if seeded
SELECT COUNT(*) FROM games;           -- Should be > 0 if seeded
SELECT COUNT(*) FROM players;         -- Should be > 0 if seeded
```

---

## Backup & Recovery

### Backup Strategy

#### Full Backup (Recommended before migrations)

```bash
#!/bin/bash
# scripts/backup-database.sh

BACKUP_DIR="./backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/gamelink_backup_${TIMESTAMP}.sql.gz"

mkdir -p ${BACKUP_DIR}

echo "Creating full database backup..."

docker-compose exec -T postgres pg_dump \
  -U gamelink \
  -d gamelink \
  --verbose \
  --no-owner \
  --no-acl \
  --format=plain \
  --encoding=UTF8 \
  2>&1 | gzip > ${BACKUP_FILE}

if [ $? -eq 0 ]; then
  echo "✅ Backup created: ${BACKUP_FILE}"
  ls -lh ${BACKUP_FILE}
else
  echo "❌ Backup failed!"
  exit 1
fi
```

#### Schema-Only Backup (For testing)

```bash
# Backup schema only (no data)
docker-compose exec postgres pg_dump \
  -U gamelink \
  -d gamelink \
  --schema-only \
  --no-owner \
  --no-acl \
  > ./backups/schema_$(date +%Y%m%d).sql
```

#### Data-Only Backup (For migration testing)

```bash
# Backup data only (no schema)
docker-compose exec postgres pg_dump \
  -U gamelink \
  -d gamelink \
  --data-only \
  --no-owner \
  --no-acl \
  > ./backups/data_$(date +%Y%m%d).sql
```

### Backup Retention Policy

| Environment | Frequency | Retention | Storage |
|-------------|-----------|-----------|---------|
| **Development** | Daily (before migrations) | 7 days | Local |
| **Staging** | Daily | 30 days | Local |
| **Production** | Hourly + Daily | 90 days | Remote (S3/Glacier) |

### Recovery Procedures

#### Full Database Restore

```bash
#!/bin/bash
# scripts/restore-database.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: $0 <backup_file.sql.gz>"
  exit 1
fi

echo "⚠️  WARNING: This will replace the entire database!"
echo "Backup file: ${BACKUP_FILE}"
read -p "Continue? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
  echo "Restore cancelled."
  exit 0
fi

echo "Stopping backend service..."
docker-compose stop backend

echo "Restoring database..."
gunzip -c ${BACKUP_FILE} | docker-compose exec -T postgres psql \
  -U gamelink \
  -d gamelink \
  --quiet \
  --set ON_ERROR_STOP=on

if [ $? -eq 0 ]; then
  echo "✅ Database restored successfully"
  echo "Starting backend service..."
  docker-compose start backend
else
  echo "❌ Restore failed!"
  exit 1
fi
```

#### Point-in-Time Recovery (PITR)

For production environments with WAL archiving enabled:

```bash
# 1. Identify the target time
TARGET_TIME="2026-02-09 10:00:00"

# 2. Restore from base backup
# 3. Replay WAL logs up to target time
docker-compose exec postgres psql \
  -U gamelink \
  -d gamelink \
  -c "SELECT pg_wal_replay_resume();"
```

---

## Rollback Procedures

### ⚠️ Important: No Automatic Rollback

GORM AutoMigrate does **NOT** support automatic rollbacks. Rollback requires:

1. Restore from backup (recommended)
2. Manual SQL scripts (risky, not recommended)

### Rollback Strategy

#### Option 1: Backup Restore (Recommended)

```bash
# 1. Stop application
docker-compose stop backend

# 2. Restore from pre-migration backup
./scripts/restore-database.sh ./backups/gamelink_backup_PRE_MIGRATION.sql.gz

# 3. Verify rollback
docker-compose exec postgres psql -U gamelink -d gamelink \
  -c "SELECT value FROM seed_metadata WHERE key = 'migrate_version';"

# 4. Start application
docker-compose start backend
```

#### Option 2: Manual Schema Reversion (Advanced)

**Only use if you understand the schema changes!**

```sql
-- Example: Revert adding a column
ALTER TABLE payments DROP COLUMN IF EXISTS prepay_id;

-- Example: Revert table creation
DROP TABLE IF EXISTS payment_callback_logs;

-- Example: Revert index changes
DROP INDEX IF EXISTS idx_payments_prepay_id;
```

### Rollback Testing

```bash
# Test rollback in Staging before production
#!/bin/bash
# scripts/test-rollback.sh

echo "=== Migration Rollback Test ==="

# 1. Create test backup
echo "1. Creating pre-test backup..."
./scripts/backup-database.sh
PRE_BACKUP="./backups/gamelink_backup_$(date +%Y%m%d_%H%M%S).sql.gz"

# 2. Run migration
echo "2. Running migration..."
docker-compose exec backend /app/api migrate

# 3. Verify migration
echo "3. Verifying migration..."
VERSION=$(docker-compose exec postgres psql -U gamelink -d gamelink \
  -t -c "SELECT value FROM seed_metadata WHERE key = 'migrate_version';")
echo "Migration version: ${VERSION}"

# 4. Restore backup
echo "4. Testing rollback..."
./scripts/restore-database.sh ${PRE_BACKUP}

# 5. Verify rollback
echo "5. Verifying rollback..."
ROLLBACK_VERSION=$(docker-compose exec postgres psql -U gamelink -d gamelink \
  -t -c "SELECT value FROM seed_metadata WHERE key = 'migrate_version';")

if [ "${ROLLBACK_VERSION}" != "${VERSION}" ]; then
  echo "✅ Rollback test PASSED"
else
  echo "❌ Rollback test FAILED"
  exit 1
fi
```

---

## Validation & Testing

### Pre-Migration Validation

```sql
-- scripts/validate-pre-migration.sql

-- 1. Check database version
\echo '=== Database Version ==='
SELECT version();

-- 2. Check migration version
\echo '=== Current Migration Version ==='
SELECT value FROM seed_metadata WHERE key = 'migrate_version';

-- 3. Check table count
\echo '=== Table Count ==='
SELECT COUNT(*) AS table_count
FROM information_schema.tables
WHERE table_schema = 'public';

-- 4. Check for orphaned records
\echo '=== Orphaned Records Check ==='
SELECT COUNT(*) AS orphaned_orders
FROM orders o
LEFT JOIN users u ON o.user_id = u.id
WHERE u.id IS NULL;

-- 5. Check data integrity
\echo '=== Data Integrity Check ==='
SELECT
  (SELECT COUNT(*) FROM users) AS user_count,
  (SELECT COUNT(*) FROM players) AS player_count,
  (SELECT COUNT(*) FROM orders) AS order_count,
  (SELECT COUNT(*) FROM payments) AS payment_count;

-- 6. Check index usage
\echo '=== Index Usage ==='
SELECT
  schemaname,
  tablename,
  indexname,
  idx_scan AS index_scans,
  idx_tup_read AS tuples_read,
  idx_tup_fetch AS tuples_fetched
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC
LIMIT 20;

-- 7. Check table sizes
\echo '=== Table Sizes ==='
SELECT
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 10;
```

### Post-Migration Validation

```sql
-- scripts/validate-post-migration.sql

-- 1. Verify migration version
\echo '=== Migration Version ==='
SELECT value FROM seed_metadata WHERE key = 'migrate_version';

-- 2. Verify all tables exist
\echo '=== Critical Tables Check ==='
SELECT
  table_name,
  CASE
    WHEN table_name IN (
      'users', 'players', 'games', 'orders', 'payments',
      'wallets', 'reviews', 'chat_groups', 'chat_messages'
    ) THEN '✅ Critical'
    ELSE '✓ Standard'
  END AS importance
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY
  CASE
    WHEN table_name IN ('users', 'players', 'games', 'orders', 'payments') THEN 1
    ELSE 2
  END,
  table_name;

-- 3. Verify indexes
\echo '=== Performance Indexes Check ==='
SELECT COUNT(*) AS index_count
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname LIKE 'idx_%';

-- 4. Verify foreign keys
\echo '=== Foreign Keys Check ==='
SELECT COUNT(*) AS fk_count
FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY'
  AND table_schema = 'public';

-- 5. Data integrity checks
\echo '=== Data Integrity ==='
SELECT
  'users' AS table_name,
  COUNT(*) AS record_count,
  COUNT(CASE WHEN email IS NULL AND phone IS NULL THEN 1 END) AS null_contacts,
  COUNT(CASE WHEN password_hash IS NULL THEN 1 END) AS null_passwords
FROM users
UNION ALL
SELECT
  'players',
  COUNT(*),
  COUNT(CASE WHEN user_id IS NULL THEN 1 END),
  NULL
FROM players
UNION ALL
SELECT
  'orders',
  COUNT(*),
  COUNT(CASE WHEN user_id IS NULL THEN 1 END),
  COUNT(CASE WHEN order_no IS NULL OR order_no = '' THEN 1 END)
FROM orders;

-- 6. Performance test queries
\echo '=== Performance Test ==='
EXPLAIN ANALYZE
SELECT * FROM orders
WHERE status = 'pending'
ORDER BY created_at DESC
LIMIT 20;

EXPLAIN ANALYZE
SELECT * FROM payments
WHERE user_id = 1
ORDER BY created_at DESC
LIMIT 20;
```

### Automated Validation Script

```bash
#!/bin/bash
# scripts/validate-migration.sh

echo "=== Migration Validation ==="

# 1. Pre-migration backup
echo "1. Creating pre-migration backup..."
./scripts/backup-database.sh

# 2. Run pre-migration validation
echo "2. Running pre-migration validation..."
docker-compose exec postgres psql -U gamelink -d gamelink \
  -f /dev/stdin < scripts/validate-pre-migration.sql

# 3. Execute migration
echo "3. Executing migration..."
docker-compose exec backend /app/api migrate

# 4. Run post-migration validation
echo "4. Running post-migration validation..."
docker-compose exec postgres psql -U gamelink -d gamelink \
  -f /dev/stdin < scripts/validate-post-migration.sql

# 5. Create post-migration backup
echo "5. Creating post-migration backup..."
./scripts/backup-database.sh

echo "✅ Migration validation complete!"
```

---

## Troubleshooting

### Common Issues

#### Issue 1: Migration Stuck

**Symptoms:**
- Migration phase hanging
- No log output for > 5 minutes

**Solution:**
```bash
# 1. Check PostgreSQL connections
docker-compose exec postgres psql -U gamelink -d gamelink \
  -c "SELECT count(*) FROM pg_stat_activity;"

# 2. Check for blocking locks
docker-compose exec postgres psql -U gamelink -d gamelink \
  -c "SELECT * FROM pg_stat_activity WHERE waiting IS NOT NULL;"

# 3. Kill blocking queries (if safe)
docker-compose exec postgres psql -U gamelink -d gamelink \
  -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE waiting IS NOT NULL;"

# 4. Restart migration
docker-compose restart backend
```

#### Issue 2: Foreign Key Constraint Failure

**Symptoms:**
```
Error: insert or update on table "orders" violates foreign key constraint
```

**Solution:**
```sql
-- 1. Identify the constraint
SELECT
  tc.constraint_name,
  tc.table_name,
  kcu.column_name,
  ccu.table_name AS foreign_table_name,
  ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
  ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
  ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_name = 'orders';

-- 2. Check for orphaned records
SELECT * FROM orders
WHERE user_id NOT IN (SELECT id FROM users);

-- 3. Fix orphaned records (if safe)
DELETE FROM orders
WHERE user_id NOT IN (SELECT id FROM users);
```

#### Issue 3: Index Creation Timeout

**Symptoms:**
```
Error: creating index timeout
```

**Solution:**
```sql
-- 1. Create index CONCURRENTLY (non-blocking)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_status_created
ON orders (status, created_at DESC);

-- 2. Check index creation progress
SELECT
  relname,
  idx_scan,
  idx_tup_read,
  idx_tup_fetch
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;
```

#### Issue 4: Version Conflict

**Symptoms:**
```
Error: migration version mismatch
```

**Solution:**
```sql
-- 1. Check current version
SELECT value FROM seed_metadata WHERE key = 'migrate_version';

-- 2. Force update version (only if you're sure!)
UPDATE seed_metadata
SET value = '2026-02-07-v1'
WHERE key = 'migrate_version';

-- 3. Re-run migration
docker-compose restart backend
```

### Getting Help

If you encounter issues not covered here:

1. **Check logs:**
   ```bash
   docker-compose logs backend | grep -i error
   docker-compose logs postgres | grep -i error
   ```

2. **Verify backup:**
   ```bash
   ls -lh ./backups/
   ```

3. **Contact Database-Architect:**
   - Provide error messages
   - Share migration logs
   - Confirm backup exists

---

## Best Practices

### ✅ DO

1. **Always backup before migration**
   ```bash
   ./scripts/backup-database.sh
   ```

2. **Test in Staging first**
   - Run migration in Staging environment
   - Validate all functionality
   - Test rollback procedure

3. **Monitor migration progress**
   ```bash
   docker-compose logs -f backend | grep migrate
   ```

4. **Validate post-migration**
   ```bash
   ./scripts/validate-migration.sh
   ```

5. **Keep backups for retention period**
   - Development: 7 days
   - Staging: 30 days
   - Production: 90 days

### ❌ DON'T

1. **Never skip backup**
   - Always backup before any migration

2. **Never migrate during peak hours**
   - Schedule migrations for low-traffic periods

3. **Never ignore migration errors**
   - All errors should be investigated

4. **Never manually modify schema**
   - Use GORM AutoMigrate
   - Manual changes can break future migrations

5. **Never delete old backups immediately**
   - Keep at least 3 backup generations

---

## Appendices

### Appendix A: Migration File Locations

```
api/pkg/db/
├── migrate.go              # Main migration logic
├── postgres.go             # PostgreSQL initialization
├── seed.go                 # Seed data functions
├── seedContent.go          # Content seed data
├── seed_demo_extensions.go # Demo extensions
└── seed_flow.go            # Workflow seed data
```

### Appendix B: Environment Variables

```bash
# Database Configuration
POSTGRES_USER=gamelink
POSTGRES_PASSWORD=gamelink123
POSTGRES_DB=gamelink
POSTGRES_PORT=5432

# Connection Pool
DB_MAX_CONNS=50
DB_MAX_IDLE=25

# Seed Data
SEED_ENABLED=true          # Set to false in production
```

### Appendix C: Useful Queries

```sql
-- Check migration history
SELECT * FROM seed_metadata;

-- List all tables
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;

-- Check table sizes
SELECT
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Find slow queries
SELECT
  query,
  calls,
  total_time,
  mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-09
**Next Review:** After next migration
