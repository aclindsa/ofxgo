package ofxgo

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

// SecurityID identifies a security by its CUSIP (for US-based FI's, others may
// use UniqueID types other than CUSIP)
type SecurityID struct {
	XMLName      xml.Name `xml:"SECID"`
	UniqueID     String   `xml:"UNIQUEID"`     // CUSIP for US FI's
	UniqueIDType String   `xml:"UNIQUEIDTYPE"` // Should always be "CUSIP" for US FI's
}

// SecurityRequest represents a request for one security. It is specified with
// a SECID aggregate, a ticker symbol, or an FI assigned identifier (but no
// more than one of them at a time)
type SecurityRequest struct {
	XMLName xml.Name `xml:"SECRQ"`
	// Only one of the next three should be present
	SecID  *SecurityID `xml:"SECID,omitempty"`
	Ticker String      `xml:"TICKER,omitempty"`
	FiID   String      `xml:"FIID,omitempty"`
}

// SecListRequest represents a request for information (namely price) about one
// or more securities
type SecListRequest struct {
	XMLName   xml.Name `xml:"SECLISTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	Securities []SecurityRequest `xml:"SECLISTRQ>SECRQ,omitempty"`
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *SecListRequest) Name() string {
	return "SECLISTTRNRQ"
}

// Valid returns (true, nil) if this struct would be valid OFX if marshalled
// into XML/SGML
func (r *SecListRequest) Valid(version ofxVersion) (bool, error) {
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Request
// element of type []Message it should appended to)
func (r *SecListRequest) Type() messageType {
	return SecListRq
}

// SecListResponse is always empty (except for the transaction UID, status, and
// optional client cookie). Its presence signifies that the SecurityList (a
// different element from this one) immediately after this element in
// Response.SecList was been generated in response to the same SecListRequest
// this is a response to.
type SecListResponse struct {
	XMLName   xml.Name `xml:"SECLISTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	// SECLISTRS is always empty, so we don't parse it here. The actual securities list will be in a top-level element parallel to SECLISTTRNRS
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *SecListResponse) Name() string {
	return "SECLISTTRNRS"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (r *SecListResponse) Valid(version ofxVersion) (bool, error) {
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (r *SecListResponse) Type() messageType {
	return SecListRs
}

// Security is satisfied by all *Info elements providing information about
// securities for SecurityList
type Security interface {
	SecurityType() string
}

// SecInfo represents the generic information about a security. It is included
// in most other *Info elements.
type SecInfo struct {
	XMLName   xml.Name   `xml:"SECINFO"`
	SecID     SecurityID `xml:"SECID"`
	SecName   String     `xml:"SECNAME"`          // Full name of security
	Ticker    String     `xml:"TICKER,omitempty"` // Ticker symbol
	FiID      String     `xml:"FIID,omitempty"`
	Rating    String     `xml:"RATING,omitempty"`
	UnitPrice Amount     `xml:"UNITPRICE,omitempty"` // Current price, as of DTASOF
	DtAsOf    *Date      `xml:"DTASOF,omitempty"`    // Date UNITPRICE was for
	Currency  *Currency  `xml:"CURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Memo      String     `xml:"MEMO,omitempty"`
}

// DebtInfo provides information about a debt security
type DebtInfo struct {
	XMLName      xml.Name   `xml:"DEBTINFO"`
	SecInfo      SecInfo    `xml:"SECINFO"`
	ParValue     Amount     `xml:"PARVALUE"`
	DebtType     debtType   `xml:"DEBTTYPE"`               // One of COUPON, ZERO (zero coupon)
	DebtClass    debtClass  `xml:"DEBTCLASS,omitempty"`    // One of TREASURY, MUNICIPAL, CORPORATE, OTHER
	CouponRate   Amount     `xml:"COUPONRT,omitempty"`     // Bond coupon rate for next closest call date
	DtCoupon     *Date      `xml:"DTCOUPON,omitempty"`     // Maturity date for next coupon
	CouponFreq   couponFreq `xml:"COUPONFREQ,omitempty"`   // When coupons mature - one of MONTHLY, QUARTERLY, SEMIANNUAL, ANNUAL, or OTHER
	CallPrice    Amount     `xml:"CALLPRICE,omitempty"`    // Bond call price
	YieldToCall  Amount     `xml:"YIELDTOCALL,omitempty"`  // Yield to next call
	DtCall       *Date      `xml:"DTCALL,omitempty"`       // Next call date
	CallType     callType   `xml:"CALLTYPE,omitempt"`      // Type of next call. One of CALL, PUT, PREFUND, MATURITY
	YieldToMat   Amount     `xml:"YIELDTOMAT,omitempty"`   // Yield to maturity
	DtMat        *Date      `xml:"DTMAT,omitempty"`        // Debt maturity date
	AssetClass   assetClass `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String     `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

// SecurityType returns a string representation of this security's type
func (i DebtInfo) SecurityType() string {
	return "DEBTINFO"
}

// AssetPortion represents the percentage of a mutual fund with the given asset
// classification
type AssetPortion struct {
	XMLName    xml.Name   `xml:"PORTION"`
	AssetClass assetClass `xml:"ASSETCLASS"` // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	Percent    Amount     `xml:"PERCENT"`    // Percentage of the fund that falls under this asset class
}

// FiAssetPortion represents the percentage of a mutual fund with the given
// FI-defined asset classification (AssetPortion should be used for all asset
// classifications defined by the assetClass enum)
type FiAssetPortion struct {
	XMLName      xml.Name `xml:"FIPORTION"`
	FiAssetClass String   `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
	Percent      Amount   `xml:"PERCENT"`                // Percentage of the fund that falls under this asset class
}

// MFInfo provides information about a mutual fund
type MFInfo struct {
	XMLName        xml.Name         `xml:"MFINFO"`
	SecInfo        SecInfo          `xml:"SECINFO"`
	MfType         mfType           `xml:"MFTYPE"`                // One of OPEN, END, CLOSEEND, OTHER
	Yield          Amount           `xml:"YIELD,omitempty"`       // Current yield reported as the dividend expressed as a portion of the current stock price
	DtYieldAsOf    *Date            `xml:"DTYIELDASOF,omitempty"` // Date YIELD is valid for
	AssetClasses   []AssetPortion   `xml:"MFASSETCLASS>PORTION"`
	FiAssetClasses []FiAssetPortion `xml:"FIMFASSETCLASS>FIPORTION"`
}

// SecurityType returns a string representation of this security's type
func (i MFInfo) SecurityType() string {
	return "MFINFO"
}

// OptInfo provides information about an option
type OptInfo struct {
	XMLName      xml.Name    `xml:"OPTINFO"`
	SecInfo      SecInfo     `xml:"SECINFO"`
	OptType      optType     `xml:"OPTTYPE"` // One of PUT, CALL
	StrikePrice  Amount      `xml:"STRIKEPRICE"`
	DtExpire     Date        `xml:"DTEXPIRE"`               // Expiration date
	ShPerCtrct   Int         `xml:"SHPERCTRCT"`             // Shares per contract
	SecID        *SecurityID `xml:"SECID,omitempty"`        // Security ID of the underlying security
	AssetClass   assetClass  `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String      `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

// SecurityType returns a string representation of this security's type
func (i OptInfo) SecurityType() string {
	return "OPTINFO"
}

// OtherInfo provides information about a security type not covered by the
// other *Info elements
type OtherInfo struct {
	XMLName      xml.Name   `xml:"OTHERINFO"`
	SecInfo      SecInfo    `xml:"SECINFO"`
	TypeDesc     String     `xml:"TYPEDESC,omitempty"`     // Description of security type
	AssetClass   assetClass `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String     `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

// SecurityType returns a string representation of this security's type
func (i OtherInfo) SecurityType() string {
	return "OTHERINFO"
}

// StockInfo provides information about a security type
type StockInfo struct {
	XMLName      xml.Name   `xml:"STOCKINFO"`
	SecInfo      SecInfo    `xml:"SECINFO"`
	StockType    stockType  `xml:"STOCKTYPE,omitempty"`    // One of COMMON, PREFERRED, CONVERTIBLE, OTHER
	Yield        Amount     `xml:"YIELD,omitempty"`        // Current yield reported as the dividend expressed as a portion of the current stock price
	DtYieldAsOf  *Date      `xml:"DTYIELDASOF,omitempty"`  // Date YIELD is valid for
	AssetClass   assetClass `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String     `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

// SecurityType returns a string representation of this security's type
func (i StockInfo) SecurityType() string {
	return "STOCKINFO"
}

// SecurityList is a container for Security objects containaing information
// about securities
type SecurityList struct {
	Securities []Security
}

// Name returns the name of the top-level transaction XML/SGML element
func (r *SecurityList) Name() string {
	return "SECLIST"
}

// Valid returns (true, nil) if this struct was valid OFX when unmarshalled
func (r *SecurityList) Valid(version ofxVersion) (bool, error) {
	// TODO implement
	return true, nil
}

// Type returns which message set this message belongs to (which Response
// element of type []Message it belongs to)
func (r *SecurityList) Type() messageType {
	return SecListRs
}

// UnmarshalXML handles unmarshalling a SecurityList from an SGML/XML string
func (r *SecurityList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := nextNonWhitespaceToken(d)
		if err != nil {
			return err
		} else if end, ok := tok.(xml.EndElement); ok && end.Name.Local == start.Name.Local {
			// If we found the end of our starting element, we're done parsing
			return nil
		} else if startElement, ok := tok.(xml.StartElement); ok {
			switch startElement.Name.Local {
			case "DEBTINFO":
				var info DebtInfo
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return err
				}
				r.Securities = append(r.Securities, Security(info))
			case "MFINFO":
				var info MFInfo
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return err
				}
				r.Securities = append(r.Securities, Security(info))
			case "OPTINFO":
				var info OptInfo
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return err
				}
				r.Securities = append(r.Securities, Security(info))
			case "OTHERINFO":
				var info OtherInfo
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return err
				}
				r.Securities = append(r.Securities, Security(info))
			case "STOCKINFO":
				var info StockInfo
				if err := d.DecodeElement(&info, &startElement); err != nil {
					return err
				}
				r.Securities = append(r.Securities, Security(info))
			default:
				return errors.New("Invalid SECLIST child tag: " + startElement.Name.Local)
			}
		} else {
			return errors.New("Didn't find an opening element")
		}
	}
}
