package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alexchebotarsky/social/social-media-aggregator/client"
	"github.com/alexchebotarsky/social/social-media-aggregator/client/database"
	"github.com/alexchebotarsky/social/social-media-aggregator/client/poststream"
	"github.com/alexchebotarsky/social/social-media-aggregator/client/pubsub"
	"github.com/alexchebotarsky/social/social-media-aggregator/env"
	"github.com/alexchebotarsky/social/social-media-aggregator/service/processor"
	"github.com/alexchebotarsky/social/social-media-aggregator/service/server"
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
	var errors []error

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, service := range app.Services {
		err := service.Stop(ctx)
		if err != nil {
			errors = append(errors, fmt.Errorf("error stopping a service: %v", err))
		}
	}

	// Shut down all clients gracefully
	err := app.Clients.Close(ctx)
	if err != nil {
		errors = append(errors, fmt.Errorf("error closing app clients: %v", err))
	}

	if len(errors) > 0 {
		log.Printf("Error gracefully shutting down: %v", &client.ErrMultiple{Errs: errors})
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

	// Main HTTP server for communicating with the app
	server := server.New(env.Host, env.Port, server.Clients{
		Database:   clients.Database,
		PostStream: clients.PostStream,
	})
	services = append(services, server)

	// Processor that handles incoming messages from ingestors
	processor := processor.New(processor.Clients{
		PubSub:     clients.PubSub,
		Database:   clients.Database,
		PostStream: clients.PostStream,
	})
	services = append(services, processor)

	return services, nil
}

// Clients holds implementations of all external clients used in the app
type Clients struct {
	PubSub     *pubsub.Client
	Database   *database.Client
	PostStream *poststream.Client
}

func setupClients(ctx context.Context, env *env.Config) (*Clients, error) {
	var c Clients
	var err error

	c.PubSub, err = pubsub.New(ctx, env.PubSubHost, env.PubSubPort, env.PubSubClientID, env.PubSubQoS)
	if err != nil {
		return nil, fmt.Errorf("error creating PubSub client: %v", err)
	}

	c.Database, err = database.New(ctx, env.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %v", err)
	}

	c.PostStream = poststream.New()

	return &c, nil
}

func (c *Clients) Close(ctx context.Context) error {
	var errors []error

	err := c.PostStream.Close(ctx)
	if err != nil {
		errors = append(errors, fmt.Errorf("error closing post stream client: %v", err))
	}

	if len(errors) > 0 {
		return &client.ErrMultiple{Errs: errors}
	}

	return nil
}
