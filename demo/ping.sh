#!/bin/bash

# filer gateway connection detail
[ -z $API_URL ] && API_URL="http://localhost:8080/v1"

[ -z $API_AUTH_URL ] && echo "missing API_AUTH_URL env. variable." 1>&2 && exit 1
[ -z $API_CLIENT_SECRET ] && echo "missing API_CLIENT_SECRET env. variable." 1>&2 && exit 1
[ -z $API_CLIENT_ID ] && API_CLIENT_ID="filer-gateway-client"

auth_host=${API_AUTH_URL}
auth_client_id="${API_CLIENT_ID}"
auth_client_secret="${API_CLIENT_SECRET}"
auth_scopes="urn:dccn:filer-gateway:*"

token_cache=/tmp/.pdb_token_$(whoami)

request_token()
{
    local response token expires

    if [[ -f ${token_cache} ]]; then
        if read -r expires token < ${token_cache}; then
            if [[ "$(date '+%s')" < "$expires" ]]; then
                >&2 echo "Using cached access token."
                echo -n "$token"
                return 0
            fi
        fi
        rm -f ${token_cache} 
    fi

    >&2 echo "Requesting new access token."

    response="$(curl -sSf \
        -X POST \
        -H 'Content-Type: application/x-www-form-urlencoded' \
        --data-urlencode 'grant_type=client_credentials' \
        --data-urlencode "client_id=$auth_client_id" \
        --data-urlencode "client_secret=$auth_client_secret" \
        --data-urlencode "scope=$auth_scopes" \
        "$auth_host/connect/token")"

    token="$(echo -n "$response" | jq -r '.access_token')"
    expires="$(echo -n "$response" | jq -r 'now + .expires_in | floor')"

    echo "$expires" "$token" > ${token_cache} 

    echo -n "$token"
}

token="$(request_token)"
curl -sS \
     -X GET "${API_URL}/ping" \
     -H 'Content-Type: application/json' \
     -H "Authorization: Bearer $token" \
     -w '\n%{content_type}\n%{http_code}'
