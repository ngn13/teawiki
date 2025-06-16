#!/bin/bash

has_tag() {
  # not the exact regex but whatever
  for tag in $(pcre2grep -M 'tags:$\n(  - .*$\n)*' "${1}" | tail -n+2); do
    [[ "$(echo "${tag}" | sed 's/^  - //g')" == "${2}" ]] && return 0
  done

  return 1
}

single_check() {
  # obtain pages listed under /_/tag/
  local pages=$(curl -s "${TW_URL}/_/tag/${1}" | htmlq -a href 'ul li a')

  # check if the pages that are listed actually have that tag
  for page in ${pages[@]}; do
    has_tag "${page:1}" "${1}" || return 1
  done

  # check if the list is missing any pages
  for page in $(find tests/content -type f -name '*.md'); do
    has_tag "${page}" "${1}" || continue
    contains "${pages}" "/${page}" || return 1
  done

  return 0
}

single_check "test"       || fail
single_check "example"    || fail
single_check "lorem"      || fail
single_check "idontexist" || fail

ok
