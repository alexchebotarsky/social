package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alexchebotarsky/social/mastodon-aggregator/aggregator/event"
	"github.com/alexchebotarsky/social/mastodon-aggregator/post"
)

func PostSave() event.Handler {
	return func(ctx context.Context, data []byte) {
		var post post.Post
		err := json.Unmarshal(data, &post)
		if err != nil {
			log.Printf("Error unmarshalling post data: %v", err)
			return
		}

		log.Printf("Post saved: %s", post.URL)
	}
}
