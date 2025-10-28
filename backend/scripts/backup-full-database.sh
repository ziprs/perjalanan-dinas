#!/bin/bash

# ============================================================================
# Full Database Backup Script - Before Cleanup
# ============================================================================

# Configuration
DB_NAME="perjalanan_dinas"
DB_USER="postgres"
BACKUP_DIR="/var/backups/perjalanan-dinas"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/perjalanan_dinas_before_cleanup_${TIMESTAMP}.sql"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "============================================================================"
echo "Database Backup Script - Perjalanan Dinas SPD"
echo "============================================================================"
echo ""

# Create backup directory if not exists
echo -e "${YELLOW}[1/5]${NC} Creating backup directory..."
mkdir -p "$BACKUP_DIR"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Backup directory ready: $BACKUP_DIR"
else
    echo -e "${RED}✗${NC} Failed to create backup directory"
    exit 1
fi

echo ""

# Check if database exists
echo -e "${YELLOW}[2/5]${NC} Checking database connection..."
psql -U "$DB_USER" -d "$DB_NAME" -c "SELECT version();" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Database connection successful"
else
    echo -e "${RED}✗${NC} Cannot connect to database"
    exit 1
fi

echo ""

# Show database statistics before backup
echo -e "${YELLOW}[3/5]${NC} Database statistics before backup:"
psql -U "$DB_USER" -d "$DB_NAME" -c "
SELECT
    'travel_requests' AS table_name,
    COUNT(*) AS total_records
FROM travel_requests
UNION ALL
SELECT 'travel_reports', COUNT(*) FROM travel_reports
UNION ALL
SELECT 'employees', COUNT(*) FROM employees
UNION ALL
SELECT 'positions', COUNT(*) FROM positions
UNION ALL
SELECT 'admins', COUNT(*) FROM admins;
"

echo ""

# Create backup
echo -e "${YELLOW}[4/5]${NC} Creating full database backup..."
pg_dump -U "$DB_USER" -d "$DB_NAME" -F p -f "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Backup created successfully"
else
    echo -e "${RED}✗${NC} Backup failed"
    exit 1
fi

echo ""

# Verify backup file
echo -e "${YELLOW}[5/5]${NC} Verifying backup file..."
if [ -f "$BACKUP_FILE" ]; then
    FILE_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}✓${NC} Backup file exists"
    echo "   Location: $BACKUP_FILE"
    echo "   Size: $FILE_SIZE"
else
    echo -e "${RED}✗${NC} Backup file not found"
    exit 1
fi

echo ""
echo "============================================================================"
echo -e "${GREEN}BACKUP COMPLETED SUCCESSFULLY${NC}"
echo "============================================================================"
echo ""
echo "Backup Details:"
echo "  - Database: $DB_NAME"
echo "  - File: $BACKUP_FILE"
echo "  - Size: $FILE_SIZE"
echo "  - Timestamp: $TIMESTAMP"
echo ""
echo "To restore this backup, run:"
echo "  psql -U $DB_USER -d $DB_NAME -f $BACKUP_FILE"
echo ""
echo "============================================================================"

# Create a restore script for convenience
RESTORE_SCRIPT="${BACKUP_DIR}/restore_${TIMESTAMP}.sh"
cat > "$RESTORE_SCRIPT" << EOF
#!/bin/bash
# Restore script for backup created on $TIMESTAMP

echo "============================================"
echo "Database Restore Script"
echo "============================================"
echo ""
echo "WARNING: This will OVERWRITE the current database!"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirm

if [ "\$confirm" != "yes" ]; then
    echo "Restore cancelled."
    exit 0
fi

echo ""
echo "Restoring database from backup..."
psql -U $DB_USER -d $DB_NAME -f $BACKUP_FILE

if [ \$? -eq 0 ]; then
    echo ""
    echo "✓ Database restored successfully!"
else
    echo ""
    echo "✗ Restore failed!"
    exit 1
fi
EOF

chmod +x "$RESTORE_SCRIPT"
echo "Restore script created: $RESTORE_SCRIPT"
echo ""

exit 0
