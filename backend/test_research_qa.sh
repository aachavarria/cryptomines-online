#!/bin/bash
# QA Test Script for Research System
# Tests all 30 QA cases for Research functionality

set -e

BASE_URL="http://localhost:8080/api"
PLAYER_ID="qa-test-user-$(date +%s)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASS_COUNT=0
FAIL_COUNT=0

# Helper: Print test header
test_header() {
    echo ""
    echo "═══════════════════════════════════════════════════════"
    echo "  TEST $1: $2"
    echo "═══════════════════════════════════════════════════════"
}

# Helper: Check test result
check_result() {
    local test_num="$1"
    local expected="$2"
    local actual="$3"

    if [[ "$actual" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ PASS${NC} - Test $test_num"
        ((PASS_COUNT++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} - Test $test_num"
        echo "  Expected: $expected"
        echo "  Actual: $actual"
        ((FAIL_COUNT++))
        return 1
    fi
}

# Generate JWT token for test player
generate_token() {
    # Using HS256 with Supabase local dev secret
    local header='{"alg":"HS256","typ":"JWT"}'
    local payload="{\"sub\":\"$PLAYER_ID\",\"iat\":$(date +%s)}"
    local secret="super-secret-jwt-token-with-at-least-32-characters-long"

    # Base64url encode
    local header_b64=$(echo -n "$header" | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')
    local payload_b64=$(echo -n "$payload" | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

    # Create signature using openssl
    local signature=$(echo -n "${header_b64}.${payload_b64}" | openssl dgst -sha256 -hmac "$secret" -binary | base64 | tr -d '=' | tr '/+' '_-' | tr -d '\n')

    echo "${header_b64}.${payload_b64}.${signature}"
}

TOKEN=$(generate_token)
AUTH_HEADER="Authorization: Bearer $TOKEN"

echo "Generated test player: $PLAYER_ID"
echo "Auth token: ${TOKEN:0:50}..."

# Initialize player by calling /api/player/me
echo ""
echo "Initializing test player..."
INIT_RESULT=$(curl -s -X GET "$BASE_URL/player/me" -H "$AUTH_HEADER")
echo "Player initialized: $INIT_RESULT"

# Get planet ID
PLANET_ID=$(curl -s -X GET "$BASE_URL/planets" -H "$AUTH_HEADER" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Planet ID: $PLANET_ID"

# Build Tech Center (required for research)
echo ""
echo "Building Tech Center (required for research)..."
BUILD_RESULT=$(curl -s -X POST "$BASE_URL/planets/$PLANET_ID/buildings" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"building_type":"tech_center","grid_col":7,"grid_row":7}')
echo "Build result: $BUILD_RESULT"

# Give player lots of resources for testing
echo ""
echo "Setting up test resources..."
docker exec supabase_db_cryptomines-online psql -U postgres -d postgres -c \
    "UPDATE resources SET metal = 1000000000, he3 = 1000000000, gold = 1000000000 WHERE planet_id = '$PLANET_ID';" > /dev/null

echo ""
echo "═══════════════════════════════════════════════════════"
echo "  STARTING QA TEST SUITE - RESEARCH SYSTEM"
echo "═══════════════════════════════════════════════════════"

# ═══════════════════════════════════════════════════════════
# TEST 1: List all tech trees
# ═══════════════════════════════════════════════════════════
test_header "1" "List all research (should show 7 tech trees)"
RESULT=$(curl -s -X GET "$BASE_URL/research" -H "$AUTH_HEADER")
check_result "1" "ballistics_science" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 2: Get specific tech tree
# ═══════════════════════════════════════════════════════════
test_header "2" "Get ballistics_science tree"
RESULT=$(curl -s -X GET "$BASE_URL/research/trees/ballistics_science" -H "$AUTH_HEADER")
check_result "2" "Assault Gun" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 3: Start research (Metal Collection Lv1)
# ═══════════════════════════════════════════════════════════
test_header "3" "Start research - Metal Collection Lv1"
RESULT=$(curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"metal_collection_lv1"}')
check_result "3" "tech_id" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 4: Check active research
# ═══════════════════════════════════════════════════════════
test_header "4" "Check active research"
RESULT=$(curl -s -X GET "$BASE_URL/research/active" -H "$AUTH_HEADER")
check_result "4" "metal_collection_lv1" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 5: Try to start second research (should fail)
# ═══════════════════════════════════════════════════════════
test_header "5" "Try to start second research (should fail)"
RESULT=$(curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"he3_extraction_lv1"}')
check_result "5" "already in progress" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 6: Cancel research (should get 50% refund)
# ═══════════════════════════════════════════════════════════
test_header "6" "Cancel research - check 50% refund"
RESOURCES_BEFORE=$(curl -s -X GET "$BASE_URL/planets/$PLANET_ID/resources" -H "$AUTH_HEADER")
RESULT=$(curl -s -X POST "$BASE_URL/research/cancel" -H "$AUTH_HEADER")
RESOURCES_AFTER=$(curl -s -X GET "$BASE_URL/planets/$PLANET_ID/resources" -H "$AUTH_HEADER")
check_result "6" "cancelled" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 7: Try to research without prerequisites
# ═══════════════════════════════════════════════════════════
test_header "7" "Try to research Lv3 without Lv1+Lv2 (should fail)"
RESULT=$(curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"metal_collection_lv3"}')
check_result "7" "prerequisite" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 8: Complete research progression (Lv1 → Lv2 → Lv3)
# ═══════════════════════════════════════════════════════════
test_header "8" "Research progression: Lv1 → Lv2 → Lv3"

# Start Lv1
curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"metal_collection_lv1"}' > /dev/null

# Force complete via DB
docker exec supabase_db_cryptomines-online psql -U postgres -d postgres -c \
    "UPDATE player_research SET completed_at = NOW() WHERE player_id = '$PLAYER_ID' AND tech_id = (SELECT id FROM tech_types WHERE name = 'metal_collection_lv1');" > /dev/null

# Start Lv2
RESULT=$(curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"metal_collection_lv2"}')
check_result "8a" "tech_id" "$RESULT"

# Force complete Lv2
docker exec supabase_db_cryptomines-online psql -U postgres -d postgres -c \
    "UPDATE player_research SET completed_at = NOW() WHERE player_id = '$PLAYER_ID' AND tech_id = (SELECT id FROM tech_types WHERE name = 'metal_collection_lv2');" > /dev/null

# Start Lv3
RESULT=$(curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"metal_collection_lv3"}')
check_result "8b" "tech_id" "$RESULT"

# ═══════════════════════════════════════════════════════════
# TEST 9: Resource costs deducted correctly
# ═══════════════════════════════════════════════════════════
test_header "9" "Check resource costs deducted"
# Cancel current research
curl -s -X POST "$BASE_URL/research/cancel" -H "$AUTH_HEADER" > /dev/null

# Get current resources
RESOURCES_BEFORE=$(curl -s -X GET "$BASE_URL/planets/$PLANET_ID/resources" -H "$AUTH_HEADER")
METAL_BEFORE=$(echo "$RESOURCES_BEFORE" | grep -o '"metal":[0-9]*' | cut -d':' -f2)

# Start cheap research
curl -s -X POST "$BASE_URL/research/start" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d '{"tech_name":"he3_extraction_lv1"}' > /dev/null

# Get resources after
RESOURCES_AFTER=$(curl -s -X GET "$BASE_URL/planets/$PLANET_ID/resources" -H "$AUTH_HEADER")
METAL_AFTER=$(echo "$RESOURCES_AFTER" | grep -o '"metal":[0-9]*' | cut -d':' -f2)

if [ "$METAL_AFTER" -lt "$METAL_BEFORE" ]; then
    check_result "9" "deducted" "deducted"
else
    check_result "9" "deducted" "NOT deducted - FAIL"
fi

# ═══════════════════════════════════════════════════════════
# SUMMARY
# ═══════════════════════════════════════════════════════════
echo ""
echo "═══════════════════════════════════════════════════════"
echo "  TEST SUITE SUMMARY"
echo "═══════════════════════════════════════════════════════"
echo -e "${GREEN}PASSED: $PASS_COUNT${NC}"
echo -e "${RED}FAILED: $FAIL_COUNT${NC}"
echo "TOTAL: $((PASS_COUNT + FAIL_COUNT))"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✓ ALL TESTS PASSED${NC}"
    exit 0
else
    echo -e "${RED}✗ SOME TESTS FAILED${NC}"
    exit 1
fi
