package ofxgo

import (
	"github.com/aclindsa/go/src/encoding/xml"
)

// StatementRequest represents a request for a bank statement. It is used to
// request balances and/or transactions for checking, savings, money market,
// and line of credit accounts. See CCStatementRequest for the analog for
// credit card accounts.
type StatementRequest struct {
	XMLName   xml.Name `xml:"STMTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	BankAcctFrom   BankAcct `xml:"STMTRQ>BANKACCTFROM"`
	DtStart        *Date    `xml:"STMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd          *Date    `xml:"STMTRQ>INCTRAN>DTEND,omitempty"`
	Include        Boolean  `xml:"STMTRQ>INCTRAN>INCLUDE"`          // Include transactions (instead of just balance)
	IncludePending Boolean  `xml:"STMTRQ>INCLUDEPENDING,omitempty"` // Include pending transactions
	IncTranImg     Boolean  `xml:"STMTRQ>INCTRANIMG,omitempty"`     // Include transaction images
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *StatementRequest) Name() string {
	return "STMTTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *StatementRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *StatementRequest) Type() messageType {
	return BankRq
}

// Payee specifies a complete billing address for a payee
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

// ImageData represents the metadata surrounding a check or other image file,
// including how to retrieve the image
type ImageData struct {
	XMLName      xml.Name     `xml:"IMAGEDATA"`
	ImageType    imageType    `xml:"IMAGETYPE"`    // One of STATEMENT, TRANSACTION, TAX
	ImageRef     String       `xml:"IMAGEREF"`     // URL or identifier, depending on IMAGEREFTYPE
	ImageRefType imageRefType `xml:"IMAGEREFTYPE"` // One of OPAQUE, URL, FORMURL (see spec for more details on how to access images of each of these types)
	// Only one of the next two should be valid at any given time
	ImageDelay   Int      `xml:"IMAGEDELAY,omitempty"`   // Number of calendar days from DTSERVER (for statement images) or DTPOSTED (for transaction image) the image will become available
	DtImageAvail *Date    `xml:"DTIMAGEAVAIL,omitempty"` // Date image will become available
	ImageTTL     Int      `xml:"IMAGETTL,omitempty"`     // Number of days after image becomes available that it will remain available
	CheckSup     checkSup `xml:"CHECKSUP,omitempty"`     // What is contained in check images. One of FRONTONLY, BACKONLY, FRONTANDBACK
}

// Transaction represents a single banking transaction. At a minimum, it
// identifies the type of transaction (TrnType) and the date it was posted
// (DtPosted). Ideally it also provides metadata to help the user recognize
// this transaction (i.e. CheckNum, Name or Payee, Memo, etc.)
type Transaction struct {
	XMLName       xml.Name      `xml:"STMTTRN"`
	TrnType       trnType       `xml:"TRNTYPE"` // One of CREDIT, DEBIT, INT (interest earned or paid. Note: Depends on signage of amount), DIV, FEE, SRVCHG (service charge), DEP (deposit), ATM (Note: Depends on signage of amount), POS (Note: Depends on signage of amount), XFER, CHECK, PAYMENT, CASH, DIRECTDEP, DIRECTDEBIT, REPEATPMT, OTHER
	DtPosted      Date          `xml:"DTPOSTED"`
	DtUser        *Date         `xml:"DTUSER,omitempty"`
	DtAvail       *Date         `xml:"DTAVAIL,omitempty"`
	TrnAmt        Amount        `xml:"TRNAMT"`
	FiTId         String        `xml:"FITID"`                   // Client uses FITID to detect whether it has previously downloaded the transaction
	CorrectFiTId  String        `xml:"CORRECTFITID,omitempty"`  // Transaction Id that this transaction corrects, if present
	CorrectAction correctAction `xml:"CORRECTACTION,omitempty"` // One of DELETE, REPLACE
	SrvrTId       String        `xml:"SRVRTID,omitempty"`
	CheckNum      String        `xml:"CHECKNUM,omitempty"`
	RefNum        String        `xml:"REFNUM,omitempty"`
	SIC           Int           `xml:"SIC,omitempty"` // Standard Industrial Code
	PayeeID       String        `xml:"PAYEEID,omitempty"`
	// Note: Servers should provide NAME or PAYEE, but not both
	Name          String        `xml:"NAME,omitempty"`
	Payee         *Payee        `xml:"PAYEE,omitempty"`
	ExtdName      String        `xml:"EXTDNAME,omitempty"`   // Extended name of payee or transaction description
	BankAcctTo    *BankAcct     `xml:"BANKACCTTO,omitempty"` // If the transfer was to a bank account we have the account information for
	CCAcctTo      *CCAcct       `xml:"CCACCTTO,omitempty"`   // If the transfer was to a credit card account we have the account information for
	Memo          String        `xml:"MEMO,omitempty"`       // Extra information (not in NAME)
	ImageData     []ImageData   `xml:"IMAGEDATA,omitempty"`
	Currency      String        `xml:"CURRENCY,omitempty"`      // If different from CURDEF in STMTTRS
	OrigCurrency  String        `xml:"ORIGCURRENCY,omitempty"`  // If different from CURDEF in STMTTRS
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST (Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST.)
}

// TransactionList represents a list of bank transactions, and also includes
// the date range its transactions cover.
type TransactionList struct {
	XMLName      xml.Name      `xml:"BANKTRANLIST"`
	DtStart      Date          `xml:"DTSTART"` // Start date for transaction data
	DtEnd        Date          `xml:"DTEND"`   // Value that client should send in next <DTSTART> request to ensure that it does not miss any transactions
	Transactions []Transaction `xml:"STMTTRN,omitempty"`
}

// PendingTransaction represents a single pending transaction. It is similar to
// Transaction, but is not finalized (and may never be). For instance, it lacks
// FiTId and DtPosted fields.
type PendingTransaction struct {
	XMLName      xml.Name    `xml:"STMTTRNP"`
	TrnType      trnType     `xml:"TRNTYPE"` // One of CREDIT, DEBIT, INT (interest earned or paid. Note: Depends on signage of amount), DIV, FEE, SRVCHG (service charge), DEP (deposit), ATM (Note: Depends on signage of amount), POS (Note: Depends on signage of amount), XFER, CHECK, PAYMENT, CASH, DIRECTDEP, DIRECTDEBIT, REPEATPMT, HOLD, OTHER
	DtTran       Date        `xml:"DTTRAN"`
	DtExpire     *Date       `xml:"DTEXPIRE,omitempty"` // only valid for TrnType==HOLD, the date the hold will expire
	TrnAmt       Amount      `xml:"TRNAMT"`
	RefNum       String      `xml:"REFNUM,omitempty"`
	Name         String      `xml:"NAME,omitempty"`
	ExtdName     String      `xml:"EXTDNAME,omitempty"` // Extended name of payee or transaction description
	Memo         String      `xml:"MEMO,omitempty"`     // Extra information (not in NAME)
	ImageData    []ImageData `xml:"IMAGEDATA,omitempty"`
	Currency     String      `xml:"CURRENCY,omitempty"`     // If different from CURDEF in STMTTRS
	OrigCurrency String      `xml:"ORIGCURRENCY,omitempty"` // If different from CURDEF in STMTTRS
}

// PendingTransactionList represents a list of pending transactions, along with
// the date they were generated
type PendingTransactionList struct {
	XMLName      xml.Name             `xml:"BANKTRANLISTP"`
	DtAsOf       Date                 `xml:"DTASOF"` // Date and time this set of pending transactions was generated
	Transactions []PendingTransaction `xml:"STMTTRNP,omitempty"`
}

// Balance represents a generic (free-form) balance defined by an FI.
type Balance struct {
	XMLName xml.Name `xml:"BAL"`
	Name    String   `xml:"NAME"`
	Desc    String   `xml:"DESC"`

	// Balance type:
	// DOLLAR = dollar (value formatted DDDD.cc)
	// PERCENT = percentage (value formatted XXXX.YYYY)
	// NUMBER = number (value formatted as is)
	BalType balType `xml:"BALTYPE"`

	Value    Amount    `xml:"VALUE"`
	DtAsOf   *Date     `xml:"DTASOF,omitempty"`
	Currency *Currency `xml:"CURRENCY,omitempty"` // if BALTYPE is DOLLAR
}

// StatementResponse represents a bank account statement, including its
// balances and possibly transactions. It is a response to StatementRequest, or
// sometimes provided as part of an OFX file downloaded manually from an FI.
type StatementResponse struct {
	XMLName   xml.Name `xml:"STMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	CurDef        String                  `xml:"STMTRS>CURDEF"`
	BankAcctFrom  BankAcct                `xml:"STMTRS>BANKACCTFROM"`
	BankTranList  *TransactionList        `xml:"STMTRS>BANKTRANLIST,omitempty"`
	BankTranListP *PendingTransactionList `xml:"STMTRS>BANKTRANLISTP,omitempty"`
	BalAmt        Amount                  `xml:"STMTRS>LEDGERBAL>BALAMT"`
	DtAsOf        Date                    `xml:"STMTRS>LEDGERBAL>DTASOF"`
	AvailBalAmt   *Amount                 `xml:"STMTRS>AVAILBAL>BALAMT,omitempty"`
	AvailDtAsOf   *Date                   `xml:"STMTRS>AVAILBAL>DTASOF,omitempty"`
	CashAdvBalAmt *Amount                 `xml:"STMTRS>CASHADVBALAMT,omitempty"` // Only for CREDITLINE accounts, available balance for cash advances
	IntRate       *Amount                 `xml:"STMTRS>INTRATE,omitempty"`       // Current interest rate
	BalList       []Balance               `xml:"STMTRS>BALLIST>BAL,omitempty"`
	MktgInfo      String                  `xml:"STMTRS>MKTGINFO,omitempty"` // Marketing information
}

// Name returns the name of the top-level transaction XML/SGML element
func (sr *StatementResponse) Name() string {
	return "STMTTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (sr *StatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (sr *StatementResponse) Type() messageType {
	return BankRs
}
