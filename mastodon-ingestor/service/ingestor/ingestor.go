package ingestor

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/event"
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/handler"
)

type Ingestor struct {
	Clients Clients

	events []event.Event
}

type Clients struct {
	Mastodon MastodonClient
	PubSub   PubSubClient
}

type MastodonClient interface {
	SubscribeEvent(ctx context.Context, eventType string, handler func(ctx context.Context, data []byte))
	Listen() error
}

type PubSubClient interface {
	handler.PostSavePublisher
	handler.PostDeletePublisher
}

func New(clients Clients) *Ingestor {
	var i Ingestor

	i.Clients = clients

	// Setup event handlers
	i.setupEvents()

	return &i
}

func (i *Ingestor) Start(ctx context.Context, errc chan<- error) {
	for _, e := range i.events {
		i.Clients.Mastodon.SubscribeEvent(ctx, e.Type, func(ctx context.Context, data []byte) {
			err := e.Handler(ctx, data)
			if err != nil {
				errc <- fmt.Errorf("error handling event %s: %v", e.Type, err)
			}
		})
	}

	log.Printf("Ingestor is listening to %d events", len(i.events))

	err := i.Clients.Mastodon.Listen()
	if !errors.Is(err, context.Canceled) {
		errc <- fmt.Errorf("error listening to Mastodon: %v", err)
	}
}

func (i *Ingestor) Stop(ctx context.Context) error {
	return nil
}

func (i *Ingestor) handle(e event.Event) {
	i.events = append(i.events, e)
}
