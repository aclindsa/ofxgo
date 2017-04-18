package ofxgo

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
	"golang.org/x/text/currency"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Int provides helper methods to unmarshal int64 values from SGML/XML
type Int int64

// UnmarshalXML handles unmarshalling an Int from an SGML/XML string. Leading
// and trailing whitespace is ignored.
func (i *Int) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string

	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	value = strings.TrimSpace(value)

	i2, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	*i = Int(i2)
	return nil
}

// Equal returns true if the two Ints are equal in value
func (i Int) Equal(o Int) bool {
	return i == o
}

// Amount represents non-integer values (or at least values for fields that may
// not necessarily be integers)
type Amount struct {
	big.Rat
}

// UnmarshalXML handles unmarshalling an Amount from an SGML/XML string.
// Leading and trailing whitespace is ignored.
func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string

	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	value = strings.TrimSpace(value)

	// The OFX spec allows the start of the fractional amount to be delineated
	// by a comma, so fix that up before attempting to parse it into big.Rat
	value = strings.Replace(value, ",", ".", 1)

	if _, ok := a.SetString(value); !ok {
		return errors.New("Failed to parse OFX amount")
	}
	return nil
}

// String prints a string representation of an Amount
func (a Amount) String() string {
	return strings.TrimRight(strings.TrimRight(a.FloatString(100), "0"), ".")
}

// MarshalXML marshals an Amount to SGML/XML
func (a *Amount) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(a.String(), start)
}

// Equal returns true if two Amounts are equal in value
func (a Amount) Equal(o Amount) bool {
	return (&a).Cmp(&o.Rat) == 0
}

// Date represents OFX date/time values
type Date struct {
	time.Time
}

var ofxDateFormats = []string{
	"20060102150405.000",
	"20060102150405",
	"200601021504",
	"2006010215",
	"20060102",
}
var ofxDateZoneRegex = regexp.MustCompile(`^([+-]?[0-9]+)(\.([0-9]{2}))?(:([A-Z]+))?$`)

// UnmarshalXML handles unmarshalling a Date from an SGML/XML string. It
// attempts to unmarshal the valid date formats in order of decreasing length
// and defaults to GMT if a time zone is not provided, as per the OFX spec.
// Leading and trailing whitespace is ignored.
func (od *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value, zone, zoneFormat string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	value = strings.SplitN(value, "]", 2)[0]
	value = strings.TrimSpace(value)

	// Split the time zone off, if any
	split := strings.SplitN(value, "[", 2)
	if len(split) == 2 {
		value = split[0]
		zoneFormat = " -0700"
		zone = strings.TrimRight(split[1], "]")

		matches := ofxDateZoneRegex.FindStringSubmatch(zone)
		if matches == nil {
			return errors.New("Invalid OFX Date timezone format: " + zone)
		}
		var err error
		var zonehours, zoneminutes int
		zonehours, err = strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		if len(matches[3]) > 0 {
			zoneminutes, err = strconv.Atoi(matches[3])
			if err != nil {
				return err
			}
			zoneminutes = zoneminutes * 60 / 100
		}
		zone = fmt.Sprintf(" %+03d%02d", zonehours, zoneminutes)

		// Get the time zone name if it's there, default to GMT if the offset
		// is 0 and a name isn't supplied
		if len(matches[5]) > 0 {
			zone = zone + " " + matches[5]
			zoneFormat = zoneFormat + " MST"
		} else if zonehours == 0 && zoneminutes == 0 {
			zone = zone + " GMT"
			zoneFormat = zoneFormat + " MST"
		}
	} else {
		// Default to GMT if no time zone was specified
		zone = " +0000 GMT"
		zoneFormat = " -0700 MST"
	}

	// Try all the date formats, from longest to shortest
	for _, format := range ofxDateFormats {
		t, err := time.Parse(format+zoneFormat, value+zone)
		if err == nil {
			od.Time = t
			return nil
		}
	}
	return errors.New("OFX: Couldn't parse date:" + value)
}

// String returns a string representation of the Date abiding by the OFX spec
func (od Date) String() string {
	format := od.Format(ofxDateFormats[0])
	zonename, zoneoffset := od.Zone()
	if zoneoffset < 0 {
		format += "[" + fmt.Sprintf("%+d", zoneoffset/3600)
	} else {
		format += "[" + fmt.Sprintf("%d", zoneoffset/3600)
	}
	fractionaloffset := (zoneoffset % 3600) / 36
	if fractionaloffset > 0 {
		format += "." + fmt.Sprintf("%02d", fractionaloffset)
	} else if fractionaloffset < 0 {
		format += "." + fmt.Sprintf("%02d", -fractionaloffset)
	}

	if len(zonename) > 0 {
		return format + ":" + zonename + "]"
	}
	return format + "]"
}

// MarshalXML marshals a Date to XML
func (od *Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(od.String(), start)
}

// Equal returns true if the two Dates represent the same time (time zones are
// accounted for when comparing, but are not required to match)
func (od Date) Equal(o Date) bool {
	return od.Time.Equal(o.Time)
}

// NewDate returns a new Date object with the provided date, time, and timezone
func NewDate(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) *Date {
	return &Date{Time: time.Date(year, month, day, hour, min, sec, nsec, loc)}
}

var gmt = time.FixedZone("GMT", 0)

// NewDateGMT returns a new Date object with the provided date and time in the
// GMT timezone
func NewDateGMT(year int, month time.Month, day, hour, min, sec, nsec int) *Date {
	return &Date{Time: time.Date(year, month, day, hour, min, sec, nsec, gmt)}
}

// String provides helper methods to unmarshal OFX string values from SGML/XML
type String string

// UnmarshalXML handles unmarshalling a String from an SGML/XML string. Leading
// and trailing whitespace is ignored.
func (os *String) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	*os = String(strings.TrimSpace(value))
	return nil
}

// String returns the string
func (os *String) String() string {
	return string(*os)
}

// Equal returns true if the two Strings are equal in value
func (os String) Equal(o String) bool {
	return os == o
}

// Boolean provides helper methods to unmarshal bool values from OFX SGML/XML
type Boolean bool

// UnmarshalXML handles unmarshalling a Boolean from an SGML/XML string.
// Leading and trailing whitespace is ignored.
func (ob *Boolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	tmpob := strings.TrimSpace(value)
	switch tmpob {
	case "Y":
		*ob = Boolean(true)
	case "N":
		*ob = Boolean(false)
	default:
		return errors.New("Invalid OFX Boolean")
	}
	return nil
}

// MarshalXML marshals a Boolean to XML
func (ob *Boolean) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if *ob {
		return e.EncodeElement("Y", start)
	}
	return e.EncodeElement("N", start)
}

// String returns a string representation of a Boolean value
func (ob *Boolean) String() string {
	return fmt.Sprintf("%v", *ob)
}

// Equal returns true if the two Booleans are the same
func (ob Boolean) Equal(o Boolean) bool {
	return ob == o
}

// UID represents an UID according to the OFX spec
type UID string

// UnmarshalXML handles unmarshalling an UID from an SGML/XML string. Leading
// and trailing whitespace is ignored.
func (ou *UID) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	*ou = UID(strings.TrimSpace(value))
	return nil
}

// RecommendedFormat returns true iff this UID meets the OFX specification's
// recommendation that UIDs follow the standard UUID 36-character format
func (ou UID) RecommendedFormat() (bool, error) {
	if len(ou) != 36 {
		return false, errors.New("UID not 36 characters long")
	}
	if ou[8] != '-' || ou[13] != '-' || ou[18] != '-' || ou[23] != '-' {
		return false, errors.New("UID missing hyphens at the appropriate places")
	}
	return true, nil
}

// Valid returns true, nil if the UID is valid. This is less strict than
// RecommendedFormat, and will always return true, nil if it does.
func (ou UID) Valid() (bool, error) {
	if len(ou) == 0 || len(ou) > 36 {
		return false, errors.New("UID invalid length")
	}
	return true, nil
}

// Equal returns true if the two UIDs are the same
func (ou UID) Equal(o UID) bool {
	return ou == o
}

// RandomUID creates a new randomly-generated UID
func RandomUID() (*UID, error) {
	uidbytes := make([]byte, 16)
	n, err := rand.Read(uidbytes[:])
	if err != nil {
		return nil, err
	}
	if n != 16 {
		return nil, errors.New("RandomUID failed to read 16 random bytes")
	}
	uid := UID(fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", uidbytes[:4], uidbytes[4:6], uidbytes[6:8], uidbytes[8:10], uidbytes[10:]))
	return &uid, nil
}

// CurrSymbol represents an ISO-4217 currency
type CurrSymbol struct {
	currency.Unit
}

// UnmarshalXML handles unmarshalling a CurrSymbol from an SGML/XML string.
// Leading and trailing whitespace is ignored.
func (c *CurrSymbol) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string

	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	value = strings.TrimSpace(value)

	unit, err := currency.ParseISO(value)
	if err != nil {
		errors.New("Error parsing CurrSymbol:" + err.Error())
	}
	c.Unit = unit
	return nil
}

// MarshalXML marshals a CurrSymbol to SGML/XML
func (c *CurrSymbol) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(c.String(), start)
}

// Equal returns true if the two Currencies are the same
func (c CurrSymbol) Equal(o CurrSymbol) bool {
	return c.String() == o.String()
}

// Valid returns true, nil if the CurrSymbol is valid.
func (c CurrSymbol) Valid() (bool, error) {
	if c.String() == "XXX" {
		return false, fmt.Errorf("Invalid CurrSymbol: %s", c.Unit)
	}
	return true, nil
}

func NewCurrSymbol(s string) (*CurrSymbol, error) {
	unit, err := currency.ParseISO(s)
	if err != nil {
		return nil, errors.New("Error parsing string to create new CurrSymbol:" + err.Error())
	}
	return &CurrSymbol{unit}, nil
}
