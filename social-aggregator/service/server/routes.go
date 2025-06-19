package server

import (
	"fmt"

	"github.com/alexchebotarsky/social/social-aggregator/service/server/handler"
)

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("GET /_healthz", handler.Health)
	s.Router.HandleFunc(fmt.Sprintf("GET %s/posts", v1API), handler.GetPosts(s.Clients.Database))
}

const v1API = "/api/v1"
