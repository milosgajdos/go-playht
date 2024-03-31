package playht

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-playht/request"
)

// CreateTTSStreamReq is used to create TTS stream.
type CreateTTSStreamReq struct {
	Text          string       `json:"text"`
	Voice         string       `json:"voice"`
	Quality       Quality      `json:"quality"`
	OutputFormat  OutputFormat `json:"output_format,omitempty"`
	VoiceEngine   VoiceEngine  `json:"voice_engine,omitempty"`
	Emotion       Emotion      `json:"emotion,omitempty"`
	SampleRate    int32        `json:"sample_rate"`
	Seed          uint8        `json:"seed,omitempty"`
	VoiceGuidance uint8        `json:"voice_guidance,omitempty"`
	StyleGuidance uint8        `json:"style_guidance,omitempty"`
	TextGuidance  uint8        `json:"text_guidance,omitempty"`
	Temperature   float32      `json:"temperature,omitempty"`
	Speed         float32      `json:"speed"`
}

// TTSStreamURL is returned when the stream URL is requested.
type TTSStreamURL struct {
	HRef   string `json:"href"`
	Method string `json:"method"`
	CType  string `json:"contentType"`
	Rel    string `json:"rel"`
	Desc   string `json:"description"`
}

// TTSStream creates a new TTS stream and streams the audio bytes immediately.
func (c *Client) TTSStream(ctx context.Context, w io.Writer, createReq *CreateTTSStreamReq) error {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts/stream")
	if err != nil {
		return err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(createReq); err != nil {
		return err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithSetHeader("Content-Type", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return err
	}

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if _, err := io.Copy(w, resp.Body); err != nil {
			return err
		}
		return nil
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		var apiErr APIErrInternal
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return errors.Join(err, jsonErr)
		}
		return apiErr
	default:
		return err
	}
}

// TTSStreamURL creates a new TTS stream and returns data containing an URL that is immediately streamable.
func (c *Client) TTSStreamURL(ctx context.Context, createReq *CreateTTSStreamReq) (*TTSStreamURL, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts/stream")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(createReq); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
		request.WithSetHeader("Content-Type", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
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
		ttsResp := new(TTSStreamURL)
		if err := json.NewDecoder(resp.Body).Decode(ttsResp); err != nil {
			return nil, err
		}
		return ttsResp, nil
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
