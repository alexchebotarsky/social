package aggregator

import (
	"github.com/alexchebotarsky/social/mastodon-aggregator/service/aggregator/event"
	"github.com/alexchebotarsky/social/mastodon-aggregator/service/aggregator/handler"
)

func (a *Aggregator) setupEvents() {
	a.handle(event.Event{
		Name:    "PostSave",
		Type:    "update",
		Handler: handler.PostSave(a.Clients.PubSub),
	})
}
