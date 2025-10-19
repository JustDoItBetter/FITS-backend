#!/bin/bash

# FITS Backend - Complete Reset Script
# Deletes database and generated files for fresh initialization

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}        FITS Backend - Complete System Reset                    ${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

# Configuration from config.toml or environment
DB_NAME="${DB_NAME:-fits_db}"
DB_USER="${DB_USER:-fits_user}"
POSTGRES_USER="${POSTGRES_USER:-postgres}"

echo -e "${YELLOW}WARNING: This will delete all data!${NC}"
echo ""
echo "The following will be deleted:"
echo "  * PostgreSQL database: $DB_NAME"
echo "  * Admin RSA keys: configs/keys/"
echo "  * Uploads: uploads/"
echo "  * Coverage data: coverage.out"
echo ""
read -p "Continue? (yes/no): " -r
echo ""

if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${RED}Aborted.${NC}"
    exit 1
fi

echo -e "${BLUE}[1/5]${NC} Stopping server (if running)..."
# Try to stop server (ignore errors)
pkill -f "fits-server" 2>/dev/null || true
pkill -f "go run cmd/server/main.go" 2>/dev/null || true
echo -e "${GREEN}[OK]${NC} Server stopped"
echo ""

echo -e "${BLUE}[2/5]${NC} Deleting database..."
if command -v psql &> /dev/null; then
    # Check if database exists
    if psql -U "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
        echo "Dropping database: $DB_NAME"
        psql -U "$POSTGRES_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;" 2>/dev/null || {
            echo -e "${YELLOW}[WARN]${NC} Could not drop database (permissions?)"
            echo "Trying with sudo:"
            sudo -u postgres psql -c "DROP DATABASE IF EXISTS $DB_NAME;" || true
        }
        echo -e "${GREEN}[OK]${NC} Database deleted"
    else
        echo -e "${YELLOW}[INFO]${NC} Database does not exist (already clean)"
    fi
else
    echo -e "${YELLOW}[WARN]${NC} PostgreSQL CLI not found - delete database manually"
fi
echo ""

echo -e "${BLUE}[3/5]${NC} Deleting generated RSA keys..."
if [ -d "configs/keys" ]; then
    rm -f configs/keys/admin.key
    rm -f configs/keys/admin.pub
    rm -f configs/keys/*.key
    rm -f configs/keys/*.pub
    echo -e "${GREEN}[OK]${NC} RSA keys deleted"
else
    echo -e "${YELLOW}[INFO]${NC} No keys directory found"
fi
echo ""

echo -e "${BLUE}[4/5]${NC} Deleting uploads..."
if [ -d "uploads" ]; then
    rm -rf uploads/*
    echo -e "${GREEN}[OK]${NC} Uploads deleted"
else
    echo -e "${YELLOW}[INFO]${NC} No uploads directory"
fi
echo ""

echo -e "${BLUE}[5/5]${NC} Deleting test artifacts..."
rm -f coverage.out
rm -f coverage.html
rm -f *.test
echo -e "${GREEN}[OK]${NC} Test artifacts deleted"
echo ""

echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}              Reset completed successfully!                     ${NC}"
echo -e "${GREEN}================================================================${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo ""
echo "1. Start server:"
echo -e "   ${YELLOW}go run cmd/server/main.go${NC}"
echo ""
echo "2. Initialize admin:"
echo -e "   ${YELLOW}curl -X POST http://localhost:8080/api/v1/bootstrap/init${NC}"
echo ""
echo "3. Save and use the admin token!"
echo ""
