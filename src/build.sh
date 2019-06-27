#!/bin/bash
GOPATH=$(pwd)/.. go build ./rogue.go
./rogue &
