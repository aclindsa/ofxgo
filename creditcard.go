package ofxgo

import (
	"github.com/aclindsa/go/src/encoding/xml"
)

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

func (r *CCStatementRequest) Name() string {
	return "CCSTMTTRNRQ"
}

func (r *CCStatementRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

func (r *CCStatementRequest) Type() messageType {
	return CreditCardRq
}

type CCStatementResponse struct {
	XMLName   xml.Name `xml:"CCSTMTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	CurDef       String           `xml:"CCSTMTRS>CURDEF"`
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

func (sr *CCStatementResponse) Name() string {
	return "CCSTMTTRNRS"
}

func (sr *CCStatementResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func (sr *CCStatementResponse) Type() messageType {
	return CreditCardRs
}
