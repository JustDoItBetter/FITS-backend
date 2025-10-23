#!/bin/bash
# ============================================================================
# FITS Backend - Database Seed Script
# ============================================================================
# Populates database with development/test data
# Creates sample teachers, students, and invitations
# Safe to run on fresh or existing database
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BLUE='\033[0;34m'
RESET='\033[0m'

# Server configuration
SERVER_URL="${SERVER_URL:-http://localhost:8080}"

echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${CYAN}FITS Backend - Database Seed${RESET}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""

# Check if server is running
if ! curl -s "$SERVER_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}✗${RESET} Server is not running on $SERVER_URL"
    echo ""
    echo "Please start the server first:"
    echo -e "  ${YELLOW}make run${RESET}"
    echo ""
    exit 1
fi

echo -e "${GREEN}✓${RESET} Server is running"
echo ""

# ============================================================================
# Step 1: Bootstrap Admin (if not already done)
# ============================================================================

echo -e "${BLUE}[1/4]${RESET} Bootstrapping admin account..."
BOOTSTRAP_RESPONSE=$(curl -s -X POST "$SERVER_URL/api/v1/bootstrap/init" \
    -H "Content-Type: application/json" 2>&1)

if echo "$BOOTSTRAP_RESPONSE" | grep -q "admin_token"; then
    ADMIN_TOKEN=$(echo "$BOOTSTRAP_RESPONSE" | grep -o '"admin_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Admin account created"
elif echo "$BOOTSTRAP_RESPONSE" | grep -q "already initialized"; then
    echo -e "${YELLOW}⚠${RESET}  Admin already exists (skipping)"
    echo ""
    echo -e "${RED}✗${RESET} Cannot seed without fresh admin token"
    echo ""
    echo "To get a fresh environment:"
    echo -e "  ${YELLOW}make fresh${RESET}   - Complete fresh start"
    echo ""
    exit 1
else
    echo -e "${RED}✗${RESET} Failed to bootstrap admin"
    echo "Response: $BOOTSTRAP_RESPONSE"
    exit 1
fi

echo "  Token: ${ADMIN_TOKEN:0:20}..."
echo ""

# ============================================================================
# Step 2: Create Teachers
# ============================================================================

echo -e "${BLUE}[2/4]${RESET} Creating sample teachers..."

# Teacher 1: Dr. Anna Schmidt (Computer Science)
TEACHER1=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
        "email": "anna.schmidt@fits.example.com",
        "first_name": "Anna",
        "last_name": "Schmidt",
        "role": "teacher",
        "department": "Computer Science"
    }')

if echo "$TEACHER1" | grep -q "invitation_token"; then
    TEACHER1_TOKEN=$(echo "$TEACHER1" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Dr. Anna Schmidt (Computer Science)"

    # Complete teacher 1 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$TEACHER1_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "anna.schmidt",
            "password": "SecurePassword123!"
        }' > /dev/null
fi

# Teacher 2: Prof. Michael Weber (Mathematics)
TEACHER2=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
        "email": "michael.weber@fits.example.com",
        "first_name": "Michael",
        "last_name": "Weber",
        "role": "teacher",
        "department": "Mathematics"
    }')

if echo "$TEACHER2" | grep -q "invitation_token"; then
    TEACHER2_TOKEN=$(echo "$TEACHER2" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Prof. Michael Weber (Mathematics)"

    # Complete teacher 2 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$TEACHER2_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "michael.weber",
            "password": "SecurePassword123!"
        }' > /dev/null
fi

# Teacher 3: Dr. Sarah Müller (Physics)
TEACHER3=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d '{
        "email": "sarah.mueller@fits.example.com",
        "first_name": "Sarah",
        "last_name": "Müller",
        "role": "teacher",
        "department": "Physics"
    }')

if echo "$TEACHER3" | grep -q "invitation_token"; then
    TEACHER3_TOKEN=$(echo "$TEACHER3" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Dr. Sarah Müller (Physics)"

    # Complete teacher 3 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$TEACHER3_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "sarah.mueller",
            "password": "SecurePassword123!"
        }' > /dev/null
fi

echo ""

# ============================================================================
# Step 3: Login as teacher to get teacher UUID for students
# ============================================================================

echo -e "${BLUE}[3/4]${RESET} Getting teacher UUID for student assignments..."

# Login as first teacher
TEACHER1_LOGIN=$(curl -s -X POST "$SERVER_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "anna.schmidt",
        "password": "SecurePassword123!"
    }')

TEACHER1_UUID=$(echo "$TEACHER1_LOGIN" | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TEACHER1_UUID" ]; then
    echo -e "${GREEN}✓${RESET} Teacher UUID: $TEACHER1_UUID"
else
    echo -e "${RED}✗${RESET} Could not get teacher UUID"
    TEACHER1_UUID="00000000-0000-0000-0000-000000000000"  # Fallback
fi

echo ""

# ============================================================================
# Step 4: Create Student Invitations
# ============================================================================

echo -e "${BLUE}[4/4]${RESET} Creating sample student invitations..."

# Student 1: Max Mustermann
STUDENT1=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
        \"email\": \"max.mustermann@fits.example.com\",
        \"first_name\": \"Max\",
        \"last_name\": \"Mustermann\",
        \"role\": \"student\",
        \"teacher_uuid\": \"$TEACHER1_UUID\"
    }")

if echo "$STUDENT1" | grep -q "invitation_token"; then
    STUDENT1_TOKEN=$(echo "$STUDENT1" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Max Mustermann (assigned to Anna Schmidt)"

    # Complete student 1 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$STUDENT1_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "max.mustermann",
            "password": "StudentPass123!"
        }' > /dev/null
fi

# Student 2: Lisa Schneider
STUDENT2=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
        \"email\": \"lisa.schneider@fits.example.com\",
        \"first_name\": \"Lisa\",
        \"last_name\": \"Schneider\",
        \"role\": \"student\",
        \"teacher_uuid\": \"$TEACHER1_UUID\"
    }")

if echo "$STUDENT2" | grep -q "invitation_token"; then
    STUDENT2_TOKEN=$(echo "$STUDENT2" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Lisa Schneider (assigned to Anna Schmidt)"

    # Complete student 2 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$STUDENT2_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "lisa.schneider",
            "password": "StudentPass123!"
        }' > /dev/null
fi

# Student 3: Tom Fischer
STUDENT3=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
        \"email\": \"tom.fischer@fits.example.com\",
        \"first_name\": \"Tom\",
        \"last_name\": \"Fischer\",
        \"role\": \"student\",
        \"teacher_uuid\": \"$TEACHER1_UUID\"
    }")

if echo "$STUDENT3" | grep -q "invitation_token"; then
    STUDENT3_TOKEN=$(echo "$STUDENT3" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Tom Fischer (assigned to Anna Schmidt)"

    # Complete student 3 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$STUDENT3_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "tom.fischer",
            "password": "StudentPass123!"
        }' > /dev/null
fi

# Student 4: Emma Wagner
STUDENT4=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
        \"email\": \"emma.wagner@fits.example.com\",
        \"first_name\": \"Emma\",
        \"last_name\": \"Wagner\",
        \"role\": \"student\",
        \"teacher_uuid\": \"$TEACHER1_UUID\"
    }")

if echo "$STUDENT4" | grep -q "invitation_token"; then
    STUDENT4_TOKEN=$(echo "$STUDENT4" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Emma Wagner (assigned to Anna Schmidt)"

    # Complete student 4 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$STUDENT4_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "emma.wagner",
            "password": "StudentPass123!"
        }' > /dev/null
fi

# Student 5: Paul Becker
STUDENT5=$(curl -s -X POST "$SERVER_URL/api/v1/admin/invite" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -d "{
        \"email\": \"paul.becker@fits.example.com\",
        \"first_name\": \"Paul\",
        \"last_name\": \"Becker\",
        \"role\": \"student\",
        \"teacher_uuid\": \"$TEACHER1_UUID\"
    }")

if echo "$STUDENT5" | grep -q "invitation_token"; then
    STUDENT5_TOKEN=$(echo "$STUDENT5" | grep -o '"invitation_token":"[^"]*"' | cut -d'"' -f4)
    echo -e "${GREEN}✓${RESET} Paul Becker (assigned to Anna Schmidt)"

    # Complete student 5 invitation
    curl -s -X POST "$SERVER_URL/api/v1/invite/$STUDENT5_TOKEN/complete" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "paul.becker",
            "password": "StudentPass123!"
        }' > /dev/null
fi

echo ""

# ============================================================================
# Summary
# ============================================================================

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${GREEN}✓ Database seeded successfully!${RESET}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${CYAN}Created Accounts:${RESET}"
echo ""
echo -e "${YELLOW}Admin:${RESET}"
echo "  • Uses bootstrap token (from /api/v1/bootstrap/init)"
echo ""
echo -e "${YELLOW}Teachers (3):${RESET}"
echo "  • anna.schmidt     / SecurePassword123!  (Computer Science)"
echo "  • michael.weber    / SecurePassword123!  (Mathematics)"
echo "  • sarah.mueller    / SecurePassword123!  (Physics)"
echo ""
echo -e "${YELLOW}Students (5):${RESET}"
echo "  • max.mustermann   / StudentPass123!"
echo "  • lisa.schneider   / StudentPass123!"
echo "  • tom.fischer      / StudentPass123!"
echo "  • emma.wagner      / StudentPass123!"
echo "  • paul.becker      / StudentPass123!"
echo ""
echo -e "${CYAN}Test the API:${RESET}"
echo -e "  ${YELLOW}make docs-serve${RESET}  - Open Swagger UI"
echo -e "  ${YELLOW}make test-integration${RESET}  - Run integration tests"
echo ""
echo -e "${CYAN}Login Example:${RESET}"
echo '  curl -X POST http://localhost:8080/api/v1/auth/login \'
echo '    -H "Content-Type: application/json" \'
echo '    -d '\''{"username":"anna.schmidt","password":"SecurePassword123!"}'\'''
echo ""
