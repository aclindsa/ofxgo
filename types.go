package ofxgo

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang/go/src/encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Int int64

func (oi *Int) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return err
	}
	*oi = (Int)(i)
	return nil
}

type Amount string

// TODO parse Amount into big.Rat?

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
			zoneminutes, err = strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
			zoneminutes = zoneminutes * 60 / 100
		}
		zone = fmt.Sprintf(" %+03d%02d", zonehours, zoneminutes)
	}

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
	format += "[" + fmt.Sprintf("%+d", zoneoffset/3600)
	fractionaloffset := (zoneoffset % 3600) / 360
	if fractionaloffset > 0 {
		format += "." + fmt.Sprintf("%02d", fractionaloffset)
	} else if fractionaloffset < 0 {
		format += "." + fmt.Sprintf("%02d", -fractionaloffset)
	}
	return format + ":" + zonename + "]"
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

type UID string

func (ou *UID) Valid() (bool, error) {
	if len(*ou) != 36 {
		return false, errors.New("UID not 36 characters long")
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
