package ofxgo

import (
	"bytes"
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
	"time"
)

type Request struct {
	URL         string
	Version     string        // OFX version string, overwritten in Client.Request()
	Signon      SignonRequest //<SIGNONMSGSETV1>
	Signup      []Message     //<SIGNUPMSGSETV1>
	Banking     []Message     //<BANKMSGSETV1>
	CreditCards []Message     //<CREDITCARDMSGSETV1>
	//<LOANMSGSETV1>
	Investments []Message //<INVSTMTMSGSETV1>
	//<INTERXFERMSGSETV1>
	//<WIREXFERMSGSETV1>
	//<BILLPAYMSGSETV1>
	//<EMAILMSGSETV1>
	Securities []Message //<SECLISTMSGSETV1>
	//<PRESDIRMSGSETV1>
	//<PRESDLVMSGSETV1>
	Profile []Message //<PROFMSGSETV1>
	//<IMAGEMSGSETV1>

	indent bool // Whether to indent the marshaled XML
}

func marshalMessageSet(e *xml.Encoder, requests []Message, set messageType) error {
	if len(requests) > 0 {
		messageSetElement := xml.StartElement{Name: xml.Name{Local: set.String()}}
		if err := e.EncodeToken(messageSetElement); err != nil {
			return err
		}

		for _, request := range requests {
			if request.Type() != set {
				return errors.New("Expected " + set.String() + " message , found " + request.Type().String())
			}
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

// Overwrite the fields in this Request object controlled by the Client
func (oq *Request) SetClientFields(c *Client) {
	oq.Signon.DtClient = Date(time.Now())

	// Overwrite fields that the client controls
	oq.Version = c.OfxVersion()
	oq.Signon.AppId = c.Id()
	oq.Signon.AppVer = c.Version()
	oq.indent = c.IndentRequests()
}

// Marshal this Request into its SGML/XML representation held in a bytes.Buffer
//
// If error is non-nil, this bytes.Buffer is ready to be sent to an OFX server
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
	signonMsgSet := xml.StartElement{Name: xml.Name{Local: SignonRq.String()}}
	if err := encoder.EncodeToken(signonMsgSet); err != nil {
		return nil, err
	}
	if err := encoder.Encode(&oq.Signon); err != nil {
		return nil, err
	}
	if err := encoder.EncodeToken(signonMsgSet.End()); err != nil {
		return nil, err
	}

	if err := marshalMessageSet(encoder, oq.Signup, SignupRq); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Banking, BankRq); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.CreditCards, CreditCardRq); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Investments, InvStmtRq); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Securities, SecListRq); err != nil {
		return nil, err
	}
	if err := marshalMessageSet(encoder, oq.Profile, ProfileRq); err != nil {
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
