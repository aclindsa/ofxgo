package ofxgo

import (
	"errors"
	"github.com/aclindsa/xml"
)

// ProfileRequest represents a request for a server to provide a profile of its
// capabilities (which message sets and versions it supports, how to access
// them, which languages and which types of synchronization they support, etc.)
type ProfileRequest struct {
	XMLName   xml.Name `xml:"PROFTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	ClientRouting String `xml:"PROFRQ>CLIENTROUTING"` // Forced to NONE
	DtProfUp      Date   `xml:"PROFRQ>DTPROFUP"`      // Date and time client last received a profile update
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *ProfileRequest) Name() string {
	return "PROFTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *ProfileRequest) Valid(version ofxVersion) (bool, error) {
	if ok, err := r.TrnUID.Valid(); !ok {
		return false, err
	}
	// TODO implement
	r.ClientRouting = "NONE"
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *ProfileRequest) Type() messageType {
	return ProfRq
}

// SignonInfo provides the requirements to login to a single signon realm. A
// signon realm consists of all MessageSets which can be accessed using one set
// of login credentials. Most FI's only use one signon realm to make it easier
// and less confusing for the user.
type SignonInfo struct {
	XMLName           xml.Name `xml:"SIGNONINFO"`
	SignonRealm       String   `xml:"SIGNONREALM"`              // The SignonRealm for which this SignonInfo provides information. This SignonInfo is valid for all MessageSets with SignonRealm fields matching this one
	Min               Int      `xml:"MIN"`                      // Minimum number of password characters
	Max               Int      `xml:"MAX"`                      // Maximum number of password characters
	CharType          charType `xml:"CHARTYPE"`                 // One of ALPHAONLY, NUMERICONLY, ALPHAORNUMERIC, ALPHAANDNUMERIC
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

// MessageSet represents one message set supported by an FI and its
// capabilities
type MessageSet struct {
	XMLName     xml.Name // <xxxMSGSETVn>
	Name        string   // <xxxMSGSETVn> (copy of XMLName.Local)
	Ver         Int      `xml:"MSGSETCORE>VER"`                   // Message set version - should always match 'n' in <xxxMSGSETVn> of Name
	URL         String   `xml:"MSGSETCORE>URL"`                   // URL where messages in this set are to be set
	OfxSec      ofxSec   `xml:"MSGSETCORE>OFXSEC"`                // NONE or 'TYPE 1'
	TranspSec   Boolean  `xml:"MSGSETCORE>TRANSPSEC"`             // Transport-level security must be used
	SignonRealm String   `xml:"MSGSETCORE>SIGNONREALM"`           // Used to identify which SignonInfo to use for to this MessageSet
	Language    []String `xml:"MSGSETCORE>LANGUAGE"`              // List of supported languages
	SyncMode    syncMode `xml:"MSGSETCORE>SYNCMODE"`              // One of FULL, LITE
	RefreshSupt Boolean  `xml:"MSGSETCORE>REFRESHSUPT,omitempty"` // Y if server supports <REFRESH>Y within synchronizations. This option is irrelevant for full synchronization servers. Clients must ignore <REFRESHSUPT> (or its absence) if the profile also specifies <SYNCMODE>FULL. For lite synchronization, the default is N. Without <REFRESHSUPT>Y, lite synchronization servers are not required to support <REFRESH>Y requests
	RespFileER  Boolean  `xml:"MSGSETCORE>RESPFILEER"`            // server supports file-based error recovery
	SpName      String   `xml:"MSGSETCORE>SPNAME"`                // Name of service provider
	// TODO MessageSet-specific stuff?
}

// MessageSetList is a list of MessageSets (necessary because they must be
// manually parsed)
type MessageSetList []MessageSet

// UnmarshalXML handles unmarshalling a MessageSetList element from an XML string
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

// ProfileResponse contains a requested profile of the server's capabilities
// (which message sets and versions it supports, how to access them, which
// languages and which types of synchronization they support, etc.). Note that
// if the server does not support ClientRouting=NONE (as we always send with
// ProfileRequest), this may be an error)
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

// Name returns the name of the top-level transaction XML/SGML element
func (pr *ProfileResponse) Name() string {
	return "PROFTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (pr *ProfileResponse) Valid(version ofxVersion) (bool, error) {
	if ok, err := pr.TrnUID.Valid(); !ok {
		return false, err
	}
	//TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (pr *ProfileResponse) Type() messageType {
	return ProfRs
}
