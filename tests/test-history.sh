#!/bin/bash

file="tests/content/b/README.md"
commits="$(git log -20 --pretty=format:%h -- "${file}")"
history="$(
  curl -s "${url}/_/history/${file}" | \
  htmlq --text '.content table td:first-child'
)"

[[ "${commits}" == "${history}" ]] && ok
fail
