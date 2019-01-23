#!/bin/sh

# git >= 2.9
git config core.hooksPath .gitHooks
# git < 2.9
# find .git/hooks -type l -exec rm {} \; && find .gitHooks -type f -exec ln -sf ../../{} .git/hooks/ \;
