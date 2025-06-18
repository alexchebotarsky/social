package aggregator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/alexchebotarsky/social/mastodon-aggregator/aggregator/event"
)

type Aggregator struct {
	Clients Clients

	events []event.Event
}

type Clients struct {
	Mastodon MastodonClient
}

type MastodonClient interface {
	SubscribeEvent(ctx context.Context, eventType string, handler event.Handler)
	Listen() error
}

func New(clients Clients) (*Aggregator, error) {
	var a Aggregator

	a.Clients = clients

	// Setup event handlers
	a.setupEvents()

	return &a, nil
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
