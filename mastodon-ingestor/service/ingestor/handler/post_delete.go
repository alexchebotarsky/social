package handler

import (
	"context"
	"log"

	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/event"
)

type PostDeletePublisher interface {
	PublishPostDelete(ctx context.Context, postID string) error
}

func PostDelete(publisher PostDeletePublisher) event.Handler {
	return func(ctx context.Context, data []byte) {
		postID := string(data)

		err := publisher.PublishPostDelete(ctx, postID)
		if err != nil {
			log.Printf("Error publishing post for delete: %v", err)
			return
		}
	}
}
