package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type CounterTypeEnum string

const (
	CounterTypeGAUGE   CounterTypeEnum = "GAUGE"
	CounterTypeCOUNTER CounterTypeEnum = "COUNTER"
)

type Metric struct {
	Endpoint    string          `json:"endpoint"`
	Metric      string          `json:"metric"`
	Timestamp   int64           `json:"timestamp"`
	Step        int64           `json:"step"`
	Value       float64         `json:"value"`
	CounterType CounterTypeEnum `json:"counterType"` // GAUGE | COUNTER
	Tags        string          `json:"tags"`
}

type Client struct {
	url        string
	token      string
	httpClient *http.Client
}

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

func New(url string, opts ...Option) *Client {
	c := &Client{
		url:        url,
		httpClient: &http.Client{Timeout: 10 * time.Second}, // 默认 client
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) PushMetrics(metrics []Metric) error {
	bodyBytes, _ := json.Marshal(metrics)

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return errors.WithMessage(err, "n9e client new http request failed")
	}

	req.Header.Set("Content-Type", "application/json")
	if len(c.token) > 0 {
		req.Header.Set("X-User-Token", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.WithMessage(err, "n9e client http do request failed")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("n9e client push metrics failed, status: %d", resp.StatusCode)
	}
	return nil
}
