package aggregator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alexchebotarsky/social/mastodon-aggregator/service/aggregator/event"
	"github.com/alexchebotarsky/social/mastodon-aggregator/service/aggregator/handler"
)

type Aggregator struct {
	Clients Clients

	events []event.Event
}

type Clients struct {
	Mastodon MastodonClient
	PubSub   PubSubClient
}

type MastodonClient interface {
	SubscribeEvent(ctx context.Context, eventType string, handler event.Handler)
	Listen() error
}

type PubSubClient interface {
	handler.PostSavePublisher
}

func New(clients Clients) *Aggregator {
	var a Aggregator

	a.Clients = clients

	// Setup event handlers
	a.setupEvents()

	return &a
}

func (a *Aggregator) Start(ctx context.Context, errc chan<- error) {
	log.Printf("Aggregator service started")

	for _, e := range a.events {
		a.Clients.Mastodon.SubscribeEvent(ctx, e.Type, e.Handler)
	}

	err := a.Clients.Mastodon.Listen()
	if !errors.Is(err, context.Canceled) {
		errc <- fmt.Errorf("error listening to Mastodon: %v", err)
	}
}

func (a *Aggregator) Stop(ctx context.Context) error {
	return nil
}

func (p *Aggregator) handle(e event.Event) {
	p.events = append(p.events, e)
}
