package main

import (
	"fmt"
	"os"
	"strconv"
)

var commands = []command{
	getAccountsCommand,
	downloadCommand,
	ccDownloadCommand,
	invDownloadCommand,
	bankTransactionsCommand,
	ccTransactionsCommand,
	invTransactionsCommand,
	detectSettingsCommand,
}

func usage() {
	fmt.Println(`The ofxgo command-line client provides a simple interface to
query, parse, and display financial data via the OFX specification.

Usage:
	ofx command [arguments]

The commands are:`)

	maxlen := 0
	for _, cmd := range commands {
		if len(cmd.Name) > maxlen {
			maxlen = len(cmd.Name)
		}
	}
	formatString := "    %-" + strconv.Itoa(maxlen) + "s    %s\n"

	for _, cmd := range commands {
		fmt.Printf(formatString, cmd.Name, cmd.Description)
	}
}

func runCmd(c *command) {
	err := c.Flags.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error parsing flags: %s\n", err)
		c.usage()
		os.Exit(1)
	}

	if !c.CheckFlags() {
		fmt.Println()
		c.usage()
		os.Exit(1)
	}

	c.Do()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Error: Please supply a sub-command. Usage:\n\n")
		usage()
		os.Exit(1)
	}
	cmdName := os.Args[1]
	for _, cmd := range commands {
		if cmd.Name == cmdName {
			runCmd(&cmd)
			os.Exit(0)
		}
	}

	switch cmdName {
	case "-h", "-help", "--help", "help":
		usage()
	default:
		fmt.Println("Error: Invalid sub-command. Usage:")
		usage()
		os.Exit(1)
	}
}
