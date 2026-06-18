package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rcli/feedback/internal/feedback"
)

type APIClient struct {
	baseURL string
	apiKey  string
	app     string
	http    *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL: envOr("FEEDBACK_API_URL", "http://localhost:8080"),
		apiKey:  os.Getenv("FEEDBACK_API_KEY"),
		app:     os.Getenv("FEEDBACK_APP"),
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *APIClient) do(method, path string, body any, result any) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	if c.app != "" {
		req.Header.Set("X-App", c.app)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("%s", errResp.Error)
		}
		return fmt.Errorf("request failed (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, result)
	}
	return nil
}

func (c *APIClient) Submit(req feedback.SubmitRequest) (*feedback.Ticket, error) {
	var ticket feedback.Ticket
	if err := c.do(http.MethodPost, "/api/v1/feedback", req, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (c *APIClient) List(filter feedback.ListFilter) ([]feedback.Ticket, error) {
	path := fmt.Sprintf("/api/v1/feedback?app=%s&status=%s&state=%s",
		filter.App, filter.Status, filter.State)
	if filter.Type != "" {
		path += "&type=" + filter.Type
	}

	var tickets []feedback.Ticket
	if err := c.do(http.MethodGet, path, nil, &tickets); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (c *APIClient) Get(number int) (*feedback.Ticket, error) {
	var ticket feedback.Ticket
	path := fmt.Sprintf("/api/v1/feedback/%d", number)
	if err := c.do(http.MethodGet, path, nil, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (c *APIClient) Comment(number int, body string) error {
	path := fmt.Sprintf("/api/v1/feedback/%d/comments", number)
	return c.do(http.MethodPost, path, map[string]string{"body": body}, nil)
}

func (c *APIClient) UpdateStatus(number int, status feedback.Status) (*feedback.Ticket, error) {
	var ticket feedback.Ticket
	path := fmt.Sprintf("/api/v1/feedback/%d", number)
	if err := c.do(http.MethodPatch, path, map[string]string{"status": string(status)}, &ticket); err != nil {
		return nil, err
	}
	return &ticket, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}