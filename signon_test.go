package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"testing"
)

func TestMarshalInvalidSignons(t *testing.T) {
	var client = ofxgo.Client{
		AppID:       "OFXGO",
		AppVer:      "0001",
		SpecVersion: "203",
	}

	var request ofxgo.Request
	request.Signon.UserID = "myusername"
	request.Signon.UserPass = "Pa$$word"
	request.Signon.Org = "BNK"
	request.Signon.Fid = "1987"

	request.SetClientFields(&client)
	_, err := request.Marshal()
	if err != nil {
		t.Fatalf("Unexpected error marshalling signon: %s\n", err)
	}

	request.Signon.UserKey = "mykey"
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to key and password both being specified\n")
	}

	request.Signon.UserPass = ""
	_, err = request.Marshal()
	if err != nil {
		t.Fatalf("Unexpected error marshalling signon: %s\n", err)
	}

	request.Signon.UserID = ""
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to unspecified UserID\n")
	}
	request.Signon.UserID = "lakhgdlsakhgdlkahdglkhsadlkghaslkdghsalkdghalsdhg"
	if err == nil {
		t.Fatalf("Expected error due to UserID too long\n")
	}
	request.Signon.UserID = "myusername"

	request.Signon.UserKey = "adlfahdslkgahdweoihadf98agrha87rghasdf9hawhra2hrkwahhaguhwaoefajkei23hff"
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to UserKey too long\n")
	}
	request.Signon.UserKey = ""

	request.Signon.UserPass = "adlfahdslkgahdweoihadf98agrha87rghasdf9hawhra2hrkwahhaguhwaoefajkei23hff"
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to UserPass too long\n")
	}
	request.Signon.UserPass = "lakhgdlkahd"

	request.Signon.Language = "English"
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to Language too long\n")
	}
	request.Signon.Language = "EN"
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to Language too short\n")
	}
	request.Signon.Language = ""
	_, err = request.Marshal()
	if err != nil || request.Signon.Language != "ENG" {
		t.Fatalf("Empty Language expected to default to ENG: %s\n", err)
	}
	request.Signon.Language = "ENG"

	request.Signon.AppID = ""
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to missing AppID\n")
	}
	request.SetClientFields(&client)
	_, err = request.Marshal()
	if err != nil {
		t.Fatalf("Client expected to set empty AppID: %s\n", err)
	}
	client.AppID = "ALKHGDH"
	request.SetClientFields(&client)
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to AppID too long\n")
	}
	client.AppID = "OFXGO"

	request.Signon.AppVer = ""
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to missing AppVer\n")
	}
	request.SetClientFields(&client)
	_, err = request.Marshal()
	if err != nil {
		t.Fatalf("Client expected to set empty AppVer: %s\n", err)
	}
	client.AppVer = "00002"
	request.SetClientFields(&client)
	_, err = request.Marshal()
	if err == nil {
		t.Fatalf("Expected error due to AppVer too long\n")
	}
	client.AppVer = "0001"

	request.SetClientFields(&client)
	_, err = request.Marshal()
	if err != nil {
		t.Fatalf("Unexpected error after resetting all fields to reasonable values: %s\n", err)
	}
}
