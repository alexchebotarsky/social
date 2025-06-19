package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/alexchebotarsky/social/social-aggregator/model/post"
)

type PostsSelector interface {
	SelectPosts(ctx context.Context, limit int) ([]post.Post, error)
}

func GetPosts(selector PostsSelector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		limit := 0
		if r.URL.Query().Get("limit") != "" {
			var err error
			limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				handleError(ctx, w, fmt.Errorf("invalid `limit` parameter: %v", err), http.StatusBadRequest, false)
				return
			}
		}

		posts, err := selector.SelectPosts(ctx, limit)
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
