package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
)

type command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	CheckFlags  func() bool // Check the flag values after they're parsed, printing errors and returning false if they're incorrect
	Do          func()      // Run the command (only called if CheckFlags returns true)
}

func (c *command) usage() {
	fmt.Printf("Usage of %s:\n", c.Name)
	c.Flags.PrintDefaults()
}

// flags common to all server transactions
var serverURL, username, password, org, fid, appID, appVer, ofxVersion, clientUID string
var noIndentRequests bool
var carriageReturn bool

func defineServerFlags(f *flag.FlagSet) {
	f.StringVar(&serverURL, "url", "", "Financial institution's OFX Server URL (see ofxhome.com if you don't know it)")
	f.StringVar(&username, "username", "", "Your username at financial institution")
	f.StringVar(&password, "password", "", "Your password at financial institution")
	f.StringVar(&org, "org", "", "'ORG' for your financial institution")
	f.StringVar(&fid, "fid", "", "'FID' for your financial institution")
	f.StringVar(&appID, "appid", "QWIN", "'APPID' to pretend to be")
	f.StringVar(&appVer, "appver", "2400", "'APPVER' to pretend to be")
	f.StringVar(&ofxVersion, "ofxversion", "203", "OFX version to use")
	f.StringVar(&clientUID, "clientuid", "", "Client UID (only required by a few FIs, like Chase)")
	f.BoolVar(&noIndentRequests, "noindent", false, "Don't indent OFX requests")
	f.BoolVar(&carriageReturn, "carriagereturn", false, "Use carriage return as line separator")
}

func checkServerFlags() bool {
	var ret bool = true
	if len(serverURL) == 0 {
		fmt.Println("Error: Server URL empty")
		ret = false
	}
	if len(username) == 0 {
		fmt.Println("Error: Username empty")
		ret = false
	}

	if ret && len(password) == 0 {
		fmt.Printf("Password for %s: ", username)
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Printf("Error reading password: %s\n", err)
			ret = false
		} else {
			password = string(pass)
		}
	}
	return ret
}
