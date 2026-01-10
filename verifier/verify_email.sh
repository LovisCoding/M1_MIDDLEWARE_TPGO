#!/bin/bash

# Base URL for Config API
API_URL="http://localhost:8080"
EMAIL=${1:-"arthur.lecomte@etu.uca.fr"}


echo "Fetching Mail ID..."
MAIL_ID=$(curl -s "$API_URL/mails" | jq -r ".[] | select(.mail==\"$EMAIL\") | .id" | tail -n 1)
echo "Mail ID: $MAIL_ID"

echo "2. Creating Test Resource..."
# Create Resource
curl -X POST "$API_URL/resources" \
     -H "Content-Type: application/json" \
     -d '{"name": "Test Resource", "type": "Test"}'
echo ""

echo "Fetching Resource ID..."
RESOURCE_ID=$(curl -s "$API_URL/resources" | jq -r '.[] | select(.name=="Test Resource") | .id' | tail -n 1)
echo "Resource ID: $RESOURCE_ID"

if [ -z "$MAIL_ID" ] || [ -z "$RESOURCE_ID" ]; then
    echo "Failed to get Mail ID or Resource ID. Exiting."
    exit 1
fi

echo "3. Creating Alert (Subscription)..."
# Create Alert: POST /resources/{id}/alerts/{mail_id}
curl -X POST "$API_URL/resources/$RESOURCE_ID/alerts/$MAIL_ID"
echo ""

echo "4. Running Verifier Tool..."
# Navigate to script directory to find main.go
cd "$(dirname "$0")" && go run main.go
