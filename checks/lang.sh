#!/bin/bash

# check html tag of the view templates for the "lang" attribute

ret=0

for view in $(find views -maxdepth 1 -type f -name '*.html'); do
  [ ! -z "$(grep '<html.*>' "${view}" | grep 'lang="{{.*}}"')" ] && continue
  echo "${view} is missing lang attribute in HTML tag" && ret=1
done

exit $ret
