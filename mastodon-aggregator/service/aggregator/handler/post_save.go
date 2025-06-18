package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alexchebotarsky/social/mastodon-aggregator/model/post"
	"github.com/alexchebotarsky/social/mastodon-aggregator/service/aggregator/event"
)

type PostSavePublisher interface {
	PublishPostSave(ctx context.Context, post *post.Post) error
}

func PostSave(publisher PostSavePublisher) event.Handler {
	return func(ctx context.Context, data []byte) {
		var post post.Post
		err := json.Unmarshal(data, &post)
		if err != nil {
			log.Printf("Error unmarshalling post data: %v", err)
			return
		}

		err = publisher.PublishPostSave(ctx, &post)
		if err != nil {
			log.Printf("Error publishing post for save: %v", err)
			return
		}
	}
}
