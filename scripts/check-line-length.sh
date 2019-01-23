#!/bin/sh

echo "Checking *.go files for max line length (120)"
! ( grep -nr --include="*.go" '.\{120\}' . ) && echo "OK!"
