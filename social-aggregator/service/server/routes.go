package server

import (
	"github.com/alexchebotarsky/social/social-aggregator/service/server/handler"
)

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("GET /_healthz", handler.Health)
}
