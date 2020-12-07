package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"io"
	"os"
)

var invDownloadCommand = command{
	Name:        "download-inv",
	Description: "Download a investment account statement to a file",
	Flags:       flag.NewFlagSet("download-inv", flag.ExitOnError),
	CheckFlags:  invDownloadCheckFlags,
	Do:          invDownload,
}

var brokerID string

func init() {
	defineServerFlags(invDownloadCommand.Flags)
	invDownloadCommand.Flags.StringVar(&filename, "filename", "./response.ofx", "The file to save to")
	invDownloadCommand.Flags.StringVar(&acctID, "acctid", "", "AcctID (from `get-accounts` subcommand)")
	invDownloadCommand.Flags.StringVar(&brokerID, "brokerid", "", "BrokerID (from `get-accounts` subcommand)")
}

func invDownloadCheckFlags() bool {
	ret := checkServerFlags()

	if len(filename) == 0 {
		fmt.Println("Error: Filename empty")
		return false
	}

	return ret
}

func invDownload() {
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
