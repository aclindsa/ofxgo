package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"io"
	"os"
	"time"
)

var invDownloadCommand = Command{
	Name:        "download-inv",
	Description: "Download a investment account statement to a file",
	Flags:       flag.NewFlagSet("download-inv", flag.ExitOnError),
	CheckFlags:  invDownloadCheckFlags,
	Do:          invDownload,
}

var brokerId string

func init() {
	defineServerFlags(invDownloadCommand.Flags)
	invDownloadCommand.Flags.StringVar(&filename, "filename", "./download.ofx", "The file to save to")
	invDownloadCommand.Flags.StringVar(&acctId, "acctid", "", "AcctId (from `get-accounts` subcommand)")
	invDownloadCommand.Flags.StringVar(&brokerId, "brokerid", "", "BrokerId (from `get-accounts` subcommand)")
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
