package ofxgo

import (
	"bytes"
	"github.com/aclindsa/xml"
)

// Returns the next available Token from the xml.Decoder that is not CharData
// made up entirely of whitespace. This is useful to skip whitespace when
// manually unmarshaling XML.
func nextNonWhitespaceToken(decoder *xml.Decoder) (xml.Token, error) {
	for {
		tok, err := decoder.Token()
		if err != nil {
			return nil, err
		} else if chars, ok := tok.(xml.CharData); ok {
			strippedBytes := bytes.TrimSpace(chars)
			if len(strippedBytes) != 0 {
				return tok, nil
			}
		} else {
			return tok, nil
		}
	}
}
