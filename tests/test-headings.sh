#!/bin/bash

file='tests/content/b/bellaque.md'
content=$(curl -s "${TW_URL}/${file}" | htmlq --text '.headings a')

while read -r heading; do
  heading="$(echo "${heading}" | sed 's/#//g')"
  contains "${content}" "${heading:1}" || fail
done < <(grep '#' "${file}")

ok
