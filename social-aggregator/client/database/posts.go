package database

import (
	"context"
	"fmt"

	"github.com/alexchebotarsky/social/social-aggregator/model/post"
)

func (c *Client) SelectPosts(ctx context.Context) ([]post.Post, error) {
	query := "SELECT id, created_at, url, language, content FROM posts"

	posts := []post.Post{}
	err := c.DB.SelectContext(ctx, &posts, query)
	if err != nil {
		return nil, fmt.Errorf("error selecting posts: %v", err)
	}

	return posts, nil
}
