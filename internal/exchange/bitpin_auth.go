package exchange

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (b *BitpinExchange) AuthenticateBitpin(ctx context.Context) (*TokenResponse, error, int) {
	urlToken := "https://api.bitpin.ir/api/v1/usr/authenticate/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]string{
		"api_key":    b.apiKey,
		"secret_key": b.secretKey,
	}

	body, status, err := b.client.PostJSON(ctx, urlToken, data, headers)
	if err != nil {
		b.logger.Error("failed to authenticate", zap.Error(err))
		return nil, err, status
	}
	if status < 200 || status >= 300 {
		b.logger.Error("authentication failed", zap.Int("status", status), zap.ByteString("body", body))
		return nil, fmt.Errorf("authentication failed with status %d", status), status

	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		b.logger.Error("failed to unmarshal auth response", zap.Error(err))
		return nil, err, status
	}
	if tokenResp.Access == "" {
		return nil, fmt.Errorf("empty token received from auth"), status
	}
	return &tokenResp, nil, status
}
