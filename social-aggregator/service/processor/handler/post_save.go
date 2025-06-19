package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alexchebotarsky/social/social-aggregator/model/post"
	"github.com/alexchebotarsky/social/social-aggregator/service/processor/event"
)

func PostSave() event.Handler {
	return func(ctx context.Context, payload []byte) error {
		var post post.Post
		err := json.Unmarshal(payload, &post)
		if err != nil {
			return fmt.Errorf("error unmarshalling post: %v", err)
		}

		log.Printf("Received post: %s", post.URL)

		return nil
	}
}
