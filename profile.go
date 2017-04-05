package ofxgo

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

type ProfileRequest struct {
	XMLName   xml.Name `xml:"PROFTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	ClientRouting String `xml:"PROFRQ>CLIENTROUTING"` // Forced to NONE
	DtProfUp      Date   `xml:"PROFRQ>DTPROFUP"`
}

func (r *ProfileRequest) Name() string {
	return "PROFTRNRQ"
}

func (r *ProfileRequest) Valid() (bool, error) {
	// TODO implement
	r.ClientRouting = "NONE"
	return true, nil
}

func (r *ProfileRequest) Type() messageType {
	return ProfRq
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
	PinCh             Boolean  `xml:"PINCH"`                    // Pin change <PINCHRQ> requests allowed
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
	Name        string   // <xxxMSGSETVn> (copy of XMLName.Local)
	Ver         Int      `xml:"MSGSETCORE>VER"`                   // Message set version - should always match 'n' in <xxxMSGSETVn>
	Url         String   `xml:"MSGSETCORE>URL"`                   // URL where messages in this set are to be set
	OfxSec      String   `xml:"MSGSETCORE>OFXSEC"`                // NONE or 'TYPE 1'
	TranspSec   Boolean  `xml:"MSGSETCORE>TRANSPSEC"`             // Transport-level security must be used
	SignonRealm String   `xml:"MSGSETCORE>SIGNONREALM"`           // Used to identify which SignonInfo to use for to this MessageSet
	Language    []String `xml:"MSGSETCORE>LANGUAGE"`              // List of supported languages
	SyncMode    String   `xml:"MSGSETCORE>SYNCMODE"`              // One of FULL, LITE
	RefreshSupt Boolean  `xml:"MSGSETCORE>REFRESHSUPT,omitempty"` // Y if server supports <REFRESH>Y within synchronizations. This option is irrelevant for full synchronization servers. Clients must ignore <REFRESHSUPT> (or its absence) if the profile also specifies <SYNCMODE>FULL. For lite synchronization, the default is N. Without <REFRESHSUPT>Y, lite synchronization servers are not required to support <REFRESH>Y requests
	RespFileER  Boolean  `xml:"MSGSETCORE>RESPFILEER"`            // server supports file-based error recovery
	SpName      String   `xml:"MSGSETCORE>SPNAME"`                // Name of service provider
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
			msgset.Name = msgset.XMLName.Local

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
	XMLName   xml.Name `xml:"PROFTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
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

func (pr *ProfileResponse) Name() string {
	return "PROFTRNRS"
}

func (pr *ProfileResponse) Valid() (bool, error) {
	//TODO implement
	return true, nil
}

func (pr *ProfileResponse) Type() messageType {
	return ProfRs
}
