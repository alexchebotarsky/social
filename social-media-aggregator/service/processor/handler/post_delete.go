package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/social/social-media-aggregator/client"
	"github.com/alexchebotarsky/social/social-media-aggregator/service/processor/event"
)

type PostDeletePayload struct {
	ID string `json:"id"`
}

type PostsDeleter interface {
	DeletePost(ctx context.Context, postID string) error
}

func PostDelete(deleter PostsDeleter) event.Handler {
	return func(ctx context.Context, payload []byte) error {
		var postDeletePayload PostDeletePayload
		err := json.Unmarshal(payload, &postDeletePayload)
		if err != nil {
			return fmt.Errorf("error unmarshalling post delete payload: %v", err)
		}

		err = deleter.DeletePost(ctx, postDeletePayload.ID)
		if err != nil {
			switch err.(type) {
			case *client.ErrNotFound:
				// If the post is not found, we can ignore the error
			default:
				return fmt.Errorf("error deleting post: %v", err)
			}
		}

		return nil
	}
}
