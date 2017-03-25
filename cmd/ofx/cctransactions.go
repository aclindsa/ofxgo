package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"os"
	"time"
)

var ccTransactionsCommand = Command{
	Name:        "transactions-cc",
	Description: "Print credit card transactions and balance",
	Flags:       flag.NewFlagSet("transactions-cc", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          ccTransactions,
}

func init() {
	defineServerFlags(ccTransactionsCommand.Flags)
	ccTransactionsCommand.Flags.StringVar(&acctId, "acctid", "", "AcctId (from `get-accounts` subcommand)")
}

func ccTransactions() {
	client, query := NewRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	statementRequest := ofxgo.CCStatementRequest{
		TrnUID: *uid,
		CCAcctFrom: ofxgo.CCAcct{
			AcctId: ofxgo.String(acctId),
		},
		DtStart: ofxgo.Date(time.Now().AddDate(-1, 0, 0)),
		DtEnd:   ofxgo.Date(time.Now()),
		Include: true,
	}
	query.CreditCards = append(query.CreditCards, &statementRequest)

	response, err := client.Request(query)
	if err != nil {
		fmt.Println("Error requesting account statement:", err)
		os.Exit(1)
	}

	if response.Signon.Status.Code != 0 {
		meaning, _ := response.Signon.Status.CodeMeaning()
		fmt.Printf("Nonzero signon status (%d: %s) with message: %s\n", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
		os.Exit(1)
	}

	if len(response.CreditCards) < 1 {
		fmt.Println("No banking messages received")
		return
	}

	if stmt, ok := response.CreditCards[0].(ofxgo.CCStatementResponse); ok {
		fmt.Printf("Balance: %s %s (as of %s)\n", stmt.BalAmt, stmt.CurDef, stmt.DtAsOf)
		fmt.Println("Transactions:")
		for _, tran := range stmt.BankTranList.Transactions {
			currency := stmt.CurDef
			if len(tran.Currency) > 0 {
				currency = tran.Currency
			} else if len(tran.OrigCurrency) > 0 {
				currency = tran.Currency
			}

			var name string
			if len(tran.Name) > 0 {
				name = string(tran.Name)
			} else {
				name = string(tran.Payee.Name)
			}

			if len(tran.Memo) > 0 {
				name = name + " - " + string(tran.Memo)
			}

			fmt.Printf("%s %-15s %-11s %s\n", tran.DtPosted, tran.TrnAmt.String()+" "+string(currency), tran.TrnType, name)
		}
	}
}
