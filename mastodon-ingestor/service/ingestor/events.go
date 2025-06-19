package ingestor

import (
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/event"
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor/handler"
)

func (i *Ingestor) setupEvents() {
	i.handle(event.Event{
		Name:    "PostSave",
		Type:    "update",
		Handler: handler.PostSave(i.Clients.PubSub),
	})
}
