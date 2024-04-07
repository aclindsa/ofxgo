package ofxgo

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aclindsa/xml"
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

	var client = BasicClient{
		AppID:       "MYAPP",
		AppVer:      "1234",
		SpecVersion: OfxVersion203,
	}

	var request Request
	request.Signon.UserID = "1998124"
	request.Signon.UserPass = "Sup3eSekrit"
	request.Signon.Org = "First Bank"
	request.Signon.Fid = "01"

	EST := time.FixedZone("EST", -5*60*60)

	statementRequest := InvStatementRequest{
		TrnUID: "382827d6-e2d0-4396-bf3b-665979285420",
		InvAcctFrom: InvAcct{
			BrokerID: "fi.example.com",
			AcctID:   "82736664",
		},
		DtStart:        NewDate(2016, 1, 1, 0, 0, 0, 0, EST),
		Include:        true,
		IncludeOO:      true,
		IncludePos:     true,
		IncludeBalance: true,
	}
	request.InvStmt = append(request.InvStmt, &statementRequest)

	request.SetClientFields(&client)
	// Overwrite the DtClient value set by SetClientFields to time.Now()
	request.Signon.DtClient = *NewDate(2016, 2, 24, 13, 19, 5, 0, EST)

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
			<DEBTINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>99182828</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>Someone's Class B Debt</SECNAME>
				</SECINFO>
				<PARVALUE>100.29</PARVALUE>
				<DEBTTYPE>COUPON</DEBTTYPE>
				<DTCOUPON>20170901</DTCOUPON>
				<COUPONFREQ>QUARTERLY</COUPONFREQ>
			</DEBTINFO>
			<OTHERINFO>
				<SECINFO>
					<SECID>
						<UNIQUEID>88181818</UNIQUEID>
						<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
					</SECID>
					<SECNAME>Foo Bar</SECNAME>
				</SECINFO>
				<TYPEDESC>Don't know what this is</TYPEDESC>
			</OTHERINFO>
		</SECLIST>
	</SECLISTMSGSRSV1>
</OFX>`)
	var expected Response

	expected.Version = OfxVersion203
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = *NewDateGMT(2017, 4, 1, 20, 12, 44, 0)
	expected.Signon.Language = "ENG"
	expected.Signon.Org = "INVSTRUS"
	expected.Signon.Fid = "9999"

	var units1, unitprice1, commission1, total1, amount2 Amount
	units1.SetFrac64(100, 1)
	unitprice1.SetFrac64(229, 1)
	commission1.SetFrac64(9, 1)
	total1.SetFrac64(-22909, 1)
	amount2.SetFrac64(22000, 1)

	invtranlist := InvTranList{
		DtStart: *NewDateGMT(2017, 1, 1, 0, 0, 0, 0),
		DtEnd:   *NewDateGMT(2017, 3, 31, 0, 0, 0, 0),
		InvTransactions: []InvTransaction{
			BuyStock{
				InvBuy: InvBuy{
					InvTran: InvTran{
						FiTID:    "729483191",
						DtTrade:  *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
						DtSettle: NewDateGMT(2017, 2, 7, 0, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Commission:  commission1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
				BuyType: BuyTypeBuy,
			},
		},
		BankTransactions: []InvBankTransaction{
			{
				Transactions: []Transaction{
					{
						TrnType:  TrnTypeCredit,
						DtPosted: *NewDateGMT(2017, 1, 20, 0, 0, 0, 0),
						DtUser:   NewDateGMT(2017, 1, 18, 0, 0, 0, 0),
						DtAvail:  NewDateGMT(2017, 1, 23, 0, 0, 0, 0),

						TrnAmt: amount2,
						FiTID:  "993838",
						Name:   "DEPOSIT",
						Memo:   "CHECK 19980",
					},
				},
				SubAcctFund: SubAcctTypeCash,
			},
		},
	}

	var availcash, marginbalance, shortbalance, balvalue Amount
	availcash.SetFrac64(1673, 100)
	marginbalance.SetFrac64(-8192, 10)
	shortbalance.SetFrac64(0, 1)
	balvalue.SetFrac64(25, 100)

	invbalance := InvBalance{
		AvailCash:     availcash,
		MarginBalance: marginbalance,
		ShortBalance:  shortbalance,
		BalList: []Balance{
			{
				Name:    "Sweep Int Rate",
				Desc:    "Current interest rate for sweep account balances",
				BalType: BalTypePercent,
				Value:   balvalue,
				DtAsOf:  NewDateGMT(2017, 4, 1, 0, 0, 0, 0),
			},
		},
	}

	var balamt, availbalamt, posunits1, posunitprice1, posmktval1, posunits2, posunitprice2, posmktval2, oounits1, oolimitprice1, oounits2, oolimitprice2 Amount
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

	usd, err := NewCurrSymbol("USD")
	if err != nil {
		t.Fatalf("Unexpected error creating CurrSymbol for USD\n")
	}

	statementResponse := InvStatementResponse{
		TrnUID: "1a0117ad-692b-4c6a-a21b-020d37d34d49",
		Status: Status{
			Code:     0,
			Severity: "INFO",
		},
		DtAsOf: *NewDateGMT(2017, 3, 31, 0, 0, 0, 0),
		CurDef: *usd,
		InvAcctFrom: InvAcct{
			BrokerID: "invstrus.com",
			AcctID:   "91827364",
		},
		InvTranList: &invtranlist,
		InvPosList: PositionList{
			StockPosition{
				InvPos: InvPosition{
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					HeldInAcct:  SubAcctTypeCash,
					PosType:     PosTypeLong,
					Units:       posunits1,
					UnitPrice:   posunitprice1,
					MktVal:      posmktval1,
					DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
					Memo:        "Price as of previous close",
				},
			},
			OptPosition{
				InvPos: InvPosition{
					SecID: SecurityID{
						UniqueID:     "129887339",
						UniqueIDType: "CUSIP",
					},
					HeldInAcct:  SubAcctTypeCash,
					PosType:     PosTypeLong,
					Units:       posunits2,
					UnitPrice:   posunitprice2,
					MktVal:      posmktval2,
					DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
				},
			},
		},
		InvBal: &invbalance,
		InvOOList: OOList{
			OOBuyMF{
				OO: OO{
					FiTID: "76464632",
					SecID: SecurityID{
						UniqueID:     "922908645",
						UniqueIDType: "CUSIP",
					},
					DtPlaced:    *NewDateGMT(2017, 3, 10, 12, 44, 45, 0),
					Units:       oounits1,
					SubAcct:     SubAcctTypeCash,
					Duration:    DurationGoodTilCancel,
					Restriction: RestrictionNone,
					LimitPrice:  oolimitprice1,
				},
				BuyType:  BuyTypeBuy,
				UnitType: UnitTypeShares,
			},
			OOBuyStock{
				OO: OO{
					FiTID: "999387423",
					SecID: SecurityID{
						UniqueID:     "899422348",
						UniqueIDType: "CUSIP",
					},
					DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
					Units:       oounits2,
					SubAcct:     SubAcctTypeCash,
					Duration:    DurationGoodTilCancel,
					Restriction: RestrictionAllOrNone,
					LimitPrice:  oolimitprice2,
				},
				BuyType: BuyTypeBuy,
			},
		},
	}
	expected.InvStmt = append(expected.InvStmt, &statementResponse)

	var yield1, yield2, strikeprice, parvalue Amount
	yield1.SetFrac64(192, 100)
	yield2.SetFrac64(17, 1)
	strikeprice.SetFrac64(79, 1)
	parvalue.SetFrac64(10029, 100)

	seclist := SecurityList{
		Securities: []Security{
			StockInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					SecName: "S&P 500 ETF",
					Ticker:  "SPY",
					FiID:    "99184",
				},
				Yield:      yield1,
				AssetClass: AssetClassOther,
			},
			OptInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "129887339",
						UniqueIDType: "CUSIP",
					},
					SecName: "John's Fertilizer Puts",
					Ticker:  "FERTP",
					FiID:    "882919",
				},
				OptType:     OptTypePut,
				StrikePrice: strikeprice,
				DtExpire:    *NewDateGMT(2017, 9, 1, 0, 0, 0, 0),
				ShPerCtrct:  100,
				SecID: &SecurityID{
					UniqueID:     "983322180",
					UniqueIDType: "CUSIP",
				},
				AssetClass: AssetClassLargeStock,
			},
			StockInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "899422348",
						UniqueIDType: "CUSIP",
					},
					SecName: "Whatchamacallit, Inc.",
					Ticker:  "WHAT",
					FiID:    "883897",
				},
				Yield:      yield2,
				AssetClass: AssetClassSmallStock,
			},
			MFInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "922908645",
						UniqueIDType: "CUSIP",
					},
					SecName: "Mid-Cap Index Fund Admiral Shares",
					Ticker:  "VIMAX",
				},
			},
			DebtInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "99182828",
						UniqueIDType: "CUSIP",
					},
					SecName: "Someone's Class B Debt",
				},
				ParValue:   parvalue,
				DebtType:   DebtTypeCoupon,
				DtCoupon:   NewDateGMT(2017, 9, 1, 0, 0, 0, 0),
				CouponFreq: CouponFreqQuarterly,
			},
			OtherInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "88181818",
						UniqueIDType: "CUSIP",
					},
					SecName: "Foo Bar",
				},
				TypeDesc: "Don't know what this is",
			},
		},
	}
	expected.SecList = append(expected.SecList, &seclist)

	response, err := ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
	checkResponseRoundTrip(t, response)
}

func TestUnmarshalInvStatementResponse102(t *testing.T) {
	responseReader := strings.NewReader(`OFXHEADER: 100
DATA: OFXSGML
VERSION: 102
SECURITY: NONE
ENCODING: USASCII
CHARSET: 1252
COMPRESSION: NONE
OLDFILEUID: NONE
NEWFILEUID: NONE

<OFX>
 <SIGNONMSGSRSV1>
  <SONRS>
   <STATUS>
    <CODE>0
    <SEVERITY>INFO
   </STATUS>
   <DTSERVER>20170403120000
   <LANGUAGE>ENG
   <FI>
    <ORG>VV
    <FID>1000
   </FI>
   <INTU.BID>1000
  </SONRS>
 </SIGNONMSGSRSV1>
 <INVSTMTMSGSRSV1>
  <INVSTMTTRNRS>
   <TRNUID>1283719872
   <STATUS>
    <CODE>0
    <SEVERITY>INFO
   </STATUS>
   <INVSTMTRS>
    <DTASOF>20170403120000
    <CURDEF>USD
    <INVACCTFROM>
     <BROKERID>www.exampletrader.com
     <ACCTID>12341234
    </INVACCTFROM>
    <INVTRANLIST>
     <DTSTART>20161206120000
     <DTEND>20170403120000
     <SELLOPT>
      <INVSELL>
       <INVTRAN>
        <FITID>12341234-20161207-1
        <DTTRADE>20161207120000
        <DTSETTLE>20161208120000
       </INVTRAN>
       <SECID>
        <UNIQUEID>SPY161216C00226000
        <UNIQUEIDTYPE>CUSIP
       </SECID>
       <UNITS>-1.0000
       <UNITPRICE>0.3500
       <COMMISSION>8.8500
       <FEES>0.2600
       <TOTAL>200.8900
       <SUBACCTSEC>CASH
       <SUBACCTFUND>CASH
      </INVSELL>
      <OPTSELLTYPE>SELLTOOPEN
      <SHPERCTRCT>100
     </SELLOPT>
     <CLOSUREOPT>
      <INVTRAN>
       <FITID>12341234-20161215-1
       <DTTRADE>20161215120000
       <DTSETTLE>20161220120000
      </INVTRAN>
      <SECID>
       <UNIQUEID>78462F10
       <UNIQUEIDTYPE>CUSIP
      </SECID>
      <OPTACTION>ASSIGN
      <UNITS>-100.0000
      <SHPERCTRCT>100
      <SUBACCTSEC>CASH
     </CLOSUREOPT>
     <CLOSUREOPT>
      <INVTRAN>
       <FITID>12341234-20161215-2
       <DTTRADE>20161215120000
       <DTSETTLE>20161215120000
      </INVTRAN>
      <SECID>
       <UNIQUEID>SPY161216C00226000
       <UNIQUEIDTYPE>CUSIP
      </SECID>
      <OPTACTION>ASSIGN
      <UNITS>1.0000
      <SHPERCTRCT>100
      <SUBACCTSEC>CASH
     </CLOSUREOPT>
    </INVTRANLIST>
    <INVPOSLIST>
     <POSSTOCK>
      <INVPOS>
       <SECID>
        <UNIQUEID>04956010
        <UNIQUEIDTYPE>CUSIP
       </SECID>
       <HELDINACCT>CASH
       <POSTYPE>LONG
       <UNITS>100
       <UNITPRICE>79.0000
       <MKTVAL>79000
       <DTPRICEASOF>20170403120000
      </INVPOS>
     </POSSTOCK>
     <POSSTOCK>
      <INVPOS>
       <SECID>
        <UNIQUEID>36960410
        <UNIQUEIDTYPE>CUSIP
       </SECID>
       <HELDINACCT>CASH
       <POSTYPE>LONG
       <UNITS>100.00
       <UNITPRICE>29.8700
       <MKTVAL>2987.00
       <DTPRICEASOF>20170403120000
      </INVPOS>
     </POSSTOCK>
    </INVPOSLIST>
    <INVBAL>
     <AVAILCASH>0.0
     <MARGINBALANCE>-0.00
     <SHORTBALANCE>0.00
    </INVBAL>
   </INVSTMTRS>
  </INVSTMTTRNRS>
 </INVSTMTMSGSRSV1>
 <SECLISTMSGSRSV1>
  <SECLIST>
   <STOCKINFO>
    <SECINFO>
     <SECID>
      <UNIQUEID>78462F10
      <UNIQUEIDTYPE>CUSIP
     </SECID>
     <SECNAME>SPDR S&amp;P 500 ETF TRUST
     <TICKER>SPY
    </SECINFO>
   </STOCKINFO>
   <OPTINFO>
    <SECINFO>
     <SECID>
      <UNIQUEID>SPY161216C00226000
      <UNIQUEIDTYPE>CUSIP
     </SECID>
     <SECNAME>SPY Dec 16 2016 226.00 Call
     <TICKER>SPY   161216C00226000
    </SECINFO>
    <OPTTYPE>CALL
    <STRIKEPRICE>226.00
    <DTEXPIRE>20161216120000
    <SHPERCTRCT>100
   </OPTINFO>
  </SECLIST>
 </SECLISTMSGSRSV1>
</OFX>`)
	var expected Response

	expected.Version = OfxVersion102
	expected.Signon.Status.Code = 0
	expected.Signon.Status.Severity = "INFO"
	expected.Signon.DtServer = *NewDateGMT(2017, 4, 3, 12, 0, 0, 0)
	expected.Signon.Language = "ENG"
	expected.Signon.Org = "VV"
	expected.Signon.Fid = "1000"
	// Ignored <INTU.BID>1000

	var units1, unitprice1, commission1, fees1, total1, units2, units3 Amount
	units1.SetFrac64(-1, 1)
	unitprice1.SetFrac64(35, 100)
	commission1.SetFrac64(885, 100)
	fees1.SetFrac64(26, 100)
	total1.SetFrac64(20089, 100)
	units2.SetFrac64(-100, 1)
	units3.SetFrac64(1, 1)

	invtranlist := InvTranList{
		DtStart: *NewDateGMT(2016, 12, 6, 12, 0, 0, 0),
		DtEnd:   *NewDateGMT(2017, 4, 3, 12, 0, 0, 0),
		InvTransactions: []InvTransaction{
			SellOpt{
				InvSell: InvSell{
					InvTran: InvTran{
						FiTID:    "12341234-20161207-1",
						DtTrade:  *NewDateGMT(2016, 12, 7, 12, 0, 0, 0),
						DtSettle: NewDateGMT(2016, 12, 8, 12, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "SPY161216C00226000",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Commission:  commission1,
					Fees:        fees1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
				OptSellType: OptSellTypeSellToOpen,
				ShPerCtrct:  100,
			},
			ClosureOpt{
				InvTran: InvTran{
					FiTID:    "12341234-20161215-1",
					DtTrade:  *NewDateGMT(2016, 12, 15, 12, 0, 0, 0),
					DtSettle: NewDateGMT(2016, 12, 20, 12, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "78462F10",
					UniqueIDType: "CUSIP",
				},
				OptAction:  OptActionAssign,
				Units:      units2,
				ShPerCtrct: 100,
				SubAcctSec: SubAcctTypeCash,
			},
			ClosureOpt{
				InvTran: InvTran{
					FiTID:    "12341234-20161215-2",
					DtTrade:  *NewDateGMT(2016, 12, 15, 12, 0, 0, 0),
					DtSettle: NewDateGMT(2016, 12, 15, 12, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "SPY161216C00226000",
					UniqueIDType: "CUSIP",
				},
				OptAction:  OptActionAssign,
				Units:      units3,
				ShPerCtrct: 100,
				SubAcctSec: SubAcctTypeCash,
			},
		},
	}

	var availcash, marginbalance, shortbalance Amount
	availcash.SetFrac64(0, 1)
	marginbalance.SetFrac64(-0, 1)
	shortbalance.SetFrac64(0, 1)

	invbalance := InvBalance{
		AvailCash:     availcash,
		MarginBalance: marginbalance,
		ShortBalance:  shortbalance,
	}

	var posunits1, posunitprice1, posmktval1, posunits2, posunitprice2, posmktval2 Amount
	posunits1.SetFrac64(100, 1)
	posunitprice1.SetFrac64(79, 1)
	posmktval1.SetFrac64(79000, 1)
	posunits2.SetFrac64(100, 1)
	posunitprice2.SetFrac64(2987, 100)
	posmktval2.SetFrac64(2987, 1)

	usd, err := NewCurrSymbol("USD")
	if err != nil {
		t.Fatalf("Unexpected error creating CurrSymbol for USD\n")
	}

	statementResponse := InvStatementResponse{
		TrnUID: "1283719872",
		Status: Status{
			Code:     0,
			Severity: "INFO",
		},
		DtAsOf: *NewDateGMT(2017, 4, 3, 12, 0, 0, 0),
		CurDef: *usd,
		InvAcctFrom: InvAcct{
			BrokerID: "www.exampletrader.com",
			AcctID:   "12341234",
		},
		InvTranList: &invtranlist,
		InvPosList: PositionList{
			StockPosition{
				InvPos: InvPosition{
					SecID: SecurityID{
						UniqueID:     "04956010",
						UniqueIDType: "CUSIP",
					},
					HeldInAcct:  SubAcctTypeCash,
					PosType:     PosTypeLong,
					Units:       posunits1,
					UnitPrice:   posunitprice1,
					MktVal:      posmktval1,
					DtPriceAsOf: *NewDateGMT(2017, 4, 3, 12, 0, 0, 0),
				},
			},
			StockPosition{
				InvPos: InvPosition{
					SecID: SecurityID{
						UniqueID:     "36960410",
						UniqueIDType: "CUSIP",
					},
					HeldInAcct:  SubAcctTypeCash,
					PosType:     PosTypeLong,
					Units:       posunits2,
					UnitPrice:   posunitprice2,
					MktVal:      posmktval2,
					DtPriceAsOf: *NewDateGMT(2017, 4, 3, 12, 0, 0, 0),
				},
			},
		},
		InvBal: &invbalance,
	}
	expected.InvStmt = append(expected.InvStmt, &statementResponse)

	var strikeprice Amount
	strikeprice.SetFrac64(226, 1)

	seclist := SecurityList{
		Securities: []Security{
			StockInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "78462F10",
						UniqueIDType: "CUSIP",
					},
					SecName: "SPDR S&P 500 ETF TRUST",
					Ticker:  "SPY",
				},
			},
			OptInfo{
				SecInfo: SecInfo{
					SecID: SecurityID{
						UniqueID:     "SPY161216C00226000",
						UniqueIDType: "CUSIP",
					},
					SecName: "SPY Dec 16 2016 226.00 Call",
					Ticker:  "SPY   161216C00226000",
				},
				OptType:     OptTypeCall,
				StrikePrice: strikeprice,
				DtExpire:    *NewDateGMT(2016, 12, 16, 12, 0, 0, 0),
				ShPerCtrct:  100,
			},
		},
	}
	expected.SecList = append(expected.SecList, &seclist)

	response, err := ParseResponse(responseReader)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %s\n", err)
	}

	checkResponsesEqual(t, &expected, response)
	checkResponseRoundTrip(t, response)
}

func TestUnmarshalInvTranList(t *testing.T) {
	input := `<INVTRANLIST>
	<DTSTART>20170101000000</DTSTART>
	<DTEND>20170331000000</DTEND>
	<BUYDEBT>
		<INVBUY>
			<INVTRAN>
				<FITID>81818</FITID>
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
			<FEES>.26</FEES>
			<TOTAL>-22090.26</TOTAL>
			<SUBACCTSEC>CASH</SUBACCTSEC>
			<SUBACCTFUND>CASH</SUBACCTFUND>
		</INVBUY>
		<ACCRDINT>101.2</ACCRDINT>
	</BUYDEBT>
	<BUYOPT>
		<INVBUY>
			<INVTRAN>
				<FITID>81818</FITID>
				<DTTRADE>20170203</DTTRADE>
				<MEMO>Something to make a memo about</MEMO>
			</INVTRAN>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<UNITS>100</UNITS>
			<UNITPRICE>229.00</UNITPRICE>
			<TOTAL>-22090.26</TOTAL>
			<SUBACCTSEC>CASH</SUBACCTSEC>
			<SUBACCTFUND>CASH</SUBACCTFUND>
		</INVBUY>
		<OPTBUYTYPE>BUYTOOPEN</OPTBUYTYPE>
		<SHPERCTRCT>100</SHPERCTRCT>
	</BUYOPT>
	<INVEXPENSE>
		<INVTRAN>
			<FITID>129837-1111</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<SECID>
			<UNIQUEID>78462F103</UNIQUEID>
			<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
		</SECID>
		<TOTAL>0.26</TOTAL>
		<SUBACCTSEC>CASH</SUBACCTSEC>
		<SUBACCTFUND>CASH</SUBACCTFUND>
	</INVEXPENSE>
	<JRNLSEC>
		<INVTRAN>
			<FITID>129837-1112</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<SECID>
			<UNIQUEID>78462F103</UNIQUEID>
			<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
		</SECID>
		<UNITS>2300</UNITS>
		<SUBACCTTO>CASH</SUBACCTTO>
		<SUBACCTFROM>CASH</SUBACCTFROM>
	</JRNLSEC>
	<JRNLFUND>
		<INVTRAN>
			<FITID>129837-1112</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<TOTAL>2300</TOTAL>
		<SUBACCTTO>CASH</SUBACCTTO>
		<SUBACCTFROM>CASH</SUBACCTFROM>
	</JRNLFUND>
	<BUYOTHER>
		<INVBUY>
			<INVTRAN>
				<FITID>81818</FITID>
				<DTTRADE>20170203</DTTRADE>
			</INVTRAN>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<UNITS>100</UNITS>
			<UNITPRICE>229.00</UNITPRICE>
			<TOTAL>-22090.26</TOTAL>
			<SUBACCTSEC>CASH</SUBACCTSEC>
			<SUBACCTFUND>CASH</SUBACCTFUND>
		</INVBUY>
	</BUYOTHER>
	<MARGININTEREST>
		<INVTRAN>
			<FITID>129837-1112</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<TOTAL>2300</TOTAL>
		<SUBACCTFUND>CASH</SUBACCTFUND>
	</MARGININTEREST>
	<SELLDEBT>
		<INVSELL>
			<INVTRAN>
				<FITID>129837-1111</FITID>
				<DTTRADE>20170203</DTTRADE>
			</INVTRAN>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<UNITS>100</UNITS>
			<UNITPRICE>229.00</UNITPRICE>
			<TOTAL>-22090.26</TOTAL>
			<SUBACCTSEC>CASH</SUBACCTSEC>
			<SUBACCTFUND>CASH</SUBACCTFUND>
		</INVSELL>
		<SELLREASON>SELL</SELLREASON>
	</SELLDEBT>
	<RETOFCAP>
		<INVTRAN>
			<FITID>129837-1111</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<SECID>
			<UNIQUEID>78462F103</UNIQUEID>
			<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
		</SECID>
		<TOTAL>2300.00</TOTAL>
		<SUBACCTSEC>CASH</SUBACCTSEC>
		<SUBACCTFUND>CASH</SUBACCTFUND>
	</RETOFCAP>
	<SPLIT>
		<INVTRAN>
			<FITID>129837-1111</FITID>
			<DTTRADE>20170203</DTTRADE>
		</INVTRAN>
		<SECID>
			<UNIQUEID>78462F103</UNIQUEID>
			<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
		</SECID>
		<SUBACCTSEC>CASH</SUBACCTSEC>
		<OLDUNITS>100</OLDUNITS>
		<NEWUNITS>200</NEWUNITS>
		<NUMERATOR>2</NUMERATOR>
		<DENOMINATOR>1</DENOMINATOR>
	</SPLIT>
	<SELLOTHER>
		<INVSELL>
			<INVTRAN>
				<FITID>129837-1111</FITID>
				<DTTRADE>20170203</DTTRADE>
			</INVTRAN>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<UNITS>100</UNITS>
			<UNITPRICE>229.00</UNITPRICE>
			<TOTAL>-22090.26</TOTAL>
			<SUBACCTSEC>CASH</SUBACCTSEC>
			<SUBACCTFUND>CASH</SUBACCTFUND>
		</INVSELL>
	</SELLOTHER>
</INVTRANLIST>`

	var units1, unitprice1, commission1, fees1, total1, accrdint, total2, oldunits1, newunits1 Amount
	units1.SetFrac64(100, 1)
	unitprice1.SetFrac64(229, 1)
	commission1.SetFrac64(9, 1)
	fees1.SetFrac64(26, 100)
	total1.SetFrac64(-2209026, 100)
	accrdint.SetFrac64(1012, 10)
	total2.SetFrac64(2300, 1)
	oldunits1.SetFrac64(100, 1)
	newunits1.SetFrac64(200, 1)

	expected := InvTranList{
		DtStart: *NewDateGMT(2017, 1, 1, 0, 0, 0, 0),
		DtEnd:   *NewDateGMT(2017, 3, 31, 0, 0, 0, 0),
		InvTransactions: []InvTransaction{
			BuyDebt{
				InvBuy: InvBuy{
					InvTran: InvTran{
						FiTID:    "81818",
						DtTrade:  *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
						DtSettle: NewDateGMT(2017, 2, 7, 0, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Commission:  commission1,
					Fees:        fees1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
				AccrdInt: accrdint,
			},
			BuyOpt{
				InvBuy: InvBuy{
					InvTran: InvTran{
						FiTID:   "81818",
						DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
						Memo:    "Something to make a memo about",
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
				OptBuyType: OptBuyTypeBuyToOpen,
				ShPerCtrct: 100,
			},
			InvExpense{
				InvTran: InvTran{
					FiTID:   "129837-1111",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				Total:       fees1,
				SubAcctSec:  SubAcctTypeCash,
				SubAcctFund: SubAcctTypeCash,
			},
			JrnlSec{
				InvTran: InvTran{
					FiTID:   "129837-1112",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				Units:       total2,
				SubAcctTo:   SubAcctTypeCash,
				SubAcctFrom: SubAcctTypeCash,
			},
			JrnlFund{
				InvTran: InvTran{
					FiTID:   "129837-1112",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				Total:       total2,
				SubAcctTo:   SubAcctTypeCash,
				SubAcctFrom: SubAcctTypeCash,
			},
			BuyOther{
				InvBuy: InvBuy{
					InvTran: InvTran{
						FiTID:   "81818",
						DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
			},
			MarginInterest{
				InvTran: InvTran{
					FiTID:   "129837-1112",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				Total:       total2,
				SubAcctFund: SubAcctTypeCash,
			},
			SellDebt{
				InvSell: InvSell{
					InvTran: InvTran{
						FiTID:   "129837-1111",
						DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
				SellReason: SellReasonSell,
			},
			RetOfCap{
				InvTran: InvTran{
					FiTID:   "129837-1111",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				Total:       total2,
				SubAcctSec:  SubAcctTypeCash,
				SubAcctFund: SubAcctTypeCash,
			},
			Split{
				InvTran: InvTran{
					FiTID:   "129837-1111",
					DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
				},
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				SubAcctSec:  SubAcctTypeCash,
				OldUnits:    oldunits1,
				NewUnits:    newunits1,
				Numerator:   2,
				Denominator: 1,
			},
			SellOther{
				InvSell: InvSell{
					InvTran: InvTran{
						FiTID:   "129837-1111",
						DtTrade: *NewDateGMT(2017, 2, 3, 0, 0, 0, 0),
					},
					SecID: SecurityID{
						UniqueID:     "78462F103",
						UniqueIDType: "CUSIP",
					},
					Units:       units1,
					UnitPrice:   unitprice1,
					Total:       total1,
					SubAcctSec:  SubAcctTypeCash,
					SubAcctFund: SubAcctTypeCash,
				},
			},
		},
	}

	var actual InvTranList
	err := xml.Unmarshal([]byte(input), &actual)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling InvTranList: %s\n", err)
	}
	checkEqual(t, "InvTranList", reflect.ValueOf(&expected), reflect.ValueOf(&actual))
}

func TestUnmarshalPositionList(t *testing.T) {
	input := `<INVPOSLIST>
	<POSOTHER>
		<INVPOS>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<HELDINACCT>CASH</HELDINACCT>
			<POSTYPE>LONG</POSTYPE>
			<UNITS>1</UNITS>
			<UNITPRICE>3</UNITPRICE>
			<MKTVAL>300</MKTVAL>
			<DTPRICEASOF>20170331160000</DTPRICEASOF>
		</INVPOS>
	</POSOTHER>
	<POSSTOCK>
		<INVPOS>
			<SECID>
				<UNIQUEID>78462F103</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<HELDINACCT>CASH</HELDINACCT>
			<POSTYPE>SHORT</POSTYPE>
			<UNITS>200</UNITS>
			<UNITPRICE>235.74</UNITPRICE>
			<MKTVAL>47148.00</MKTVAL>
			<DTPRICEASOF>20170331160000</DTPRICEASOF>
			<MEMO>Price as of previous close</MEMO>
		</INVPOS>
		<REINVDIV>Y</REINVDIV>
	</POSSTOCK>
	<POSDEBT>
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
	</POSDEBT>
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
	<POSMF>
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
		<REINVDIV>Y</REINVDIV>
		<REINVCG>N</REINVCG>
	</POSMF>
</INVPOSLIST>`

	var posunits1, posunitprice1, posmktval1, posunits2, posunitprice2, posmktval2 Amount
	posunits1.SetFrac64(200, 1)
	posunitprice1.SetFrac64(23574, 100)
	posmktval1.SetFrac64(47148, 1)
	posunits2.SetFrac64(1, 1)
	posunitprice2.SetFrac64(3, 1)
	posmktval2.SetFrac64(300, 1)

	expected := PositionList{
		OtherPosition{
			InvPos: InvPosition{
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				HeldInAcct:  SubAcctTypeCash,
				PosType:     PosTypeLong,
				Units:       posunits2,
				UnitPrice:   posunitprice2,
				MktVal:      posmktval2,
				DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
			},
		},
		StockPosition{
			InvPos: InvPosition{
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				HeldInAcct:  SubAcctTypeCash,
				PosType:     PosTypeShort,
				Units:       posunits1,
				UnitPrice:   posunitprice1,
				MktVal:      posmktval1,
				DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
				Memo:        "Price as of previous close",
			},
			ReinvDiv: true,
		},
		DebtPosition{
			InvPos: InvPosition{
				SecID: SecurityID{
					UniqueID:     "129887339",
					UniqueIDType: "CUSIP",
				},
				HeldInAcct:  SubAcctTypeCash,
				PosType:     PosTypeLong,
				Units:       posunits2,
				UnitPrice:   posunitprice2,
				MktVal:      posmktval2,
				DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
			},
		},
		OptPosition{
			InvPos: InvPosition{
				SecID: SecurityID{
					UniqueID:     "129887339",
					UniqueIDType: "CUSIP",
				},
				HeldInAcct:  SubAcctTypeCash,
				PosType:     PosTypeLong,
				Units:       posunits2,
				UnitPrice:   posunitprice2,
				MktVal:      posmktval2,
				DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
			},
		},
		MFPosition{
			InvPos: InvPosition{
				SecID: SecurityID{
					UniqueID:     "78462F103",
					UniqueIDType: "CUSIP",
				},
				HeldInAcct:  SubAcctTypeCash,
				PosType:     PosTypeLong,
				Units:       posunits1,
				UnitPrice:   posunitprice1,
				MktVal:      posmktval1,
				DtPriceAsOf: *NewDateGMT(2017, 3, 31, 16, 0, 0, 0),
				Memo:        "Price as of previous close",
			},
			ReinvDiv: true,
			ReinvCG:  false,
		},
	}

	var actual PositionList
	err := xml.Unmarshal([]byte(input), &actual)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling PositionList: %s\n", err)
	}
	checkEqual(t, "PositionList", reflect.ValueOf(&expected), reflect.ValueOf(&actual))
}

func TestUnmarshalOOList(t *testing.T) {
	input := `<INVOOLIST>
	<OOBUYDEBT>
		<OO>
			<FITID>76464632</FITID>
			<SECID>
				<UNIQUEID>922908645</UNIQUEID>
				<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
			</SECID>
			<DTPLACED>20170310124445</DTPLACED>
			<UNITS>10</UNITS>
			<SUBACCT>CASH</SUBACCT>
			<DURATION>DAY</DURATION>
			<RESTRICTION>NONE</RESTRICTION>
		</OO>
		<AUCTION>Y</AUCTION>
	</OOBUYDEBT>
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
	<OOBUYOPT>
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
		<OPTBUYTYPE>BUYTOCLOSE</OPTBUYTYPE>
	</OOBUYOPT>
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
	<OOBUYOTHER>
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
		<UNITTYPE>CURRENCY</UNITTYPE>
	</OOBUYOTHER>
	<OOSELLDEBT>
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
		</OO>
	</OOSELLDEBT>
	<OOSELLMF>
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
		</OO>
		<SELLTYPE>SELLSHORT</SELLTYPE>
		<UNITTYPE>SHARES
		</UNITTYPE>
		<SELLALL>Y</SELLALL>
	</OOSELLMF>
	<OOSELLOPT>
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
		</OO>
		<OPTSELLTYPE>SELLTOOPEN</OPTSELLTYPE>
	</OOSELLOPT>
	<OOSELLOTHER>
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
		</OO>
		<UNITTYPE>SHARES</UNITTYPE>
	</OOSELLOTHER>
	<OOSELLSTOCK>
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
		</OO>
		<SELLTYPE>SELL</SELLTYPE>
	</OOSELLSTOCK>
	<SWITCHMF>
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
		</OO>
		<SECID>
			<UNIQUEID>899422389</UNIQUEID>
			<UNIQUEIDTYPE>CUSIP</UNIQUEIDTYPE>
		</SECID>
		<UNITTYPE>CURRENCY</UNITTYPE>
		<SWITCHALL>N</SWITCHALL>
	</SWITCHMF>
</INVOOLIST>`

	var oounits1, oolimitprice1, oounits2, oolimitprice2 Amount
	oounits1.SetFrac64(10, 1)
	oolimitprice1.SetFrac64(16850, 100)
	oounits2.SetFrac64(25, 1)
	oolimitprice2.SetFrac64(1975, 100)

	expected := OOList{
		OOBuyDebt{
			OO: OO{
				FiTID: "76464632",
				SecID: SecurityID{
					UniqueID:     "922908645",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 10, 12, 44, 45, 0),
				Units:       oounits1,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationDay,
				Restriction: RestrictionNone,
			},
			Auction: true,
		},
		OOBuyMF{
			OO: OO{
				FiTID: "76464632",
				SecID: SecurityID{
					UniqueID:     "922908645",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 10, 12, 44, 45, 0),
				Units:       oounits1,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionNone,
				LimitPrice:  oolimitprice1,
			},
			BuyType:  BuyTypeBuy,
			UnitType: UnitTypeShares,
		},
		OOBuyOpt{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
				LimitPrice:  oolimitprice2,
			},
			OptBuyType: OptBuyTypeBuyToClose,
		},
		OOBuyStock{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
				LimitPrice:  oolimitprice2,
			},
			BuyType: BuyTypeBuy,
		},
		OOBuyOther{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
				LimitPrice:  oolimitprice2,
			},
			UnitType: UnitTypeCurrency,
		},
		OOSellDebt{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
		},
		OOSellMF{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
			SellType: SellTypeSellShort,
			UnitType: UnitTypeShares,
			SellAll:  true,
		},
		OOSellOpt{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
			OptSellType: OptSellTypeSellToOpen,
		},
		OOSellOther{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
			UnitType: UnitTypeShares,
		},
		OOSellStock{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
			SellType: SellTypeSell,
		},
		OOSwitchMF{
			OO: OO{
				FiTID: "999387423",
				SecID: SecurityID{
					UniqueID:     "899422348",
					UniqueIDType: "CUSIP",
				},
				DtPlaced:    *NewDateGMT(2017, 3, 24, 3, 19, 0, 0),
				Units:       oounits2,
				SubAcct:     SubAcctTypeCash,
				Duration:    DurationGoodTilCancel,
				Restriction: RestrictionAllOrNone,
			},
			SecID: SecurityID{
				UniqueID:     "899422389",
				UniqueIDType: "CUSIP",
			},
			UnitType:  UnitTypeCurrency,
			SwitchAll: false,
		},
	}

	var actual OOList
	err := xml.Unmarshal([]byte(input), &actual)
	if err != nil {
		t.Fatalf("Unexpected error unmarshalling OOList: %s\n", err)
	}
	checkEqual(t, "OOList", reflect.ValueOf(&expected), reflect.ValueOf(&actual))
}

func TestSecurityInfo(t *testing.T) {
	secInfo := SecInfo{
		Ticker: "ABC",
	}
	tests := []Security{
		DebtInfo{SecInfo: secInfo},
		MFInfo{SecInfo: secInfo},
		OptInfo{SecInfo: secInfo},
		OtherInfo{SecInfo: secInfo},
		StockInfo{SecInfo: secInfo},
	}

	for _, tc := range tests {
		t.Run(tc.SecurityType(), func(t *testing.T) {
			info := tc.SecurityInfo()
			if info.Ticker != secInfo.Ticker {
				t.Errorf("got %v, want %v", info, secInfo)
			}
		})
	}
}

func TestInvPosition(t *testing.T) {
	invPos := InvPosition{
		Memo: "stuff",
	}
	tests := []Position{
		DebtPosition{InvPos: invPos},
		MFPosition{InvPos: invPos},
		OptPosition{InvPos: invPos},
		OtherPosition{InvPos: invPos},
		StockPosition{InvPos: invPos},
	}

	for _, tc := range tests {
		t.Run(tc.PositionType(), func(t *testing.T) {
			pos := tc.InvPosition()
			if pos.Memo != invPos.Memo {
				t.Errorf("got %v, want %v", pos, invPos)
			}
		})
	}
}

func TestInvTransaction(t *testing.T) {
	invTran := InvTran{
		Memo: "stuff",
	}
	tests := []InvTransaction{
		BuyDebt{InvBuy: InvBuy{InvTran: invTran}},
		BuyMF{InvBuy: InvBuy{InvTran: invTran}},
		BuyOpt{InvBuy: InvBuy{InvTran: invTran}},
		BuyOther{InvBuy: InvBuy{InvTran: invTran}},
		BuyStock{InvBuy: InvBuy{InvTran: invTran}},
		ClosureOpt{InvTran: invTran},
		Income{InvTran: invTran},
		InvExpense{InvTran: invTran},
		JrnlFund{InvTran: invTran},
		JrnlSec{InvTran: invTran},
		MarginInterest{InvTran: invTran},
		Reinvest{InvTran: invTran},
		RetOfCap{InvTran: invTran},
		SellDebt{InvSell: InvSell{InvTran: invTran}},
		SellMF{InvSell: InvSell{InvTran: invTran}},
		SellOpt{InvSell: InvSell{InvTran: invTran}},
		SellOther{InvSell: InvSell{InvTran: invTran}},
		SellStock{InvSell: InvSell{InvTran: invTran}},
		Split{InvTran: invTran},
		Transfer{InvTran: invTran},
	}

	for _, tc := range tests {
		t.Run(tc.TransactionType(), func(t *testing.T) {
			tran := tc.InvTransaction()
			if tran.Memo != invTran.Memo {
				t.Errorf("got %v, want %v", tran, invTran)
			}
		})
	}
}
