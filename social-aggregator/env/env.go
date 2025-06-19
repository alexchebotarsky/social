package env

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host string `env:"HOST,default=0.0.0.0"`
	Port uint16 `env:"PORT,default=8000"`

	PubSubHost     string `env:"PUBSUB_HOST,default=localhost"`
	PubSubPort     uint16 `env:"PUBSUB_PORT,default=1883"`
	PubSubClientID string `env:"PUBSUB_CLIENT_ID,default=social-aggregator"`
	PubSubQoS      byte   `env:"PUBSUB_QOS,default=1"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var c Config

	// We are loading env variables from .env file only for local development
	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file for local development: %v", err)
	}

	err = envconfig.Process(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("error processing environment variables: %v", err)
	}

	return &c, nil
}
