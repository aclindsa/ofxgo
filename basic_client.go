package ofxgo

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

// BasicClient provides a standard Client implementation suitable for most
// financial institutions. BasicClient uses default, non-zero settings, even if
// its fields are not initialized.
type BasicClient struct {
	// Request fields to overwrite with the client's values. If nonempty,
	// defaults are used
	SpecVersion ofxVersion // VERSION in header
	AppID       string     // SONRQ>APPID
	AppVer      string     // SONRQ>APPVER

	// Don't insert newlines or indentation when marshalling to SGML/XML
	NoIndent bool
}

// OfxVersion returns the OFX specification version this BasicClient will marshal
// Requests as. Defaults to "203" if the client's SpecVersion field is empty.
func (c *BasicClient) OfxVersion() ofxVersion {
	if c.SpecVersion.Valid() {
		return c.SpecVersion
	}
	return OfxVersion203
}

// ID returns this BasicClient's OFX AppID field, defaulting to "OFXGO" if
// unspecified.
func (c *BasicClient) ID() String {
	if len(c.AppID) > 0 {
		return String(c.AppID)
	}
	return String("OFXGO")
}

// Version returns this BasicClient's version number as a string, defaulting to
// "0001" if unspecified.
func (c *BasicClient) Version() String {
	if len(c.AppVer) > 0 {
		return String(c.AppVer)
	}
	return String("0001")
}

// IndentRequests returns true if the marshaled XML should be indented (and
// contain newlines, since the two are linked in the current implementation)
func (c *BasicClient) IndentRequests() bool {
	return !c.NoIndent
}

func (c *BasicClient) RawRequest(URL string, r io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(URL, "https://") {
		return nil, errors.New("Refusing to send OFX request with possible plain-text password over non-https protocol")
	}

	response, err := http.Post(URL, "application/x-ofx", r)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("OFXQuery request status: " + response.Status)
	}

	return response, nil
}

func (c *BasicClient) RequestNoParse(r *Request) (*http.Response, error) {
	return clientRequestNoParse(c, r)
}

func (c *BasicClient) Request(r *Request) (*Response, error) {
	return clientRequest(c, r)
}
