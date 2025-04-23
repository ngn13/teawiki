#!/bin/bash

files=()
sorted=()

for file in $(git ls-files tests/content/); do
  last=$(git log -1 --pretty=format:%aI -- "${file}")
  files+=("${last}@${file}")
done

for file in $(printf '%s\n' "${files[@]}" | sort -nr | uniq -w 25); do
  sorted+=("$(echo "${file}" | cut -d@ -f2)")
done

recent="$(curl -s "${url}" | htmlq -a href '.latest div a')"
prev=0

for file in "${sorted[@]}"; do
  index "${recent}" "${file}"
  cur=$?

  if [ $cur -lt $prev ]; then
    fail
  fi

  prev=$cur
done

ok
