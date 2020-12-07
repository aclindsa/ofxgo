package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"io"
	"os"
)

var profileDownloadCommand = command{
	Name:        "download-profile",
	Description: "Download a FI profile to a file",
	Flags:       flag.NewFlagSet("download-profile", flag.ExitOnError),
	CheckFlags:  downloadProfileCheckFlags,
	Do:          downloadProfile,
}

func init() {
	defineServerFlags(profileDownloadCommand.Flags)
	profileDownloadCommand.Flags.StringVar(&filename, "filename", "./response.ofx", "The file to save to")
}

func downloadProfileCheckFlags() bool {
	// Assume if the user didn't specify username that we should use anonymous
	// values for it and password
	if len(username) == 0 {
		username = "anonymous00000000000000000000000"
		password = "anonymous00000000000000000000000"
	}

	ret := checkServerFlags()

	if len(filename) == 0 {
		fmt.Println("Error: Filename empty")
		return false
	}

	return ret
}

func downloadProfile() {
	client, query := newRequest()

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	profileRequest := ofxgo.ProfileRequest{
		TrnUID: *uid,
	}

	query.Prof = append(query.Prof, &profileRequest)

	response, err := client.RequestNoParse(query)
	if err != nil {
		fmt.Println("Error requesting FI profile:", err)
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
