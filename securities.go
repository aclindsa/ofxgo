package ofxgo

import (
	"errors"
	"github.com/aclindsa/go/src/encoding/xml"
)

type SecurityId struct {
	XMLName      xml.Name `xml:"SECID"`
	UniqueId     String   `xml:"UNIQUEID"`     // CUSIP for US FI's
	UniqueIdType String   `xml:"UNIQUEIDTYPE"` // Should always be "CUSIP" for US FI's
}

type SecurityRequest struct {
	XMLName xml.Name `xml:"SECRQ"`
	// Only one of the next three should be present
	SecId  *SecurityId `xml:"SECID,omitempty"`
	Ticker String      `xml:"TICKER,omitempty"`
	FiId   String      `xml:"FIID,omitempty"`
}

type SecListRequest struct {
	XMLName   xml.Name `xml:"SECLISTTRNRQ"`
	TrnUID    UID      `xml:"TRNUID"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	TAN       String   `xml:"TAN,omitempty"` // Transaction authorization number
	// TODO `xml:"OFXEXTENSION,omitempty"`
	Securities []SecurityRequest `xml:"SECLISTRQ>SECRQ,omitempty"`
}

func (r *SecListRequest) Name() string {
	return "SECLISTTRNRQ"
}

func (r *SecListRequest) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

func (r *SecListRequest) Type() messageType {
	return SecListRq
}

type SecListResponse struct {
	XMLName   xml.Name `xml:"SECLISTTRNRS"`
	TrnUID    UID      `xml:"TRNUID"`
	Status    Status   `xml:"STATUS"`
	CltCookie String   `xml:"CLTCOOKIE,omitempty"`
	// TODO `xml:"OFXEXTENSION,omitempty"`
	// SECLISTRS is always empty, so we don't parse it here. The actual securities list will be in a top-level element parallel to SECLISTTRNRS
}

func (r *SecListResponse) Name() string {
	return "SECLISTTRNRS"
}

func (r *SecListResponse) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

func (r *SecListResponse) Type() messageType {
	return SecListRs
}

type Security interface {
	SecurityType() string
}

type SecInfo struct {
	XMLName   xml.Name   `xml:"SECINFO"`
	SecId     SecurityId `xml:"SECID"`
	SecName   String     `xml:"SECNAME"`          // Full name of security
	Ticker    String     `xml:"TICKER,omitempty"` // Ticker symbol
	FiId      String     `xml:"FIID,omitempty"`
	Rating    String     `xml:"RATING,omitempty"`
	UnitPrice Amount     `xml:"UNITPRICE,omitempty"` // Current price, as of DTASOF
	DtAsOf    *Date      `xml:"DTASOF,omitempty"`    // Date UNITPRICE was for
	Currency  *Currency  `xml:"CURRENCY,omitempty"`  // Overriding currency for UNITPRICE
	Memo      String     `xml:"MEMO,omitempty"`
}

type DebtInfo struct {
	XMLName      xml.Name `xml:"DEBTINFO"`
	SecInfo      SecInfo  `xml:"SECINFO"`
	ParValue     Amount   `xml:"PARVALUE"`
	DebtType     String   `xml:"DEBTTYPE"`               // One of COUPON, ZERO (zero coupon)
	DebtClass    String   `xml:"DEBTCLASS,omitempty"`    // One of TREASURY, MUNICIPAL, CORPORATE, OTHER
	CouponRate   Amount   `xml:"COUPONRT,omitempty"`     // Bond coupon rate for next closest call date
	DtCoupon     *Date    `xml:"DTCOUPON,omitempty"`     // Maturity date for next coupon
	CouponFreq   String   `xml:"COUPONFREQ,omitempty"`   // When coupons mature - one of MONTHLY, QUARTERLY, SEMIANNUAL, ANNUAL, or OTHER
	CallPrice    Amount   `xml:"CALLPRICE,omitempty"`    // Bond call price
	YieldToCall  Amount   `xml:"YIELDTOCALL,omitempty"`  // Yield to next call
	DtCall       *Date    `xml:"DTCALL,omitempty"`       // Next call date
	CallType     String   `xml:"CALLTYPE,omitempt"`      // Type of next call. One of CALL, PUT, PREFUND, MATURITY
	YieldToMat   Amount   `xml:"YIELDTOMAT,omitempty"`   // Yield to maturity
	DtMat        *Date    `xml:"DTMAT,omitempty"`        // Debt maturity date
	AssetClass   String   `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String   `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

func (i DebtInfo) SecurityType() string {
	return "DEBTINFO"
}

type AssetPortion struct {
	XMLName    xml.Name `xml:"PORTION"`
	AssetClass String   `xml:"ASSETCLASS"` // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	Percent    Amount   `xml:"PERCENT"`    // Percentage of the fund that falls under this asset class
}

type FiAssetPortion struct {
	XMLName      xml.Name `xml:"FIPORTION"`
	FiAssetClass String   `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
	Percent      Amount   `xml:"PERCENT"`                // Percentage of the fund that falls under this asset class
}

type MFInfo struct {
	XMLName        xml.Name         `xml:"MFINFO"`
	SecInfo        SecInfo          `xml:"SECINFO"`
	MfType         String           `xml:"MFTYPE"`                // One of OPEN, END, CLOSEEND, OTHER
	Yield          Amount           `xml:"YIELD,omitempty"`       // Current yield reported as the dividend expressed as a portion of the current stock price
	DtYieldAsOf    *Date            `xml:"DTYIELDASOF,omitempty"` // Date YIELD is valid for
	AssetClasses   []AssetPortion   `xml:"MFASSETCLASS>PORTION"`
	FiAssetClasses []FiAssetPortion `xml:"FIMFASSETCLASS>FIPORTION"`
}

func (i MFInfo) SecurityType() string {
	return "MFINFO"
}

type OptInfo struct {
	XMLName      xml.Name    `xml:"OPTINFO"`
	SecInfo      SecInfo     `xml:"SECINFO"`
	OptType      String      `xml:"OPTTYPE"` // One of PUT, CALL
	StrikePrice  Amount      `xml:"STRIKEPRICE"`
	DtExpire     Date        `xml:"DTEXPIRE"`               // Expiration date
	ShPerCtrct   Int         `xml:"SHPERCTRCT"`             // Shares per contract
	SecId        *SecurityId `xml:"SECID,omitempty"`        // Security ID of the underlying security
	AssetClass   String      `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String      `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

func (i OptInfo) SecurityType() string {
	return "OPTINFO"
}

type OtherInfo struct {
	XMLName      xml.Name `xml:"OTHERINFO"`
	SecInfo      SecInfo  `xml:"SECINFO"`
	TypeDesc     String   `xml:"TYPEDESC,omitempty"`     // Description of security type
	AssetClass   String   `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String   `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

func (i OtherInfo) SecurityType() string {
	return "OTHERINFO"
}

type StockInfo struct {
	XMLName      xml.Name `xml:"STOCKINFO"`
	SecInfo      SecInfo  `xml:"SECINFO"`
	StockType    String   `xml:"STOCKTYPE,omitempty"`    // One of COMMON, PREFERRED, CONVERTIBLE, OTHER
	Yield        Amount   `xml:"YIELD,omitempty"`        // Current yield reported as the dividend expressed as a portion of the current stock price
	DtYieldAsOf  *Date    `xml:"DTYIELDASOF,omitempty"`  // Date YIELD is valid for
	AssetClass   String   `xml:"ASSETCLASS,omitempty"`   // One of DOMESTICBOND, INTLBOND, LARGESTOCK, SMALLSTOCK, INTLSTOCK, MONEYMRKT, OTHER
	FiAssetClass String   `xml:"FIASSETCLASS,omitempty"` // FI-defined asset class
}

func (i StockInfo) SecurityType() string {
	return "STOCKINFO"
}

type SecurityList struct {
	Securities []Security
}

func (r *SecurityList) Name() string {
	return "SECLIST"
}

func (r *SecurityList) Valid() (bool, error) {
	// TODO implement
	return true, nil
}

func (r *SecurityList) Type() messageType {
	return SecListRs
}

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
