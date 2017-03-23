package main

import (
	"flag"
	"fmt"
	"github.com/aclindsa/ofxgo"
	"os"
	"time"
)

var detectSettingsCommand = Command{
	Name:        "detect-settings",
	Description: "Attempt to guess client settings needed for a particular financial institution",
	Flags:       flag.NewFlagSet("detect-settings", flag.ExitOnError),
	CheckFlags:  checkServerFlags,
	Do:          detectSettings,
}

var delay uint64

func init() {
	detectSettingsCommand.Flags.StringVar(&serverURL, "url", "", "Financial institution's OFX Server URL (see ofxhome.com if you don't know it)")
	detectSettingsCommand.Flags.StringVar(&username, "username", "", "Your username at financial institution")
	detectSettingsCommand.Flags.StringVar(&password, "password", "", "Your password at financial institution")
	detectSettingsCommand.Flags.StringVar(&org, "org", "", "'ORG' for your financial institution")
	detectSettingsCommand.Flags.StringVar(&fid, "fid", "", "'FID' for your financial institution")
	detectSettingsCommand.Flags.Uint64Var(&delay, "delay", 500, "How long to delay between two subsequent requests, in milliseconds")
}

// We keep a separate list of APPIDs to preserve the ordering (ordering isn't
// guaranteed in maps). We want to try them in order from 'best' and most
// likely to work to 'worse' and least likely to work
var appIds = []string{
	"OFXGO", // ofxgo (this library)
	"QWIN",  // Intuit Quicken Windows
	"QMOFX", // Intuit Quicken Mac
	"QB",    // Intuit QuickBooks Windows
	"Money", // Microsoft Money 2007
}

var appVersions = map[string][]string{
	"OFXGO": []string{ // ofxgo (this library)
		"0001",
	},
	"QWIN": []string{ // Intuit Quicken Windows
		"2400", // 2015
		"2300", // 2014
		"2200", // 2013
		"2100", // 2012
		"2000", // 2011
		"1900", // 2010
		"1800", // 2009
		"1700", // 2008
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
	},
	"QMOFX": []string{ // Intuit Quicken Mac
		"1700", // 2008
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
	},
	"QB": []string{ // Intuit QuickBooks Windows
		"1800", // 2008
		"1700", // 2007
		"1600", // 2006
		"1500", // 2005
	},
	"Money": []string{ // Microsoft Money 2007
		"1600", // 2007
		"1500", // 2006
		"1400", // 2005
		"1200", // 2004
		"1100", // 2003
	},
}

var versions = []string{
	"203",
	"103",
	"200",
	"201",
	"202",
	"210",
	"211",
	"102",
	"151",
	"160",
	"220",
}

func detectSettings() {
	var attempts uint
	for _, appId := range appIds {
		for _, appVer := range appVersions[appId] {
			for _, version := range versions {
				for _, noIndent := range []bool{false, true} {
					if tryProfile(appId, appVer, version, noIndent) {
						fmt.Println("The following settings were found to work:")
						fmt.Printf("AppId: %s\n", appId)
						fmt.Printf("AppVer: %s\n", appVer)
						fmt.Printf("OFX Version: %s\n", version)
						fmt.Printf("noindent: %t\n", noIndent)
						os.Exit(0)
					} else {
						attempts += 1
						var noIndentString string
						if noIndent {
							noIndentString = " noindent"
						}
						fmt.Printf("Attempt %d failed (%s %s %s%s), trying again after %dms...\n", attempts, appId, appVer, version, noIndentString, delay)
						time.Sleep(time.Duration(delay) * time.Millisecond)
					}
				}
			}
		}
	}
}

const anonymous = "anonymous00000000000000000000000"

func tryProfile(appId, appVer, version string, noindent bool) bool {
	var client = ofxgo.Client{
		AppId:       appId,
		AppVer:      appVer,
		SpecVersion: version,
		NoIndent:    noindent,
	}

	var query ofxgo.Request
	query.URL = serverURL
	query.Signon.ClientUID = ofxgo.UID(clientUID)
	query.Signon.UserId = ofxgo.String(username)
	query.Signon.UserPass = ofxgo.String(password)
	query.Signon.Org = ofxgo.String(org)
	query.Signon.Fid = ofxgo.String(fid)

	uid, err := ofxgo.RandomUID()
	if err != nil {
		fmt.Println("Error creating uid for transaction:", err)
		os.Exit(1)
	}

	profileRequest := ofxgo.ProfileRequest{
		TrnUID:   *uid,
		DtProfUp: ofxgo.Date(time.Unix(0, 0)),
	}
	query.Profile = append(query.Profile, &profileRequest)

	_, err = client.Request(&query)
	if err == nil {
		return true
	}

	// try again with anonymous logins
	query.Signon.UserId = ofxgo.String(anonymous)
	query.Signon.UserPass = ofxgo.String(anonymous)

	_, err = client.Request(&query)
	return err == nil
}
