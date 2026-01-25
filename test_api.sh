#!/bin/bash

# ============================================
# Kasir API Test Script
# Test all endpoints of the Kasir API
# ============================================
#
# Usage: ./test_api.sh [BASE_URL]
# Example: ./test_api.sh http://localhost:8080
#
# Prerequisites:
# - curl (for HTTP requests)
# - jq (optional, for pretty JSON output)
#

# Configuration
BASE_URL="${1:-http://localhost:8080}"
CONTENT_TYPE="Content-Type: application/json"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Print functions
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_test() {
    echo -e "${YELLOW}TEST:${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓ Success${NC}\n"
}

print_error() {
    echo -e "${RED}✗ Failed${NC}\n"
}

# Test runner function
run_test() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_code="${5:-2xx}" # Default expect success (2xx)

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    print_test "$test_name"
    echo "Request: $method $endpoint"

    if [ -z "$data" ]; then
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$method" "$BASE_URL$endpoint")
    else
        echo "Data: $data"
        response=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$method" -H "$CONTENT_TYPE" -d "$data" "$BASE_URL$endpoint")
    fi

    http_code=$(echo "$response" | grep "HTTP_CODE:" | cut -d':' -f2)
    body=$(echo "$response" | sed '/HTTP_CODE:/d')

    echo "Response Code: $http_code"
    if command -v jq &> /dev/null; then
        echo "Response Body:"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo "Response Body: $body"
    fi

    # Check if response matches expected code
    local test_passed=false
    case "$expected_code" in
        "2xx")
            if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
                test_passed=true
            fi
            ;;
        "4xx")
            if [ "$http_code" -ge 400 ] && [ "$http_code" -lt 500 ]; then
                test_passed=true
            fi
            ;;
        *)
            if [ "$http_code" -eq "$expected_code" ]; then
                test_passed=true
            fi
            ;;
    esac

    if [ "$test_passed" = true ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        print_success
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        print_error
    fi
}

# Start
clear
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     KASIR API TEST SCRIPT              ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo -e "Base URL: $BASE_URL\n"

# Health Check
print_header "1. HEALTH CHECK"
run_test "Health Check" "GET" "/health"

# Category Tests
print_header "2. CATEGORY TESTS - CREATE"
run_test "Create Category - Electronics" "POST" "/api/categories" \
'{"name":"Electronics","description":"Electronic devices and gadgets"}'

run_test "Create Category - Food" "POST" "/api/categories" \
'{"name":"Food","description":"Food and beverages"}'

run_test "Create Category - Clothing" "POST" "/api/categories" \
'{"name":"Clothing","description":"Clothes and accessories"}'

print_header "3. CATEGORY TESTS - READ"
run_test "Get All Categories" "GET" "/api/categories"
run_test "Get Category by ID (1)" "GET" "/api/categories/1"

print_header "4. CATEGORY TESTS - UPDATE"
run_test "Update Category (1)" "PUT" "/api/categories/1" \
'{"name":"Electronics & Gadgets","description":"All electronic devices and accessories"}'

run_test "Get Updated Category (1)" "GET" "/api/categories/1"

# Product Tests
print_header "5. PRODUCT TESTS - CREATE"
run_test "Create Product - Laptop" "POST" "/api/product" \
'{"name":"Laptop","price":15000000,"stock":10}'

run_test "Create Product - Mouse" "POST" "/api/product" \
'{"name":"Wireless Mouse","price":250000,"stock":50}'

run_test "Create Product - Keyboard" "POST" "/api/product" \
'{"name":"Mechanical Keyboard","price":850000,"stock":30}'

run_test "Create Product - Monitor" "POST" "/api/product" \
'{"name":"Monitor 24 inch","price":2500000,"stock":15}'

print_header "6. PRODUCT TESTS - READ"
run_test "Get All Products" "GET" "/api/product"
run_test "Get Product by ID (1)" "GET" "/api/product/1"

print_header "7. PRODUCT TESTS - UPDATE"
run_test "Update Product (1)" "PUT" "/api/product/1" \
'{"name":"Gaming Laptop","price":18000000,"stock":8}'

run_test "Get Updated Product (1)" "GET" "/api/product/1"

# Error Handling Tests
print_header "8. ERROR HANDLING TESTS"
run_test "Get Non-existent Product (999)" "GET" "/api/product/999" "" "4xx"
run_test "Get Non-existent Category (999)" "GET" "/api/categories/999" "" "4xx"
run_test "Get Product with Invalid ID" "GET" "/api/product/invalid" "" "4xx"
run_test "Get Category with Invalid ID" "GET" "/api/categories/invalid" "" "4xx"
run_test "Create Product with Invalid JSON" "POST" "/api/product" '{"name":"Test","price":"invalid"}' "4xx"
run_test "Create Category with Invalid JSON" "POST" "/api/categories" '{"name":123}' "4xx"
run_test "Update Non-existent Product (999)" "PUT" "/api/product/999" '{"name":"Test","price":1000,"stock":10}' "4xx"
run_test "Update Non-existent Category (999)" "PUT" "/api/categories/999" '{"name":"Test","description":"Test"}' "4xx"

# Delete Tests
print_header "9. DELETE TESTS"
run_test "Delete Product (2)" "DELETE" "/api/product/2"
run_test "Verify Product Deleted (2)" "GET" "/api/product/2" "" "404"
run_test "Delete Category (2)" "DELETE" "/api/categories/2"
run_test "Verify Category Deleted (2)" "GET" "/api/categories/2" "" "404"
run_test "Delete Non-existent Product (999)" "DELETE" "/api/product/999" "" "4xx"
run_test "Delete Non-existent Category (999)" "DELETE" "/api/categories/999" "" "4xx"

# Final Results
print_header "TEST RESULTS"
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"
echo -e "Success Rate: ${BLUE}$((PASSED_TESTS * 100 / TOTAL_TESTS))%${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}╔════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║   ALL TESTS PASSED! ✓                  ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════╝${NC}\n"
    exit 0
else
    echo -e "\n${RED}╔════════════════════════════════════════╗${NC}"
    echo -e "${RED}║   SOME TESTS FAILED! ✗                 ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════╝${NC}\n"
    exit 1
fi
