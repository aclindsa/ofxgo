package ofxgo

import (
	"errors"
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
)

// SignonRequest identifies and authenticates a user to their FI and is
// provided with every Request
type SignonRequest struct {
	XMLName   xml.Name `xml:"SONRQ"`
	DtClient  Date     `xml:"DTCLIENT"` // Current time on client, overwritten in Client.Request()
	UserID    String   `xml:"USERID"`
	UserPass  String   `xml:"USERPASS,omitempty"`
	UserKey   String   `xml:"USERKEY,omitempty"`
	Language  String   `xml:"LANGUAGE"` // Defaults to ENG
	Org       String   `xml:"FI>ORG"`
	Fid       String   `xml:"FI>FID"`
	AppID     String   `xml:"APPID"`  // Overwritten in Client.Request()
	AppVer    String   `xml:"APPVER"` // Overwritten in Client.Request()
	ClientUID UID      `xml:"CLIENTUID,omitempty"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *SignonRequest) Name() string {
	return "SONRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *SignonRequest) Valid(version ofxVersion) (bool, error) {
	if len(r.UserID) < 1 || len(r.UserID) > 32 {
		return false, errors.New("SONRQ>USERID invalid length")
	}
	if (len(r.UserPass) == 0) == (len(r.UserKey) == 0) {
		return false, errors.New("One and only one of SONRQ>USERPASS and USERKEY must be supplied")
	}
	if len(r.UserPass) > 32 {
		return false, errors.New("SONRQ>USERPASS invalid length")
	}
	if len(r.UserKey) > 64 {
		return false, errors.New("SONRQ>USERKEY invalid length")
	}
	if len(r.Language) == 0 {
		r.Language = "ENG"
	} else if len(r.Language) != 3 {
		return false, fmt.Errorf("SONRQ>LANGUAGE invalid length: \"%s\"", r.Language)
	}
	if len(r.AppID) < 1 || len(r.AppID) > 5 {
		return false, errors.New("SONRQ>APPID invalid length")
	}
	if len(r.AppVer) < 1 || len(r.AppVer) > 4 {
		return false, errors.New("SONRQ>APPVER invalid length")
	}
	return true, nil
}

// SignonResponse is provided with every Response and indicates the success or
// failure of the SignonRequest in the corresponding Request
type SignonResponse struct {
	XMLName     xml.Name `xml:"SONRS"`
	Status      Status   `xml:"STATUS"`
	DtServer    Date     `xml:"DTSERVER"`
	UserKey     String   `xml:"USERKEY,omitempty"`
	TsKeyExpire *Date    `xml:"TSKEYEXPIRE,omitempty"`
	Language    String   `xml:"LANGUAGE"`
	DtProfUp    *Date    `xml:"DTPROFUP,omitempty"`
	DtAcctUp    *Date    `xml:"DTACCTUP,omitempty"`
	Org         String   `xml:"FI>ORG"`
	Fid         String   `xml:"FI>FID"`
	SessCookie  String   `xml:"SESSCOOKIE,omitempty"`
	AccessKey   String   `xml:"ACCESSKEY,omitempty"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *SignonResponse) Name() string {
	return "SONRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (r *SignonResponse) Valid(version ofxVersion) (bool, error) {
	if len(r.Language) != 3 {
		return false, fmt.Errorf("SONRS>LANGUAGE invalid length: \"%s\"", r.Language)
	}
	return r.Status.Valid()
}
