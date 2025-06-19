package main

import (
	"context"
	"log"
	"os"

	"github.com/alexchebotarsky/social/social-media-aggregator/app"
	"github.com/alexchebotarsky/social/social-media-aggregator/env"
)

func main() {
	ctx := context.Background()

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
