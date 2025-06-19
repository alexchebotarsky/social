package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
	"github.com/alexchebotarsky/social/social-media-aggregator/service/processor/event"
)

type PostsInserter interface {
	InsertPost(ctx context.Context, post *post.Post) error
}

func PostSave(inserter PostsInserter) event.Handler {
	return func(ctx context.Context, payload []byte) error {
		var post post.Post
		err := json.Unmarshal(payload, &post)
		if err != nil {
			return fmt.Errorf("error unmarshalling post: %v", err)
		}

		err = inserter.InsertPost(ctx, &post)
		if err != nil {
			return fmt.Errorf("error inserting post: %v", err)
		}

		return nil
	}
}
