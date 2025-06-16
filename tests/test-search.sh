#!/bin/bash

term='au'
body=$(curl -s --data-urlencode "term=${term}" "${TW_URL}/_/search")

# get the list of title match links
links="$(echo "${body}" | htmlq -a href '.content .titles a')"
IFS=$'\n' links=(${links})
index=0

# obtain and check the title matches
while read -r title; do
  path="${links[$index]:1}"
  index=$((index+1))

  grep "title: ${title}" "${path}" &> /dev/null || fail
done < <(echo "${body}" | htmlq -t '.content .titles a')

# get a list of heading match links
links=$(echo "${body}" | htmlq -a href '.content .headings a')
IFS=$'\n' links=(${links})
index=0

# obtain and check the heading matches
while read -r heading; do
  [ -z "${heading}" ] && continue

  path="$(echo "${links[$index]:1}" | cut -d# -f1)"
  index=$((index+1))

  # extract the title
  title="$(echo "${heading}" | grep -oE '^.*:')"
  title="${title::-1}"

  # extract the heading
  heading="$(echo "${heading}" | sed "s/${title}: //g")"

  grep -E "^title: ${title}$" "${path}" &> /dev/null || fail
  grep -E "^#* ${heading}$" "${path}" &> /dev/null || fail
done < <(echo "${body}" | htmlq -t '.content .headings a' | sed 's/   *//g')

ok
