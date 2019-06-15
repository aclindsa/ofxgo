package ofxgo_test

import (
	"github.com/aclindsa/ofxgo"
	"regexp"
	"strings"
	"testing"
)

// match leading and trailing whitespace on each line
var ignoreSpacesRe = regexp.MustCompile("(?m)^[ \t]+|$[\r\n]+")

func marshalCheckRequest(t *testing.T, request *ofxgo.Request, expected string) {
	t.Helper()
	buf, err := request.Marshal()
	if err != nil {
		t.Fatalf("%s: Unexpected error marshalling request: %s\n", t.Name(), err)
	}
	actualString := buf.String()

	// Ignore spaces between XML elements
	expectedString := ignoreSpacesRe.ReplaceAllString(expected, "")
	actualString = ignoreSpacesRe.ReplaceAllString(actualString, "")

	if expectedString != actualString {
		compareLength := len(expectedString)
		if len(actualString) < compareLength {
			compareLength = len(actualString)
		}

		for i := 0; i < compareLength; i++ {
			if expectedString[i] != actualString[i] {
				firstDifferencePosition := 13
				displayStart := i - 10
				prefix := "..."
				suffix := "..."
				if displayStart < 0 {
					prefix = ""
					firstDifferencePosition = i
					displayStart = 0
				}
				displayEnd := displayStart + 40
				if displayEnd > compareLength {
					suffix = ""
					displayEnd = compareLength
				}
				t.Fatalf("%s expected '%s%s%s',\ngot '%s%s%s'\n     %s^ first difference\n", t.Name(), prefix, expectedString[displayStart:displayEnd], suffix, prefix, actualString[displayStart:displayEnd], suffix, strings.Repeat(" ", firstDifferencePosition))
			}
		}

		if len(actualString) > compareLength {
			t.Fatalf("%s: Actual string longer than expected string\n", t.Name())
		} else {
			t.Fatalf("%s: Actual string shorter than expected string\n", t.Name())
		}
	}
}
