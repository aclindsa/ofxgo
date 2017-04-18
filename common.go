package ofxgo

//go:generate ./generate_constants.py

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

// Message represents an OFX message in a message set. it is used to ease
// marshalling and unmarshalling.
type Message interface {
	Name() string                           // The name of the OFX transaction wrapper element this represents
	Valid(version ofxVersion) (bool, error) // Called before a Message is marshaled and after it's unmarshaled to ensure the request or response is valid
	Type() messageType                      // The message set this message belongs to
}

type messageType uint

// These constants are returned by Messages' Type() functions to determine
// which message set they belong to
const (
	// Requests
	SignonRq messageType = iota
	SignupRq
	BankRq
	CreditCardRq
	LoanRq
	InvStmtRq
	InterXferRq
	WireXferRq
	BillpayRq
	EmailRq
	SecListRq
	PresDirRq
	PresDlvRq
	ProfRq
	ImageRq

	// Responses
	SignonRs
	SignupRs
	BankRs
	CreditCardRs
	LoanRs
	InvStmtRs
	InterXferRs
	WireXferRs
	BillpayRs
	EmailRs
	SecListRs
	PresDirRs
	PresDlvRs
	ProfRs
	ImageRs
)

func (t messageType) String() string {
	switch t {
	case SignonRq:
		return "SIGNONMSGSRQV1"
	case SignupRq:
		return "SIGNUPMSGSRQV1"
	case BankRq:
		return "BANKMSGSRQV1"
	case CreditCardRq:
		return "CREDITCARDMSGSRQV1"
	case LoanRq:
		return "LOANMSGSRQV1"
	case InvStmtRq:
		return "INVSTMTMSGSRQV1"
	case InterXferRq:
		return "INTERXFERMSGSRQV1"
	case WireXferRq:
		return "WIREXFERMSGSRQV1"
	case BillpayRq:
		return "BILLPAYMSGSRQV1"
	case EmailRq:
		return "EMAILMSGSRQV1"
	case SecListRq:
		return "SECLISTMSGSRQV1"
	case PresDirRq:
		return "PRESDIRMSGSRQV1"
	case PresDlvRq:
		return "PRESDLVMSGSRQV1"
	case ProfRq:
		return "PROFMSGSRQV1"
	case ImageRq:
		return "IMAGEMSGSRQV1"
	case SignonRs:
		return "SIGNONMSGSRSV1"
	case SignupRs:
		return "SIGNUPMSGSRSV1"
	case BankRs:
		return "BANKMSGSRSV1"
	case CreditCardRs:
		return "CREDITCARDMSGSRSV1"
	case LoanRs:
		return "LOANMSGSRSV1"
	case InvStmtRs:
		return "INVSTMTMSGSRSV1"
	case InterXferRs:
		return "INTERXFERMSGSRSV1"
	case WireXferRs:
		return "WIREXFERMSGSRSV1"
	case BillpayRs:
		return "BILLPAYMSGSRSV1"
	case EmailRs:
		return "EMAILMSGSRSV1"
	case SecListRs:
		return "SECLISTMSGSRSV1"
	case PresDirRs:
		return "PRESDIRMSGSRSV1"
	case PresDlvRs:
		return "PRESDLVMSGSRSV1"
	case ProfRs:
		return "PROFMSGSRSV1"
	case ImageRs:
		return "IMAGEMSGSRSV1"
	}
	panic("Invalid messageType")
}

// Map of error codes to their meanings, SEVERITY, and conditions under which
// OFX servers are expected to return them
var statusMeanings = map[Int][3]string{
	0:     {"Success", "INFO", "The server successfully processed the request."},
	1:     {"Client is up-to-date", "INFO", "Based on the client timestamp, the client has the latest information. The response does not supply any additional information."},
	2000:  {"General error", "ERROR", "Error other than those specified by the remaining error codes. Note: Servers should provide a more specific error whenever possible. Error code 2000 should be reserved for cases in which a more specific code is not available."},
	2001:  {"Invalid account", "ERROR", ""},
	2002:  {"General account error", "ERROR", "Account error not specified by the remaining error codes."},
	2003:  {"Account not found", "ERROR", "The specified account number does not correspond to one of the user’s accounts."},
	2004:  {"Account closed", "ERROR", "The specified account number corresponds to an account that has been closed."},
	2005:  {"Account not authorized", "ERROR", "The user is not authorized to perform this action on the account, or the server does not allow this type of action to be performed on the account."},
	2006:  {"Source account not found", "ERROR", "The specified account number does not correspond to one of the user’s accounts."},
	2007:  {"Source account closed", "ERROR", "The specified account number corresponds to an account that has been closed."},
	2008:  {"Source account not authorized", "ERROR", "The user is not authorized to perform this action on the account, or the server does not allow this type of action to be performed on the account."},
	2009:  {"Destination account not found", "ERROR", "The specified account number does not correspond to one of the user’s accounts."},
	2010:  {"Destination account closed", "ERROR", "The specified account number corresponds to an account that has been closed."},
	2011:  {"Destination account not authorized", "ERROR", "The user is not authorized to perform this action on the account, or the server does not allow this type of action to be performed on the account."},
	2012:  {"Invalid amount", "ERROR", "The specified amount is not valid for this action; for example, the user specified a negative payment amount."},
	2014:  {"Date too soon", "ERROR", "The server cannot process the requested action by the date specified by the user."},
	2015:  {"Date too far in future", "ERROR", "The server cannot accept requests for an action that far in the future."},
	2016:  {"Transaction already committed", "ERROR", "Transaction has entered the processing loop and cannot be modified/cancelled using OFX.  The transaction may still be cancelled or modified using other means (for example, a phone call to Customer Service)."},
	2017:  {"Already canceled", "ERROR", "The transaction cannot be canceled or modified because it has already been canceled."},
	2018:  {"Unknown server ID", "ERROR", "The specified server ID does not exist or no longer exists."},
	2019:  {"Duplicate request", "ERROR", "A request with this <TRNUID> has already been received and processed."},
	2020:  {"Invalid date", "ERROR", "The specified datetime stamp cannot be parsed; for instance, the datetime stamp specifies 25:00 hours."},
	2021:  {"Unsupported version", "ERROR", "The server does not support the requested version. The version of the message set specified by the client is not supported by this server."},
	2022:  {"Invalid TAN", "ERROR", "The server was unable to validate the TAN sent in the request."},
	2023:  {"Unknown FITID", "ERROR", "The specified FITID/BILLID does not exist or no longer exists.  [BILLID not found (ERROR) in the billing message sets]"},
	2025:  {"Branch ID missing", "ERROR", "A <BRANCHID> value must be provided in the <BANKACCTFROM> aggregate for this country system, but this field is missing."},
	2026:  {"Bank name doesn’t match bank ID", "ERROR", "The value of <BANKNAME> in the <EXTBANKACCTTO> aggregate is inconsistent with the value of <BANKID> in the <BANKACCTTO> aggregate."},
	2027:  {"Invalid date range", "ERROR", "Response for non-overlapping dates, date ranges in the future, et cetera."},
	2028:  {"Requested element unknown", "WARN", "One or more elements of the request were not recognized by the server or the server (as noted in the FI Profile) does not support the elements. The server executed the element transactions it understood and supported. For example, the request file included private tags in a <PMTRQ> but the server was able to execute the rest of the request."},
	3000:  {"MFA Challenge authentication required", "ERROR", "User credentials are correct, but further authentication required. Client should send <MFACHALLENGERQ> in next request."},
	3001:  {"MFA Challenge information is invalid", "ERROR", "User or client information sent in MFACHALLENGEA contains invalid information"},
	6500:  {"<REJECTIFMISSING>Y invalid without <TOKEN>", "ERROR", "This error code may appear in the <SYNCERROR> element of an <xxxSYNCRS> wrapper (in <PRESDLVMSGSRSV1> and V2 message set responses) or the <CODE> contained in any embedded transaction wrappers within a sync response. The corresponding sync request wrapper included <REJECTIFMISSING>Y with <REFRESH>Y or <TOKENONLY>Y, which is illegal."},
	6501:  {"Embedded transactions in request failed to process: Out of date", "WARN", "<REJECTIFMISSING>Y and embedded transactions appeared in the request sync wrapper and the provided <TOKEN> was out of date. This code should be used in the <SYNCERROR> of the response sync wrapper."},
	6502:  {"Unable to process embedded transaction due to out-of-date <TOKEN>", "ERROR", "Used in response transaction wrapper for embedded transactions when <SYNCERROR>6501 appears in the surrounding sync wrapper."},
	10000: {"Stop check in process", "INFO", "Stop check is already in process."},
	10500: {"Too many checks to process", "ERROR", "The stop-payment request <STPCHKRQ> specifies too many checks."},
	10501: {"Invalid payee", "ERROR", "Payee error not specified by the remaining error codes."},
	10502: {"Invalid payee address", "ERROR", "Some portion of the payee’s address is incorrect or unknown."},
	10503: {"Invalid payee account number", "ERROR", "The account number <PAYACCT> of the requested payee is invalid."},
	10504: {"Insufficient funds", "ERROR", "The server cannot process the request because the specified account does not have enough funds."},
	10505: {"Cannot modify element", "ERROR", "The server does not allow modifications to one or more values in a modification request."},
	10506: {"Cannot modify source account", "ERROR", "Reserved for future use."},
	10507: {"Cannot modify destination account", "ERROR", "Reserved for future use."},
	10508: {"Invalid frequency", "ERROR", "The specified frequency <FREQ> does not match one of the accepted frequencies for recurring transactions."},
	10509: {"Model already canceled", "ERROR", "The server has already canceled the specified recurring model."},
	10510: {"Invalid payee ID", "ERROR", "The specified payee ID does not exist or no longer exists."},
	10511: {"Invalid payee city", "ERROR", "The specified city is incorrect or unknown."},
	10512: {"Invalid payee state", "ERROR", "The specified state is incorrect or unknown."},
	10513: {"Invalid payee postal code", "ERROR", "The specified postal code is incorrect or unknown."},
	10514: {"Transaction already processed", "ERROR", "Transaction has already been sent or date due is past"},
	10515: {"Payee not modifiable by client", "ERROR", "The server does not allow clients to change payee information."},
	10516: {"Wire beneficiary invalid", "ERROR", "The specified wire beneficiary does not exist or no longer exists."},
	10517: {"Invalid payee name", "ERROR", "The server does not recognize the specified payee name."},
	10518: {"Unknown model ID", "ERROR", "The specified model ID does not exist or no longer exists."},
	10519: {"Invalid payee list ID", "ERROR", "The specified payee list ID does not exist or no longer exists."},
	10600: {"Table type not found", "ERROR", "The specified table type is not recognized or does not exist."},
	12250: {"Investment transaction download not supported", "WARN", "The server does not support investment transaction download."},
	12251: {"Investment position download not supported", "WARN", "The server does not support investment position download."},
	12252: {"Investment positions for specified date not available", "WARN", "The server does not support investment positions for the specified date."},
	12253: {"Investment open order download not supported", "WARN", "The server does not support open order download."},
	12254: {"Investment balances download not supported", "WARN", "The server does not support investment balances download."},
	12255: {"401(k) not available for this account", "ERROR", "401(k) information requested from a non- 401(k) account."},
	12500: {"One or more securities not found", "ERROR", "The server could not find the requested securities."},
	13000: {"User ID & password will be sent out-of-band", "INFO", "The server will send the user ID and password via postal mail, e-mail, or another means. The accompanying message will provide details."},
	13500: {"Unable to enroll user", "ERROR", "The server could not enroll the user."},
	13501: {"User already enrolled", "ERROR", "The server has already enrolled the user."},
	13502: {"Invalid service", "ERROR", "The server does not support the service <SVC> specified in the service-activation request."},
	13503: {"Cannot change user information", "ERROR", "The server does not support the <CHGUSERINFORQ> request."},
	13504: {"<FI> Missing or Invalid in <SONRQ>", "ERROR", "The FI requires the client to provide the <FI> aggregate in the <SONRQ> request, but either none was provided, or the one provided was invalid."},
	14500: {"1099 forms not available", "ERROR", "1099 forms are not yet available for the tax year requested."},
	14501: {"1099 forms not available for user ID", "ERROR", "This user does not have any 1099 forms available."},
	14600: {"W2 forms not available", "ERROR", "W2 forms are not yet available for the tax year requested."},
	14601: {"W2 forms not available for user ID", "ERROR", "The user does not have any W2 forms available."},
	14700: {"1098 forms not available", "ERROR", "1098 forms are not yet available for the tax year requested."},
	14701: {"1098 forms not available for user ID", "ERROR", "The user does not have any 1098 forms available."},
	15000: {"Must change USERPASS", "INFO", "The user must change his or her <USERPASS> number as part of the next OFX request."},
	15500: {"Signon invalid", "ERROR", "The user cannot signon because he or she entered an invalid user ID or password."},
	15501: {"Customer account already in use", "ERROR", "The server allows only one connection at a time, and another user is already signed on.  Please try again later."},
	15502: {"USERPASS lockout", "ERROR", "The server has received too many failed signon attempts for this user. Please call the FI’s technical support number."},
	15503: {"Could not change USERPASS", "ERROR", "The server does not support the <PINCHRQ> request."},
	15504: {"Could not provide random data", "ERROR", "The server could not generate random data as requested by the <CHALLENGERQ>."},
	15505: {"Country system not supported", "ERROR", "The server does not support the country specified in the <COUNTRY> field of the <SONRQ> aggregate."},
	15506: {"Empty signon not supported", "ERROR", "The server does not support signons not accompanied by some other transaction."},
	15507: {"Signon invalid without supporting pin change request", "ERROR", "The OFX block associated with the signon does not contain a pin change request and should."},
	15508: {"Transaction not authorized. ", "ERROR", "Current user is not authorized to perform this action on behalf of the <USERID>."},
	15510: {"CLIENTUID error", "ERROR", "The CLIENTUID sent by the client was incorrect. User must register the Client UID."},
	15511: {"MFA error", "ERROR", "User should contact financial institution."},
	15512: {"AUTHTOKEN required", "ERROR", "User needs to contact financial institution to obtain AUTHTOKEN. Client should send it in the next request."},
	15513: {"AUTHTOKEN invalid", "ERROR", "The AUTHTOKEN sent by the client was invalid."},
	16500: {"HTML not allowed", "ERROR", "The server does not accept HTML formatting in the request."},
	16501: {"Unknown mail To:", "ERROR", "The server was unable to send mail to the specified Internet address."},
	16502: {"Invalid URL", "ERROR", "The server could not parse the URL."},
	16503: {"Unable to get URL", "ERROR", "The server was unable to retrieve the information at this URL (e.g., an HTTP 400 or 500 series error)."},
}

// Status represents the status of a Response (both top-level Request objects,
// and *Response objects)
type Status struct {
	XMLName  xml.Name `xml:"STATUS"`
	Code     Int      `xml:"CODE"`
	Severity String   `xml:"SEVERITY"`
	Message  String   `xml:"MESSAGE,omitempty"`
}

// Valid returns whether the Status is valid according to the OFX spec
func (s *Status) Valid() (bool, error) {
	switch s.Severity {
	case "INFO", "WARN", "ERROR":
	default:
		return false, errors.New("Invalid STATUS>SEVERITY")
	}

	if arr, ok := statusMeanings[s.Code]; ok {
		if arr[1] != string(s.Severity) {
			return false, errors.New("Unexpected SEVERITY for STATUS>CODE")
		}
	} else {
		return false, errors.New("Unknown OFX status code")
	}

	return true, nil
}

// CodeMeaning returns the meaning of the current status Code
func (s *Status) CodeMeaning() (string, error) {
	if arr, ok := statusMeanings[s.Code]; ok {
		return arr[0], nil
	}
	return "", errors.New("Unknown OFX status code")
}

// CodeConditions returns the conditions under which an OFX server is expected
// to return the current status Code
func (s *Status) CodeConditions() (string, error) {
	if arr, ok := statusMeanings[s.Code]; ok {
		return arr[2], nil
	}
	return "", errors.New("Unknown OFX status code")
}

// BankAcct represents the identifying information for one bank account
type BankAcct struct {
	XMLName  xml.Name // BANKACCTTO or BANKACCTFROM
	BankID   String   `xml:"BANKID"`
	BranchID String   `xml:"BRANCHID,omitempty"` // Unused in USA
	AcctID   String   `xml:"ACCTID"`
	AcctType acctType `xml:"ACCTTYPE"`          // One of CHECKING, SAVINGS, MONEYMRKT, CREDITLINE, CD
	AcctKey  String   `xml:"ACCTKEY,omitempty"` // Unused in USA
}

// Valid returns whether the BankAcct is valid according to the OFX spec
func (b BankAcct) Valid() (bool, error) {
	if len(b.BankID) == 0 {
		return false, errors.New("BankAcct.BankID empty")
	}
	if len(b.AcctID) == 0 {
		return false, errors.New("BankAcct.AcctID empty")
	}
	if !b.AcctType.Valid() {
		return false, errors.New("Invalid or unspecified BankAcct.AcctType")
	}
	return true, nil
}

// CCAcct represents the identifying information for one checking account
type CCAcct struct {
	XMLName xml.Name // CCACCTTO or CCACCTFROM
	AcctID  String   `xml:"ACCTID"`
	AcctKey String   `xml:"ACCTKEY,omitempty"` // Unused in USA
}

// Valid returns whether the CCAcct is valid according to the OFX spec
func (c CCAcct) Valid() (bool, error) {
	if len(c.AcctID) == 0 {
		return false, errors.New("CCAcct.AcctID empty")
	}
	return true, nil
}

// InvAcct represents the identifying information for one investment account
type InvAcct struct {
	XMLName  xml.Name // INVACCTTO or INVACCTFROM
	BrokerID String   `xml:"BROKERID"`
	AcctID   String   `xml:"ACCTID"`
}

// Currency represents one ISO-4217 currency
type Currency struct {
	XMLName xml.Name   // CURRENCY or ORIGCURRENCY
	CurRate Amount     `xml:"CURRATE"` // Ratio of <CURDEF> currency to <CURSYM> currency
	CurSym  CurrSymbol `xml:"CURSYM"`  // ISO-4217 3-character currency identifier
}
