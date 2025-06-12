#!/bin/bash

file="tests/content/a/aurea.md"
content=$(curl -s "${TW_URL}/${file}" | htmlq --text '.tags a')
tag_start=0

# read tags from the file
while read -r line; do
  if [[ "${line}" == "tags:" ]]; then
    tag_start=1
    continue
  elif [ $tag_start -eq 1 ] && [[ "${line:0:2}" != "- " ]]; then
    tag_start=0
    break
  fi

  [ $tag_start -eq 0 ] && continue
  tags+=("${line:2}")
done < "${file}"

# check if tags exist in the content
for tag in ${tags[@]}; do
  contains "${content}" "${tag}" || fail
done

ok
