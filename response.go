package ofxgo

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
	"io"
	"strings"
)

type Response struct {
	Version     string         // String for OFX header, defaults to 203
	Signon      SignonResponse //<SIGNONMSGSETV1>
	Signup      []Message      //<SIGNUPMSGSETV1>
	Banking     []Message      //<BANKMSGSETV1>
	CreditCards []Message      //<CREDITCARDMSGSETV1>
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
	var tok xml.Token
	tok, err := nextNonWhitespaceToken(decoder)
	if err != nil {
		return err
	} else if xmlElem, ok := tok.(xml.ProcInst); !ok || xmlElem.Target != "xml" {
		return errors.New("Missing xml processing instruction")
	}

	// parse the OFX header
	tok, err = nextNonWhitespaceToken(decoder)
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

// Number of bytes of response to read when attempting to figure out whether
// we're using OFX or SGML
const guessVersionCheckBytes = 1024

// Defaults to XML if it can't determine the version or if there is any
// ambiguity
func guessVersion(r *bufio.Reader) (bool, error) {
	b, _ := r.Peek(guessVersionCheckBytes)
	if b == nil {
		return false, errors.New("Failed to read OFX header")
	}
	sgmlIndex := bytes.Index(b, []byte("OFXHEADER:"))
	xmlIndex := bytes.Index(b, []byte("OFXHEADER="))
	if sgmlIndex < 0 {
		return true, nil
	} else if xmlIndex < 0 {
		return false, nil
	} else {
		return xmlIndex <= sgmlIndex, nil
	}
}

// ParseResponse parses an OFX response in SGML or XML into a Response object
// from the given io.Reader
//
// It is commonly used as part of Client.Request(), but may be used on its own
// to parse already-downloaded OFX files (such as those from 'Web Connect'). It
// performs version autodetection if it can and attempts to be as forgiving as
// possible about the input format.
func ParseResponse(reader io.Reader) (*Response, error) {
	var or Response

	r := bufio.NewReaderSize(reader, guessVersionCheckBytes)
	xmlVersion, err := guessVersion(r)
	if err != nil {
		return nil, err
	}

	// parse SGML headers before creating XML decoder
	if !xmlVersion {
		if err := or.readSGMLHeaders(r); err != nil {
			return nil, err
		}
	}

	decoder := xml.NewDecoder(r)
	if !xmlVersion {
		decoder.Strict = false
		decoder.AutoCloseAfterCharData = ofxLeafElements
	}
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}

	if xmlVersion {
		// parse the xml header
		if err := or.readXMLHeaders(decoder); err != nil {
			return nil, err
		}
	}

	tok, err := nextNonWhitespaceToken(decoder)
	if err != nil {
		return nil, err
	} else if ofxStart, ok := tok.(xml.StartElement); !ok || ofxStart.Name.Local != "OFX" {
		return nil, errors.New("Missing opening OFX xml element")
	}

	// Unmarshal the signon message
	tok, err = nextNonWhitespaceToken(decoder)
	if err != nil {
		return nil, err
	} else if signonStart, ok := tok.(xml.StartElement); ok && signonStart.Name.Local == "SIGNONMSGSRSV1" {
		if err := decoder.Decode(&or.Signon); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Missing opening SIGNONMSGSRSV1 xml element")
	}

	tok, err = nextNonWhitespaceToken(decoder)
	if err != nil {
		return nil, err
	} else if signonEnd, ok := tok.(xml.EndElement); !ok || signonEnd.Name.Local != "SIGNONMSGSRSV1" {
		return nil, errors.New("Missing closing SIGNONMSGSRSV1 xml element")
	}
	if ok, err := or.Signon.Valid(); !ok {
		return nil, err
	}

	for {
		tok, err = nextNonWhitespaceToken(decoder)
		if err != nil {
			return nil, err
		} else if ofxEnd, ok := tok.(xml.EndElement); ok && ofxEnd.Name.Local == "OFX" {
			return &or, nil // found closing XML element, so we're done
		} else if start, ok := tok.(xml.StartElement); ok {
			// TODO decode other types
			switch start.Name.Local {
			case "SIGNUPMSGSRSV1":
				msgs, err := DecodeSignupMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.Signup = msgs
			case "BANKMSGSRSV1":
				msgs, err := DecodeBankingMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.Banking = msgs
			case "CREDITCARDMSGSRSV1":
				msgs, err := DecodeCCMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.CreditCards = msgs
			//case "LOANMSGSRSV1":
			case "INVSTMTMSGSRSV1":
				msgs, err := DecodeInvestmentsMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.Investments = msgs
			//case "INTERXFERMSGSRSV1":
			//case "WIREXFERMSGSRSV1":
			//case "BILLPAYMSGSRSV1":
			//case "EMAILMSGSRSV1":
			case "SECLISTMSGSRSV1":
				msgs, err := DecodeSecuritiesMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.Securities = msgs
			//case "PRESDIRMSGSRSV1":
			//case "PRESDLVMSGSRSV1":
			case "PROFMSGSRSV1":
				msgs, err := DecodeProfileMessageSet(decoder, start)
				if err != nil {
					return nil, err
				}
				or.Profile = msgs
			//case "IMAGEMSGSRSV1":
			default:
				return nil, errors.New("Unsupported message set: " + start.Name.Local)
			}
		} else {
			return nil, errors.New("Found unexpected token")
		}
	}
}
