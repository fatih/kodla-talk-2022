#!/usr/bin/env bash

set -eou pipefail

repos=("structtag" "vim-go" "gomodifytags")

for repo in "${repos[@]}" ; do
  url="https://github.com/fatih/$repo"
  echo "cloning $url"
  git clone $url
done
