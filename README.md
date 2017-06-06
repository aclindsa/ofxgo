# OFXGo

[![Go Report Card](https://goreportcard.com/badge/github.com/aclindsa/ofxgo)](https://goreportcard.com/report/github.com/aclindsa/ofxgo)
[![Build Status](https://travis-ci.org/aclindsa/ofxgo.svg?branch=master)](https://travis-ci.org/aclindsa/ofxgo)
[![Coverage Status](https://coveralls.io/repos/github/aclindsa/ofxgo/badge.svg?branch=master)](https://coveralls.io/github/aclindsa/ofxgo?branch=master)
[![GoDoc](https://godoc.org/github.com/aclindsa/ofxgo?status.svg)](https://godoc.org/github.com/aclindsa/ofxgo)

**OFXGo** is a library for querying OFX servers and/or parsing the responses. It
also provides an example command-line client to demonstrate the use of the
library.

## Goals

The main purpose of this project is to provide a library to make it easier to
query financial information with OFX from the comfort of Golang, without having
to marshal/unmarshal to SGML or XML. The library does *not* intend to abstract
away all of the details of the OFX specification, which would be difficult to do
well. Instead, it exposes the OFX SGML/XML hierarchy as structs which mostly
resemble it. Its primary goal is to enable the creation of other personal
finance software in Go (as it was created to allow me to fetch OFX transactions
for my own project, [MoneyGo](https://github.com/aclindsa/moneygo)).

Because the OFX specification is rather... 'comprehensive,' it can be difficult
for those unfamiliar with it to figure out where to start. To that end, I have
created a sample command-line client which uses the library to do simple tasks
(currently it does little more than list accounts and query for balances and
transactions). My hope is that by studying its code, new users will be able to
figure out how to use the library much faster than staring at the OFX
specification (or this library's API documentation). The command-line client
also serves as an easy way for me to test/debug the library with actual
financial institutions, which frequently have 'quirks' in their implementations.
The command-line client can be found in the [cmd/ofx
directory](https://github.com/aclindsa/ofxgo/tree/master/cmd/ofx) of this
repository.

## Library documentation

Documentation can be found with the `go doc` tool, or at
https://godoc.org/github.com/aclindsa/ofxgo

## Using the command-line client

To install the command-line client and test it out, you may do the following:

$ go get -v github.com/aclindsa/ofxgo/cmd/ofx && go install -v github.com/aclindsa/ofxgo/cmd/ofx

Once installed (at ~/go/bin/ofx by default, if you haven't set $GOPATH), the
command's usage should help you to use it (`./ofx --help` for a listing of the
available subcommands and their purposes, `./ofx subcommand --help` for
individual subcommand usage).
