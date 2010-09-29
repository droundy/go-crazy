#!/bin/sh

set -ev

grep 'hello(' inline-compiled.go

./inline > noinline.temp

../go-crazy --inline hello inline.go

grep 'hello(' inline-compiled.go && exit 1

./inline > inline.temp

diff inline.temp noinline.temp

echo Inlining works!
