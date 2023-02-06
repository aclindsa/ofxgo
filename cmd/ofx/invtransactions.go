package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"os"
)

var invTransactionsCommand = command{
	Name:        "transactions-inv",
	Description: "Print investment transactions",
	Flags:       flag.NewFlagSet("transactions-inv", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          invTransactions,
}

func init() {
	defineServerFlags(invTransactionsCommand.Flags)
	invTransactionsCommand.Flags.StringVar(&acctID, "acctid", "", "AcctID (from `get-accounts` subcommand)")
	invTransactionsCommand.Flags.StringVar(&brokerID, "brokerid", "", "BrokerID (from `get-accounts` subcommand)")
}

func invTransactions() {
	client, query := newRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	statementRequest := ofxgo.InvStatementRequest{
		TrnUID: *uid,
		InvAcctFrom: ofxgo.InvAcct{
			BrokerID: ofxgo.String(brokerID),
			AcctID:   ofxgo.String(acctID),
		},
		Include:        true,
		IncludeOO:      true,
		IncludePos:     true,
		IncludeBalance: true,
		Include401K:    true,
		Include401KBal: true,
	}
	query.InvStmt = append(query.InvStmt, &statementRequest)

	if dryrun {
		printRequest(client, query)
		return
	}

	response, err := client.Request(query)
	if err != nil {
		os.Exit(1)
	}

	if response.Signon.Status.Code != 0 {
		meaning, _ := response.Signon.Status.CodeMeaning()
		fmt.Printf("Nonzero signon status (%d: %s) with message: %s\n", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
		os.Exit(1)
	}

	if len(response.InvStmt) < 1 {
		fmt.Println("No investment messages received")
		return
	}

	if stmt, ok := response.InvStmt[0].(*ofxgo.InvStatementResponse); ok {
		availCash := stmt.InvBal.AvailCash
		if availCash.IsInt() && availCash.Num().Int64() != 0 {
			fmt.Printf("Balance: %s %s (as of %s)\n", stmt.InvBal.AvailCash, stmt.CurDef, stmt.DtAsOf)
		}
		for _, banktrans := range stmt.InvTranList.BankTransactions {
			fmt.Printf("\nBank Transactions for %s subaccount:\n", banktrans.SubAcctFund)
			for _, tran := range banktrans.Transactions {
				printTransaction(stmt.CurDef, &tran)
			}
		}
		fmt.Printf("\nInvestment Transactions:\n")
		for _, t := range stmt.InvTranList.InvTransactions {
			fmt.Printf("%-14s", t.TransactionType())
			switch tran := t.(type) {
			case ofxgo.BuyDebt:
				printInvBuy(stmt.CurDef, &tran.InvBuy)
			case ofxgo.BuyMF:
				printInvBuy(stmt.CurDef, &tran.InvBuy)
			case ofxgo.BuyOpt:
				printInvBuy(stmt.CurDef, &tran.InvBuy)
			case ofxgo.BuyOther:
				printInvBuy(stmt.CurDef, &tran.InvBuy)
			case ofxgo.BuyStock:
				printInvBuy(stmt.CurDef, &tran.InvBuy)
			case ofxgo.ClosureOpt:
				printInvTran(&tran.InvTran)
				fmt.Printf("%s %s contracts (%d shares each)\n", tran.OptAction, tran.Units, tran.ShPerCtrct)
			case ofxgo.Income:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s %s (%s %s)\n", tran.IncomeType, tran.Total, currency, tran.SecID.UniqueIDType, tran.SecID.UniqueID)
				// TODO print ticker instead of CUSIP
			case ofxgo.InvExpense:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s (%s %s)\n", tran.Total, currency, tran.SecID.UniqueIDType, tran.SecID.UniqueID)
				// TODO print ticker instead of CUSIP
			case ofxgo.JrnlFund:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s %s (%s -> %s)\n", tran.Total, stmt.CurDef, tran.SubAcctFrom, tran.SubAcctTo)
			case ofxgo.JrnlSec:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s %s %s (%s -> %s)\n", tran.Units, tran.SecID.UniqueIDType, tran.SecID.UniqueID, tran.SubAcctFrom, tran.SubAcctTo)
				// TODO print ticker instead of CUSIP
			case ofxgo.MarginInterest:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s\n", tran.Total, currency)
			case ofxgo.Reinvest:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s (%s %s)@%s %s (Total: %s)\n", tran.Units, tran.SecID.UniqueIDType, tran.SecID.UniqueID, tran.UnitPrice, currency, tran.Total)
				// TODO print ticker instead of CUSIP
			case ofxgo.RetOfCap:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s (%s %s)\n", tran.Total, currency, tran.SecID.UniqueIDType, tran.SecID.UniqueID)
				// TODO print ticker instead of CUSIP
			case ofxgo.SellDebt:
				printInvSell(stmt.CurDef, &tran.InvSell)
			case ofxgo.SellMF:
				printInvSell(stmt.CurDef, &tran.InvSell)
			case ofxgo.SellOpt:
				printInvSell(stmt.CurDef, &tran.InvSell)
			case ofxgo.SellOther:
				printInvSell(stmt.CurDef, &tran.InvSell)
			case ofxgo.SellStock:
				printInvSell(stmt.CurDef, &tran.InvSell)
			case ofxgo.Split:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if ok, _ := tran.Currency.Valid(); ok {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %d/%d %s -> %s shares of %s %s (%s %s for fractional shares)\n", tran.Numerator, tran.Denominator, tran.OldUnits, tran.NewUnits, tran.SecID.UniqueIDType, tran.SecID.UniqueID, tran.FracCash, currency)
				// TODO print ticker instead of CUSIP
			case ofxgo.Transfer:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s (%s %s) %s\n", tran.Units, tran.SecID.UniqueIDType, tran.SecID.UniqueID, tran.TferAction)
				// TODO print ticker instead of CUSIP
			}
		}
	}
}

func printInvTran(it *ofxgo.InvTran) {
	fmt.Printf("%s", it.DtTrade)
}

func printInvBuy(defCurrency ofxgo.CurrSymbol, ib *ofxgo.InvBuy) {
	printInvTran(&ib.InvTran)
	currency := defCurrency
	if ok, _ := ib.Currency.Valid(); ok {
		currency = ib.Currency.CurSym
	}

	fmt.Printf("%s (%s %s)@%s %s (Total: %s)\n", ib.Units, ib.SecID.UniqueIDType, ib.SecID.UniqueID, ib.UnitPrice, currency, ib.Total)
	// TODO print ticker instead of CUSIP
}

func printInvSell(defCurrency ofxgo.CurrSymbol, is *ofxgo.InvSell) {
	printInvTran(&is.InvTran)
	currency := defCurrency
	if ok, _ := is.Currency.Valid(); ok {
		currency = is.Currency.CurSym
	}

	fmt.Printf(" %s (%s %s)@%s %s (Total: %s)\n", is.Units, is.SecID.UniqueIDType, is.SecID.UniqueID, is.UnitPrice, currency.String(), is.Total)
	// TODO print ticker instead of CUSIP
}
