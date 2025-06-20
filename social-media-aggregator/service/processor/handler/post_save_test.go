package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alexchebotarsky/social/social-media-aggregator/model/post"
)

type fakePostInserter struct {
	Posts []post.Post

	shouldFail bool
}

func (f *fakePostInserter) InsertPost(ctx context.Context, post *post.Post) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	if post != nil {
		f.Posts = append(f.Posts, *post)
	}

	return nil
}

func TestPostSave(t *testing.T) {
	type args struct {
		inserter *fakePostInserter
		payload  []byte
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantPosts []post.Post
	}{
		{
			name: "should save post",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
					},
					shouldFail: false,
				},
				payload: []byte(`{"id":"another-post-id","content":"Another content"}`),
			},
			wantErr: false,
			wantPosts: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
				{ID: "another-post-id", Content: "Another content"},
			},
		},
		{
			name: "should return error if failed to save post",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
					},
					shouldFail: true,
				},
				payload: []byte(`{"id":"another-post-id","content":"Another content"}`),
			},
			wantErr: true,
			wantPosts: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
			},
		},
		{
			name: "should return error if payload is invalid",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{ID: "test-post-id", Content: "Test content"},
					},
					shouldFail: false,
				},
				payload: []byte(`{"id":"another-post-id","content":}`),
			},
			wantErr: true,
			wantPosts: []post.Post{
				{ID: "test-post-id", Content: "Test content"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PostSave(tt.args.inserter)
			err := handler(context.Background(), tt.args.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostSave() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.inserter.Posts, tt.wantPosts) {
				t.Errorf("PostSave() posts = %v, want %v", tt.args.inserter.Posts, tt.wantPosts)
			}
		})
	}
}
