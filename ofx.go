package ofxgo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/go/src/encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"
)

type Message interface {
	Name() string
	Valid() (bool, error)
}

type Request struct {
	URL     string
	Version string        // String for OFX header, defaults to 203
	Signon  SignonRequest //<SIGNONMSGSETV1>
	Signup  []Message     //<SIGNUPMSGSETV1>
	//<BANKMSGSETV1>
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
}

func (oq *Request) marshalMessageSet(e *xml.Encoder, requests []Message, setname string) error {
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

	if len(oq.Version) == 0 {
		oq.Version = "203"
	}

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
	encoder.Indent("", "    ")

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

	if err := oq.marshalMessageSet(encoder, oq.Signup, "SIGNUPMSGSRQV1"); err != nil {
		return nil, err
	}
	if err := oq.marshalMessageSet(encoder, oq.Profile, "PROFMSGSRQV1"); err != nil {
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

func (oq *Request) Request() (*Response, error) {
	oq.Signon.Dtclient = Date(time.Now())

	b, err := oq.Marshal()
	if err != nil {
		return nil, err
	}
	response, err := http.Post(oq.URL, "application/x-ofx", b)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("OFXQuery request status: " + response.Status)
	}

	// Help the parser out by giving it a clue about what header format to
	// expect
	var xmlVersion bool = true
	switch oq.Version {
	case "102", "103", "151", "160":
		xmlVersion = false
	}

	var ofxresp Response
	if err := ofxresp.Unmarshal(response.Body, xmlVersion); err != nil {
		return nil, err
	}

	return &ofxresp, nil
}

type Response struct {
	Version string         // String for OFX header, defaults to 203
	Signon  SignonResponse //<SIGNONMSGSETV1>
	Signup  []Message      //<SIGNUPMSGSETV1>
	//<BANKMSGSETV1>
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
}

func (or *Response) readSGMLHeaders(r *bufio.Reader) error {
	var seenHeader, seenVersion bool = false, false
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		// r.ReadString leaves the '\n' on the end...
		line = strings.TrimSpace(line)

		if len(line) == 0 {
			if seenHeader {
				break
			} else {
				continue
			}
		}
		header := strings.SplitN(line, ":", 2)
		if header == nil || len(header) != 2 {
			return errors.New("OFX headers malformed")
		}

		switch header[0] {
		case "OFXHEADER":
			if header[1] != "100" {
				return errors.New("OFXHEADER is not 100")
			}
			seenHeader = true
		case "DATA":
			if header[1] != "OFXSGML" {
				return errors.New("OFX DATA header does not contain OFXSGML")
			}
		case "VERSION":
			switch header[1] {
			case "102", "103", "151", "160":
				seenVersion = true
				or.Version = header[1]
			default:
				return errors.New("Invalid OFX VERSION in header")
			}
		case "SECURITY":
			if header[1] != "NONE" {
				return errors.New("OFX SECURITY header not NONE")
			}
		case "COMPRESSION":
			if header[1] != "NONE" {
				return errors.New("OFX COMPRESSION header not NONE")
			}
		case "ENCODING", "CHARSET", "OLDFILEUID", "NEWFILEUID":
			// TODO check/handle these headers?
		default:
			return errors.New("Invalid OFX header: " + header[0])
		}
	}

	if !seenVersion {
		return errors.New("OFX VERSION header missing")
	}
	return nil
}

func (or *Response) readXMLHeaders(decoder *xml.Decoder) error {
	tok, err := decoder.Token()
	if err != nil {
		return err
	} else if xmlElem, ok := tok.(xml.ProcInst); !ok || xmlElem.Target != "xml" {
		return errors.New("Missing xml processing instruction")
	}

	// parse the OFX header
	tok, err = decoder.Token()
	if err != nil {
		return err
	} else if ofxElem, ok := tok.(xml.ProcInst); ok && ofxElem.Target == "OFX" {
		var seenHeader, seenVersion bool = false, false

		headers := bytes.TrimSpace(ofxElem.Inst)
		for len(headers) > 0 {
			tmp := bytes.SplitN(headers, []byte("=\""), 2)
			if len(tmp) != 2 {
				return errors.New("Malformed OFX header")
			}
			header := string(tmp[0])
			headers = tmp[1]
			tmp = bytes.SplitN(headers, []byte("\""), 2)
			if len(tmp) != 2 {
				return errors.New("Malformed OFX header")
			}
			value := string(tmp[0])
			headers = bytes.TrimSpace(tmp[1])

			switch header {
			case "OFXHEADER":
				if value != "200" {
					return errors.New("OFXHEADER is not 200")
				}
				seenHeader = true
			case "VERSION":
				switch value {
				case "200", "201", "202", "203", "210", "211", "220":
					seenVersion = true
					or.Version = value
				default:
					return errors.New("Invalid OFX VERSION in header")
				}
			case "SECURITY":
				if value != "NONE" {
					return errors.New("OFX SECURITY header not NONE")
				}
			case "OLDFILEUID", "NEWFILEUID":
				// TODO check/handle these headers?
			default:
				return errors.New("Invalid OFX header: " + header)
			}
		}

		if !seenHeader {
			return errors.New("OFXHEADER version missing")
		}
		if !seenVersion {
			return errors.New("OFX VERSION header missing")
		}

	} else {
		return errors.New("Missing xml 'OFX' processing instruction")
	}
	return nil
}

func (or *Response) Unmarshal(reader io.Reader, xmlVersion bool) error {
	r := bufio.NewReader(reader)

	// parse SGML headers before creating XML decoder
	if !xmlVersion {
		if err := or.readSGMLHeaders(r); err != nil {
			return err
		}
	}

	decoder := xml.NewDecoder(r)
	if !xmlVersion {
		decoder.Strict = false
		decoder.AutoCloseAfterCharData = ofxLeafElements
	}

	if xmlVersion {
		// parse the xml header
		if err := or.readXMLHeaders(decoder); err != nil {
			return err
		}
	}

	tok, err := decoder.Token()
	if err != nil {
		return err
	} else if ofxStart, ok := tok.(xml.StartElement); !ok || ofxStart.Name.Local != "OFX" {
		return errors.New("Missing opening OFX xml element")
	}

	// Unmarshal the signon message
	tok, err = decoder.Token()
	if err != nil {
		return err
	} else if signonStart, ok := tok.(xml.StartElement); ok && signonStart.Name.Local == "SIGNONMSGSRSV1" {
		if err := decoder.Decode(&or.Signon); err != nil {
			return err
		}
	} else {
		return errors.New("Missing opening SIGNONMSGSRSV1 xml element")
	}

	tok, err = decoder.Token()
	if err != nil {
		return err
	} else if signonEnd, ok := tok.(xml.EndElement); !ok || signonEnd.Name.Local != "SIGNONMSGSRSV1" {
		return errors.New("Missing closing SIGNONMSGSRSV1 xml element")
	}
	if ok, err := or.Signon.Valid(); !ok {
		return err
	}

	for {
		tok, err = decoder.Token()
		if err != nil {
			return err
		} else if ofxEnd, ok := tok.(xml.EndElement); ok && ofxEnd.Name.Local == "OFX" {
			return nil // found closing XML element, so we're done
		} else if start, ok := tok.(xml.StartElement); ok {
			// TODO decode other types
			fmt.Println("Found starting element for: " + start.Name.Local)
		} else {
			return errors.New("Found unexpected token")
		}

		decoder.Skip()
	}
}
