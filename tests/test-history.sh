#!/bin/bash

file="tests/content/b/README.md"
commits="$(git log -20 --pretty=format:%h -- "${file}")"
body=$(curl -s "${url}/_/history/${file}")
history="$(
  echo "${body}" | \
  htmlq --text '.content table td:first-child'
)"

[[ "${commits}" == "${history}" ]] && ok

echo "${commits} does not match with ${history}"
echo "--- body ---"
echo "${body}"
echo "------------"

fail
