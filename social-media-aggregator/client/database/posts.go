package database

import (
	"context"
	"fmt"

	"github.com/alexchebotarsky/social/social-media-aggregator/client"
	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
)

// SelectPosts retrieves posts in chronological order from the database with a
// limit. If limit is 0, it retrieves all posts.
func (c *Client) SelectPosts(ctx context.Context, limit int) ([]post.Post, error) {
	query := "SELECT id, created_at, url, language, content FROM posts ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

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

// DeletePost deletes a post by its ID from the database. If the post does not
// exist, it returns a not found error.
func (c *Client) DeletePost(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = :id`

	args := struct {
		ID string `db:"id"`
	}{
		ID: id,
	}

	result, err := c.DB.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error deleting post with id %s: %v", id, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return &client.ErrNotFound{Err: fmt.Errorf("no post with id %s to delete", id)}
	}

	return nil
}
