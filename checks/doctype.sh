#!/bin/bash

# check for doctype at the start of the view templates

ret=0

for view in $(find views -maxdepth 1 -type f -name '*.html'); do
  [[ "$(sed '1q;d' "${view}")" == '<!doctype html>' ]] && continue
  echo "${view} is missing doctype" && ret=1
done

exit $ret
