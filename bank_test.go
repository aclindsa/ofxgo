package ofxgo

import (
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

	statementRequest := StatementRequest{
		TrnUID: "123",
		BankAcctFrom: BankAcct{
			BankID:   "318398732",
			AcctID:   "78346129",
			AcctType: AcctTypeChecking,
		},
		Include: true,
	}
	request.Bank = append(request.Bank, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	EST := time.FixedZone("EST", -5*60*60)
	request.Signon.DtClient = *NewDate(2006, 1, 15, 11, 23, 0, 0, EST)

	marshalCheckRequest(t, &request, expectedString)
}

func TestMarshalBankStatementRequest103(t *testing.T) {
	var expectedString string = `OFXHEADER:100
DATA:OFXSGML
VERSION:103
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20060115112300.000[-5:EST]
			<USERID>myusername
			<USERPASS>Pa$$word
			<LANGUAGE>ENG
			<FI>
				<ORG>BNK
				<FID>1987
			</FI>
			<APPID>OFXGO
			<APPVER>0001
		</SONRQ>
	</SIGNONMSGSRQV1>
	<BANKMSGSRQV1>
		<STMTTRNRQ>
			<TRNUID>123
			<STMTRQ>
				<BANKACCTFROM>
					<BANKID>318398732
					<ACCTID>78346129
					<ACCTTYPE>CHECKING
				</BANKACCTFROM>
				<INCTRAN>
					<INCLUDE>Y
				</INCTRAN>
			</STMTRQ>
		</STMTTRNRQ>
	</BANKMSGSRQV1>
</OFX>`

	var client = BasicClient{
		AppID:       "OFXGO",
		AppVer:      "0001",
		SpecVersion: OfxVersion103,
	}

	var request Request
	request.Signon.UserID = "myusername"
	request.Signon.UserPass = "Pa$$word"
	request.Signon.Org = "BNK"
	request.Signon.Fid = "1987"

	statementRequest := StatementRequest{
		TrnUID: "123",
		BankAcctFrom: BankAcct{
			BankID:   "318398732",
			AcctID:   "78346129",
			AcctType: AcctTypeChecking,
		},
		Include: true,
	}
	request.Bank = append(request.Bank, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	EST := time.FixedZone("EST", -5*60*60)
	request.Signon.DtClient = *NewDate(2006, 1, 15, 11, 23, 0, 0, EST)

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

	var trnamt1, trnamt2 Amount
	trnamt1.SetFrac64(-20000, 100)
	trnamt2.SetFrac64(-30000, 100)

	banktranlist := TransactionList{
		DtStart: *NewDateGMT(2006, 1, 1, 0, 0, 0, 0),
		DtEnd:   *NewDateGMT(2006, 1, 15, 0, 0, 0, 0),
		Transactions: []Transaction{
			{
				TrnType:  TrnTypeCheck,
				DtPosted: *NewDateGMT(2006, 1, 4, 0, 0, 0, 0),
				TrnAmt:   trnamt1,
				FiTID:    "00592",
				CheckNum: "2002",
			},
			{
				TrnType:  TrnTypeATM,
				DtPosted: *NewDateGMT(2006, 1, 12, 0, 0, 0, 0),
				DtUser:   NewDateGMT(2006, 1, 12, 0, 0, 0, 0),
				TrnAmt:   trnamt2,
				FiTID:    "00679",
			},
		},
	}

	var balamt, availbalamt Amount
	balamt.SetFrac64(20029, 100)
	availbalamt.SetFrac64(20029, 100)

	usd, err := NewCurrSymbol("USD")
	if err != nil {
		t.Fatalf("Unexpected error creating CurrSymbol for USD\n")
	}

	statementResponse := StatementResponse{
		TrnUID: "1001",
		Status: Status{
			Code:     0,
			Severity: "INFO",
		},
		CurDef: *usd,
		BankAcctFrom: BankAcct{
			BankID:   "318398732",
			AcctID:   "78346129",
			AcctType: AcctTypeChecking,
		},
		BankTranList: &banktranlist,
		BalAmt:       balamt,
		DtAsOf:       *NewDateGMT(2006, 1, 14, 16, 0, 0, 0),
		AvailBalAmt:  &availbalamt,
		AvailDtAsOf:  NewDateGMT(2006, 1, 14, 16, 0, 0, 0),
	}
	expected.Bank = append(expected.Bank, &statementResponse)

	response, err := ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
	checkResponseRoundTrip(t, response)
}

func TestPayeeValid(t *testing.T) {
	p := Payee{
		Name:       "Jane",
		Addr1:      "Sesame Street",
		City:       "Mytown",
		State:      "AA",
		PostalCode: "12345",
		Phone:      "12345678901",
	}
	valid, err := p.Valid()
	if !valid {
		t.Fatalf("Unexpected error from calling Valid: %s\n", err)
	}

	// Ensure some empty fields trigger invalid response
	badp := p
	badp.Name = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty name\n")
	}

	badp = p
	badp.Addr1 = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty address\n")
	}

	badp = p
	badp.City = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty city\n")
	}

	badp = p
	badp.State = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty state\n")
	}

	badp = p
	badp.PostalCode = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty postal code\n")
	}

	badp = p
	badp.Phone = ""
	valid, err = badp.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty phone\n")
	}
}

func TestBalanceValid(t *testing.T) {
	var a Amount
	a.SetFrac64(8, 1)
	b := Balance{
		Name:    "Checking",
		Desc:    "Jane's Personal Checking",
		BalType: BalTypeDollar,
		Value:   a,
	}
	valid, err := b.Valid()
	if !valid {
		t.Fatalf("Unexpected error from calling Valid: %s\n", err)
	}

	badb := b
	badb.Name = ""
	valid, err = badb.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty name\n")
	}

	badb = b
	badb.Desc = ""
	valid, err = badb.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with empty description\n")
	}

	badb = Balance{
		Name:  "Checking",
		Desc:  "Jane's Personal Checking",
		Value: a,
	}
	valid, err = badb.Valid()
	if valid || err == nil {
		t.Fatalf("Expected error from calling Valid with unspecified balance type\n")
	}
}
