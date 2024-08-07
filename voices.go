package playht

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"github.com/milosgajdos/go-playht/request"
)

// Voice is the stock PlayHT voice.
type Voice struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sample   string `json:"sample"`
	Accent   string `json:"accent"`
	Age      string `json:"age"`
	Gender   string `json:"gender"`
	Language string `json:"language"`
	LangCode string `json:"language_code"`
	Loudness string `json:"loudness"`
	Style    string `json:"style"`
	Tempo    string `json:"tempo"`
	Texture  string `json:"texture"`
}

// ClonedVoice data.
type ClonedVoice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// CloneVoiceFileRequest is used to create a voice clone.
type CloneVoiceFileRequest struct {
	SampleFile string `json:"sample_file"`
	VoiceName  string `json:"voice_name"`
	MimeType   string `json:"mime_type"`
}

// CloneVoiceURLRequest is used to create a voice clone via file URL.
type CloneVoiceURLRequest struct {
	SampleFileURL string `json:"sample_file_url"`
	VoiceName     string `json:"voice_name"`
}

// DeleteVoiceRequest is used to deleted cloned voice.
type DeleteClonedVoiceRequest struct {
	VoiceID string `json:"voice_id"`
}

// DeleteClonedVoiceResp is returned when the cloned voice has been deleted.
type DeleteClonedVoiceResp struct {
	Message string      `json:"message"`
	Deleted ClonedVoice `json:"deleted"`
}

// GetVoices returns the full list of stock PlayHT GetVoices.
func (c *Client) GetVoices(ctx context.Context) ([]Voice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/voices")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
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

	var voices []Voice
	if err := json.NewDecoder(resp.Body).Decode(&voices); err != nil {
		return nil, err
	}
	return voices, nil
}

// GetClonedVoices obtains a list of all cloned voices created by the user.
func (c *Client) GetClonedVoices(ctx context.Context) ([]ClonedVoice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/cloned-voices")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
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

	var voices []ClonedVoice
	if err := json.NewDecoder(resp.Body).Decode(&voices); err != nil {
		return nil, err
	}
	return voices, nil
}

// CreateInstantVoiceCloneFromFile creates an instant voice clone by providing a sample audio file via file upload.
func (c *Client) CreateInstantVoiceCloneFromFile(ctx context.Context, cloneReq *CloneVoiceFileRequest) (*ClonedVoice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/cloned-voices/instant")
	if err != nil {
		return nil, err
	}

	f, err := os.Open(cloneReq.SampleFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := createFilePart(w, "sample_file", cloneReq.SampleFile, cloneReq.MimeType)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}
	fw, err = createFieldPart(w, "voice_name", "text/plain")
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, strings.NewReader(cloneReq.VoiceName)); err != nil {
		return nil, err
	}
	w.Close()

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
		request.WithAddHeader("Content-Type", w.FormDataContentType()),
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

	cloneVoiceResp := new(ClonedVoice)
	if err := json.NewDecoder(resp.Body).Decode(cloneVoiceResp); err != nil {
		return nil, err
	}
	return cloneVoiceResp, nil
}

// CreateInstantVoiceCloneFromURL create an instant voice clone by providing an URL for a sample audio file.
func (c *Client) CreateInstantVoiceCloneFromURL(ctx context.Context, cloneReq *CloneVoiceURLRequest) (*ClonedVoice, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/cloned-voices/instant/")
	if err != nil {
		return nil, err
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	m := map[string]io.Reader{
		"sample_file_url": strings.NewReader(cloneReq.SampleFileURL),
		"voice_name":      strings.NewReader(cloneReq.VoiceName),
	}
	for field, data := range m {
		fw, err := createFieldPart(w, field, "text/plain")
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(fw, data); err != nil {
			return nil, err
		}
	}
	w.Close()

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
		request.WithAddHeader("Content-Type", w.FormDataContentType()),
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

	cloneVoiceResp := new(ClonedVoice)
	if err := json.NewDecoder(resp.Body).Decode(cloneVoiceResp); err != nil {
		return nil, err
	}
	return cloneVoiceResp, nil
}

// DeleteClonedVoice eletes a cloned voice created by the user using the provided voice_id.
func (c *Client) DeleteClonedVoice(ctx context.Context, delReq *DeleteClonedVoiceRequest) (*DeleteClonedVoiceResp, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/cloned-voices/")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(delReq); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithAuthSecret(c.opts.SecretKey),
		request.WithSetHeader(UserIDHeader, c.opts.UserID),
		request.WithAddHeader("Accept", "application/json"),
	}

	req, err := request.NewHTTP(ctx, http.MethodDelete, u.String(), body, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[*APIError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	delVoiceResp := new(DeleteClonedVoiceResp)
	if err := json.NewDecoder(resp.Body).Decode(delVoiceResp); err != nil {
		return nil, err
	}
	return delVoiceResp, nil
}

func createFilePart(w *multipart.Writer, fieldname, filename, mimeType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", mimeType)
	return w.CreatePart(h)
}

func createFieldPart(w *multipart.Writer, fieldname, mimeType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname)))
	h.Set("Content-Type", mimeType)
	return w.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
