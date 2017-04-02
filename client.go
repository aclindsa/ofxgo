package ofxgo

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	// Request fields to overwrite with the client's values. If nonempty,
	// defaults are used
	SpecVersion string // VERSION in header
	AppId       string // SONRQ>APPID
	AppVer      string // SONRQ>APPVER

	// Don't insert newlines or indentation when marshalling to SGML/XML
	NoIndent bool
}

var defaultClient Client

func (c *Client) OfxVersion() string {
	if len(c.SpecVersion) > 0 {
		return c.SpecVersion
	} else {
		return "203"
	}
}

func (c *Client) Id() String {
	if len(c.AppId) > 0 {
		return String(c.AppId)
	} else {
		return String("OFXGO")
	}
}

func (c *Client) Version() String {
	if len(c.AppVer) > 0 {
		return String(c.AppVer)
	} else {
		return String("0001")
	}
}

// Returns true if the marshaled XML should be indented (and contain newlines,
// since the two are linked in the current implementation)
func (c *Client) IndentRequests() bool {
	return !c.NoIndent
}

// RawRequest is little more than a thin wrapper around http.Post
//
// In most cases, you should probably be using Request() instead, but
// RawRequest can be useful if you need to read the raw unparsed http response
// yourself (perhaps for downloading an OFX file for use by an external
// program, or debugging server behavior), or have a handcrafted request you'd
// like to try.
//
// Caveats: RawRequest does *not* take client settings into account as
// Request() does, so your particular server may or may not like whatever we
// read from 'r'. The caller is responsible for closing the http Response.Body
// (see the http module's documentation for more information)
func RawRequest(URL string, r io.Reader) (*http.Response, error) {
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

// RequestNoParse marshals a Request object into XML, makes an HTTP request,
// and returns the raw HTTP response. Unlike RawRequest(), it takes client
// settings into account. Unlike Request(), it doesn't parse the response into
// a Request object.
//
// Caveat: The caller is responsible for closing the http Response.Body (see
// the http module's documentation for more information)
func (c *Client) RequestNoParse(r *Request) (*http.Response, error) {
	r.SetClientFields(c)

	b, err := r.Marshal()
	if err != nil {
		return nil, err
	}

	return RawRequest(r.URL, b)
}

// Request marshals a Request object into XML, makes an HTTP request against
// it's URL, and then unmarshals the response into a Response object.
//
// Before being marshaled, some of the the Request object's values are
// overwritten, namely those dictated by the Client's configuration (Version,
// AppId, AppVer fields), and the client's curren time (DtClient). These are
// updated in place in the supplied Request object so they may later be
// inspected by the caller.
func (c *Client) Request(r *Request) (*Response, error) {
	response, err := c.RequestNoParse(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	ofxresp, err := ParseResponse(response.Body)
	if err != nil {
		return nil, err
	}
	return ofxresp, nil
}
