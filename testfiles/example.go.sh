#!/bin/sh

set -ev

grep _dot_sub example-compiled.go

./example

./example | grep 'Hello world!'

cat > Makefile <<EOF
include \$(GOROOT)/src/Make.inc

TARG=foo

GOFILES=\
	example-compiled.go\

include \$(GOROOT)/src/Make.cmd
EOF

make

./foo

rm Makefile foo
