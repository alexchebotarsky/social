package processor

import (
	"github.com/alexchebotarsky/social/social-aggregator/service/processor/event"
	"github.com/alexchebotarsky/social/social-aggregator/service/processor/handler"
)

func (p *Processor) setupEvents() {
	p.handle(event.Event{
		Topic:   "social/save-post",
		Handler: handler.PostSave(p.Clients.Database),
	})
}
