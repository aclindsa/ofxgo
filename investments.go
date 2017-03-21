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
	XMLName          xml.Name             `xml:"INVTRANLIST"`
	DtStart          Date                 `xml:"DTSTART"`
	DtEnd            Date                 `xml:"DTEND"`
	BankTransactions []InvBankTransaction `xml:"INVBANKTRAN,omitempty"`
}

type InvPosition struct {
	XMLName       xml.Name   `xml:"INVPOS"`
	SecId         SecurityId `xml:"SECID"`
	HeldInAcct    String     `xml:"HELDINACCT"`             // Sub-account type, one of CASH, MARGIN, SHORT, OTHER
	PosType       String     `xml:"POSTYPE"`                // SHORT = Writer for options, Short for all others; LONG = Holder for options, Long for all others.
	Units         Amount     `xml:"UNITS"`                  // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice     Amount     `xml:"UNITPRICE"`              // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	MktVal        Amount     `xml:"MKTVAL"`                 // Market value of this position
	AvgCostBasis  Amount     `xml:"AVGCOSTBASIS,omitempty"` //
	DtPriceAsOf   Date       `xml:"DTPRICEASOF"`            // Date and time of unit price and market value, and cost basis. If this date is unknown, use 19900101 as the placeholder; do not use 0,
	Currency      Currency   `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
	Memo          String     `xml:"MEMO,omitempty"`
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

type Position interface {
	PositionType() string
}

type DebtPosition struct {
	XMLName xml.Name    `xml:"POSDEBT"`
	InvPos  InvPosition `xml:"INVPOS"`
}

func (p DebtPosition) PositionType() string {
	return "POSDEBT"
}

type MFPosition struct {
	XMLName     xml.Name    `xml:"POSMF"`
	InvPos      InvPosition `xml:"INVPOS"`
	UnitsStreet Amount      `xml:"UNITSSTREET,omitempty"` // Units in the FI’s street name
	UnitsUser   Amount      `xml:"UNITSUSER,omitempty"`   // Units in the user's name directly
	ReinvDiv    Boolean     `xml:"REINVDIV,omitempty"`    // Reinvest dividends
	ReinCG      Boolean     `xml:"REINVCG,omitempty"`     // Reinvest capital gains
}

func (p MFPosition) PositionType() string {
	return "POSMF"
}

type OptPosition struct {
	XMLName xml.Name    `xml:"POSOPT"`
	InvPos  InvPosition `xml:"INVPOS"`
	Secured String      `xml:"SECURED,omitempty"` // One of NAKED, COVERED
}

func (p OptPosition) PositionType() string {
	return "POSOPT"
}

type OtherPosition struct {
	XMLName xml.Name    `xml:"POSOTHER"`
	InvPos  InvPosition `xml:"INVPOS"`
}

func (p OtherPosition) PositionType() string {
	return "POSOTHER"
}

type StockPosition struct {
	XMLName     xml.Name    `xml:"POSSTOCK"`
	InvPos      InvPosition `xml:"INVPOS"`
	UnitsStreet Amount      `xml:"UNITSSTREET,omitempty"` // Units in the FI’s street name
	UnitsUser   Amount      `xml:"UNITSUSER,omitempty"`   // Units in the user's name directly
	ReinvDiv    Boolean     `xml:"REINVDIV,omitempty"`    // Reinvest dividends
}

func (p StockPosition) PositionType() string {
	return "POSSTOCK"
}

type PositionList []Position

func (p PositionList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "POSDEBT":
				var position DebtPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				p = append(p, Position(position))
			case "POSMF":
				var position MFPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				p = append(p, Position(position))
			case "POSOPT":
				var position OptPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				p = append(p, Position(position))
			case "POSOTHER":
				var position OtherPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				p = append(p, Position(position))
			case "POSSTOCK":
				var position StockPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				p = append(p, Position(position))
			default:
				return errors.New("Invalid INVPOSLIST child tag: " + startElement.Name.Local)
			}
		} else {
			return errors.New("Didn't find an opening element")
		}
	}
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
	InvPosList  PositionList       `xml:"INVSTMTRS>INVPOSLIST,omitempty"`
	InvBal      InvBalance         `xml:"INVSTMTRS>INVBAL,omitempty"`
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
