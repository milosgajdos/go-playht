package request

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPRequest(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		t.Parallel()
		// nolint:staticcheck
		req, err := NewHTTP(nil, http.MethodGet, "http://foo.com", nil)
		assert.NoError(t, err)
		assert.NotNil(t, req.Context())
	})
	t.Run("nil body", func(t *testing.T) {
		t.Parallel()
		req, err := NewHTTP(context.TODO(), http.MethodGet, "http://foo.com", nil)
		assert.NoError(t, err)
		assert.NotNil(t, req.Body)
	})
	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		secret := "secret"
		options := []HTTPOption{
			WithAuthSecret(secret),
		}
		req, err := NewHTTP(context.TODO(), http.MethodGet, "http://foo.com", &bytes.Reader{}, options...)
		assert.NoError(t, err)
		assert.NotNil(t, req.Body)

		// check all default headers are set as well as the bearer one
		header := make(http.Header)
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Authorization", secret)
		// NOTE: this header is set by default
		// on every request we create via NewHTTP.
		header.Set("User-Agent", UserAgent)
		assert.Equal(t, req.Header, header)
	})
}

func TestHTTPReqOption(t *testing.T) {
	t.Parallel()
	t.Run("set bearer", func(t *testing.T) {
		t.Parallel()
		req := &http.Request{}
		token := "token"
		WithBearer(token)(req)
		assert.Equal(t, req.Header.Get("Authorization"), fmt.Sprintf("Bearer %s", token))
	})
	t.Run("set authz secret", func(t *testing.T) {
		t.Parallel()
		req := &http.Request{}
		secret := "secret"
		WithAuthSecret(secret)(req)
		assert.Equal(t, req.Header.Get("Authorization"), secret)
	})
	t.Run("set header", func(t *testing.T) {
		t.Parallel()
		req := &http.Request{}
		key, val := "foo", "bar"
		WithSetHeader(key, val)(req)
		assert.Equal(t, req.Header.Get(key), val)
	})

	t.Run("add header", func(t *testing.T) {
		t.Parallel()
		key, val := "foo", "bar"
		req := &http.Request{
			Header: make(http.Header),
		}
		req.Header.Add(key, val)
		WithAddHeader(key, val)(req)
		assert.Equal(t, req.Header.Values(key), []string{val, val})
	})
}
