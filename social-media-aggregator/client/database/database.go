package database

import (
	"context"
	"fmt"

	"github.com/alexchebotarsky/social/social-media-aggregator/client"
	"github.com/jmoiron/sqlx"

	// sqlite driver
	_ "modernc.org/sqlite"
)

type Client struct {
	DB *sqlx.DB
}

func New(ctx context.Context, path string) (*Client, error) {
	var c Client
	var err error

	c.DB, err = sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = c.initSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing posts table: %v", err)
	}

	return &c, nil
}

func (c *Client) initSchema(ctx context.Context) error {
	schema := `
		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			created_at TEXT NOT NULL,
			url TEXT NOT NULL,
			language TEXT,
			content TEXT
		);
	`

	_, err := c.DB.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("error initializing database schema: %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	errs := []error{}

	err := c.DB.Close()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return &client.ErrMultiple{Errs: errs}
	}

	return nil
}
