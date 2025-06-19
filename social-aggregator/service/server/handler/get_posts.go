package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexchebotarsky/social/social-aggregator/model/post"
)

type PostsSelector interface {
	SelectPosts(ctx context.Context) ([]post.Post, error)
}

func GetPosts(selector PostsSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get posts from the database
		posts, err := selector.SelectPosts(ctx)
		if err != nil {
			handleError(ctx, w, fmt.Errorf("error selecting posts"), http.StatusInternalServerError, true)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Printf("Error encoding posts to JSON: %v", err)
		}
	}
}
