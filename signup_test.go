package ofxgo

import (
	"strings"
	"testing"
	"time"
)

func TestMarshalAcctInfoRequest(t *testing.T) {
	var expectedString string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20160115112300.000[-5:EST]</DTCLIENT>
			<USERID>myusername</USERID>
			<USERPASS>Pa$$word</USERPASS>
			<LANGUAGE>ENG</LANGUAGE>
			<FI>
				<ORG>BNK</ORG>
				<FID>1987</FID>
			</FI>
			<APPID>OFXGO</APPID>
			<APPVER>0001</APPVER>
		</SONRQ>
	</SIGNONMSGSRQV1>
	<SIGNUPMSGSRQV1>
		<ACCTINFOTRNRQ>
			<TRNUID>e3ad9bda-38fa-4e5b-8099-1bd567ddef7a</TRNUID>
			<ACCTINFORQ>
				<DTACCTUP>20151221182945.000[-5:EST]</DTACCTUP>
			</ACCTINFORQ>
		</ACCTINFOTRNRQ>
	</SIGNUPMSGSRQV1>
</OFX>`

	EST := time.FixedZone("EST", -5*60*60)

	var client = BasicClient{
		AppID:       "OFXGO",
		AppVer:      "0001",
		SpecVersion: OfxVersion203,
	}

	var request Request
	request.Signon.UserID = "myusername"
	request.Signon.UserPass = "Pa$$word"
	request.Signon.Org = "BNK"
	request.Signon.Fid = "1987"

	acctInfoRequest := AcctInfoRequest{
		TrnUID:   "e3ad9bda-38fa-4e5b-8099-1bd567ddef7a",
		DtAcctUp: *NewDate(2015, 12, 21, 18, 29, 45, 0, EST),
	}
	request.Signup = append(request.Signup, &acctInfoRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	request.Signon.DtClient = *NewDate(2016, 1, 15, 11, 23, 0, 0, EST)

	marshalCheckRequest(t, &request, expectedString)
}

func TestUnmarshalAcctInfoResponse(t *testing.T) {
	responseReader := strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRSV1>
		<SONRS>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>INFO</SEVERITY>
			</STATUS>
			<DTSERVER>20060115112303</DTSERVER>
			<LANGUAGE>ENG</LANGUAGE>
			<DTPROFUP>20050221091300</DTPROFUP>
			<DTACCTUP>20060102160000</DTACCTUP>
			<FI>
				<ORG>BNK</ORG>
				<FID>1987</FID>
			</FI>
		</SONRS>
	</SIGNONMSGSRSV1>
	<SIGNUPMSGSRSV1>
		<ACCTINFOTRNRS>
			<TRNUID>10938754</TRNUID>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>INFO</SEVERITY>
			</STATUS>
			<ACCTINFORS>
				<DTACCTUP>20050228</DTACCTUP>
				<ACCTINFO>
					<DESC>Personal Checking</DESC>
					<PHONE>888-222-5827</PHONE>
					<BANKACCTINFO>
						<BANKACCTFROM>
							<BANKID>8367556009</BANKID>
							<ACCTID>000999847</ACCTID>
							<ACCTTYPE>MONEYMRKT</ACCTTYPE>
						</BANKACCTFROM>
						<SUPTXDL>Y</SUPTXDL>
						<XFERSRC>Y</XFERSRC>
						<XFERDEST>Y</XFERDEST>
						<SVCSTATUS>ACTIVE</SVCSTATUS>
					</BANKACCTINFO>
				</ACCTINFO>
			</ACCTINFORS>
		</ACCTINFOTRNRS>
	</SIGNUPMSGSRSV1>
</OFX>`)
	var expected Response

	expected.Version = OfxVersion203
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = *NewDateGMT(2006, 1, 15, 11, 23, 03, 0)
	expected.Signon.Language = "ENG"
	expected.Signon.DtProfUp = NewDateGMT(2005, 2, 21, 9, 13, 0, 0)
	expected.Signon.DtAcctUp = NewDateGMT(2006, 1, 2, 16, 0, 0, 0)
	expected.Signon.Org = "BNK"
	expected.Signon.Fid = "1987"

	bankacctinfo := BankAcctInfo{
		BankAcctFrom: BankAcct{
			BankID:   "8367556009",
			AcctID:   "000999847",
			AcctType: AcctTypeMoneyMrkt,
		},
		SupTxDl:   true,
		XferSrc:   true,
		XferDest:  true,
		SvcStatus: SvcStatusActive,
	}

	acctInfoResponse := AcctInfoResponse{
		TrnUID: "10938754",
		Status: Status{
			Code:     0,
			Severity: "INFO",
		},
		DtAcctUp: *NewDateGMT(2005, 2, 28, 0, 0, 0, 0),
		AcctInfo: []AcctInfo{{
			Desc:         "Personal Checking",
			Phone:        "888-222-5827",
			BankAcctInfo: &bankacctinfo,
		}},
	}
	expected.Signup = append(expected.Signup, &acctInfoResponse)

	response, err := ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
	checkResponseRoundTrip(t, response)
}
