package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CAT5NEKO/hijikiTool/internal/application/ports"
)

type MisskeyPoster struct {
	host       string
	token      string
	visibility string
	localOnly  bool
	httpClient *http.Client
}

type misskeyPostRequest struct {
	I          string `json:"i"`
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
	LocalOnly  bool   `json:"localOnly,omitempty"`
}

func NewMisskeyPoster(config ports.Config) ports.Poster {
	return &MisskeyPoster{
		host:       config.MisskeyHost,
		token:      config.MisskeyToken,
		visibility: config.Visibility,
		localOnly:  config.LocalOnly,
		httpClient: &http.Client{},
	}
}

func (p *MisskeyPoster) Post(content string) error {
	request := misskeyPostRequest{
		I:          p.token,
		Text:       content,
		Visibility: p.visibility,
		LocalOnly:  p.localOnly,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://%s/api/notes/create", p.host)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return nil
}
