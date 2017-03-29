package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"regexp"
	"testing"
)

var ignoreSpacesRe = regexp.MustCompile(">[ \t\r\n]+<")

func marshalCheckRequest(t *testing.T, request *ofxgo.Request, expected string) {
	buf, err := request.Marshal()
	if err != nil {
		t.Fatalf("Unexpected error marshalling request: %s\n", err)
	}
	actualString := buf.String()

	// Ignore spaces between XML elements
	expectedString := ignoreSpacesRe.ReplaceAllString(expected, "><")
	actualString = ignoreSpacesRe.ReplaceAllString(actualString, "><")

	if expectedString != actualString {
		compareLength := len(expectedString)
		if len(actualString) < compareLength {
			compareLength = len(actualString)
		}

		for i := 0; i < compareLength; i++ {
			if expectedString[i] != actualString[i] {
				displayStart := i - 10
				if displayStart < 0 {
					displayStart = 0
				}
				displayEnd := displayStart + 40
				if displayEnd > compareLength {
					displayEnd = compareLength
				}
				t.Fatalf("%s expected '...%s...',\ngot '...%s...'\n", t.Name(), expectedString[displayStart:displayEnd], actualString[displayStart:displayEnd])
			}
		}

		if len(actualString) > compareLength {
			t.Fatalf("%s: Actual string longer than expected string\n", t.Name())
		} else {
			t.Fatalf("%s: Actual string shorter than expected string\n", t.Name())
		}
	}
}
