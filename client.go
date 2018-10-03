package ofxgo

import (
	"io"
	"net/http"
	"strings"
)

// Client serves to aggregate OFX client settings that may be necessary to talk
// to a particular server due to quirks in that server's implementation.
// Client also provides the Request and RequestNoParse helper methods to aid in
// making and parsing requests.
type Client interface {
	// Used to fill out a Request object
	OfxVersion() ofxVersion
	ID() String
	Version() String
	IndentRequests() bool

	// Request marshals a Request object into XML, makes an HTTP request
	// against it's URL, and then unmarshals the response into a Response
	// object.
	//
	// Before being marshaled, some of the the Request object's values are
	// overwritten, namely those dictated by the BasicClient's configuration
	// (Version, AppID, AppVer fields), and the client's current time
	// (DtClient). These are updated in place in the supplied Request object so
	// they may later be inspected by the caller.
	Request(r *Request) (*Response, error)

	// RequestNoParse marshals a Request object into XML, makes an HTTP
	// request, and returns the raw HTTP response. Unlike RawRequest(), it
	// takes client settings into account. Unlike Request(), it doesn't parse
	// the response into  an ofxgo.Request object.
	//
	// Caveat: The caller is responsible for closing the http Response.Body
	// (see the http module's documentation for more information)
	RequestNoParse(r *Request) (*http.Response, error)

	// RawRequest is little more than a thin wrapper around http.Post
	//
	// In most cases, you should probably be using Request() instead, but
	// RawRequest can be useful if you need to read the raw unparsed http
	// response yourself (perhaps for downloading an OFX file for use by an
	// external program, or debugging server behavior), or have a handcrafted
	// request you'd like to try.
	//
	// Caveats: RawRequest does *not* take client settings into account as
	// Client.Request() does, so your particular server may or may not like
	// whatever we read from 'r'. The caller is responsible for closing the
	// http Response.Body (see the http module's documentation for more
	// information)
	RawRequest(URL string, r io.Reader) (*http.Response, error)
}

type clientCreationFunc func(*BasicClient) Client

// GetClient returns a new Client for a given URL. It attempts to find a
// specialized client for this URL, but simply returns the passed-in
// BasicClient if no such match is found.
func GetClient(URL string, bc *BasicClient) Client {
	clients := []struct {
		URL  string
		Func clientCreationFunc
	}{
		{"https://vesnc.vanguard.com/us/OfxDirectConnectServlet", NewVanguardClient},
	}
	for _, client := range clients {
		if client.URL == strings.Trim(URL, "/") {
			return client.Func(bc)
		}
	}
	return bc
}

// clientRequestNoParse can be used for building clients' RequestNoParse
// methods if they require fairly standard behavior
func clientRequestNoParse(c Client, r *Request) (*http.Response, error) {
	r.SetClientFields(c)

	b, err := r.Marshal()
	if err != nil {
		return nil, err
	}

	return c.RawRequest(r.URL, b)
}

// clientRequest can be used for building clients' Request methods if they
// require fairly standard behavior
func clientRequest(c Client, r *Request) (*Response, error) {
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
