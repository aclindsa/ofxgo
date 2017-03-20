package ofxgo

import (
	"errors"
	"github.com/golang/go/src/encoding/xml"
)

type InvStatementRequest struct {
	XMLName   xml.Name `xml:"INVSTMTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	InvAcctFrom      InvAcct `xml:"INVSTMTRQ>INVACCTFROM"`
	DtStart          Date    `xml:"INVSTMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd            Date    `xml:"INVSTMTRQ>INCTRAN>DTEND,omitempty"`
	Include          Boolean `xml:"INVSTMTRQ>INCTRAN>INCLUDE"`         // Include transactions (instead of just balance)
	IncludeOO        Boolean `xml:"INVSTMTRQ>INCOO"`                   // Include open orders
	PosDtAsOf        Date    `xml:"INVSTMTRQ>INCPOS>DTASOF,omitempty"` // Date that positions should be sent down for, if present
	IncludePos       Boolean `xml:"INVSTMTRQ>INCPOS>INCLUDE"`          // Include position data in response
	IncludeBalance   Boolean `xml:"INVSTMTRQ>INCBAL"`                  // Include investment balance in response
	Include401K      Boolean `xml:"INVSTMTRQ>INC401K,omitempty"`       // Include 401k information
	Include401KBal   Boolean `xml:"INVSTMTRQ>INC401KBAL,omitempty"`    // Include 401k balance information
	IncludeTranImage Boolean `xml:"INVSTMTRQ>INCTRANIMAGE,omitempty"`  // Include transaction images
}

func (r *InvStatementRequest) Name() string {
	return "INVSTMTTRNRQ"
}

func (r *InvStatementRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	return true, nil
}

type InvBankTransaction struct {
	XMLName      xml.Name      `xml:"INVBANKTRAN"`
	Transactions []Transaction `xml:"STMTTRN,omitempty"`
	SubAcctFund  String        `xml:"SUBACCTFUND"` // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
}

type InvTransactionList struct {
	XMLName      xml.Name             `xml:"INVTRANLIST"`
	DtStart      Date                 `xml:"DTSTART"`
	DtEnd        Date                 `xml:"DTEND"`
	Transactions []InvBankTransaction `xml:"INVBANKTRAN,omitempty"`
}

type InvBalance struct {
	XMLName       xml.Name  `xml:"INVBAL"`
	AvailCash     Amount    `xml:"AVAILCASH"`     // Available cash across all sub-accounts, including sweep funds
	MarginBalance Amount    `xml:"MARGINBALANCE"` // Negative means customer has borrowed funds
	ShortBalance  Amount    `xml:"SHORTBALANCE"`  // Always positive, market value of all short positions
	BuyPower      Amount    `xml:"BUYPOWER"`
	BalList       []Balance `xml:"BALLIST>BAL,omitempty"`
}

type InvStatementResponse struct {
	XMLName   xml.Name `xml:"INVSTMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO OFXEXTENSION
	DtAsOf      Date               `xml:"INVSTMTRS>DTASOF"`
	CurDef      String             `xml:"INVSTMTRS>CURDEF"`
	InvAcctFrom InvAcct            `xml:"INVSTMTRS>INVACCTFROM"`
	InvTranList InvTransactionList `xml:"INVSTMTRS>INVTRANLIST,omitempty"`
	// TODO INVPOSLIST
	InvBal InvBalance `xml:"INVSTMTRS>INVBAL,omitempty"`
	// TODO INVOOLIST
	MktgInfo String `xml:"INVSTMTRS>MKTGINFO,omitempty"` // Marketing information
	// TODO INV401K
	// TODO INV401KBAL
}

func (sr InvStatementResponse) Name() string {
	return "INVSTMTTRNRS"
}

func (sr InvStatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func DecodeInvestmentsMessageSet(d *xml.Decoder, start xml.StartElement) ([]Message, error) {
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
			case "INVSTMTTRNRS":
				var info InvStatementResponse
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return nil, err
				}
				msgs = append(msgs, Message(info))
			default:
				return nil, errors.New("Unsupported investments response tag: " + startElement.Name.Local)
			}
		} else {
			return nil, errors.New("Didn't find an opening element")
		}
	}
}
