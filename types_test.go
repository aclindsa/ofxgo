package ofxgo_test

import (
	"fmt"
	"github.com/aclindsa/go/src/encoding/xml"
	"github.com/aclindsa/ofxgo"
	"math/big"
	"reflect"
	"testing"
	"time"
)

func getTypeName(i interface{}) string {
	val := reflect.ValueOf(i)

	// Do the same thing that encoding/xml does to get the name
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}
	return val.Type().Name()
}

func marshalHelper(t *testing.T, expected string, i interface{}) {
	typename := getTypeName(i)
	expectedstring := fmt.Sprintf("<%s>%s</%s>", typename, expected, typename)
	b, err := xml.Marshal(i)
	if err != nil {
		t.Fatalf("Unexpected error on xml.Marshal(%T): %s\n", i, err)
	}
	if string(b) != expectedstring {
		t.Fatalf("Expected '%s', got '%s'\n", expectedstring, string(b))
	}
}

func unmarshalHelper2(t *testing.T, input string, expected interface{}, overwritten interface{}, eq func(a, b interface{}) bool) {
	typename := getTypeName(expected)
	inputstring := fmt.Sprintf("<%s>%s</%s>", typename, input, typename)
	err := xml.Unmarshal([]byte(inputstring), &overwritten)
	if err != nil {
		t.Fatalf("Unexpected error on xml.Unmarshal(%T): %s\n", expected, err)
	}
	if !eq(overwritten, expected) {
		t.Fatalf("Expected '%s', got '%s'\n", expected, overwritten)
	}
}

func unmarshalHelper(t *testing.T, input string, expected interface{}, overwritten interface{}) {
	eq := func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	}
	unmarshalHelper2(t, input, expected, overwritten, eq)
}

func TestMarshalInt(t *testing.T) {
	var i ofxgo.Int = 927
	marshalHelper(t, "927", &i)
	i = 0
	marshalHelper(t, "0", &i)
	i = -768276587425
	marshalHelper(t, "-768276587425", &i)
}

func TestUnmarshalInt(t *testing.T) {
	var i, overwritten ofxgo.Int = -48394, 0
	unmarshalHelper(t, "-48394", &i, &overwritten)
	i = 0
	unmarshalHelper(t, "0", &i, &overwritten)
	i = 198237198
	unmarshalHelper(t, "198237198", &i, &overwritten)
	// Make sure stray newlines are handled properly
	unmarshalHelper(t, "198237198\n", &i, &overwritten)
}

func TestMarshalAmount(t *testing.T) {
	var a ofxgo.Amount
	var b *big.Rat = (*big.Rat)(&a)

	b.SetFrac64(8, 1)
	marshalHelper(t, "8", &a)
	b.SetFrac64(1, 8)
	marshalHelper(t, "0.125", &a)
	b.SetFrac64(-1, 200)
	marshalHelper(t, "-0.005", &a)
	b.SetInt64(0)
	marshalHelper(t, "0", &a)
	b.SetInt64(-768276587425)
	marshalHelper(t, "-768276587425", &a)
	b.SetFrac64(1, 12)
	marshalHelper(t, "0.0833333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333333", &a)
}

func TestUnmarshalAmount(t *testing.T) {
	var a, overwritten ofxgo.Amount
	var b *big.Rat = (*big.Rat)(&a)

	// Amount/big.Rat needs a special equality test because reflect.DeepEqual
	// doesn't always return equal for two values that big.Rat.Cmp() does
	eq := func(a, b interface{}) bool {
		if amountA, ok := a.(*ofxgo.Amount); ok {
			if amountB, ok2 := b.(*ofxgo.Amount); ok2 {
				ratA := (*big.Rat)(amountA)
				return ratA.Cmp((*big.Rat)(amountB)) == 0
			}
		}
		return false
	}

	b.SetFrac64(12, 1)
	unmarshalHelper2(t, "12", &a, &overwritten, eq)
	b.SetFrac64(-21309, 100)
	unmarshalHelper2(t, "-213.09", &a, &overwritten, eq)
	b.SetFrac64(8192, 1000)
	unmarshalHelper2(t, "8.192", &a, &overwritten, eq)
	unmarshalHelper2(t, "+8.192", &a, &overwritten, eq)
	b.SetInt64(0)
	unmarshalHelper2(t, "0", &a, &overwritten, eq)
	unmarshalHelper2(t, "+0", &a, &overwritten, eq)
	unmarshalHelper2(t, "-0", &a, &overwritten, eq)
	b.SetInt64(-19487135)
	unmarshalHelper2(t, "-19487135", &a, &overwritten, eq)
	// Make sure stray newlines are handled properly
	unmarshalHelper2(t, "-19487135\n", &a, &overwritten, eq)
}

func TestMarshalDate(t *testing.T) {
	var d ofxgo.Date
	UTC := time.FixedZone("UTC", 0)
	GMT_nodesc := time.FixedZone("", 0)
	EST := time.FixedZone("EST", -5*60*60)
	NPT := time.FixedZone("NPT", (5*60+45)*60)
	IST := time.FixedZone("IST", (5*60+30)*60)
	NST := time.FixedZone("NST", -(3*60+30)*60)

	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, NPT))
	marshalHelper(t, "20170314150926.053[5.75:NPT]", &d)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, EST))
	marshalHelper(t, "20170314150926.053[-5:EST]", &d)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, UTC))
	marshalHelper(t, "20170314150926.053[0:UTC]", &d)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, IST))
	marshalHelper(t, "20170314150926.053[5.50:IST]", &d)
	d = ofxgo.Date(time.Date(9999, 11, 1, 23, 59, 59, 1000, EST))
	marshalHelper(t, "99991101235959.000[-5:EST]", &d)
	d = ofxgo.Date(time.Date(0, 1, 1, 0, 0, 0, 0, IST))
	marshalHelper(t, "00000101000000.000[5.50:IST]", &d)
	d = ofxgo.Date(time.Unix(0, 0).In(UTC))
	marshalHelper(t, "19700101000000.000[0:UTC]", &d)
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 26, 53*1000*1000, EST))
	marshalHelper(t, "20170314000026.053[-5:EST]", &d)
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 26, 53*1000*1000, NST))
	marshalHelper(t, "20170314000026.053[-3.50:NST]", &d)

	// Time zone without textual description
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, GMT_nodesc))
	marshalHelper(t, "20170314150926.053[0]", &d)
}

func TestUnmarshalDate(t *testing.T) {
	var d, overwritten ofxgo.Date
	GMT := time.FixedZone("GMT", 0)
	EST := time.FixedZone("EST", -5*60*60)
	NPT := time.FixedZone("NPT", (5*60+45)*60)
	IST := time.FixedZone("IST", (5*60+30)*60)
	NST := time.FixedZone("NST", -(3*60+30)*60)
	NST_nodesc := time.FixedZone("", -(3*60+30)*60)

	eq := func(a, b interface{}) bool {
		if dateA, ok := a.(*ofxgo.Date); ok {
			if dateB, ok2 := b.(*ofxgo.Date); ok2 {
				timeA := (*time.Time)(dateA)
				return timeA.Equal((time.Time)(*dateB))
			}
		}
		return false
	}

	// Ensure omitted fields default to the correct values
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, GMT))
	unmarshalHelper2(t, "20170314150926.053[0]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 0, 0, GMT))
	unmarshalHelper2(t, "20170314", &d, &overwritten, eq)

	// Ensure all signs on time zone offsets are properly handled
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, GMT))
	unmarshalHelper2(t, "20170314150926.053[0:GMT]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053[+0:GMT]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053[-0:GMT]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053[0]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053[+0]", &d, &overwritten, eq)
	unmarshalHelper2(t, "20170314150926.053[-0]", &d, &overwritten, eq)

	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, NPT))
	unmarshalHelper2(t, "20170314150926.053[5.75:NPT]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, EST))
	unmarshalHelper2(t, "20170314150926.053[-5:EST]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, GMT))
	unmarshalHelper2(t, "20170314150926.053[0:GMT]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, IST))
	unmarshalHelper2(t, "20170314150926.053[5.50:IST]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2018, 11, 1, 23, 59, 58, 0, EST))
	unmarshalHelper2(t, "20181101235958.000[-5:EST]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(0, 1, 1, 0, 0, 0, 0, IST))
	unmarshalHelper2(t, "00000101000000.000[5.50:IST]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Unix(0, 0).In(GMT))
	unmarshalHelper2(t, "19700101000000.000[0:GMT]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 26, 53*1000*1000, EST))
	unmarshalHelper2(t, "20170314000026.053[-5:EST]", &d, &overwritten, eq)
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 26, 53*1000*1000, NST))
	unmarshalHelper2(t, "20170314000026.053[-3.50:NST]", &d, &overwritten, eq)

	// Autopopulate zone without textual description for GMT
	d = ofxgo.Date(time.Date(2017, 3, 14, 15, 9, 26, 53*1000*1000, GMT))
	unmarshalHelper2(t, "20170314150926.053[0]", &d, &overwritten, eq)
	// but not for others:
	d = ofxgo.Date(time.Date(2017, 3, 14, 0, 0, 26, 53*1000*1000, NST_nodesc))
	unmarshalHelper2(t, "20170314000026.053[-3.50]", &d, &overwritten, eq)

	// Make sure we handle poorly-formatted dates (from Vanguard)
	d = ofxgo.Date(time.Date(2016, 12, 7, 16, 0, 0, 0, EST))
	unmarshalHelper2(t, "20161207160000.000[-5:EST]610900.500[-9:BST]", &d, &overwritten, eq) // extra part intentionally different to ensure the first timezone is parsed

	// Make sure we properly handle ending newlines
	d = ofxgo.Date(time.Date(2018, 11, 1, 23, 59, 58, 0, EST))
	unmarshalHelper2(t, "20181101235958.000[-5:EST]\n", &d, &overwritten, eq)
}

func TestMarshalString(t *testing.T) {
	var s ofxgo.String = ""
	marshalHelper(t, "", &s)
	s = "foo&bar"
	marshalHelper(t, "foo&amp;bar", &s)
	s = "\n"
	marshalHelper(t, "&#xA;", &s)
	s = "Some Name"
	marshalHelper(t, "Some Name", &s)
}

func TestUnmarshalString(t *testing.T) {
	var s, overwritten ofxgo.String = "", ""
	unmarshalHelper(t, "", &s, &overwritten)
	s = "foo&bar"
	unmarshalHelper(t, "foo&amp;bar", &s, &overwritten)
	// whitespace intentionally stripped because some OFX servers add newlines
	// inside tags
	s = "new\nline"
	unmarshalHelper(t, " new&#xA;line&#xA;", &s, &overwritten)
	s = "Some Name"
	unmarshalHelper(t, "Some Name", &s, &overwritten)
}

func TestMarshalBoolean(t *testing.T) {
	var b ofxgo.Boolean = true
	marshalHelper(t, "Y", &b)
	b = false
	marshalHelper(t, "N", &b)
}

func TestUnmarshalBoolean(t *testing.T) {
	var b, overwritten ofxgo.Boolean = true, false
	unmarshalHelper(t, "Y", &b, &overwritten)
	b = false
	unmarshalHelper(t, "N", &b, &overwritten)
}

func TestMarshalUID(t *testing.T) {
	var u ofxgo.UID = "d1cf3d3d-9ef9-4a97-b180-81706829cb04"
	marshalHelper(t, "d1cf3d3d-9ef9-4a97-b180-81706829cb04", &u)
}

func TestUnmarshalUID(t *testing.T) {
	var u, overwritten ofxgo.UID = "d1cf3d3d-9ef9-4a97-b180-81706829cb04", ""
	unmarshalHelper(t, "d1cf3d3d-9ef9-4a97-b180-81706829cb04", &u, &overwritten)
}

func TestUIDRecommendedFormat(t *testing.T) {
	var u ofxgo.UID = "d1cf3d3d-9ef9-4a97-b180-81706829cb04"
	if ok, err := u.RecommendedFormat(); !ok || err != nil {
		t.Fatalf("UID unexpectedly failed validation\n")
	}
	u = "d1cf3d3d-9ef9-4a97-b180-81706829cb0"
	if ok, err := u.RecommendedFormat(); ok || err == nil {
		t.Fatalf("UID should have failed validation because it's too short\n")
	}
	u = "d1cf3d3d-9ef94a97-b180-81706829cb04"
	if ok, err := u.RecommendedFormat(); ok || err == nil {
		t.Fatalf("UID should have failed validation because it's missing hyphens\n")
	}
	u = "d1cf3d3d-9ef9-4a97-b180981706829cb04"
	if ok, err := u.RecommendedFormat(); ok || err == nil {
		t.Fatalf("UID should have failed validation because its hyphens have been replaced\n")
	}
}
