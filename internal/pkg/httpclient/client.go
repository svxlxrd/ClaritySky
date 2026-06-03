package httpclient

import (
	"claritysky/internal/pkg/retry"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client

	retryCfg retry.Config
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},

		retryCfg: retry.Config{
			MaxRetries: 3,
			BaseDelay:  100 * time.Millisecond,
			MaxDelay:   3 * time.Second,
		},
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var response *http.Response

	err := retry.Do(
		ctx,
		c.retryCfg,
		func(ctx context.Context) error {
			if response != nil {
				response.Body.Close()
				response = nil
			}

			resp, err := c.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("%w: %v", retry.ErrRetryable, err,)
			}

			response = resp

			// Ретраим только 429 (Too Many Requests) и 5xx (Ошибки сервера)
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
				return fmt.Errorf("%w: server returned status %d", retry.ErrRetryable, resp.StatusCode)
			}

			return nil
		},
	)

	return response, err
}
