package ofxgo

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestBasicClient_HTTPClient(t *testing.T) {
	c := &BasicClient{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return nil, errors.New("bad test client")
				},
			},
		},
	}
	_, err := c.Request(&Request{
		URL: "https://test",
		Signon: SignonRequest{
			UserID:   "test",
			UserPass: "test",
		},
	})
	if err == nil || !strings.Contains(err.Error(), "bad test client") {
		t.Fatalf("expected error containing 'bad test client', got: %v", err)
	}
}
