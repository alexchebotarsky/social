package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
)

type PostDeletePayload struct {
	ID string `json:"id"`
}

func (ps *Client) PublishPostDelete(ctx context.Context, id string) error {
	payload := PostDeletePayload{
		ID: id,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling post delete payload: %v", err)
	}

	err = ps.Publish(ctx, "social/delete-post", data)
	if err != nil {
		return fmt.Errorf("error publishing post delete: %v", err)
	}

	return nil
}
