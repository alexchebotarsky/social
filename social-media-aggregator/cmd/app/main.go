package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/alexchebotarsky/social/social-media-aggregator/app"
	"github.com/alexchebotarsky/social/social-media-aggregator/env"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	env, err := env.LoadConfig(ctx)
	if err != nil {
		log.Printf("Error loading env config: %v", err)
		os.Exit(1)
	}

	app, err := app.New(ctx, env)
	if err != nil {
		log.Printf("Error creating app: %v", err)
		os.Exit(1)
	}

	app.Launch(ctx)
}
