package ofxgo

import (
	"github.com/golang/go/src/encoding/xml"
)

type OfxAcctInfoRequest struct {
	XMLName   xml.Name `xml:"ACCTINFOTRNRQ"`
	TrnUID    OfxUID   `xml:"TRNUID"`
	CltCookie OfxInt   `xml:"CLTCOOKIE"`
	DtAcctup  OfxDate  `xml:"ACCTINFORQ>DTACCTUP"`
}

func (r *OfxAcctInfoRequest) Name() string {
	return "ACCTINFOTRNRQ"
}

func (r *OfxAcctInfoRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	return true, nil
}
