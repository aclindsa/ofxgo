package ofxgo

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

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

func (r *InvStatementRequest) Name() string {
	return "INVSTMTTRNRQ"
}

func (r *InvStatementRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

type InvTran struct {
	XMLName       xml.Name `xml:"INVTRAN"`
	FiTId         String   `xml:"FITID"`
	SrvrTId       String   `xml:"SRVRTID,omitempty"`
	DtTrade       Date     `xml:"DTTRADE"`                 // trade date; for stock splits, day of record
	DtSettle      *Date    `xml:"DTSETTLE,omitempty"`      // settlement date; for stock splits, execution date
	ReversalFiTId String   `xml:"REVERSALFITID,omitempty"` // For a reversal transaction, the FITID of the transaction that is being reversed.
	Memo          String   `xml:"MEMO,omitempty"`
}

type InvBuy struct {
	XMLName      xml.Name   `xml:"INVBUY"`
	InvTran      InvTran    `xml:"INVTRAN"`
	SecId        SecurityId `xml:"SECID"`
	Units        Amount     `xml:"UNITS"`            // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice    Amount     `xml:"UNITPRICE"`        // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Markup       Amount     `xml:"MARKUP,omitempty"` // Portion of UNITPRICE that is attributed to the dealer markup
	Commission   Amount     `xml:"COMMISSION,omitempty"`
	Taxes        Amount     `xml:"TAXES,omitempty"`
	Fees         Amount     `xml:"FEES,omitempty"`
	Load         Amount     `xml:"LOAD,omitempty"`
	Total        Amount     `xml:"TOTAL"`                  // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	Currency     *Currency  `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
	OrigCurrency *Currency  `xml:"ORIGCURRENCY,omitempty"` // Overriding currency for UNITPRICE
	SubAcctSec   String     `xml:"SUBACCTSEC"`             // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund  String     `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER

	// The next three elements must either all be provided, or none of  htem
	LoanId        String `xml:"LOANID,omitempty"`        // For 401(k) accounts only. Indicates that the transaction was due to a loan or a loan repayment, and which loan it was
	LoanPrincipal Amount `xml:"LOANPRINCIPAL,omitempty"` // For 401(k) accounts only. Indicates how much of the loan repayment was principal
	LoanInterest  Amount `xml:"LOANINTEREST,omitempty"`  // For 401(k) accounts only. Indicates how much of the loan repayment was interest

	Inv401kSource    String  `xml:"INV401KSOURCE,omitempty"`    // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
	DtPayroll        *Date   `xml:"DTPAYROLL,omitempty"`        // For 401(k)accounts, date the funds for this transaction was obtained via payroll deduction
	PriorYearContrib Boolean `xml:"PRIORYEARCONTRIB,omitempty"` // For 401(k) accounts, indicates that this Buy was made with a prior year contribution
}

type InvSell struct {
	XMLName      xml.Name   `xml:"INVSELL"`
	InvTran      InvTran    `xml:"INVTRAN"`
	SecId        SecurityId `xml:"SECID"`
	Units        Amount     `xml:"UNITS"`              // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice    Amount     `xml:"UNITPRICE"`          // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Markdown     Amount     `xml:"MARKDOWN,omitempty"` // Portion of UNITPRICE that is attributed to the dealer markdown
	Commission   Amount     `xml:"COMMISSION,omitempty"`
	Taxes        Amount     `xml:"TAXES,omitempty"`
	Fees         Amount     `xml:"FEES,omitempty"`
	Load         Amount     `xml:"LOAD,omitempty"`
	Witholding   Amount     `xml:"WITHHOLDING,omitempty"`  // Federal tax witholdings
	TaxExempt    Boolean    `xml:"TAXEXEMPT,omitempty"`    // Tax-exempt transaction
	Total        Amount     `xml:"TOTAL"`                  // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	Gain         Amount     `xml:"GAIN,omitempty"`         // Total gain
	Currency     *Currency  `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
	OrigCurrency *Currency  `xml:"ORIGCURRENCY,omitempty"` // Overriding currency for UNITPRICE
	SubAcctSec   String     `xml:"SUBACCTSEC"`             // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund  String     `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER

	LoanId          String `xml:"LOANID,omitempty"`           // For 401(k) accounts only. Indicates that the transaction was due to a loan or a loan repayment, and which loan it was
	StateWitholding Amount `xml:"STATEWITHHOLDING,omitempty"` // State tax witholdings
	Penalty         Amount `xml:"PENALTY,omitempty"`          // Amount witheld due to penalty

	Inv401kSource String `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

type BuyDebt struct {
	XMLName  xml.Name `xml:"BUYDEBT"`
	InvBuy   InvBuy   `xml:"INVBUY"`
	AccrdInt Amount   `xml:"ACCRDINT,omitempty"` // Accrued interest. This amount is not reflected in the <TOTAL> field of a containing aggregate.
}

func (t BuyDebt) TransactionType() string {
	return "BUYDEBT"
}

type BuyMF struct {
	XMLName  xml.Name `xml:"BUYMF"`
	InvBuy   InvBuy   `xml:"INVBUY"`
	BuyType  String   `xml:"BUYTYPE"`            // One of BUY, BUYTOCOVER (BUYTOCOVER used to close short sales.)
	RelFiTId String   `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
}

func (t BuyMF) TransactionType() string {
	return "BUYMF"
}

type BuyOpt struct {
	XMLName    xml.Name `xml:"BUYOPT"`
	InvBuy     InvBuy   `xml:"INVBUY"`
	OptBuyType String   `xml:"OPTBUYTYPE"` // type of purchase: BUYTOOPEN, BUYTOCLOSE (The BUYTOOPEN buy type is like “ordinary” buying of option and works like stocks.)
	ShPerCtrct Int      `xml:"SHPERCTRCT"` // Shares per contract
}

func (t BuyOpt) TransactionType() string {
	return "BUYOPT"
}

type BuyOther struct {
	XMLName xml.Name `xml:"BUYOTHER"`
	InvBuy  InvBuy   `xml:"INVBUY"`
}

func (t BuyOther) TransactionType() string {
	return "BUYOTHER"
}

type BuyStock struct {
	XMLName xml.Name `xml:"BUYSTOCK"`
	InvBuy  InvBuy   `xml:"INVBUY"`
	BuyType String   `xml:"BUYTYPE"` // One of BUY, BUYTOCOVER (BUYTOCOVER used to close short sales.)
}

func (t BuyStock) TransactionType() string {
	return "BUYSTOCK"
}

type ClosureOpt struct {
	XMLName    xml.Name   `xml:"CLOSUREOPT"`
	InvTran    InvTran    `xml:"INVTRAN"`
	SecId      SecurityId `xml:"SECID"`
	OptAction  String     `xml:"OPTACTION"`          // One of EXERCISE, ASSIGN, EXPIRE. The EXERCISE action is used to close out an option that is exercised. The ASSIGN action is used when an option writer is assigned. The EXPIRE action is used when the option’s expired date is reached
	Units      Amount     `xml:"UNITS"`              // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	ShPerCtrct Int        `xml:"SHPERCTRCT"`         // Shares per contract
	SubAcctSec String     `xml:"SUBACCTSEC"`         // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	RelFiTId   String     `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
	Gain       Amount     `xml:"GAIN,omitempty"`     // Total gain
}

func (t ClosureOpt) TransactionType() string {
	return "CLOSUREOPT"
}

// Investment income is realized as cash into the investment account
type Income struct {
	XMLName       xml.Name   `xml:"INCOME"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	IncomeType    String     `xml:"INCOMETYPE"` // Type of investment income: CGLONG (capital gains-long term), CGSHORT (capital gains-short term), DIV (dividend), INTEREST, MISC
	Total         Amount     `xml:"TOTAL"`
	SubAcctSec    String     `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   String     `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	TaxExempt     Boolean    `xml:"TAXEXEMPT,omitempty"`     // Tax-exempt transaction
	Witholding    Amount     `xml:"WITHHOLDING,omitempty"`   // Federal tax witholdings
	Currency      *Currency  `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	OrigCurrency  *Currency  `xml:"ORIGCURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t Income) TransactionType() string {
	return "INCOME"
}

// Expense associated with an investment
type InvExpense struct {
	XMLName       xml.Name   `xml:"INVEXPENSE"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	Total         Amount     `xml:"TOTAL"`
	SubAcctSec    String     `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   String     `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency      *Currency  `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	OrigCurrency  *Currency  `xml:"ORIGCURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t InvExpense) TransactionType() string {
	return "INVEXPENSE"
}

// Journaling cash holdings between sub-accounts within the same investment account.
type JrnlFund struct {
	XMLName     xml.Name `xml:"JRNLFUND"`
	InvTran     InvTran  `xml:"INVTRAN"`
	Total       Amount   `xml:"TOTAL"`
	SubAcctFrom String   `xml:"SUBACCTFROM"` // Sub-account cash is being transferred from: CASH, MARGIN, SHORT, OTHER
	SubAcctTo   String   `xml:"SUBACCTTO"`   // Sub-account cash is being transferred to: CASH, MARGIN, SHORT, OTHER
}

func (t JrnlFund) TransactionType() string {
	return "JRNLFUND"
}

// Journaling security holdings between sub-accounts within the same investment account.
type JrnlSec struct {
	XMLName     xml.Name   `xml:"JRNLSEC"`
	InvTran     InvTran    `xml:"INVTRAN"`
	SecId       SecurityId `xml:"SECID"`
	SubAcctFrom String     `xml:"SUBACCTFROM"` // Sub-account cash is being transferred from: CASH, MARGIN, SHORT, OTHER
	SubAcctTo   String     `xml:"SUBACCTTO"`   // Sub-account cash is being transferred to: CASH, MARGIN, SHORT, OTHER
	Units       Amount     `xml:"UNITS"`       // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
}

func (t JrnlSec) TransactionType() string {
	return "JRNLSEC"
}

type MarginInterest struct {
	XMLName      xml.Name  `xml:"MARGININTEREST"`
	InvTran      InvTran   `xml:"INVTRAN"`
	Total        Amount    `xml:"TOTAL"`
	SubAcctFund  String    `xml:"SUBACCTFUND"`            // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency     *Currency `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
	OrigCurrency *Currency `xml:"ORIGCURRENCY,omitempty"` // Overriding currency for UNITPRICE
}

func (t MarginInterest) TransactionType() string {
	return "MARGININTEREST"
}

// REINVEST is a single transaction that contains both income and an investment transaction. If servers can’t track this as a single transaction they should return an INCOME transaction and an INVTRAN.
type Reinvest struct {
	XMLName       xml.Name   `xml:"REINVEST"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	IncomeType    String     `xml:"INCOMETYPE"` // Type of investment income: CGLONG (capital gains-long term), CGSHORT (capital gains-short term), DIV (dividend), INTEREST, MISC
	Total         Amount     `xml:"TOTAL"`      // Transaction total. Buys, sells, etc.:((quan. * (price +/- markup/markdown)) +/-(commission + fees + load + taxes + penalty + withholding + statewithholding)). Distributions, interest, margin interest, misc. expense, etc.: amount. Return of cap: cost basis
	SubAcctSec    String     `xml:"SUBACCTSEC"` // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	Units         Amount     `xml:"UNITS"`      // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	UnitPrice     Amount     `xml:"UNITPRICE"`  // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	Commission    Amount     `xml:"COMMISSION,omitempty"`
	Taxes         Amount     `xml:"TAXES,omitempty"`
	Fees          Amount     `xml:"FEES,omitempty"`
	Load          Amount     `xml:"LOAD,omitempty"`
	TaxExempt     Boolean    `xml:"TAXEXEMPT,omitempty"`     // Tax-exempt transaction
	Currency      *Currency  `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	OrigCurrency  *Currency  `xml:"ORIGCURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t Reinvest) TransactionType() string {
	return "REINVEST"
}

type RetOfCap struct {
	XMLName       xml.Name   `xml:"RETOFCAP"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	Total         Amount     `xml:"TOTAL"`
	SubAcctSec    String     `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	SubAcctFund   String     `xml:"SUBACCTFUND"`             // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Currency      *Currency  `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	OrigCurrency  *Currency  `xml:"ORIGCURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t RetOfCap) TransactionType() string {
	return "RETOFCAP"
}

type SellDebt struct {
	XMLName    xml.Name `xml:"SELLDEBT"`
	InvSell    InvSell  `xml:"INVSELL"`
	SellReason String   `xml:"SELLREASON"`         // CALL (the debt was called), SELL (the debt was sold), MATURITY (the debt reached maturity)
	AccrdInt   Amount   `xml:"ACCRDINT,omitempty"` // Accrued interest
}

func (t SellDebt) TransactionType() string {
	return "SELLDEBT"
}

type SellMF struct {
	XMLName      xml.Name `xml:"SELLMF"`
	InvSell      InvSell  `xml:"INVSELL"`
	SellType     String   `xml:"SELLTYPE"` // Type of sell. SELL, SELLSHORT
	AvgCostBasis Amount   `xml:"AVGCOSTBASIS"`
	RelFiTId     String   `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
}

func (t SellMF) TransactionType() string {
	return "SELLMF"
}

type SellOpt struct {
	XMLName     xml.Name `xml:"SELLOPT"`
	InvSell     InvSell  `xml:"INVSELL"`
	OptSellType String   `xml:"OPTSELLTYPE"`        // For options, type of sell: SELLTOCLOSE, SELLTOOPEN. The SELLTOCLOSE action is selling a previously bought option. The SELLTOOPEN action is writing an option
	ShPerCtrct  Int      `xml:"SHPERCTRCT"`         // Shares per contract
	RelFiTId    String   `xml:"RELFITID,omitempty"` // used to relate transactions associated with mutual fund exchanges
	RelType     String   `xml:"RELTYPE,omitempty"`  // Related option transaction type: SPREAD, STRADDLE, NONE, OTHER
	Secured     String   `xml:"SECURED,omitempty"`  // NAKED, COVERED
}

func (t SellOpt) TransactionType() string {
	return "SELLOPT"
}

type SellOther struct {
	XMLName xml.Name `xml:"SELLOTHER"`
	InvSell InvSell  `xml:"INVSELL"`
}

func (t SellOther) TransactionType() string {
	return "SELLOTHER"
}

type SellStock struct {
	XMLName  xml.Name `xml:"SELLSTOCK"`
	InvSell  InvSell  `xml:"INVSELL"`
	SellType String   `xml:"SELLTYPE"` // Type of sell. SELL, SELLSHORT
}

func (t SellStock) TransactionType() string {
	return "SELLSTOCK"
}

type Split struct {
	XMLName       xml.Name   `xml:"SPLIT"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	SubAcctSec    String     `xml:"SUBACCTSEC"`              // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	OldUnits      Amount     `xml:"OLDUNITS"`                // number of shares before the split
	NewUnits      Amount     `xml:"NEWUNITS"`                // number of shares after the split
	Numerator     Int        `xml:"NUMERATOR"`               // split ratio numerator
	Denominator   Int        `xml:"DENOMINATOR"`             // split ratio denominator
	Currency      *Currency  `xml:"CURRENCY,omitempty"`      // Overriding currency for UNITPRICE
	OrigCurrency  *Currency  `xml:"ORIGCURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	FracCash      Amount     `xml:"FRACCASH,omitempty"`      // cash for fractional units
	SubAcctFund   String     `xml:"SUBACCTFUND,omitempty"`   // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t Split) TransactionType() string {
	return "SPLIT"
}

type Transfer struct {
	XMLName       xml.Name   `xml:"TRANSFER"`
	InvTran       InvTran    `xml:"INVTRAN"`
	SecId         SecurityId `xml:"SECID"`
	SubAcctSec    String     `xml:"SUBACCTSEC"` // Sub-account type for this security. One of CASH, MARGIN, SHORT, OTHER
	Units         Amount     `xml:"UNITS"`      // For stocks, MFs, other, number of shares held. Bonds = face value. Options = number of contracts
	TferAction    String     `xml:"TFERACTION"` // One of IN, OUT
	PosType       String     `xml:"POSTYPE"`    // Position type. One of LONG, SHORT
	InvAcctFrom   InvAcct    `xml:"INVACCTFROM,omitempty"`
	AvgCostBasis  Amount     `xml:"AVGCOSTBASIS,omitempty"`
	UnitPrice     Amount     `xml:"UNITPRICE,omitempty"` // For stocks, MFs, other, price per share. Bonds = percentage of par. Option = premium per share of underlying security
	DtPurchase    *Date      `xml:"DTPURCHASE,omitempty"`
	Inv401kSource String     `xml:"INV401KSOURCE,omitempty"` // Source of money for this transaction. One of PRETAX, AFTERTAX, MATCH, PROFITSHARING, ROLLOVER, OTHERVEST, OTHERNONVEST for 401(k) accounts. Default if not present is OTHERNONVEST. The following cash source types are subject to vesting: MATCH, PROFITSHARING, and OTHERVEST
}

func (t Transfer) TransactionType() string {
	return "TRANSFER"
}

type InvTransaction interface {
	TransactionType() string
}

type InvBankTransaction struct {
	XMLName      xml.Name      `xml:"INVBANKTRAN"`
	Transactions []Transaction `xml:"STMTTRN,omitempty"`
	SubAcctFund  String        `xml:"SUBACCTFUND"` // Where did the money for the transaction come from or go to? CASH, MARGIN, SHORT, OTHER
}

// Must be unmarshalled manually due to the structure (don't know what kind of
// InvTransaction is coming next)
type InvTranList struct {
	DtStart          Date
	DtEnd            Date
	InvTransactions  []InvTransaction
	BankTransactions []InvBankTransaction
}

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
	Currency      *Currency  `xml:"CURRENCY,omitempty"`     // Overriding currency for UNITPRICE
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

type InvBalance struct {
	XMLName       xml.Name  `xml:"INVBAL"`
	AvailCash     Amount    `xml:"AVAILCASH"`     // Available cash across all sub-accounts, including sweep funds
	MarginBalance Amount    `xml:"MARGINBALANCE"` // Negative means customer has borrowed funds
	ShortBalance  Amount    `xml:"SHORTBALANCE"`  // Always positive, market value of all short positions
	BuyPower      Amount    `xml:"BUYPOWER"`
	BalList       []Balance `xml:"BALLIST>BAL,omitempty"`
}

type ContribSecurity struct {
	XMLName                 xml.Name   `xml:"CONTRIBSECURITY"`
	SecId                   SecurityId `xml:"SECID"`
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

type VestInfo struct {
	XMLName  xml.Name `xml:"VESTINFO"`
	VestDate *Date    `xml:"VESTDATE,omitempty"` // Date at which vesting percentage changes. Default (if empty) is that the vesting percentage below applies to the current date
	VestPct  Amount   `xml:"VESTPCT"`
}

type LoanInfo struct {
	XMLName               xml.Name `xml:"VESTINFO"`
	LoanID                String   `xml:"LOANID"`                          // Identifier of this loan
	LoanDesc              String   `xml:"LOANDESC,omitempty"`              // Loan description
	InitialLoanBal        Amount   `xml:"INITIALLOANBAL,omitempty"`        // Initial loan balance
	LoanStartDate         *Date    `xml:"LOANSTARTDATE,omitempty"`         // Start date of loan
	CurrentLoanBal        Amount   `xml:"CURRENTLOANBAL"`                  // Current loan principal balance
	DtAsOf                *Date    `xml:"DTASOF"`                          // Date and time of the current loan balance
	LoanRate              Amount   `xml:"LOANRATE,omitempty"`              // Loan annual interest rate
	LoanPmtAmt            Amount   `xml:"LOANPMTAMT,omitempty"`            // Loan payment amount
	LoanPmtFreq           String   `xml:"LOANPMTFREQ,omitempty"`           // Frequency of loan repayments: WEEKLY, BIWEEKLY, TWICEMONTHLY, MONTHLY, FOURWEEKS, BIMONTHLY, QUARTERLY, SEMIANNUALLY, ANNUALLY, OTHER. See section 10.2.1 for calculation rules.
	LoanPmtsInitial       Int      `xml:"LOANPMTSINITIAL,omitempty"`       // Initial number of loan payments.
	LoanPmtsRemaining     Int      `xml:"LOANPMTSREMAINING,omitempty"`     // Remaining number of loan payments
	LoanMaturityDate      *Date    `xml:"LOANMATURITYDATE,omitempty"`      // Expected loan end date
	LoanTotalProjInterest Amount   `xml:"LOANTOTALPROJINTEREST,omitempty"` // Total projected interest to be paid on this loan
	LoanInterestToDate    Amount   `xml:"LOANINTERESTTODATE,omitempty"`    // Total interested paid to date on this loan
	LoanExtPmtDate        *Date    `xml:"LOANNEXTPMTDATE,omitempty"`       // Next payment due date
}

type Inv401KSummaryAggregate struct {
	XMLName       xml.Name // One of CONTRIBUTIONS, WITHDRAWALS, EARNINGS
	PreTax        Amount   `xml:"PRETAX,omitempty"`        // Pretax withdrawals.
	AfterTax      Amount   `xml:"AFTERTAX,omitempty"`      // After tax withdrawals.
	Match         Amount   `xml:"MATCH,omitempty"`         // Employer matching withdrawals.
	ProfitSharing Amount   `xml:"PROFITSHARING,omitempty"` // Profit sharing withdrawals.
	Rollover      Amount   `xml:"ROLLOVER,omitempty"`      // Rollover withdrawals.
	OtherVest     Amount   `xml:"OTHERVEST,omitempty"`     // Other vesting withdrawals.
	OtherNonVest  Amount   `xml:"OTHERNONVEST,omitempty"`  // Other non-vesting withdrawals.
	Total         Amount   `xml:"TOTAL"`                   // Sum of withdrawals from all fund sources.
}

type Inv401KSummaryPeriod struct {
	XMLName       xml.Name                 // One of YEARTODATE, INCEPTODATE, or PERIODTODATE
	DtStart       Date                     `xml:"DTSTART"`
	DtEnd         Date                     `xml:"DTEND"`
	Contributions *Inv401KSummaryAggregate `xml:"CONTRIBUTIONS,omitempty"` // 401(k) contribution aggregate. Note: this includes loan payments.
	Withdrawls    *Inv401KSummaryAggregate `xml:"WITHDRAWLS,omitempty"`    // 401(k) withdrawals aggregate. Note: this includes loan withdrawals.
	Earnings      *Inv401KSummaryAggregate `xml:"EARNINGS,omitempty"`      // 401(k) earnings aggregate. This is the market value change. It includes dividends/interest, and capital gains - realized and unrealized.
}

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

type InvStatementResponse struct {
	XMLName   xml.Name `xml:"INVSTMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO OFXEXTENSION
	DtAsOf      Date         `xml:"INVSTMTRS>DTASOF"`
	CurDef      String       `xml:"INVSTMTRS>CURDEF"`
	InvAcctFrom InvAcct      `xml:"INVSTMTRS>INVACCTFROM"`
	InvTranList *InvTranList `xml:"INVSTMTRS>INVTRANLIST,omitempty"`
	InvPosList  PositionList `xml:"INVSTMTRS>INVPOSLIST,omitempty"`
	InvBal      *InvBalance  `xml:"INVSTMTRS>INVBAL,omitempty"`
	// TODO INVOOLIST
	MktgInfo   String      `xml:"INVSTMTRS>MKTGINFO,omitempty"` // Marketing information
	Inv401K    *Inv401K    `xml:"INVSTMTRS>INV401K,omitempty"`
	Inv401KBal *Inv401KBal `xml:"INVSTMTRS>INV401KBAL,omitempty"`
}

func (sr InvStatementResponse) Name() string {
	return "INVSTMTTRNRS"
}

func (sr InvStatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func decodeInvestmentsMessageSet(d *xml.Decoder, start xml.StartElement) ([]Message, error) {
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
