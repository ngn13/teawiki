#!/bin/bash

page="tests/content/b/README.md"
dir="$(dirname "${page}")"

list="$(curl -s "${TW_URL}/${page}" | htmlq -a href '.list div a')"

for file in $(ls -1 "${dir}"); do
  contains "${list}" "/${dir}/${file}" || fail
done

ok
