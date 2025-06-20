package processor

import (
	"github.com/alexchebotarsky/social/social-media-aggregator/service/processor/event"
	"github.com/alexchebotarsky/social/social-media-aggregator/service/processor/handler"
)

func (p *Processor) setupEvents() {
	p.handle(event.Event{
		Topic:   "social/save-post",
		Handler: handler.PostSave(p.Clients.Database, p.Clients.PostStream),
	})
	p.handle(event.Event{
		Topic:   "social/delete-post",
		Handler: handler.PostDelete(p.Clients.Database),
	})
}
