#!/bin/bash

body='{}'
sign="$(
  echo -n "${body}" | \
  openssl sha256 -hex -mac HMAC -macopt "key:${TW_WEBHOOK_SECRET}" | \
  awk '{print $2}'
)"

code=$(curl -s "${TW_URL}/_/webhook/gitea" \
  -X POST --data "${body}" \
  -H 'Content-Type: application/json' \
  -H 'X-Gitea-Event: push' \
  -H "HTTP_X_GITEA_SIGNATURE: ${sign}" \
  -o /dev/null \
  -w "%{http_code}")

[[ "${code}" == "202" ]] && ok
fail
