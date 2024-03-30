package playht

import (
	"context"
	"net/http"
)

// HTTP is an HTTP client.
type HTTP struct {
	client  *http.Client
	limiter Limiter
}

// Options configure the HTTP client.
type Options struct {
	HTTPClient *http.Client
	Limiter    Limiter
}

// Option is functional option.
type Option func(*Options)

// Limiter is used to apply API rate limits.
// NOTE: you can use off the shelf limiter from
// https://pkg.go.dev/golang.org/x/time/rate#Limiter
type Limiter interface {
	// Wait must block until limiter
	// permits another request to proceed.
	Wait(context.Context) error
}

// DefaultTransport returns a new http.Transport
// which is a clione of the http.DefaultTransport
// This is to avoid accidental transport overrides.
func DefaultTransport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	return transport
}

// NewHTTP creates a new HTTP client and returns it.
func NewHTTP(opts ...Option) *HTTP {
	options := Options{
		HTTPClient: &http.Client{
			Transport: DefaultTransport(),
		},
	}
	for _, apply := range opts {
		apply(&options)
	}

	return &HTTP{
		client:  options.HTTPClient,
		limiter: options.Limiter,
	}
}

// Do dispatches the HTTP request to the remote endpoint.
func (h *HTTP) Do(req *http.Request) (*http.Response, error) {
	if h.limiter != nil {
		err := h.limiter.Wait(req.Context()) // This is a blocking call which honors the rate limit
		if err != nil {
			return nil, err
		}
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WithHTTPClient sets the HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(o *Options) {
		o.HTTPClient = c
	}
}

// WithLimiter sets the http rate limiter.
func WithLimiter(l Limiter) Option {
	return func(o *Options) {
		o.Limiter = l
	}
}
