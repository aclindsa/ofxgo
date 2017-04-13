package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"strings"
	"testing"
	"time"
)

func TestMarshalProfileRequest(t *testing.T) {
	var expectedString string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20160614073400.000[-5:EST]</DTCLIENT>
			<USERID>anonymous00000000000000000000000</USERID>
			<USERPASS>anonymous00000000000000000000000</USERPASS>
			<LANGUAGE>ENG</LANGUAGE>
			<FI>
				<ORG>BNK</ORG>
				<FID>1987</FID>
			</FI>
			<APPID>OFXGO</APPID>
			<APPVER>0001</APPVER>
		</SONRQ>
	</SIGNONMSGSRQV1>
	<PROFMSGSRQV1>
		<PROFTRNRQ>
			<TRNUID>983373</TRNUID>
			<PROFRQ>
				<CLIENTROUTING>NONE</CLIENTROUTING>
				<DTPROFUP>20160101000000.000[-5:EST]</DTPROFUP>
			</PROFRQ>
		</PROFTRNRQ>
	</PROFMSGSRQV1>
</OFX>`

	var client = ofxgo.Client{
		AppID:       "OFXGO",
		AppVer:      "0001",
		SpecVersion: "203",
	}

	var request ofxgo.Request
	request.Signon.UserID = "anonymous00000000000000000000000"
	request.Signon.UserPass = "anonymous00000000000000000000000"
	request.Signon.Org = "BNK"
	request.Signon.Fid = "1987"

	EST := time.FixedZone("EST", -5*60*60)

	profileRequest := ofxgo.ProfileRequest{
		TrnUID:   "983373",
		DtProfUp: *ofxgo.NewDate(2016, 1, 1, 0, 0, 0, 0, EST),
	}
	request.Prof = append(request.Prof, &profileRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	request.Signon.DtClient = *ofxgo.NewDate(2016, 6, 14, 7, 34, 0, 0, EST)

	marshalCheckRequest(t, &request, expectedString)
}

func TestUnmarshalProfileResponse102(t *testing.T) {
	responseReader := strings.NewReader(`OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
<SIGNONMSGSRSV1>
<SONRS>
<STATUS>
<CODE>0
<SEVERITY>INFO
</STATUS>
<DTSERVER>20170403093458.000
<LANGUAGE>ENG
<DTPROFUP>20021119140000
</SONRS>
</SIGNONMSGSRSV1>
<PROFMSGSRSV1>
<PROFTRNRS>
<TRNUID>0f94ce83-13b7-7568-e4fc-c02c7b47e7ab
<STATUS>
<CODE>0
<SEVERITY>INFO
</STATUS>
<PROFRS>
<MSGSETLIST>
<SIGNONMSGSET>
<SIGNONMSGSETV1>
<MSGSETCORE>
<VER>1
<URL>https://ofx.example.com/cgi-ofx/exampleofx
<OFXSEC>NONE
<TRANSPSEC>Y
<SIGNONREALM>Example Trade
<LANGUAGE>ENG
<SYNCMODE>LITE
<RESPFILEER>N
<INTU.TIMEOUT>300
</MSGSETCORE>
</SIGNONMSGSETV1>
</SIGNONMSGSET>
<SIGNUPMSGSET>
<SIGNUPMSGSETV1>
<MSGSETCORE>
<VER>1
<URL>https://ofx.example.com/cgi-ofx/exampleofx
<OFXSEC>NONE
<TRANSPSEC>Y
<SIGNONREALM>Example Trade
<LANGUAGE>ENG
<SYNCMODE>LITE
<RESPFILEER>N
<INTU.TIMEOUT>300
</MSGSETCORE>
<CLIENTENROLL>
<ACCTREQUIRED>Y
</CLIENTENROLL>
<CHGUSERINFO>N
<AVAILACCTS>Y
<CLIENTACTREQ>Y
</SIGNUPMSGSETV1>
</SIGNUPMSGSET>
<INVSTMTMSGSET>
<INVSTMTMSGSETV1>
<MSGSETCORE>
<VER>1
<URL>https://ofx.example.com/cgi-ofx/exampleofx
<OFXSEC>NONE
<TRANSPSEC>Y
<SIGNONREALM>Example Trade
<LANGUAGE>ENG
<SYNCMODE>LITE
<RESPFILEER>N
<INTU.TIMEOUT>300
</MSGSETCORE>
<TRANDNLD>Y
<OODNLD>N
<POSDNLD>Y
<BALDNLD>Y
<CANEMAIL>N
</INVSTMTMSGSETV1>
</INVSTMTMSGSET>
<SECLISTMSGSET>
<SECLISTMSGSETV1>
<MSGSETCORE>
<VER>1
<URL>https://ofx.example.com/cgi-ofx/exampleofx
<OFXSEC>NONE
<TRANSPSEC>Y
<SIGNONREALM>Example Trade
<LANGUAGE>ENG
<SYNCMODE>LITE
<RESPFILEER>N
<INTU.TIMEOUT>300
</MSGSETCORE>
<SECLISTRQDNLD>Y
</SECLISTMSGSETV1>
</SECLISTMSGSET>
<PROFMSGSET>
<PROFMSGSETV1>
<MSGSETCORE>
<VER>1
<URL>https://ofx.example.com/cgi-ofx/exampleofx
<OFXSEC>NONE
<TRANSPSEC>Y
<SIGNONREALM>Example Trade
<LANGUAGE>ENG
<SYNCMODE>LITE
<RESPFILEER>N
<INTU.TIMEOUT>300
</MSGSETCORE>
</PROFMSGSETV1>
</PROFMSGSET>
</MSGSETLIST>
<SIGNONINFOLIST>
<SIGNONINFO>
<SIGNONREALM>Example Trade
<MIN>1
<MAX>32
<CHARTYPE>ALPHAORNUMERIC
<CASESEN>N
<SPECIAL>Y
<SPACES>N
<PINCH>N
<CHGPINFIRST>N
</SIGNONINFO>
</SIGNONINFOLIST>
<DTPROFUP>20021119140000
<FINAME>Example Trade Financial
<ADDR1>5555 Buhunkus Drive
<CITY>Someville
<STATE>NC
<POSTALCODE>28801
<COUNTRY>USA
<CSPHONE>1-800-234-5678
<TSPHONE>1-800-234-5678
<FAXPHONE>1-888-234-5678
<URL>http://www.example.com
<EMAIL>service@example.com
<INTU.BROKERID>example.com
</PROFRS>
</PROFTRNRS>
</PROFMSGSRSV1>
</OFX>`)
	var expected ofxgo.Response

	expected.Version = "102"
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = *ofxgo.NewDateGMT(2017, 4, 3, 9, 34, 58, 0)
	expected.Signon.Language = "ENG"
	expected.Signon.DtProfUp = ofxgo.NewDateGMT(2002, 11, 19, 14, 0, 0, 0)

	profileResponse := ofxgo.ProfileResponse{
		TrnUID: "0f94ce83-13b7-7568-e4fc-c02c7b47e7ab",
		Status: ofxgo.Status{
			Code:     0,
			Severity: "INFO",
		},
		MessageSetList: ofxgo.MessageSetList{
			ofxgo.MessageSet{
				Name:        "SIGNONMSGSETV1",
				Ver:         1,
				URL:         "https://ofx.example.com/cgi-ofx/exampleofx",
				OfxSec:      ofxgo.OfxSecNone,
				TranspSec:   true,
				SignonRealm: "Example Trade",
				Language:    []ofxgo.String{"ENG"},
				SyncMode:    ofxgo.SyncModeLite,
				RespFileER:  false,
				// Ignored: <INTU.TIMEOUT>300
			},
			ofxgo.MessageSet{
				Name:        "SIGNUPMSGSETV1",
				Ver:         1,
				URL:         "https://ofx.example.com/cgi-ofx/exampleofx",
				OfxSec:      ofxgo.OfxSecNone,
				TranspSec:   true,
				SignonRealm: "Example Trade",
				Language:    []ofxgo.String{"ENG"},
				SyncMode:    ofxgo.SyncModeLite,
				RespFileER:  false,
				// Ignored: <INTU.TIMEOUT>300
			},
			ofxgo.MessageSet{
				Name:        "INVSTMTMSGSETV1",
				Ver:         1,
				URL:         "https://ofx.example.com/cgi-ofx/exampleofx",
				OfxSec:      ofxgo.OfxSecNone,
				TranspSec:   true,
				SignonRealm: "Example Trade",
				Language:    []ofxgo.String{"ENG"},
				SyncMode:    ofxgo.SyncModeLite,
				RespFileER:  false,
				// Ignored: <INTU.TIMEOUT>300
			},
			ofxgo.MessageSet{
				Name:        "SECLISTMSGSETV1",
				Ver:         1,
				URL:         "https://ofx.example.com/cgi-ofx/exampleofx",
				OfxSec:      ofxgo.OfxSecNone,
				TranspSec:   true,
				SignonRealm: "Example Trade",
				Language:    []ofxgo.String{"ENG"},
				SyncMode:    ofxgo.SyncModeLite,
				RespFileER:  false,
				// Ignored: <INTU.TIMEOUT>300
			},
			ofxgo.MessageSet{
				Name:        "PROFMSGSETV1",
				Ver:         1,
				URL:         "https://ofx.example.com/cgi-ofx/exampleofx",
				OfxSec:      ofxgo.OfxSecNone,
				TranspSec:   true,
				SignonRealm: "Example Trade",
				Language:    []ofxgo.String{"ENG"},
				SyncMode:    ofxgo.SyncModeLite,
				RespFileER:  false,
				// Ignored: <INTU.TIMEOUT>300
			},
		},
		SignonInfoList: []ofxgo.SignonInfo{
			{
				SignonRealm: "Example Trade",
				Min:         1,
				Max:         32,
				CharType:    ofxgo.CharTypeAlphaOrNumeric,
				CaseSen:     false,
				Special:     true,
				Spaces:      false,
				PinCh:       false,
				ChgPinFirst: false,
			},
		},
		DtProfUp:   *ofxgo.NewDateGMT(2002, 11, 19, 14, 0, 0, 0),
		FiName:     "Example Trade Financial",
		Addr1:      "5555 Buhunkus Drive",
		City:       "Someville",
		State:      "NC",
		PostalCode: "28801",
		Country:    "USA",
		CsPhone:    "1-800-234-5678",
		TsPhone:    "1-800-234-5678",
		FaxPhone:   "1-888-234-5678",
		URL:        "http://www.example.com",
		Email:      "service@example.com",
		// Ignored: <INTU.BROKERID>example.com
	}
	expected.Prof = append(expected.Prof, &profileResponse)

	response, err := ofxgo.ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
}
