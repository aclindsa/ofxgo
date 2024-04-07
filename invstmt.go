package ofxgo

import (
	"errors"
	"github.com/aclindsa/xml"
)

// InvStatementRequest allows a customer to request transactions, positions,
// open orders, and balances. It specifies what types of information to include
// in hte InvStatementResponse and which account to include it for.
type InvStatementRequest struct {
	XMLName   xml.Name `xml:"INVSTMTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	InvAcctFrom      InvAcct `xml:"INVSTMTRQ>INVACCTFROM"`
	DtStart          *Date   `xml:"INVSTMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd            *Date   `xml:"INVSTMTRQ>INCTRAN>DTEND,omitempty"`
	Include          Boolean `xml:"INVSTMTRQ>INCTRAN>INCLUDE"`         // Include transactions (instead of just balance)
	IncludeOO        Boolean `xml:"INVSTMTRQ>INCOO"`                   // Include open orders
	PosDtAsOf        *Date   `xml:"INVSTMTRQ>INCPOS>DTASOF,omitempty"` // Date that positions should be sent down for, if present
	IncludePos       Boolean `xml:"INVSTMTRQ>INCPOS>INCLUDE"`          // Include position data in response
	IncludeBalance   Boolean `xml:"INVSTMTRQ>INCBAL"`                  // Include investment balance in response
	Include401K      Boolean `xml:"INVSTMTRQ>INC401K,omitempty"`       // Include 401k information
	Include401KBal   Boolean `xml:"INVSTMTRQ>INC401KBAL,omitempty"`    // Include 401k balance information
	IncludeTranImage Boolean `xml:"INVSTMTRQ>INCTRANIMAGE,omitempty"`  // Include transaction images
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *InvStatementRequest) Name() string {
	return "INVSTMTTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *InvStatementRequest) Valid(version ofxVersion) (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *InvStatementRequest) Type() messageType {
	return InvStmtRq
}

// InvTran represents generic investment transaction. It is included in both
// InvBuy and InvSell as well as many of the more specific transaction
// aggregates.
type InvTran struct {
	XMLName       xml.Name `xml:"INVTRAN"`
	FiTID         String   `xml:"FITID"`                   // Unique FI-assigned transaction ID. This ID is used to detect duplicate downloads
	SrvrTID       String   `xml:"SRVRTID,omitempty"`       // Server assigned transaction ID
	DtTrade       Date     `xml:"DTTRADE"`                 // trade date; for stock splits, day of record
	DtSettle      *Date    `xml:"DTSETTLE,omitempty"`      // settlement date; for stock splits, execution date
	ReversalFiTID String   `xml:"REVERSALFITID,omitempty"` // For a reversal transaction, the FITID of the transaction that is being reversed.
	Memo          String   `xml:"MEMO,omitempty"`
}

// InvBuy represents generic investment purchase transaction. It is included
// in many of the more specific transaction Buy* aggregates below.
type InvBuy struct {
	XMLName      xml.Name    `xml:"INVBUY"`
	InvTran      InvTran     `xml:"INVTRAN"`
	SecID        SecurityID  `xml:"SECID"`
	Units        Amount      `xml:"UNITS"`            // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice    Amount      `xml:"UNITPRICE"`        // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Markup       Amount      `xml:"MARKUP,omitempty"` // Portion of UNITPRICE that is attributed to the dealer markup
	Commission   Amount      `xml:"COMMISSION,omitempty"`
	Taxes        Amount      `xml:"TAXES,omitempty"`
	Fees         Amount      `xml:"FEES,omitempty"`
	Load         Amount      `xml:"LOAD,omitempty"`
	Total        Amount      `xml:"TOTAL"`                  // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	Currency     Currency    `xml:"CURRENCY,omitempty"`     // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency Currency    `xml:"ORIGCURRENCY,omitempty"` // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	SubAcctSec   subAcctType `xml:"SUBACCTSEC"`             // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund  subAcctType `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER

	// The next three elements must either all be provided, or none of them
	LoanID        String `xml:"LOANID,omitempty"`        // For 401(k) accounts only. Indicates that the transaction was due to a loan or a loan repayment, and which loan it was
	LoanPrincipal Amount `xml:"LOANPRINCIPAL,omitempty"` // For 401(k) accounts only. Indicates how much of the loan repayment was principal
	LoanInterest  Amount `xml:"LOANINTEREST,omitempty"`  // For 401(k) accounts only. Indicates how much of the loan repayment was interest

	Inv401kSource    inv401kSource `xml:"INV401KSOURCE,omitempty"`    // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
	DtPayroll        *Date         `xml:"DTPAYROLL,omitempty"`        // For 401(k)accounts, date the funds for this transaction was obtained via payroll deduction
	PriorYearContrib Boolean       `xml:"PRIORYEARCONTRIB,omitempty"` // For 401(k) accounts, indicates that this Buy was made with a prior year contribution
}

// InvSell represents generic investment sale transaction. It is included in
// many of the more specific transaction Sell* aggregates below.
type InvSell struct {
	XMLName      xml.Name    `xml:"INVSELL"`
	InvTran      InvTran     `xml:"INVTRAN"`
	SecID        SecurityID  `xml:"SECID"`
	Units        Amount      `xml:"UNITS"`              // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice    Amount      `xml:"UNITPRICE"`          // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Markdown     Amount      `xml:"MARKDOWN,omitempty"` // Portion of UNITPRICE that is attributed to the dealer markdown
	Commission   Amount      `xml:"COMMISSION,omitempty"`
	Taxes        Amount      `xml:"TAXES,omitempty"`
	Fees         Amount      `xml:"FEES,omitempty"`
	Load         Amount      `xml:"LOAD,omitempty"`
	Withholding  Amount      `xml:"WITHHOLDING,omitempty"`  // Federal tax withholdings
	TaxExempt    Boolean     `xml:"TAXEXEMPT,omitempty"`    // Tax-exempt transaction
	Total        Amount      `xml:"TOTAL"`                  // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	Gain         Amount      `xml:"GAIN,omitempty"`         // Total gain
	Currency     Currency    `xml:"CURRENCY,omitempty"`     // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency Currency    `xml:"ORIGCURRENCY,omitempty"` // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	SubAcctSec   subAcctType `xml:"SUBACCTSEC"`             // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund  subAcctType `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER

	LoanID           String `xml:"LOANID,omitempty"`           // For 401(k) accounts only. Indicates that the transaction was due to a loan or a loan repayment, and which loan it was
	StateWithholding Amount `xml:"STATEWITHHOLDING,omitempty"` // State tax withholdings
	Penalty          Amount `xml:"PENALTY,omitempty"`          // Amount withheld due to penalty

	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// BuyDebt represents a transaction purchasing a debt security
type BuyDebt struct {
	XMLName  xml.Name `xml:"BUYDEBT"`
	InvBuy   InvBuy   `xml:"INVBUY"`
	AccrdInt Amount   `xml:"ACCRDINT,omitempty"` // Accrued interest. This amount is not reflected in the <TOTAL> field of a containing aggregate.
}

// TransactionType returns a string representation of this transaction's type
func (t BuyDebt) TransactionType() string {
	return "BUYDEBT"
}

func (t BuyDebt) InvTransaction() InvTran {
	return t.InvBuy.InvTran
}

// BuyMF represents a transaction purchasing a mutual fund
type BuyMF struct {
	XMLName  xml.Name `xml:"BUYMF"`
	InvBuy   InvBuy   `xml:"INVBUY"`
	BuyType  buyType  `xml:"BUYTYPE"`            // One of BUY, BUYTOCOVER (BUYTOCOVER used to close short sales.)
	RelFiTID String   `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
}

// TransactionType returns a string representation of this transaction's type
func (t BuyMF) TransactionType() string {
	return "BUYMF"
}

func (t BuyMF) InvTransaction() InvTran {
	return t.InvBuy.InvTran
}

// BuyOpt represents a transaction purchasing an option
type BuyOpt struct {
	XMLName    xml.Name   `xml:"BUYOPT"`
	InvBuy     InvBuy     `xml:"INVBUY"`
	OptBuyType optBuyType `xml:"OPTBUYTYPE"` // type of purchase: BUYTOOPEN, BUYTOCLOSE (The BUYTOOPEN buy type is like “ordinary” buying of option and works like stocks.)
	ShPerCtrct Int        `xml:"SHPERCTRCT"` // Shares per contract
}

// TransactionType returns a string representation of this transaction's type
func (t BuyOpt) TransactionType() string {
	return "BUYOPT"
}

func (t BuyOpt) InvTransaction() InvTran {
	return t.InvBuy.InvTran
}

// BuyOther represents a transaction purchasing a type of security not covered
// by the other Buy* structs
type BuyOther struct {
	XMLName xml.Name `xml:"BUYOTHER"`
	InvBuy  InvBuy   `xml:"INVBUY"`
}

// TransactionType returns a string representation of this transaction's type
func (t BuyOther) TransactionType() string {
	return "BUYOTHER"
}

func (t BuyOther) InvTransaction() InvTran {
	return t.InvBuy.InvTran
}

// BuyStock represents a transaction purchasing stock
type BuyStock struct {
	XMLName xml.Name `xml:"BUYSTOCK"`
	InvBuy  InvBuy   `xml:"INVBUY"`
	BuyType buyType  `xml:"BUYTYPE"` // One of BUY, BUYTOCOVER (BUYTOCOVER used to close short sales.)
}

// TransactionType returns a string representation of this transaction's type
func (t BuyStock) TransactionType() string {
	return "BUYSTOCK"
}

func (t BuyStock) InvTransaction() InvTran {
	return t.InvBuy.InvTran
}

// ClosureOpt represents a transaction closing a position for an option
type ClosureOpt struct {
	XMLName    xml.Name    `xml:"CLOSUREOPT"`
	InvTran    InvTran     `xml:"INVTRAN"`
	SecID      SecurityID  `xml:"SECID"`
	OptAction  optAction   `xml:"OPTACTION"`          // One of EXERCISE, ASSIGN, EXPIRE. The EXERCISE action is used to close out an option that is exercised. The ASSIGN action is used when an option writer is assigned. The EXPIRE action is used when the option’s expired date is reached
	Units      Amount      `xml:"UNITS"`              // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	ShPerCtrct Int         `xml:"SHPERCTRCT"`         // Shares per contract
	SubAcctSec subAcctType `xml:"SUBACCTSEC"`         // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	RelFiTID   String      `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
	Gain       Amount      `xml:"GAIN,omitempty"`     // Total gain
}

// TransactionType returns a string representation of this transaction's type
func (t ClosureOpt) TransactionType() string {
	return "CLOSUREOPT"
}

func (t ClosureOpt) InvTransaction() InvTran {
	return t.InvTran
}

// Income represents a transaction where investment income is being realized as
// cash into the investment account
type Income struct {
	XMLName       xml.Name      `xml:"INCOME"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	IncomeType    incomeType    `xml:"INCOMETYPE"` // Type of investment income: CGLONG (capital gains-long term), CGSHORT (capital gains-short term), DIV (dividend), INTEREST, MISC
	Total         Amount        `xml:"TOTAL"`
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   subAcctType   `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	TaxExempt     Boolean       `xml:"TAXEXEMPT,omitempty"`     // Tax-exempt transaction
	Withholding   Amount        `xml:"WITHHOLDING,omitempty"`   // Federal tax withholdings
	Currency      Currency      `xml:"CURRENCY,omitempty"`      // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency  Currency      `xml:"ORIGCURRENCY,omitempty"`  // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t Income) TransactionType() string {
	return "INCOME"
}

func (t Income) InvTransaction() InvTran {
	return t.InvTran
}

// InvExpense represents a transaction realizing an expense associated with an
// investment
type InvExpense struct {
	XMLName       xml.Name      `xml:"INVEXPENSE"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	Total         Amount        `xml:"TOTAL"`
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   subAcctType   `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency      Currency      `xml:"CURRENCY,omitempty"`      // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency  Currency      `xml:"ORIGCURRENCY,omitempty"`  // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t InvExpense) TransactionType() string {
	return "INVEXPENSE"
}

func (t InvExpense) InvTransaction() InvTran {
	return t.InvTran
}

// JrnlFund represents a transaction journaling cash holdings between
// sub-accounts within the same investment account
type JrnlFund struct {
	XMLName     xml.Name    `xml:"JRNLFUND"`
	InvTran     InvTran     `xml:"INVTRAN"`
	Total       Amount      `xml:"TOTAL"`
	SubAcctFrom subAcctType `xml:"SUBACCTFROM"` // Sub-account cash is being transferred from: CASH, MARGIN, SHORT, OTHER
	SubAcctTo   subAcctType `xml:"SUBACCTTO"`   // Sub-account cash is being transferred to: CASH, MARGIN, SHORT, OTHER
}

// TransactionType returns a string representation of this transaction's type
func (t JrnlFund) TransactionType() string {
	return "JRNLFUND"
}

func (t JrnlFund) InvTransaction() InvTran {
	return t.InvTran
}

// JrnlSec represents a transaction journaling security holdings between
// sub-accounts within the same investment account
type JrnlSec struct {
	XMLName     xml.Name    `xml:"JRNLSEC"`
	InvTran     InvTran     `xml:"INVTRAN"`
	SecID       SecurityID  `xml:"SECID"`
	SubAcctFrom subAcctType `xml:"SUBACCTFROM"` // Sub-account cash is being transferred from: CASH, MARGIN, SHORT, OTHER
	SubAcctTo   subAcctType `xml:"SUBACCTTO"`   // Sub-account cash is being transferred to: CASH, MARGIN, SHORT, OTHER
	Units       Amount      `xml:"UNITS"`       // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
}

// TransactionType returns a string representation of this transaction's type
func (t JrnlSec) TransactionType() string {
	return "JRNLSEC"
}

func (t JrnlSec) InvTransaction() InvTran {
	return t.InvTran
}

// MarginInterest represents a transaction realizing a margin interest expense
type MarginInterest struct {
	XMLName      xml.Name    `xml:"MARGININTEREST"`
	InvTran      InvTran     `xml:"INVTRAN"`
	Total        Amount      `xml:"TOTAL"`
	SubAcctFund  subAcctType `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency     Currency    `xml:"CURRENCY,omitempty"`     // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency Currency    `xml:"ORIGCURRENCY,omitempty"` // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
}

// TransactionType returns a string representation of this transaction's type
func (t MarginInterest) TransactionType() string {
	return "MARGININTEREST"
}

func (t MarginInterest) InvTransaction() InvTran {
	return t.InvTran
}

// Reinvest is a single transaction that contains both income and an investment
// transaction. If servers can’t track this as a single transaction they should
// return an Income transaction and an InvTran.
type Reinvest struct {
	XMLName       xml.Name      `xml:"REINVEST"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	IncomeType    incomeType    `xml:"INCOMETYPE"` // Type of investment income: CGLONG (capital gains-long term), CGSHORT (capital gains-short term), DIV (dividend), INTEREST, MISC
	Total         Amount        `xml:"TOTAL"`      // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"` // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	Units         Amount        `xml:"UNITS"`      // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice     Amount        `xml:"UNITPRICE"`  // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Commission    Amount        `xml:"COMMISSION,omitempty"`
	Taxes         Amount        `xml:"TAXES,omitempty"`
	Fees          Amount        `xml:"FEES,omitempty"`
	Load          Amount        `xml:"LOAD,omitempty"`
	TaxExempt     Boolean       `xml:"TAXEXEMPT,omitempty"`     // Tax-exempt transaction
	Currency      Currency      `xml:"CURRENCY,omitempty"`      // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency  Currency      `xml:"ORIGCURRENCY,omitempty"`  // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t Reinvest) TransactionType() string {
	return "REINVEST"
}

func (t Reinvest) InvTransaction() InvTran {
	return t.InvTran
}

// RetOfCap represents a transaction where capital is being returned to the
// account holder
type RetOfCap struct {
	XMLName       xml.Name      `xml:"RETOFCAP"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	Total         Amount        `xml:"TOTAL"`
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   subAcctType   `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency      Currency      `xml:"CURRENCY,omitempty"`      // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency  Currency      `xml:"ORIGCURRENCY,omitempty"`  // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t RetOfCap) TransactionType() string {
	return "RETOFCAP"
}

func (t RetOfCap) InvTransaction() InvTran {
	return t.InvTran
}

// SellDebt represents the sale of a debt security. Used when debt is sold,
// called, or reaches maturity.
type SellDebt struct {
	XMLName    xml.Name   `xml:"SELLDEBT"`
	InvSell    InvSell    `xml:"INVSELL"`
	SellReason sellReason `xml:"SELLREASON"`         // CALL (the debt was called), SELL (the debt was sold), MATURITY (the debt reached maturity)
	AccrdInt   Amount     `xml:"ACCRDINT,omitempty"` // Accrued interest
}

// TransactionType returns a string representation of this transaction's type
func (t SellDebt) TransactionType() string {
	return "SELLDEBT"
}

func (t SellDebt) InvTransaction() InvTran {
	return t.InvSell.InvTran
}

// SellMF represents a transaction selling a mutual fund
type SellMF struct {
	XMLName      xml.Name `xml:"SELLMF"`
	InvSell      InvSell  `xml:"INVSELL"`
	SellType     sellType `xml:"SELLTYPE"` // Type of sell. SELL, SELLSHORT
	AvgCostBasis Amount   `xml:"AVGCOSTBASIS"`
	RelFiTID     String   `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
}

// TransactionType returns a string representation of this transaction's type
func (t SellMF) TransactionType() string {
	return "SELLMF"
}

func (t SellMF) InvTransaction() InvTran {
	return t.InvSell.InvTran
}

// SellOpt represents a transaction selling an option. Depending on the value
// of OptSellType, can be used to sell a previously bought option or write a
// new option.
type SellOpt struct {
	XMLName     xml.Name    `xml:"SELLOPT"`
	InvSell     InvSell     `xml:"INVSELL"`
	OptSellType optSellType `xml:"OPTSELLTYPE"`        // For options, type of sell: SELLTOCLOSE, SELLTOOPEN. The SELLTOCLOSE action is selling a previously bought option. The SELLTOOPEN action is writing an option
	ShPerCtrct  Int         `xml:"SHPERCTRCT"`         // Shares per contract
	RelFiTID    String      `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
	RelType     relType     `xml:"RELTYPE,omitempty"`  // Related option transaction type: SPREAD, STRADDLE, NONE, OTHER
	Secured     secured     `xml:"SECURED,omitempty"`  // NAKED, COVERED
}

// TransactionType returns a string representation of this transaction's type
func (t SellOpt) TransactionType() string {
	return "SELLOPT"
}

func (t SellOpt) InvTransaction() InvTran {
	return t.InvSell.InvTran
}

// SellOther represents a transaction selling a security type not covered by
// the other Sell* structs
type SellOther struct {
	XMLName xml.Name `xml:"SELLOTHER"`
	InvSell InvSell  `xml:"INVSELL"`
}

// TransactionType returns a string representation of this transaction's type
func (t SellOther) TransactionType() string {
	return "SELLOTHER"
}

func (t SellOther) InvTransaction() InvTran {
	return t.InvSell.InvTran
}

// SellStock represents a transaction selling stock
type SellStock struct {
	XMLName  xml.Name `xml:"SELLSTOCK"`
	InvSell  InvSell  `xml:"INVSELL"`
	SellType sellType `xml:"SELLTYPE"` // Type of sell. SELL, SELLSHORT
}

// TransactionType returns a string representation of this transaction's type
func (t SellStock) TransactionType() string {
	return "SELLSTOCK"
}

func (t SellStock) InvTransaction() InvTran {
	return t.InvSell.InvTran
}

// Split represents a stock or mutual fund split
type Split struct {
	XMLName       xml.Name      `xml:"SPLIT"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	OldUnits      Amount        `xml:"OLDUNITS"`                // number of shares before the split
	NewUnits      Amount        `xml:"NEWUNITS"`                // number of shares after the split
	Numerator     Int           `xml:"NUMERATOR"`               // split ratio numerator
	Denominator   Int           `xml:"DENOMINATOR"`             // split ratio denominator
	Currency      Currency      `xml:"CURRENCY,omitempty"`      // Represents the currency this transaction is in (instead of CURDEF in INVSTMTRS) if Valid()
	OrigCurrency  Currency      `xml:"ORIGCURRENCY,omitempty"`  // Represents the currency this transaction was converted to INVSTMTRS' CURDEF from if Valid
	FracCash      Amount        `xml:"FRACCASH,omitempty"`      // cash for fractional units
	SubAcctFund   subAcctType   `xml:"SUBACCTFUND,omitempty"`   // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t Split) TransactionType() string {
	return "SPLIT"
}

func (t Split) InvTransaction() InvTran {
	return t.InvTran
}

// Transfer represents the transfer of securities into or out of an account
type Transfer struct {
	XMLName       xml.Name      `xml:"TRANSFER"`
	InvTran       InvTran       `xml:"INVTRAN"`
	SecID         SecurityID    `xml:"SECID"`
	SubAcctSec    subAcctType   `xml:"SUBACCTSEC"` // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	Units         Amount        `xml:"UNITS"`      // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	TferAction    tferAction    `xml:"TFERACTION"` // One of IN, OUT
	PosType       posType       `xml:"POSTYPE"`    // Position type. One of LONG, SHORT
	InvAcctFrom   InvAcct       `xml:"INVACCTFROM,omitempty"`
	AvgCostBasis  Amount        `xml:"AVGCOSTBASIS,omitempty"`
	UnitPrice     Amount        `xml:"UNITPRICE,omitempty"` // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	DtPurchase    *Date         `xml:"DTPURCHASE,omitempty"`
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// TransactionType returns a string representation of this transaction's type
func (t Transfer) TransactionType() string {
	return "TRANSFER"
}

func (t Transfer) InvTransaction() InvTran {
	return t.InvTran
}

// InvTransaction is a generic interface met by all investment transactions
// (Buy*, Sell*, & co.)
type InvTransaction interface {
	TransactionType() string
	InvTransaction() InvTran
}

// InvBankTransaction is a banking transaction performed in an investment
// account. This represents all transactions not related to securities - for
// instance, funding the account using cash from another bank.
type InvBankTransaction struct {
	XMLName      xml.Name      `xml:"INVBANKTRAN"`
	Transactions []Transaction `xml:"STMTTRN,omitempty"`
	SubAcctFund  subAcctType   `xml:"SUBACCTFUND"` // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
}

// InvTranList represents a list of investment account transactions. It
// includes the date range its transactions cover, as well as the bank- and
// security-related transactions themselves. It must be unmarshalled manually
// due to the structure (don't know what kind of InvTransaction is coming next)
type InvTranList struct {
	XMLName          xml.Name `xml:"INVTRANLIST"`
	DtStart          Date
	DtEnd            Date // This is the value that should be sent as <DTSTART> in the next InvStatementRequest to ensure that no transactions are missed
	InvTransactions  []InvTransaction
	BankTransactions []InvBankTransaction
}

// UnmarshalXML handles unmarshalling an InvTranList element from an SGML/XML
// string
func (l *InvTranList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "DTSTART":
				var dtstart Date
				if err := d.DecodeElement(&dtstart, &startElement); err != nil {
					return err
				}
				l.DtStart = dtstart
			case "DTEND":
				var dtend Date
				if err := d.DecodeElement(&dtend, &startElement); err != nil {
					return err
				}
				l.DtEnd = dtend
			case "BUYDEBT":
				var tran BuyDebt
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "BUYMF":
				var tran BuyMF
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "BUYOPT":
				var tran BuyOpt
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "BUYOTHER":
				var tran BuyOther
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "BUYSTOCK":
				var tran BuyStock
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "CLOSUREOPT":
				var tran ClosureOpt
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "INCOME":
				var tran Income
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "INVEXPENSE":
				var tran InvExpense
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "JRNLFUND":
				var tran JrnlFund
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "JRNLSEC":
				var tran JrnlSec
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "MARGININTEREST":
				var tran MarginInterest
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "REINVEST":
				var tran Reinvest
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "RETOFCAP":
				var tran RetOfCap
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SELLDEBT":
				var tran SellDebt
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SELLMF":
				var tran SellMF
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SELLOPT":
				var tran SellOpt
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SELLOTHER":
				var tran SellOther
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SELLSTOCK":
				var tran SellStock
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "SPLIT":
				var tran Split
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "TRANSFER":
				var tran Transfer
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.InvTransactions = append(l.InvTransactions, InvTransaction(tran))
			case "INVBANKTRAN":
				var tran InvBankTransaction
				if err := d.DecodeElement(&tran, &startElement); err != nil {
					return err
				}
				l.BankTransactions = append(l.BankTransactions, tran)
			default:
				return errors.New("Invalid INVTRANLIST child tag: " + startElement.Name.Local)
			}
		} else {
			return errors.New("Didn't find an opening element")
		}
	}
}

// MarshalXML handles marshalling an InvTranList element to an SGML/XML string
func (l *InvTranList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	invTranListElement := xml.StartElement{Name: xml.Name{Local: "INVTRANLIST"}}
	if err := e.EncodeToken(invTranListElement); err != nil {
		return err
	}
	err := e.EncodeElement(&l.DtStart, xml.StartElement{Name: xml.Name{Local: "DTSTART"}})
	if err != nil {
		return err
	}
	err = e.EncodeElement(&l.DtEnd, xml.StartElement{Name: xml.Name{Local: "DTEND"}})
	if err != nil {
		return err
	}
	for _, t := range l.InvTransactions {
		start := xml.StartElement{Name: xml.Name{Local: t.TransactionType()}}
		switch tran := t.(type) {
		case BuyDebt:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case BuyMF:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case BuyOpt:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case BuyOther:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case BuyStock:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case ClosureOpt:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case Income:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case InvExpense:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case JrnlFund:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case JrnlSec:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case MarginInterest:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case Reinvest:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case RetOfCap:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case SellDebt:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case SellMF:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case SellOpt:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case SellOther:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case SellStock:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case Split:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		case Transfer:
			if err := e.EncodeElement(&tran, start); err != nil {
				return err
			}
		default:
			return errors.New("Invalid INVTRANLIST child type: " + tran.TransactionType())
		}
	}
	for _, tran := range l.BankTransactions {
		err = e.EncodeElement(&tran, xml.StartElement{Name: xml.Name{Local: "INVBANKTRAN"}})
		if err != nil {
			return err
		}
	}
	if err := e.EncodeToken(invTranListElement.End()); err != nil {
		return err
	}
	return nil
}

// InvPosition contains generic position information included in each of the
// other *Position types
type InvPosition struct {
	XMLName       xml.Name      `xml:"INVPOS"`
	SecID         SecurityID    `xml:"SECID"`
	HeldInAcct    subAcctType   `xml:"HELDINACCT"`             // Sub-account type, one of CASH, MARGIN, SHORT, OTHER
	PosType       posType       `xml:"POSTYPE"`                // SHORT = Writer for options, Short for all others; LONG = Holder for options, Long for all others.
	Units         Amount        `xml:"UNITS"`                  // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice     Amount        `xml:"UNITPRICE"`              // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	MktVal        Amount        `xml:"MKTVAL"`                 // Market value of this position
	AvgCostBasis  Amount        `xml:"AVGCOSTBASIS,omitempty"` //
	DtPriceAsOf   Date          `xml:"DTPRICEASOF"`            // Date and time of unit price and market value, and cost basis. If this date is unknown, use 19900101 as the placeholder; do not use 0,
	Currency      *Currency     `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
	Memo          String        `xml:"MEMO,omitempty"`
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// Position is an interface satisfied by all the other *Position types
type Position interface {
	PositionType() string
	InvPosition() InvPosition
}

// DebtPosition represents a position held in a debt security
type DebtPosition struct {
	XMLName xml.Name    `xml:"POSDEBT"`
	InvPos  InvPosition `xml:"INVPOS"`
}

// PositionType returns a string representation of this position's type
func (p DebtPosition) PositionType() string {
	return "POSDEBT"
}

// InvPosition returns InvPos
func (p DebtPosition) InvPosition() InvPosition {
	return p.InvPos
}

// MFPosition represents a position held in a mutual fund
type MFPosition struct {
	XMLName     xml.Name    `xml:"POSMF"`
	InvPos      InvPosition `xml:"INVPOS"`
	UnitsStreet Amount      `xml:"UNITSSTREET,omitempty"` // Units in the FI’s street name
	UnitsUser   Amount      `xml:"UNITSUSER,omitempty"`   // Units in the user's name directly
	ReinvDiv    Boolean     `xml:"REINVDIV,omitempty"`    // Reinvest dividends
	ReinvCG     Boolean     `xml:"REINVCG,omitempty"`     // Reinvest capital gains
}

// PositionType returns a string representation of this position's type
func (p MFPosition) PositionType() string {
	return "POSMF"
}

// InvPosition returns InvPos
func (p MFPosition) InvPosition() InvPosition {
	return p.InvPos
}

// OptPosition represents a position held in an option
type OptPosition struct {
	XMLName xml.Name    `xml:"POSOPT"`
	InvPos  InvPosition `xml:"INVPOS"`
	Secured secured     `xml:"SECURED,omitempty"` // One of NAKED, COVERED
}

// PositionType returns a string representation of this position's type
func (p OptPosition) PositionType() string {
	return "POSOPT"
}

// InvPosition returns InvPos
func (p OptPosition) InvPosition() InvPosition {
	return p.InvPos
}

// OtherPosition represents a position held in a security type not covered by
// the other *Position elements
type OtherPosition struct {
	XMLName xml.Name    `xml:"POSOTHER"`
	InvPos  InvPosition `xml:"INVPOS"`
}

// PositionType returns a string representation of this position's type
func (p OtherPosition) PositionType() string {
	return "POSOTHER"
}

// InvPosition returns InvPos
func (p OtherPosition) InvPosition() InvPosition {
	return p.InvPos
}

// StockPosition represents a position held in a stock
type StockPosition struct {
	XMLName     xml.Name    `xml:"POSSTOCK"`
	InvPos      InvPosition `xml:"INVPOS"`
	UnitsStreet Amount      `xml:"UNITSSTREET,omitempty"` // Units in the FI’s street name
	UnitsUser   Amount      `xml:"UNITSUSER,omitempty"`   // Units in the user's name directly
	ReinvDiv    Boolean     `xml:"REINVDIV,omitempty"`    // Reinvest dividends
}

// PositionType returns a string representation of this position's type
func (p StockPosition) PositionType() string {
	return "POSSTOCK"
}

// InvPosition returns InvPos
func (p StockPosition) InvPosition() InvPosition {
	return p.InvPos
}

// PositionList represents a list of positions held in securities in an
// investment account
type PositionList []Position

// UnmarshalXML handles unmarshalling a PositionList from an XML string
func (p *PositionList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				*p = append(*p, Position(position))
			case "POSMF":
				var position MFPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				*p = append(*p, Position(position))
			case "POSOPT":
				var position OptPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				*p = append(*p, Position(position))
			case "POSOTHER":
				var position OtherPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				*p = append(*p, Position(position))
			case "POSSTOCK":
				var position StockPosition
				if err := d.DecodeElement(&position, &startElement); err != nil {
					return err
				}
				*p = append(*p, Position(position))
			default:
				return errors.New("Invalid INVPOSLIST child tag: " + startElement.Name.Local)
			}
		} else {
			return errors.New("Didn't find an opening element")
		}
	}
}

// MarshalXML handles marshalling a PositionList to an XML string
func (p PositionList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	invPosListElement := xml.StartElement{Name: xml.Name{Local: "INVPOSLIST"}}
	if err := e.EncodeToken(invPosListElement); err != nil {
		return err
	}
	for _, position := range p {
		start := xml.StartElement{Name: xml.Name{Local: position.PositionType()}}
		switch pos := position.(type) {
		case DebtPosition:
			if err := e.EncodeElement(&pos, start); err != nil {
				return err
			}
		case MFPosition:
			if err := e.EncodeElement(&pos, start); err != nil {
				return err
			}
		case OptPosition:
			if err := e.EncodeElement(&pos, start); err != nil {
				return err
			}
		case OtherPosition:
			if err := e.EncodeElement(&pos, start); err != nil {
				return err
			}
		case StockPosition:
			if err := e.EncodeElement(&pos, start); err != nil {
				return err
			}
		default:
			return errors.New("Invalid INVPOSLIST child type: " + pos.PositionType())
		}
	}
	if err := e.EncodeToken(invPosListElement.End()); err != nil {
		return err
	}
	return nil
}

// InvBalance contains three (or optionally four) specified balances as well as
// a free-form list of generic balance information which may be provided by an
// FI.
type InvBalance struct {
	XMLName       xml.Name  `xml:"INVBAL"`
	AvailCash     Amount    `xml:"AVAILCASH"`     // Available cash across all sub-accounts, including sweep funds
	MarginBalance Amount    `xml:"MARGINBALANCE"` // Negative means customer has borrowed funds
	ShortBalance  Amount    `xml:"SHORTBALANCE"`  // Always positive, market value of all short positions
	BuyPower      Amount    `xml:"BUYPOWER, omitempty"`
	BalList       []Balance `xml:"BALLIST>BAL,omitempty"`
}

// OO represents a generic open investment order. It is included in the other
// OO* elements.
type OO struct {
	XMLName       xml.Name      `xml:"OO"`
	FiTID         String        `xml:"FITID"`
	SrvrTID       String        `xml:"SRVRTID,omitempty"`
	SecID         SecurityID    `xml:"SECID"`
	DtPlaced      Date          `xml:"DTPLACED"`           // Date the order was placed
	Units         Amount        `xml:"UNITS"`              // Quantity of the security the open order is for
	SubAcct       subAcctType   `xml:"SUBACCT"`            // One of CASH, MARGIN, SHORT, OTHER
	Duration      duration      `xml:"DURATION"`           // How long the order is good for. One of DAY, GOODTILCANCEL, IMMEDIATE
	Restriction   restriction   `xml:"RESTRICTION"`        // Special restriction on the order: One of ALLORNONE, MINUNITS, NONE
	MinUnits      Amount        `xml:"MINUNITS,omitempty"` // Minimum number of units that must be filled for the order
	LimitPrice    Amount        `xml:"LIMITPRICE,omitempty"`
	StopPrice     Amount        `xml:"STOPPRICE,omitempty"`
	Memo          String        `xml:"MEMO,omitempty"`
	Currency      *Currency     `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	Inv401kSource inv401kSource `xml:"INV401KSOURCE,omitempty"` // One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

// OpenOrder is an interface satisfied by all the OO* elements.
type OpenOrder interface {
	OrderType() string
}

// OOBuyDebt represents an open order to purchase a debt security
type OOBuyDebt struct {
	XMLName   xml.Name `xml:"OOBUYDEBT"`
	OO        OO       `xml:"OO"`
	Auction   Boolean  `xml:"AUCTION"` // whether the debt should be purchased at the auction
	DtAuction *Date    `xml:"DTAUCTION,omitempty"`
}

// OrderType returns a string representation of this order's type
func (o OOBuyDebt) OrderType() string {
	return "OOBUYDEBT"
}

// OOBuyMF represents an open order to purchase a mutual fund
type OOBuyMF struct {
	XMLName  xml.Name `xml:"OOBUYMF"`
	OO       OO       `xml:"OO"`
	BuyType  buyType  `xml:"BUYTYPE"`  // One of BUY, BUYTOCOVER
	UnitType unitType `xml:"UNITTYPE"` // What the units represent: one of SHARES, CURRENCY
}

// OrderType returns a string representation of this order's type
func (o OOBuyMF) OrderType() string {
	return "OOBUYMF"
}

// OOBuyOpt represents an open order to purchase an option
type OOBuyOpt struct {
	XMLName    xml.Name   `xml:"OOBUYOPT"`
	OO         OO         `xml:"OO"`
	OptBuyType optBuyType `xml:"OPTBUYTYPE"` // One of BUYTOOPEN, BUYTOCLOSE
}

// OrderType returns a string representation of this order's type
func (o OOBuyOpt) OrderType() string {
	return "OOBUYOPT"
}

// OOBuyOther represents an open order to purchase a security type not covered
// by the other OOBuy* elements
type OOBuyOther struct {
	XMLName  xml.Name `xml:"OOBUYOTHER"`
	OO       OO       `xml:"OO"`
	UnitType unitType `xml:"UNITTYPE"` // What the units represent: one of SHARES, CURRENCY
}

// OrderType returns a string representation of this order's type
func (o OOBuyOther) OrderType() string {
	return "OOBUYOTHER"
}

// OOBuyStock represents an open order to purchase stock
type OOBuyStock struct {
	XMLName xml.Name `xml:"OOBUYSTOCK"`
	OO      OO       `xml:"OO"`
	BuyType buyType  `xml:"BUYTYPE"` // One of BUY, BUYTOCOVER
}

// OrderType returns a string representation of this order's type
func (o OOBuyStock) OrderType() string {
	return "OOBUYSTOCK"
}

// OOSellDebt represents an open order to sell a debt security
type OOSellDebt struct {
	XMLName xml.Name `xml:"OOSELLDEBT"`
	OO      OO       `xml:"OO"`
}

// OrderType returns a string representation of this order's type
func (o OOSellDebt) OrderType() string {
	return "OOSELLDEBT"
}

// OOSellMF represents an open order to sell a mutual fund
type OOSellMF struct {
	XMLName  xml.Name `xml:"OOSELLMF"`
	OO       OO       `xml:"OO"`
	SellType sellType `xml:"SELLTYPE"` // One of SELL, SELLSHORT
	UnitType unitType `xml:"UNITTYPE"` // What the units represent: one of SHARES, CURRENCY
	SellAll  Boolean  `xml:"SELLALL"`  // Sell entire holding
}

// OrderType returns a string representation of this order's type
func (o OOSellMF) OrderType() string {
	return "OOSELLMF"
}

// OOSellOpt represents an open order to sell an option
type OOSellOpt struct {
	XMLName     xml.Name    `xml:"OOSELLOPT"`
	OO          OO          `xml:"OO"`
	OptSellType optSellType `xml:"OPTSELLTYPE"` // One of SELLTOOPEN, SELLTOCLOSE
}

// OrderType returns a string representation of this order's type
func (o OOSellOpt) OrderType() string {
	return "OOSELLOPT"
}

// OOSellOther represents an open order to sell a security type not covered by
// the other OOSell* elements
type OOSellOther struct {
	XMLName  xml.Name `xml:"OOSELLOTHER"`
	OO       OO       `xml:"OO"`
	UnitType unitType `xml:"UNITTYPE"` // What the units represent: one of SHARES, CURRENCY
}

// OrderType returns a string representation of this order's type
func (o OOSellOther) OrderType() string {
	return "OOSELLOTHER"
}

// OOSellStock represents an open order to sell stock
type OOSellStock struct {
	XMLName  xml.Name `xml:"OOSELLSTOCK"`
	OO       OO       `xml:"OO"`
	SellType sellType `xml:"SELLTYPE"` // One of SELL, SELLSHORT
}

// OrderType returns a string representation of this order's type
func (o OOSellStock) OrderType() string {
	return "OOSELLSTOCK"
}

// OOSwitchMF represents an open order to switch to or purchase a different
// mutual fund
type OOSwitchMF struct {
	XMLName   xml.Name   `xml:"SWITCHMF"`
	OO        OO         `xml:"OO"`
	SecID     SecurityID `xml:"SECID"`     // Security ID of the fund to switch to or purchase
	UnitType  unitType   `xml:"UNITTYPE"`  // What the units represent: one of SHARES, CURRENCY
	SwitchAll Boolean    `xml:"SWITCHALL"` // Switch entire holding
}

// OrderType returns a string representation of this order's type
func (o OOSwitchMF) OrderType() string {
	return "SWITCHMF"
}

// OOList represents a list of open orders (OO* elements)
type OOList []OpenOrder

// UnmarshalXML handles unmarshalling an OOList element from an XML string
func (o *OOList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "OOBUYDEBT":
				var oo OOBuyDebt
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOBUYMF":
				var oo OOBuyMF
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOBUYOPT":
				var oo OOBuyOpt
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOBUYOTHER":
				var oo OOBuyOther
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOBUYSTOCK":
				var oo OOBuyStock
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOSELLDEBT":
				var oo OOSellDebt
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOSELLMF":
				var oo OOSellMF
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOSELLOPT":
				var oo OOSellOpt
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOSELLOTHER":
				var oo OOSellOther
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "OOSELLSTOCK":
				var oo OOSellStock
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			case "SWITCHMF":
				var oo OOSwitchMF
				if err := d.DecodeElement(&oo, &startElement); err != nil {
					return err
				}
				*o = append(*o, OpenOrder(oo))
			default:
				return errors.New("Invalid OOList child tag: " + startElement.Name.Local)
			}
		} else {
			return errors.New("Didn't find an opening element")
		}
	}
}

// MarshalXML handles marshalling an OOList to an XML string
func (o OOList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	ooListElement := xml.StartElement{Name: xml.Name{Local: "INVOOLIST"}}
	if err := e.EncodeToken(ooListElement); err != nil {
		return err
	}
	for _, openorder := range o {
		start := xml.StartElement{Name: xml.Name{Local: openorder.OrderType()}}
		switch oo := openorder.(type) {
		case OOBuyDebt:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOBuyMF:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOBuyOpt:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOBuyOther:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOBuyStock:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSellDebt:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSellMF:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSellOpt:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSellOther:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSellStock:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		case OOSwitchMF:
			if err := e.EncodeElement(&oo, start); err != nil {
				return err
			}
		default:
			return errors.New("Invalid OOLIST child type: " + oo.OrderType())
		}
	}
	if err := e.EncodeToken(ooListElement.End()); err != nil {
		return err
	}
	return nil
}

// ContribSecurity identifies current contribution allocation for a security in
// a 401(k) account
type ContribSecurity struct {
	XMLName                 xml.Name   `xml:"CONTRIBSECURITY"`
	SecID                   SecurityID `xml:"SECID"`
	PreTaxContribPct        Amount     `xml:"PRETAXCONTRIBPCT,omitempty"`        // Percentage of each new employee pretax contribution allocated to this security, rate.
	PreTaxContribAmt        Amount     `xml:"PRETAXCONTRIBAMT,omitempty"`        // Fixed amount of each new employee pretax contribution allocated to this security, amount
	AfterTaxContribPct      Amount     `xml:"AFTERTAXCONTRIBPCT,omitempty"`      // Percentage of each new employee after tax contribution allocated to this security, rate.
	AfterTaxContribAmt      Amount     `xml:"AFTERTAXCONTRIBAMT,omitempty"`      // Fixed amount of each new employee pretax contribution allocated to this security, amount.
	MatchContribPct         Amount     `xml:"MATCHCONTRIBPCT,omitempty"`         // Percentage of each new employer match contribution allocated to this security, rate.
	MatchContribAmt         Amount     `xml:"MATCHCONTRIBAMT,omitempty"`         // Fixed amount of each new employer match contribution allocated to this security, amount.
	ProfitSharingContribPct Amount     `xml:"PROFITSHARINGCONTRIBPCT,omitempty"` // Percentage of each new employer profit sharing contribution allocated to this security, rate.
	ProfitSharingContribAmt Amount     `xml:"PROFITSHARINGCONTRIBAMT,omitempty"` // Fixed amount of each new employer profit sharing contribution allocated to this security, amount.
	RolloverContribPct      Amount     `xml:"ROLLOVERCONTRIBPCT,omitempty"`      // Percentage of new rollover contributions allocated to this security, rate.
	RolloverContribAmt      Amount     `xml:"ROLLOVERCONTRIBAMT,omitempty"`      // Fixed amount of new rollover contributions allocated to this security, amount.
	OtherVestPct            Amount     `xml:"OTHERVESTPCT,omitempty"`            // Percentage of each new other employer contribution allocated to this security, rate.
	OtherVestAmt            Amount     `xml:"OTHERVESTAMT,omitempty"`            // Fixed amount of each new other employer contribution allocated to this security, amount.
	OtherNonVestPct         Amount     `xml:"OTHERNONVESTPCT,omitempty"`         // Percentage of each new other employee contribution allocated to this security, rate.
	OtherNonVestAmt         Amount     `xml:"OTHERNONVESTAMT,omitempty"`         // Fixed amount of each new other employee contribution allocated to this security, amount
}

// VestInfo provides the vesting percentage of a 401(k) account as of a
// particular date (past, present, or future)
type VestInfo struct {
	XMLName  xml.Name `xml:"VESTINFO"`
	VestDate *Date    `xml:"VESTDATE,omitempty"` // Date at which vesting percentage changes. Default (if empty) is that the vesting percentage below applies to the current date
	VestPct  Amount   `xml:"VESTPCT"`
}

// LoanInfo represents a loan outstanding against this 401(k) account
type LoanInfo struct {
	XMLName               xml.Name    `xml:"VESTINFO"`
	LoanID                String      `xml:"LOANID"`                          // Identifier of this loan
	LoanDesc              String      `xml:"LOANDESC,omitempty"`              // Loan description
	InitialLoanBal        Amount      `xml:"INITIALLOANBAL,omitempty"`        // Initial loan balance
	LoanStartDate         *Date       `xml:"LOANSTARTDATE,omitempty"`         // Start date of loan
	CurrentLoanBal        Amount      `xml:"CURRENTLOANBAL"`                  // Current loan principal balance
	DtAsOf                *Date       `xml:"DTASOF"`                          // Date and time of the current loan balance
	LoanRate              Amount      `xml:"LOANRATE,omitempty"`              // Loan annual interest rate
	LoanPmtAmt            Amount      `xml:"LOANPMTAMT,omitempty"`            // Loan payment amount
	LoanPmtFreq           loanPmtFreq `xml:"LOANPMTFREQ,omitempty"`           // Frequency of loan repayments: WEEKLY, BIWEEKLY, TWICEMONTHLY, MONTHLY, FOURWEEKS, BIMONTHLY, QUARTERLY, SEMIANNUALLY, ANNUALLY, OTHER. See section 10.2.1 for calculation rules.
	LoanPmtsInitial       Int         `xml:"LOANPMTSINITIAL,omitempty"`       // Initial number of loan payments.
	LoanPmtsRemaining     Int         `xml:"LOANPMTSREMAINING,omitempty"`     // Remaining number of loan payments
	LoanMaturityDate      *Date       `xml:"LOANMATURITYDATE,omitempty"`      // Expected loan end date
	LoanTotalProjInterest Amount      `xml:"LOANTOTALPROJINTEREST,omitempty"` // Total projected interest to be paid on this loan
	LoanInterestToDate    Amount      `xml:"LOANINTERESTTODATE,omitempty"`    // Total interested paid to date on this loan
	LoanExtPmtDate        *Date       `xml:"LOANNEXTPMTDATE,omitempty"`       // Next payment due date
}

// Inv401KSummaryAggregate represents the total of either contributions,
// withdrawals, or earnings made in each contribution type in a given period
// (dates specified in a containing Inv401KSummaryPeriod)
type Inv401KSummaryAggregate struct {
	XMLName       xml.Name // One of CONTRIBUTIONS, WITHDRAWALS, EARNINGS
	PreTax        Amount   `xml:"PRETAX,omitempty"`        // Pretax contributions, withdrawals, or earlings.
	AfterTax      Amount   `xml:"AFTERTAX,omitempty"`      // After tax contributions, withdrawals, or earlings.
	Match         Amount   `xml:"MATCH,omitempty"`         // Employer matching contributions, withdrawals, or earlings.
	ProfitSharing Amount   `xml:"PROFITSHARING,omitempty"` // Profit sharing contributions, withdrawals, or earlings.
	Rollover      Amount   `xml:"ROLLOVER,omitempty"`      // Rollover contributions, withdrawals, or earlings.
	OtherVest     Amount   `xml:"OTHERVEST,omitempty"`     // Other vesting contributions, withdrawals, or earlings.
	OtherNonVest  Amount   `xml:"OTHERNONVEST,omitempty"`  // Other non-vesting contributions, withdrawals, or earlings.
	Total         Amount   `xml:"TOTAL"`                   // Sum of contributions, withdrawals, or earlings from all fund sources.
}

// Inv401KSummaryPeriod contains the total contributions, withdrawals, and
// earnings made in the given date range
type Inv401KSummaryPeriod struct {
	XMLName       xml.Name                 // One of YEARTODATE, INCEPTODATE, or PERIODTODATE
	DtStart       Date                     `xml:"DTSTART"`
	DtEnd         Date                     `xml:"DTEND"`
	Contributions *Inv401KSummaryAggregate `xml:"CONTRIBUTIONS,omitempty"` // 401(k) contribution aggregate. Note: this includes loan payments.
	Withdrawls    *Inv401KSummaryAggregate `xml:"WITHDRAWLS,omitempty"`    // 401(k) withdrawals aggregate. Note: this includes loan withdrawals.
	Earnings      *Inv401KSummaryAggregate `xml:"EARNINGS,omitempty"`      // 401(k) earnings aggregate. This is the market value change. It includes dividends/interest, and capital gains - realized and unrealized.
}

// Inv401K is included in InvStatementResponse for 401(k) accounts and provides
// a summary of the 401(k) specific information about the user's account.
type Inv401K struct {
	XMLName             xml.Name `xml:"INV401K"`
	EmployerName        String   `xml:"EMPLOYERNAME"`
	PlanID              String   `xml:"PLANID,omitempty"`              // Plan number
	PlanJoinDate        *Date    `xml:"PLANJOINDATE,omitempty"`        // Date the employee joined the plan
	EmployerContactInfo String   `xml:"EMPLOYERCONTACTINFO,omitempty"` // Name of contact person at employer, plus any available contact information, such as phone number
	BrokerContactInfo   String   `xml:"BROKERCONTACTINFO,omitempty"`   // Name of contact person at broker, plus any available contact information, such as phone number
	DeferPctPreTax      Amount   `xml:"DEFERPCTPRETAX,omitempty"`      // Percent of employee salary deferred before tax
	DeferPctAfterTax    Amount   `xml:"DEFERPCTAFTERTAX,omitempty"`    // Percent of employee salary deferred after tax

	//<MATCHINFO> Aggregate containing employer match information. Absent if employer does not contribute matching funds.
	MatchPct            Amount                `xml:"MATCHINFO>MATCHPCT,omitempty"`          // Percent of employee contribution matched, e.g., 75% if contribution rate is $0.75/$1.00
	MaxMatchAmt         Amount                `xml:"MATCHINFO>MAXMATCHAMT,omitempty"`       // Maximum employer contribution amount in any year
	MaxMatchPct         Amount                `xml:"MATCHINFO>MAXMATCHPCT,omitempty"`       // Current maximum employer contribution percentage. Maximum match in a year is MAXMATCHPCT up to the MAXMATCHAMT, if provided
	StartOfYear         *Date                 `xml:"MATCHINFO>STARTOFYEAR,omitempty"`       // Specifies when the employer contribution max is reset. Some plans have a maximum based on the company fiscal year rather than calendar year. Assume calendar year if omitted. Only the month and day (MMDD) are used; year (YYYY) and time are ignored
	BaseMatchAmt        Amount                `xml:"MATCHINFO>BASEMATCHAMT"`                // Specifies a fixed dollar amount contributed by the employer if the employee participates in the plan at all. This may be present in addition to the <MATCHPCT>. $0 if omitted
	BaseMatchPct        Amount                `xml:"MATCHINFO>BASEMATCHPCT"`                // Specifies a fixed percent of employee salary matched if the employee participates in the plan at all. This may be present in addition to the MATCHPCT>. 0% if omitted. Base match in a year is BASEMATCHPCT up to the BASEMATCHAMT,if provided
	ContribInfo         []ContribSecurity     `xml:"CONTRIBINTO>CONTRIBSECURITY"`           // Aggregate to describe how new contributions are distributed among the available securities.
	CurrentVestPct      Amount                `xml:"CURRENTVESTPCT,omitempty"`              // Estimated percentage of employer contributions vested as of the current date. If omitted, assume 100%
	VestInfo            []VestInfo            `xml:"VESTINFO,omitempty"`                    // Vest change dates. Provides the vesting percentage as of any particular past, current, or future date. 0 or more.
	LoanInfo            []LoanInfo            `xml:"LOANINFO,omitempty"`                    // List of any loans outstanding against this account
	YearToDateSummary   Inv401KSummaryPeriod  `xml:"INV401KSUMMARY>YEARTODATE"`             // Contributions to date for this calendar year.
	InceptToDateSummary *Inv401KSummaryPeriod `xml:"INV401KSUMMARY>INCEPTODATE,omitempty"`  // Total contributions to date (since inception)
	PeriodToDate        *Inv401KSummaryPeriod `xml:"INV401KSUMMARY>PERIODTODATE,omitempty"` // Total contributions this contribution period
}

// Inv401KBal provides the balances for different 401(k) subaccount types, as
// well as the total cash value of the securities held
type Inv401KBal struct {
	XMLName       xml.Name  `xml:"INV401KBAL"`
	CashBal       Amount    `xml:"CASHBAL,omitempty"`       // Available cash balance
	PreTax        Amount    `xml:"PRETAX,omitempty"`        // Current value of all securities purchased with Before Tax Employee contributions
	AfterTax      Amount    `xml:"AFTERTAX,omitempty"`      // Current value of all securities purchased with After Tax Employee contributions
	Match         Amount    `xml:"MATCH,omitempty"`         // Current value of all securities purchased with Employer Match contributions
	ProfitSharing Amount    `xml:"PROFITSHARING,omitempty"` // Current value of all securities purchased with Employer Profit Sharing contributions
	Rollover      Amount    `xml:"ROLLOVER,omitempty"`      // Current value of all securities purchased with Rollover contributions
	OtherVest     Amount    `xml:"OTHERVEST,omitempty"`     // Current value of all securities purchased with Other (vesting) Employer contributions
	OtherNonVest  Amount    `xml:"OTHERNONVEST,omitempty"`  // Current value of all securities purchased with Other (non-vesting) Employer contributions
	Total         Amount    `xml:"TOTAL"`                   // Current value of all securities purchased with all contributions
	BalList       []Balance `xml:"BALLIST>BAL,omitempty"`
}

// InvStatementResponse includes requested transaction, position, open order,
// and balance information for an investment account. It is in response to an
// InvStatementRequest or sometimes provided as part of an OFX file downloaded
// manually from an FI.
type InvStatementResponse struct {
	XMLName   xml.Name `xml:"INVSTMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	DtAsOf      Date         `xml:"INVSTMTRS>DTASOF"`
	CurDef      CurrSymbol   `xml:"INVSTMTRS>CURDEF"`
	InvAcctFrom InvAcct      `xml:"INVSTMTRS>INVACCTFROM"`
	InvTranList *InvTranList `xml:"INVSTMTRS>INVTRANLIST,omitempty"`
	InvPosList  PositionList `xml:"INVSTMTRS>INVPOSLIST,omitempty"`
	InvBal      *InvBalance  `xml:"INVSTMTRS>INVBAL,omitempty"`
	InvOOList   OOList       `xml:"INVSTMTRS>INVOOLIST,omitempty"`
	MktgInfo    String       `xml:"INVSTMTRS>MKTGINFO,omitempty"` // Marketing information
	Inv401K     *Inv401K     `xml:"INVSTMTRS>INV401K,omitempty"`
	Inv401KBal  *Inv401KBal  `xml:"INVSTMTRS>INV401KBAL,omitempty"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (sr *InvStatementResponse) Name() string {
	return "INVSTMTTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (sr *InvStatementResponse) Valid(version ofxVersion) (bool, error) {
	if ok, err := sr.TrnUID.Valid(); !ok {
		return false, err
	}
	//TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (sr *InvStatementResponse) Type() messageType {
	return InvStmtRs
}
