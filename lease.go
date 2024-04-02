package playht

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/milosgajdos/go-playht/request"
)

const (
	// HTEpoch is the HT Lease epoch.
	// I've no idea why but whatever.
	HTEpoch int64 = 1519257480 // 2018-02-21 23:58:00 UTC
)

// Lease for gRPC stream.
type Lease struct {
	Data     []byte
	Created  time.Time
	Duration time.Duration
	Metadata map[string]any
}

// Expires returns a timestamp when the Lease expires
func (l *Lease) Expires() time.Time {
	return l.Created.Add(l.Duration)
}

// CreateLeaseReq is used to create a nw Lease.
type CreateLeaseReq struct{}

// CreateLease creates a new lease and returns it.
func (c *Client) CreateLease(ctx context.Context, _ *CreateLeaseReq) (*Lease, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/leases")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
		request.WithSetHeader("Content-Type", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var created, duration int32
		buf := bytes.NewReader(data[64:68])
		if err := binary.Read(buf, binary.BigEndian, &created); err != nil {
			return nil, fmt.Errorf("failed reading created data: %v", err)
		}
		buf = bytes.NewReader(data[68:72])
		if err := binary.Read(buf, binary.BigEndian, &duration); err != nil {
			return nil, fmt.Errorf("failed reading duration data: %v", err)
		}
		md := map[string]any{}
		if err := json.Unmarshal(data[72:], &md); err != nil {
			return nil, fmt.Errorf("failed reading lease metadata: %v", err)
		}

		return &Lease{
			Data:     data,
			Created:  time.Unix(HTEpoch, 0).Add(time.Duration(created) * time.Second),
			Duration: time.Duration(duration) * time.Second,
			Metadata: md,
		}, nil
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusInternalServerError:
		var apiErr APIErrInternal
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	default:
		return nil, err
	}
}

// RefreshLease refreshes the existing Lease and returns it.
// nolint:revive
func (c *Client) RefreshLease(ctx context.Context, req *CreateLeaseReq) (*Lease, error) {
	return nil, nil
}
