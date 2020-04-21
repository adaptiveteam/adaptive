package reimpl

import (
	// "bytes"
	"context"
	// "errors"
	// "net/http"
	"net/url"
	// "strings"
	// "sync"
)

// Exchange1 converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The provided context optionally controls which HTTP client is used. See the HTTPClient variable.
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
//
// Opts may include the PKCE verifier code if previously used in AuthCodeURL.
// See https://www.oauth.com/oauth2-servers/pkce/ for more info.
func (c *Config) Exchange1(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	return retrieveToken(ctx, c, v)
}
