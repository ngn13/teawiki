#!/bin/bash

dir="tests/content/b"
list="$(curl -s "${url}/${dir}" | htmlq -a href '.list div a')"

for file in $(ls -1 "${dir}"); do
  contains "${list}" "/${dir}/${file}" || fail
done

ok
