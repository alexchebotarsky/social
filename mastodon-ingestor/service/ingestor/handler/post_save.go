package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/social/mastodon-ingestor/model/post"
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/event"
)

type PostSavePublisher interface {
	PublishPostSave(ctx context.Context, post *post.Post) error
}

func PostSave(publisher PostSavePublisher) event.Handler {
	return func(ctx context.Context, data []byte) error {
		var post post.Post
		err := json.Unmarshal(data, &post)
		if err != nil {
			return fmt.Errorf("error unmarshalling post data: %v", err)
		}

		err = publisher.PublishPostSave(ctx, &post)
		if err != nil {
			return fmt.Errorf("error publishing post for save: %v", err)
		}

		return nil
	}
}
