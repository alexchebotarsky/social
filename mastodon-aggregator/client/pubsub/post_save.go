package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexchebotarsky/social/mastodon-aggregator/model/post"
)

func (ps *Client) PublishPostSave(ctx context.Context, post *post.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("error marshalling post: %v", err)
	}

	err = ps.Publish(ctx, "social/save-post", data)
	if err != nil {
		return fmt.Errorf("error publishing post save: %v", err)
	}

	return nil
}
