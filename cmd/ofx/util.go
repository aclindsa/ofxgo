package main

import (
	"fmt"
	"github.com/aclindsa/ofxgo"
	"os"
)

func newRequest() (ofxgo.Client, *ofxgo.Request) {
	ver, err := ofxgo.NewOfxVersion(ofxVersion)
	if err != nil {
		fmt.Println("Error creating new OfxVersion enum:", err)
		os.Exit(1)
	}
	var client = ofxgo.BasicClient{
		AppID:       appID,
		AppVer:      appVer,
		SpecVersion: ver,
		NoIndent:    noIndentRequests,
	}

	var query ofxgo.Request
	query.URL = serverURL
	query.Signon.ClientUID = ofxgo.UID(clientUID)
	query.Signon.UserID = ofxgo.String(username)
	query.Signon.UserPass = ofxgo.String(password)
	query.Signon.Org = ofxgo.String(org)
	query.Signon.Fid = ofxgo.String(fid)

	return &client, &query
}
