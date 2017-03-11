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

type OfxInt int64

func (oi *OfxInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return err
	}
	*oi = (OfxInt)(i)
	return nil
}

var ofxDateFormats = []string{
	"20060102150405.000",
	"20060102150405",
	"200601021504",
	"2006010215",
	"20060102",
}
var ofxDateZoneFormat = "20060102150405.000 -0700"
var ofxDateZoneRegex = regexp.MustCompile(`^\[([+-]?[0-9]+)(\.([0-9]{2}))?(:([A-Z]+))?\]$`)

type OfxDate time.Time

func (od *OfxDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	value = strings.TrimSpace(value)

	if len(value) > len(ofxDateFormats[0]) {
		matches := ofxDateZoneRegex.FindStringSubmatch(value[len(ofxDateFormats[0]):])
		if matches == nil {
			return errors.New("Invalid OFX Date Format")
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
		value = value[:len(ofxDateFormats[0])] + " " + fmt.Sprintf("%+d%02d", zonehours, zoneminutes)
		t, err := time.Parse(ofxDateZoneFormat, value)
		if err == nil {
			tmpod := OfxDate(t)
			*od = tmpod
			return nil
		}
	}

	for _, format := range ofxDateFormats {
		t, err := time.Parse(format, value)
		if err == nil {
			tmpod := OfxDate(t)
			*od = tmpod
			return nil
		}
	}
	return errors.New("OFX: Couldn't parse date:" + value)
}

func (od *OfxDate) String() string {
	t := time.Time(*od)
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

func (od *OfxDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(od.String(), start)
}

type OfxString string

func (os *OfxString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	*os = OfxString(strings.TrimSpace(value))
	return nil
}

type OfxBoolean bool

func (ob *OfxBoolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}
	tmpob := strings.TrimSpace(value)
	switch tmpob {
	case "Y":
		*ob = OfxBoolean(true)
	case "N":
		*ob = OfxBoolean(false)
	default:
		return errors.New("Invalid OFX Boolean")
	}
	return nil
}

func (ob *OfxBoolean) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if *ob {
		return e.EncodeElement("Y", start)
	}
	return e.EncodeElement("N", start)
}

type OfxUID string

func (ou *OfxUID) Valid() (bool, error) {
	if len(*ou) != 36 {
		return false, errors.New("UID not 36 characters long")
	}
	return true, nil
}

func RandomUID() (*OfxUID, error) {
	uidbytes := make([]byte, 16)
	n, err := rand.Read(uidbytes[:])
	if err != nil {
		return nil, err
	}
	if n != 16 {
		return nil, errors.New("RandomUID failed to read 16 random bytes")
	}
	uid := OfxUID(fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", uidbytes[:4], uidbytes[4:6], uidbytes[6:8], uidbytes[8:10], uidbytes[10:]))
	return &uid, nil
}
