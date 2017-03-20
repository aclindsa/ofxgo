package ofxgo

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang/go/src/encoding/xml"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Int int64

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

type Amount big.Rat

func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	var b big.Rat

	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	value = strings.TrimSpace(value)

	// The OFX spec allows the start of the fractional amount to be delineated
	// by a comma, so fix that up before attempting to parse it into big.Rat
	value = strings.Replace(value, ",", ".", 1)

	if _, ok := b.SetString(value); !ok {
		return errors.New("Failed to parse OFX amount into big.Rat")
	}
	*a = Amount(b)
	return nil
}

func (a Amount) String() string {
	var b big.Rat = big.Rat(a)
	return strings.TrimRight(strings.TrimRight(b.FloatString(100), "0"), ".")
}

func (a *Amount) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(a.String(), start)
}

type Date time.Time

var ofxDateFormats = []string{
	"20060102150405.000",
	"20060102150405",
	"200601021504",
	"2006010215",
	"20060102",
}
var ofxDateZoneRegex = regexp.MustCompile(`^([+-]?[0-9]+)(\.([0-9]{2}))?(:([A-Z]+))?$`)

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
			tmpod := Date(t)
			*od = tmpod
			return nil
		}
	}
	return errors.New("OFX: Couldn't parse date:" + value)
}

func (od Date) String() string {
	t := time.Time(od)
	format := t.Format(ofxDateFormats[0])
	zonename, zoneoffset := t.Zone()
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
	} else {
		return format + "]"
	}
}

func (od *Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(od.String(), start)
}

type String string

func (os *String) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	*os = String(strings.TrimSpace(value))
	return nil
}

func (os *String) String() string {
	return string(*os)
}

type Boolean bool

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

func (ob *Boolean) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if *ob {
		return e.EncodeElement("Y", start)
	}
	return e.EncodeElement("N", start)
}

func (ob *Boolean) String() string {
	return fmt.Sprintf("%v", *ob)
}

type UID string

func (ou UID) Valid() (bool, error) {
	if len(ou) != 36 {
		return false, errors.New("UID not 36 characters long")
	}
	if ou[8] != '-' || ou[13] != '-' || ou[18] != '-' || ou[23] != '-' {
		return false, errors.New("UID missing hyphens at the appropriate places")
	}
	return true, nil
}

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
