package ofxgo

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/aclindsa/xml"
	"io"
	"reflect"
	"strings"
)

// Response is the top-level object returned from a parsed OFX response file.
// It can be inspected by using type assertions or switches on the message set
// you're interested in.
type Response struct {
	Version    ofxVersion     // OFX header version
	Signon     SignonResponse //<SIGNONMSGSETV1>
	Signup     []Message      //<SIGNUPMSGSETV1>
	Bank       []Message      //<BANKMSGSETV1>
	CreditCard []Message      //<CREDITCARDMSGSETV1>
	Loan       []Message      //<LOANMSGSETV1>
	InvStmt    []Message      //<INVSTMTMSGSETV1>
	InterXfer  []Message      //<INTERXFERMSGSETV1>
	WireXfer   []Message      //<WIREXFERMSGSETV1>
	Billpay    []Message      //<BILLPAYMSGSETV1>
	Email      []Message      //<EMAILMSGSETV1>
	SecList    []Message      //<SECLISTMSGSETV1>
	PresDir    []Message      //<PRESDIRMSGSETV1>
	PresDlv    []Message      //<PRESDLVMSGSETV1>
	Prof       []Message      //<PROFMSGSETV1>
	Image      []Message      //<IMAGEMSGSETV1>
}

func (or *Response) readSGMLHeaders(r *bufio.Reader) error {
	var seenHeader, seenVersion bool = false, false
	for {
		// Some financial institutions do not properly leave an empty line after the last header.
		// Avoid attempting to read another header in that case.
		next, err := r.Peek(1)
		if err != nil {
			return err
		}
		if next[0] == '<' {
			break
		}

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

		// Some OFX servers put a space after the colon
		headervalue := strings.TrimSpace(header[1])

		switch header[0] {
		case "OFXHEADER":
			if headervalue != "100" {
				return errors.New("OFXHEADER is not 100")
			}
			seenHeader = true
		case "DATA":
			if headervalue != "OFXSGML" {
				return errors.New("OFX DATA header does not contain OFXSGML")
			}
		case "VERSION":
			err := or.Version.FromString(headervalue)
			if err != nil {
				return err
			}
			seenVersion = true

			if or.Version > OfxVersion160 {
				return errors.New("OFX VERSION > 160 in SGML header")
			}
		case "SECURITY":
			if headervalue != "NONE" {
				return errors.New("OFX SECURITY header not NONE")
			}
		case "COMPRESSION":
			if headervalue != "NONE" {
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
				err := or.Version.FromString(value)
				if err != nil {
					return err
				}
				seenVersion = true

				if or.Version < OfxVersion200 {
					return errors.New("OFX VERSION < 200 in XML header")
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

// A map of message set tags to a map of transaction wrapper tags to the
// reflect.Type of the struct for that transaction type. Used when decoding
// Responses. Newly-implemented response transaction types *must* be added to
// this map in order to be unmarshalled.
var responseTypes = map[string]map[string]reflect.Type{
	SignupRs.String(): {
		(&AcctInfoResponse{}).Name(): reflect.TypeOf(AcctInfoResponse{})},
	BankRs.String(): {
		(&StatementResponse{}).Name(): reflect.TypeOf(StatementResponse{})},
	CreditCardRs.String(): {
		(&CCStatementResponse{}).Name(): reflect.TypeOf(CCStatementResponse{})},
	LoanRs.String(): {},
	InvStmtRs.String(): {
		(&InvStatementResponse{}).Name(): reflect.TypeOf(InvStatementResponse{})},
	InterXferRs.String(): {},
	WireXferRs.String():  {},
	BillpayRs.String():   {},
	EmailRs.String():     {},
	SecListRs.String(): {
		(&SecListResponse{}).Name(): reflect.TypeOf(SecListResponse{}),
		(&SecurityList{}).Name():    reflect.TypeOf(SecurityList{})},
	PresDirRs.String(): {},
	PresDlvRs.String(): {},
	ProfRs.String(): {
		(&ProfileResponse{}).Name(): reflect.TypeOf(ProfileResponse{})},
	ImageRs.String(): {},
}

func decodeMessageSet(d *xml.Decoder, start xml.StartElement, msgs *[]Message, version ofxVersion) error {
	setTypes, ok := responseTypes[start.Name.Local]
	if !ok {
		return errors.New("Invalid message set: " + start.Name.Local)
	}
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			responseType, ok := setTypes[startElement.Name.Local]
			if !ok {
				// If you are a developer and received this message after you
				// thought you added a new transaction type, make sure you
				// added it to the responseTypes map above
				return errors.New("Unsupported response transaction for " +
					start.Name.Local + ": " + startElement.Name.Local)
			}
			response := reflect.New(responseType).Interface()
			responseMessage := response.(Message)
			if err := d.DecodeElement(responseMessage, &startElement); err != nil {
				return err
			}
			if ok, err := responseMessage.Valid(version); !ok {
				return err
			}
			*msgs = append(*msgs, responseMessage)
		} else {
			return errors.New("Didn't find an opening element")
		}
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
	} else if signonStart, ok := tok.(xml.StartElement); ok && signonStart.Name.Local == SignonRs.String() {
		if err := decoder.Decode(&or.Signon); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Missing opening SIGNONMSGSRSV1 xml element")
	}

	tok, err = nextNonWhitespaceToken(decoder)
	if err != nil {
		return nil, err
	} else if signonEnd, ok := tok.(xml.EndElement); !ok || signonEnd.Name.Local != SignonRs.String() {
		return nil, errors.New("Missing closing SIGNONMSGSRSV1 xml element")
	}
	if ok, err := or.Signon.Valid(or.Version); !ok {
		return nil, err
	}

	var messageSlices = map[string]*[]Message{
		SignupRs.String():     &or.Signup,
		BankRs.String():       &or.Bank,
		CreditCardRs.String(): &or.CreditCard,
		LoanRs.String():       &or.Loan,
		InvStmtRs.String():    &or.InvStmt,
		InterXferRs.String():  &or.InterXfer,
		WireXferRs.String():   &or.WireXfer,
		BillpayRs.String():    &or.Billpay,
		EmailRs.String():      &or.Email,
		SecListRs.String():    &or.SecList,
		PresDirRs.String():    &or.PresDir,
		PresDlvRs.String():    &or.PresDlv,
		ProfRs.String():       &or.Prof,
		ImageRs.String():      &or.Image,
	}

	for {
		tok, err = nextNonWhitespaceToken(decoder)
		if err != nil {
			return nil, err
		} else if ofxEnd, ok := tok.(xml.EndElement); ok && ofxEnd.Name.Local == "OFX" {
			return &or, nil // found closing XML element, so we're done
		} else if start, ok := tok.(xml.StartElement); ok {
			slice, ok := messageSlices[start.Name.Local]
			if !ok {
				return nil, errors.New("Invalid message set: " + start.Name.Local)
			}
			if err := decodeMessageSet(decoder, start, slice, or.Version); err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("Found unexpected token")
		}
	}
}
