package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"io"
	"os"
)

var ccDownloadCommand = Command{
	Name:        "download-cc",
	Description: "Download a credit card account statement to a file",
	Flags:       flag.NewFlagSet("download-cc", flag.ExitOnError),
	CheckFlags:  ccDownloadCheckFlags,
	Do:          ccDownload,
}

func init() {
	defineServerFlags(ccDownloadCommand.Flags)
	ccDownloadCommand.Flags.StringVar(&filename, "filename", "./download.ofx", "The file to save to")
	ccDownloadCommand.Flags.StringVar(&acctID, "acctid", "", "AcctID (from `get-accounts` subcommand)")
}

func ccDownloadCheckFlags() bool {
	ret := checkServerFlags()

	if len(filename) == 0 {
		fmt.Println("Error: Filename empty")
		return false
	}

	return ret
}

func ccDownload() {
	client, query := NewRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	statementRequest := ofxgo.CCStatementRequest{
		TrnUID: *uid,
		CCAcctFrom: ofxgo.CCAcct{
			AcctID: ofxgo.String(acctID),
		},
		Include: true,
	}
	query.CreditCard = append(query.CreditCard, &statementRequest)

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
