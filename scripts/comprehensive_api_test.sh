#!/bin/bash

# FITS Backend - Comprehensive API Test Suite
# Systematically tests all endpoints and documents issues

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
WARNINGS=0

# Problem lists
declare -a PROBLEMS=()
declare -a WARNINGS_LIST=()

# Base URL
BASE_URL="http://localhost:8080"

# Test data
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_STUDENT_EMAIL="student_$(date +%s)@example.com"
TEST_TEACHER_EMAIL="teacher_$(date +%s)@example.com"

# Variables for test flow
ADMIN_TOKEN=""
INVITE_TOKEN=""
STUDENT_UUID=""
TEACHER_UUID=""
ACCESS_TOKEN=""
REFRESH_TOKEN=""

# Helper functions
print_header() {
    echo -e "\n${BLUE}================================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}================================================================${NC}\n"
}

print_test() {
    echo -e "${YELLOW}[TEST $(($TOTAL_TESTS + 1))]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

print_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    PROBLEMS+=("$1")
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
    WARNINGS_LIST+=("$1")
    ((WARNINGS++))
}

check_response() {
    local response="$1"
    local expected_status="$2"
    local test_name="$3"

    local status=$(echo "$response" | jq -r '.status // empty')
    local http_code=$(echo "$response" | jq -r '.http_code // empty')

    if [[ -z "$http_code" ]]; then
        http_code=$(echo "$response" | head -n1)
    fi

    if [[ "$http_code" == "$expected_status" ]] || [[ "$status" == "ok" && "$expected_status" == "200" ]]; then
        print_success "$test_name"
        return 0
    else
        print_error "$test_name - Expected $expected_status, got $http_code"
        echo "Response: $response" | head -n5
        return 1
    fi
}

validate_json() {
    local json="$1"
    if echo "$json" | jq . >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# =============================================================================
# Test-Suites
# =============================================================================

test_health_and_system() {
    print_header "1. System Health & Monitoring Tests"

    print_test "Health Check"
    local response=$(curl -s "$BASE_URL/health")
    if echo "$response" | jq -e '.status == "ok"' >/dev/null; then
        print_success "Health check - Server is healthy"
    else
        print_error "Health check failed"
    fi

    print_test "Swagger UI available"
    local swagger_response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/swagger/index.html")
    if [[ "$swagger_response" == "200" ]]; then
        print_success "Swagger UI accessible"
    else
        print_error "Swagger UI not accessible (HTTP $swagger_response)"
    fi
}

test_authentication_flow() {
    print_header "2. Authentication & Authorization Tests"

    # Bootstrap Admin
    print_test "Bootstrap Admin (POST /api/v1/bootstrap/init)"
    local bootstrap_response=$(curl -s -X POST "$BASE_URL/api/v1/bootstrap/init")

    if echo "$bootstrap_response" | jq -e '.success == true' >/dev/null 2>&1; then
        ADMIN_TOKEN=$(echo "$bootstrap_response" | jq -r '.data.admin_token')
        print_success "Admin bootstrap successful"

        # Validate token format
        if [[ -z "$ADMIN_TOKEN" || "$ADMIN_TOKEN" == "null" ]]; then
            print_error "Admin token is empty or null"
        else
            print_success "Admin token received (length: ${#ADMIN_TOKEN})"
        fi
    else
        print_error "Admin bootstrap failed"
        echo "Response: $bootstrap_response"
    fi

    # Test Login without credentials (should fail)
    print_test "Login without credentials (should fail)"
    local invalid_login=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d '{}')

    if echo "$invalid_login" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Invalid login correctly rejected"
    else
        print_warning "Invalid login should return error"
    fi

    # Test Token Refresh without token
    print_test "Token refresh without refresh token"
    local invalid_refresh=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
        -H "Content-Type: application/json" \
        -d '{}')

    if echo "$invalid_refresh" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Invalid refresh correctly rejected"
    else
        print_warning "Invalid refresh should return error"
    fi

    # Test Protected Endpoint without auth
    print_test "Protected endpoint without authorization header"
    local unauth_response=$(curl -s -X POST "$BASE_URL/api/v1/auth/logout")

    if echo "$unauth_response" | grep -q "unauthorized\|Unauthorized\|401"; then
        print_success "Unauthorized access correctly blocked"
    else
        print_warning "Unauthorized access should return 401"
    fi
}

test_invitation_system() {
    print_header "3. Invitation System Tests"

    if [[ -z "$ADMIN_TOKEN" ]]; then
        print_error "Admin token not available - skipping invitation tests"
        return
    fi

    # Create Student Invitation
    print_test "Create Student Invitation (Admin only)"
    local invite_response=$(curl -s -X POST "$BASE_URL/api/v1/admin/invite" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$TEST_EMAIL\",
            \"first_name\": \"Test\",
            \"last_name\": \"Student\",
            \"role\": \"student\"
        }")

    if echo "$invite_response" | jq -e '.success == true' >/dev/null 2>&1; then
        INVITE_TOKEN=$(echo "$invite_response" | jq -r '.data.token')
        print_success "Student invitation created"

        # Validate invitation link
        local invite_link=$(echo "$invite_response" | jq -r '.data.link')
        if [[ -n "$invite_link" && "$invite_link" != "null" ]]; then
            print_success "Invitation link generated: $invite_link"
        else
            print_warning "Invitation link is empty"
        fi
    else
        print_error "Failed to create invitation"
        echo "Response: $invite_response"
    fi

    # Get Invitation Details
    if [[ -n "$INVITE_TOKEN" ]]; then
        print_test "Get Invitation Details (GET /api/v1/invite/:token)"
        local invite_details=$(curl -s "$BASE_URL/api/v1/invite/$INVITE_TOKEN")

        if echo "$invite_details" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Invitation details retrieved"

            # Validate fields
            local email=$(echo "$invite_details" | jq -r '.data.email')
            local role=$(echo "$invite_details" | jq -r '.data.role')

            if [[ "$email" == "$TEST_EMAIL" ]]; then
                print_success "Email matches"
            else
                print_error "Email mismatch: expected $TEST_EMAIL, got $email"
            fi

            if [[ "$role" == "student" ]]; then
                print_success "Role matches"
            else
                print_error "Role mismatch: expected student, got $role"
            fi
        else
            print_error "Failed to get invitation details"
        fi
    fi

    # Test Invalid Token
    print_test "Get Invitation with invalid token"
    local invalid_invite=$(curl -s "$BASE_URL/api/v1/invite/invalid-token-123")

    if echo "$invalid_invite" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Invalid invitation token correctly rejected"
    else
        print_warning "Invalid token should return error"
    fi

    # Complete Invitation (Register)
    if [[ -n "$INVITE_TOKEN" ]]; then
        print_test "Complete Invitation (Register new user)"
        local complete_response=$(curl -s -X POST "$BASE_URL/api/v1/invite/$INVITE_TOKEN/complete" \
            -H "Content-Type: application/json" \
            -d "{
                \"username\": \"testuser_$(date +%s)\",
                \"password\": \"SecurePassword123!\"
            }")

        if echo "$complete_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "User registration completed"

            # Extract tokens
            ACCESS_TOKEN=$(echo "$complete_response" | jq -r '.data.access_token')
            REFRESH_TOKEN=$(echo "$complete_response" | jq -r '.data.refresh_token')

            if [[ -n "$ACCESS_TOKEN" && "$ACCESS_TOKEN" != "null" ]]; then
                print_success "Access token received"
            else
                print_error "Access token missing"
            fi
        else
            print_error "Failed to complete invitation"
            echo "Response: $complete_response"
        fi
    fi
}

test_student_management() {
    print_header "4. Student Management Tests"

    if [[ -z "$ADMIN_TOKEN" ]]; then
        print_error "Admin token not available - skipping student tests"
        return
    fi

    # Create Student
    print_test "Create Student (PUT /api/v1/student)"
    local create_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"Max\",
            \"last_name\": \"Mustermann\",
            \"email\": \"$TEST_STUDENT_EMAIL\"
        }")

    if echo "$create_response" | jq -e '.success == true' >/dev/null 2>&1; then
        STUDENT_UUID=$(echo "$create_response" | jq -r '.data.uuid')
        print_success "Student created (UUID: $STUDENT_UUID)"
    else
        print_error "Failed to create student"
        echo "Response: $create_response"
    fi

    # Get Student
    if [[ -n "$STUDENT_UUID" ]]; then
        print_test "Get Student by UUID (GET /api/v1/student/:uuid)"
        local get_response=$(curl -s "$BASE_URL/api/v1/student/$STUDENT_UUID")

        if echo "$get_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Student retrieved"

            # Validate data
            local email=$(echo "$get_response" | jq -r '.data.email')
            if [[ "$email" == "$TEST_STUDENT_EMAIL" ]]; then
                print_success "Student data correct"
            else
                print_error "Student data mismatch"
            fi
        else
            print_error "Failed to get student"
        fi
    fi

    # List Students
    print_test "List All Students (GET /api/v1/student)"
    local list_response=$(curl -s "$BASE_URL/api/v1/student")

    if echo "$list_response" | jq -e '.success == true' >/dev/null 2>&1; then
        local count=$(echo "$list_response" | jq '.data | length')
        print_success "Student list retrieved ($count students)"

        # Check pagination
        if echo "$list_response" | jq -e '.pagination' >/dev/null 2>&1; then
            print_success "Pagination metadata present"
        else
            print_warning "No pagination metadata (KNOWN ISSUE)"
        fi
    else
        print_error "Failed to list students"
    fi

    # Update Student
    if [[ -n "$STUDENT_UUID" ]]; then
        print_test "Update Student (POST /api/v1/student/:uuid)"
        local update_response=$(curl -s -X POST "$BASE_URL/api/v1/student/$STUDENT_UUID" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"first_name\": \"Updated\",
                \"last_name\": \"Name\"
            }")

        if echo "$update_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Student updated"
        else
            print_error "Failed to update student"
        fi
    fi

    # Test Duplicate Email
    print_test "Create Student with duplicate email (should fail)"
    local duplicate_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"Another\",
            \"last_name\": \"Student\",
            \"email\": \"$TEST_STUDENT_EMAIL\"
        }")

    if echo "$duplicate_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Duplicate email correctly rejected"
    else
        print_error "Duplicate email should be rejected"
    fi
}

test_teacher_management() {
    print_header "5. Teacher Management Tests"

    if [[ -z "$ADMIN_TOKEN" ]]; then
        print_error "Admin token not available - skipping teacher tests"
        return
    fi

    # Create Teacher
    print_test "Create Teacher (POST /api/v1/teacher)"
    local create_response=$(curl -s -X POST "$BASE_URL/api/v1/teacher" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"Anna\",
            \"last_name\": \"Schmidt\",
            \"email\": \"$TEST_TEACHER_EMAIL\",
            \"department\": \"Computer Science\"
        }")

    if echo "$create_response" | jq -e '.success == true' >/dev/null 2>&1; then
        TEACHER_UUID=$(echo "$create_response" | jq -r '.data.uuid')
        print_success "Teacher created (UUID: $TEACHER_UUID)"
    else
        print_error "Failed to create teacher"
        echo "Response: $create_response"
    fi

    # Get Teacher
    if [[ -n "$TEACHER_UUID" ]]; then
        print_test "Get Teacher by UUID (GET /api/v1/teacher/:uuid)"
        local get_response=$(curl -s "$BASE_URL/api/v1/teacher/$TEACHER_UUID")

        if echo "$get_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Teacher retrieved"
        else
            print_error "Failed to get teacher"
        fi
    fi

    # List Teachers
    print_test "List All Teachers (GET /api/v1/teacher)"
    local list_response=$(curl -s "$BASE_URL/api/v1/teacher")

    if echo "$list_response" | jq -e '.success == true' >/dev/null 2>&1; then
        local count=$(echo "$list_response" | jq '.data | length')
        print_success "Teacher list retrieved ($count teachers)"

        # Check pagination
        if echo "$list_response" | jq -e '.pagination' >/dev/null 2>&1; then
            print_success "Pagination metadata present"
        else
            print_warning "No pagination metadata (KNOWN ISSUE)"
        fi
    else
        print_error "Failed to list teachers"
    fi

    # Update Teacher
    if [[ -n "$TEACHER_UUID" ]]; then
        print_test "Update Teacher (POST /api/v1/teacher/:uuid)"
        local update_response=$(curl -s -X POST "$BASE_URL/api/v1/teacher/$TEACHER_UUID" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"department\": \"Mathematics\"
            }")

        if echo "$update_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Teacher updated"
        else
            print_error "Failed to update teacher"
        fi
    fi

    # Test HTTP Method Inconsistency
    print_warning "HTTP Method Inconsistency: Student uses PUT, Teacher uses POST for creation"
}

test_signing_endpoints() {
    print_header "6. Digital Signing Tests"

    if [[ -z "$ACCESS_TOKEN" && -z "$ADMIN_TOKEN" ]]; then
        print_error "No auth token available - skipping signing tests"
        return
    fi

    local token="${ACCESS_TOKEN:-$ADMIN_TOKEN}"

    # Test Upload
    print_test "Upload Report (POST /api/v1/signing/upload)"
    echo "test data" > /tmp/test_report.txt
    local upload_response=$(curl -s -X POST "$BASE_URL/api/v1/signing/upload" \
        -H "Authorization: Bearer $token" \
        -F "file=@/tmp/test_report.txt")

    if echo "$upload_response" | jq -e '.code == 501' >/dev/null 2>&1; then
        print_warning "Upload endpoint returns 501 Not Implemented (STUB)"
    elif echo "$upload_response" | jq -e '.success == true' >/dev/null 2>&1; then
        print_success "Upload successful"
    else
        print_error "Upload failed unexpectedly"
        echo "Response: $upload_response"
    fi

    # Test Get Sign Requests
    print_test "Get Sign Requests (GET /api/v1/signing/sign_requests)"
    local requests_response=$(curl -s "$BASE_URL/api/v1/signing/sign_requests" \
        -H "Authorization: Bearer $token")

    if echo "$requests_response" | jq -e '.code == 501' >/dev/null 2>&1; then
        print_warning "Sign requests endpoint returns 501 Not Implemented (STUB)"
    elif echo "$requests_response" | jq -e '.success == true' >/dev/null 2>&1; then
        print_success "Sign requests retrieved"
    else
        print_error "Failed to get sign requests"
    fi

    # Test Sign Uploads
    print_test "Upload Signed Requests (POST /api/v1/signing/sign_uploads)"
    local sign_response=$(curl -s -X POST "$BASE_URL/api/v1/signing/sign_uploads" \
        -H "Authorization: Bearer $token" \
        -F "file=@/tmp/test_report.txt")

    if echo "$sign_response" | jq -e '.code == 501' >/dev/null 2>&1; then
        print_warning "Sign uploads endpoint returns 501 Not Implemented (STUB)"
    elif echo "$sign_response" | jq -e '.success == true' >/dev/null 2>&1; then
        print_success "Signed uploads processed"
    else
        print_error "Failed to process signed uploads"
    fi

    rm -f /tmp/test_report.txt
}

test_input_validation() {
    print_header "7. Input Validation & Security Tests"

    if [[ -z "$ADMIN_TOKEN" ]]; then
        print_error "Admin token not available - skipping validation tests"
        return
    fi

    # Test XSS in name fields
    print_test "XSS Prevention in Student Name"
    local xss_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"<script>alert('xss')</script>\",
            \"last_name\": \"Test\",
            \"email\": \"xss_$(date +%s)@test.com\"
        }")

    if echo "$xss_response" | jq -e '.data.first_name' | grep -q "<script>"; then
        print_warning "XSS characters not sanitized in first_name"
    else
        print_success "XSS characters handled"
    fi

    # Test SQL Injection
    print_test "SQL Injection Prevention"
    local sql_response=$(curl -s "$BASE_URL/api/v1/student/'; DROP TABLE students; --")

    if echo "$sql_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "SQL injection prevented (invalid UUID)"
    else
        print_warning "Unexpected response to SQL injection attempt"
    fi

    # Test Very Long Strings
    print_test "Long String Handling"
    local long_string=$(python3 -c "print('a' * 10000)")
    local long_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"$long_string\",
            \"last_name\": \"Test\",
            \"email\": \"long_$(date +%s)@test.com\"
        }")

    if echo "$long_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Long strings rejected"
    else
        print_warning "Long strings should be rejected"
    fi

    # Test Invalid Email Format
    print_test "Email Validation"
    local invalid_email_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\": \"Test\",
            \"last_name\": \"User\",
            \"email\": \"not-an-email\"
        }")

    if echo "$invalid_email_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Invalid email format rejected"
    else
        print_error "Invalid email should be rejected"
    fi
}

test_error_handling() {
    print_header "8. Error Handling Tests"

    # Test 404 - Non-existent endpoint
    print_test "404 Not Found"
    local not_found=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/v1/nonexistent")
    if [[ "$not_found" == "404" ]]; then
        print_success "404 correctly returned for non-existent endpoint"
    else
        print_warning "Expected 404, got $not_found"
    fi

    # Test Non-existent UUID
    print_test "Get Non-existent Student"
    local fake_uuid="00000000-0000-0000-0000-000000000000"
    local not_found_response=$(curl -s "$BASE_URL/api/v1/student/$fake_uuid")

    if echo "$not_found_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Non-existent student correctly returns error"
    else
        print_warning "Non-existent resource should return error"
    fi

    # Test Malformed JSON
    print_test "Malformed JSON Handling"
    local malformed_response=$(curl -s -X PUT "$BASE_URL/api/v1/student" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{invalid json")

    if echo "$malformed_response" | jq -e '.success == false' >/dev/null 2>&1; then
        print_success "Malformed JSON correctly rejected"
    else
        print_warning "Malformed JSON should return error"
    fi
}

test_database_validation() {
    print_header "9. Database Validation Tests"

    print_test "Check if created students persist in database"

    # This would require direct DB access, which we'll skip for now
    # but it's an important test to add later
    print_warning "Database validation tests require direct DB access - skipped"
    print_warning "Recommendation: Add integration tests with DB queries"
}

cleanup_test_data() {
    print_header "10. Cleanup Test Data"

    if [[ -z "$ADMIN_TOKEN" ]]; then
        print_warning "Admin token not available - skipping cleanup"
        return
    fi

    # Delete Test Student
    if [[ -n "$STUDENT_UUID" ]]; then
        print_test "Delete Test Student"
        local delete_response=$(curl -s -X DELETE "$BASE_URL/api/v1/student/$STUDENT_UUID" \
            -H "Authorization: Bearer $ADMIN_TOKEN")

        if echo "$delete_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Test student deleted"
        else
            print_warning "Failed to delete test student"
        fi
    fi

    # Delete Test Teacher
    if [[ -n "$TEACHER_UUID" ]]; then
        print_test "Delete Test Teacher"
        local delete_response=$(curl -s -X DELETE "$BASE_URL/api/v1/teacher/$TEACHER_UUID" \
            -H "Authorization: Bearer $ADMIN_TOKEN")

        if echo "$delete_response" | jq -e '.success == true' >/dev/null 2>&1; then
            print_success "Test teacher deleted"
        else
            print_warning "Failed to delete test teacher"
        fi
    fi
}

# =============================================================================
# Test Summary & Report
# =============================================================================

print_summary() {
    print_header "TEST SUMMARY"

    echo -e "${BLUE}Total Tests:${NC}    $TOTAL_TESTS"
    echo -e "${GREEN}Passed:${NC}         $PASSED_TESTS"
    echo -e "${RED}Failed:${NC}         $FAILED_TESTS"
    echo -e "${YELLOW}Warnings:${NC}       $WARNINGS"

    local success_rate=0
    if [[ $TOTAL_TESTS -gt 0 ]]; then
        success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    fi
    echo -e "${BLUE}Success Rate:${NC}   $success_rate%"

    if [[ ${#PROBLEMS[@]} -gt 0 ]]; then
        echo -e "\n${RED}=== CRITICAL PROBLEMS ===${NC}"
        for problem in "${PROBLEMS[@]}"; do
            echo -e "${RED}[FAIL]${NC} $problem"
        done
    fi

    if [[ ${#WARNINGS_LIST[@]} -gt 0 ]]; then
        echo -e "\n${YELLOW}=== WARNINGS & KNOWN ISSUES ===${NC}"
        for warning in "${WARNINGS_LIST[@]}"; do
            echo -e "${YELLOW}[WARN]${NC} $warning"
        done
    fi
}

# =============================================================================
# Main Execution
# =============================================================================

main() {
    echo -e "${BLUE}"
    echo "================================================================"
    echo "  FITS Backend - Comprehensive API Test Suite"
    echo "  Version: 1.0"
    echo "  Date: $(date)"
    echo "================================================================"
    echo -e "${NC}"

    # Check dependencies
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}[ERROR] jq is not installed. Please install jq first.${NC}"
        exit 1
    fi

    if ! command -v curl &> /dev/null; then
        echo -e "${RED}[ERROR] curl is not installed. Please install curl first.${NC}"
        exit 1
    fi

    # Check if server is running
    if ! curl -s "$BASE_URL/health" >/dev/null 2>&1; then
        echo -e "${RED}[ERROR] Server is not running on $BASE_URL${NC}"
        echo "Please start the server with: make run"
        exit 1
    fi

    # Run all test suites
    test_health_and_system
    test_authentication_flow
    test_invitation_system
    test_student_management
    test_teacher_management
    test_signing_endpoints
    test_input_validation
    test_error_handling
    test_database_validation
    cleanup_test_data

    # Print summary
    print_summary

    # Exit code based on results
    if [[ $FAILED_TESTS -gt 0 ]]; then
        exit 1
    else
        exit 0
    fi
}

main
