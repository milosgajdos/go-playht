package playht

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/milosgajdos/go-playht/request"
)

// CreateTTSJobReq is used to create a new TTS.
type CreateTTSJobReq struct {
	Text          string       `json:"text"`
	Voice         string       `json:"voice"`
	Quality       Quality      `json:"quality"`
	OutputFormat  OutputFormat `json:"output_format"`
	VoiceEngine   VoiceEngine  `json:"voice_engine,omitempty"`
	Emotion       Emotion      `json:"emotion,omitempty"`
	Speed         float32      `json:"speed"`
	Temperature   float32      `json:"temperature,omitempty"`
	SampleRate    int32        `json:"sample_rate"`
	Seed          uint8        `json:"seed,omitempty"`
	VoiceGuidance float32      `json:"voice_guidance,omitempty"`
	StyleGuidance float32      `json:"style_guidance,omitempty"`
}

type Link struct {
	ContentType string `json:"content_type,omitempty"`
	Description string `json:"description,omitempty"`
	Href        string `json:"href,omitempty"`
	Method      string `json:"method,omitempty"`
	Rel         string `json:"rel,omitempty"`
}

// TTSJob is returned when a new TTS async job has been created.
type TTSJob struct {
	ID      string           `json:"id"`
	Created time.Time        `json:"created"`
	Input   *CreateTTSJobReq `json:"input"`
	Output  struct {
		Size     int     `json:"size"`
		URL      string  `json:"url"`
		Duration float64 `json:"duration"`
	} `json:"output"`
	Status string `json:"status,omitempty"`
	Links  []Link `json:"_links,omitempty"`
}

// CreateTTSJob creates a new Text-to-Speech (TTS) job that converts input text into audio asynchronously
func (c *Client) CreateTTSJob(ctx context.Context, createReq *CreateTTSJobReq) (*TTSJob, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts")
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

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ttsResp := new(TTSJob)
	if err := json.NewDecoder(resp.Body).Decode(ttsResp); err != nil {
		return nil, err
	}
	return ttsResp, nil
}

// GetTTSJob retrieves information about an async text-to-speech job.
func (c *Client) GetTTSJob(ctx context.Context, id string) (*TTSJob, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts/" + id)
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ttsResp := new(TTSJob)
	if err := json.NewDecoder(resp.Body).Decode(ttsResp); err != nil {
		return nil, err
	}
	return ttsResp, nil
}

// GetTTSJobAudioStream retrieves the TTS job audio stream from the job with the given id.
// It streams audio in the MP3 format or returns error if the file was not generated as MP3.
func (c *Client) GetTTSJobAudioStream(ctx context.Context, w io.Writer, id string) error {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts/" + id)
	if err != nil {
		return err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return nil
}

// CreateTTSJobWithProgressStream creates a new Text-to-Speech (TTS) SSE stream that converts input text into audio
// asynchronously and returns the job progress SSE stream URL. If w is not nil, the events are streamed into it.
func (c *Client) CreateTTSJobWithProgressStream(ctx context.Context, w io.Writer, createReq *CreateTTSJobReq) (string, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts")
	if err != nil {
		return "", err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(createReq); err != nil {
		return "", err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "text/event-stream"),
		request.WithSetHeader("Content-Type", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return "", err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	streamURL := resp.Header.Get("Content-Location")
	if w != nil {
		if _, err := io.Copy(w, resp.Body); err != nil {
			return streamURL, err
		}
		return streamURL, nil
	}
	return streamURL, nil
}

// GetTTSJobProgressStream retrieves the TTS job progress SSE stream for the job with the given id and streams it into w.
func (c *Client) GetTTSJobProgressStream(ctx context.Context, w io.Writer, id string) error {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/tts/" + id)
	if err != nil {
		return err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "text/event-stream"),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return nil
}
