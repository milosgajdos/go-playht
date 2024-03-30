package playht

import (
	"os"

	"github.com/milosgajdos/go-playht/client"
)

const (
	// BaseURL is OpenAI HTTP API base URL.
	BaseURL = "https://api.play.ht/api"
	// APIV2 V2 version.
	APIV2 = "v2"
	// APIV2 V1 version.
	APIV1 = "v1"
	// UserIDHeader
	UserIDHeader = "X-USER-ID"
)

// Client is an OpenAI HTTP API client.
type Client struct {
	opts Options
}

type Options struct {
	SecretKey  string
	UserID     string
	BaseURL    string
	Version    string
	HTTPClient *client.HTTP
}

// Option is functional graph option.
type Option func(*Options)

// NewClient creates a new HTTP API client and returns it.
// By default it reads the secret key from PLAYHT_SECRET_KEY env var
// and user ID from PLAYHT_USER_ID env var and uses
// the default http client for making the HTTP api requests.
func NewClient(opts ...Option) *Client {
	options := Options{
		SecretKey:  os.Getenv("PLAYHT_SECRET_KEY"),
		UserID:     os.Getenv("PLAYHT_USER_ID"),
		BaseURL:    BaseURL,
		Version:    APIV2,
		HTTPClient: client.NewHTTP(),
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &Client{
		opts: options,
	}
}

// WithSecretKey sets the secret key.
func WithSecretKey(apiKey string) Option {
	return func(o *Options) {
		o.SecretKey = apiKey
	}
}

// WithUserID sets the user ID.
func WithUserID(userID string) Option {
	return func(o *Options) {
		o.UserID = userID
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

// WithVersion sets the API version.
func WithVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(httpClient *client.HTTP) Option {
	return func(o *Options) {
		o.HTTPClient = httpClient
	}
}
