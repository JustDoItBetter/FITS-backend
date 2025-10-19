#!/bin/bash

# FITS Backend - Complete End-to-End Test
# Tests complete auth flow from bootstrap to logout

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}        FITS Backend - End-to-End Test                          ${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

API_URL="${API_URL:-http://localhost:8080}"

# Check if server is running
echo -e "${YELLOW}[1/12]${NC} Checking if server is running..."
if ! curl -s "$API_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}[FAIL]${NC} Server is not running on $API_URL"
    echo ""
    echo "Please start the server:"
    echo "  go run cmd/server/main.go"
    echo ""
    exit 1
fi
echo -e "${GREEN}[PASS]${NC} Server is running"
echo ""

# 1. Bootstrap Admin
echo -e "${YELLOW}[2/12]${NC} Bootstrap admin..."
BOOTSTRAP_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/bootstrap/init")

if echo "$BOOTSTRAP_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    ADMIN_TOKEN=$(echo "$BOOTSTRAP_RESPONSE" | jq -r '.data.admin_token')
    echo -e "${GREEN}[PASS]${NC} Admin initialized"
    echo "   Token: ${ADMIN_TOKEN:0:50}..."
else
    ERROR=$(echo "$BOOTSTRAP_RESPONSE" | jq -r '.details // .error')
    if [[ "$ERROR" == *"already initialized"* ]]; then
        echo -e "${YELLOW}[WARN]${NC} Admin already initialized"
        echo ""
        echo "Run to reset:"
        echo "  ./scripts/reset_db_only.sh"
        echo ""
        exit 1
    else
        echo -e "${RED}[FAIL]${NC} Bootstrap failed: $ERROR"
        exit 1
    fi
fi
echo ""

# 2. Create Student Invitation
echo -e "${YELLOW}[3/12]${NC} Creating student invitation..."
STUDENT_INVITE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/admin/invite" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test.student@example.com",
        "first_name": "Test",
        "last_name": "Student",
        "role": "student"
    }')

if echo "$STUDENT_INVITE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    STUDENT_INVITE_TOKEN=$(echo "$STUDENT_INVITE_RESPONSE" | jq -r '.data.invitation_token')
    echo -e "${GREEN}[PASS]${NC} Student invitation created"
    echo "   Email: test.student@example.com"
else
    echo -e "${RED}[FAIL]${NC} Student invitation failed"
    echo "$STUDENT_INVITE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 3. Create Teacher Invitation
echo -e "${YELLOW}[4/12]${NC} Creating teacher invitation..."
TEACHER_INVITE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/admin/invite" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test.teacher@example.com",
        "first_name": "Test",
        "last_name": "Teacher",
        "role": "teacher",
        "department": "Computer Science"
    }')

if echo "$TEACHER_INVITE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    TEACHER_INVITE_TOKEN=$(echo "$TEACHER_INVITE_RESPONSE" | jq -r '.data.invitation_token')
    echo -e "${GREEN}[PASS]${NC} Teacher invitation created"
    echo "   Email: test.teacher@example.com"
else
    echo -e "${RED}[FAIL]${NC} Teacher invitation failed"
    echo "$TEACHER_INVITE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 4. Register Student
echo -e "${YELLOW}[5/12]${NC} Registering student..."
STUDENT_REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"invitation_token\": \"$STUDENT_INVITE_TOKEN\",
        \"username\": \"test.student\",
        \"password\": \"StudentPass123!\"
    }")

if echo "$STUDENT_REGISTER_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    STUDENT_ACCESS_TOKEN=$(echo "$STUDENT_REGISTER_RESPONSE" | jq -r '.data.access_token')
    STUDENT_REFRESH_TOKEN=$(echo "$STUDENT_REGISTER_RESPONSE" | jq -r '.data.refresh_token')
    echo -e "${GREEN}[PASS]${NC} Student registered"
    echo "   Username: test.student"
else
    echo -e "${RED}[FAIL]${NC} Student registration failed"
    echo "$STUDENT_REGISTER_RESPONSE" | jq .
    exit 1
fi
echo ""

# 5. Register Teacher
echo -e "${YELLOW}[6/12]${NC} Registering teacher..."
TEACHER_REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"invitation_token\": \"$TEACHER_INVITE_TOKEN\",
        \"username\": \"test.teacher\",
        \"password\": \"TeacherPass123!\"
    }")

if echo "$TEACHER_REGISTER_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    TEACHER_ACCESS_TOKEN=$(echo "$TEACHER_REGISTER_RESPONSE" | jq -r '.data.access_token')
    TEACHER_REFRESH_TOKEN=$(echo "$TEACHER_REGISTER_RESPONSE" | jq -r '.data.refresh_token')
    echo -e "${GREEN}[PASS]${NC} Teacher registered"
    echo "   Username: test.teacher"
else
    echo -e "${RED}[FAIL]${NC} Teacher registration failed"
    echo "$TEACHER_REGISTER_RESPONSE" | jq .
    exit 1
fi
echo ""

# 6. Student Login
echo -e "${YELLOW}[7/12]${NC} Student login..."
STUDENT_LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "test.student",
        "password": "StudentPass123!"
    }')

if echo "$STUDENT_LOGIN_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Student login successful"
else
    echo -e "${RED}[FAIL]${NC} Student login failed"
    echo "$STUDENT_LOGIN_RESPONSE" | jq .
    exit 1
fi
echo ""

# 7. Teacher Login
echo -e "${YELLOW}[8/12]${NC} Teacher login..."
TEACHER_LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "test.teacher",
        "password": "TeacherPass123!"
    }')

if echo "$TEACHER_LOGIN_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Teacher login successful"
else
    echo -e "${RED}[FAIL]${NC} Teacher login failed"
    echo "$TEACHER_LOGIN_RESPONSE" | jq .
    exit 1
fi
echo ""

# 8. Test wrong password
echo -e "${YELLOW}[9/12]${NC} Testing wrong password..."
WRONG_LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "username": "test.student",
        "password": "WrongPassword123!"
    }')

if echo "$WRONG_LOGIN_RESPONSE" | jq -e '.success == false' >/dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Wrong password correctly rejected"
else
    echo -e "${RED}[FAIL]${NC} Wrong password should be rejected"
    exit 1
fi
echo ""

# 9. Refresh Student Token
echo -e "${YELLOW}[10/12]${NC} Refreshing student token..."
STUDENT_REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{
        \"refresh_token\": \"$STUDENT_REFRESH_TOKEN\"
    }")

if echo "$STUDENT_REFRESH_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Token refresh successful"
else
    echo -e "${RED}[FAIL]${NC} Token refresh failed"
    echo "$STUDENT_REFRESH_RESPONSE" | jq .
    exit 1
fi
echo ""

# 10. Student Logout
echo -e "${YELLOW}[11/12]${NC} Student logout..."
STUDENT_LOGOUT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/logout" \
    -H "Authorization: Bearer $STUDENT_ACCESS_TOKEN")

if echo "$STUDENT_LOGOUT_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Student logout successful"
else
    echo -e "${RED}[FAIL]${NC} Student logout failed"
    echo "$STUDENT_LOGOUT_RESPONSE" | jq .
    exit 1
fi
echo ""

# 11. Teacher Logout
echo -e "${YELLOW}[12/12]${NC} Teacher logout..."
TEACHER_LOGOUT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/logout" \
    -H "Authorization: Bearer $TEACHER_ACCESS_TOKEN")

if echo "$TEACHER_LOGOUT_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Teacher logout successful"
else
    echo -e "${RED}[FAIL]${NC} Teacher logout failed"
    echo "$TEACHER_LOGOUT_RESPONSE" | jq .
    exit 1
fi
echo ""

# Success
echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}              All tests passed successfully!                    ${NC}"
echo -e "${GREEN}================================================================${NC}"
echo ""
echo -e "${BLUE}Test Summary:${NC}"
echo ""
echo "  [PASS] Admin bootstrap"
echo "  [PASS] Student invitation"
echo "  [PASS] Teacher invitation"
echo "  [PASS] Student registration"
echo "  [PASS] Teacher registration"
echo "  [PASS] Student login"
echo "  [PASS] Teacher login"
echo "  [PASS] Security (wrong password)"
echo "  [PASS] Token refresh"
echo "  [PASS] Student logout"
echo "  [PASS] Teacher logout"
echo ""
