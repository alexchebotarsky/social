package handler

import (
	"context"
	"fmt"

	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/event"
)

type PostDeletePublisher interface {
	PublishPostDelete(ctx context.Context, postID string) error
}

func PostDelete(publisher PostDeletePublisher) event.Handler {
	return func(ctx context.Context, data []byte) error {
		postID := string(data)

		err := publisher.PublishPostDelete(ctx, postID)
		if err != nil {
			return fmt.Errorf("error publishing post for delete: %v", err)
		}

		return nil
	}
}
