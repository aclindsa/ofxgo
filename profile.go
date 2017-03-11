package ofxgo

import (
	"github.com/golang/go/src/encoding/xml"
)

type OfxProfileRequest struct {
	XMLName       xml.Name  `xml:"PROFTRNRQ"`
	TrnUID        OfxUID    `xml:"TRNUID"`
	ClientRouting OfxString `xml:"PROFRQ>CLIENTROUTING"` // Forced to NONE
	DtProfup      OfxDate   `xml:"PROFRQ>DTPROFUP"`
}

func (r *OfxProfileRequest) Name() string {
	return "PROFTRNRQ"
}

func (r *OfxProfileRequest) Valid() (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	r.ClientRouting = "NONE"
	return true, nil
}
