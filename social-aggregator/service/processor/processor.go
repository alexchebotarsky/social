package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/alexchebotarsky/social/social-aggregator/service/processor/event"
)

type Processor struct {
	Events  []event.Event
	Clients Clients
}

type Clients struct {
	PubSub PubSubClient
}

type PubSubClient interface {
	Subscribe(ctx context.Context, topic string, handler event.Handler) error
}

func New(clients Clients) *Processor {
	var p Processor

	p.Clients = clients

	p.setupEvents()

	return &p
}

func (p *Processor) Start(ctx context.Context, errc chan<- error) {
	for _, e := range p.Events {
		err := p.Clients.PubSub.Subscribe(ctx, e.Topic, e.Handler)
		if err != nil {
			errc <- fmt.Errorf("error subscribing to topic %s: %v", e.Topic, err)
			return
		}
	}

	log.Printf("PubSub event processor listening to %d events", len(p.Events))
}

func (p *Processor) Stop(ctx context.Context) error {
	return nil
}

func (p *Processor) handle(e event.Event) {
	p.Events = append(p.Events, e)
}
