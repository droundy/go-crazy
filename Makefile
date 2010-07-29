# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=goop

GOFILES=\
	goop.go\
	dummy.go\

include $(GOROOT)/src/Make.cmd

dummy.go: transform/*.go
	cd transform && make install
	echo package main > dummy.go
