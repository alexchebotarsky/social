package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
)

type fakePostDeleter struct {
	Posts []post.Post

	shouldFail bool
}

func (f *fakePostDeleter) DeletePost(ctx context.Context, postID string) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	for i, p := range f.Posts {
		if p.ID == postID {
			f.Posts = append(f.Posts[:i], f.Posts[i+1:]...)
			return nil
		}
	}

	return nil
}

func TestPostDelete(t *testing.T) {
	type args struct {
		deleter *fakePostDeleter
		payload []byte
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPosts []post.Post
	}{
		{
			name: "should delete post",
			args: args{
				deleter: &fakePostDeleter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
					},
					shouldFail: false,
				},
				payload: []byte(`{"id":"test-post-id"}`),
			},
			wantErr: false,
			wantPosts: []post.Post{
				{ID: "another-post-id", Content: "Another content"},
			},
		},
		{
			name: "should return error if failed to delete post",
			args: args{
				deleter: &fakePostDeleter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
					},
					shouldFail: true,
				},
				payload: []byte(`{"id":"test-post-id"}`),
			},
			wantErr: true,
			wantPosts: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
				{ID: "another-post-id", Content: "Another content"},
			},
		},
		{
			name: "should return error if payload is invalid",
			args: args{
				deleter: &fakePostDeleter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
						{ID: "another-post-id", Content: "Another content"},
					},
					shouldFail: true,
				},
				payload: []byte(`test-post-id`),
			},
			wantErr: true,
			wantPosts: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
				{ID: "another-post-id", Content: "Another content"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PostDelete(tt.args.deleter)
			err := handler(context.Background(), tt.args.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostDelete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.deleter.Posts, tt.wantPosts) {
				t.Errorf("PostDelete() posts = %v, want %v", tt.args.deleter.Posts, tt.wantPosts)
			}
		})
	}
}
