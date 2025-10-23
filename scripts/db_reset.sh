#!/bin/bash
# ============================================================================
# FITS Backend - Database Reset Script
# ============================================================================
# Truncates all tables for a fresh start without dropping the database
# Preserves database structure and migrations
# Safe to run multiple times
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
RESET='\033[0m'

# Database configuration (from config.toml)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-fits_db}"
DB_USER="${DB_USER:-fits_user}"
DB_PASSWORD="${DB_PASSWORD:-fits_password}"

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${CYAN}FITS Backend - Database Reset${RESET}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${YELLOW}⚠  This will DELETE ALL DATA in the database!${RESET}"
echo ""

# Truncate all tables in the correct order (respecting foreign keys)
echo -e "${CYAN}→${RESET} Resetting database tables..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'SQL'
-- Disable triggers temporarily for faster truncation
SET session_replication_role = replica;

-- Truncate all tables (CASCADE handles foreign keys)
TRUNCATE TABLE signatures CASCADE;
TRUNCATE TABLE reports CASCADE;
TRUNCATE TABLE teacher_keys CASCADE;
TRUNCATE TABLE refresh_tokens CASCADE;
TRUNCATE TABLE invitations CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE students CASCADE;
TRUNCATE TABLE teachers CASCADE;

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Verify tables are empty
DO $$
DECLARE
    table_record RECORD;
    row_count INTEGER;
BEGIN
    FOR table_record IN
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public'
        AND table_name NOT IN ('schema_migrations')
        ORDER BY table_name
    LOOP
        EXECUTE 'SELECT COUNT(*) FROM ' || table_record.table_name INTO row_count;
        RAISE NOTICE '  ✓ % (% rows)', table_record.table_name, row_count;
    END LOOP;
END $$;
SQL

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓${RESET} Database reset successful!"
    echo ""
    echo -e "${CYAN}Next steps:${RESET}"
    echo -e "  ${YELLOW}make run${RESET}       - Start server (migrations run automatically)"
    echo -e "  ${YELLOW}make db-seed${RESET}   - Populate with test data"
    echo ""
else
    echo ""
    echo -e "${RED}✗${RESET} Database reset failed"
    echo ""
    echo "Possible causes:"
    echo "  1. PostgreSQL is not running → make db-up"
    echo "  2. Wrong credentials in config.toml"
    echo "  3. Database does not exist → make run (auto-creates)"
    echo ""
    exit 1
fi
