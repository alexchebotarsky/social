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

type mockPostPublisher struct {
	PublishCounter int

	shouldFail bool
}

func (f *mockPostPublisher) PublishPost(post *post.Post) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	f.PublishCounter += 1

	return nil
}

func TestPostSave(t *testing.T) {
	type args struct {
		inserter  *fakePostInserter
		publisher *mockPostPublisher
		payload   []byte
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		wantPublishCounter int
		wantPosts          []post.Post
	}{
		{
			name: "should save post",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{
							ID:        "test-post-id",
							Content:   "Test content",
							CreatedAt: "2025-06-20T07:33:19.000Z",
							URL:       "https://mastodon.social/@UserName/114714528936167563",
							Language:  "en",
						},
					},
					shouldFail: false,
				},
				publisher: &mockPostPublisher{
					shouldFail: false,
				},
				payload: []byte(`{
					"id": "another-post-id",
					"content": "Another content",
					"created_at": "2025-06-20T07:33:19.000Z",
					"url": "https://mastodon.social/@UserName/114714528936167563",
					"language": "en"
				}`),
			},
			wantErr:            false,
			wantPublishCounter: 1,
			wantPosts: []post.Post{
				{
					ID:        "test-post-id",
					Content:   "Test content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
				{
					ID:        "another-post-id",
					Content:   "Another content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
			},
		},
		{
			name: "should return error if failed to save post",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{
							ID:        "test-post-id",
							Content:   "Test content",
							CreatedAt: "2025-06-20T07:33:19.000Z",
							URL:       "https://mastodon.social/@UserName/114714528936167563",
							Language:  "en",
						},
					},
					shouldFail: true,
				},
				publisher: &mockPostPublisher{
					shouldFail: false,
				},
				payload: []byte(`{
					"id": "another-post-id",
					"content": "Another content",
					"created_at": "2025-06-20T07:33:19.000Z",
					"url": "https://mastodon.social/@UserName/114714528936167563",
					"language": "en"
				}`),
			},
			wantErr:            true,
			wantPublishCounter: 0,
			wantPosts: []post.Post{
				{
					ID:        "test-post-id",
					Content:   "Test content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
			},
		},
		{
			name: "should return error if failed to publish post",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{
							ID:        "test-post-id",
							Content:   "Test content",
							CreatedAt: "2025-06-20T07:33:19.000Z",
							URL:       "https://mastodon.social/@UserName/114714528936167563",
							Language:  "en",
						},
					},
					shouldFail: false,
				},
				publisher: &mockPostPublisher{
					shouldFail: true,
				},
				payload: []byte(`{
					"id": "another-post-id",
					"content": "Another content",
					"created_at": "2025-06-20T07:33:19.000Z",
					"url": "https://mastodon.social/@UserName/114714528936167563",
					"language": "en"
				}`),
			},
			wantErr:            true,
			wantPublishCounter: 0,
			wantPosts: []post.Post{
				{
					ID:        "test-post-id",
					Content:   "Test content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
				{
					ID:        "another-post-id",
					Content:   "Another content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
			},
		},
		{
			name: "should return error if payload is invalid",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{
							ID:        "test-post-id",
							Content:   "Test content",
							CreatedAt: "2025-06-20T07:33:19.000Z",
							URL:       "https://mastodon.social/@UserName/114714528936167563",
							Language:  "en",
						},
					},
					shouldFail: false,
				},
				publisher: &mockPostPublisher{
					shouldFail: false,
				},
				payload: []byte(`{
					"id": "another-post-id",
					"content": "Another content",
					"created_at": "2025-06-20T07:33:19.000Z",
					"url": "https://mastodon.social/@UserName/114714528936167563",
					"language":
				}`),
			},
			wantErr:            true,
			wantPublishCounter: 0,
			wantPosts: []post.Post{
				{
					ID:        "test-post-id",
					Content:   "Test content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
			},
		},
		{
			name: "should return error if post is invalid",
			args: args{
				inserter: &fakePostInserter{
					Posts: []post.Post{
						{
							ID:        "test-post-id",
							Content:   "Test content",
							CreatedAt: "2025-06-20T07:33:19.000Z",
							URL:       "https://mastodon.social/@UserName/114714528936167563",
							Language:  "en",
						},
					},
					shouldFail: false,
				},
				publisher: &mockPostPublisher{
					shouldFail: false,
				},
				payload: []byte(`{
					"id": "",
					"content": "",
					"created_at": "",
					"url": "",
					"language": ""
				}`),
			},
			wantErr:            true,
			wantPublishCounter: 0,
			wantPosts: []post.Post{
				{
					ID:        "test-post-id",
					Content:   "Test content",
					CreatedAt: "2025-06-20T07:33:19.000Z",
					URL:       "https://mastodon.social/@UserName/114714528936167563",
					Language:  "en",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PostSave(tt.args.inserter, tt.args.publisher)
			err := handler(context.Background(), tt.args.payload)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostSave() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.inserter.Posts, tt.wantPosts) {
				t.Errorf("PostSave() posts = %v, want %v", tt.args.inserter.Posts, tt.wantPosts)
			}

			if tt.args.publisher.PublishCounter != tt.wantPublishCounter {
				t.Errorf("PostSave() publish counter = %d, want %d", tt.args.publisher.PublishCounter, tt.wantPublishCounter)
			}
		})
	}
}
