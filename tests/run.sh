#!/bin/bash

# this script runs all the other test scripts after starting
# the wiki in the background, and checks all of their results

# do not run this script directly, instead run "make test"

source tests/common.sh

depends=(
  curl
  git
  htmlq
)

for depend in "${depends[@]}"; do
  command -v "${depend}" &> /dev/null && continue
  echo "${depend} not found"
  exit
done

run() {
  ./"${1}"

  if [ $? -eq 0 ]; then
    echo "$(basename "${script}"): OK"
  else
    echo "$(basename "${script}"): FAIL"
    ret=1
  fi
}

script=""
ret=0

if [ $# -eq 1 ]; then
  script="tests/test-${1}.sh"
  if [ ! -f "${script}" ]; then
    echo "test script not found"
    exit 1
  fi
fi

export TW_URL="http://127.0.0.1:8080"
export TW_REPO_PATH="."
export TW_WEBHOOK_SOURCE="gitea"
export TW_WEBHOOK_SECRET="topsecret"
./teawiki.elf &
sleep 3

if [[ "200" != "$(curl -s "${TW_URL}" -o /dev/null -w '%{http_code}')" ]]; then
  echo "server did not respond with 200"
  ret=1
elif [ ! -z "${script}" ]; then
  run "${script}"
else
  for script in tests/test-*.sh; do
    run "${script}"
  done
fi

pkill -9 teawiki.elf
exit $ret
