package mastodon

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/tmaxmax/go-sse"
)

type Client struct {
	URL         string
	AccessToken string

	sseConn *sse.Connection
}

func New(url, accessToken string) (*Client, error) {
	var c Client

	c.URL = url
	c.AccessToken = accessToken

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	c.sseConn = sse.DefaultClient.NewConnection(req)

	return &c, nil
}

func (c *Client) Listen() error {
	err := c.sseConn.Connect()
	if !errors.Is(err, context.Canceled) {
		return fmt.Errorf("error connecting to SSE stream: %v", err)
	}

	return nil
}

func (c *Client) SubscribeEvent(ctx context.Context, eventType string, handler func(context.Context, []byte)) {
	c.sseConn.SubscribeEvent(eventType, func(sseEvent sse.Event) {
		handler(ctx, []byte(sseEvent.Data))
	})
}
