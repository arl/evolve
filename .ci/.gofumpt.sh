#!/bin/sh

if [ -n "$(gofumpt -l .)" ]; then
  echo "Following files are not formatted with gofumpt:"
  gofumpt -l .
  exit 1
fi
