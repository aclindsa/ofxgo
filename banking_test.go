package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"testing"
	"time"
)

func TestMarshalBankStatementRequest(t *testing.T) {
	var expectedString string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20060115112300.000[-5:EST]</DTCLIENT>
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
	<BANKMSGSRQV1>
		<STMTTRNRQ>
			<TRNUID>123</TRNUID>
			<STMTRQ>
				<BANKACCTFROM>
					<BANKID>318398732</BANKID>
					<ACCTID>78346129</ACCTID>
					<ACCTTYPE>CHECKING</ACCTTYPE>
				</BANKACCTFROM>
				<INCTRAN>
					<INCLUDE>Y</INCLUDE>
				</INCTRAN>
			</STMTRQ>
		</STMTTRNRQ>
	</BANKMSGSRQV1>
</OFX>`

	var client = ofxgo.Client{
		AppId:       "OFXGO",
		AppVer:      "0001",
		SpecVersion: "203",
	}

	var request ofxgo.Request
	request.Signon.UserId = "myusername"
	request.Signon.UserPass = "Pa$$word"
	request.Signon.Org = "BNK"
	request.Signon.Fid = "1987"

	statementRequest := ofxgo.StatementRequest{
		TrnUID: "123",
		BankAcctFrom: ofxgo.BankAcct{
			BankId:   "318398732",
			AcctId:   "78346129",
			AcctType: "CHECKING",
		},
		Include: true,
	}
	request.Banking = append(request.Banking, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	EST := time.FixedZone("EST", -5*60*60)
	request.Signon.DtClient = ofxgo.Date(time.Date(2006, 1, 15, 11, 23, 0, 0, EST))

	marshalCheckRequest(t, &request, expectedString)
}
