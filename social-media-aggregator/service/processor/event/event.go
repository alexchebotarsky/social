package event

import "context"

type Event struct {
	Topic   string
	Handler Handler
}

type Handler = func(ctx context.Context, payload []byte) error
