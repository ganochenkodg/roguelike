#!/bin/bash
GOPATH=$(pwd)/.. go build ./rogue.go
/home/dganochenko/go/rl/src/rogue &
