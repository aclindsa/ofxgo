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
	// Use carriage returns on new lines
	CarriageReturn bool
	// Set User-Agent header to this string, if not empty
	UserAgent string
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

// CarriageReturnNewLines returns true if carriage returns should be used on new lines, false otherwise
func (c *BasicClient) CarriageReturnNewLines() bool {
	return c.CarriageReturn
}

// RawRequest is a convenience wrapper around http.Post. It is exposed only for
// when you need to read/inspect the raw HTTP response yourself.
func (c *BasicClient) RawRequest(URL string, r io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(URL, "https://") {
		return nil, errors.New("Refusing to send OFX request with possible plain-text password over non-https protocol")
	}

	request, err := http.NewRequest("POST", URL, r)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-ofx")
	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return response, errors.New("OFXQuery request status: " + response.Status)
	}

	return response, nil
}

// RequestNoParse marshals a Request to XML, makes an HTTP request, and returns
// the raw HTTP response
func (c *BasicClient) RequestNoParse(r *Request) (*http.Response, error) {
	return clientRequestNoParse(c, r)
}

// Request marshals a Request to XML, makes an HTTP request, and then
// unmarshals the response into a Response object.
func (c *BasicClient) Request(r *Request) (*Response, error) {
	return clientRequest(c, r)
}
