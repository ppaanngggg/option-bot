#!/usr/bin/env bash

# 1. func to buf lint and buf generate
gen_buf() {
  cd ./proto || exit
  buf lint
  buf generate
  cd ..
}

# 2. func to generate ent
gen_ent() {
  go run entgo.io/ent/cmd/ent generate ./ent/schema
}

# if no args, then run all
if [ $# -eq 0 ]; then
  gen_buf
  gen_ent
  exit 0
fi
# if args, iter args and run
for arg in "$@"; do
  case $arg in
  buf)
    gen_buf
    ;;
  ent)
    gen_ent
    ;;
  *)
    echo "unknown arg: $arg"
    ;;
  esac
done
