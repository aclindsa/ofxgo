package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"os"
	"time"
)

var bankTransactionsCommand = Command{
	Name:        "transactions-bank",
	Description: "Print bank transactions and balance",
	Flags:       flag.NewFlagSet("transactions-bank", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          bankTransactions,
}

func init() {
	defineServerFlags(bankTransactionsCommand.Flags)
	bankTransactionsCommand.Flags.StringVar(&bankId, "bankid", "", "BankId (from `get-accounts` subcommand)")
	bankTransactionsCommand.Flags.StringVar(&acctId, "acctid", "", "AcctId (from `get-accounts` subcommand)")
	bankTransactionsCommand.Flags.StringVar(&acctType, "accttype", "CHECKING", "AcctType (from `get-accounts` subcommand)")
}

func bankTransactions() {
	client, query := NewRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	statementRequest := ofxgo.StatementRequest{
		TrnUID: *uid,
		BankAcctFrom: ofxgo.BankAcct{
			BankId:   ofxgo.String(bankId),
			AcctId:   ofxgo.String(acctId),
			AcctType: ofxgo.String(acctType),
		},
		DtStart: ofxgo.Date(time.Now().AddDate(-1, 0, 0)),
		DtEnd:   ofxgo.Date(time.Now()),
		Include: true,
	}
	query.Banking = append(query.Banking, &statementRequest)

	response, err := client.Request(query)
	if err != nil {
		fmt.Println("Error requesting account statement:", err)
		os.Exit(1)
	}

	if response.Signon.Status.Code != 0 {
		fmt.Printf("Nonzero signon status with message: %s\n", response.Signon.Status.Message)
		os.Exit(1)
	}

	if len(response.Banking) < 1 {
		fmt.Println("No banking messages received")
		return
	}

	if stmt, ok := response.Banking[0].(ofxgo.StatementResponse); ok {
		fmt.Printf("Balance: %s %s (as of %s)\n", stmt.BalAmt, stmt.CurDef, stmt.DtAsOf)
		fmt.Println("Transactions:")
		for _, tran := range stmt.BankTranList.Transactions {
			printTransaction(stmt.CurDef, &tran)
		}
	}
}

func printTransaction(defCurrency ofxgo.String, tran *ofxgo.Transaction) {
	currency := defCurrency
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
