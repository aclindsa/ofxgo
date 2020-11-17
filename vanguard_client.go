package ofxgo

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

// VanguardClient provides a Client implementation which handles Vanguard's
// cookie-passing requirements. VanguardClient uses default, non-zero settings,
// if its fields are not initialized.
type VanguardClient struct {
	*BasicClient
}

// NewVanguardClient returns a Client interface configured to handle Vanguard's
// brand of idiosyncrasy
func NewVanguardClient(bc *BasicClient) Client {
	return &VanguardClient{bc}
}

// rawRequestCookies is RawRequest with the added feature of sending cookies
func rawRequestCookies(URL string, r io.Reader, cookies []*http.Cookie) (*http.Response, error) {
	if !strings.HasPrefix(URL, "https://") {
		return nil, errors.New("Refusing to send OFX request with possible plain-text password over non-https protocol")
	}

	request, err := http.NewRequest("POST", URL, r)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-ofx")
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	response, err := http.DefaultClient.Do(request)
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
func (c *VanguardClient) RequestNoParse(r *Request) (*http.Response, error) {
	r.SetClientFields(c)

	b, err := r.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := c.RawRequest(r.URL, b)

	// Some financial institutions (cough, Vanguard, cough), require a cookie
	// to be set on the http request, or they return empty responses.
	// Fortunately, the initial response contains the cookie we need, so if we
	// detect an empty response with cookies set that didn't have any errors,
	// re-try the request while sending their cookies back to them.
	if response != nil && response.ContentLength <= 0 && len(response.Cookies()) > 0 {
		b, err = r.Marshal()
		if err != nil {
			return nil, err
		}

		return rawRequestCookies(r.URL, b, response.Cookies())
	}

	return response, err
}

// Request marshals a Request to XML, makes an HTTP request, and then
// unmarshals the response into a Response object.
func (c *VanguardClient) Request(r *Request) (*Response, error) {
	return clientRequest(c, r)
}
