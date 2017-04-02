package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestMarshalInvStatementRequest(t *testing.T) {
	var expectedString string = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRQV1>
		<SONRQ>
			<DTCLIENT>20160224131905.000[-5:EST]</DTCLIENT>
			<USERID>1998124</USERID>
			<USERPASS>Sup3eSekrit</USERPASS>
			<LANGUAGE>ENG</LANGUAGE>
			<FI>
				<ORG>First Bank</ORG>
				<FID>01</FID>
			</FI>
			<APPID>MYAPP</APPID>
			<APPVER>1234</APPVER>
		</SONRQ>
	</SIGNONMSGSRQV1>
	<INVSTMTMSGSRQV1>
		<INVSTMTTRNRQ>
			<TRNUID>382827d6-e2d0-4396-bf3b-665979285420</TRNUID>
			<INVSTMTRQ>
				<INVACCTFROM>
					<BROKERID>fi.example.com</BROKERID>
					<ACCTID>82736664</ACCTID>
				</INVACCTFROM>
				<INCTRAN>
					<DTSTART>20160101000000.000[-5:EST]</DTSTART>
					<INCLUDE>Y</INCLUDE>
				</INCTRAN>
				<INCOO>Y</INCOO>
				<INCPOS>
					<INCLUDE>Y</INCLUDE>
				</INCPOS>
				<INCBAL>Y</INCBAL>
			</INVSTMTRQ>
		</INVSTMTTRNRQ>
	</INVSTMTMSGSRQV1>
</OFX>`

	var client = ofxgo.Client{
		AppId:       "MYAPP",
		AppVer:      "1234",
		SpecVersion: "203",
	}

	var request ofxgo.Request
	request.Signon.UserId = "1998124"
	request.Signon.UserPass = "Sup3eSekrit"
	request.Signon.Org = "First Bank"
	request.Signon.Fid = "01"

	EST := time.FixedZone("EST", -5*60*60)
	dtstart := ofxgo.Date(time.Date(2016, 1, 1, 0, 0, 0, 0, EST))

	statementRequest := ofxgo.InvStatementRequest{
		TrnUID: "382827d6-e2d0-4396-bf3b-665979285420",
		InvAcctFrom: ofxgo.InvAcct{
			BrokerId: "fi.example.com",
			AcctId:   "82736664",
		},
		DtStart:        &dtstart,
		Include:        true,
		IncludeOO:      true,
		IncludePos:     true,
		IncludeBalance: true,
	}
	request.InvStmt = append(request.InvStmt, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	request.Signon.DtClient = ofxgo.Date(time.Date(2016, 2, 24, 13, 19, 5, 0, EST))

	marshalCheckRequest(t, &request, expectedString)
}

func TestUnmarshalInvStatementResponse(t *testing.T) {
	responseReader := strings.NewReader(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<?OFX OFXHEADER="200" VERSION="203" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>
<OFX>
	<SIGNONMSGSRSV1>
		<SONRS>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>INFO</SEVERITY>
			</STATUS>
			<DTSERVER>20170401201244</DTSERVER>
			<LANGUAGE>ENG</LANGUAGE>
			<FI>
				<ORG>INVSTRUS</ORG>
				<FID>9999</FID>
			</FI>
		</SONRS>
	</SIGNONMSGSRSV1>
	<INVSTMTMSGSRSV1>
		<INVSTMTTRNRS>
			<TRNUID>1a0117ad-692b-4c6a-a21b-020d37d34d49</TRNUID>
			<STATUS>
				<CODE>0</CODE>
				<SEVERITY>INFO</SEVERITY>
			</STATUS>
			<INVSTMTRS>
				<DTASOF>20170331000000</DTASOF>
				<CURDEF>USD</CURDEF>
				<INVACCTFROM>
					<BROKERID>invstrus.com</BROKERID>
					<ACCTID>91827364</ACCTID>
				</INVACCTFROM>
				<INVTRANLIST>
					<DTSTART>20170101000000</DTSTART>
					<DTEND>20170331000000</DTEND>
					<BUYSTOCK>
						<INVBUY>
							<INVTRAN>
								<FITID>729483191</FITID>
								<DTTRADE>20170203</DTTRADE>
								<DTSETTLE>20170207</DTSETTLE>
							</INVTRAN>
							<SECID>
								<UNIQUEID>78462F103</UNIQUEID>
								<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
							</SECID>
							<UNITS>100</UNITS>
							<UNITPRICE>229.00</UNITPRICE>
							<COMMISSION>9.00</COMMISSION>
							<TOTAL>-22909.00</TOTAL>
							<SUBACCTSEC>CASH</SUBACCTSEC>
							<SUBACCTFUND>CASH</SUBACCTFUND>
						</INVBUY>
						<BUYTYPE>BUY</BUYTYPE>
					</BUYSTOCK>
					<INVBANKTRAN>
						<STMTTRN>
							<TRNTYPE>CREDIT</TRNTYPE>
							<DTPOSTED>20170120</DTPOSTED>
							<DTUSER>20170118</DTUSER>
							<DTAVAIL>20170123</DTAVAIL>
							<TRNAMT>22000.00</TRNAMT>
							<FITID>993838</FITID>
							<NAME>DEPOSIT</NAME>
							<MEMO>CHECK 19980</MEMO>
						</STMTTRN>
						<SUBACCTFUND>CASH</SUBACCTFUND>
					</INVBANKTRAN>
				</INVTRANLIST>
				<INVPOSLIST>
					<POSSTOCK>
						<INVPOS>
							<SECID>
								<UNIQUEID>78462F103</UNIQUEID>
								<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
							</SECID>
							<HELDINACCT>CASH</HELDINACCT>
							<POSTYPE>LONG</POSTYPE>
							<UNITS>200</UNITS>
							<UNITPRICE>235.74</UNITPRICE>
							<MKTVAL>47148.00</MKTVAL>
							<DTPRICEASOF>20170331160000</DTPRICEASOF>
							<MEMO>Price as of previous close</MEMO>
						</INVPOS>
					</POSSTOCK>
					<POSOPT>
						<INVPOS>
							<SECID>
								<UNIQUEID>129887339</UNIQUEID>
								<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
							</SECID>
							<HELDINACCT>CASH</HELDINACCT>
							<POSTYPE>LONG</POSTYPE>
							<UNITS>1</UNITS>
							<UNITPRICE>3</UNITPRICE>
							<MKTVAL>300</MKTVAL>
							<DTPRICEASOF>20170331160000</DTPRICEASOF>
						</INVPOS>
					</POSOPT>
				</INVPOSLIST>
				<INVBAL>
					<AVAILCASH>16.73</AVAILCASH>
					<MARGINBALANCE>-819.20</MARGINBALANCE>
					<SHORTBALANCE>0</SHORTBALANCE>
					<BALLIST>
						<BAL>
							<NAME>Sweep Int Rate</NAME>
							<DESC>Current interest rate for sweep account balances</DESC>
							<BALTYPE>PERCENT</BALTYPE>
							<VALUE>0.25</VALUE>
							<DTASOF>20170401</DTASOF>
						</BAL>
					</BALLIST>
				</INVBAL>
				<INVOOLIST>
					<OOBUYMF>
						<OO>
							<FITID>76464632</FITID>
							<SECID>
								<UNIQUEID>922908645</UNIQUEID>
								<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
							</SECID>
							<DTPLACED>20170310124445</DTPLACED>
							<UNITS>10</UNITS>
							<SUBACCT>CASH</SUBACCT>
							<DURATION>GOODTILCANCEL</DURATION>
							<RESTRICTION>NONE</RESTRICTION>
							<LIMITPRICE>168.50</LIMITPRICE>
						</OO>
						<BUYTYPE>BUY</BUYTYPE>
						<UNITTYPE>SHARES</UNITTYPE>
					</OOBUYMF>
					<OOBUYSTOCK>
						<OO>
							<FITID>999387423</FITID>
							<SECID>
								<UNIQUEID>899422348</UNIQUEID>
								<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
							</SECID>
							<DTPLACED>20170324031900</DTPLACED>
							<UNITS>25</UNITS>
							<SUBACCT>CASH</SUBACCT>
							<DURATION>GOODTILCANCEL</DURATION>
							<RESTRICTION>ALLORNONE</RESTRICTION>
							<LIMITPRICE>19.75</LIMITPRICE>
						</OO>
						<BUYTYPE>BUY</BUYTYPE>
					</OOBUYSTOCK>
				</INVOOLIST>
			</INVSTMTRS>
		</INVSTMTTRNRS>
	</INVSTMTMSGSRSV1>
	<SECLISTMSGSRSV1>
		<SECLIST>
			<STOCKINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>78462F103</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>S&amp;P 500 ETF</SECNAME>
					<TICKER>SPY</TICKER>
					<FIID>99184</FIID>
				</SECINFO>
				<YIELD>1.92</YIELD>
				<ASSETCLASS>OTHER</ASSETCLASS>
			</STOCKINFO>
			<OPTINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>129887339</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>John's Fertilizer Puts</SECNAME>
					<TICKER>FERTP</TICKER>
					<FIID>882919</FIID>
				</SECINFO>
				<OPTTYPE>PUT</OPTTYPE>
				<STRIKEPRICE>79.00</STRIKEPRICE>
				<DTEXPIRE>20170901</DTEXPIRE>
				<SHPERCTRCT>100</SHPERCTRCT>
				<SECID>
					<UNIQUEID>983322180</UNIQUEID>
					<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
				</SECID>
				<ASSETCLASS>LARGESTOCK</ASSETCLASS>
			</OPTINFO>
			<STOCKINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>899422348</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>Whatchamacallit, Inc.</SECNAME>
					<TICKER>WHAT</TICKER>
					<FIID>883897</FIID>
				</SECINFO>
				<YIELD>17</YIELD>
				<ASSETCLASS>SMALLSTOCK</ASSETCLASS>
			</STOCKINFO>
			<MFINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>922908645</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>Mid-Cap Index Fund Admiral Shares</SECNAME>
					<TICKER>VIMAX</TICKER>
				</SECINFO>
			</MFINFO>
		</SECLIST>
	</SECLISTMSGSRSV1>
</OFX>`)
	var expected ofxgo.Response
	GMT := time.FixedZone("GMT", 0)

	expected.Version = "203"
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = ofxgo.Date(time.Date(2017, 4, 1, 20, 12, 44, 0, GMT))
	expected.Signon.Language = "ENG"
	expected.Signon.Org = "INVSTRUS"
	expected.Signon.Fid = "9999"

	var units1, unitprice1, commission1, total1, amount2 big.Rat
	units1.SetFrac64(100, 1)
	unitprice1.SetFrac64(229, 1)
	commission1.SetFrac64(9, 1)
	total1.SetFrac64(-22909, 1)
	amount2.SetFrac64(22000, 1)

	dtuser := ofxgo.Date(time.Date(2017, 1, 18, 0, 0, 0, 0, GMT))
	dtavail := ofxgo.Date(time.Date(2017, 1, 23, 0, 0, 0, 0, GMT))
	dtsettle := ofxgo.Date(time.Date(2017, 2, 7, 0, 0, 0, 0, GMT))

	invtranlist := ofxgo.InvTranList{
		DtStart: ofxgo.Date(time.Date(2017, 1, 1, 0, 0, 0, 0, GMT)),
		DtEnd:   ofxgo.Date(time.Date(2017, 3, 31, 0, 0, 0, 0, GMT)),
		InvTransactions: []ofxgo.InvTransaction{
			ofxgo.BuyStock{
				InvBuy: ofxgo.InvBuy{
					InvTran: ofxgo.InvTran{
						FiTId:    "729483191",
						DtTrade:  ofxgo.Date(time.Date(2017, 2, 3, 0, 0, 0, 0, GMT)),
						DtSettle: &dtsettle,
					},
					SecId: ofxgo.SecurityId{
						UniqueId:     "78462F103",
						UniqueIdType: "CUSIP",
					},
					Units:       ofxgo.Amount(units1),
					UnitPrice:   ofxgo.Amount(unitprice1),
					Commission:  ofxgo.Amount(commission1),
					Total:       ofxgo.Amount(total1),
					SubAcctSec:  "CASH",
					SubAcctFund: "CASH",
				},
				BuyType: "BUY",
			},
		},
		BankTransactions: []ofxgo.InvBankTransaction{
			ofxgo.InvBankTransaction{
				Transactions: []ofxgo.Transaction{
					ofxgo.Transaction{
						TrnType:  "CREDIT",
						DtPosted: ofxgo.Date(time.Date(2017, 1, 20, 0, 0, 0, 0, GMT)),
						DtUser:   &dtuser,
						DtAvail:  &dtavail,
						TrnAmt:   ofxgo.Amount(amount2),
						FiTId:    "993838",
						Name:     "DEPOSIT",
						Memo:     "CHECK 19980",
					},
				},
				SubAcctFund: "CASH",
			},
		},
	}

	var availcash, marginbalance, shortbalance, balvalue big.Rat
	availcash.SetFrac64(1673, 100)
	marginbalance.SetFrac64(-8192, 10)
	shortbalance.SetFrac64(0, 1)
	balvalue.SetFrac64(25, 100)

	baldtasof := ofxgo.Date(time.Date(2017, 4, 1, 0, 0, 0, 0, GMT))

	invbalance := ofxgo.InvBalance{
		AvailCash:     ofxgo.Amount(availcash),
		MarginBalance: ofxgo.Amount(marginbalance),
		ShortBalance:  ofxgo.Amount(shortbalance),
		BalList: []ofxgo.Balance{
			ofxgo.Balance{
				Name:    "Sweep Int Rate",
				Desc:    "Current interest rate for sweep account balances",
				BalType: "PERCENT",
				Value:   ofxgo.Amount(balvalue),
				DtAsOf:  &baldtasof,
			},
		},
	}

	var balamt, availbalamt, posunits1, posunitprice1, posmktval1, posunits2, posunitprice2, posmktval2, oounits1, oolimitprice1, oounits2, oolimitprice2 big.Rat
	balamt.SetFrac64(20029, 100)
	availbalamt.SetFrac64(20029, 100)
	posunits1.SetFrac64(200, 1)
	posunitprice1.SetFrac64(23574, 100)
	posmktval1.SetFrac64(47148, 1)
	posunits2.SetFrac64(1, 1)
	posunitprice2.SetFrac64(3, 1)
	posmktval2.SetFrac64(300, 1)
	oounits1.SetFrac64(10, 1)
	oolimitprice1.SetFrac64(16850, 100)
	oounits2.SetFrac64(25, 1)
	oolimitprice2.SetFrac64(1975, 100)

	statementResponse := ofxgo.InvStatementResponse{
		TrnUID: "1a0117ad-692b-4c6a-a21b-020d37d34d49",
		Status: ofxgo.Status{
			Code:     0,
			Severity: "INFO",
		},
		DtAsOf: ofxgo.Date(time.Date(2017, 3, 31, 0, 0, 0, 0, GMT)),
		CurDef: "USD",
		InvAcctFrom: ofxgo.InvAcct{
			BrokerId: "invstrus.com",
			AcctId:   "91827364",
		},
		InvTranList: &invtranlist,
		InvPosList: ofxgo.PositionList{
			ofxgo.StockPosition{
				InvPos: ofxgo.InvPosition{
					SecId: ofxgo.SecurityId{
						UniqueId:     "78462F103",
						UniqueIdType: "CUSIP",
					},
					HeldInAcct:  "CASH",
					PosType:     "LONG",
					Units:       ofxgo.Amount(posunits1),
					UnitPrice:   ofxgo.Amount(posunitprice1),
					MktVal:      ofxgo.Amount(posmktval1),
					DtPriceAsOf: ofxgo.Date(time.Date(2017, 3, 31, 16, 0, 0, 0, GMT)),
					Memo:        "Price as of previous close",
				},
			},
			ofxgo.OptPosition{
				InvPos: ofxgo.InvPosition{
					SecId: ofxgo.SecurityId{
						UniqueId:     "129887339",
						UniqueIdType: "CUSIP",
					},
					HeldInAcct:  "CASH",
					PosType:     "LONG",
					Units:       ofxgo.Amount(posunits2),
					UnitPrice:   ofxgo.Amount(posunitprice2),
					MktVal:      ofxgo.Amount(posmktval2),
					DtPriceAsOf: ofxgo.Date(time.Date(2017, 3, 31, 16, 0, 0, 0, GMT)),
				},
			},
		},
		InvBal: &invbalance,
		InvOOList: ofxgo.OOList{
			ofxgo.OOBuyMF{
				OO: ofxgo.OO{
					FiTId: "76464632",
					SecId: ofxgo.SecurityId{
						UniqueId:     "922908645",
						UniqueIdType: "CUSIP",
					},
					DtPlaced:    ofxgo.Date(time.Date(2017, 3, 10, 12, 44, 45, 0, GMT)),
					Units:       ofxgo.Amount(oounits1),
					SubAcct:     "CASH",
					Duration:    "GOODTILCANCEL",
					Restriction: "NONE",
					LimitPrice:  ofxgo.Amount(oolimitprice1),
				},
				BuyType:  "BUY",
				UnitType: "SHARES",
			},
			ofxgo.OOBuyStock{
				OO: ofxgo.OO{
					FiTId: "999387423",
					SecId: ofxgo.SecurityId{
						UniqueId:     "899422348",
						UniqueIdType: "CUSIP",
					},
					DtPlaced:    ofxgo.Date(time.Date(2017, 3, 24, 3, 19, 0, 0, GMT)),
					Units:       ofxgo.Amount(oounits2),
					SubAcct:     "CASH",
					Duration:    "GOODTILCANCEL",
					Restriction: "ALLORNONE",
					LimitPrice:  ofxgo.Amount(oolimitprice2),
				},
				BuyType: "BUY",
			},
		},
	}
	expected.InvStmt = append(expected.InvStmt, &statementResponse)

	var yield1, yield2, strikeprice big.Rat
	yield1.SetFrac64(192, 100)
	yield2.SetFrac64(17, 1)
	strikeprice.SetFrac64(79, 1)

	seclist := ofxgo.SecurityList{
		Securities: []ofxgo.Security{
			ofxgo.StockInfo{
				SecInfo: ofxgo.SecInfo{
					SecId: ofxgo.SecurityId{
						UniqueId:     "78462F103",
						UniqueIdType: "CUSIP",
					},
					SecName: "S&P 500 ETF",
					Ticker:  "SPY",
					FiId:    "99184",
				},
				Yield:      ofxgo.Amount(yield1),
				AssetClass: "OTHER",
			},
			ofxgo.OptInfo{
				SecInfo: ofxgo.SecInfo{
					SecId: ofxgo.SecurityId{
						UniqueId:     "129887339",
						UniqueIdType: "CUSIP",
					},
					SecName: "John's Fertilizer Puts",
					Ticker:  "FERTP",
					FiId:    "882919",
				},
				OptType:     "PUT",
				StrikePrice: ofxgo.Amount(strikeprice),
				DtExpire:    ofxgo.Date(time.Date(2017, 9, 1, 0, 0, 0, 0, GMT)),
				ShPerCtrct:  100,
				SecId: &ofxgo.SecurityId{
					UniqueId:     "983322180",
					UniqueIdType: "CUSIP",
				},
				AssetClass: "LARGESTOCK",
			},
			ofxgo.StockInfo{
				SecInfo: ofxgo.SecInfo{
					SecId: ofxgo.SecurityId{
						UniqueId:     "899422348",
						UniqueIdType: "CUSIP",
					},
					SecName: "Whatchamacallit, Inc.",
					Ticker:  "WHAT",
					FiId:    "883897",
				},
				Yield:      ofxgo.Amount(yield2),
				AssetClass: "SMALLSTOCK",
			},
			ofxgo.MFInfo{
				SecInfo: ofxgo.SecInfo{
					SecId: ofxgo.SecurityId{
						UniqueId:     "922908645",
						UniqueIdType: "CUSIP",
					},
					SecName: "Mid-Cap Index Fund Admiral Shares",
					Ticker:  "VIMAX",
				},
			},
		},
	}
	expected.SecList = append(expected.SecList, &seclist)

	response, err := ofxgo.ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
}
