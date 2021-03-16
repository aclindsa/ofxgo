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
	var client = ofxgo.GetClient(serverURL,
		&ofxgo.BasicClient{
			AppID:          appID,
			AppVer:         appVer,
			SpecVersion:    ver,
			NoIndent:       noIndentRequests,
			CarriageReturn: carriageReturn,
			UserAgent:      userAgent,
		})

	var query ofxgo.Request
	query.URL = serverURL
	query.Signon.ClientUID = ofxgo.UID(clientUID)
	query.Signon.UserID = ofxgo.String(username)
	query.Signon.UserPass = ofxgo.String(password)
	query.Signon.Org = ofxgo.String(org)
	query.Signon.Fid = ofxgo.String(fid)

	return client, &query
}

func printRequest(c ofxgo.Client, r *ofxgo.Request) {
	r.SetClientFields(c)

	b, err := r.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b)
}
