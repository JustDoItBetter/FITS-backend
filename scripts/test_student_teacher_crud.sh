#!/bin/bash

# FITS Backend - Student & Teacher CRUD End-to-End Test
# Tests complete CRUD operations for Student and Teacher management

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}     FITS Backend - Student/Teacher CRUD End-to-End Test       ${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

API_URL="${API_URL:-http://localhost:8080}"

# Check if server is running
echo -e "${YELLOW}[1/15]${NC} Checking if server is running..."
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
echo -e "${YELLOW}[2/15]${NC} Bootstrap admin..."
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

# 2. Create Student
echo -e "${YELLOW}[3/15]${NC} Creating student..."
STUDENT_CREATE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/student" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "first_name": "Max",
        "last_name": "Mustermann",
        "email": "max.mustermann@test.com"
    }')

if echo "$STUDENT_CREATE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    STUDENT_UUID=$(echo "$STUDENT_CREATE_RESPONSE" | jq -r '.data.uuid')
    echo -e "${GREEN}[PASS]${NC} Student created"
    echo "   UUID: $STUDENT_UUID"
    echo "   Name: Max Mustermann"
else
    echo -e "${RED}[FAIL]${NC} Student creation failed"
    echo "$STUDENT_CREATE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 3. Get Student by UUID
echo -e "${YELLOW}[4/15]${NC} Getting student by UUID..."
STUDENT_GET_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/student/$STUDENT_UUID")

if echo "$STUDENT_GET_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    STUDENT_NAME=$(echo "$STUDENT_GET_RESPONSE" | jq -r '.data.first_name + " " + .data.last_name')
    echo -e "${GREEN}[PASS]${NC} Student retrieved"
    echo "   Name: $STUDENT_NAME"
else
    echo -e "${RED}[FAIL]${NC} Student retrieval failed"
    echo "$STUDENT_GET_RESPONSE" | jq .
    exit 1
fi
echo ""

# 4. Update Student
echo -e "${YELLOW}[5/15]${NC} Updating student..."
STUDENT_UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/api/v1/student/$STUDENT_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "first_name": "Maximilian",
        "email": "maximilian.mustermann@test.com"
    }')

if echo "$STUDENT_UPDATE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Student updated"
    echo "   New name: Maximilian Mustermann"
else
    echo -e "${RED}[FAIL]${NC} Student update failed"
    echo "$STUDENT_UPDATE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 5. List Students
echo -e "${YELLOW}[6/15]${NC} Listing all students..."
STUDENT_LIST_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/student")

if echo "$STUDENT_LIST_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    STUDENT_COUNT=$(echo "$STUDENT_LIST_RESPONSE" | jq '.data | length')
    echo -e "${GREEN}[PASS]${NC} Students listed"
    echo "   Count: $STUDENT_COUNT"
else
    echo -e "${RED}[FAIL]${NC} Student list failed"
    echo "$STUDENT_LIST_RESPONSE" | jq .
    exit 1
fi
echo ""

# 6. Create Teacher
echo -e "${YELLOW}[7/15]${NC} Creating teacher..."
TEACHER_CREATE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/teacher" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "first_name": "Anna",
        "last_name": "Schmidt",
        "email": "anna.schmidt@test.com",
        "department": "Computer Science"
    }')

if echo "$TEACHER_CREATE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    TEACHER_UUID=$(echo "$TEACHER_CREATE_RESPONSE" | jq -r '.data.uuid')
    echo -e "${GREEN}[PASS]${NC} Teacher created"
    echo "   UUID: $TEACHER_UUID"
    echo "   Name: Anna Schmidt"
    echo "   Department: Computer Science"
else
    echo -e "${RED}[FAIL]${NC} Teacher creation failed"
    echo "$TEACHER_CREATE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 7. Get Teacher by UUID
echo -e "${YELLOW}[8/15]${NC} Getting teacher by UUID..."
TEACHER_GET_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/teacher/$TEACHER_UUID")

if echo "$TEACHER_GET_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    TEACHER_NAME=$(echo "$TEACHER_GET_RESPONSE" | jq -r '.data.first_name + " " + .data.last_name')
    TEACHER_DEPT=$(echo "$TEACHER_GET_RESPONSE" | jq -r '.data.department')
    echo -e "${GREEN}[PASS]${NC} Teacher retrieved"
    echo "   Name: $TEACHER_NAME"
    echo "   Department: $TEACHER_DEPT"
else
    echo -e "${RED}[FAIL]${NC} Teacher retrieval failed"
    echo "$TEACHER_GET_RESPONSE" | jq .
    exit 1
fi
echo ""

# 8. Update Teacher
echo -e "${YELLOW}[9/15]${NC} Updating teacher..."
TEACHER_UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/api/v1/teacher/$TEACHER_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "department": "Data Science"
    }')

if echo "$TEACHER_UPDATE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Teacher updated"
    echo "   New department: Data Science"
else
    echo -e "${RED}[FAIL]${NC} Teacher update failed"
    echo "$TEACHER_UPDATE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 9. List Teachers
echo -e "${YELLOW}[10/15]${NC} Listing all teachers..."
TEACHER_LIST_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/teacher")

if echo "$TEACHER_LIST_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    TEACHER_COUNT=$(echo "$TEACHER_LIST_RESPONSE" | jq '.data | length')
    echo -e "${GREEN}[PASS]${NC} Teachers listed"
    echo "   Count: $TEACHER_COUNT"
else
    echo -e "${RED}[FAIL]${NC} Teacher list failed"
    echo "$TEACHER_LIST_RESPONSE" | jq .
    exit 1
fi
echo ""

# 10. Assign Student to Teacher
echo -e "${YELLOW}[11/15]${NC} Assigning student to teacher..."
STUDENT_ASSIGN_RESPONSE=$(curl -s -X PUT "$API_URL/api/v1/student/$STUDENT_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"teacher_id\": \"$TEACHER_UUID\"
    }")

if echo "$STUDENT_ASSIGN_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Student assigned to teacher"
    echo "   Student: Max Mustermann -> Teacher: Anna Schmidt"
else
    echo -e "${RED}[FAIL]${NC} Student-teacher assignment failed"
    echo "$STUDENT_ASSIGN_RESPONSE" | jq .
    exit 1
fi
echo ""

# 11. Verify Student has Teacher
echo -e "${YELLOW}[12/15]${NC} Verifying teacher assignment..."
STUDENT_VERIFY_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/student/$STUDENT_UUID")

if echo "$STUDENT_VERIFY_RESPONSE" | jq -e '.data.teacher_id' > /dev/null 2>&1; then
    ASSIGNED_TEACHER=$(echo "$STUDENT_VERIFY_RESPONSE" | jq -r '.data.teacher_id')
    if [ "$ASSIGNED_TEACHER" == "$TEACHER_UUID" ]; then
        echo -e "${GREEN}[PASS]${NC} Teacher assignment verified"
        echo "   Teacher ID: $ASSIGNED_TEACHER"
    else
        echo -e "${RED}[FAIL]${NC} Teacher ID mismatch"
        exit 1
    fi
else
    echo -e "${RED}[FAIL]${NC} No teacher assignment found"
    exit 1
fi
echo ""

# 12. Test Duplicate Email
echo -e "${YELLOW}[13/15]${NC} Testing duplicate email validation..."
STUDENT_DUPLICATE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/student" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "first_name": "Duplicate",
        "last_name": "Test",
        "email": "maximilian.mustermann@test.com"
    }')

if echo "$STUDENT_DUPLICATE_RESPONSE" | jq -e '.success == false' > /dev/null 2>&1; then
    ERROR=$(echo "$STUDENT_DUPLICATE_RESPONSE" | jq -r '.error')
    if [[ "$ERROR" == *"already exists"* ]]; then
        echo -e "${GREEN}[PASS]${NC} Duplicate email correctly rejected"
    else
        echo -e "${RED}[FAIL]${NC} Unexpected error: $ERROR"
        exit 1
    fi
else
    echo -e "${RED}[FAIL]${NC} Duplicate email should be rejected"
    exit 1
fi
echo ""

# 13. Delete Student
echo -e "${YELLOW}[14/15]${NC} Deleting student..."
STUDENT_DELETE_RESPONSE=$(curl -s -X DELETE "$API_URL/api/v1/student/$STUDENT_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

if echo "$STUDENT_DELETE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Student deleted"
else
    echo -e "${RED}[FAIL]${NC} Student deletion failed"
    echo "$STUDENT_DELETE_RESPONSE" | jq .
    exit 1
fi
echo ""

# 14. Delete Teacher
echo -e "${YELLOW}[15/15]${NC} Deleting teacher..."
TEACHER_DELETE_RESPONSE=$(curl -s -X DELETE "$API_URL/api/v1/teacher/$TEACHER_UUID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

if echo "$TEACHER_DELETE_RESPONSE" | jq -e '.success' > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} Teacher deleted"
else
    echo -e "${RED}[FAIL]${NC} Teacher deletion failed"
    echo "$TEACHER_DELETE_RESPONSE" | jq .
    exit 1
fi
echo ""

# Success
echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}        All Student/Teacher CRUD tests passed!                 ${NC}"
echo -e "${GREEN}================================================================${NC}"
echo ""
echo -e "${BLUE}Test Summary:${NC}"
echo ""
echo "  [PASS] Admin bootstrap"
echo "  [PASS] Student create"
echo "  [PASS] Student get (UUID)"
echo "  [PASS] Student update"
echo "  [PASS] Student list"
echo "  [PASS] Teacher create"
echo "  [PASS] Teacher get (UUID)"
echo "  [PASS] Teacher update"
echo "  [PASS] Teacher list"
echo "  [PASS] Student-teacher assignment"
echo "  [PASS] Duplicate email validation"
echo "  [PASS] Student delete"
echo "  [PASS] Teacher delete"
echo ""
