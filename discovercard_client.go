package ofxgo

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// DiscoverCardClient provides a Client implementation which handles
// DiscoverCard's broken HTTP header behavior. DiscoverCardClient uses default,
// non-zero settings, if its fields are not initialized.
type DiscoverCardClient struct {
	*BasicClient
}

// NewDiscoverCardClient returns a Client interface configured to handle
// Discover Card's brand of idiosyncracy
func NewDiscoverCardClient(bc *BasicClient) Client {
	return &DiscoverCardClient{bc}
}

func discoverCardHTTPPost(URL string, r io.Reader) (*http.Response, error) {
	// Either convert or copy to a bytes.Buffer to be able to determine the
	// request length for the Content-Length header
	buf, ok := r.(*bytes.Buffer)
	if !ok {
		buf = &bytes.Buffer{}
		_, err := io.Copy(buf, r)
		if err != nil {
			return nil, err
		}
	}

	url, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	path := url.Path
	if path == "" {
		path = "/"
	}

	// Discover requires only these headers and in this exact order, or it
	// returns HTTP 403
	headers := fmt.Sprintf("POST %s HTTP/1.1\r\n"+
		"Content-Type: application/x-ofx\r\n"+
		"Host: %s\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: Keep-Alive\r\n"+
		"\r\n", path, url.Hostname(), buf.Len())

	host := url.Host
	if url.Port() == "" {
		host += ":443"
	}

	// BUGBUG: cannot do defer conn.Close() until body is read,
	// we are "leaking" a socket here, but it will be finalized
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		return nil, err
	}

	fmt.Fprint(conn, headers)
	_, err = io.Copy(conn, buf)
	if err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(conn), nil)
}

// RawRequest is a convenience wrapper around http.Post. It is exposed only for
// when you need to read/inspect the raw HTTP response yourself.
func (c *DiscoverCardClient) RawRequest(URL string, r io.Reader) (*http.Response, error) {
	if !strings.HasPrefix(URL, "https://") {
		return nil, errors.New("Refusing to send OFX request with possible plain-text password over non-https protocol")
	}

	response, err := discoverCardHTTPPost(URL, r)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("OFXQuery request status: " + response.Status)
	}

	return response, nil
}

// RequestNoParse marshals a Request to XML, makes an HTTP request, and returns
// the raw HTTP response
func (c *DiscoverCardClient) RequestNoParse(r *Request) (*http.Response, error) {
	return clientRequestNoParse(c, r)
}

// Request marshals a Request to XML, makes an HTTP request, and then
// unmarshals the response into a Response object.
func (c *DiscoverCardClient) Request(r *Request) (*Response, error) {
	return clientRequest(c, r)
}
