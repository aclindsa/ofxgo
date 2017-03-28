package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"math/big"
	"os"
	"time"
)

var invTransactionsCommand = Command{
	Name:        "transactions-inv",
	Description: "Print investment transactions",
	Flags:       flag.NewFlagSet("transactions-inv", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          invTransactions,
}

func init() {
	defineServerFlags(invTransactionsCommand.Flags)
	invTransactionsCommand.Flags.StringVar(&acctId, "acctid", "", "AcctId (from `get-accounts` subcommand)")
	invTransactionsCommand.Flags.StringVar(&brokerId, "brokerid", "", "BrokerId (from `get-accounts` subcommand)")
}

func invTransactions() {
	client, query := NewRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	statementRequest := ofxgo.InvStatementRequest{
		TrnUID: *uid,
		InvAcctFrom: ofxgo.InvAcct{
			BrokerId: ofxgo.String(brokerId),
			AcctId:   ofxgo.String(acctId),
		},
		DtStart:        ofxgo.Date(time.Now().AddDate(-1, 0, 0)), // a year ago
		DtEnd:          ofxgo.Date(time.Now().AddDate(0, 0, -1)), // Some FIs (*cough* Fidelity) return errors if DTEND is the current day
		Include:        true,
		IncludeOO:      true,
		PosDtAsOf:      ofxgo.Date(time.Now()),
		IncludePos:     true,
		IncludeBalance: true,
		Include401K:    true,
		Include401KBal: true,
	}
	query.Investments = append(query.Investments, &statementRequest)

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

	if len(response.Investments) < 1 {
		fmt.Println("No investment messages received")
		return
	}

	if stmt, ok := response.Investments[0].(ofxgo.InvStatementResponse); ok {
		availCash := big.Rat(stmt.InvBal.AvailCash)
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
				fmt.Println("%s %s contracts (%s shares each)\n", tran.OptAction, tran.Units, tran, tran.ShPerCtrct)
			case ofxgo.Income:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s %s (%s %s)\n", tran.IncomeType, tran.Total, currency, tran.SecId.UniqueIdType, tran.SecId.UniqueId)
				// TODO print ticker instead of CUSIP
			case ofxgo.InvExpense:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s (%s %s)\n", tran.Total, currency, tran.SecId.UniqueIdType, tran.SecId.UniqueId)
				// TODO print ticker instead of CUSIP
			case ofxgo.JrnlFund:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s %s (%s -> %s)\n", tran.Total, stmt.CurDef, tran.SubAcctFrom, tran.SubAcctTo)
			case ofxgo.JrnlSec:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s %s %s (%s -> %s)\n", tran.Units, tran.SecId.UniqueIdType, tran.SecId.UniqueId, tran.SubAcctFrom, tran.SubAcctTo)
				// TODO print ticker instead of CUSIP
			case ofxgo.MarginInterest:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s\n", tran.Total, currency)
			case ofxgo.Reinvest:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s (%s %s)@%s %s (Total: %s)\n", tran.Units, tran.SecId.UniqueIdType, tran.SecId.UniqueId, tran.UnitPrice, currency, tran.Total)
				// TODO print ticker instead of CUSIP
			case ofxgo.RetOfCap:
				printInvTran(&tran.InvTran)
				currency := stmt.CurDef
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s %s (%s %s)\n", tran.Total, currency, tran.SecId.UniqueIdType, tran.SecId.UniqueId)
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
				if len(tran.Currency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				} else if len(tran.OrigCurrency.CurSym) > 0 {
					currency = tran.Currency.CurSym
				}
				fmt.Printf(" %s/%s %s -> %s shares of %s %s (%s %s for fractional shares)\n", tran.Numerator, tran.Denominator, tran.OldUnits, tran.NewUnits, tran.SecId.UniqueIdType, tran.SecId.UniqueId, tran.FracCash, currency)
				// TODO print ticker instead of CUSIP
			case ofxgo.Transfer:
				printInvTran(&tran.InvTran)
				fmt.Printf(" %s (%s %s) %s\n", tran.Units, tran.SecId.UniqueIdType, tran.SecId.UniqueId, tran.TferAction)
				// TODO print ticker instead of CUSIP
			}
		}
	}
}

func printInvTran(it *ofxgo.InvTran) {
	fmt.Printf("%s", it.DtTrade)
}

func printInvBuy(defCurrency ofxgo.String, ib *ofxgo.InvBuy) {
	printInvTran(&ib.InvTran)
	currency := defCurrency
	if len(ib.Currency.CurSym) > 0 {
		currency = ib.Currency.CurSym
	} else if len(ib.OrigCurrency.CurSym) > 0 {
		currency = ib.Currency.CurSym
	}

	fmt.Printf("%s (%s %s)@%s %s (Total: %s)\n", ib.Units, ib.SecId.UniqueIdType, ib.SecId.UniqueId, ib.UnitPrice, currency, ib.Total)
	// TODO print ticker instead of CUSIP
}

func printInvSell(defCurrency ofxgo.String, is *ofxgo.InvSell) {
	printInvTran(&is.InvTran)
	currency := defCurrency
	if len(is.Currency.CurSym) > 0 {
		currency = is.Currency.CurSym
	} else if len(is.OrigCurrency.CurSym) > 0 {
		currency = is.Currency.CurSym
	}

	fmt.Printf(" %s (%s %s)@%s %s (Total: %s)\n", is.Units, is.SecId.UniqueIdType, is.SecId.UniqueId, is.UnitPrice, currency, is.Total)
	// TODO print ticker instead of CUSIP
}
