#!/bin/bash
# GameLink Database Backup Script
# Author: Database-Architect
# Version: 1.0.0

set -e

ENVIRONMENT=${1:-dev}
BACKUP_DIR="./backups/${ENVIRONMENT}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/gamelink_backup_${TIMESTAMP}.sql.gz"
LOG_FILE="${BACKUP_DIR}/backup_${TIMESTAMP}.log"

mkdir -p ${BACKUP_DIR}

echo "=== GameLink Database Backup ==="
echo "Environment: ${ENVIRONMENT}"
echo "Backup file: ${BACKUP_FILE}"
echo ""

# Check if Docker Compose is running
if ! docker-compose ps postgres 2>&1 | grep -q "Up"; then
    echo "ERROR: PostgreSQL container is not running"
    exit 1
fi

# Start timer
START_TIME=$(date +%s)

# Create backup
docker-compose exec -T postgres pg_dump \
    -U gamelink \
    -d gamelink \
    --verbose \
    --no-owner \
    --no-acl \
    --format=plain \
    --encoding=UTF8 \
    2>&1 | gzip > ${BACKUP_FILE}

BACKUP_EXIT_CODE=${PIPESTATUS[0]}

# End timer
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

if [ ${BACKUP_EXIT_CODE} -eq 0 ]; then
    BACKUP_SIZE=$(ls -lh ${BACKUP_FILE} | awk '{print $5}')
    echo "✅ Backup completed successfully!"
    echo "Backup file: ${BACKUP_FILE}"
    echo "Backup size: ${BACKUP_SIZE}"
    echo "Duration: ${DURATION}s"
else
    echo "❌ Backup failed!"
    exit 1
fi
