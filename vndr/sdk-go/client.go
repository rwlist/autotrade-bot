package chatexsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const accessTokenRenewGap = time.Minute

type Client struct {
	client       *http.Client
	baseURL      string
	refreshToken string

	accessToken           string
	accessTokenExpiration time.Time
	accessTokenLock       sync.Locker
}

func NewClient(baseURL, refreshToken string, client ...*http.Client) *Client {
	c := http.DefaultClient
	if len(client) > 0 {
		c = client[0]
	}

	return &Client{
		client:          c,
		baseURL:         strings.TrimRight(baseURL, "/"),
		refreshToken:    refreshToken,
		accessTokenLock: &sync.Mutex{},
	}
}

// Body should be already added to req
func (c *Client) sendRequest(ctx context.Context, req *http.Request, result interface{}, overrideToken ...string) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	var token string
	if len(overrideToken) > 0 {
		token = overrideToken[0]
	} else {
		accessToken, err := c.getAccessToken(ctx)
		if err != nil {
			return err
		}

		token = accessToken
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return c.parseError(res)
	}

	body, _ := ioutil.ReadAll(res.Body)

	logrus.WithFields(logrus.Fields{
		"status_code": res.StatusCode,
		"url":         req.URL.String(),
	}).Info(string(body))

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	return nil
}

func (c *Client) parseError(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusUnprocessableEntity:
		return ErrUnprocessableEntity
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusBadRequest:
		body, _ := ioutil.ReadAll(res.Body)
		var errors map[string]interface{}

		if err := json.Unmarshal(body, &errors); err != nil {
			return NewValidationError(err.Error(), nil)
		}

		return NewValidationError("validation error", errors)
	}
	return ErrInternalServer
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	c.accessTokenLock.Lock()
	defer c.accessTokenLock.Unlock()

	now := time.Now().Add(accessTokenRenewGap)
	if c.accessTokenExpiration.After(now) {
		return c.accessToken, nil
	}

	token, err := c.CreateAccessToken(ctx)
	if err != nil {
		return "", err
	}

	c.accessToken = token.Token
	c.accessTokenExpiration = time.Unix(token.ExpiresAt, 0)

	return c.accessToken, nil
}
