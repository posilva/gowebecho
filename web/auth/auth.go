package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

//OktaProvider olds the data that allows to perform a Okta API request
type OktaProvider struct {
	CustomURL string
	Client    *http.Client
	Values    url.Values
}

type OktaConfig struct {
	Domain       string
	ClientId     string
	ClientSecret string
	RedirectURI  string
	Scope        string
	Nonce        string
	State        string
	ResponseType string
	ResponseMode string
}

func defaultOktaConfig() OktaConfig {

	return OktaConfig{
		Domain:       os.Getenv("AUTH_DOMAIN"),
		ClientId:     os.Getenv("AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("AUTH_REDIRECT_URL"),
		Scope:        os.Getenv("AUTH_SCOPE"),
		Nonce:        os.Getenv("AUTH_NONCE"),
		State:        os.Getenv("AUTH_STATE"),
		ResponseType: os.Getenv("AUTH_RESPONSE_TYPE"),
		ResponseMode: os.Getenv("AUTH_RESPONSE_MODE"),
	}
}

// NewOktaProviderWithConfig creates a new instance of Okta provider based on the configuration provided
func NewOktaProviderWithConfig(c OktaConfig) (*OktaProvider, error) {

	baseURL := fmt.Sprintf("https://%s/oauth2/default/v1/authorize", c.Domain)

	vals := url.Values{}

	vals.Add("client_id", c.ClientId)
	vals.Add("redirect_uri", c.RedirectURI)
	vals.Add("response_type", c.ResponseType)
	vals.Add("response_mode", c.ResponseMode)
	vals.Add("scope", c.Scope)
	vals.Add("state", c.State)
	vals.Add("nonce", c.Nonce)

	// no premature optimization: https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	return &OktaProvider{
		CustomURL: baseURL,
		Client:    client,
		Values:    vals,
	}, nil
}

// NewOktaProvider creates an Okta provider with default config
func NewOktaProvider() (*OktaProvider, error) {
	return NewOktaProviderWithConfig(defaultOktaConfig())
}

// Authorize executes a authorize method on Okta OpenID connect API
// see: https://developer.okta.com/docs/reference/api/oidc/#authorize
func (p *OktaProvider) Authorize(c echo.Context, state string, nonce string) error {
	v := p.Values
	v.Set("nonce", nonce)
	v.Set("state", state)

	req, err := http.NewRequest(http.MethodGet, p.CustomURL, nil)
	if err != nil {
		c.Logger().Errorf("%v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to create request: %v", err))
	}

	req.URL.RawQuery = v.Encode()
	return c.Redirect(http.StatusMovedPermanently, req.URL.String())
}
