package ofxgo

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strings"
)

// VanguardClient provides a Client implementation which handles Vanguard's
// cookie-passing requirements and also enables older, disabled-by-default
// cipher suites. VanguardClient uses default, non-zero settings, if its fields
// are not initialized.
type VanguardClient struct {
	*BasicClient
}

// NewVanguardClient returns a Client interface configured to handle Vanguard's
// brand of idiosyncrasy
func NewVanguardClient(bc *BasicClient) Client {
	return &VanguardClient{bc}
}

// vanguardHttpClient returns an http.Client with the default supported
// ciphers plus the insecure ciphers Vanguard still uses.
func vanguardHttpClient() *http.Client {
	var clientCiphers []uint16

	vanguardCiphers := []uint16{
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	}
	defaultCiphers := tls.CipherSuites()
	for _, cipher := range defaultCiphers {
		clientCiphers = append(clientCiphers, cipher.ID)
	}
	clientCiphers = append(clientCiphers, vanguardCiphers...)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				CipherSuites: clientCiphers,
			},
		},
	}
	return client
}

// RawRequest is a copy of BasicClient RawRequest with a custom http.Client
// that enables older cipher suites.
func (c *VanguardClient) RawRequest(URL string, r io.Reader) (*http.Response, error) {
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
	client := vanguardHttpClient()
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return response, errors.New("OFXQuery request status: " + response.Status)
	}

	return response, nil
}

// rawRequestCookiesInsecureCiphers is RawRequest with the added features of
// sending cookies and using a custom http.Client that enables older cipher
// suites which are disabled by default
func rawRequestCookiesInsecureCiphers(URL string, r io.Reader, cookies []*http.Cookie) (*http.Response, error) {
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

	client := vanguardHttpClient()
	response, err := client.Do(request)
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

		return rawRequestCookiesInsecureCiphers(r.URL, b, response.Cookies())
	}

	return response, err
}

// Request marshals a Request to XML, makes an HTTP request, and then
// unmarshals the response into a Response object.
func (c *VanguardClient) Request(r *Request) (*Response, error) {
	return clientRequest(c, r)
}
