package ofxgo

import (
	"errors"
	"github.com/golang/go/src/encoding/xml"
)

type StatementRequest struct {
	XMLName          xml.Name `xml:"STMTTRNRQ"`
	TrnUID           UID      `xml:"TRNUID"`
	BankAcctFrom     BankAcct `xml:"STMTRQ>BANKACCTFROM"`
	DtStart          Date     `xml:"STMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd            Date     `xml:"STMTRQ>INCTRAN>DTEND,omitempty"`
	Include          Boolean  `xml:"STMTRQ>INCTRAN>INCLUDE"`            // Include transactions (instead of just balance)
	IncludePending   Boolean  `xml:"STMTRQ>INCLUDEPENDING,omitempty"`   // Include pending transactions
	IncludeTranImage Boolean  `xml:"STMTRQ>INCLUDETRANIMAGE,omitempty"` // Include transaction images
}

func (r *StatementRequest) Name() string {
	return "STMTTRNRQ"
}

func (r *StatementRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	return true, nil
}

type Payee struct {
	XMLName    xml.Name `xml:"PAYEE"`
	Name       String   `xml:"NAME"`
	Addr1      String   `xml:"ADDR1"`
	Addr2      String   `xml:"ADDR2,omitempty"`
	Addr3      String   `xml:"ADDR3,omitempty"`
	City       String   `xml:"CITY"`
	State      String   `xml:"STATE"`
	PostalCode String   `xml:"POSTALCODE"`
	Country    String   `xml:"COUNTRY,omitempty"`
	Phone      String   `xml:"PHONE"`
}

type Transaction struct {
	XMLName       xml.Name `xml:"STMTTRN"`
	TrnType       String   `xml:"TRNTYPE"` // One of CREDIT, DEBIT, INT (interest earned or paid. Note: Depends on signage of amount), DIV, FEE, SRVCHG (service charge), DEP (deposit), ATM (Note: Depends on signage of amount), POS (Note: Depends on signage of amount), XFER, CHECK, PAYMENT, CASH, DIRECTDEP, DIRECTDEBIT, REPEATPMT, HOLD (Only valid in <STMTTRNP>), OTHER
	DtPosted      Date     `xml:"DTPOSTED"`
	DtUser        Date     `xml:"DTUSER,omitempty"`
	DtAvail       Date     `xml:"DTAVAIL,omitempty"`
	TrnAmt        Amount   `xml:"TRNAMT"`
	FiTId         String   `xml:"FITID"`
	CorrectFiTId  String   `xml:"CORRECTFITID,omitempty"`  // Transaction Id that this transaction corrects, if present
	CorrectAction String   `xml:"CORRECTACTION,omitempty"` // One of DELETE, REPLACE
	SrvrTId       String   `xml:"SRVRTID,omitempty"`
	CheckNum      String   `xml:"CHECKNUM,omitempty"`
	RefNum        String   `xml:"REFNUM,omitempty"`
	SIC           Int      `xml:"SIC,omitempty"` // Standard Industrial Code
	PayeeId       String   `xml:"PAYEEID,omitempty"`
	// Note: Servers should provide NAME or PAYEE, but not both
	Name     String `xml:"NAME,omitempty"`
	Payee    Payee  `xml:"PAYEE,omitempty"`
	ExtdName String `xml:"EXTDNAME,omitempty"` // Extended name of payee or transaction description
	// TODO BANKACCTTO
	// TODO CCACCTTO
	Memo String `xml:"MEMO,omitempty"` // Extra information (not in NAME)
	// TODO IMAGEDATA
	Currency      String `xml:"CURRENCY,omitempty"`      // If different from CURDEF in STMTTRS
	OrigCurrency  String `xml:"ORIGCURRENCY,omitempty"`  // If different from CURDEF in STMTTRS
	Inv401kSource String `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST (Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST.)
}

type TransactionList struct {
	XMLName      xml.Name      `xml:"BANKTRANLIST"`
	DtStart      Date          `xml:"DTSTART"`
	DtEnd        Date          `xml:"DTEND"`
	Transactions []Transaction `xml:"STMTTRN"`
}

type StatementResponse struct {
	XMLName      xml.Name        `xml:"STMTTRNRS"`
	TrnUID       UID             `xml:"TRNUID"`
	CurDef       String          `xml:"STMTRS>CURDEF"`
	BankAcctFrom BankAcct        `xml:"STMTRS>BANKACCTFROM"`
	BankTranList TransactionList `xml:"STMTRS>BANKTRANLIST,omitempty"`
	BalAmt       Amount          `xml:"STMTRS>LEDGERBAL>BALAMT"`
	DtAsOf       Date            `xml:"STMTRS>LEDGERBAL>DTASOF"`
	// TODO AVAILBAL et. al.
}

func (sr StatementResponse) Name() string {
	return "STMTTRNRS"
}

func (sr StatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}
func DecodeBankingMessageSet(d *xml.Decoder, start xml.StartElement) ([]Message, error) {
	var msgs []Message
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return nil, err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return msgs, nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "STMTTRNRS":
				var info StatementResponse
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return nil, err
				}
				msgs = append(msgs, Message(info))
			default:
				return nil, errors.New("Unsupported banking response tag: " + startElement.Name.Local)
			}
		} else {
			return nil, errors.New("Didn't find an opening element")
		}
	}
}
