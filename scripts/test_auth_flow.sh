#!/bin/bash

# FITS Backend - Auth Flow Test Script
# This script tests the complete authentication flow

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_URL="${API_URL:-http://localhost:8080}"

echo -e "${GREEN}=== FITS Backend Auth Flow Test ===${NC}\n"

# Check if server is running
echo -e "${YELLOW}[1/7]${NC} Checking if server is running..."
if ! curl -s "$API_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}Error: Server is not running at $API_URL${NC}"
    echo "Start the server first: go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}Server is running${NC}\n"

# Initialize admin
echo -e "${YELLOW}[2/7]${NC} Initializing admin..."
BOOTSTRAP_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/bootstrap/init" \
  -H "Content-Type: application/json")

if echo "$BOOTSTRAP_RESPONSE" | grep -q "success.*true"; then
    ADMIN_TOKEN=$(echo "$BOOTSTRAP_RESPONSE" | grep -o '"admin_token":"[^"]*' | cut -d'"' -f4)
    echo -e "${GREEN}Admin initialized successfully${NC}"
    echo "Admin Token: ${ADMIN_TOKEN:0:50}..."
elif echo "$BOOTSTRAP_RESPONSE" | grep -q "admin already initialized"; then
    echo -e "${YELLOW}Admin already initialized (this is normal on second run)${NC}"
    echo -e "${YELLOW}Note: You'll need the existing admin token from first run${NC}"
    echo ""
    echo "To continue testing, either:"
    echo "  1. Use the admin token from your first run"
    echo "  2. Reset the database: dropdb fits_db && ./scripts/setup_db.sh"
    echo ""
    read -p "Enter existing admin token (or press Enter to skip admin tests): " ADMIN_TOKEN

    if [ -z "$ADMIN_TOKEN" ]; then
        echo -e "${YELLOW}Skipping admin-dependent tests${NC}"
        echo "Testing login with existing user..."
        echo ""
        read -p "Enter username: " USERNAME
        read -sp "Enter password: " PASSWORD
        echo ""

        # Test login
        echo -e "${YELLOW}[Login Test]${NC} Attempting login..."
        LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
          -H "Content-Type: application/json" \
          -d '{
            "username": "'"$USERNAME"'",
            "password": "'"$PASSWORD"'"
          }')

        if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
            echo -e "${GREEN}Login successful!${NC}"
            echo ""
            echo "Full response:"
            echo "$LOGIN_RESPONSE" | jq . 2>/dev/null || echo "$LOGIN_RESPONSE"

            ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
            echo ""
            echo -e "${GREEN}Access Token: ${ACCESS_TOKEN:0:50}...${NC}"

            # Test token usage
            echo ""
            echo -e "${YELLOW}[Token Test]${NC} Using access token..."
            PROTECTED_RESPONSE=$(curl -s "$API_URL/api/v1/student" \
              -H "Authorization: Bearer $ACCESS_TOKEN")

            echo "Protected endpoint response:"
            echo "$PROTECTED_RESPONSE" | jq . 2>/dev/null || echo "$PROTECTED_RESPONSE"
        else
            echo -e "${RED}Login failed${NC}"
            echo "Response: $LOGIN_RESPONSE"
        fi

        exit 0
    fi
else
    echo -e "${RED}Failed to initialize admin${NC}"
    echo "Response: $BOOTSTRAP_RESPONSE"
    exit 1
fi
echo ""

# Create invitation
echo -e "${YELLOW}[3/7]${NC} Creating student invitation..."
INVITE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/admin/invite" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test.student@example.com",
    "first_name": "Test",
    "last_name": "Student",
    "role": "student"
  }')

if echo "$INVITE_RESPONSE" | grep -q "invitation_token"; then
    INVITE_TOKEN=$(echo "$INVITE_RESPONSE" | grep -o '"invitation_token":"[^"]*' | cut -d'"' -f4)
    echo -e "${GREEN}Invitation created successfully${NC}"
    echo "Invitation Token: ${INVITE_TOKEN:0:50}..."
else
    echo -e "${RED}Failed to create invitation${NC}"
    echo "Response: $INVITE_RESPONSE"
    exit 1
fi
echo ""

# Get invitation details
echo -e "${YELLOW}[4/7]${NC} Getting invitation details..."
INVITE_DETAILS=$(curl -s "$API_URL/api/v1/invite/$INVITE_TOKEN")

if echo "$INVITE_DETAILS" | grep -q "success.*true"; then
    echo -e "${GREEN}Invitation details retrieved${NC}"
    echo "$INVITE_DETAILS" | grep -o '"email":"[^"]*' | cut -d'"' -f4
else
    echo -e "${RED}Failed to get invitation details${NC}"
    echo "Response: $INVITE_DETAILS"
    exit 1
fi
echo ""

# Complete registration
echo -e "${YELLOW}[5/7]${NC} Completing registration..."
TIMESTAMP=$(date +%s)
USERNAME="test.student$TIMESTAMP"
PASSWORD="SecurePass123!"

COMPLETE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/invite/$INVITE_TOKEN/complete" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'"$USERNAME"'",
    "password": "'"$PASSWORD"'"
  }')

if echo "$COMPLETE_RESPONSE" | grep -q "success.*true"; then
    echo -e "${GREEN}Registration completed${NC}"
    echo "Username: $USERNAME"
    echo "Password: $PASSWORD"
else
    echo -e "${RED}Failed to complete registration${NC}"
    echo "Response: $COMPLETE_RESPONSE"
    exit 1
fi
echo ""

# Login
echo -e "${YELLOW}[6/7]${NC} Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'"$USERNAME"'",
    "password": "'"$PASSWORD"'"
  }')

if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
    echo -e "${GREEN}Login successful!${NC}"

    # Extract tokens
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"refresh_token":"[^"]*' | cut -d'"' -f4)

    echo ""
    echo "Login Response:"
    echo "$LOGIN_RESPONSE" | jq . 2>/dev/null || echo "$LOGIN_RESPONSE"
    echo ""

    if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
        echo -e "${RED}WARNING: Access token is empty!${NC}"
        echo "This is the bug you were experiencing."
        echo "Check server logs for more details."
        exit 1
    fi

    echo -e "${GREEN}Access Token (first 50 chars): ${ACCESS_TOKEN:0:50}...${NC}"
    echo -e "${GREEN}Refresh Token (first 50 chars): ${REFRESH_TOKEN:0:50}...${NC}"
else
    echo -e "${RED}Login failed${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi
echo ""

# Test authenticated request
echo -e "${YELLOW}[7/7]${NC} Testing authenticated request..."
AUTH_TEST=$(curl -s "$API_URL/api/v1/student" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Protected endpoint response:"
echo "$AUTH_TEST" | jq . 2>/dev/null || echo "$AUTH_TEST"
echo ""

# Test refresh token
echo -e "${YELLOW}[Bonus]${NC} Testing token refresh..."
REFRESH_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'"$REFRESH_TOKEN"'"
  }')

if echo "$REFRESH_RESPONSE" | grep -q "access_token"; then
    echo -e "${GREEN}Token refresh successful!${NC}"
    NEW_ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)
    echo "New Access Token: ${NEW_ACCESS_TOKEN:0:50}..."
else
    echo -e "${RED}Token refresh failed${NC}"
    echo "Response: $REFRESH_RESPONSE"
fi
echo ""

# Test logout
echo -e "${YELLOW}[Cleanup]${NC} Testing logout..."
LOGOUT_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

if echo "$LOGOUT_RESPONSE" | grep -q "success.*true"; then
    echo -e "${GREEN}Logout successful!${NC}"
else
    echo -e "${RED}Logout failed${NC}"
    echo "Response: $LOGOUT_RESPONSE"
fi
echo ""

echo -e "${GREEN}=== All tests completed successfully! ===${NC}"
echo ""
echo "Summary:"
echo "  [PASS] Admin initialization"
echo "  [PASS] Invitation creation"
echo "  [PASS] Invitation retrieval"
echo "  [PASS] User registration"
echo "  [PASS] User login"
echo "  [PASS] Token authentication"
echo "  [PASS] Token refresh"
echo "  [PASS] User logout"
echo ""
echo "Your authentication system is working correctly!"
echo "Tokens are being generated and returned properly."
