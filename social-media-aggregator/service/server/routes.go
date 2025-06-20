package server

import (
	"fmt"

	"github.com/alexchebotarsky/social/social-media-aggregator/service/server/handler"
)

func (s *Server) setupRoutes() {
	s.Router.HandleFunc("GET /_healthz", handler.Health)
	s.Router.HandleFunc("/openapi.yaml", handler.OpenAPIYaml)
	s.Router.HandleFunc("/docs", handler.SwaggerUI)

	s.Router.HandleFunc(fmt.Sprintf("GET %s/posts", v1API), handler.GetPosts(s.Clients.Database))
	s.Router.Handle(fmt.Sprintf("GET %s/posts/stream", v1API), s.Clients.PostStream.Handler())
}

const v1API = "/api/v1"
