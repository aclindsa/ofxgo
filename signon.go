package ofxgo

import (
	"errors"
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
)

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

func (r *SignonRequest) Name() string {
	return "SONRQ"
}

func (r *SignonRequest) Valid() (bool, error) {
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
		return false, fmt.Errorf("SONRQ>LANGUAGE invalid length: \"%s\"\n", r.Language)
	}
	if len(r.AppID) < 1 || len(r.AppID) > 5 {
		return false, errors.New("SONRQ>APPID invalid length")
	}
	if len(r.AppVer) < 1 || len(r.AppVer) > 4 {
		return false, errors.New("SONRQ>APPVER invalid length")
	}
	return true, nil
}

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

func (r *SignonResponse) Name() string {
	return "SONRS"
}

func (r *SignonResponse) Valid() (bool, error) {
	if len(r.Language) != 3 {
		return false, fmt.Errorf("SONRS>LANGUAGE invalid length: \"%s\"\n", r.Language)
	}
	return r.Status.Valid()
}
