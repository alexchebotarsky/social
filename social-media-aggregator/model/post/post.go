package post

import (
	"errors"
	"fmt"
)

type Post struct {
	ID        string `json:"id" db:"id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	URL       string `json:"url" db:"url"`
	Language  string `json:"language" db:"language"`
	Content   string `json:"content" db:"content"`
}

func (p *Post) Validate() error {
	if p.ID == "" {
		return errors.New("missing required field: id")
	}
	if p.CreatedAt == "" {
		return errors.New("missing required field: created_at")
	}
	if p.URL == "" {
		return fmt.Errorf("missing required field: url")
	}
	if p.Language == "" {
		return fmt.Errorf("missing required field: language")
	}
	if p.Content == "" {
		return fmt.Errorf("missing required field: content")
	}
	return nil
}
