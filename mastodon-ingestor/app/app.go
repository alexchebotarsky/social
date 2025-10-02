package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alexchebotarsky/social/mastodon-ingestor/client/mastodon"
	"github.com/alexchebotarsky/social/mastodon-ingestor/client/pubsub"
	"github.com/alexchebotarsky/social/mastodon-ingestor/env"
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/ingestor"
	"github.com/alexchebotarsky/social/mastodon-ingestor/service/server"
)

type App struct {
	Services []Service
	Clients  *Clients
}

func New(ctx context.Context, env *env.Config) (*App, error) {
	var app App
	var err error

	app.Clients, err = setupClients(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("error setting up clients: %v", err)
	}

	app.Services, err = setupServices(ctx, env, app.Clients)
	if err != nil {
		return nil, fmt.Errorf("error setting up services: %v", err)
	}

	return &app, nil
}

func (app *App) Launch(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	errc := make(chan error, 1)

	// Start all services in separate goroutines
	for _, service := range app.Services {
		go service.Start(ctx, errc)
	}

	// Wait for context cancellation or a critical error from any service
	select {
	case <-ctx.Done():
		log.Print("Context is cancelled")
	case err := <-errc:
		log.Printf("Critical service error: %v", err)
	}

	// Shut down all services gracefully
	var errs []error

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, service := range app.Services {
		err := service.Stop(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("error stopping a service: %v", err))
		}
	}

	// Shut down all clients gracefully
	err := app.Clients.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing app clients: %v", err))
	}

	if len(errs) > 0 {
		log.Printf("Error gracefully shutting down: %v", errors.Join(errs...))
	} else {
		log.Print("App has been gracefully shut down")
	}
}

// Service is an interface for main actors of the app, it runs in a separate goroutine
type Service interface {
	Start(context.Context, chan<- error)
	Stop(context.Context) error
}

func setupServices(ctx context.Context, env *env.Config, clients *Clients) ([]Service, error) {
	var services []Service

	// Internal server for health checks and other utility endpoints
	server := server.New(env.Host, env.Port)
	services = append(services, server)

	// Ingestor service that collects and processes data
	ingestor := ingestor.New(ingestor.Clients{
		Mastodon: clients.Mastodon,
		PubSub:   clients.PubSub,
	})
	services = append(services, ingestor)

	return services, nil
}

// Clients holds implementations of all external clients used in the app
type Clients struct {
	Mastodon *mastodon.Client
	PubSub   *pubsub.Client
}

func setupClients(ctx context.Context, env *env.Config) (*Clients, error) {
	var c Clients
	var err error

	c.Mastodon, err = mastodon.New(env.MastodonStreamingURL, env.MastodonAccessToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Mastodon client: %v", err)
	}

	c.PubSub, err = pubsub.New(ctx, env.PubSubHost, env.PubSubPort, env.PubSubClientID, env.PubSubQoS)
	if err != nil {
		return nil, fmt.Errorf("error creating PubSub client: %v", err)
	}

	return &c, nil
}

func (c *Clients) Close() error {
	var errs []error

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
