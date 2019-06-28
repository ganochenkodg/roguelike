#!/bin/bash
GOPATH=$(pwd)/.. go build ./rogue.go
LD_LIBRARY_PATH=$(pwd)/bearlibterminal ./rogue &
