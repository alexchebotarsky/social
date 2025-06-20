package poststream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
	"github.com/tmaxmax/go-sse"
)

type Client struct {
	server *sse.Server
}

func New() *Client {
	var c Client

	c.server = &sse.Server{}

	return &c
}

func (c *Client) PublishPost(post *post.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("error marshalling post: %w", err)
	}

	msg := &sse.Message{}
	msg.AppendData(string(data))
	msg.Type = sse.Type("post")
	msg.ID = sse.ID(post.ID)

	err = c.server.Publish(msg)
	if err != nil {
		return fmt.Errorf("error publishing post: %w", err)
	}

	return nil
}

// Handler for serving the SSE stream
func (c *Client) Handler() http.Handler {
	return c.server
}

func (c *Client) Close(ctx context.Context) error {
	err := c.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down post stream server: %w", err)
	}

	return nil
}
