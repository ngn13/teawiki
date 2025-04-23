#!/bin/bash

export webhook_secret="topsecret"
export url="http://127.0.0.1:8080"

fail() {
  exit 1
}

ok() {
  exit 0
}

contains() {
  if [ $# -ne 2 ]; then
    return 1
  fi

  echo "${1}" | grep "${2}" &> /dev/null
  return $?
}

index() {
  if [ $# -ne 2 ]; then
    return 0
  fi

  return $(echo "${1}" | grep "${2}" -n | cut -d: -f1)
}

export -f contains
export -f index
export -f fail
export -f ok
