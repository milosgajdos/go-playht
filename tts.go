package playht

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/milosgajdos/go-playht/request"
)

// CreateTTSReq is used to create a new TTS.
type CreateTTSReq struct {
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
	VoiceGuidance uint8        `json:"voice_guidance,omitempty"`
	StyleGuidance uint8        `json:"style_guidance,omitempty"`
}

// TTSJob is returned when a new TTS async job has been created.
type TTSJob struct {
	ID      string        `json:"id"`
	Created time.Time     `json:"created"`
	Input   *CreateTTSReq `json:"input"`
	Output  struct {
		Size     int     `json:"size"`
		URL      string  `json:"url"`
		Duration float64 `json:"duration"`
	} `json:"output"`
	// NOTE: this does not seem to work in line with the docs
	//Links []string `json:"_links,omitempty"`
}

type VoiceEngine string

const (
	PlayHTv2 VoiceEngine = "PlayHT2.0"
	PlayHTv1 VoiceEngine = "PlayHT1.0"
)

func (v VoiceEngine) String() string {
	return string(v)
}

type OutputFormat string

const (
	Mp3   OutputFormat = "mp3"
	Wav   OutputFormat = "wav"
	Ogg   OutputFormat = "ogg"
	Flac  OutputFormat = "flac"
	Mulaw OutputFormat = "mulaw"
)

func (o OutputFormat) String() string {
	return string(o)
}

type Quality string

const (
	Draft   Quality = "draft"
	Low     Quality = "low"
	Medium  Quality = "medium"
	High    Quality = "high"
	Premium Quality = "premium"
)

func (q Quality) String() string {
	return string(q)
}

type Emotion string

const (
	FemaleHappy     Emotion = "female_happy"
	FemaleSad       Emotion = "female_sad"
	FemaleAngry     Emotion = "female_angry"
	FemaleFearful   Emotion = "female_fearful"
	FemaleDisgust   Emotion = "female_disgust"
	FemaleSurprised Emotion = "female_surprised"
	MaleHappy       Emotion = "male_happy"
	MaleSad         Emotion = "male_sad"
	MaleAngry       Emotion = "male_angry"
	MaleFearful     Emotion = "male_fearful"
	MaleDisgust     Emotion = "male_disgust"
	MaleSurprised   Emotion = "male_surprised"
)

func (e Emotion) String() string {
	return string(e)
}

// CreateTTSJob creates a new Text-to-Speech (TTS) job that converts input text into audio asynchronously
func (c *Client) CreateTTSJob(ctx context.Context, createReq *CreateTTSReq) (*TTSJob, error) {
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

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		ttsResp := new(TTSJob)
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

	resp, err := request.Do[APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		ttsResp := new(TTSJob)
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
		request.WithAddHeader("Accept", "audio/mpeg"),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
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

// CreateTTSJobWithProgressStream creates a new Text-to-Speech (TTS) SSE stream that converts input text into audio
// asynchronously and returns the job progress SSE stream URL.
// nolint:revive
func (c *Client) CreateTTSJobWithProgressStream(ctx context.Context, req *CreateTTSReq) (string, error) {
	panic("not implemented")
}

// GetTTSJobProgressStream retrieves the TTS job progress SSE stream from the job with the given id.
// nolint:revive
func (c *Client) GetTTSJobProgressStream(ctx context.Context, w io.Writer, id string) error {
	panic("not implemented")
}
