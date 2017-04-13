package ofxgo

import (
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
)

// AcctInfoRequest represents a request for the server to provide information
// for all of the user's available accounts at this FI
type AcctInfoRequest struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	DtAcctUp Date `xml:"ACCTINFORQ>DTACCTUP"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *AcctInfoRequest) Name() string {
	return "ACCTINFOTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *AcctInfoRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *AcctInfoRequest) Type() messageType {
	return SignupRq
}

// HolderInfo contains the information a FI has about an account-holder
type HolderInfo struct {
	XMLName    xml.Name
	FirstName  String     `xml:"FIRSTNAME"`
	MiddleName String     `xml:"MIDDLENAME,omitempty"`
	LastName   String     `xml:"LASTNAME"`
	Addr1      String     `xml:"ADDR1"`
	Addr2      String     `xml:"ADDR2,omitempty"`
	Addr3      String     `xml:"ADDR3,omitempty"`
	City       String     `xml:"CITY"`
	State      String     `xml:"STATE"`
	PostalCode String     `xml:"POSTALCODE"`
	Country    String     `xml:"COUNTRY,omitempty"`
	DayPhone   String     `xml:"DAYPHONE,omitempty"`
	EvePhone   String     `xml:"EVEPHONE,omitempty"`
	Email      String     `xml:"EMAIL,omitempty"`
	HolderType holderType `xml:"HOLDERTYPE,omitempty"` // One of INDIVIDUAL, JOINT, CUSTODIAL, TRUST, OTHER
}

// BankAcctInfo contains information about a bank account, including how to
// access it (BankAcct), and whether it supports downloading transactions
// (SupTxDl).
type BankAcctInfo struct {
	XMLName            xml.Name           `xml:"BANKACCTINFO"`
	BankAcctFrom       BankAcct           `xml:"BANKACCTFROM"`
	SupTxDl            Boolean            `xml:"SUPTXDL"`                      // Supports downloading transactions (as opposed to balance only)
	XferSrc            Boolean            `xml:"XFERSRC"`                      // Enabled as source for intra/interbank transfer
	XferDest           Boolean            `xml:"XFERDEST"`                     // Enabled as destination for intra/interbank transfer
	MaturityDate       Date               `xml:"MATURITYDATE,omitempty"`       // Maturity date for CD, if CD
	MaturityAmt        Amount             `xml:"MATURITYAMOUNT,omitempty"`     // Maturity amount for CD, if CD
	MinBalReq          Amount             `xml:"MINBALREQ,omitempty"`          // Minimum balance required to avoid service fees
	AcctClassification acctClassification `xml:"ACCTCLASSIFICATION,omitempty"` // One of PERSONAL, BUSINESS, CORPORATE, OTHER
	OverdraftLimit     Amount             `xml:"OVERDRAFTLIMIT,omitempty"`
	SvcStatus          svcStatus          `xml:"SVCSTATUS"` // One of AVAIL (available, but not yet requested), PEND (requested, but not yet available), ACTIVE
}

// String makes pointers to BankAcctInfo structs print nicely
func (bai *BankAcctInfo) String() string {
	return fmt.Sprintf("%+v", *bai)
}

// CCAcctInfo contains information about a credit card account, including how
// to access it (CCAcct), and whether it supports downloading transactions
// (SupTxDl).
type CCAcctInfo struct {
	XMLName            xml.Name           `xml:"CCACCTINFO"`
	CCAcctFrom         CCAcct             `xml:"CCACCTFROM"`
	SupTxDl            Boolean            `xml:"SUPTXDL"`                      // Supports downloading transactions (as opposed to balance only)
	XferSrc            Boolean            `xml:"XFERSRC"`                      // Enabled as source for intra/interbank transfer
	XferDest           Boolean            `xml:"XFERDEST"`                     // Enabled as destination for intra/interbank transfer
	AcctClassification acctClassification `xml:"ACCTCLASSIFICATION,omitempty"` // One of PERSONAL, BUSINESS, CORPORATE, OTHER
	SvcStatus          svcStatus          `xml:"SVCSTATUS"`                    // One of AVAIL (available, but not yet requested), PEND (requested, but not yet available), ACTIVE
}

// String makes pointers to CCAcctInfo structs print nicely
func (ci *CCAcctInfo) String() string {
	return fmt.Sprintf("%+v", *ci)
}

// InvAcctInfo contains information about an investment account, including how
// to access it (InvAcct), and whether it supports downloading transactions
// (SupTxDl).
type InvAcctInfo struct {
	XMLName       xml.Name      `xml:"INVACCTINFO"`
	InvAcctFrom   InvAcct       `xml:"INVACCTFROM"`
	UsProductType usProductType `xml:"USPRODUCTTYPE"`         // One of 401K, 403B, IRA, KEOGH, OTHER, SARSEP, SIMPLE, NORMAL, TDA, TRUST, UGMA
	Checking      Boolean       `xml:"CHECKING"`              // Has check-writing privileges
	SvcStatus     svcStatus     `xml:"SVCSTATUS"`             // One of AVAIL (available, but not yet requested), PEND (requested, but not yet available), ACTIVE
	InvAcctType   holderType    `xml:"INVACCTTYPE,omitempty"` // One of INDIVIDUAL, JOINT, TRUST, CORPORATE
	OptionLevel   String        `xml:"OPTIONLEVEL,omitempty"` // Text desribing option trading privileges
}

// String makes pointers to InvAcctInfo structs print nicely
func (iai *InvAcctInfo) String() string {
	return fmt.Sprintf("%+v", *iai)
}

// AcctInfo represents generic account information. It should contain one (and
// only one) *AcctInfo element corresponding to the tyep of account it
// represents.
type AcctInfo struct {
	XMLName         xml.Name   `xml:"ACCTINFO"`
	Name            String     `xml:"NAME,omitempty"`
	Desc            String     `xml:"DESC,omitempty"`
	Phone           String     `xml:"PHONE,omitempty"`
	PrimaryHolder   HolderInfo `xml:"HOLDERINFO>PRIMARYHOLDER,omitempty"`
	SecondaryHolder HolderInfo `xml:"HOLDERINFO>SECONDARYHOLDER,omitempty"`

	// Only one of the rest of the fields will be valid for any given AcctInfo
	BankAcctInfo *BankAcctInfo `xml:"BANKACCTINFO,omitempty"`
	CCAcctInfo   *CCAcctInfo   `xml:"CCACCTINFO,omitempty"`
	InvAcctInfo  *InvAcctInfo  `xml:"INVACCTINFO,omitempty"`
	// TODO LOANACCTINFO
	// TODO BPACCTINFO?
}

// AcctInfoResponse contains the information about all a user's accounts
// accessible from this FI
type AcctInfoResponse struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	DtAcctUp Date       `xml:"ACCTINFORS>DTACCTUP"`
	AcctInfo []AcctInfo `xml:"ACCTINFORS>ACCTINFO,omitempty"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (air *AcctInfoResponse) Name() string {
	return "ACCTINFOTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (air *AcctInfoResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (air *AcctInfoResponse) Type() messageType {
	return SignupRs
}
