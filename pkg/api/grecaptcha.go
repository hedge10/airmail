package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	GOOGLE_SITE_VERIFY = "https://www.google.com"
)

type Client struct {
	baseUri    string
	httpClient *http.Client
}

type Option func(*Client) error

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func BaseUri(uri string) Option {
	return func(c *Client) error {
		c.baseUri = uri
		return nil
	}
}

func CreateClient(opts ...Option) (*Client, error) {
	c := &Client{
		baseUri: GOOGLE_SITE_VERIFY,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	if err := c.parseOptions(opts...); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) ValidateGrecaptcha(secret string, g_captcha_response string, ip string) error {
	url := fmt.Sprintf("%s%s", c.baseUri, "/recaptcha/api/siteverify")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", g_captcha_response)
	q.Add("remoteip", ip)

	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body SiteVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}
	if !body.Success {
		return errors.New("unsuccessful recaptcha verify request")
	}

	return nil
}

func (c *Client) parseOptions(opts ...Option) error {
	// Range over each options function and apply it to our API type to
	// configure it. Options functions are applied in order, with any
	// conflicting options overriding earlier calls.
	for _, option := range opts {
		err := option(c)
		if err != nil {
			return err
		}
	}

	return nil
}
