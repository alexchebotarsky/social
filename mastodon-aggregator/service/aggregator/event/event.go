package event

import (
	"context"
)

type Event struct {
	Name    string
	Type    string
	Handler Handler
}

type Handler = func(context.Context, []byte)
