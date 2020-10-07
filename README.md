# OFXGo

[![Go Report Card](https://goreportcard.com/badge/github.com/aclindsa/ofxgo)](https://goreportcard.com/report/github.com/aclindsa/ofxgo)
[![Build Status](https://travis-ci.com/aclindsa/ofxgo.svg?branch=master)](https://travis-ci.com/aclindsa/ofxgo)
[![Coverage Status](https://coveralls.io/repos/github/aclindsa/ofxgo/badge.svg?branch=master)](https://coveralls.io/github/aclindsa/ofxgo?branch=master)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/aclindsa?ofxgo)](https://pkg.go.dev/github.com/aclindsa/ofxgo)

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
https://pkg.go.dev/github.com/aclindsa/ofxgo

## Example Usage

The following code snippet demonstrates how to use OFXGo to query and parse
OFX code from a checking account, printing the balance and returned transactions:

```go
client := ofxgo.BasicClient{} // Accept the default Client settings

// These values are specific to your bank
var query ofxgo.Request
query.URL = "https://secu.example.com/ofx"
query.Signon.Org = ofxgo.String("SECU")
query.Signon.Fid = ofxgo.String("1234")

// Set your username/password
query.Signon.UserID = ofxgo.String("username")
query.Signon.UserPass = ofxgo.String("hunter2")

uid, _ := ofxgo.RandomUID() // Handle error in real code
query.Bank = append(query.Bank, &ofxgo.StatementRequest{
	TrnUID: *uid,
	BankAcctFrom: ofxgo.BankAcct{
		BankID:   ofxgo.String("123456789"),   // Possibly your routing number
		AcctID:   ofxgo.String("00011122233"), // Possibly your account number
		AcctType: ofxgo.AcctTypeChecking,
	},
	Include: true, // Include transactions (instead of only balance information)
})

response, _ := client.Request(&query) // Handle error in real code

// Was there an OFX error while processing our request?
if response.Signon.Status.Code != 0 {
	meaning, _ := response.Signon.Status.CodeMeaning()
	fmt.Printf("Nonzero signon status (%d: %s) with message: %s\n", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
	os.Exit(1)
}

if len(response.Bank) < 1 {
	fmt.Println("No banking messages received")
	os.Exit(1)
}

if stmt, ok := response.Bank[0].(*ofxgo.StatementResponse); ok {
	fmt.Printf("Balance: %s %s (as of %s)\n", stmt.BalAmt, stmt.CurDef, stmt.DtAsOf)
	fmt.Println("Transactions:")
	for _, tran := range stmt.BankTranList.Transactions {
		currency := stmt.CurDef
		if ok, _ := tran.Currency.Valid(); ok {
			currency = tran.Currency.CurSym
		}
		fmt.Printf("%s %-15s %-11s %s%s%s\n", tran.DtPosted, tran.TrnAmt.String()+" "+currency.String(), tran.TrnType, tran.Name, tran.Payee.Name, tran.Memo)
	}
}
```

## Requirements

OFXGo requires go >= 1.9

## Using the command-line client

To install the command-line client and test it out, you may do the following:

$ go get -v github.com/aclindsa/ofxgo/cmd/ofx && go install -v github.com/aclindsa/ofxgo/cmd/ofx

Once installed (at ~/go/bin/ofx by default, if you haven't set $GOPATH), the
command's usage should help you to use it (`./ofx --help` for a listing of the
available subcommands and their purposes, `./ofx subcommand --help` for
individual subcommand usage).
