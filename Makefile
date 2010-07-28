# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=goop

GOFILES=\
	goop.go\
	modifywalk.go\

include $(GOROOT)/src/Make.cmd
