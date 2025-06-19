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

// InsertPost inserts a new post into the database or updates it if it already
// exists.
func (c *Client) InsertPost(ctx context.Context, p *post.Post) error {
	query := `
		INSERT INTO posts (id, created_at, url, language, content)
		VALUES (:id, :created_at, :url, :language, :content)
		ON CONFLICT(id) DO UPDATE SET
			created_at = excluded.created_at,
			url = excluded.url,
			language = excluded.language,
			content = excluded.content
	`

	_, err := c.DB.NamedExecContext(ctx, query, p)
	if err != nil {
		return fmt.Errorf("error inserting post: %v", err)
	}

	return nil
}
