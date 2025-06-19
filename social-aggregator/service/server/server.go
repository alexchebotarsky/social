package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Host   string
	Port   uint16
	Router *http.ServeMux
	HTTP   *http.Server
}

func New(host string, port uint16) *Server {
	var s Server

	s.Host = host
	s.Port = port
	s.Router = http.NewServeMux()
	s.HTTP = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler:      s.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	s.setupRoutes()

	return &s
}

func (s *Server) Start(ctx context.Context, errc chan<- error) {
	log.Printf("Server is listening at %s:%d", s.Host, s.Port)
	err := s.HTTP.ListenAndServe()
	if err != http.ErrServerClosed {
		errc <- fmt.Errorf("Error listening and serving: %v", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.HTTP.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down http server: %v", err)
	}

	return nil
}
