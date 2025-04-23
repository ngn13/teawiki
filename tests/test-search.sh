#!/bin/bash

term='bel'
results="$(
  curl -s --data-urlencode "term=${term}" "${url}/_/search" | \
  htmlq -a href '.content ul li a'
)"

for file in $(grep -l -r "title: *${term}*" tests/content); do
  contains "${results}" "${file}" || fail
done

for file in $(find tests/content -type f -name "*${term}*"); do
  contains "${results}" "${file}" || fail
done

ok
