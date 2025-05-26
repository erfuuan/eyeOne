package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	httpClient *http.Client
	logger     *zap.Logger
}

func New(logger *zap.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("failed to create GET request", zap.Error(err))
		return nil, 0, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("GET request failed", zap.Error(err))
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read GET response body", zap.Error(err))
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func (c *Client) PostJSON(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, int, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("failed to marshal body", zap.Error(err))
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.logger.Error("failed to create POST request", zap.Error(err))
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("POST request failed", zap.Error(err))
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("failed to read response body", zap.Error(err))
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}
