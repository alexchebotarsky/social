package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/alexchebotarsky/social/social-aggregator/client"
	"github.com/alexchebotarsky/social/social-aggregator/client/pubsub"
	"github.com/alexchebotarsky/social/social-aggregator/env"
	"github.com/alexchebotarsky/social/social-aggregator/service/processor"
	"github.com/alexchebotarsky/social/social-aggregator/service/server"
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
	err := app.Clients.Close()
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
	server, err := server.New(env.Host, env.Port)
	if err != nil {
		return nil, fmt.Errorf("error creating server: %v", err)
	}
	services = append(services, server)

	// Processor that handles incoming messages from ingestors
	processor, err := processor.New(processor.Clients{PubSub: clients.PubSub})
	if err != nil {
		return nil, fmt.Errorf("error creating processor: %v", err)
	}
	services = append(services, processor)

	return services, nil
}

// Clients holds implementations of all external clients used in the app
type Clients struct {
	PubSub *pubsub.PubSub
}

func setupClients(ctx context.Context, env *env.Config) (*Clients, error) {
	var c Clients
	var err error

	c.PubSub, err = pubsub.New(ctx, env.PubSubHost, env.PubSubPort, env.PubSubClientID, env.PubSubQoS)
	if err != nil {
		return nil, fmt.Errorf("error creating PubSub client: %v", err)
	}

	return &c, nil
}

func (c *Clients) Close() error {
	var errors []error

	if len(errors) > 0 {
		return &client.ErrMultiple{Errs: errors}
	}

	return nil
}
