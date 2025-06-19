package server

import (
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/server/handler"
)

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("GET /_healthz", handler.Health)
}
