package ofxgo

import (
	"strings"
	"testing"
	"time"
)

func TestMarshalCCStatementRequest(t *testing.T) {
	var expectedString string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20170331153848.000[0:GMT]</DTCLIENT>
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
	<CREDITCARDMSGSRQV1>
		<CCSTMTTRNRQ>
			<TRNUID>913846</TRNUID>
			<CCSTMTRQ>
				<CCACCTFROM>
					<ACCTID>XXXXXXXXXXXX1234</ACCTID>
				</CCACCTFROM>
				<INCTRAN>
					<DTSTART>20170101000000.000[0:GMT]</DTSTART>
					<INCLUDE>Y</INCLUDE>
				</INCTRAN>
			</CCSTMTRQ>
		</CCSTMTTRNRQ>
	</CREDITCARDMSGSRQV1>
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

	statementRequest := CCStatementRequest{
		TrnUID: "913846",
		CCAcctFrom: CCAcct{
			AcctID: "XXXXXXXXXXXX1234",
		},
		DtStart: NewDateGMT(2017, 1, 1, 0, 0, 0, 0),
		Include: true,
	}
	request.CreditCard = append(request.CreditCard, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	request.Signon.DtClient = *NewDateGMT(2017, 3, 31, 15, 38, 48, 0)

	marshalCheckRequest(t, &request, expectedString)
}

func TestUnmarshalCCStatementResponse102(t *testing.T) {
	responseReader := strings.NewReader(`OFXHEADER:100
DATA:OFXSGML
VERSION:102
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX><SIGNONMSGSRSV1><SONRS><STATUS><CODE>0<SEVERITY>INFO<MESSAGE>SUCCESS</STATUS><DTSERVER>20170331154648.331[-4:EDT]<LANGUAGE>ENG<FI><ORG>01<FID>81729</FI></SONRS></SIGNONMSGSRSV1><CREDITCARDMSGSRSV1><CCSTMTTRNRS><TRNUID>59e850ad-7448-b4ce-4b71-29057763b306<STATUS><CODE>0<SEVERITY>INFO</STATUS><CCSTMTRS><CURDEF>USD<CCACCTFROM><ACCTID>9283744488463775</CCACCTFROM><BANKTRANLIST><DTSTART>20161201154648.688[-5:EST]<DTEND>20170331154648.688[-4:EDT]<STMTTRN><TRNTYPE>DEBIT<DTPOSTED>20170209120000[0:GMT]<TRNAMT>-7.96<FITID>2017020924435657040207171600195<NAME>SLICE OF NY</STMTTRN><STMTTRN><TRNTYPE>CREDIT<DTPOSTED>20161228120000[0:GMT]<TRNAMT>3830.46<FITID>2016122823633637200000258482730<NAME>Payment Thank You Electro</STMTTRN><STMTTRN><TRNTYPE>DEBIT<DTPOSTED>20170327120000[0:GMT]<TRNAMT>-17.7<FITID>2017032724445727085300442885680<NAME>KROGER FUEL #9999</STMTTRN></BANKTRANLIST><LEDGERBAL><BALAMT>-9334<DTASOF>20170331080000.000[-4:EDT]</LEDGERBAL><AVAILBAL><BALAMT>7630.17<DTASOF>20170331080000.000[-4:EDT]</AVAILBAL></CCSTMTRS></CCSTMTTRNRS></CREDITCARDMSGSRSV1></OFX>`)
	var expected Response
	EDT := time.FixedZone("EDT", -4*60*60)
	EST := time.FixedZone("EST", -5*60*60)

	expected.Version = OfxVersion102
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.Status.Message = "SUCCESS"
	expected.Signon.DtServer = *NewDate(2017, 3, 31, 15, 46, 48, 331000000, EDT)
	expected.Signon.Language = "ENG"
	expected.Signon.Org = "01"
	expected.Signon.Fid = "81729"

	var trnamt1, trnamt2, trnamt3 Amount
	trnamt1.SetFrac64(-796, 100)
	trnamt2.SetFrac64(383046, 100)
	trnamt3.SetFrac64(-1770, 100)

	banktranlist := TransactionList{
		DtStart: *NewDate(2016, 12, 1, 15, 46, 48, 688000000, EST),
		DtEnd:   *NewDate(2017, 3, 31, 15, 46, 48, 688000000, EDT),
		Transactions: []Transaction{
			{
				TrnType:  TrnTypeDebit,
				DtPosted: *NewDateGMT(2017, 2, 9, 12, 0, 0, 0),
				TrnAmt:   trnamt1,
				FiTID:    "2017020924435657040207171600195",
				Name:     "SLICE OF NY",
			},
			{
				TrnType:  TrnTypeCredit,
				DtPosted: *NewDateGMT(2016, 12, 28, 12, 0, 0, 0),
				TrnAmt:   trnamt2,
				FiTID:    "2016122823633637200000258482730",
				Name:     "Payment Thank You Electro",
			},
			{
				TrnType:  TrnTypeDebit,
				DtPosted: *NewDateGMT(2017, 3, 27, 12, 0, 0, 0),
				TrnAmt:   trnamt3,
				FiTID:    "2017032724445727085300442885680",
				Name:     "KROGER FUEL #9999",
			},
		},
	}

	var balamt, availbalamt Amount
	balamt.SetFrac64(-933400, 100)
	availbalamt.SetFrac64(763017, 100)

	usd, err := NewCurrSymbol("USD")
	if err != nil {
		t.Fatalf("Unexpected error creating CurrSymbol for USD\n")
	}

	statementResponse := CCStatementResponse{
		TrnUID: "59e850ad-7448-b4ce-4b71-29057763b306",
		Status: Status{
			Code:     0,
			Severity: "INFO",
		},
		CurDef: *usd,
		CCAcctFrom: CCAcct{
			AcctID: "9283744488463775",
		},
		BankTranList: &banktranlist,
		BalAmt:       balamt,
		DtAsOf:       *NewDate(2017, 3, 31, 8, 0, 0, 0, EDT),
		AvailBalAmt:  &availbalamt,
		AvailDtAsOf:  NewDate(2017, 3, 31, 8, 0, 0, 0, EDT),
	}
	expected.CreditCard = append(expected.CreditCard, &statementResponse)

	response, err := ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
	checkResponseRoundTrip(t, response)
}
