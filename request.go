package ofxgo

import (
	"bytes"
	"errors"
	"github.com/aclindsa/xml"
	"time"
)

// Request is the top-level object marshalled and sent to OFX servers. It is
// constructed by appending one or more request objects to the message set they
// correspond to (i.e. appending StatementRequest to Request.Bank to get a bank
// statemement). If a *Request object is appended to the wrong message set, an
// error will be returned when Marshal() is called on this Request.
type Request struct {
	URL        string
	Version    ofxVersion    // OFX version, overwritten in Client.Request()
	Signon     SignonRequest //<SIGNONMSGSETV1>
	Signup     []Message     //<SIGNUPMSGSETV1>
	Bank       []Message     //<BANKMSGSETV1>
	CreditCard []Message     //<CREDITCARDMSGSETV1>
	Loan       []Message     //<LOANMSGSETV1>
	InvStmt    []Message     //<INVSTMTMSGSETV1>
	InterXfer  []Message     //<INTERXFERMSGSETV1>
	WireXfer   []Message     //<WIREXFERMSGSETV1>
	Billpay    []Message     //<BILLPAYMSGSETV1>
	Email      []Message     //<EMAILMSGSETV1>
	SecList    []Message     //<SECLISTMSGSETV1>
	PresDir    []Message     //<PRESDIRMSGSETV1>
	PresDlv    []Message     //<PRESDLVMSGSETV1>
	Prof       []Message     //<PROFMSGSETV1>
	Image      []Message     //<IMAGEMSGSETV1>

	indent         bool // Whether to indent the marshaled XML
	carriageReturn bool // Whether to user carriage returns in new lines for marshaled XML
}

func encodeMessageSet(e *xml.Encoder, requests []Message, set messageType, version ofxVersion) error {
	if len(requests) > 0 {
		messageSetElement := xml.StartElement{Name: xml.Name{Local: set.String()}}
		if err := e.EncodeToken(messageSetElement); err != nil {
			return err
		}

		for _, request := range requests {
			if request.Type() != set {
				return errors.New("Expected " + set.String() + " message , found " + request.Type().String())
			}
			if ok, err := request.Valid(version); !ok {
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

// SetClientFields overwrites the fields in this Request object controlled by
// the Client
func (oq *Request) SetClientFields(c Client) {
	oq.Signon.DtClient.Time = time.Now()

	// Overwrite fields that the client controls
	oq.Version = c.OfxVersion()
	oq.Signon.AppID = c.ID()
	oq.Signon.AppVer = c.Version()
	oq.indent = c.IndentRequests()
	oq.carriageReturn = c.CarriageReturnNewLines()
}

// Marshal this Request into its SGML/XML representation held in a bytes.Buffer
//
// If error is non-nil, this bytes.Buffer is ready to be sent to an OFX server
func (oq *Request) Marshal() (*bytes.Buffer, error) {
	var b bytes.Buffer

	// Write the header appropriate to our version
	writeHeader(&b, oq.Version, oq.carriageReturn)

	encoder := xml.NewEncoder(&b)
	if oq.indent {
		encoder.Indent("", "    ")
	}
	if oq.carriageReturn {
		encoder.CarriageReturn(true)
	}

	ofxElement := xml.StartElement{Name: xml.Name{Local: "OFX"}}

	if err := encoder.EncodeToken(ofxElement); err != nil {
		return nil, err
	}

	if ok, err := oq.Signon.Valid(oq.Version); !ok {
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

	messageSets := []struct {
		Messages []Message
		Type     messageType
	}{
		{oq.Signup, SignupRq},
		{oq.Bank, BankRq},
		{oq.CreditCard, CreditCardRq},
		{oq.Loan, LoanRq},
		{oq.InvStmt, InvStmtRq},
		{oq.InterXfer, InterXferRq},
		{oq.WireXfer, WireXferRq},
		{oq.Billpay, BillpayRq},
		{oq.Email, EmailRq},
		{oq.SecList, SecListRq},
		{oq.PresDir, PresDirRq},
		{oq.PresDlv, PresDlvRq},
		{oq.Prof, ProfRq},
		{oq.Image, ImageRq},
	}
	for _, set := range messageSets {
		if err := encodeMessageSet(encoder, set.Messages, set.Type, oq.Version); err != nil {
			return nil, err
		}
	}

	if err := encoder.EncodeToken(ofxElement.End()); err != nil {
		return nil, err
	}

	if err := encoder.Flush(); err != nil {
		return nil, err
	}
	return &b, nil
}
