package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"strings"
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
	request.Bank = append(request.Bank, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	EST := time.FixedZone("EST", -5*60*60)
	request.Signon.DtClient = *ofxgo.NewDate(2006, 1, 15, 11, 23, 0, 0, EST)

	marshalCheckRequest(t, &request, expectedString)
}

func TestUnmarshalBankStatementResponse(t *testing.T) {
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
	<BANKMSGSRSV1>
		<STMTTRNRS>
			<TRNUID>1001</TRNUID>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>INFO</SEVERITY>
			</STATUS>
			<STMTRS>
				<CURDEF>USD</CURDEF>
				<BANKACCTFROM>
					<BANKID>318398732</BANKID>
					<ACCTID>78346129</ACCTID>
					<ACCTTYPE>CHECKING</ACCTTYPE>
				</BANKACCTFROM>
				<BANKTRANLIST>
					<DTSTART>20060101</DTSTART>
					<DTEND>20060115</DTEND>
					<STMTTRN>
						<TRNTYPE>CHECK</TRNTYPE>
						<DTPOSTED>20060104</DTPOSTED>
						<TRNAMT>-200.00</TRNAMT>
						<FITID>00592</FITID>
						<CHECKNUM>2002</CHECKNUM>
					</STMTTRN>
					<STMTTRN>
						<TRNTYPE>ATM</TRNTYPE>
						<DTPOSTED>20060112</DTPOSTED>
						<DTUSER>20060112</DTUSER>
						<TRNAMT>-300.00</TRNAMT>
						<FITID>00679</FITID>
					</STMTTRN>
				</BANKTRANLIST>
				<LEDGERBAL>
					<BALAMT>200.29</BALAMT>
					<DTASOF>200601141600</DTASOF>
				</LEDGERBAL>
				<AVAILBAL>
					<BALAMT>200.29</BALAMT>
					<DTASOF>200601141600</DTASOF>
				</AVAILBAL>
			</STMTRS>
		</STMTTRNRS>
	</BANKMSGSRSV1>
</OFX>`)
	var expected ofxgo.Response

	expected.Version = "203"
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = *ofxgo.NewDateGMT(2006, 1, 15, 11, 23, 03, 0)
	expected.Signon.Language = "ENG"
	expected.Signon.DtProfUp = ofxgo.NewDateGMT(2005, 2, 21, 9, 13, 0, 0)
	expected.Signon.DtAcctUp = ofxgo.NewDateGMT(2006, 1, 2, 16, 0, 0, 0)
	expected.Signon.Org = "BNK"
	expected.Signon.Fid = "1987"

	var trnamt1, trnamt2 ofxgo.Amount
	trnamt1.SetFrac64(-20000, 100)
	trnamt2.SetFrac64(-30000, 100)

	banktranlist := ofxgo.TransactionList{
		DtStart: *ofxgo.NewDateGMT(2006, 1, 1, 0, 0, 0, 0),
		DtEnd:   *ofxgo.NewDateGMT(2006, 1, 15, 0, 0, 0, 0),
		Transactions: []ofxgo.Transaction{
			{
				TrnType:  "CHECK",
				DtPosted: *ofxgo.NewDateGMT(2006, 1, 4, 0, 0, 0, 0),
				TrnAmt:   trnamt1,
				FiTId:    "00592",
				CheckNum: "2002",
			},
			{
				TrnType:  "ATM",
				DtPosted: *ofxgo.NewDateGMT(2006, 1, 12, 0, 0, 0, 0),
				DtUser:   ofxgo.NewDateGMT(2006, 1, 12, 0, 0, 0, 0),
				TrnAmt:   trnamt2,
				FiTId:    "00679",
			},
		},
	}

	var balamt, availbalamt ofxgo.Amount
	balamt.SetFrac64(20029, 100)
	availbalamt.SetFrac64(20029, 100)

	statementResponse := ofxgo.StatementResponse{
		TrnUID: "1001",
		Status: ofxgo.Status{
			Code:     0,
			Severity: "INFO",
		},
		CurDef: "USD",
		BankAcctFrom: ofxgo.BankAcct{
			BankId:   "318398732",
			AcctId:   "78346129",
			AcctType: "CHECKING",
		},
		BankTranList: &banktranlist,
		BalAmt:       balamt,
		DtAsOf:       *ofxgo.NewDateGMT(2006, 1, 14, 16, 0, 0, 0),
		AvailBalAmt:  &availbalamt,
		AvailDtAsOf:  ofxgo.NewDateGMT(2006, 1, 14, 16, 0, 0, 0),
	}
	expected.Bank = append(expected.Bank, &statementResponse)

	response, err := ofxgo.ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
}
