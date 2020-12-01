package chatexsdk

import (
	"context"
	"fmt"
	"net/http"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

func (c *Client) CreateAccessToken(ctx context.Context) (*AccessToken, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/access-token", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	var token AccessToken
	if err := c.sendRequest(ctx, req, &token, c.refreshToken); err != nil {
		return nil, err
	}

	return &token, nil
}
