package handler

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type fakePostDeletePublisher struct {
	PublishedIDs []string

	shouldFail bool
}

func (f *fakePostDeletePublisher) PublishPostDelete(ctx context.Context, postID string) error {
	if f.shouldFail {
		return errors.New("test error")
	}

	f.PublishedIDs = append(f.PublishedIDs, postID)

	return nil
}

func TestPostDelete(t *testing.T) {
	type args struct {
		publisher *fakePostDeletePublisher
		data      []byte
	}
	tests := []struct {
		name        string
		args        args
		wantStatus  int
		wantErr     bool
		wantPostIDs []string
	}{
		{
			name: "should publish post id for delete",
			args: args{
				publisher: &fakePostDeletePublisher{
					PublishedIDs: []string{},
					shouldFail:   false,
				},
				data: []byte("test-post-id"),
			},
			wantErr:     false,
			wantPostIDs: []string{"test-post-id"},
		},
		{
			name: "should return error if failed publishing post id for delete",
			args: args{
				publisher: &fakePostDeletePublisher{
					PublishedIDs: []string{},
					shouldFail:   true,
				},
				data: []byte("test-post-id"),
			},
			wantErr:     true,
			wantPostIDs: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := PostDelete(tt.args.publisher)
			err := handler(context.Background(), tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("PostDelete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.publisher.PublishedIDs, tt.wantPostIDs) {
				t.Errorf("PostDelete() published IDs = %v, want %v", tt.args.publisher.PublishedIDs, tt.wantPostIDs)
			}
		})
	}
}
