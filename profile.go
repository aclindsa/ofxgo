package ofxgo

import (
	"github.com/golang/go/src/encoding/xml"
)

type ProfileRequest struct {
	XMLName       xml.Name `xml:"PROFTRNRQ"`
	TrnUID        UID      `xml:"TRNUID"`
	ClientRouting String   `xml:"PROFRQ>CLIENTROUTING"` // Forced to NONE
	DtProfup      Date     `xml:"PROFRQ>DTPROFUP"`
}

func (r *ProfileRequest) Name() string {
	return "PROFTRNRQ"
}

func (r *ProfileRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	r.ClientRouting = "NONE"
	return true, nil
}
