#!/usr/bin/env bash
fmt:
	gofmt -l -w ./

test:
	go test

bench:
	go test -test.benchmem -test.bench=".*" -count=3
