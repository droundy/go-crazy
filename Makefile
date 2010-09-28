# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=go-crazy

GOFILES=\
	go-crazy.go\
	dummy.go\

include $(GOROOT)/src/Make.cmd

dummy.go: parser/*.go scanner/*.go transform/*.go
	cd scanner && make install
	cd parser && make install
	cd transform && make install
	echo package main > dummy.go

cleanall:
	make clean
	rm -f dummy.go
	cd scanner && make clean
	cd parser && make clean
	cd transform && make clean
