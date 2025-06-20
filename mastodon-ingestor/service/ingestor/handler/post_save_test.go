package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alexchebotarsky/social/mastodon-ingestor/model/post"
)

type fakePostSavePublisher struct {
	PublishedPosts []post.Post

	shouldFail bool
}

func (f *fakePostSavePublisher) PublishPostSave(ctx context.Context, post *post.Post) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	if post != nil {
		f.PublishedPosts = append(f.PublishedPosts, *post)
	}

	return nil
}

func TestPostSave(t *testing.T) {
	type args struct {
		publisher *fakePostSavePublisher
		data      []byte
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
		wantPosts  []post.Post
	}{
		{
			name: "should publish post for save",
			args: args{
				publisher: &fakePostSavePublisher{
					PublishedPosts: []post.Post{},
					shouldFail:     false,
				},
				data: []byte(`{"id":"test-post-id","content":"Test content"}`),
			},
			wantErr: false,
			wantPosts: []post.Post{
				{
					ID:      "test-post-id",
					Content: "Test content",
				},
			},
		},
		{
			name: "should return error if failed publishing post for save",
			args: args{
				publisher: &fakePostSavePublisher{
					PublishedPosts: []post.Post{},
					shouldFail:     true,
				},
				data: []byte(`{"id":"test-post-id","content":"Test content"}`),
			},
			wantErr:   true,
			wantPosts: []post.Post{},
		},
		{
			name: "should return error if data is invalid",
			args: args{
				publisher: &fakePostSavePublisher{
					PublishedPosts: []post.Post{},
					shouldFail:     false,
				},
				data: []byte(`{"id":"test-post-id","content"}`),
			},
			wantErr:   true,
			wantPosts: []post.Post{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PostSave(tt.args.publisher)
			err := handler(context.Background(), tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostSave() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.publisher.PublishedPosts, tt.wantPosts) {
				t.Errorf("PostSave() published posts = %v, want %v", tt.args.publisher.PublishedPosts, tt.wantPosts)
			}
		})
	}
}
