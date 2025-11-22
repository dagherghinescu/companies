#!/bin/bash
set -euo pipefail

# Config
BASE_URL="http://localhost:8080"

# This is a mock user.
# Production users will be handled in a separate database. 
USERNAME="admin"
PASSWORD="admin123"

# Helper to check HTTP status
check_status() {
  local status=$1
  local expected=$2
  local msg=$3
  if [ "$status" -ne "$expected" ]; then
    echo "❌ $msg — expected HTTP $expected, got $status"
    exit 1
  else
    echo "✅ $msg"
  fi
}

echo "=== Logging in ==="
LOGIN_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\", \"password\":\"$PASSWORD\"}")

HTTP_STATUS=$(echo "$LOGIN_RESP" | tail -n1)
BODY=$(echo "$LOGIN_RESP" | head -n -1)

check_status "$HTTP_STATUS" 200 "Login"

TOKEN=$(echo "$BODY" | jq -r '.token')
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo "❌ Failed to get JWT token"
  exit 1
fi
echo "Token: $TOKEN"

echo "=== Creating company ==="
CREATE_RESP=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/companies" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Acme Corp v22",
    "description": "A sample company",
    "amount_of_employees": 100,
    "registered": true,
    "type": "Corporations"
  }')

HTTP_STATUS=$(echo "$CREATE_RESP" | tail -n1)
BODY=$(echo "$CREATE_RESP" | head -n -1)
check_status "$HTTP_STATUS" 201 "Create company"

COMPANY_ID=$(echo "$BODY" | jq -r '.id')
echo "Company ID: $COMPANY_ID"

echo "=== Getting company ==="
GET_RESP=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/companies/$COMPANY_ID" \
  -H "Content-Type: application/json")

HTTP_STATUS=$(echo "$GET_RESP" | tail -n1)
BODY=$(echo "$GET_RESP" | head -n -1)
check_status "$HTTP_STATUS" 200 "Get company"
echo "$BODY" | jq

echo "=== Patching company ==="
PATCH_RESP=$(curl -s -w "\n%{http_code}" -X PATCH "$BASE_URL/companies/$COMPANY_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "description": "Updated description",
    "amount_of_employees": 120
  }')

HTTP_STATUS=$(echo "$PATCH_RESP" | tail -n1)
BODY=$(echo "$PATCH_RESP" | head -n -1)
check_status "$HTTP_STATUS" 200 "Patch company"
echo "$BODY" | jq

echo "=== Deleting company ==="
DELETE_RESP=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/companies/$COMPANY_ID" \
  -H "Authorization: Bearer $TOKEN")

HTTP_STATUS=$(echo "$DELETE_RESP" | tail -n1)
check_status "$HTTP_STATUS" 204 "Delete company"

echo "✅ All integration tests passed!"
