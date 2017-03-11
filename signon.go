package ofxgo

import (
	"errors"
	"github.com/golang/go/src/encoding/xml"
)

type OfxSignonRequest struct {
	XMLName   xml.Name  `xml:"SONRQ"`
	Dtclient  OfxDate   `xml:"DTCLIENT"` // Overridden in OfxRequest.Request()
	UserId    OfxString `xml:"USERID"`
	UserPass  OfxString `xml:"USERPASS,omitempty"`
	UserKey   OfxString `xml:"USERKEY,omitempty"`
	Language  OfxString `xml:"LANGUAGE"` // Defaults to ENG
	Org       OfxString `xml:"FI>ORG"`
	Fid       OfxString `xml:"FI>FID"`
	AppId     OfxString `xml:"APPID"`  // Defaults to OFXGO
	AppVer    OfxString `xml:"APPVER"` // Defaults to 0001
	ClientUID OfxUID    `xml:"CLIENTUID,omitempty"`
}

func (r *OfxSignonRequest) Name() string {
	return "SONRQ"
}

func (r *OfxSignonRequest) Valid() (bool, error) {
	if len(r.UserId) < 1 || len(r.UserId) > 32 {
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
		return false, errors.New("SONRQ>LANGUAGE invalid length")
	}
	if len(r.AppId) == 0 {
		r.AppId = "OFXGO"
	} else if len(r.AppId) > 5 {
		return false, errors.New("SONRQ>APPID invalid length")
	}
	if len(r.AppVer) == 0 {
		r.AppVer = "0001"
	} else if len(r.AppVer) > 4 {
		return false, errors.New("SONRQ>APPVER invalid length")
	}
	if ok, err := r.ClientUID.Valid(); !ok {
		if len(r.ClientUID) > 0 { // ClientUID isn't required
			return false, err
		}
	}
	return true, nil
}

type OfxStatus struct {
	XMLName  xml.Name  `xml:"STATUS"`
	Code     OfxInt    `xml:"CODE"`
	Severity OfxString `xml:"SEVERITY"`
	Message  OfxString `xml:"MESSAGE,omitempty"`
}

func (s *OfxStatus) Valid() (bool, error) {
	switch s.Severity {
	case "INFO", "WARN", "ERROR":
		return true, nil
	default:
		return false, errors.New("Invalid STATUS>SEVERITY")
	}
}

type OfxSignonResponse struct {
	XMLName     xml.Name  `xml:"SONRS"`
	Status      OfxStatus `xml:"STATUS"`
	Dtserver    OfxDate   `xml:"DTSERVER"`
	UserKey     OfxString `xml:"USERKEY,omitempty"`
	TsKeyExpire OfxDate   `xml:"TSKEYEXPIRE,omitempty"`
	Language    OfxString `xml:"LANGUAGE"`
	Dtprofup    OfxDate   `xml:"DTPROFUP,omitempty"`
	Dtacctup    OfxDate   `xml:"DTACCTUP,omitempty"`
	Org         OfxString `xml:"FI>ORG"`
	Fid         OfxString `xml:"FI>FID"`
	SessCookie  OfxString `xml:"SESSCOOKIE,omitempty"`
	AccessKey   OfxString `xml:"ACCESSKEY,omitempty"`
}

func (r *OfxSignonResponse) Name() string {
	return "SONRS"
}

func (r *OfxSignonResponse) Valid() (bool, error) {
	if len(r.Language) != 3 {
		return false, errors.New("SONRS>LANGUAGE invalid length: " + string(r.Language))
	}
	return r.Status.Valid()
}
