#!/bin/bash

# check the general syntax of the view templates

if grep -Er '{{  *|  *}}' views; then
  echo "found spaced brackets in the view templates"
  exit 1
fi

if grep -Pr '\t' views; then
  echo "found tab character in the view templates"
  exit 1
fi

exit 0
