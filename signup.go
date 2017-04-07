package ofxgo

import (
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
)

type AcctInfoRequest struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	DtAcctUp Date `xml:"ACCTINFORQ>DTACCTUP"`
}

func (r *AcctInfoRequest) Name() string {
	return "ACCTINFOTRNRQ"
}

func (r *AcctInfoRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

func (r *AcctInfoRequest) Type() messageType {
	return SignupRq
}

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

// Make pointers to these structs print nicely
func (bai *BankAcctInfo) String() string {
	return fmt.Sprintf("%+v", *bai)
}

type CCAcctInfo struct {
	XMLName            xml.Name           `xml:"CCACCTINFO"`
	CCAcctFrom         CCAcct             `xml:"CCACCTFROM"`
	SupTxDl            Boolean            `xml:"SUPTXDL"`                      // Supports downloading transactions (as opposed to balance only)
	XferSrc            Boolean            `xml:"XFERSRC"`                      // Enabled as source for intra/interbank transfer
	XferDest           Boolean            `xml:"XFERDEST"`                     // Enabled as destination for intra/interbank transfer
	AcctClassification acctClassification `xml:"ACCTCLASSIFICATION,omitempty"` // One of PERSONAL, BUSINESS, CORPORATE, OTHER
	SvcStatus          svcStatus          `xml:"SVCSTATUS"`                    // One of AVAIL (available, but not yet requested), PEND (requested, but not yet available), ACTIVE
}

// Make pointers to these structs print nicely
func (ci *CCAcctInfo) String() string {
	return fmt.Sprintf("%+v", *ci)
}

type InvAcctInfo struct {
	XMLName       xml.Name      `xml:"INVACCTINFO"`
	InvAcctFrom   InvAcct       `xml:"INVACCTFROM"`
	UsProductType usProductType `xml:"USPRODUCTTYPE"`         // One of 401K, 403B, IRA, KEOGH, OTHER, SARSEP, SIMPLE, NORMAL, TDA, TRUST, UGMA
	Checking      Boolean       `xml:"CHECKING"`              // Has check-writing privileges
	SvcStatus     svcStatus     `xml:"SVCSTATUS"`             // One of AVAIL (available, but not yet requested), PEND (requested, but not yet available), ACTIVE
	InvAcctType   holderType    `xml:"INVACCTTYPE,omitempty"` // One of INDIVIDUAL, JOINT, TRUST, CORPORATE
	OptionLevel   String        `xml:"OPTIONLEVEL,omitempty"` // Text desribing option trading privileges
}

// Make pointers to these structs print nicely
func (iai *InvAcctInfo) String() string {
	return fmt.Sprintf("%+v", *iai)
}

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

type AcctInfoResponse struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	DtAcctUp Date       `xml:"ACCTINFORS>DTACCTUP"`
	AcctInfo []AcctInfo `xml:"ACCTINFORS>ACCTINFO,omitempty"`
}

func (air *AcctInfoResponse) Name() string {
	return "ACCTINFOTRNRS"
}

func (air *AcctInfoResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func (air *AcctInfoResponse) Type() messageType {
	return SignupRs
}
