package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"io"
	"os"
)

var downloadCommand = Command{
	Name:        "download-bank",
	Description: "Download a bank account statement to a file",
	Flags:       flag.NewFlagSet("download-bank", flag.ExitOnError),
	CheckFlags:  downloadCheckFlags,
	Do:          download,
}

var filename, bankId, acctId, acctType string

func init() {
	defineServerFlags(downloadCommand.Flags)
	downloadCommand.Flags.StringVar(&filename, "filename", "./download.ofx", "The file to save to")
	downloadCommand.Flags.StringVar(&bankId, "bankid", "", "BankId (from `get-accounts` subcommand)")
	downloadCommand.Flags.StringVar(&acctId, "acctid", "", "AcctId (from `get-accounts` subcommand)")
	downloadCommand.Flags.StringVar(&acctType, "accttype", "CHECKING", "AcctType (from `get-accounts` subcommand)")
}

func downloadCheckFlags() bool {
	ret := checkServerFlags()

	if len(filename) == 0 {
		fmt.Println("Error: Filename empty")
		return false
	}

	return ret
}

func download() {
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
		Include: true,
	}
	query.Bank = append(query.Bank, &statementRequest)

	response, err := client.RequestNoParse(query)
	if err != nil {
		fmt.Println("Error requesting account statement:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file to write to:", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Error writing response to file:", err)
		os.Exit(1)
	}
}
