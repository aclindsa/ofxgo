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

type CCStatementRequest struct {
	XMLName          xml.Name `xml:"CCSTMTTRNRQ"`
	TrnUID           UID      `xml:"TRNUID"`
	CCAcctFrom       CCAcct   `xml:"CCSTMTRQ>CCACCTFROM"`
	DtStart          Date     `xml:"CCSTMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd            Date     `xml:"CCSTMTRQ>INCTRAN>DTEND,omitempty"`
	Include          Boolean  `xml:"CCSTMTRQ>INCTRAN>INCLUDE"`            // Include transactions (instead of just balance)
	IncludePending   Boolean  `xml:"CCSTMTRQ>INCLUDEPENDING,omitempty"`   // Include pending transactions
	IncludeTranImage Boolean  `xml:"CCSTMTRQ>INCLUDETRANIMAGE,omitempty"` // Include transaction images
}

func (r *CCStatementRequest) Name() string {
	return "CCSTMTTRNRQ"
}

func (r *CCStatementRequest) Valid() (bool, error) {
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

type ImageData struct {
	XMLName      xml.Name `xml:"IMAGEDATA"`
	ImageType    String   `xml:"IMAGETYPE"`    // One of STATEMENT, TRANSACTION, TAX
	ImageRef     String   `xml:"IMAGEREF"`     // URL or identifier, depending on IMAGEREFTYPE
	ImageRefType String   `xml:"IMAGEREFTYPE"` // One of OPAQUE, URL, FORMURL (see spec for more details on how to access images of each of these types)
	// Only one of the next two should be valid at any given time
	ImageDelay   Int    `xml:"IMAGEDELAY,omitempty"`   // Number of calendar days from DTSERVER (for statement images) or DTPOSTED (for transaction image) the image will become available
	DtImageAvail Date   `xml:"DTIMAGEAVAIL,omitempty"` // Date image will become available
	ImageTTL     Int    `xml:"IMAGETTL,omitempty"`     // Number of days after image becomes available that it will remain available
	CheckSup     String `xml:"CHECKSUP,omitempty"`     // What is contained in check images. One of FRONTONLY, BACKONLY, FRONTANDBACK
}

type Transaction struct {
	XMLName       xml.Name `xml:"STMTTRN"`
	TrnType       String   `xml:"TRNTYPE"` // One of CREDIT, DEBIT, INT (interest earned or paid. Note: Depends on signage of amount), DIV, FEE, SRVCHG (service charge), DEP (deposit), ATM (Note: Depends on signage of amount), POS (Note: Depends on signage of amount), XFER, CHECK, PAYMENT, CASH, DIRECTDEP, DIRECTDEBIT, REPEATPMT, OTHER
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
	Name          String      `xml:"NAME,omitempty"`
	Payee         Payee       `xml:"PAYEE,omitempty"`
	ExtdName      String      `xml:"EXTDNAME,omitempty"`   // Extended name of payee or transaction description
	BankAcctTo    BankAcct    `xml:"BANKACCTTO,omitempty"` // If the transfer was to a bank account we have the account information for
	CCAcctTo      CCAcct      `xml:"CCACCTTO,omitempty"`   // If the transfer was to a credit card account we have the account information for
	Memo          String      `xml:"MEMO,omitempty"`       // Extra information (not in NAME)
	ImageData     []ImageData `xml:"IMAGEDATA,omitempty"`
	Currency      String      `xml:"CURRENCY,omitempty"`      // If different from CURDEF in STMTTRS
	OrigCurrency  String      `xml:"ORIGCURRENCY,omitempty"`  // If different from CURDEF in STMTTRS
	Inv401kSource String      `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST (Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST.)
}

type TransactionList struct {
	XMLName      xml.Name      `xml:"BANKTRANLIST"`
	DtStart      Date          `xml:"DTSTART"`
	DtEnd        Date          `xml:"DTEND"`
	Transactions []Transaction `xml:"STMTTRN,omitempty"`
}

type PendingTransaction struct {
	XMLName      xml.Name    `xml:"STMTTRN"`
	TrnType      String      `xml:"TRNTYPE"` // One of CREDIT, DEBIT, INT (interest earned or paid. Note: Depends on signage of amount), DIV, FEE, SRVCHG (service charge), DEP (deposit), ATM (Note: Depends on signage of amount), POS (Note: Depends on signage of amount), XFER, CHECK, PAYMENT, CASH, DIRECTDEP, DIRECTDEBIT, REPEATPMT, HOLD, OTHER
	DtTran       Date        `xml:"DTTRAN"`
	DtExpire     Date        `xml:"DTEXPIRE,omitempty"` // only valid for TrnType==HOLD, the date the hold will expire
	TrnAmt       Amount      `xml:"TRNAMT"`
	RefNum       String      `xml:"REFNUM,omitempty"`
	Name         String      `xml:"NAME,omitempty"`
	ExtdName     String      `xml:"EXTDNAME,omitempty"` // Extended name of payee or transaction description
	Memo         String      `xml:"MEMO,omitempty"`     // Extra information (not in NAME)
	ImageData    []ImageData `xml:"IMAGEDATA,omitempty"`
	Currency     String      `xml:"CURRENCY,omitempty"`     // If different from CURDEF in STMTTRS
	OrigCurrency String      `xml:"ORIGCURRENCY,omitempty"` // If different from CURDEF in STMTTRS
}

// List of pending transactions
type PendingTransactionList struct {
	XMLName      xml.Name             `xml:"BANKTRANLISTP"`
	DtAsOf       Date                 `xml:"DTASOF"`
	Transactions []PendingTransaction `xml:"STMTTRNP,omitempty"`
}

type Currency struct {
	XMLName xml.Name // CURRENCY or ORIGCURRENCY
	CurRate Amount   `xml:"CURRATE"` // Ratio of <CURDEF> currency to <CURSYM> currency
	CurSym  String   `xml:"CURSYM"`  // ISO-4217 3-character currency identifier
}

type Balance struct {
	XMLName xml.Name `xml:"BAL"`
	Name    String   `xml:"NAME"`
	Desc    String   `xml:"DESC"`

	// Balance type:
	// DOLLAR = dollar (value formatted DDDD.cc)
	// PERCENT = percentage (value formatted XXXX.YYYY)
	// NUMBER = number (value formatted as is)
	BalType String `xml:"BALTYPE"`

	Value    Amount   `xml:"VALUE"`
	DtAsOf   Date     `xml:"DTASOF,omitempty"`
	Currency Currency `xml:"CURRENCY,omitempty"` // if BALTYPE is DOLLAR
}

type StatementResponse struct {
	XMLName       xml.Name               `xml:"STMTTRNRS"`
	TrnUID        UID                    `xml:"TRNUID"`
	CurDef        String                 `xml:"STMTRS>CURDEF"`
	BankAcctFrom  BankAcct               `xml:"STMTRS>BANKACCTFROM"`
	BankTranList  TransactionList        `xml:"STMTRS>BANKTRANLIST,omitempty"`
	BankTranListP PendingTransactionList `xml:"STMTRS>BANKTRANLISTP,omitempty"`
	BalAmt        Amount                 `xml:"STMTRS>LEDGERBAL>BALAMT"`
	DtAsOf        Date                   `xml:"STMTRS>LEDGERBAL>DTASOF"`
	AvailBalAmt   Amount                 `xml:"STMTRS>AVAILBAL>BALAMT,omitempty"`
	AvailDtAsOf   Date                   `xml:"STMTRS>AVAILBAL>DTASOF,omitempty"`
	CashAdvBalAmt Amount                 `xml:"STMTRS>CASHADVBALAMT,omitempty"` // Only for CREDITLINE accounts, available balance for cash advances
	IntRate       Amount                 `xml:"STMTRS>INTRATE,omitempty"`       // Current interest rate
	BalList       []Balance              `xml:"STMTRS>BALLIST>BAL,omitempty"`
	MktgInfo      String                 `xml:"STMTRS>MKTGINFO"` // Marketing information
}

func (sr StatementResponse) Name() string {
	return "STMTTRNRS"
}

func (sr StatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

type CCStatementResponse struct {
	XMLName      xml.Name        `xml:"CCSTMTTRNRS"`
	TrnUID       UID             `xml:"TRNUID"`
	CurDef       String          `xml:"CCSTMTRS>CURDEF"`
	CCAcctFrom   CCAcct          `xml:"CCSTMTRS>CCACCTFROM"`
	BankTranList TransactionList `xml:"CCSTMTRS>BANKTRANLIST,omitempty"`
	//BANKTRANLISTP
	BalAmt        Amount    `xml:"CCSTMTRS>LEDGERBAL>BALAMT"`
	DtAsOf        Date      `xml:"CCSTMTRS>LEDGERBAL>DTASOF"`
	AvailBalAmt   Amount    `xml:"CCSTMTRS>AVAILBAL>BALAMT,omitempty"`
	AvailDtAsOf   Date      `xml:"CCSTMTRS>AVAILBAL>DTASOF,omitempty"`
	CashAdvBalAmt Amount    `xml:"CCSTMTRS>CASHADVBALAMT,omitempty"`           // Only for CREDITLINE accounts, available balance for cash advances
	IntRatePurch  Amount    `xml:"CCSTMTRS>INTRATEPURCH,omitempty"`            // Current interest rate for purchases
	IntRateCash   Amount    `xml:"CCSTMTRS>INTRATECASH,omitempty"`             // Current interest rate for cash advances
	IntRateXfer   Amount    `xml:"CCSTMTRS>INTRATEXFER,omitempty"`             // Current interest rate for cash advances
	RewardName    String    `xml:"CCSTMTRS>REWARDINFO>NAME,omitempty"`         // Name of the reward program referred to by the next two elements
	RewardBal     Amount    `xml:"CCSTMTRS>REWARDINFO>REWARDBAL,omitempty"`    // Current balance of the reward program
	RewardEarned  Amount    `xml:"CCSTMTRS>REWARDINFO>REWARDEARNED,omitempty"` // Reward amount earned YTD
	BalList       []Balance `xml:"CCSTMTRS>BALLIST>BAL,omitempty"`
	MktgInfo      String    `xml:"CCSTMTRS>MKTGINFO"` // Marketing information
}

func (sr CCStatementResponse) Name() string {
	return "CCSTMTTRNRS"
}

func (sr CCStatementResponse) Valid() (bool, error) {
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
			case "CCSTMTTRNRS":
				var info CCStatementResponse
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
