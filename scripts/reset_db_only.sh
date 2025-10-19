#!/bin/bash

# FITS Backend - Database Quick Reset
# Truncates all tables for fresh bootstrap without dropping database

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}        FITS Backend - Database Reset (Quick)                   ${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-fits_db}"
DB_USER="${DB_USER:-fits_user}"
DB_PASSWORD="${DB_PASSWORD:-fits_password}"

echo -e "${YELLOW}Truncating all tables in database...${NC}"
echo ""

# SQL to truncate all data
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" << 'EOF'
-- Truncate all tables (in correct order due to foreign keys)
TRUNCATE TABLE signatures CASCADE;
TRUNCATE TABLE reports CASCADE;
TRUNCATE TABLE teacher_keys CASCADE;
TRUNCATE TABLE refresh_tokens CASCADE;
TRUNCATE TABLE invitations CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE students CASCADE;
TRUNCATE TABLE teachers CASCADE;

-- Confirmation
SELECT 'All tables truncated successfully!' as status;
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}[OK] Database reset successful!${NC}"
    echo ""
    echo "All tables are empty - you can now start fresh:"
    echo ""
    echo -e "${YELLOW}curl -X POST http://localhost:8080/api/v1/bootstrap/init${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}[ERROR] Failed to reset database${NC}"
    echo ""
    echo "Possible causes:"
    echo "  1. PostgreSQL is not running"
    echo "  2. Wrong credentials in configuration"
    echo "  3. Database does not exist"
    echo ""
    echo "Try:"
    echo "  - Stop and restart server (creates DB automatically)"
    echo "  - Check PostgreSQL: pg_isready -h localhost -p 5432"
    exit 1
fi
