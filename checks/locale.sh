#!/bin/bash

# OPTIONAL

names() {
  yq -r '[paths | join(".")] | .[]' "${1}"
}

en=$(names "locale/en.yaml")
code=0

for locale in $(find locale -type f -name '*.yaml' -not -name en.yaml); do
  keys=$(names "${locale}")

  for key in ${en[@]}; do
    echo "${keys}" | grep "${key}" &> /dev/null && continue
    echo "'${key}' is missing in locale ${locale}"
    code=1
  done
done

exit $code
