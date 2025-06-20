package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
)

type fakePostsSelector struct {
	Posts []post.Post

	shouldFail bool
}

func (f *fakePostsSelector) SelectPosts(ctx context.Context, limit int) ([]post.Post, error) {
	if f.shouldFail {
		return nil, errors.New("test error")
	}

	if limit > 0 && len(f.Posts) > limit {
		return f.Posts[:limit], nil
	}

	return f.Posts, nil
}

func TestGetPosts(t *testing.T) {
	type args struct {
		selector PostsSelector
		req      *http.Request
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
		wantBody   []post.Post
	}{
		{
			name: "should return all posts with no limit",
			args: args{
				selector: &fakePostsSelector{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
					},
					shouldFail: false,
				},
				req: httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantBody: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
				{ID: "another-post-id", Content: "Another content"},
			},
		},
		{
			name: "should return error if failed to select posts",
			args: args{
				selector: &fakePostsSelector{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
					},
					shouldFail: true,
				},
				req: httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil),
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
			wantBody:   nil,
		},
		{
			name: "should return limited number of posts",
			args: args{
				selector: &fakePostsSelector{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
						{ID: "third-post-id", Content: "Third content"},
					},
					shouldFail: false,
				},
				req: httptest.NewRequest(http.MethodGet, "/api/v1/posts?limit=2", nil),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
			wantBody: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
				{ID: "another-post-id", Content: "Another content"},
			},
		},
		{
			name: "should return error if invalid limit parameter",
			args: args{
				selector: &fakePostsSelector{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
						{ID: "third-post-id", Content: "Third content"},
					},
					shouldFail: false,
				},
				req: httptest.NewRequest(http.MethodGet, "/api/v1/posts?limit=text", nil),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
			wantBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := GetPosts(tt.args.selector)
			handler(w, tt.args.req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetPosts() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// If we expect an error, we just check that response body is not empty
			if tt.wantErr {
				if w.Body.Len() == 0 {
					t.Errorf("GetPosts() response body is empty, want error")
				}
				return
			}

			// Decode the response body into struct for comparison
			var resBody []post.Post
			if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Errorf("GetPosts() error json decoding response body: %v", err)
			}

			if !reflect.DeepEqual(resBody, tt.wantBody) {
				t.Errorf("GetPosts() response body = %v, want %v", resBody, tt.wantBody)
			}
		})
	}
}
