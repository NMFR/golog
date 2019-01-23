#!/bin/sh

STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
[ -z "$STAGED_GO_FILES" ] && echo >&2 "No files stagged" && exit 0

UNFORMATTED_GO_FILES=$(gofmt -l $STAGED_GO_FILES)
[ -z "$UNFORMATTED_GO_FILES" ] &&  echo >&2 "No unformated files" && exit 0

echo >&2 "Formating files:"
for fn in $UNFORMATTED_GO_FILES; do
    gofmt -w $PWD/$fn
    git add $PWD/$fn
    echo >&2 $fn
done
