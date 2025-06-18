package aggregator

import (
	"github.com/alexchebotarsky/social/mastodon-aggregator/aggregator/event"
	"github.com/alexchebotarsky/social/mastodon-aggregator/aggregator/handler"
)

func (a *Aggregator) setupEvents() {
	a.handle(event.Event{
		Name:    "PostSave",
		Type:    "update",
		Handler: handler.PostSave(),
	})
}
