#!/bin/bash

pages=$(curl -s "${TW_URL}/sitemap.xml" | grep -oP "(?<=<loc>)[^<]+")
urllen=${#TW_URL}
count=0

for page in ${pages[@]}; do
  [[ "${page}" == "${TW_URL}" ]] && continue
  [ ! -f "./${page:$urllen}" ] && fail
  ((count++))
done

for file in $(find tests/content -type f -name '*.md'); do
  ((count--))
done

[ $count -ne 0 ] && fail
ok
