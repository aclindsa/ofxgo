package ofxgo

import (
	"github.com/aclindsa/xml"
)

// CCStatementRequest represents a request for a credit card statement. It is
// used to request balances and/or transactions. See StatementRequest for the
// analog for all other bank accounts.
type CCStatementRequest struct {
	XMLName   xml.Name `xml:"CCSTMTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"`
	// TODO OFXEXTENSION
	CCAcctFrom     CCAcct  `xml:"CCSTMTRQ>CCACCTFROM"`
	DtStart        *Date   `xml:"CCSTMTRQ>INCTRAN>DTSTART,omitempty"`
	DtEnd          *Date   `xml:"CCSTMTRQ>INCTRAN>DTEND,omitempty"`
	Include        Boolean `xml:"CCSTMTRQ>INCTRAN>INCLUDE"`          // Include transactions (instead of just balance)
	IncludePending Boolean `xml:"CCSTMTRQ>INCLUDEPENDING,omitempty"` // Include pending transactions
	IncTranImg     Boolean `xml:"CCSTMTRQ>INCTRANIMG,omitempty"`     // Include transaction images
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *CCStatementRequest) Name() string {
	return "CCSTMTTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *CCStatementRequest) Valid(version ofxVersion) (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *CCStatementRequest) Type() messageType {
	return CreditCardRq
}

// CCStatementResponse represents a credit card statement, including its
// balances and possibly transactions. It is a response to CCStatementRequest,
// or sometimes provided as part of an OFX file downloaded manually from an FI.
type CCStatementResponse struct {
	XMLName   xml.Name `xml:"CCSTMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	CurDef       CurrSymbol       `xml:"CCSTMTRS>CURDEF"`
	CCAcctFrom   CCAcct           `xml:"CCSTMTRS>CCACCTFROM"`
	BankTranList *TransactionList `xml:"CCSTMTRS>BANKTRANLIST,omitempty"`
	//BANKTRANLISTP
	BalAmt        Amount    `xml:"CCSTMTRS>LEDGERBAL>BALAMT"`
	DtAsOf        Date      `xml:"CCSTMTRS>LEDGERBAL>DTASOF"`
	AvailBalAmt   *Amount   `xml:"CCSTMTRS>AVAILBAL>BALAMT,omitempty"`
	AvailDtAsOf   *Date     `xml:"CCSTMTRS>AVAILBAL>DTASOF,omitempty"`
	CashAdvBalAmt Amount    `xml:"CCSTMTRS>CASHADVBALAMT,omitempty"`           // Only for CREDITLINE accounts, available balance for cash advances
	IntRatePurch  Amount    `xml:"CCSTMTRS>INTRATEPURCH,omitempty"`            // Current interest rate for purchases
	IntRateCash   Amount    `xml:"CCSTMTRS>INTRATECASH,omitempty"`             // Current interest rate for cash advances
	IntRateXfer   Amount    `xml:"CCSTMTRS>INTRATEXFER,omitempty"`             // Current interest rate for cash advances
	RewardName    String    `xml:"CCSTMTRS>REWARDINFO>NAME,omitempty"`         // Name of the reward program referred to by the next two elements
	RewardBal     Amount    `xml:"CCSTMTRS>REWARDINFO>REWARDBAL,omitempty"`    // Current balance of the reward program
	RewardEarned  Amount    `xml:"CCSTMTRS>REWARDINFO>REWARDEARNED,omitempty"` // Reward amount earned YTD
	BalList       []Balance `xml:"CCSTMTRS>BALLIST>BAL,omitempty"`
	MktgInfo      String    `xml:"CCSTMTRS>MKTGINFO,omitempty"` // Marketing information
}

// Name returns the name of the top-level transaction XML/SGML element
func (sr *CCStatementResponse) Name() string {
	return "CCSTMTTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (sr *CCStatementResponse) Valid(version ofxVersion) (bool, error) {
	if ok, err := sr.TrnUID.Valid(); !ok {
		return false, err
	}
	//TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (sr *CCStatementResponse) Type() messageType {
	return CreditCardRs
}
