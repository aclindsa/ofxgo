package ofxgo

import (
	"github.com/golang/go/src/encoding/xml"
)

type AcctInfoRequest struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie Int      `xml:"CLTCOOKIE"`
	DtAcctup  Date     `xml:"ACCTINFORQ>DTACCTUP"`
}

func (r *AcctInfoRequest) Name() string {
	return "ACCTINFOTRNRQ"
}

func (r *AcctInfoRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	return true, nil
}
