package ofxgo

import (
	"bytes"
	"errors"
	"github.com/golang/go/src/encoding/xml"
)

type Request struct {
	URL     string
	Version string        // OFX version string, overwritten in Client.Request()
	Signon  SignonRequest //<SIGNONMSGSETV1>
	Signup  []Message     //<SIGNUPMSGSETV1>
	Banking []Message     //<BANKMSGSETV1>
	//<CREDITCARDMSGSETV1>
	//<LOANMSGSETV1>
	//<INVSTMTMSGSETV1>
	//<INTERXFERMSGSETV1>
	//<WIREXFERMSGSETV1>
	//<BILLPAYMSGSETV1>
	//<EMAILMSGSETV1>
	//<SECLISTMSGSETV1>
	//<PRESDIRMSGSETV1>
	//<PRESDLVMSGSETV1>
	Profile []Message //<PROFMSGSETV1>
	//<IMAGEMSGSETV1>

	indent bool // Whether to indent the marshaled XML
}

func marshalMessageSet(e *xml.Encoder, requests []Message, setname string) error {
	if len(requests) > 0 {
		messageSetElement := xml.StartElement{Name: xml.Name{Local: setname}}
		if err := e.EncodeToken(messageSetElement); err != nil {
			return err
		}

		for _, request := range requests {
			if ok, err := request.Valid(); !ok {
				return err
			}
			if err := e.Encode(request); err != nil {
				return err
			}
		}

		if err := e.EncodeToken(messageSetElement.End()); err != nil {
			return err
		}
	}
	return nil
}

func (oq *Request) Marshal() (*bytes.Buffer, error) {
	var b bytes.Buffer

	// Write the header appropriate to our version
	switch oq.Version {
	case "102", "103", "151", "160":
		b.WriteString(`OFXHEADER:100
DATA:OFXSGML
VERSION:` + oq.Version + `
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

`)
	case "200", "201", "202", "203", "210", "211", "220":
		b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\n")
		b.WriteString(`<?OFX OFXHEADER="200" VERSION="` + oq.Version + `" SECURITY="NONE" OLDFILEUID="NONE" NEWFILEUID="NONE"?>` + "\n")
	default:
		return nil, errors.New(oq.Version + " is not a valid OFX version string")
	}

	encoder := xml.NewEncoder(&b)
	if oq.indent {
		encoder.Indent("", "    ")
	}

	ofxElement := xml.StartElement{Name: xml.Name{Local: "OFX"}}

	if err := encoder.EncodeToken(ofxElement); err != nil {
		return nil, err
	}

	if ok, err := oq.Signon.Valid(); !ok {
		return nil, err
	}
	signonMsgSet := xml.StartElement{Name: xml.Name{Local: "SIGNONMSGSRQV1"}}
	if err := encoder.EncodeToken(signonMsgSet); err != nil {
		return nil, err
	}
	if err := encoder.Encode(&oq.Signon); err != nil {
		return nil, err
	}
	if err := encoder.EncodeToken(signonMsgSet.End()); err != nil {
		return nil, err
	}

	if err := marshalMessageSet(encoder, oq.Signup, "SIGNUPMSGSRQV1"); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Banking, "BANKMSGSRQV1"); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Profile, "PROFMSGSRQV1"); err != nil {
		return nil, err
	}

	if err := encoder.EncodeToken(ofxElement.End()); err != nil {
		return nil, err
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}
	return &b, nil
}
