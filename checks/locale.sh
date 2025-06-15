#!/bin/bash

names() {
  yq -r '[paths | join(".")] | .[]' "${1}"
}

en=$(names "locale/en.yaml")

for locale in $(find locale -type f -name '*.yaml' -not -name en.yaml); do
  keys=$(names "${locale}")

  for key in ${en[@]}; do
    echo "${keys}" | grep "${key}" &> /dev/null && continue
    echo "'${key}' is missing in locale ${locale}"
    exit 1
  done
done

exit 0
