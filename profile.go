package ofxgo

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

type ProfileRequest struct {
	XMLName       xml.Name `xml:"PROFTRNRQ"`
	TrnUID        UID      `xml:"TRNUID"`
	ClientRouting String   `xml:"PROFRQ>CLIENTROUTING"` // Forced to NONE
	DtProfUp      Date     `xml:"PROFRQ>DTPROFUP"`
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

type SignonInfo struct {
	XMLName           xml.Name `xml:"SIGNONINFO"`
	SignonRealm       String   `xml:"SIGNONREALM"`
	Min               Int      `xml:"MIN"`                      // Minimum number of password characters
	Max               Int      `xml:"MAX"`                      // Maximum number of password characters
	CharType          String   `xml:"CHARTYPE"`                 // ALPHAONLY, NUMERICONLY, ALPHAORNUMERIC, ALPHAANDNUMERIC
	CaseSen           Boolean  `xml:"CASESEN"`                  // Password is case-sensitive?
	Special           Boolean  `xml:"SPECIAL"`                  // Special characters allowed?
	Spaces            Boolean  `xml:"SPACES"`                   // Spaces allowed?
	Pinch             Boolean  `xml:"PINCH"`                    // Pin change <PINCHRQ> requests allowed
	ChgPinFirst       Boolean  `xml:"CHGPINFIRST"`              // Server requires user to change password at first signon
	UserCred1Label    String   `xml:"USERCRED1LABEL,omitempty"` // Prompt for USERCRED1 (if this field is present, USERCRED1 is required)
	UserCred2Label    String   `xml:"USERCRED2LABEL,omitempty"` // Prompt for USERCRED2 (if this field is present, USERCRED2 is required)
	ClientUIDReq      Boolean  `xml:"CLIENTUIDREQ,omitempty"`   // CLIENTUID required?
	AuthTokenFirst    Boolean  `xml:"AUTHTOKENFIRST,omitempty"` // Server requires AUTHTOKEN as part of first signon
	AuthTokenLabel    String   `xml:"AUTHTOKENLABEL,omitempty"`
	AuthTokenInfoURL  String   `xml:"AUTHTOKENINFOURL,omitempty"`
	MFAChallengeSupt  Boolean  `xml:"MFACHALLENGESUPT,omitempty"`  // Server supports MFACHALLENGE
	MFAChallengeFIRST Boolean  `xml:"MFACHALLENGEFIRST,omitempty"` // Server requires MFACHALLENGE to be sent with first signon
	AccessTokenReq    Boolean  `xml:"ACCESSTOKENREQ,omitempty"`    // Server requires ACCESSTOKEN to be sent with all requests except profile
}

type MessageSet struct {
	XMLName     xml.Name // <xxxMSGSETVn>
	Ver         String   `xml:"MSGSETCORE>VER"`
	Url         String   `xml:"MSGSETCORE>URL"`
	OfxSec      String   `xml:"MSGSETCORE>OFXSEC"`
	TranspSec   Boolean  `xml:"MSGSETCORE>TRANSPSEC"`
	SignonRealm String   `xml:"MSGSETCORE>SIGNONREALM"` // Used to identify which SignonInfo to use for to this MessageSet
	Language    []String `xml:"MSGSETCORE>LANGUAGE"`
	SyncMode    String   `xml:"MSGSETCORE>SYNCMODE"`
	// TODO MessageSet-specific stuff?
}

type MessageSetList []MessageSet

func (msl *MessageSetList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		var msgset MessageSet
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if _, ok := tok.(xml.StartElement); ok {
			// Found starting tag for <xxxMSGSET>. Get the next one (xxxMSGSETVn) and decode that struct
			tok, err := nextNonWhitespaceToken(d)
			if err != nil {
				return err
			} else if versionStart, ok := tok.(xml.StartElement); ok {
				if err := d.DecodeElement(&msgset, &versionStart); err != nil {
					return err
				}
			} else {
				return errors.New("Invalid MSGSETLIST formatting")
			}

			// Eat ending tags for <xxxMSGSET>
			tok, err = nextNonWhitespaceToken(d)
			if err != nil {
				return err
			} else if _, ok := tok.(xml.EndElement); !ok {
				return errors.New("Invalid MSGSETLIST formatting")
			}
		} else {
			return errors.New("MSGSETLIST didn't find an opening xxxMSGSETVn element")
		}
		*msl = MessageSetList(append(*(*[]MessageSet)(msl), msgset))
	}
}

type ProfileResponse struct {
	XMLName        xml.Name       `xml:"PROFTRNRS"`
	TrnUID         UID            `xml:"TRNUID"`
	MessageSetList MessageSetList `xml:"PROFRS>MSGSETLIST"`
	SignonInfoList []SignonInfo   `xml:"PROFRS>SIGNONINFOLIST>SIGNONINFO"`
	DtProfUp       Date           `xml:"PROFRS>DTPROFUP"`
	FiName         String         `xml:"PROFRS>FINAME"`
	Addr1          String         `xml:"PROFRS>ADDR1"`
	Addr2          String         `xml:"PROFRS>ADDR2,omitempty"`
	Addr3          String         `xml:"PROFRS>ADDR3,omitempty"`
	City           String         `xml:"PROFRS>CITY"`
	State          String         `xml:"PROFRS>STATE"`
	PostalCode     String         `xml:"PROFRS>POSTALCODE"`
	Country        String         `xml:"PROFRS>COUNTRY"`
	CsPhone        String         `xml:"PROFRS>CSPHONE,omitempty"`
	TsPhone        String         `xml:"PROFRS>TSPHONE,omitempty"`
	FaxPhone       String         `xml:"PROFRS>FAXPHONE,omitempty"`
	URL            String         `xml:"PROFRS>URL,omitempty"`
	Email          String         `xml:"PROFRS>EMAIL,omitempty"`
}

func (pr ProfileResponse) Name() string {
	return "PROFTRNRS"
}

func (pr ProfileResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func DecodeProfileMessageSet(d *xml.Decoder, start xml.StartElement) ([]Message, error) {
	var msgs []Message
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return nil, err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return msgs, nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "PROFTRNRS":
				var prof ProfileResponse
				if err := d.DecodeElement(&prof, &startElement); err != nil {
					return nil, err
				}
				msgs = append(msgs, Message(prof))
			default:
				return nil, errors.New("Unsupported profile response tag: " + startElement.Name.Local)
			}
		} else {
			return nil, errors.New("Didn't find an opening element")
		}
	}
}
